package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

func ExportJSON(proxies []Proxy) ([]byte, error) {
	return json.MarshalIndent(proxies, "", "  ")
}

func ImportJSON(r io.Reader) ([]Proxy, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return nil, errors.New("empty body")
	}

	var out []Proxy
	if err := json.Unmarshal(b, &out); err == nil {
		return out, nil
	}

	var one Proxy
	if err := json.Unmarshal(b, &one); err != nil {
		return nil, err
	}
	return []Proxy{one}, nil
}

func ExportText(proxies []Proxy) []byte {
	var b strings.Builder
	for _, p := range proxies {
		b.WriteString(FormatProxyLine(p))
		b.WriteString("\n")
	}
	return []byte(b.String())
}

func FormatProxyLine(p Proxy) string {
	scheme := string(p.Type)
	if scheme == "" {
		scheme = string(TypeHTTP)
	}

	userInfo := ""
	if p.Auth == AuthBasic && p.User != "" {
		userInfo = url.UserPassword(p.User, p.Pass).String() + "@"
	}

	name := ""
	if p.Name != "" {
		name = "#" + p.Name
	}

	return fmt.Sprintf("%s://%s%s:%d%s", scheme, userInfo, p.Host, p.Port, name)
}

func ImportText(r io.Reader) ([]Proxy, error) {
	s := bufio.NewScanner(r)
	out := make([]Proxy, 0)
	lineNo := 0
	for s.Scan() {
		lineNo++
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		p, err := ParseProxyLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		out = append(out, p)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func ParseProxyLine(line string) (Proxy, error) {
	line = strings.TrimSpace(line)

	// allow an inline JSON object per line
	if strings.HasPrefix(line, "{") {
		var p Proxy
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			return Proxy{}, err
		}
		if p.Type == "" {
			p.Type = TypeHTTP
		}
		if p.Auth == "" {
			p.Auth = AuthNone
		}
		return p, nil
	}

	// allow: name=MyProxy type=socks5 host=1.2.3.4 port=1080 user=u pass=p
	if strings.Contains(line, "=") && !strings.Contains(line, "://") {
		p, err := parseKeyValueProxy(line)
		if err != nil {
			return Proxy{}, err
		}
		if p.Type == "" {
			p.Type = TypeHTTP
		}
		if p.Auth == "" {
			p.Auth = AuthNone
		}
		return p, nil
	}

	// allow URL formats: scheme://user:pass@host:port#name
	if strings.Contains(line, "://") {
		u, err := url.Parse(line)
		if err != nil {
			return Proxy{}, err
		}
		pt, err := ParseType(u.Scheme)
		if err != nil {
			return Proxy{}, err
		}
		host := u.Hostname()
		portStr := u.Port()
		if host == "" || portStr == "" {
			return Proxy{}, errors.New("host:port required")
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return Proxy{}, errors.New("invalid port")
		}

		p := Proxy{Type: pt, Host: host, Port: port, Auth: AuthNone}
		if u.User != nil {
			p.User = u.User.Username()
			pw, _ := u.User.Password()
			p.Pass = pw
			if p.User != "" {
				p.Auth = AuthBasic
			}
		}

		frag := strings.TrimPrefix(u.Fragment, "#")
		if frag != "" {
			p.Name = frag
		}
		if p.Name == "" {
			p.Name = fmt.Sprintf("%s-%s", p.Type, p.Address())
		}
		return p, nil
	}

	// allow bare host:port (defaults to http)
	host, port, err := splitHostPort(line)
	if err != nil {
		return Proxy{}, err
	}
	p := Proxy{Type: TypeHTTP, Host: host, Port: port, Auth: AuthNone}
	p.Name = fmt.Sprintf("%s-%s", p.Type, p.Address())
	return p, nil
}

func ParseType(s string) (Type, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case string(TypeHTTP):
		return TypeHTTP, nil
	case string(TypeHTTPS):
		return TypeHTTPS, nil
	case string(TypeSOCKS4):
		return TypeSOCKS4, nil
	case string(TypeSOCKS5):
		return TypeSOCKS5, nil
	default:
		return "", fmt.Errorf("unsupported proxy type: %s", s)
	}
}

func splitHostPort(s string) (string, int, error) {
	s = strings.TrimSpace(s)
	idx := strings.LastIndex(s, ":")
	if idx <= 0 || idx >= len(s)-1 {
		return "", 0, errors.New("host:port required")
	}
	host := strings.TrimSpace(s[:idx])
	portStr := strings.TrimSpace(s[idx+1:])
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, errors.New("invalid port")
	}
	return host, port, nil
}

func parseKeyValueProxy(line string) (Proxy, error) {
	parts := strings.Fields(line)
	kv := make(map[string]string, len(parts))
	for _, p := range parts {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			continue
		}
		kv[strings.ToLower(strings.TrimSpace(k))] = strings.TrimSpace(v)
	}

	name := kv["name"]
	pt := kv["type"]
	host := kv["host"]
	portStr := kv["port"]
	user := kv["user"]
	pass := kv["pass"]

	var out Proxy
	out.Name = name
	if pt != "" {
		t, err := ParseType(pt)
		if err != nil {
			return Proxy{}, err
		}
		out.Type = t
	}
	out.Host = host
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return Proxy{}, errors.New("invalid port")
		}
		out.Port = port
	}
	out.User = user
	out.Pass = pass
	out.Auth = AuthNone
	if out.User != "" {
		out.Auth = AuthBasic
	}

	if out.Host == "" || out.Port <= 0 {
		// try to parse from addr=host:port
		if addr := kv["addr"]; addr != "" {
			h, port, err := splitHostPort(addr)
			if err != nil {
				return Proxy{}, err
			}
			out.Host = h
			out.Port = port
		}
	}
	if out.Host == "" || out.Port <= 0 {
		return Proxy{}, errors.New("proxy host/port required")
	}
	if out.Name == "" {
		out.Name = fmt.Sprintf("%s-%s", out.Type, out.Address())
	}
	return out, nil
}

func DetectImportFormat(body []byte, contentType string) string {
	ct := strings.ToLower(contentType)
	if strings.Contains(ct, "application/json") {
		return "json"
	}

	trim := bytes.TrimSpace(body)
	if len(trim) == 0 {
		return ""
	}
	if trim[0] == '{' || trim[0] == '[' {
		return "json"
	}
	return "text"
}
