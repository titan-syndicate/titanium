package pluginhost

import (
	"context"

	"github.com/titan-syndicate/titanium/pkg/pluginapi"
)

// GRPCPluginServer is the gRPC server implementation
type GRPCPluginServer struct {
	pluginapi.UnimplementedPluginServer
}

func (s *GRPCPluginServer) Name(ctx context.Context, req *pluginapi.Empty) (*pluginapi.NameResponse, error) {
	return &pluginapi.NameResponse{Name: "ti-example-plugin"}, nil
}

func (s *GRPCPluginServer) Version(ctx context.Context, req *pluginapi.Empty) (*pluginapi.VersionResponse, error) {
	return &pluginapi.VersionResponse{Version: "v0.1.0"}, nil
}

func (s *GRPCPluginServer) Execute(ctx context.Context, req *pluginapi.ExecuteRequest) (*pluginapi.ExecuteResponse, error) {
	return &pluginapi.ExecuteResponse{Result: "Hello from test plugin!"}, nil
}

// GRPCPluginClient is the gRPC client implementation
type GRPCPluginClient struct {
	client pluginapi.PluginClient
}

func (c *GRPCPluginClient) Name() string {
	resp, err := c.client.Name(context.Background(), &pluginapi.Empty{})
	if err != nil {
		return ""
	}
	return resp.Name
}

func (c *GRPCPluginClient) Version() string {
	resp, err := c.client.Version(context.Background(), &pluginapi.Empty{})
	if err != nil {
		return ""
	}
	return resp.Version
}

func (c *GRPCPluginClient) Execute(args []string) (string, error) {
	resp, err := c.client.Execute(context.Background(), &pluginapi.ExecuteRequest{Args: args})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}
