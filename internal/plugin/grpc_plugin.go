package plugin

import (
	"context"

	"github.com/titan-syndicate/titanium/internal/plugin/test"
)

// GRPCPluginServer is the gRPC server implementation
type GRPCPluginServer struct {
	test.UnimplementedPluginServer
}

func (s *GRPCPluginServer) Name(ctx context.Context, req *test.Empty) (*test.NameResponse, error) {
	return &test.NameResponse{Name: "test-plugin"}, nil
}

func (s *GRPCPluginServer) Version(ctx context.Context, req *test.Empty) (*test.VersionResponse, error) {
	return &test.VersionResponse{Version: "v0.1.0"}, nil
}

func (s *GRPCPluginServer) Execute(ctx context.Context, req *test.ExecuteRequest) (*test.ExecuteResponse, error) {
	return &test.ExecuteResponse{Result: "Hello from test plugin!"}, nil
}

// GRPCPluginClient is the gRPC client implementation
type GRPCPluginClient struct {
	client test.PluginClient
}

func (c *GRPCPluginClient) Name() string {
	resp, err := c.client.Name(context.Background(), &test.Empty{})
	if err != nil {
		return ""
	}
	return resp.Name
}

func (c *GRPCPluginClient) Version() string {
	resp, err := c.client.Version(context.Background(), &test.Empty{})
	if err != nil {
		return ""
	}
	return resp.Version
}

func (c *GRPCPluginClient) Execute(args []string) (string, error) {
	resp, err := c.client.Execute(context.Background(), &test.ExecuteRequest{Args: args})
	if err != nil {
		return "", err
	}
	return resp.Result, nil
}
