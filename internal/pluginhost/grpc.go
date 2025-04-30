package pluginhost

import (
	"context"
	"log"

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
	log.Printf("[GRPC] Registering gRPC server")
	pluginapi.RegisterPluginServer(s, &pluginServer{})
	return nil
}

// GRPCClient implements the gRPC client for the plugin
func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	log.Printf("[GRPC] Creating gRPC client")
	return &pluginClient{client: pluginapi.NewPluginClient(c)}, nil
}

// pluginServer implements the gRPC server interface
type pluginServer struct {
	pluginapi.UnimplementedPluginServer
}

// Name implements the Name RPC method
func (s *pluginServer) Name(ctx context.Context, req *pluginapi.Empty) (*pluginapi.NameResponse, error) {
	log.Printf("[GRPC] Name called")
	return &pluginapi.NameResponse{Name: "test-plugin"}, nil
}

// Version implements the Version RPC method
func (s *pluginServer) Version(ctx context.Context, req *pluginapi.Empty) (*pluginapi.VersionResponse, error) {
	log.Printf("[GRPC] Version called")
	return &pluginapi.VersionResponse{Version: "1.0.0"}, nil
}

// Execute implements the Execute RPC method
func (s *pluginServer) Execute(ctx context.Context, req *pluginapi.ExecuteRequest) (*pluginapi.ExecuteResponse, error) {
	log.Printf("[GRPC] Execute called with args: %v", req.Args)
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
	log.Printf("[GRPC] Client Name called")
	resp, err := c.client.Name(context.Background(), &pluginapi.Empty{})
	if err != nil {
		log.Printf("[GRPC] Name error: %v", err)
		return ""
	}
	return resp.Name
}

// Version implements the PluginInterface
func (c *pluginClient) Version() string {
	log.Printf("[GRPC] Client Version called")
	resp, err := c.client.Version(context.Background(), &pluginapi.Empty{})
	if err != nil {
		log.Printf("[GRPC] Version error: %v", err)
		return ""
	}
	return resp.Version
}

// Execute implements the PluginInterface
func (c *pluginClient) Execute(args []string) (string, error) {
	log.Printf("[GRPC] Client Execute called with args: %v", args)
	resp, err := c.client.Execute(context.Background(), &pluginapi.ExecuteRequest{
		Args: args,
	})
	if err != nil {
		log.Printf("[GRPC] Execute error: %v", err)
		return "", err
	}
	return resp.Result, nil
}
