package main

import (
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"

	log "github.com/sirupsen/logrus"
)

var (
	connHost    *string
	connPort    *int
	tlsCert     *string
	tlsKey      *string
	optsVerbose *bool
)

func init() {

	connHost = flag.String("h", "0.0.0.0", "set the ip `address` of the server")
	connPort = flag.Int("p", 60209, "set the port `number` of the server")
	tlsCert = flag.String("cert", "/etc/pki/tls/private/torque-helper.pem", "set the `path` of the TLS certificate")
	tlsKey = flag.String("key", "/etc/pki/tls/private/torque-helper.key", "set the `path` of the TLS private key")
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

	// TODO: Check if client address is allowed to perform request.

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.Error("Error reading: ", err.Error())
	}
	id := string(buf[:reqLen])
	// Check if the id is a digit number
	if match, _ := regexp.MatchString("([0-9]+)", id); !match {
		conn.Write([]byte("Job id must be a digital number: " + id))
		return
	}
	cmd := exec.Command("checkjob", "--xml", id)
	cmd.Env = []string{
		"PATH=/opt/cluster/bin:$PATH",
		"MOABHOMEDIR=/opt/cluster/external/moab",
	}
	// Send a command output directly back to the connector.
	cmd.Stdout = conn
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}
