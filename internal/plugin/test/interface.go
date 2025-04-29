package test

import (
	"github.com/hashicorp/go-plugin"
)

// Plugin is the interface that all test plugins must implement
type Plugin interface {
	// Name returns the name of the plugin
	Name() string

	// Version returns the version of the plugin
	Version() string

	// Execute runs the plugin's main functionality
	Execute(args []string) (string, error)
}

// HandshakeConfig is the configuration for the handshake between the plugin and host
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "TEST_PLUGIN",
	MagicCookieValue: "test",
}
