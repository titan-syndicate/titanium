package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/titan-syndicate/titanium/internal/plugin/test"
	"google.golang.org/grpc"
)

// TestPlugin is a simple implementation of the Plugin interface
type TestPlugin struct{}

// Name returns the name of the plugin
func (p *TestPlugin) Name() string {
	return "test-plugin"
}

// Version returns the version of the plugin
func (p *TestPlugin) Version() string {
	return "v0.1.0"
}

// Execute runs the plugin's main functionality
func (p *TestPlugin) Execute(args []string) (string, error) {
	return fmt.Sprintf("Hello from test plugin! Args: %s", strings.Join(args, ", ")), nil
}

// TestPluginGRPC is the GRPC implementation of the plugin
type TestPluginGRPC struct {
	plugin.NetRPCUnsupportedPlugin
	Impl test.Plugin
}

func (p *TestPluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	test.RegisterPluginServer(s, &TestPluginServer{Impl: p.Impl})
	return nil
}

func (p *TestPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &TestPluginClient{client: test.NewPluginClient(c)}, nil
}

// TestPluginServer is the gRPC server implementation
type TestPluginServer struct {
	Impl test.Plugin
	test.UnimplementedPluginServer
}

func (s *TestPluginServer) Name(ctx context.Context, req *test.Empty) (*test.NameResponse, error) {
	return &test.NameResponse{Name: s.Impl.Name()}, nil
}

func (s *TestPluginServer) Version(ctx context.Context, req *test.Empty) (*test.VersionResponse, error) {
	return &test.VersionResponse{Version: s.Impl.Version()}, nil
}

func (s *TestPluginServer) Execute(ctx context.Context, req *test.ExecuteRequest) (*test.ExecuteResponse, error) {
	result, err := s.Impl.Execute(req.Args)
	if err != nil {
		return nil, err
	}
	return &test.ExecuteResponse{Result: result}, nil
}

// TestPluginClient is the gRPC client implementation
type TestPluginClient struct {
	client test.PluginClient
}

func (c *TestPluginClient) Name() string {
	resp, err := c.client.Name(context.Background(), &test.Empty{})
	if err != nil {
		return ""
	}
	return resp.Name
}

func (c *TestPluginClient) Version() string {
	resp, err := c.client.Version(context.Background(), &test.Empty{})
	if err != nil {
		return ""
	}
	return resp.Version
}

func (c *TestPluginClient) Execute(args []string) (string, error) {
	resp, err := c.client.Execute(context.Background(), &test.ExecuteRequest{Args: args})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: test.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"test": &TestPluginGRPC{Impl: &TestPlugin{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
