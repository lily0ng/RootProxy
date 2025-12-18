package proxy

import (
	"context"
	"net"
	"time"
)

type TestResult struct {
	OK      bool
	Latency time.Duration
	Error   string
}

func TestConnectivity(ctx context.Context, p Proxy) TestResult {
	start := time.Now()
	d := net.Dialer{}
	c, err := d.DialContext(ctx, "tcp", p.Address())
	if err != nil {
		return TestResult{OK: false, Error: err.Error()}
	}
	_ = c.Close()
	return TestResult{OK: true, Latency: time.Since(start)}
}
