package plugin

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Plugin is the interface that all plugins must implement
type Plugin interface {
	// Name returns the name of the plugin
	Name() string

	// Version returns the version of the plugin
	Version() string

	// Execute runs the plugin's main functionality
	Execute(args []string) (string, error)
}

// GRPCPlugin is the interface that all GRPC plugins must implement
type GRPCPlugin interface {
	Plugin
	plugin.GRPCPlugin
}

// GRPCServer is the interface that all GRPC servers must implement
type GRPCServer interface {
	plugin.GRPCPlugin
}

// GRPCClient is the interface that all GRPC clients must implement
type GRPCClient interface {
	Plugin
}

// Client returns a new GRPC client for the given plugin
func Client(client *plugin.Client) (Plugin, error) {
	// Get the plugin
	raw, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Cast to our plugin type
	p, ok := raw.(Plugin)
	if !ok {
		return nil, err
	}

	return p, nil
}

// TestPluginGRPC is the GRPC implementation of the plugin
type TestPluginGRPC struct {
	Impl Plugin
}

func (p *TestPluginGRPC) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return p.Impl, nil
}

func (p *TestPluginGRPC) Server(b *plugin.MuxBroker) (interface{}, error) {
	return p.Impl, nil
}

func (p *TestPluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return nil
}

func (p *TestPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return p.Impl, nil
}

// ServeConfig is the configuration for serving plugins
type ServeConfig struct {
	HandshakeConfig plugin.HandshakeConfig
	Plugins         map[string]Plugin
}

// Serve serves the plugins
func Serve(config *ServeConfig) {
	pluginMap := make(map[string]plugin.Plugin)
	for name, p := range config.Plugins {
		pluginMap[name] = &TestPluginGRPC{Impl: p}
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: config.HandshakeConfig,
		Plugins:         pluginMap,
	})
}
