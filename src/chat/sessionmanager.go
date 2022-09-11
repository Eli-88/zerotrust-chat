package chat

import (
	"errors"
	"sync"
	"zerotrust_chat/logger"
)

var _ SessionManager = &sessionManager{}

type sessionManager struct {
	mapping map[string]Session
	rwMux   sync.RWMutex
}

func NewSessionManager() SessionManager {
	return &sessionManager{
		mapping: make(map[string]Session),
		rwMux:   sync.RWMutex{},
	}
}

func (s *sessionManager) Add(session Session) bool {
	logger.Debug("added id:", session.GetId())
	logger.Trace()
	s.rwMux.Lock()
	defer s.rwMux.Unlock()

	id := session.GetId()
	_, ok := s.mapping[id]
	if !ok {
		s.mapping[id] = session
		return true
	} else {
		return false
	}
}

func (s *sessionManager) Remove(id string) {
	logger.Trace()
	s.rwMux.Lock()
	defer s.rwMux.Unlock()
	delete(s.mapping, id)
}

func (s *sessionManager) Write(id string, msg []byte) error {
	s.rwMux.RLock()
	defer s.rwMux.RUnlock()

	session, ok := s.mapping[id]
	if ok {
		logger.Debug("writing to id [", id, "]:", string(msg))
		return session.Write(msg)
	} else {
		return errors.New("id is not found in session manager")
	}
}

func (s *sessionManager) GetAllIds() []string {
	s.rwMux.RLock()
	defer s.rwMux.RUnlock()
	var result []string
	for id, _ := range s.mapping {
		result = append(result, id)
	}
	logger.Debug("all ids:", result)
	return result
}

func (s *sessionManager) ValidateId(id string) bool {
	s.rwMux.RLock()
	defer s.rwMux.RUnlock()
	_, ok := s.mapping[id]
	return ok
}
