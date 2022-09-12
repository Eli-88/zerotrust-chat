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
	secretKey  aes.Key
}

func NewClientHandshake(
	id string,
	conn HandshakeConn,
	keyFactory crypto.KeyFactory,
	secretKey aes.Key,
) HandShake {
	return &clientHandshake{
		id:         id,
		conn:       conn,
		buffer:     make([]byte, 1024),
		keyFactory: keyFactory,
		secretKey:  secretKey,
	}
}

func (c *clientHandshake) Handshake() error {
	pubKey, err := c.pubKeyRequest(c.id)
	if err != nil {
		return err
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

func (c *clientHandshake) shareSecretKey(pubKey rsa.PublicKey) error {

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
