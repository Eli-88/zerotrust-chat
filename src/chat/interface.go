package chat

import "zerotrust_chat/crypto/aes"

//go:generate mockgen -destination=../test/mocks/mock_server.go -package=mocks zerotrust_chat/chat Server
//go:generate mockgen -destination=../test/mocks/mock_session.go -package=mocks zerotrust_chat/chat Session
//go:generate mockgen -destination=../test/mocks/mock_sessionmanager.go -package=mocks zerotrust_chat/chat SessionManager
//go:generate mockgen -destination=../test/mocks/mock_receivehandler.go -package=mocks zerotrust_chat/chat ReceiveHandler
//go:generate mockgen -destination=../test/mocks/mock_handshakeconn.go -package=mocks zerotrust_chat/chat HandshakeConn
//go:generate mockgen -destination=../test/mocks/mock_handshake.go -package=mocks zerotrust_chat/chat HandShake

type Server interface {
	Run() error
}

type Session interface {
	Write([]byte) error
	GetId() string
}

type SessionManager interface {
	Add(Session) bool
	Remove(string)
	Write(string, []byte) error
	GetAllIds() []string
	ValidateId(string) bool
}

type ReceiveHandler interface {
	OnReceive(string)
}

type HandshakeConn interface {
	Read(b []byte) ([]byte, error)
	Write(b []byte) (int, error)
}

type HandShake interface {
	Handshake() (aes.Key, error)
}
