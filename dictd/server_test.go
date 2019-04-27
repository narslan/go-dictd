package dictd

import (
	"bytes"
	"net"
	"testing"
)

var tcp DictServer

func init() {

	// Start the new server
	tcp := NewServer(":2500")
	// Run the servers in goroutines to stop blocking
	go func() {
		tcp.Run()
	}()

}

func TestServer_Running(t *testing.T) {
	// Simply check that the server is up and can
	// accept connections.
	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", "localhost:2500"},
	}

	for _, serv := range servers {
		conn, err := net.Dial(serv.protocol, serv.addr)
		if err != nil {
			t.Error("could not connect to server: ", err)
		}
		defer conn.Close()
	}
}

func TestServerRequest(t *testing.T) {

	servers := []struct {
		protocol string
		addr     string
	}{
		{"tcp", "localhost:2500"},
	}

	tt := []struct {
		test    string
		payload []byte
		want    []byte
	}{
		{"Sending a simple request returns result", []byte("hello world\n"), []byte("Request received: hello world")},
		{"Sending another simple request works", []byte("goodbye world\n"), []byte("Request received: goodbye world")},
	}

	for _, serv := range servers {
		for _, tc := range tt {
			t.Run(tc.test, func(t *testing.T) {
				conn, err := net.Dial(serv.protocol, serv.addr)
				if err != nil {
					t.Error("could not connect to server: ", err)
				}
				defer conn.Close()

				if _, err := conn.Write(tc.payload); err != nil {
					t.Error("could not write payload to server:", err)
				}

				out := make([]byte, 1024)
				if _, err := conn.Read(out); err == nil {
					if bytes.Compare(out, tc.want) == 0 {
						t.Error("response did match expected output")
					}
				} else {
					t.Error("could not read from connection")
				}
			})
		}
	}

}
