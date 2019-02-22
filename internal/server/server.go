package server

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"syscall"

	pb "github.com/Donders-Institute/hpc-torque-helper/internal"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
)

// TorqueHelperSrv implements the gRPC interfaces exported by the TorqueHelper service running on the Torque/Moab server.
type TorqueHelperSrv struct {
	TorqueServer string
}

// validateJobID checks the validity of the given job id.
func (s *TorqueHelperSrv) validateJobID(id string) (jobFqid string, err error) {
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
		jobFqid = sid + "." + s.TorqueServer
	}
	return
}

// TraceJob returns job tracing logs available in the Torque server log.
func (s *TorqueHelperSrv) TraceJob(ctx context.Context, in *pb.JobInfoRequest) (out *pb.GeneralResponse, err error) {

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

	return
}

// GetJobBlockReason returns information from the `checkjob` command.  The output contains the reason why a job is not started
// by Moab.
func (s *TorqueHelperSrv) GetJobBlockReason(ctx context.Context, in *pb.JobInfoRequest) (out *pb.GeneralResponse, err error) {

	return
}

// GetBlockedJobsOfUser returns a list of jobs that are not started by Moab.
func (s *TorqueHelperSrv) GetBlockedJobsOfUser(ctx context.Context, in *pb.UserInfoRequest) (out *pb.GeneralResponse, err error) {

	return
}

// Qstat returns an tabular overview of all jobs in the memory of the Torque server.
func (s *TorqueHelperSrv) Qstat(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {

	return
}

// Qstatx returns XML output of all jobs in the memory of the Torque server.
func (s *TorqueHelperSrv) Qstatx(ctx context.Context, in *empty.Empty) (out *pb.GeneralResponse, err error) {

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
