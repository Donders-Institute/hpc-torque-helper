package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"dccn.nl/netutil"
	log "github.com/sirupsen/logrus"
)

var (
	connHost    *string
	connPort    *int
	maxConn     *int
	role        *string
	tlsCert     *string
	tlsKey      *string
	mdir        *string
	tdir        *string
	trqServer   *string
	optsVerbose *bool
)

func init() {

	// Command-line arguments
	connHost = flag.String("h", "0.0.0.0", "set the ip `address` of the server")
	connPort = flag.Int("p", 60209, "set the port `number` of the server")
	maxConn = flag.Int("n", 100, "set the max. client connections")
	mdir = flag.String("m", os.Getenv("MOABHOMEDIR"), "set the `path` of Moab installation, usually referred by $MOABHOMEDIR")
	tdir = flag.String("t", os.Getenv("TORQUEHOME"), "set the `path` of the Torque installation")
	role = flag.String("r", os.Getenv("TRQHELPD_ROLE"), "set the `role` of the trqhelpd service. \"SRV\" for trqhelpd running on pbs_server node; or \"MOM\" for running on pbs_mom node.")
	trqServer = flag.String("s", os.Getenv("TORQUESERVER"), "set the `hostname` of the Torque server.  It is used to construct the job's FQID.")
	tlsCert = flag.String("c", os.Getenv("TLS_CERT"), "set the `path` of the TLS certificate")
	tlsKey = flag.String("k", os.Getenv("TLS_KEY"), "set the `path` of the TLS private key")
	optsVerbose = flag.Bool("v", false, "print debug messages")
	flag.Usage = usage

	flag.Parse()

	// set logging
	log.SetOutput(os.Stderr)
	// set logging level
	llevel := log.InfoLevel
	if *optsVerbose {
		llevel = log.DebugLevel
	}
	log.SetLevel(llevel)
	// set logging
	log.SetOutput(os.Stderr)
}

