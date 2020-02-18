package client

import (
	"io/ioutil"
	"testing"
)

var (
	srvClient = TorqueHelperSrvClient{
		SrvHost:     "localhost",
		SrvPort:     60209,
		SrvCertFile: "testdata/cert/TestServer.crt",
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

func TestParseQstatXML(t *testing.T) {
	xmldata, err := ioutil.ReadFile("testdata/qstat.xml")
	if err != nil {
		t.Errorf("fail to read XML data from test file.\n")
	}
	jobinfo, err := parseQstatXML(xmldata)
	if err != nil {
		t.Errorf("fail parsing XML data: %+v\n", err)
	}
	if jobinfo.JobID != "34854814.dccn-l029.dccn.nl" || jobinfo.Host != "dccn-c360.dccn.nl" {
		t.Errorf("unexpected data: %+v\n", jobinfo)
	}
}
