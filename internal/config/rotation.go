package config

import "time"

type RotationMode string

const (
	RotationOff        RotationMode = "off"
	RotationRoundRobin RotationMode = "round_robin"
	RotationRandom     RotationMode = "random"
)

type RotationPolicy struct {
	Enabled  bool
	Mode     RotationMode
	Interval time.Duration
}
