package plugin

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-plugin"
	"github.com/titan-syndicate/titanium/internal/plugin/test"
	"github.com/titan-syndicate/titanium/internal/plugin/types"
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
	log.Printf("[MANAGER] Installing plugin from %s", pluginPath)

	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
		log.Printf("[MANAGER] Failed to create plugin directory: %v", err)
		return err
	}

	// Get plugin name from binary
	pluginName := filepath.Base(pluginPath)
	log.Printf("[MANAGER] Plugin name: %s", pluginName)

	// Open source file
	src, err := os.Open(pluginPath)
	if err != nil {
		log.Printf("[MANAGER] Failed to open source file: %v", err)
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer src.Close()

	// Create destination file
	destPath := filepath.Join(m.pluginDir, pluginName)
	log.Printf("[MANAGER] Installing to %s", destPath)
	dest, err := os.Create(destPath)
	if err != nil {
		log.Printf("[MANAGER] Failed to create destination file: %v", err)
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dest.Close()

	// Copy the file
	if _, err := io.Copy(dest, src); err != nil {
		log.Printf("[MANAGER] Failed to copy file: %v", err)
		return fmt.Errorf("failed to copy file: %v", err)
	}

	// Make the plugin executable
	if err := os.Chmod(destPath, 0755); err != nil {
		log.Printf("[MANAGER] Failed to make plugin executable: %v", err)
		return fmt.Errorf("failed to make plugin executable: %v", err)
	}

	log.Printf("[MANAGER] Successfully installed plugin %s", pluginName)
	return nil
}

// ListPlugins returns a list of installed plugins
func (m *Manager) ListPlugins() ([]string, error) {
	log.Printf("[MANAGER] Listing plugins in %s", m.pluginDir)
	entries, err := os.ReadDir(m.pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[MANAGER] Plugin directory does not exist")
			return []string{}, nil
		}
		log.Printf("[MANAGER] Failed to read plugin directory: %v", err)
		return nil, err
	}

	var plugins []string
	for _, entry := range entries {
		if !entry.IsDir() {
			plugins = append(plugins, entry.Name())
		}
	}
	log.Printf("[MANAGER] Found plugins: %v", plugins)
	return plugins, nil
}

// GetPluginPath returns the full path to a plugin binary
func (m *Manager) GetPluginPath(name string) string {
	return filepath.Join(m.pluginDir, name)
}

// LoadPlugin loads a plugin by name
func (m *Manager) LoadPlugin(name string) (types.Plugin, error) {
	log.Printf("[MANAGER] Loading plugin %s", name)
	pluginPath := m.GetPluginPath(name)
	log.Printf("[MANAGER] Plugin path: %s", pluginPath)

	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		log.Printf("[MANAGER] Plugin not found: %s", name)
		return nil, fmt.Errorf("plugin not found: %s", name)
	}

	// Create a plugin client with gRPC support
	log.Printf("[MANAGER] Creating plugin client")
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"test": &test.TestPluginGRPC{Impl: &test.TestPlugin{}},
		},
		// Cmd: exec.Command("sh", "-c", pluginPath),
		Cmd:              exec.Command("sh", "-c", "./ti-example-plugin"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	// Connect to the plugin
	log.Printf("[MANAGER] Connecting to plugin")
	rpcClient, err := client.Client()
	if err != nil {
		log.Printf("[MANAGER] Failed to connect to plugin: %v", err)
		client.Kill()
		return nil, fmt.Errorf("failed to connect to plugin: %v", err)
	}

	// Request the plugin
	log.Printf("[MANAGER] Dispensing plugin")
	raw, err := rpcClient.Dispense("test")
	if err != nil {
		log.Printf("[MANAGER] Failed to dispense plugin: %v", err)
		client.Kill()
		return nil, fmt.Errorf("failed to dispense plugin: %v", err)
	}

	// Store the client for cleanup
	m.clients[name] = client

	// Cast to our plugin interface
	log.Printf("[MANAGER] Casting to plugin interface")
	plugin, ok := raw.(types.Plugin)
	if !ok {
		log.Printf("[MANAGER] Plugin does not implement required interface")
		client.Kill()
		return nil, fmt.Errorf("plugin does not implement the required interface")
	}

	log.Printf("[MANAGER] Successfully loaded plugin %s", name)
	return plugin, nil
}

// ExecutePlugin executes a plugin with the given arguments
func (m *Manager) ExecutePlugin(name string, args []string) (string, error) {
	log.Printf("[MANAGER] Executing plugin %s with args: %v", name, args)
	plugin, err := m.LoadPlugin(name)
	if err != nil {
		log.Printf("[MANAGER] Failed to load plugin: %v", err)
		return "", err
	}

	log.Printf("[MANAGER] Calling plugin.Execute()")
	result, err := plugin.Execute(args)
	if err != nil {
		log.Printf("[MANAGER] Plugin execution failed: %v", err)
		return "", err
	}

	log.Printf("[MANAGER] Plugin execution successful")
	return result, nil
}

// Cleanup closes all plugin clients
func (m *Manager) Cleanup() {
	log.Printf("[MANAGER] Cleaning up plugin clients")
	for _, client := range m.clients {
		client.Kill()
	}
}
