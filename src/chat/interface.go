package chat

import "zerotrust_chat/crypto/aes"

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
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

type HandShake interface {
	Handshake() (aes.Key, error)
}
