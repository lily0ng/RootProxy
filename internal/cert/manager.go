package cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"sync"
)

type Certificate struct {
	Name string
	PEM  []byte
}

type Manager struct {
	mu     sync.RWMutex
	byName map[string]Certificate
}

func NewManager() *Manager {
	return &Manager{byName: make(map[string]Certificate)}
}

func (m *Manager) Add(name string, pemBytes []byte) error {
	if name == "" {
		return errors.New("certificate name required")
	}
	if len(pemBytes) == 0 {
		return errors.New("certificate PEM required")
	}
	if _, err := ParsePEM(pemBytes); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.byName[name] = Certificate{Name: name, PEM: pemBytes}
	return nil
}

func (m *Manager) List() []Certificate {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Certificate, 0, len(m.byName))
	for _, c := range m.byName {
		out = append(out, c)
	}
	return out
}

func ParsePEM(pemBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
