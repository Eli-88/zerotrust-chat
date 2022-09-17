package builder

import (
	"zerotrust_chat/chat"
	"zerotrust_chat/crypto"
)

// interface compliance
var _ Builder = builder{}

type builder struct {
	serverAddr              string
	sessionManager          chat.SessionManager
	cryptoKeyFactory        crypto.KeyFactory
	chatReaderWriterFactory chat.ChatReaderWriterFactory
}

func NewBuilder(serverAddr string) Builder {
	return &builder{
		serverAddr:              serverAddr,
		sessionManager:          chat.NewSessionManager(),
		cryptoKeyFactory:        crypto.NewKeyFactory(),
		chatReaderWriterFactory: chat.NewChatReaderFactory(),
	}
}

func (b builder) NewServer(receiveHandler chat.ReceiveHandler) chat.Server {
	return chat.NewServer(b.serverAddr, b.sessionManager, b.cryptoKeyFactory, receiveHandler, b.chatReaderWriterFactory)
}

func (b builder) NewClient(targetAddr string, receiveHandler chat.ReceiveHandler) (chat.Session, error) {
	return chat.NewClient(b.serverAddr, targetAddr, b.cryptoKeyFactory, b.sessionManager, receiveHandler, b.chatReaderWriterFactory)
}

func (b builder) GetSessionManager() chat.SessionManager {
	return b.sessionManager
}
