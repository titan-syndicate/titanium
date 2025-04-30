package types

// Plugin is the interface that all plugins must implement
type Plugin interface {
	Name() string
	Version() string
	Execute(args []string) (string, error)
}
