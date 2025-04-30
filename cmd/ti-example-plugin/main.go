package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/titan-syndicate/titanium/internal/plugin/test"
	"google.golang.org/grpc"
)

// TestPlugin is a simple implementation of the Plugin interface
type TestPlugin struct{}

// Name returns the name of the plugin
func (p *TestPlugin) Name() string {
	log.Println("[PLUGIN] Name() called")
	return "ti-example-plugin"
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
	test.RegisterPluginServer(s, &pluginServer{Impl: p.Impl})
	log.Println("[PLUGIN] Registered plugin server")
	return nil
}

func (p *TestPluginGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	log.Println("[PLUGIN] GRPCClient() called")
	client := &testPluginClient{client: test.NewPluginClient(c)}
	log.Println("[PLUGIN] Created plugin client")
	return client, nil
}

// pluginServer is the gRPC server implementation
type pluginServer struct {
	Impl *TestPlugin
	test.UnimplementedPluginServer
}

func (s *pluginServer) Name(ctx context.Context, req *test.Empty) (*test.NameResponse, error) {
	log.Println("[PLUGIN] Name RPC called")
	return &test.NameResponse{Name: s.Impl.Name()}, nil
}

func (s *pluginServer) Version(ctx context.Context, req *test.Empty) (*test.VersionResponse, error) {
	log.Println("[PLUGIN] Version RPC called")
	return &test.VersionResponse{Version: s.Impl.Version()}, nil
}

func (s *pluginServer) Execute(ctx context.Context, req *test.ExecuteRequest) (*test.ExecuteResponse, error) {
	log.Printf("[PLUGIN] Execute RPC called with args: %v", req.Args)
	result, err := s.Impl.Execute(req.Args)
	if err != nil {
		log.Printf("[PLUGIN] Execute RPC error: %v", err)
		return nil, err
	}
	return &test.ExecuteResponse{Result: result}, nil
}

// testPluginClient is the gRPC client implementation
type testPluginClient struct {
	client test.PluginClient
}

func (c *testPluginClient) Name() string {
	log.Println("[PLUGIN] Client Name() called")
	resp, err := c.client.Name(context.Background(), &test.Empty{})
	if err != nil {
		log.Printf("[PLUGIN] Client Name() error: %v", err)
		return ""
	}
	return resp.Name
}

func (c *testPluginClient) Version() string {
	log.Println("[PLUGIN] Client Version() called")
	resp, err := c.client.Version(context.Background(), &test.Empty{})
	if err != nil {
		log.Printf("[PLUGIN] Client Version() error: %v", err)
		return ""
	}
	return resp.Version
}

func (c *testPluginClient) Execute(args []string) (string, error) {
	log.Printf("[PLUGIN] Client Execute() called with args: %v", args)
	resp, err := c.client.Execute(context.Background(), &test.ExecuteRequest{Args: args})
	if err != nil {
		log.Printf("[PLUGIN] Client Execute() error: %v", err)
		return "", err
	}
	return resp.Result, nil
}

// Ensure testPluginClient implements the Plugin interface
var _ test.Plugin = (*testPluginClient)(nil)

func main() {
	log.SetPrefix("[PLUGIN] ")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Starting plugin...")

	// Create the plugin implementation
	testPlugin := &TestPlugin{}
	log.Println("[PLUGIN] Created plugin implementation")

	// Create the gRPC plugin
	grpcPlugin := &TestPluginGRPC{Impl: testPlugin}
	log.Println("[PLUGIN] Created gRPC plugin")

	// Serve the plugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "TITANIUM_PLUGIN",
			MagicCookieValue: "titanium",
		},
		Plugins: map[string]plugin.Plugin{
			"test": grpcPlugin,
		},
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			log.Println("[PLUGIN] Creating gRPC server")
			return grpc.NewServer(opts...)
		},
	})
}
