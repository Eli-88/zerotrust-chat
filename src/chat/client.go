package chat

import (
	"encoding/json"
	"net"
	"zerotrust_chat/crypto"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/logger"
)

var _ Client = &client{}

type client struct {
	buffer         []byte
	conn           net.Conn
	keyFactory     crypto.KeyFactory
	secretKey      aes.Key
	targetAddr     string
	receiveHandler ReceiveHandler
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
) (Client, error) {

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

	client := makeClient(
		targetAddr,
		conn,
		keyFactory,
		receiveHandler,
	)
	if err != nil {
		return nil, err
	}

	clientHandshake := NewClientHandshake(personalId, conn, keyFactory)
	secretKey, err := clientHandshake.Handshake()

	if err != nil {
		return nil, err
	}

	client.secretKey = secretKey

	sessionManager.Add(client)

	go func() {
		for {
			msg, err := client.Read()
			if err != nil {
				logger.Debug(err)
				break
			}
			msgCpy := make([]byte, len(msg))
			copy(msgCpy, msg)
			client.receiveHandler.OnReceive(string(msgCpy))
		}
	}()

	return client, nil
}

func makeClient(
	targetAddr string,
	conn net.Conn,
	keyFactory crypto.KeyFactory,
	receiveHandler ReceiveHandler,
) *client {
	return &client{
		conn:           conn,
		buffer:         make([]byte, 1024),
		keyFactory:     keyFactory,
		targetAddr:     targetAddr,
		receiveHandler: receiveHandler,
	}
}

func (c *client) Read() (string, error) {
	numByte, err := c.conn.Read(c.buffer)
	if err != nil {
		return "", err
	}

	chatMessage := ChatMessage{}
	json.Unmarshal(c.buffer[:numByte], &chatMessage)

	recvMsg, err := c.secretKey.Decrypt(chatMessage.Data)

	if err != nil {
		return "", err
	}
	return string(recvMsg), nil
}

func (c *client) Write(msg []byte) error {
	encryptedMsg, err := c.secretKey.Encrypt(msg)
	if err != nil {
		return err
	}

	chatMessage := ChatMessage{
		Data: encryptedMsg,
	}
	msgToBeSent, err := json.Marshal(chatMessage)

	if err != nil {
		return err
	}
	_, err = c.conn.Write(msgToBeSent)
	return err
}
