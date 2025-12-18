package proxy

import (
	"errors"
	"sort"
	"sync"
)

type Manager struct {
	mu     sync.RWMutex
	byID   map[string]Proxy
	byName map[string]string
}

func NewManager() *Manager {
	return &Manager{
		byID:   make(map[string]Proxy),
		byName: make(map[string]string),
	}
}

func (m *Manager) Add(p Proxy) error {
	if p.Name == "" {
		return errors.New("proxy name required")
	}
	if p.Host == "" || p.Port <= 0 {
		return errors.New("proxy host/port required")
	}
	if p.Type == "" {
		return errors.New("proxy type required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.byName[p.Name]; exists {
		return errors.New("proxy name already exists")
	}
	if p.ID == "" {
		p.ID = NewID()
	}
	m.byID[p.ID] = p
	m.byName[p.Name] = p.ID
	return nil
}

func (m *Manager) Update(id string, p Proxy) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	old, ok := m.byID[id]
	if !ok {
		return errors.New("proxy not found")
	}
	p.ID = id
	if p.Name == "" {
		p.Name = old.Name
	}
	if p.Type == "" {
		p.Type = old.Type
	}
	if p.Host == "" {
		p.Host = old.Host
	}
	if p.Port == 0 {
		p.Port = old.Port
	}

	if p.Name != old.Name {
		if _, exists := m.byName[p.Name]; exists {
			return errors.New("proxy name already exists")
		}
		delete(m.byName, old.Name)
		m.byName[p.Name] = id
	}

	m.byID[id] = p
	return nil
}

func (m *Manager) Remove(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.byID[id]
	if !ok {
		return errors.New("proxy not found")
	}
	delete(m.byID, id)
	delete(m.byName, p.Name)
	return nil
}

func (m *Manager) GetByName(name string) (Proxy, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	id, ok := m.byName[name]
	if !ok {
		return Proxy{}, false
	}
	p, ok := m.byID[id]
	return p, ok
}

func (m *Manager) List() []Proxy {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Proxy, 0, len(m.byID))
	for _, p := range m.byID {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
