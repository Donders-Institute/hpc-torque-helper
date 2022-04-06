package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"

	pb "github.com/Donders-Institute/hpc-torque-helper/internal/grpc"
	sys "github.com/Donders-Institute/hpc-torque-helper/internal/sys"
	"github.com/golang/protobuf/ptypes/empty"
)

// TorqueHelperSrv implements the gRPC interfaces exported by the TorqueHelper service running on the Torque/Moab server.
type TorqueHelperSrv struct {
	// TorqueServer is the hostname of the Torque/Moab server.
	TorqueServer string
}

// Ping returns a greeting message to the client.
func (s *TorqueHelperSrv) Ping(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {
	out = &pb.GeneralResponse{ExitCode: 0, ResponseData: "Hi there!", ErrorMessage: ""}
	return
}

// TraceJob returns job tracing logs available in the Torque server log.
func (s *TorqueHelperSrv) TraceJob(ctx context.Context, in *pb.JobInfoRequest) (out *pb.GeneralResponse, err error) {
	jobFqid, err := validateJobID(in.GetJid(), s.TorqueServer)
	if err != nil {
		return
	}

	stdout, stderr, ec := sys.ExecCmd("tracejob", []string{"-n", "3", jobFqid})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// TorqueConfig returns the configuration of the Torque server retrieved via 'qmgr' command.
func (s *TorqueHelperSrv) TorqueConfig(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {

	stdout, stderr, ec := sys.ExecCmd("qmgr", []string{"-c", "print server"})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// MoabConfig returns the configuration of the Moab server in the 'moab.cfg' file.
func (s *TorqueHelperSrv) MoabConfig(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {

	moabDir := os.Getenv("MOABHOMEDIR")

	if moabDir == "" {
		moabDir = "/usr/local/moab"
	}

	stdout, stderr, ec := sys.ExecCmd("cat", []string{path.Join(moabDir, "etc", "moab.cfg")})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// GetJobBlockReason returns information from the `checkjob` command.  The output contains the reason why a job is not started
// by Moab.
func (s *TorqueHelperSrv) GetJobBlockReason(ctx context.Context, in *pb.JobInfoRequest) (out *pb.GeneralResponse, err error) {
	jobFqid, err := validateJobID(in.GetJid(), s.TorqueServer)
	if err != nil {
		return
	}

	stdout, stderr, ec := sys.ExecCmd("checkjob", []string{"--xml", jobFqid})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// GetBlockedJobsOfUser returns a list of jobs that are not started by Moab.
func (s *TorqueHelperSrv) GetBlockedJobsOfUser(ctx context.Context, in *pb.UserInfoRequest) (out *pb.GeneralResponse, err error) {
	if err = validateUserID(in.GetUid()); err != nil {
		return
	}
	stdout, stderr, ec := sys.ExecCmd("showq", []string{"-b", "--xml", "-w", fmt.Sprintf("user=%s", in.GetUid())})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// Qstat returns an tabular overview of all jobs in the memory of the Torque server.
func (s *TorqueHelperSrv) Qstat(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {
	stdout, stderr, ec := sys.ExecCmd("qstat", []string{"-a", "-t", "-G", "-n", "-1"})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}
	return
}

// Qstatx returns XML output of all jobs in the memory of the Torque server.
func (s *TorqueHelperSrv) Qstatx(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {
	stdout, stderr, ec := sys.ExecCmd("qstat", []string{"-x"})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}
	return
}

// Checknode returns checknode output for a given node or ALL.
func (s *TorqueHelperSrv) Checknode(ctx context.Context, in *pb.NodeInfoRequest) (out *pb.GeneralResponse, err error) {

	var nid string
	if strings.ToLower(in.Nid) == "all" {
		// for all compute nodes
		nid = "ALL"
	} else {
		// for a specific node
		nid, err = validateNodeID(in.Nid)
		if err != nil {
			return
		}
	}

	// construct checknode command-line arguments
	args := []string{}
	if in.Xml {
		args = append(args, "--xml")
	}
	args = append(args, nid)

	// run checknode to get node information
	moabHomeDir := os.Getenv("MOABHOMEDIR")
	if moabHomeDir == "" {
		moabHomeDir = "/usr/local/moab"
	}
	stdout, stderr, ec := sys.ExecCmd(path.Join(moabHomeDir, "bin", "checknode"), args)
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}
	return
}

// TorqueHelperMom implements the gRPC interfaces exported by the TorqueHelper service running on the Mom server.
type TorqueHelperMom struct {
	// TorqueServer is the hostname of the Torque/Moab server.
	TorqueServer string
}

// JobMemInfo returns memory utilisation of a running job using the `cgget` cgroups command.
// The job must has a running process on the node the TorqueHelpMom service is running.
func (s *TorqueHelperMom) JobMemInfo(ctx context.Context, in *pb.JobInfoRequest) (out *pb.GeneralResponse, err error) {

	jobFqid, err := validateJobID(in.GetJid(), s.TorqueServer)
	if err != nil {
		return
	}

	stdout, stderr, ec := sys.ExecCmd(
		"cgget",
		[]string{
			"-r", "memory.usage_in_bytes",
			"-r", "memory.max_usage_in_bytes",
			"-r", "memory.limit_in_bytes",
			fmt.Sprintf("torque/%s", jobFqid),
		},
	)
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// TorqueHelperAcc implements the gRPC interfaces exported by the TorqueHelper service running on the access node of the Torque cluster.
type TorqueHelperAcc struct {
}

// GetVNCServers gets the VNC servers running on the server's local.
func (s *TorqueHelperAcc) GetVNCServers(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {

	stdout, stderr, ec := sys.ExecCmd("ps", []string{"h", "-o", "user:20,pid,command", "-C", "Xvnc"})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	// out = &pb.ServerListResponse{ExitCode: ec, ErrorMessage: stderr.String()}

	// if ec == 0 {
	// 	for {
	// 		l, err := stdout.ReadString('\n')
	// 		if err != nil && err != io.EOF {
	// 			// stop reading the stdout buffer in the raised error is not io.EOF
	// 			break
	// 		}
	// 		// parse the line and get the 1st and 12th field
	// 		scanner := bufio.NewScanner(strings.NewReader(l))
	// 		scanner.Split(bufio.ScanWords)
	// 		var cols []string
	// 		for scanner.Scan() {
	// 			cols = append(cols, scanner.Text())
	// 		}
	// 		if err := scanner.Err(); err != nil {
	// 			continue
	// 		}
	// 		if len(cols) < 12 {
	// 			continue
	// 		}
	// 		s := &pb.ServerListResponse_Server{Id: cols[11], Owner: cols[0]}
	// 		out.Servers = append(out.Servers, s)
	// 	}
	// }

	return
}

// validateUserID checks if the given user id is a valid system user.
func validateUserID(id string) (err error) {
	_, err = user.Lookup(id)
	return
}

// validateJobID checks the validity of the given job id, and returns
// the full job id with given torque server name as suffix.
func validateJobID(id, torqueServer string) (jobFqid string, err error) {
	// Trim the job suffix
	sid := strings.Split(id, ".")[0]
	jobFqid = id
	// Check if the short job id contains only digits
	if match, _ := regexp.MatchString("^([0-9]+)$", sid); !match {
		err = errors.New("Invalid job id: " + id)
		return
	}
	// Compose jobFqid if the provided id is a short job id
	if id == sid {
		jobFqid = sid + "." + torqueServer
	}
	return
}

// validateNodeID checks if the given node id `id` is a valid Torque node id.
// If the given id is a valid short hostname, the fully qualified hostname is constructed
// by appending `dccn.nl` to the short name.
func validateNodeID(id string) (nodeFqid string, err error) {

	// assuming compute node has hostname with prefix "dccn-c" followed by digits
	matched, _ := regexp.MatchString(`dccn-c[0-9]+`, id)
	if !matched {
		err = fmt.Errorf("invalid compute node: %s", id)
		return
	}

	// construct the fully qualified node name.
	nodeFqid = strings.Join([]string{strings.TrimSuffix(id, ".dccn.nl"), "dccn.nl"}, ".")

	return
}