func usage() {
	fmt.Printf("\nHelper service for retriving job information with leveraged privilege.\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {

	// Update PATH environment variable with paths of moab/torque executables.
	os.Setenv("PATH", fmt.Sprintf("%s/bin:%s/bin:%s", *tdir, *mdir, os.Getenv("PATH")))

	// Update TORQUEHOME and MOABHOMEDIR environment variables with the values set to this program.
	if *tdir != "" {
		os.Setenv("TORQUEHOME", *tdir)
	}
	if *mdir != "" {
		os.Setenv("MOABHOMEDIR", *mdir)
	}

	// Select command mapper based the server role.
	var cmdMapper CommandMapper
	switch *role {
	case "SRV":
		cmdMapper = CommandMapperSrv{}
	case "MOM":
		cmdMapper = CommandMapperMom{}
	default:
		log.Error("Unknown service role:", *role)
		os.Exit(1)
	}

	// Load server certificate
	cert, err := tls.LoadX509KeyPair(*tlsCert, *tlsKey)

	if err != nil {
		log.Error("Cannot load certificate:", err.Error())
		os.Exit(1)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader

	// Listen for incoming connections.
	l, err := tls.Listen("tcp", fmt.Sprintf("%s:%d", *connHost, *connPort), &config)
	if err != nil {
		log.Error("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	// Use LimitListener to restrict number of activate connections
	l = netutil.LimitListener(l, *maxConn)

	log.Infof("Listening on %s:%d\n", *connHost, *connPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Error("Error accepting:", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, cmdMapper)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, cmdMapper CommandMapper) {

	caddr := conn.RemoteAddr()
	log.Info(caddr, " connected")

	// Close the connection when this routine ends.
	defer func() {
		conn.Close()
		log.Info(caddr, " disconnected")
	}()

	// TODO: Check if client address is allowed to perform request.
	r := bufio.NewReader(conn)

	for {
		// Here is the protocol:
		// - each command starts with a command name followed by multiple command arguments
		// - command name and arguments are separated by string "++++"
		// - commands are concatenated by character '\n'
		// - connection is terminiated when receiving the command "bye"
		//
		// Example: "clusterQstat\ngetBlockedJobsOfUser++++honlee\nbye\n"
		//
		// Whenever this for loop should continue for next command, the character '\a' is
		// send to the client for the next command.

		// Set initial read timeout to 3 seconds from now
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		m, err := r.ReadString('\n')

		log.Debug("message received: ", m)

		// io.EOF received, it implies that the connection's I/O is closed (e.g. client disconnect).
		// break the loop to close the connection properly.
		if err == io.EOF {
			break
		}

		// Error reading the command message, skip it and continue with the next command.
		if err != nil {
			log.Error(err)
			conn.Write([]byte("Error: " + err.Error() + "\n\a"))
			continue
		}

		// Command message is read.  Trim the last '\n' character to get the actual command.
		m = strings.TrimSuffix(m, "\n")
		if m == "bye" {
			break
		}

		// Empty message doesn't make sense, return '\a' to client for the next command.
		if m == "" {
			conn.Write([]byte{'\a'})
			continue
		}

		// Resolve the actual command name, arguments for making system call.
		cmdName, cmdArgs, err := cmdMapper.MapCommand(m)

		// Cannot resolve the command or the command is invalid.
		// Notify the client with '\a' for the next command.
		if err != nil {
			log.Error(err)
			conn.Write([]byte("Error: " + err.Error() + "\n\a"))
			continue
		}

		// Execute command and send a command output directly back to the connector.
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Env = os.Environ()
		cmd.Stdout = conn
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error(err)
			conn.Write([]byte("Error: " + err.Error() + ": " + cmdName + "\n\a"))
			continue
		}
		cmd.Wait()

		// Notify client with '\a' for the next command.
		conn.Write([]byte{'\a'})
	}
}

func validateJobID(id string) (jobFqid string, err error) {
	// Trim the job suffix
	sid := strings.Split(id, ".")[0]
	jobFqid = id
	// Check if the id is a digit number
	if match, _ := regexp.MatchString("^([0-9]+)$", sid); !match {
		err = errors.New("Invalid job id: " + id)
		return
	}
	// Compose jobFqid if the provided id is not the FQID
	if sid != id {
		jobFqid = id + "." + *trqServer
	}
	return
}

func validateUserID(id string) (err error) {
	err = errors.New("Invalid username: " + id)
	return
}

// CommandMapper defines
type CommandMapper interface {
	// MapCommand maps incoming TCP command into system call supported on the pbs_mom node.
	MapCommand(input string) (cmdName string, cmdArgs []string, err error)
}

// CommandMapperMom implements the CommandMapper interface for trqhelpd service running on pbs_mom node.
type CommandMapperMom struct {
}

// MapCommand maps incoming TCP command into system call supported on the pbs_mom node.
func (c CommandMapperMom) MapCommand(input string) (cmdName string, cmdArgs []string, err error) {
	// Split input to get command data, as we expect command and arguments
	// are separated by 4 '+' characters.
	data := strings.Split(input, "++++")
	switch data[0] {
	case "jobMemInfo":
		// Check current job memory usage via cgroups
		if len(data) < 2 {
			err = fmt.Errorf("Job id not provided: %v", data)
			return
		}
		jobFqid, verr := validateJobID(data[1])
		if verr != nil {
			err = verr
			return
		}
		cmdName = "cgget"
		cmdArgs = []string{"-r", "memory.usage_in_bytes", "-r", "memory.max_usage_in_bytes", "torque/" + jobFqid}
	case "getVncServices":
		// Get list of active vnc services of user
	default:
		err = errors.New("Command not found: " + data[0])
		return
	}
	return
}

// CommandMapperSrv implements the CommandMapper interface for trqhelpd service running on pbs_server node.
type CommandMapperSrv struct {
}

// MapCommand maps incoming TCP command into system call supported on on the pbs_server node.
func (c CommandMapperSrv) MapCommand(input string) (cmdName string, cmdArgs []string, err error) {
	// Split input to get command data, as we expect command and arguments
	// are separated by 4 '+' characters.
	data := strings.Split(input, "++++")
	switch data[0] {
	case "torqueConfig":
		// Get torque configuration from qmgr
		cmdName = "qmgr"
		cmdArgs = []string{"-c", "print server"}
	case "moabConfig":
		// Get moab configuration from moab configuration file
		cmdName = "cat"
		cmdArgs = []string{path.Join(*mdir, "etc", "moab.cfg")}
	case "clusterQstat":
		// Get whole cluster qstat
		cmdName = "qstat"
		cmdArgs = []string{"-a", "-t", "-G", "-n", "-1"}
	case "clusterQstatX":
		// Get whole cluster qstat in XML format
		cmdName = "qstat"
		cmdArgs = []string{"-x"}
	case "clusterFaireshare":
		// Get cluster fairshare status from the diagnose command of moab
		cmdName = "diagnose"
		cmdArgs = []string{"-f"}
	case "checkBlockedJob":
		// Get details of a blocked job
		if len(data) < 2 {
			err = fmt.Errorf("Job id not provided: %v", data)
			return
		}
		jobFqid, verr := validateJobID(data[1])
		if verr != nil {
			err = verr
			return
		}
		cmdName = "checkjob"
		cmdArgs = []string{"--xml", jobFqid}
	case "getBlockedJobsOfUser":
		// Get blocked jobs of a user
		if len(data) < 2 {
			err = fmt.Errorf("user id not provided: %v", data)
			return
		}
		if err = validateUserID(data[1]); err != nil {
			return
		}
		cmdName = "showq"
		cmdArgs = []string{"-b", "--xml", "-w", "user=" + data[1]}
	case "traceJob":
		// Trace a given job from the logs in the last 3 days
		if len(data) < 2 {
			err = fmt.Errorf("job id not provided: %v", data)
			return
		}
		jobFqid, verr := validateJobID(data[1])
		if verr != nil {
			err = verr
			return
		}
		cmdName = "tracejob"
		cmdArgs = []string{"-n", "3", jobFqid}
	default:
		err = errors.New("Command not found: " + data[0])
		return
	}
	return
}
