package proxy

import "errors"

type Chain struct {
	Name string
	Hops []string
}

func (c Chain) Validate(maxHops int) error {
	if c.Name == "" {
		return errors.New("chain name required")
	}
	if len(c.Hops) == 0 {
		return errors.New("chain must include at least one hop")
	}
	if len(c.Hops) > maxHops {
		return errors.New("chain exceeds max hops")
	}
	return nil
}
