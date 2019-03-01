package client

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/Donders-Institute/hpc-torque-helper/internal/sys"

	pb "github.com/Donders-Institute/hpc-torque-helper/internal/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

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
		out, err = client.Qstatx(ctx, &empty.Empty{})
		if err != nil {
			return err
		}

	} else {
		out, err = client.Qstat(ctx, &empty.Empty{})
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

	xmldata, err := jobQstatXML(jobID)
	if err != nil {
		return err
	}

	jobInfo, err := parseQstatXML(xmldata)
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

// TorqueHelperAccClient implements client APIs for the TorqueHelperAcc service.
type TorqueHelperAccClient struct {
	SrvHost     string
	SrvPort     int
	SrvCertFile string
}

// grpcConnect establishes client connection to the TorqueHelperMom service via gPRC.
func (c *TorqueHelperAccClient) grpcConnect() (*grpc.ClientConn, error) {
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

// VNCServer defines data structure of a VNC server.
type VNCServer struct {
	// ID is the VNC server id, e.g. mentat001.dccn.nl:1
	ID string
	// Owner is the VNC server owner's user id
	Owner string
}

// GetVNCServers gets a list of VNC servers.
func (c *TorqueHelperAccClient) GetVNCServers() (servers []VNCServer, err error) {

	conn, err := c.grpcConnect()
	if err != nil {
		return
	}
	defer conn.Close()

	client := pb.NewTorqueHelperAccServiceClient(conn)

	md := metadata.Pairs("token", pb.GetSecret())
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := client.GetVNCServers(ctx, &empty.Empty{})

	return nil, printRPCOutput(out)

	// if err != nil {
	// 	return
	// }
	// for _, s := range out.GetServers() {
	// 	servers = append(servers, VNCServer{ID: fmt.Sprintf("%s%s", c.SrvHost, s.GetId()), Owner: s.GetOwner()})
	// }

	//return
}

// printRPCOutput prints output from a Unary gRPC call.
func printRPCOutput(out *pb.GeneralResponse) error {
	if out.GetExitCode() != 0 {
		return fmt.Errorf("grpc server process error: %s (ec=%d)", out.GetErrorMessage(), out.GetExitCode())
	}
	fmt.Printf("%s\n", out.GetResponseData())
	return nil
}

// Job contains information of the cluster job retrived from the `qstat -x` command.
type Job struct {
	JobID  string `xml:"Job_Id"`
	Host   string `xml:"req_information>hostlist.0"`
	Memset string `xml:"memset_string"`
}

// jobQstatXML runs `qstat -x` locally to get the full job information in XML format.
func jobQstatXML(jobID string) (xmlData []byte, err error) {

	var stdout, stderr bytes.Buffer
	stdout, stderr, ec := sys.ExecCmd("qstat", []string{"-x", jobID})
	if ec != 0 {
		err = fmt.Errorf("%s: (ec=%d)", string(stderr.Bytes()), ec)
		return
	}
	xmlData = stdout.Bytes()
	return
}

// parseQstatXML parses the output of `qstat -x` and returns the Job data structure.
func parseQstatXML(xmlData []byte) (*Job, error) {
	type data struct {
		XMLName xml.Name `xml:"Data"`
		Job     Job
	}

	d := data{}
	if err := xml.Unmarshal(xmlData, &d); err != nil {
		return nil, fmt.Errorf("cannot get job's execution host: %v", err)
	}

	// remove the eventual node attributes concatenated to the node's hostname with ":"
	d.Job.Host = strings.Split(d.Job.Host, ":")[0]
	log.Debugf("job data parsed from XML: %+v", d.Job)

	return &d.Job, nil
}
