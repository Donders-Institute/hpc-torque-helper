package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"strings"
	"syscall"

	pb "github.com/Donders-Institute/hpc-torque-helper/internal/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
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
	jobFqid, verr := validateJobID(in.GetJid(), s.TorqueServer)
	if verr != nil {
		err = verr
		return
	}

	stdout, stderr, ec := execCmd("tracejob", []string{"-n", "3", jobFqid})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// TorqueConfig returns the configuration of the Torque server retrieved via 'qmgr' command.
func (s *TorqueHelperSrv) TorqueConfig(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {

	stdout, stderr, ec := execCmd("qmgr", []string{"-c", "print server"})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// MoabConfig returns the configuration of the Moab server in the 'moab.cfg' file.
func (s *TorqueHelperSrv) MoabConfig(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {

	moabDir := os.Getenv("MOABHOMEDIR")

	if moabDir == "" {
		moabDir = "/usr/local/moab"
	}

	stdout, stderr, ec := execCmd("cat", []string{path.Join(moabDir, "etc", "moab.cfg")})
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

	stdout, stderr, ec := execCmd("checkjob", []string{"--xml", jobFqid})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// GetBlockedJobsOfUser returns a list of jobs that are not started by Moab.
func (s *TorqueHelperSrv) GetBlockedJobsOfUser(ctx context.Context, in *pb.UserInfoRequest) (out *pb.GeneralResponse, err error) {
	if err = validateUserID(in.GetUid()); err != nil {
		return
	}
	stdout, stderr, ec := execCmd("showq", []string{"-b", "--xml", "-w", fmt.Sprintf("user=%s", in.GetUid())})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// Qstat returns an tabular overview of all jobs in the memory of the Torque server.
func (s *TorqueHelperSrv) Qstat(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {
	stdout, stderr, ec := execCmd("qstat", []string{"-a", "-t", "-G", "-n", "-1"})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}
	return
}

// Qstatx returns XML output of all jobs in the memory of the Torque server.
func (s *TorqueHelperSrv) Qstatx(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {
	stdout, stderr, ec := execCmd("qstat", []string{"-x"})
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

	jobFqid, verr := validateJobID(in.GetJid(), s.TorqueServer)
	if verr != nil {
		err = verr
		return
	}

	stdout, stderr, ec := execCmd("cgget",
		[]string{"-r", "memory.usage_in_bytes", "-r", "memory.max_usage_in_bytes", fmt.Sprintf("torque/%s", jobFqid)})
	out = &pb.GeneralResponse{ExitCode: ec, ResponseData: stdout.String(), ErrorMessage: stderr.String()}

	return
}

// validateUserID checks if the given user id is a valid system user.
func validateUserID(id string) (err error) {
	_, err = user.Lookup(id)
	return
}

func execCmd(cmdName string, cmdArgs []string) (stdout, stderr bytes.Buffer, ec int32) {
	// Execute command and catch the stdout and stderr as byte buffer.
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Env = os.Environ()
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	ec = 0
	if err := cmd.Run(); err != nil {
		log.Error(err)
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			ec = int32(ws.ExitStatus())
		} else {
			ec = 1
		}
	}
	return
}

// validateJobID checks the validity of the given job id, and returns
// the full job id with given torque server name as suffix.
func validateJobID(id, torqueServer string) (jobFqid string, err error) {
	// Trim the job suffix
	sid := strings.Split(id, ".")[0]
	jobFqid = id
	// Check if the id is a digit number
	if match, _ := regexp.MatchString("^([0-9]+)$", sid); !match {
		err = errors.New("Invalid job id: " + id)
		return
	}
	// Compose jobFqid if the provided id is not the FQID
	if sid != id {
		jobFqid = sid + "." + torqueServer
	}
	return
}
