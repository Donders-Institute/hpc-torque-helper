package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"

	pb "github.com/Donders-Institute/hpc-torque-helper/internal/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var secret string

// TorqueHelperSrvClient implements client APIs for the TorqueHelperSrv service.
type TorqueHelperSrvClient struct {
	SrvHost     string
	SrvPort     int
	SrvCertFile string
}

// grpcConnect establishes client connection to the TorqueHelperMom service via gPRC.
func (c *TorqueHelperSrvClient) grpcConnect() (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile(c.SrvCertFile, c.SrvHost)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", c.SrvHost, c.SrvPort), grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Ping makes the gRPC call to the ping function on the TorqueHelperSrv service.
func (c *TorqueHelperSrvClient) Ping() error {

	conn, err := c.grpcConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTorqueHelperSrvServiceClient(conn)

	md := metadata.Pairs("token", pb.GetSecret())
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := client.Ping(ctx, &empty.Empty{})
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", out.GetResponseData())

	return nil
}

// PrintClusterConfig prints configurations of Torque (pbs_server) and Moab services.
func (c *TorqueHelperSrvClient) PrintClusterConfig() error {

	conn, err := c.grpcConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTorqueHelperSrvServiceClient(conn)

	md := metadata.Pairs("token", pb.GetSecret())
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// get torque config
	out, err := client.TorqueConfig(ctx, &empty.Empty{})
	if err != nil {
		return err
	}
	if err := printRPCOutput(out); err != nil {
		return err
	}

	// get moab config
	out, err = client.MoabConfig(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	return printRPCOutput(out)
}

// PrintClusterQstat prints all jobs in the memory of the Torque (pbs_server) service.
func (c *TorqueHelperSrvClient) PrintClusterQstat(xml bool) error {
	conn, err := c.grpcConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTorqueHelperSrvServiceClient(conn)

	md := metadata.Pairs("token", pb.GetSecret())
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var out *pb.GeneralResponse

	if xml {
		out, err = client.Qstat(ctx, &empty.Empty{})
		if err != nil {
			return err
		}

	} else {
		out, err = client.Qstatx(ctx, &empty.Empty{})
		if err != nil {
			return err
		}
	}

	return printRPCOutput(out)
}

// PrintClusterTracejob prints job tracing logs available on the Torque (pbs_server) server.
func (c *TorqueHelperSrvClient) PrintClusterTracejob(jobID string) error {
	conn, err := c.grpcConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTorqueHelperSrvServiceClient(conn)

	md := metadata.Pairs("token", pb.GetSecret())
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := client.TraceJob(ctx, &pb.JobInfoRequest{Jid: jobID, Xml: false})
	if err != nil {
		return err
	}
	return printRPCOutput(out)
}

// TorqueHelperMomClient implements client APIs for the TorqueHelperMom service.
type TorqueHelperMomClient struct {
	SrvHost     string
	SrvPort     int
	SrvCertFile string
}

// grpcConnect establishes client connection to the TorqueHelperMom service via gPRC.
func (c *TorqueHelperMomClient) grpcConnect() (*grpc.ClientConn, error) {
	creds, err := credentials.NewClientTLSFromFile(c.SrvCertFile, c.SrvHost)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", c.SrvHost, c.SrvPort), grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// PrintJobMemoryInfo prints the memory usage of a running job.
func (c *TorqueHelperMomClient) PrintJobMemoryInfo(jobID string) error {
	jobInfo, err := getJobQstatInfo(jobID)
	if err != nil {
		return err
	}

	// force Mom service host to the one the job is running.
	c.SrvHost = jobInfo.Host

	conn, err := c.grpcConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTorqueHelperMomServiceClient(conn)

	md := metadata.Pairs("token", pb.GetSecret())
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := client.JobMemInfo(ctx, &pb.JobInfoRequest{Jid: jobID, Xml: false})
	if err != nil {
		return err
	}
	return printRPCOutput(out)
}

// printRPCOutput prints output from a Unary gRPC call.
func printRPCOutput(out *pb.GeneralResponse) error {
	if out.GetExitCode() != 0 {
		return fmt.Errorf("grpc server process error: %+v (ec=%d)", out.GetErrorMessage(), out.GetExitCode())
	}
	fmt.Printf("%s\n", out.GetResponseData())
	return nil
}

// JobInfo contains information of the cluster job retrived from the `qstat -x` command.
type JobInfo struct {
	JobID  string `xml:"Job_Id"`
	Host   string `xml:"req_information>hostlist.0"`
	Memset string `xml:"memset_string"`
}

// getJobQstatInfo gets job information using `qstat -x` command directory.
func getJobQstatInfo(jobID string) (*JobInfo, error) {

	type Data struct {
		XMLName xml.Name `xml:"Data"`
		JobInfo JobInfo
	}

	// get the job's execution host
	cmd := exec.Command("qstat", "-x", jobID)
	cmd.Env = os.Environ()
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("cannot get job's execution host: %s", err)
	}
	log.Debug(string(b))

	data := Data{}
	if err := xml.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("cannot get job's execution host: %v", err)
	}

	// remove the eventual node attributes concatenated to the node's hostname with ":"
	data.JobInfo.Host = strings.Split(data.JobInfo.Host, ":")[0]
	log.Debugf("job data parsed from XML: %+v", data.JobInfo)

	jdata := strings.Split(data.JobInfo.Memset, ":")
	if jdata[0] == "" {
		return nil, fmt.Errorf("Invalid job's execution host: %+v", data.JobInfo)
	}

	return &data.JobInfo, nil
}
