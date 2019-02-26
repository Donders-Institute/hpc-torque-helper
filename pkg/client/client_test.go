package client

import (
	"os"
	"path"
	"testing"
)

var (
	srvClient = TorqueHelperSrvClient{
		SrvHost:     "localhost",
		SrvPort:     60209,
		SrvCertFile: path.Join(os.Getenv("GOPATH"), "src/github.com/Donders-Institute/hpc-torque-helper/test/cert/TestServer.crt"),
	}
)

// TestSrvPing performs test on the ping function of the TorqueHelperSrv service.
func TestSrvPing(t *testing.T) {
	if err := srvClient.Ping(); err != nil {
		t.Errorf("fail to ping: %+v\n", err)
	}
}

func TestSrvPrintClusterQstat(t *testing.T) {
	if err := srvClient.PrintClusterQstat(true); err != nil {
		t.Errorf("fail to print qstat: %+v\n", err)
	}
}

func TestSrvPrintClusterConfig(t *testing.T) {
	if err := srvClient.PrintClusterConfig(); err != nil {
		t.Errorf("fail to print cluster status: %+v\n", err)
	}
}
