package proxy

import (
	"errors"
	"sort"
	"sync"
)

type ChainStore struct {
	mu     sync.RWMutex
	byName map[string]Chain
}

func NewChainStore() *ChainStore {
	return &ChainStore{byName: make(map[string]Chain)}
}

func (s *ChainStore) Upsert(c Chain, maxHops int) error {
	if err := c.Validate(maxHops); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byName[c.Name] = c
	return nil
}

func (s *ChainStore) Remove(name string) error {
	if name == "" {
		return errors.New("chain name required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byName[name]; !ok {
		return errors.New("chain not found")
	}
	delete(s.byName, name)
	return nil
}

func (s *ChainStore) Get(name string) (Chain, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.byName[name]
	return c, ok
}

func (s *ChainStore) List() []Chain {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Chain, 0, len(s.byName))
	for _, c := range s.byName {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
