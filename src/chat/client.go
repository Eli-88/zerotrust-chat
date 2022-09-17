package chat

import (
	"net"
	"zerotrust_chat/crypto"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/logger"
)

type client struct {
	buffer           []byte
	conn             net.Conn
	keyFactory       crypto.KeyFactory
	secretKey        aes.Key
	targetAddr       string
	receiveHandler   ReceiveHandler
	chatReaderWriter ChatReaderWriter
}

func (c client) GetId() string {
	return c.targetAddr
}

func NewClient(
	personalId string,
	targetAddr string,
	keyFactory crypto.KeyFactory,
	sessionManager SessionManager,
	receiveHandler ReceiveHandler,
	chatReaderWriterFactory ChatReaderWriterFactory,
) (Session, error) {

	logger.Debug("connecting to:", targetAddr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if conn == nil {
		logger.Fatal("conn cannot be nil")
	}
	clientHandshake := NewClientHandshake(personalId, NewConn(conn), keyFactory)
	secretKey, err := clientHandshake.Handshake()

	if err != nil {
		return nil, err
	}
	client := makeClient(
		targetAddr,
		conn,
		keyFactory,
		receiveHandler,
		secretKey,
		chatReaderWriterFactory.Create(secretKey, NewConn(conn)),
	)

	sessionManager.Add(client)

	go func() {
		for {
			msg, err := client.Read()
			if err != nil {
				logger.Debug(err)
				break
			}
			client.receiveHandler.OnReceive(msg)
		}
	}()

	return client, nil
}

func makeClient(
	targetAddr string,
	conn net.Conn,
	keyFactory crypto.KeyFactory,
	receiveHandler ReceiveHandler,
	secretKey aes.Key,
	chatReaderWriter ChatReaderWriter,
) *client {
	return &client{
		conn:             conn,
		buffer:           make([]byte, 1024),
		keyFactory:       keyFactory,
		targetAddr:       targetAddr,
		receiveHandler:   receiveHandler,
		secretKey:        secretKey,
		chatReaderWriter: chatReaderWriter,
	}
}

func (c *client) Read() ([]ChatMessage, error) {
	numByte, err := c.conn.Read(c.buffer)
	if err != nil {
		logger.Debug(err)
		return nil, err
	}
	return c.chatReaderWriter.Read(c.buffer[:numByte])
}

func (c *client) Write(msg []byte) error {
	return c.chatReaderWriter.Write(msg)
}
