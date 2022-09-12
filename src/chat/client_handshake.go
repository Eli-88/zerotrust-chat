package chat

import (
	"encoding/json"
	"zerotrust_chat/crypto"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
	"zerotrust_chat/logger"
)

var _ HandShake = &clientHandshake{}

type clientHandshake struct {
	id         string
	conn       HandshakeConn
	buffer     []byte
	keyFactory crypto.KeyFactory
}

func NewClientHandshake(
	id string,
	conn HandshakeConn,
	keyFactory crypto.KeyFactory,
) HandShake {
	return &clientHandshake{
		id:         id,
		conn:       conn,
		buffer:     make([]byte, 1024),
		keyFactory: keyFactory,
	}
}

func (c *clientHandshake) Handshake() (aes.Key, error) {
	pubKey, err := c.pubKeyRequest(c.id)
	if err != nil {
		return nil, err
	}

	return c.shareSecretKey(pubKey)
}

func (c *clientHandshake) pubKeyRequest(id string) (rsa.PublicKey, error) {
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

	resp, err := c.read()
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

func (c *clientHandshake) shareSecretKey(pubKey rsa.PublicKey) (aes.Key, error) {
	secretKey, err := c.keyFactory.GenerateAesSecretKey()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	cipherSecretKey, err := pubKey.Encrypt([]byte(secretKey.ToString()))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	keyExchangeResponse := keyExchangeResponse{SecretKey: cipherSecretKey}

	request, err := json.Marshal(keyExchangeResponse)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug("sharing secret:", string(request))

	err = c.write(request)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return secretKey, nil
}

func (c *clientHandshake) write(msg []byte) error {
	_, err := c.conn.Write(msg)
	return err
}

func (c *clientHandshake) read() ([]byte, error) {
	numByte, err := c.conn.Read(c.buffer)
	if err != nil {
		return nil, err
	}
	return c.buffer[:numByte], nil
}
