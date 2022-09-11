package chat

import (
	"encoding/json"
	"net"
	"zerotrust_chat/crypto"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/crypto/rsa"
	"zerotrust_chat/logger"
)

var _ Server = &server{}
var _ Session = &session{}

type server struct {
	ipAddr          string
	sessionManager  SessionManager
	keyFactory      crypto.KeyFactory
	receiverHandler ReceiveHandler
}

type session struct {
	id             string
	conn           net.Conn
	buffer         []byte
	keyFactory     crypto.KeyFactory
	secretKey      aes.Key
	sessionManager SessionManager
	receiveHandler ReceiveHandler
}

func (s session) GetId() string {
	return s.id
}

func NewServer(
	ipAddr string,
	sessionManager SessionManager,
	keyFactory crypto.KeyFactory,
	receiveHandler ReceiveHandler,
) Server {
	return &server{
		ipAddr:          ipAddr,
		sessionManager:  sessionManager,
		keyFactory:      keyFactory,
		receiverHandler: receiveHandler,
	}
}

func (s server) Run() error {
	listener, err := net.Listen("tcp", s.ipAddr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Debug("accept error:", err)
			continue
		}

		go s.onConnection(conn)
	}
}

func (s server) onConnection(conn net.Conn) {
	session := NewSession(
		conn,
		s.keyFactory,
		s.sessionManager,
		s.receiverHandler,
	)
	err := session.run() // blocking run, return when connection closes
	if err != nil {
		logger.Error(err)
	}
	logger.Debug("session ended")
	s.sessionManager.Remove(session.GetId()) // remove session once it returns
}

func NewSession(
	conn net.Conn,
	keyFactory crypto.KeyFactory,
	sessionManager SessionManager,
	receiveHandler ReceiveHandler,
) *session {
	return &session{
		id:             "",
		conn:           conn,
		buffer:         make([]byte, 1024),
		keyFactory:     keyFactory,
		secretKey:      nil,
		sessionManager: sessionManager,
		receiveHandler: receiveHandler,
	}
}

func (s *session) run() error {
	defer s.conn.Close()

	err := s.handshake()
	if err != nil {
		return err
	}
	logger.Debug("handshake successful!!!")

	for {
		numBytes, err := s.conn.Read(s.buffer)
		if err != nil {
			logger.Debug(err)
			break
		}

		msg := ChatMessage{}
		err = json.Unmarshal(s.buffer[:numBytes], &msg)
		if err != nil {
			logger.Debug(err)
			break
		}

		data, err := s.secretKey.Decrypt(msg.Data)
		if err != nil {
			logger.Debug(err)
			break
		}
		s.receiveHandler.OnReceive(string(data))
	}

	return nil
}

func (s *session) handshake() error {
	logger.Trace()

	// extract the secret key and encrypt your reply before sending to client
	secretKey, priKey, err := s.keyExchangeRequest()
	if err != nil {
		return err
	}

	err = s.startCommRequest(secretKey, priKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *session) internalWrite(msg []byte) error {
	_, err := s.conn.Write(msg)
	return err
}

func (s *session) Write(msg []byte) error {
	encryptedMsg, err := s.secretKey.Encrypt(msg)
	if err != nil {
		return err
	}

	chatMessage := ChatMessage{
		Data: encryptedMsg,
	}

	toBeSent, err := json.Marshal(chatMessage)
	if err != nil {
		return err
	}

	_, err = s.conn.Write(toBeSent)
	return err
}

func (s *session) read() ([]byte, error) {
	n, err := s.conn.Read(s.buffer)
	if err != nil {
		return nil, err
	}
	return s.buffer[:n], nil
}

func (s *session) keyExchangeRequest() (string, rsa.PrivateKey, error) {
	startConRequest, err := s.read()
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}

	// extract id and save them into session manager
	startConnectionRequest := startConnectionRequest{}
	err = json.Unmarshal(startConRequest, &startConnectionRequest)
	if err != nil {
		logger.Error(err)
		return "", nil, err
	}
	s.id = startConnectionRequest.Id
	s.sessionManager.Add(s)

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
	err = s.internalWrite(req)
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

func (s *session) startCommRequest(secretKey string, priKey rsa.PrivateKey) error {
	decryptedSecretKey, err := priKey.Decrypt(secretKey)
	if err != nil {
		logger.Error(err)
		return err
	}

	key, err := s.keyFactory.ConstructAesSecretKey(string(decryptedSecretKey))
	if err != nil {
		logger.Error(err)
		return err
	}
	s.secretKey = key
	return nil
}
