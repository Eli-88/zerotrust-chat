package chat

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

type Client interface {
	Read() (string, error)
	Session
}
