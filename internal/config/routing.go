package config

import (
	"errors"
	"sort"
	"sync"
	"time"
)

type RoutingAction string

const (
	RouteDirect  RoutingAction = "direct"
	RouteProxy   RoutingAction = "proxy"
	RouteChain   RoutingAction = "chain"
	RouteProfile RoutingAction = "profile"
)

type MatchType string

const (
	MatchDomainGlob   MatchType = "domain_glob"
	MatchDomainSuffix MatchType = "domain_suffix"
	MatchCIDR         MatchType = "cidr"
)

type RoutingRule struct {
	ID        string
	Name      string
	Enabled   bool
	Priority  int
	Match     MatchType
	Pattern   string
	Action    RoutingAction
	Target    string
	UpdatedAt time.Time
}

type RoutingStore struct {
	mu   sync.RWMutex
	byID map[string]RoutingRule
}

func NewRoutingStore() *RoutingStore {
	return &RoutingStore{byID: make(map[string]RoutingRule)}
}

func (s *RoutingStore) Upsert(r RoutingRule) error {
	if r.ID == "" {
		return errors.New("routing rule id required")
	}
	if r.Name == "" {
		return errors.New("routing rule name required")
	}
	if r.Match == "" {
		return errors.New("routing rule match required")
	}
	if r.Pattern == "" {
		return errors.New("routing rule pattern required")
	}
	if r.Action == "" {
		return errors.New("routing rule action required")
	}
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = time.Now().UTC()
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[r.ID] = r
	return nil
}

func (s *RoutingStore) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id == "" {
		return errors.New("routing rule id required")
	}
	if _, ok := s.byID[id]; !ok {
		return errors.New("routing rule not found")
	}
	delete(s.byID, id)
	return nil
}

func (s *RoutingStore) List() []RoutingRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]RoutingRule, 0, len(s.byID))
	for _, r := range s.byID {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].Name < out[j].Name
		}
		return out[i].Priority < out[j].Priority
	})
	return out
}
