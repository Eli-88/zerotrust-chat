package chat

import (
	"encoding/json"
	"net"
	"zerotrust_chat/crypto"
	"zerotrust_chat/crypto/aes"
	"zerotrust_chat/logger"
)

var _ Server = &server{}
var _ Session = &session{}

type server struct {
	ipAddr                  string
	sessionManager          SessionManager
	keyFactory              crypto.KeyFactory
	receiverHandler         ReceiveHandler
	chatReaderWriterFactory ChatReaderWriterFactory
}

type session struct {
	id                      string
	conn                    net.Conn
	buffer                  []byte
	keyFactory              crypto.KeyFactory
	secretKey               aes.Key
	sessionManager          SessionManager
	receiveHandler          ReceiveHandler
	chatReaderWriterFactory ChatReaderWriterFactory
	chatReaderWriter        ChatReaderWriter
}

func (s session) GetId() string {
	return s.id
}

func NewServer(
	ipAddr string,
	sessionManager SessionManager,
	keyFactory crypto.KeyFactory,
	receiveHandler ReceiveHandler,
	chatReaderWriterFactory ChatReaderWriterFactory,
) Server {
	return &server{
		ipAddr:                  ipAddr,
		sessionManager:          sessionManager,
		keyFactory:              keyFactory,
		receiverHandler:         receiveHandler,
		chatReaderWriterFactory: chatReaderWriterFactory,
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
		s.chatReaderWriterFactory,
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
	chatReaderWriterFactory ChatReaderWriterFactory,
) *session {
	return &session{
		id:                      "",
		conn:                    conn,
		buffer:                  make([]byte, 1024),
		keyFactory:              keyFactory,
		secretKey:               nil,
		sessionManager:          sessionManager,
		receiveHandler:          receiveHandler,
		chatReaderWriterFactory: chatReaderWriterFactory,
		chatReaderWriter:        nil,
	}
}

func (s *session) run() error {
	defer s.conn.Close()

	id, err := s.extractId()
	if err != nil {
		logger.Error(err)
		return err
	}

	s.id = id

	serverHandsake := NewServerHandshake(NewConn(s.conn), s.keyFactory)
	secretKey, err := serverHandsake.Handshake()
	if err != nil {
		return err
	}
	s.secretKey = secretKey

	s.sessionManager.Add(s)
	logger.Debug("handshake successful!!!")

	s.chatReaderWriter = s.chatReaderWriterFactory.Create(s.secretKey, NewConn(s.conn))
	for {
		numBytes, err := s.conn.Read(s.buffer)
		if err != nil {
			logger.Debug(err)
			break
		}
		msg, err := s.chatReaderWriter.Read(s.buffer[:numBytes])
		if err != nil {
			logger.Debug(err)
			break
		}
		s.receiveHandler.OnReceive(msg)
	}

	return nil
}

func (s *session) Write(msg []byte) error {
	return s.chatReaderWriter.Write(msg)
}

func (s *session) extractId() (string, error) {
	numBytes, err := s.conn.Read(s.buffer)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	startConRequest := s.buffer[:numBytes]

	// extract id and save them into session manager
	startConnectionRequest := startConnectionRequest{}
	err = json.Unmarshal(startConRequest, &startConnectionRequest)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return startConnectionRequest.Id, nil
}
