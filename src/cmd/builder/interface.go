package builder

import "zerotrust_chat/chat"

type Builder interface {
	NewServer() chat.Server
	NewClient(string) (chat.Client, error)
	GetSessionManager() chat.SessionManager
}
