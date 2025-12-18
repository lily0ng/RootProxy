package proxy

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Type string

const (
	TypeHTTP   Type = "http"
	TypeHTTPS  Type = "https"
	TypeSOCKS4 Type = "socks4"
	TypeSOCKS5 Type = "socks5"
)

type AuthType string

const (
	AuthNone  AuthType = "none"
	AuthBasic AuthType = "basic"
)

type Proxy struct {
	ID   string
	Name string
	Type Type
	Host string
	Port int
	Auth AuthType
	User string
	Pass string
}

func (p Proxy) Address() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

func NewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
