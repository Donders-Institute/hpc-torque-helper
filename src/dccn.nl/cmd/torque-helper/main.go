package main

import (
    "net"
    "os"
    "os/exec"
    log "github.com/sirupsen/logrus"
)

const (
    CONN_HOST = "0.0.0.0"
    CONN_PORT = "60209"
    CONN_TYPE = "tcp"
)

func init() {
    // set logging
    log.SetOutput(os.Stderr)
}

func main() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        log.Error("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    log.Info("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            log.Error("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    // Close the connection when this routine ends.
    defer conn.Close()

    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)
    // Read the incoming connection into the buffer.
    reqLen, err := conn.Read(buf)
    if err != nil {
        log.Error("Error reading:", err.Error())
    }
    id := string(buf[:reqLen])
    cmd := exec.Command("checkjob","--xml",id)
    cmd.Env = []string {
        "PATH=/opt/cluster/bin:$PATH",
	"MOABHOMEDIR=/opt/cluster/external/moab",
    }
    // Send a command output directly back to the connector.
    cmd.Stdout = conn
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(err)
    }
}
