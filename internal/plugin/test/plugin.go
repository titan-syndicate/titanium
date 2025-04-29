package test

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-plugin"
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
	Impl *TestPlugin
}

func (p *TestPluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterPluginServer(s, &pluginServer{Impl: p.Impl})
	return nil
}

func (p *TestPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &testPluginClient{client: NewPluginClient(c)}, nil
}

// pluginServer is the gRPC server implementation
type pluginServer struct {
	Impl *TestPlugin
	UnimplementedPluginServer
}

func (s *pluginServer) Name(ctx context.Context, req *Empty) (*NameResponse, error) {
	return &NameResponse{Name: s.Impl.Name()}, nil
}

func (s *pluginServer) Version(ctx context.Context, req *Empty) (*VersionResponse, error) {
	return &VersionResponse{Version: s.Impl.Version()}, nil
}

func (s *pluginServer) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	result, err := s.Impl.Execute(req.Args)
	if err != nil {
		return nil, err
	}
	return &ExecuteResponse{Result: result}, nil
}

// testPluginClient is the gRPC client implementation
type testPluginClient struct {
	client PluginClient
}

func (c *testPluginClient) Name() string {
	resp, err := c.client.Name(context.Background(), &Empty{})
	if err != nil {
		return ""
	}
	return resp.Name
}

func (c *testPluginClient) Version() string {
	resp, err := c.client.Version(context.Background(), &Empty{})
	if err != nil {
		return ""
	}
	return resp.Version
}

func (c *testPluginClient) Execute(args []string) (string, error) {
	resp, err := c.client.Execute(context.Background(), &ExecuteRequest{Args: args})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"test": &TestPluginGRPC{Impl: &TestPlugin{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
