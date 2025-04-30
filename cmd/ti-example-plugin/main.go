package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/go-plugin"
	"github.com/titan-syndicate/titanium/pkg/pluginapi"
	"google.golang.org/grpc"
)

// GRPCPlugin is the gRPC implementation of the plugin
type GRPCPlugin struct {
	plugin.NetRPCUnsupportedPlugin
}

// GRPCServer implements the gRPC server for the plugin
func (p *GRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	log.Printf("[PLUGIN] Registering gRPC server")
	pluginapi.RegisterPluginServer(s, &pluginServer{})
	return nil
}

// GRPCClient implements the gRPC client for the plugin
func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	log.Printf("[PLUGIN] Creating gRPC client")
	return &pluginClient{client: pluginapi.NewPluginClient(c)}, nil
}

// pluginServer implements the gRPC server interface
type pluginServer struct {
	pluginapi.UnimplementedPluginServer
}

// Name implements the Name RPC method
func (s *pluginServer) Name(ctx context.Context, req *pluginapi.Empty) (*pluginapi.NameResponse, error) {
	log.Printf("[PLUGIN] Name called")
	return &pluginapi.NameResponse{Name: "ti-example-plugin"}, nil
}

// Version implements the Version RPC method
func (s *pluginServer) Version(ctx context.Context, req *pluginapi.Empty) (*pluginapi.VersionResponse, error) {
	log.Printf("[PLUGIN] Version called")
	return &pluginapi.VersionResponse{Version: "1.0.0"}, nil
}

// Execute implements the Execute RPC method
func (s *pluginServer) Execute(ctx context.Context, req *pluginapi.ExecuteRequest) (*pluginapi.ExecuteResponse, error) {
	log.Printf("[PLUGIN] Execute called with args: %v", req.Args)
	return &pluginapi.ExecuteResponse{
		Result: "Plugin executed successfully",
	}, nil
}

// pluginClient implements the PluginInterface for the gRPC client
type pluginClient struct {
	client pluginapi.PluginClient
}

// Name implements the PluginInterface
func (c *pluginClient) Name() string {
	log.Printf("[PLUGIN] Client Name called")
	resp, err := c.client.Name(context.Background(), &pluginapi.Empty{})
	if err != nil {
		log.Printf("[PLUGIN] Name error: %v", err)
		return ""
	}
	return resp.Name
}

// Version implements the PluginInterface
func (c *pluginClient) Version() string {
	log.Printf("[PLUGIN] Client Version called")
	resp, err := c.client.Version(context.Background(), &pluginapi.Empty{})
	if err != nil {
		log.Printf("[PLUGIN] Version error: %v", err)
		return ""
	}
	return resp.Version
}

// Execute implements the PluginInterface
func (c *pluginClient) Execute(args []string) (string, error) {
	log.Printf("[PLUGIN] Client Execute called with args: %v", args)
	resp, err := c.client.Execute(context.Background(), &pluginapi.ExecuteRequest{
		Args: args,
	})
	if err != nil {
		log.Printf("[PLUGIN] Execute error: %v", err)
		return "", err
	}
	return resp.Result, nil
}

func main() {
	log.SetPrefix("[PLUGIN] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting plugin...")

	// Create the plugin map
	pluginMap := map[string]plugin.Plugin{
		"plugin": &GRPCPlugin{},
	}

	// Serve the plugin
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "TITANIUM_PLUGIN",
			MagicCookieValue: "titanium",
		},
		Plugins:    pluginMap,
		GRPCServer: plugin.DefaultGRPCServer,
	})

	os.Exit(0)
}
