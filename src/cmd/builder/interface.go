package builder

import "zerotrust_chat/chat"

type Builder interface {
	NewServer(chat.ReceiveHandler) chat.Server
	NewClient(string, chat.ReceiveHandler) (chat.Client, error)
	GetSessionManager() chat.SessionManager
}
