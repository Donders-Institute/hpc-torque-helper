package client

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// PrintClusterConfig prints configurations of Torque (pbs_server) and Moab services.
func PrintClusterConfig(trqhelpdHost string, trqhelpdPort int) {

	config := tls.Config{}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", trqhelpdHost, trqhelpdPort), &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()

	for _, m := range []string{"torqueConfig", "moabConfig", "bye"} {
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

// PrintJobMemoryInfo prints the memory usage of a running job.
func PrintJobMemoryInfo(jobID string, trqhelpdPort int) {
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
	cmd := exec.Command("qstat", "-x", jobID)
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

	// remove the eventual node attributes concatenated to the node's hostname with ":"
	data.Job.Host = strings.Split(data.Job.Host, ":")[0]
	log.Debugf("job data parsed from XML: %+v", data.Job)

	jdata := strings.Split(data.Job.Memset, ":")
	if jdata[0] == "" {
		log.Fatalf("Invalid job's execution host: %+v", data.Job)
	}

	config := tls.Config{}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", jdata[0], trqhelpdPort), &config)
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

// PrintClusterQstat prints all jobs in the memory of the Torque (pbs_server) service.
func PrintClusterQstat(trqhelpdHost string, trqhelpdPort int, xml bool) {
	config := tls.Config{}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", trqhelpdHost, trqhelpdPort), &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()

	cmd := "clusterQstat"
	if xml {
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

// PrintClusterTracejob prints job tracing logs available on the Torque (pbs_server) server.
func PrintClusterTracejob(jobID string, trqhelpdHost string, trqhelpdPort int) {

	config := tls.Config{}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", trqhelpdHost, trqhelpdPort), &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()

	cmds := []string{
		fmt.Sprintf("traceJob++++%s", jobID),
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
