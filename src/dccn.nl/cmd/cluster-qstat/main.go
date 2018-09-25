package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	trqhelpdSrv *string
	optsXML     *bool
	optsVerbose *bool
)

func init() {

	// Command-line arguments
	trqhelpdSrv = flag.String("s", "torque.dccn.nl:60209", "set the service `endpoint` of the trqhelpd")
	optsXML = flag.Bool("x", false, "print jobs in XML format")
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
	fmt.Printf("\nList jobs of the HPC cluster.\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {
	config := tls.Config{}

	conn, err := tls.Dial("tcp", *trqhelpdSrv, &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()

	cmd := "clusterQstat"
	if *optsXML {
		cmd += "X"
	}

	for _, m := range []string{cmd, "bye"} {
		_, err := conn.Write(append([]byte(m), '\n'))
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}

		term := false
		reply := make([]byte, 4096)
		for {

			n, err := conn.Read(reply)

			// Error in reading command output or io.EOF
			if err != nil {
				if err != io.EOF {
					log.Error(err)
				}
				term = true
				break
			}

			// Received '\a' from server indicating the end of the command output
			if reply[n-1] == '\a' {
				if n > 0 {
					fmt.Printf("%s", reply[:n-1])
				}
				break
			}

			// Received a part of the command output
			fmt.Printf("%s", reply[:n])
		}

		// stop sending more command if the connection has been terminated.
		if term {
			break
		}
	}
}
