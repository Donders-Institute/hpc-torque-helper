package main

import (
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	trqHelpdPort *int
	optsVerbose  *bool
)

func init() {

	// Command-line arguments
	trqHelpdPort = flag.Int("p", 60209, "set the `port` number of the trqhelpd")
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
	fmt.Printf("\nGet the memory usage information of a job.\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS] jobId\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {

	// command-line arguments
	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		log.Fatal(fmt.Sprintf("require one job id: %+v", args))
	}

	// defining data structure for unmarshalling qstat's XML document
	type Job struct {
		JobID  string `xml:"Job_Id"`
		Host   string `xml:"req_information>hostlist.0"`
		Memset string `xml:"memset_string"`
	}

	type Data struct {
		XMLName xml.Name `xml:"Data"`
		Job     Job
	}

	// get the job's execution host
	cmd := exec.Command("qstat", "-x", args[0])
	cmd.Env = os.Environ()
	b, err := cmd.Output()
	if err != nil {
		log.Fatalf("cannot get job's execution host: %s", err)
	}
	log.Debug(string(b))

	data := Data{}
	if err := xml.Unmarshal(b, &data); err != nil {
		log.Fatalf("cannot get job's execution host: %v", err)
	}
	log.Debugf("job exec host: %+v", data.Job)

	jdata := strings.Split(data.Job.Memset, ":")
	if jdata[0] == "" {
		log.Fatalf("Invalid job's execution host: %v", data.Job)
	}

	config := tls.Config{}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", jdata[0], *trqHelpdPort), &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()

	cmds := []string{
		fmt.Sprintf("jobMemInfo++++%s", data.Job.JobID),
		"bye",
	}

	for _, m := range cmds {
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
