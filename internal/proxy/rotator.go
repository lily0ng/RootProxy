package proxy

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/lily0ng/RootProxy/internal/config"
)

type Rotator struct {
	mu        sync.Mutex
	positions map[string]int
	rnd       *rand.Rand
}

func NewRotator() *Rotator {
	return &Rotator{
		positions: make(map[string]int),
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Rotator) Rotate(profileName string, chain []string, policy config.RotationPolicy, mgr *Manager) (string, error) {
	if mgr == nil {
		return "", errors.New("proxy manager required")
	}
	if !policy.Enabled || policy.Mode == config.RotationOff {
		return "", errors.New("rotation disabled")
	}
	if len(chain) == 0 {
		return "", errors.New("profile chain is empty")
	}

	// Only consider proxies that exist
	valid := make([]string, 0, len(chain))
	for _, name := range chain {
		if _, ok := mgr.GetByName(name); ok {
			valid = append(valid, name)
		}
	}
	if len(valid) == 0 {
		return "", errors.New("no valid proxies in chain")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	switch policy.Mode {
	case config.RotationRoundRobin:
		idx := r.positions[profileName] % len(valid)
		chosen := valid[idx]
		r.positions[profileName] = (idx + 1) % len(valid)
		if err := mgr.SetActive(chosen); err != nil {
			return "", err
		}
		return chosen, nil
	case config.RotationRandom:
		chosen := valid[r.rnd.Intn(len(valid))]
		if err := mgr.SetActive(chosen); err != nil {
			return "", err
		}
		return chosen, nil
	default:
		return "", errors.New("unsupported rotation mode")
	}
}
