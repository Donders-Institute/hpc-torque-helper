package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	pb "github.com/Donders-Institute/hpc-torque-helper/internal/grpc"
	"github.com/Donders-Institute/hpc-torque-helper/internal/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	connHost    *string
	connPort    *int
	role        *string
	tlsCert     *string
	tlsKey      *string
	mdir        *string
	tdir        *string
	trqServer   *string
	optsVerbose *bool

	// secret specivied externally
	secret string
)

func init() {
	// Command-line arguments
	connHost = flag.String("h", "0.0.0.0", "set the ip `address` of the server")
	connPort = flag.Int("p", 60209, "set the port `number` of the server")
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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *connPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// load TLS credential
	creds, err := credentials.NewServerTLSFromFile(*tlsCert, *tlsKey)
	if err != nil {
		log.Fatalf("failed to setup tls: %v", err)
	}

	log.Debugf("accepting client secret: %s\n", pb.GetSecret())

	opts := []grpc.ServerOption{
		// Enable TLS for all incoming connections.
		grpc.Creds(creds),
		// Enable AuthInterceptor for token validation.
		grpc.UnaryInterceptor(pb.UnarySecretValidator),
	}

	grpcServer := grpc.NewServer(opts...)

	srv := server.TorqueHelperSrv{TorqueServer: *trqServer}

	pb.RegisterTorqueHelperSrvServiceServer(grpcServer, &srv)

	grpcServer.Serve(lis)
}
