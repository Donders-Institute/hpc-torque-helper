package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sync"

	pb "github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

var (
	optsVerbose *bool
	logfiles    []string
	optsWorkers *int
	optsOutCsv  *bool
	optsOutput  *string
)

func init() {
	// Command-line arguments
	optsVerbose = flag.Bool("v", false, "print debug messages")
	optsWorkers = flag.Int("p", 1, "`number` of parallel workers for parsing multiple XML files")
	optsOutCsv = flag.Bool("c", false, "print accounting table in CSV format")
	optsOutput = flag.String("o", "stdout", "write accounting output to `file`")

	flag.Usage = usage

	flag.Parse()

	logfiles = flag.Args()

	if len(logfiles) < 1 {
		usage()
		log.Fatal("missing log file")
	}

	// set logging level
	llevel := log.InfoLevel
	if *optsVerbose {
		llevel = log.DebugLevel
	}
	log.SetLevel(llevel)

	// set logging output
	log.SetOutput(os.Stdout)
}

func usage() {
	fmt.Printf("\nSimple HPC job accounting by parsing Torque's log files.\n")
	fmt.Printf("\nUSAGE: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("\nOPTIONS:\n")
	flag.PrintDefaults()
	fmt.Printf("\n")
}

func main() {

	nworkers := *optsWorkers

	ichan := make(chan string, 2*nworkers)
	ochan := make(chan JobInfo, 1000*nworkers)

	// channel for monitoring progress
	pchan := make(chan int)
	bar := pb.Default(int64(len(logfiles)))

	var wg sync.WaitGroup
	for i := 0; i < nworkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range ichan {
				jobs, err := ParseJobLog(f)
				if err != nil {
					log.Errorf("[%s] %s", f, err)
				} else {
					for _, jinfo := range jobs.JobInfo {
						ochan <- jinfo
					}
				}
				pchan <- 1
			}
		}()
	}

	// wait until all workers are finished, close the output and monitor channel
	go func() {
		wg.Wait()
		close(ochan)
		close(pchan)
	}()

	log.Debugf("processing %d log files ...", len(logfiles))
	go func() {
		for _, f := range logfiles {
			ichan <- f
		}
		close(ichan)
	}()

	gaccounts := []Account{}
	uaccounts := []Account{}
loop:
	for {
		select {
		case jinfo := <-ochan:
			log.Debugf("%s: %+v\n", jinfo.ID, jinfo)

			// group accounting
			if idx := FindAccount(gaccounts, jinfo.Group); idx == -1 {
				gaccounts = append(gaccounts, Account{
					ID:       jinfo.Group,
					Njobs:    1,
					Walltime: uint64(jinfo.Used.Walltime),
					Memory:   uint64(jinfo.Used.Memory),
				})
			} else {
				gaccounts[idx].Njobs += 1
				gaccounts[idx].Walltime += uint64(jinfo.Used.Walltime)
				gaccounts[idx].Memory += uint64(jinfo.Requested.Memory)
			}

			// user accounting
			if idx := FindAccount(uaccounts, jinfo.User); idx == -1 {
				uaccounts = append(uaccounts, Account{
					ID:       jinfo.User,
					Njobs:    1,
					Walltime: uint64(jinfo.Used.Walltime),
					Memory:   uint64(jinfo.Used.Memory),
				})
			} else {
				uaccounts[idx].Njobs += 1
				uaccounts[idx].Walltime += uint64(jinfo.Used.Walltime)
				uaccounts[idx].Memory += uint64(jinfo.Requested.Memory)
			}
		case p, ok := <-pchan:
			if !ok {
				break loop
			}
			bar.Add(p)
		}
	}

	log.Debugf("producing accounting data ...")
	f := os.Stdout
	if *optsOutput != "stdout" {
		var err error
		f, err = os.OpenFile(*optsOutput, os.O_RDWR|os.O_CREATE, 0640)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	switch {
	case *optsOutCsv:
		w := csv.NewWriter(f)
		w.Write([]string{
			"type",
			"id",
			"njobs",
			"walltime",
			"memory",
		})
		for _, account := range gaccounts {
			w.Write([]string{
				"group",
				account.ID,
				fmt.Sprintf("%d", account.Njobs),
				fmt.Sprintf("%d", account.Walltime),
				fmt.Sprintf("%d", account.Memory),
			})
		}
		for _, account := range uaccounts {
			w.Write([]string{
				"user",
				account.ID,
				fmt.Sprintf("%d", account.Njobs),
				fmt.Sprintf("%d", account.Walltime),
				fmt.Sprintf("%d", account.Memory),
			})
		}

		// Write any buffered data to the underlying writer (standard output).
		w.Flush()

		if err := w.Error(); err != nil {
			log.Fatal(err)
		}

	default:
		w := bufio.NewWriter(f)
		fmt.Fprintf(w, "Group Accounts:\n")
		fmt.Fprintf(w, "%s\n", AccountTabletHeader())
		for _, account := range gaccounts {
			fmt.Fprintf(w, "%s\n", account)
		}
		fmt.Fprintf(w, "%s\n", AccountTabletSeparator())

		fmt.Fprintf(w, "User Accounts:\n")
		fmt.Fprintf(w, "%s\n", AccountTabletHeader())
		for _, account := range uaccounts {
			fmt.Fprintf(w, "%s\n", account)
		}
		fmt.Fprintf(w, "%s\n", AccountTabletSeparator())
		w.Flush()
	}
}
