package chat

import (
	"encoding/json"
	"zerotrust_chat/crypto"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
	"zerotrust_chat/logger"
)

var _ HandShake = &serverHandshake{}

type serverHandshake struct {
	conn       HandshakeConn
	keyFactory crypto.KeyFactory
	buffer     []byte
}

func NewServerHandshake(conn HandshakeConn, keyFactory crypto.KeyFactory) HandShake {
	return &serverHandshake{
		conn:       conn,
		keyFactory: keyFactory,
		buffer:     make([]byte, 1024),
	}
}

func (s *serverHandshake) Handshake() (aes.Key, error) {
	logger.Trace()

	// extract the secret key and encrypt your reply before sending to client
	secretKey, priKey, err := s.keyExchangeRequest()
	if err != nil {
		return nil, err
	}
	return s.startCommRequest(secretKey, priKey)
}

func (s *serverHandshake) keyExchangeRequest() (string, rsa.PrivateKey, error) {
	// generate rsa key pair and send the public key to client
	priKey, err := s.keyFactory.GenerateRsaPrivateKey()
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}
	pubKey := priKey.GetPublicKey().ToString()
	keyRequest := keyExchangeRequest{
		PubKey: pubKey,
	}

	req, err := json.Marshal(keyRequest)
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}

	logger.Debug("server sending pub key:", string(req))
	err = s.write(req)
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}

	response, err := s.read()
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}

	logger.Debug("receiving secret:", string(response))

	// extract the secret key and store in memory
	keyResponse := keyExchangeResponse{}
	err = json.Unmarshal(response, &keyResponse)
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}

	return keyResponse.SecretKey, priKey, nil
}

func (s *serverHandshake) startCommRequest(secretKey string, priKey rsa.PrivateKey) (aes.Key, error) {
	decryptedSecretKey, err := priKey.Decrypt(secretKey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	key, err := s.keyFactory.ConstructAesSecretKey(string(decryptedSecretKey))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return key, nil
}

func (s *serverHandshake) write(msg []byte) error {
	_, err := s.conn.Write(msg)
	return err
}

func (s *serverHandshake) read() ([]byte, error) {
	n, err := s.conn.Read(s.buffer)
	if err != nil {
		return nil, err
	}
	return s.buffer[:n], nil
}
