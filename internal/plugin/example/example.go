package main

import (
	"fmt"

	"github.com/titan-syndicate/titanium/internal/plugin"
)

// ExamplePlugin is a sample plugin implementation
type ExamplePlugin struct{}

// Name returns the name of the plugin
func (p *ExamplePlugin) Name() string {
	return "example"
}

// Version returns the version of the plugin
func (p *ExamplePlugin) Version() string {
	return "1.0.0"
}

// Description returns a description of the plugin
func (p *ExamplePlugin) Description() string {
	return "An example plugin for Titanium"
}

// Execute runs the plugin with the given arguments
func (p *ExamplePlugin) Execute(args []string) (string, error) {
	if len(args) == 0 {
		return "Hello from the example plugin!", nil
	}
	return fmt.Sprintf("Hello %s!", args[0]), nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"example": &ExamplePlugin{},
		},
	})
}
