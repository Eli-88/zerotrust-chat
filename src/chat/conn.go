package chat

import "net"

var _ Conn = &conn{}

type conn struct {
	connection net.Conn
}

func NewConn(connection net.Conn) Conn {
	return &conn{
		connection: connection,
	}
}

func (c *conn) Read(b []byte) ([]byte, error) {
	numByte, err := c.connection.Read(b)
	if err != nil {
		return nil, err
	}
	return b[:numByte], nil
}

func (c *conn) Write(b []byte) (int, error) {
	return c.connection.Write(b)
}
