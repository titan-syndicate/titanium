package test

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/titan-syndicate/titanium/internal/plugin/types"
	"google.golang.org/grpc"
)

// TestPlugin is a simple implementation of the Plugin interface
type TestPlugin struct{}

// Name returns the name of the plugin
func (p *TestPlugin) Name() string {
	log.Println("[PLUGIN] Name() called")
	return "test-plugin"
}

// Version returns the version of the plugin
func (p *TestPlugin) Version() string {
	log.Println("[PLUGIN] Version() called")
	return "v0.1.0"
}

// Execute runs the plugin's main functionality
func (p *TestPlugin) Execute(args []string) (string, error) {
	log.Printf("[PLUGIN] Execute() called with args: %v", args)
	return fmt.Sprintf("Hello from test plugin! Args: %s", strings.Join(args, ", ")), nil
}

// TestPluginGRPC is the GRPC implementation of the plugin
type TestPluginGRPC struct {
	plugin.NetRPCUnsupportedPlugin
	Impl *TestPlugin
}

func (p *TestPluginGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	log.Println("[PLUGIN] GRPCServer() called")
	RegisterPluginServer(s, &pluginServer{Impl: p.Impl})
	return nil
}

func (p *TestPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	log.Println("[PLUGIN] GRPCClient() called")
	return &testPluginClient{client: NewPluginClient(c)}, nil
}

// pluginServer is the gRPC server implementation
type pluginServer struct {
	Impl *TestPlugin
	UnimplementedPluginServer
}

func (s *pluginServer) Name(ctx context.Context, req *Empty) (*NameResponse, error) {
	log.Println("[PLUGIN] Name RPC called")
	return &NameResponse{Name: s.Impl.Name()}, nil
}

func (s *pluginServer) Version(ctx context.Context, req *Empty) (*VersionResponse, error) {
	log.Println("[PLUGIN] Version RPC called")
	return &VersionResponse{Version: s.Impl.Version()}, nil
}

func (s *pluginServer) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	log.Printf("[PLUGIN] Execute RPC called with args: %v", req.Args)
	result, err := s.Impl.Execute(req.Args)
	if err != nil {
		log.Printf("[PLUGIN] Execute RPC error: %v", err)
		return nil, err
	}
	return &ExecuteResponse{Result: result}, nil
}

// testPluginClient is the gRPC client implementation
type testPluginClient struct {
	client PluginClient
}

func (c *testPluginClient) Name() string {
	log.Println("[PLUGIN] Client Name() called")
	resp, err := c.client.Name(context.Background(), &Empty{})
	if err != nil {
		log.Printf("[PLUGIN] Client Name() error: %v", err)
		return ""
	}
	return resp.Name
}

func (c *testPluginClient) Version() string {
	log.Println("[PLUGIN] Client Version() called")
	resp, err := c.client.Version(context.Background(), &Empty{})
	if err != nil {
		log.Printf("[PLUGIN] Client Version() error: %v", err)
		return ""
	}
	return resp.Version
}

func (c *testPluginClient) Execute(args []string) (string, error) {
	log.Printf("[PLUGIN] Client Execute() called with args: %v", args)
	resp, err := c.client.Execute(context.Background(), &ExecuteRequest{Args: args})
	if err != nil {
		log.Printf("[PLUGIN] Client Execute() error: %v", err)
		return "", err
	}
	return resp.Result, nil
}

// Ensure testPluginClient implements the Plugin interface
var _ types.Plugin = (*testPluginClient)(nil)

func main() {
	log.SetPrefix("[PLUGIN] ")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Starting plugin...")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "TITANIUM_PLUGIN",
			MagicCookieValue: "titanium",
		},
		Plugins: map[string]plugin.Plugin{
			"test": &TestPluginGRPC{Impl: &TestPlugin{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
