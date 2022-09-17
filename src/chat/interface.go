package chat

import "zerotrust_chat/crypto/aes"

//go:generate mockgen -destination=../test/mocks/mock_server.go -package=mocks zerotrust_chat/chat Server
//go:generate mockgen -destination=../test/mocks/mock_session.go -package=mocks zerotrust_chat/chat Session
//go:generate mockgen -destination=../test/mocks/mock_sessionmanager.go -package=mocks zerotrust_chat/chat SessionManager
//go:generate mockgen -destination=../test/mocks/mock_receivehandler.go -package=mocks zerotrust_chat/chat ReceiveHandler
//go:generate mockgen -destination=../test/mocks/mock_handshakeconn.go -package=mocks zerotrust_chat/chat HandshakeConn
//go:generate mockgen -destination=../test/mocks/mock_handshake.go -package=mocks zerotrust_chat/chat HandShake
//go:generate mockgen -destination=../test/mocks/mock_chat_reader.go -package=mocks zerotrust_chat/chat ChatReader

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
	OnReceive([]ChatMessage)
}

type Conn interface {
	Read(b []byte) ([]byte, error)
	Write(b []byte) (int, error)
}

type HandShake interface {
	Handshake() (aes.Key, error)
}

type ChatReaderWriter interface {
	Read([]byte) ([]ChatMessage, error)
	Write(msg []byte) error
}

type ChatReaderWriterFactory interface {
	Create(aes.Key, Conn) ChatReaderWriter
}
