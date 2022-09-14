package chat

import "net"

var _ HandshakeConn = &handshakeConn{}

type handshakeConn struct {
	conn net.Conn
}

func MakeHandshakeConn(conn net.Conn) HandshakeConn {
	return &handshakeConn{
		conn: conn,
	}
}

func (h *handshakeConn) Read(b []byte) ([]byte, error) {
	numByte, err := h.conn.Read(b)
	if err != nil {
		return nil, err
	}
	return b[:numByte], nil
}

func (h *handshakeConn) Write(b []byte) (int, error) {
	return h.conn.Write(b)
}
