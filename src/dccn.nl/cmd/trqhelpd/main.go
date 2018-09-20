package main

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"strings"
	"time"
	"bufio"

	log "github.com/sirupsen/logrus"
)

var (
	connHost    *string
	connPort    *int
	tlsCert     *string
	tlsKey      *string
	mdir        *string
	tdir        *string
	optsVerbose *bool
)

func init() {

	// Command-line arguments
	connHost = flag.String("h", "0.0.0.0", "set the ip `address` of the server")
	connPort = flag.Int("p", 60209, "set the port `number` of the server")
	mdir = flag.String("m", os.Getenv("MOABHOMEDIR"), "set the `path` of Moab installation, usually referred by $MOABHOMEDIR")
	tdir = flag.String("t", os.Getenv("TORQUEHOME"), "set the `path` of the Torque installation")
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

	log.Infof("Listening on %s:%d\n", *connHost, *connPort)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Error("Error accepting:", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

	caddr := conn.RemoteAddr()
	log.Info(caddr, " connected")

	// Close the connection when this routine ends.
	defer func() {
		conn.Close()
		log.Info(caddr, " disconnected")
	}()

	// Set read timeout to 5 seconds from now
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	// TODO: Check if client address is allowed to perform request.

	// Make a buffer to hold incoming data.
	buf := bufio.NewReader(conn)

	for {
		// Read the incoming message until the first '\n'.
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Error("Error reading: ", err.Error())
			conn.Write([]byte("Error reading: " + err.Error()))
			return
		}
		// Leave loop when the message received is 'bye'
		if msg == "bye" {
			break
		}
		// Switch to right command based on client input
		cmdName, cmdArgs, err := switchCommand(msg))
		if err != nil {
			log.Error(err)
			conn.Write([]byte(err.Error()))
			return
		}
		// Execute command and send a command output directly back to the connector.
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s/bin:%s/bin:$PATH", *tdir, *mdir))
		if *tdir != "" {
			cmd.Env = append(cmd.Env, "TORQUEHOME="+*tdir)
		}
		if *mdir != "" {
			cmd.Env = append(cmd.Env, "MOABHOMEDIR="+*mdir)
		}
		cmd.Stdout = conn
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error(err)
			conn.Write([]byte("Error checkjob: " + err.Error()))
			return
		}
	}
}

func switchCommand(input string) (cmdName string, cmdArgs []string, err error) {
	// Split input to get command data, as we expect command and arguments
	// are separated by 4 '+' characters.
	data := strings.Split(input, "++++")
	switch data[0] {
	case "torqueConfig":
		// Get torque configuration from qmgr
		cmdName = "qmgr"
		cmdArgs = []string{"-c", "'print server'"}
	case "moabConfig":
		mdir, ok := os.LookupEnv("MOABHOMEDIR")
		if !ok {
			mdir = "/usr/local/moab"
		}
		// Get moab configuration from moab configuration file
		cmdName = "cat"
		cmdArgs = []string{path.Join(mdir, "moab.cfg")}
	case "clusterQstat":
		// Get whole cluster qstat
		cmdName = "qstat"
		cmdArgs = []string{"-atGn1"}
	case "clusterFaireshare":
		// Get cluster fairshare status from the diagnose command of moab
		cmdName = "diagnose"
		cmdArgs = []string{"-f"}
	case "checkBlockedJob":
		// Check if the id is a digit number
		if match, _ := regexp.MatchString("([0-9]+)", data[1]); !match {
			err = errors.New("Invalid job id: " + data[1])
			return
		}
		cmdName = "checkjob"
		cmdArgs = []string{"--xml", data[1]}
	case "getBlockedJobsOfUser":
		// TODO: Check if the uid is a valid user
		if _, ierr := user.Lookup(data[1]); ierr != nil {
			err = errors.New("Invalid username: " + data[1])
			return
		}
		cmdName = "showq"
		cmdArgs = []string{"-b", "--xml", "-w", "user=" + data[1]}
	default:
		err = errors.New("Command not found: " + data[0])
		return
	}
	return
}

// ConnWriter modifies commandline output before writing to socket connection.
type ConnWriter struct {
	Conn net.Conn
}

func (c ConnWriter) Write(p []byte) (n int, err error) {
	if len(p) > 0 && p[0] == '\n' {
		return c.Conn.Write(p[1:])
	}
	return c.Conn.Write(p)
}
