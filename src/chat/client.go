package chat

import (
	"encoding/json"
	"net"
	"zerotrust_chat/crypto"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
	"zerotrust_chat/logger"
)

var _ Client = &client{}

type client struct {
	buffer     []byte
	conn       net.Conn
	keyFactory crypto.KeyFactory
	secretKey  aes.Key
	targetAddr string
}

func (c client) GetId() string {
	return c.targetAddr
}

func NewClient(
	personalId string,
	targetAddr string,
	keyFactory crypto.KeyFactory,
	sessionManager SessionManager,
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

	client, err := makeClient(targetAddr, conn, keyFactory)
	if err != nil {
		return nil, err
	}

	err = client.handshake(personalId)
	if err != nil {
		return nil, err
	}

	sessionManager.Add(client)

	go func(c Client) {
		for {
			msg, err := c.Read()
			if err != nil {
				logger.Debug(err)
				break
			}
			println("recv:", msg) // TODO: implement observer pattern here to notify all subscriber
		}
	}(client)

	return client, nil
}

func makeClient(targetAddr string, conn net.Conn, keyFactory crypto.KeyFactory) (*client, error) {
	client := &client{
		conn:       conn,
		buffer:     make([]byte, 1024),
		keyFactory: keyFactory,
		targetAddr: targetAddr,
	}

	secretKey, err := keyFactory.GenerateAesSecretKey()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	client.secretKey = secretKey
	return client, nil
}

func (c *client) handshake(id string) error {
	pubKey, err := c.pubKeyRequest(id)
	if err != nil {
		return err
	}

	return c.shareSecretKey(pubKey)
}

func (c *client) internalRead() ([]byte, error) {
	numByte, err := c.conn.Read(c.buffer)
	if err != nil {
		return nil, err
	}
	return c.buffer[:numByte], nil
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

func (c *client) write(msg []byte) error {
	_, err := c.conn.Write(msg)
	return err
}

func (c *client) pubKeyRequest(id string) (rsa.PublicKey, error) {
	startConnectionRequest := startConnectionRequest{
		Id: id,
	}

	startRequest, err := json.Marshal(startConnectionRequest)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	_, err = c.conn.Write(startRequest)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	resp, err := c.internalRead()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	keyExchangeRequest := keyExchangeRequest{}
	err = json.Unmarshal(resp, &keyExchangeRequest)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return c.keyFactory.ConstructRsaPublicKey(
		keyExchangeRequest.PubKey,
	)
}

func (c *client) shareSecretKey(pubKey rsa.PublicKey) error {

	cipherSecretKey, err := pubKey.Encrypt([]byte(c.secretKey.ToString()))
	if err != nil {
		logger.Error(err)
		return err
	}
	keyExchangeResponse := keyExchangeResponse{SecretKey: cipherSecretKey}

	request, err := json.Marshal(keyExchangeResponse)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Debug("sharing secret:", string(request))

	err = c.write(request)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
