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
	optsVerbose *bool
)

func init() {

	// Command-line arguments
	trqhelpdSrv = flag.String("s", "torque.dccn.nl:60209", "set the service `endpoint` of the trqhelpd")
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
	log.Println("client: connected to: ", conn.RemoteAddr())

	message := "clusterQstat++++"
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 4096)

	for {
		n, err = conn.Read(reply)
		fmt.Printf("%s", reply[:n])
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
