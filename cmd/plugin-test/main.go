package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-plugin"
	"github.com/titan-syndicate/titanium/internal/plugin/test"
)

func main() {
	// Get the absolute path to the plugin binary
	pluginPath, err := filepath.Abs("bin/test-plugin")
	if err != nil {
		log.Fatalf("Failed to get plugin path: %v", err)
	}

	// Create a plugin client
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: test.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"test": &test.TestPluginGRPC{},
		},
		Cmd:              exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect to the plugin
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatalf("Failed to connect to plugin: %v", err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("test")
	if err != nil {
		log.Fatalf("Failed to dispense plugin: %v", err)
	}

	// Cast the plugin to our interface
	plugin := raw.(test.Plugin)

	// Use the plugin
	fmt.Printf("Plugin Name: %s\n", plugin.Name())
	fmt.Printf("Plugin Version: %s\n", plugin.Version())

	result, err := plugin.Execute(os.Args[1:])
	if err != nil {
		log.Fatalf("Plugin execution failed: %v", err)
	}

	fmt.Printf("Plugin Result: %s\n", result)
}
