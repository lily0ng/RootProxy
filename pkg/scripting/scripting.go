package scripting

type Engine interface {
	Run(script string) error
}
