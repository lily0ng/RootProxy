package integrations

type Integration interface {
	Name() string
	HealthCheck() error
}
