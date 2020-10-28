package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
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

// GetNodeResourceStatus returns node resource status extracted from the
// output of the `checknode`.
func (c *TorqueHelperSrvClient) GetNodeResourceStatus(nodeID string) ([]NodeResourceStatus, error) {
	conn, err := c.grpcConnect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTorqueHelperSrvServiceClient(conn)

	md := metadata.Pairs("token", pb.GetSecret())
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := client.Checknode(ctx, &pb.NodeInfoRequest{Nid: nodeID, Xml: false})
	if err != nil {
		return nil, err
	}

	// parsing the out to return node resource
	return parseChecknodeXML([]byte(out.ResponseData))
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

	log.Debugf("jobInfo: %+v\n", jobInfo)

	// check if job's Host attribute is available
	if jobInfo.Host == "" {
		return fmt.Errorf("unknown job execution host (%s)", jobInfo.Host)
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
	if err != nil {
		return
	}

	// The code below parses output to VNCServer object.  An example output is below:
	// user1                 1552 /usr/bin/Xvnc :51 -auth ...
	// user2                 2050 /usr/bin/Xvnc :92 -auth ...
	// user3                 2862 /usr/bin/Xvnc :11 -auth  ...
	s := bufio.NewScanner(strings.NewReader(out.GetResponseData()))
	for s.Scan() {
		ws := bufio.NewScanner(strings.NewReader(s.Text()))
		ws.Split(bufio.ScanWords)
		vnc := VNCServer{}
		cnt := 0
		for ws.Scan() {
			cnt++
			switch cnt {
			case 1: // colume 1 - owner
				vnc.Owner = ws.Text()
			case 4: // colume 4 - vnc display number
				vnc.ID = fmt.Sprintf("%s%s", c.SrvHost, ws.Text())
				servers = append(servers, vnc)
			default:
				// do nothing here!!
			}
		}
		if err := ws.Err(); err != nil {
			log.Warnf("error parsing vnc owner and display: %+v", err)
		}
	}
	if err := s.Err(); err != nil {
		log.Warnf("error parsing vnc owner and display: %+v\n", err)
	}

	return
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

// NodeResourceStatus defines the data structure of the node resource status extracted from
// the XML output of `checknode --xml`.
type NodeResourceStatus struct {
	ID          string
	State       string
	Features    []string
	TotalProcs  int
	AvailProcs  int
	TotalMemGB  int
	AvailMemGB  int
	TotalDiskGB int
	AvailDiskGB int
	TotalGPUS   int
	AvailGPUS   int
}

// UnmarshalXML implemented `xml.Unmarshaler` for node resource data.
func (c *NodeResourceStatus) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "NODEID":
			c.ID = attr.Value
		case "NODESTATE":
			c.State = attr.Value
		case "FEATURES":
			c.Features = strings.Split(attr.Value, ",")
		case "RCPROC":
			nproc, err := strconv.Atoi(attr.Value)
			if err != nil {
				return err
			}
			c.TotalProcs = nproc
		case "RAPROC":
			nproc, err := strconv.Atoi(attr.Value)
			if err != nil {
				return err
			}
			c.AvailProcs = nproc
		case "RCDISK":
			sizeMB, err := strconv.Atoi(attr.Value)
			if err != nil {
				return err
			}
			c.TotalDiskGB = sizeMB / 1024
		case "RADISK":
			sizeMB, err := strconv.Atoi(attr.Value)
			if err != nil {
				return err
			}
			c.AvailDiskGB = sizeMB / 1024
		case "RCMEM":
			sizeMB, err := strconv.Atoi(attr.Value)
			if err != nil {
				return err
			}
			c.TotalMemGB = sizeMB / 1024
		case "RAMEM":
			sizeMB, err := strconv.Atoi(attr.Value)
			if err != nil {
				return err
			}
			c.AvailMemGB = sizeMB / 1024
		case "GRES":
			for _, r := range strings.Split(attr.Value, ";") {
				if strings.HasPrefix(r, "gpus=") {
					c.TotalGPUS, _ = strconv.Atoi(strings.TrimPrefix(r, "gpus="))
				}
			}
		case "AGRES":
			for _, r := range strings.Split(attr.Value, ";") {
				if strings.HasPrefix(r, "gpus=") {
					c.AvailGPUS, _ = strconv.Atoi(strings.TrimPrefix(r, "gpus="))
				}
			}
		default:
		}
	}

	return d.Skip()
}

// parseChecknodeXML parses the output of `checknode --xml` and returns the Node data structure.
func parseChecknodeXML(xmlData []byte) ([]NodeResourceStatus, error) {

	type data struct {
		XMLName xml.Name             `xml:"Data"`
		Nodes   []NodeResourceStatus `xml:"node"`
	}

	r := data{}
	if err := xml.Unmarshal(xmlData, &r); err != nil {
		return nil, fmt.Errorf("cannot get node resource status: %v", err)
	}

	return r.Nodes, nil
}
