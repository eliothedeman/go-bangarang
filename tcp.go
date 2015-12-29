package client

import (
	"net"

	"github.com/eliothedeman/bangarang/src/github.com/eliothedeman/bangarang/event"
	"github.com/eliothedeman/newman"
)

// A client which maintains an open tcp connection to the server
type TcpClient struct {
	rAddr string
	conn  *newman.Conn
}

// Create and return a new tcp client with it's tcp connection initilized
func NewTcpClient(addr string) (*TcpClient, error) {
	c := &TcpClient{
		rAddr: addr,
	}

	return c, c.dial()
}

// establish a tcp connection with the remote server
func (t *TcpClient) dial() error {
	conn, err := net.Dial("tcp", t.rAddr)
	if err != nil {
		return err
	}
	t.conn = newman.NewConn(conn)
	return nil
}

// Send the given event over the client's tcp connection
func (t *TcpClient) Send(e *event.Event) error {
	return t.conn.Write(e)
}
