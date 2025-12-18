package config

import (
	"errors"
	"sort"
	"sync"
	"time"
)

type Profile struct {
	Name      string
	Chain     []string
	UpdatedAt time.Time
}

type ProfileStore struct {
	mu         sync.RWMutex
	activeName string
	byName     map[string]Profile
}

func NewProfileStore(defaultActive string) *ProfileStore {
	return &ProfileStore{
		activeName: defaultActive,
		byName:     make(map[string]Profile),
	}
}

func (s *ProfileStore) Upsert(p Profile) error {
	if p.Name == "" {
		return errors.New("profile name required")
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now().UTC()
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.byName[p.Name] = p
	if s.activeName == "" {
		s.activeName = p.Name
	}
	return nil
}

func (s *ProfileStore) List() []Profile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Profile, 0, len(s.byName))
	for _, p := range s.byName {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (s *ProfileStore) Active() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeName
}

func (s *ProfileStore) SetActive(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if name == "" {
		return errors.New("profile name required")
	}
	if _, ok := s.byName[name]; !ok {
		return errors.New("profile not found")
	}
	s.activeName = name
	return nil
}
