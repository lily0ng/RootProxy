package config

import "sync"

type SecuritySettings struct {
	DoH            bool
	DoT            bool
	LeakProtection bool
	KillSwitch     bool
}

type SecurityStore struct {
	mu  sync.RWMutex
	cur SecuritySettings
}

func NewSecurityStore() *SecurityStore {
	return &SecurityStore{cur: SecuritySettings{}}
}

func (s *SecurityStore) Get() SecuritySettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cur
}

func (s *SecurityStore) Set(v SecuritySettings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cur = v
}
