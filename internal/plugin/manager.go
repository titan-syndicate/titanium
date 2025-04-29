package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-plugin"
)

// HandshakeConfig is the configuration for plugin handshake
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "TITANIUM_PLUGIN",
	MagicCookieValue: "titanium",
}

// Manager handles plugin discovery and loading
type Manager struct {
	pluginDir string
	clients   map[string]*plugin.Client
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	homeDir, _ := os.UserHomeDir()
	return &Manager{
		pluginDir: filepath.Join(homeDir, ".titanium", "plugins"),
		clients:   make(map[string]*plugin.Client),
	}
}

// InstallPlugin installs a plugin from a local binary
func (m *Manager) InstallPlugin(pluginPath string) error {
	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
		return err
	}

	// Get plugin name from binary
	pluginName := filepath.Base(pluginPath)

	// Copy plugin to plugin directory
	destPath := filepath.Join(m.pluginDir, pluginName)
	return os.Rename(pluginPath, destPath)
}

// ListPlugins returns a list of installed plugins
func (m *Manager) ListPlugins() ([]string, error) {
	entries, err := os.ReadDir(m.pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var plugins []string
	for _, entry := range entries {
		if !entry.IsDir() {
			plugins = append(plugins, entry.Name())
		}
	}
	return plugins, nil
}

// GetPluginPath returns the full path to a plugin binary
func (m *Manager) GetPluginPath(name string) string {
	return filepath.Join(m.pluginDir, name)
}

// LoadPlugin loads a plugin by name
func (m *Manager) LoadPlugin(name string) (interface{}, error) {
	pluginPath := m.GetPluginPath(name)
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin not found: %s", name)
	}

	// Create a plugin client
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"test": &TestPluginGRPC{},
		},
		Cmd: exec.Command(pluginPath),
	})

	// Connect to the plugin
	rpcClient, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to plugin: %v", err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("test")
	if err != nil {
		return nil, fmt.Errorf("failed to dispense plugin: %v", err)
	}

	// Store the client for cleanup
	m.clients[name] = client

	return raw, nil
}

// Cleanup closes all plugin clients
func (m *Manager) Cleanup() {
	for _, client := range m.clients {
		client.Kill()
	}
}
