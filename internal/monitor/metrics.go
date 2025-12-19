package monitor

import (
	"sort"
	"sync"
	"time"

	"github.com/lily0ng/RootProxy/internal/proxy"
)

type ProxyMetrics struct {
	Name        string        `json:"name"`
	Successes   uint64        `json:"successes"`
	Failures    uint64        `json:"failures"`
	LastOK      bool          `json:"last_ok"`
	LastLatency time.Duration `json:"last_latency"`
	LastError   string        `json:"last_error"`
	LastTestAt  time.Time     `json:"last_test_at"`
}

type Store struct {
	mu      sync.RWMutex
	started time.Time
	byProxy map[string]*ProxyMetrics
}

func NewStore() *Store {
	return &Store{started: time.Now().UTC(), byProxy: make(map[string]*ProxyMetrics)}
}

func (s *Store) StartedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.started
}

func (s *Store) RecordTest(proxyName string, tr proxy.TestResult) {
	if proxyName == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.byProxy[proxyName]
	if !ok {
		m = &ProxyMetrics{Name: proxyName}
		s.byProxy[proxyName] = m
	}
	m.LastOK = tr.OK
	m.LastLatency = tr.Latency
	m.LastError = tr.Error
	m.LastTestAt = time.Now().UTC()
	if tr.OK {
		m.Successes++
	} else {
		m.Failures++
	}
}

func (s *Store) Snapshot() []ProxyMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]ProxyMetrics, 0, len(s.byProxy))
	for _, m := range s.byProxy {
		out = append(out, *m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
