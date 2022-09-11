package builder

import (
	"zerotrust_chat/chat"
	"zerotrust_chat/crypto"
)

// interface compliance
var _ Builder = builder{}

type builder struct {
	serverAddr       string
	sessionManager   chat.SessionManager
	cryptoKeyFactory crypto.KeyFactory
}

func NewBuilder(serverAddr string) Builder {
	return &builder{
		serverAddr:       serverAddr,
		sessionManager:   chat.NewSessionManager(),
		cryptoKeyFactory: crypto.NewKeyFactory(),
	}
}

func (b builder) NewServer() chat.Server {
	return chat.NewServer(b.serverAddr, b.sessionManager, b.cryptoKeyFactory)
}

func (b builder) NewClient(targetAddr string) (chat.Client, error) {
	return chat.NewClient(b.serverAddr, targetAddr, b.cryptoKeyFactory, b.sessionManager)
}

func (b builder) GetSessionManager() chat.SessionManager {
	return b.sessionManager
}
