package pluginhost

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/titan-syndicate/titanium-plugin-sdk/pkg/logger"
)

// PluginInterface is the interface that all plugins must implement
type PluginInterface interface {
	Name() string
	Version() string
	Execute(args []string) (string, error)
}

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

// InstallPlugin installs a plugin from a local binary or GitHub release
func (m *Manager) InstallPlugin(pluginPath string) error {
	logger.Log.Info("[MANAGER] Installing plugin from %s", pluginPath)

	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
		logger.Log.Error("[MANAGER] Failed to create plugin directory: %v", err)
		return err
	}

	// Get plugin name from binary
	pluginName := filepath.Base(pluginPath)
	logger.Log.Info("[MANAGER] Plugin name: %s", pluginName)

	var src io.ReadCloser
	var err error

	// Check if this is a GitHub repository
	if strings.HasPrefix(pluginPath, "github.com/") {
		logger.Log.Info("[MANAGER] Detected GitHub repository: %s", pluginPath)
		src, err = m.downloadFromGitHub(pluginPath)
		if err != nil {
			logger.Log.Error("[MANAGER] Failed to download from GitHub: %v", err)
			return fmt.Errorf("failed to download from GitHub: %v", err)
		}
		defer src.Close()
	} else {
		// Open source file
		src, err = os.Open(pluginPath)
		if err != nil {
			logger.Log.Error("[MANAGER] Failed to open source file: %v", err)
			return fmt.Errorf("failed to open source file: %v", err)
		}
		defer src.Close()
	}

	// Create destination file
	destPath := filepath.Join(m.pluginDir, pluginName)
	logger.Log.Info("[MANAGER] Installing to %s", destPath)
	dest, err := os.Create(destPath)
	if err != nil {
		logger.Log.Error("[MANAGER] Failed to create destination file: %v", err)
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dest.Close()

	// Copy the file
	if _, err := io.Copy(dest, src); err != nil {
		logger.Log.Error("[MANAGER] Failed to copy file: %v", err)
		return fmt.Errorf("failed to copy file: %v", err)
	}

	// Make the plugin executable
	if err := os.Chmod(destPath, 0755); err != nil {
		logger.Log.Error("[MANAGER] Failed to make plugin executable: %v", err)
		return fmt.Errorf("failed to make plugin executable: %v", err)
	}

	logger.Log.Info("[MANAGER] Successfully installed plugin %s", pluginName)
	return nil
}

// UninstallPlugin removes a plugin by name
func (m *Manager) UninstallPlugin(name string) error {
	logger.Log.Info("[MANAGER] Uninstalling plugin %s", name)
	pluginPath := m.GetPluginPath(name)

	// Check if plugin exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		logger.Log.Warn("[MANAGER] Plugin not found: %s", name)
		return fmt.Errorf("plugin not found: %s", name)
	}

	// Remove the plugin
	if err := os.Remove(pluginPath); err != nil {
		logger.Log.Error("[MANAGER] Failed to remove plugin: %v", err)
		return fmt.Errorf("failed to remove plugin: %v", err)
	}

	logger.Log.Info("[MANAGER] Successfully uninstalled plugin %s", name)
	return nil
}

// UninstallAllPlugins removes all installed plugins
func (m *Manager) UninstallAllPlugins() error {
	logger.Log.Info("[MANAGER] Uninstalling all plugins")
	plugins, err := m.ListPlugins()
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		if err := m.UninstallPlugin(plugin); err != nil {
			logger.Log.Error("[MANAGER] Failed to uninstall plugin %s: %v", plugin, err)
			// Continue with other plugins even if one fails
		}
	}

	logger.Log.Info("[MANAGER] Successfully uninstalled all plugins")
	return nil
}

// downloadFromGitHub downloads the latest release from a GitHub repository
func (m *Manager) downloadFromGitHub(repo string) (io.ReadCloser, error) {
	// Split the repository path
	parts := strings.Split(repo, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid GitHub repository format: %s", repo)
	}
	owner, repo := parts[1], parts[2]

	// Get the latest release
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get latest release: %s", resp.Status)
	}

	// Parse the release response
	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release response: %v", err)
	}

	// Get plugin name from repo
	pluginName := repo

	// Determine the correct archive name for this platform
	var archiveName string
	switch {
	case runtime.GOOS == "darwin" && runtime.GOARCH == "arm64":
		archiveName = fmt.Sprintf("%s_%s_darwin_arm64.tar.gz", pluginName, strings.TrimPrefix(release.TagName, "v"))
	case runtime.GOOS == "darwin" && runtime.GOARCH == "amd64":
		archiveName = fmt.Sprintf("%s_%s_darwin_amd64.tar.gz", pluginName, strings.TrimPrefix(release.TagName, "v"))
	case runtime.GOOS == "linux" && runtime.GOARCH == "amd64":
		archiveName = fmt.Sprintf("%s_%s_linux_amd64.tar.gz", pluginName, strings.TrimPrefix(release.TagName, "v"))
	default:
		return nil, fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Find the archive asset
	var archiveURL string
	for _, asset := range release.Assets {
		if asset.Name == archiveName {
			archiveURL = asset.URL
			break
		}
	}
	if archiveURL == "" {
		return nil, fmt.Errorf("no archive found for platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Download the archive
	resp, err = http.Get(archiveURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download archive: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download archive: %s", resp.Status)
	}
	defer resp.Body.Close()

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "titanium-plugin-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save the archive to a temporary file
	archivePath := filepath.Join(tmpDir, archiveName)
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create archive file: %v", err)
	}
	if _, err := io.Copy(archiveFile, resp.Body); err != nil {
		archiveFile.Close()
		return nil, fmt.Errorf("failed to save archive: %v", err)
	}
	archiveFile.Close()

	// Extract the archive
	cmd := exec.Command("tar", "-xzf", archivePath, "-C", tmpDir)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to extract archive: %v", err)
	}

	// Open the extracted binary
	binaryPath := filepath.Join(tmpDir, pluginName)
	binaryFile, err := os.Open(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open extracted binary: %v", err)
	}

	return binaryFile, nil
}

// ListPlugins returns a list of installed plugins
func (m *Manager) ListPlugins() ([]string, error) {
	logger.Log.Info("[MANAGER] Listing plugins in %s", m.pluginDir)
	entries, err := os.ReadDir(m.pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Log.Info("[MANAGER] Plugin directory does not exist")
			return []string{}, nil
		}
		logger.Log.Error("[MANAGER] Failed to read plugin directory: %v", err)
		return nil, err
	}

	var plugins []string
	for _, entry := range entries {
		if !entry.IsDir() {
			plugins = append(plugins, entry.Name())
		}
	}
	logger.Log.Info("[MANAGER] Found plugins: %v", plugins)
	return plugins, nil
}

// GetPluginPath returns the full path to a plugin binary
func (m *Manager) GetPluginPath(name string) string {
	return filepath.Join(m.pluginDir, name)
}

// LoadPlugin loads a plugin by name
func (m *Manager) LoadPlugin(name string) (PluginInterface, error) {
	logger.Log.Info("[MANAGER] Loading plugin %s", name)
	pluginPath := m.GetPluginPath(name)
	logger.Log.Info("[MANAGER] Plugin path: %s", pluginPath)

	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		logger.Log.Info("[MANAGER] Plugin not found: %s", name)
		return nil, fmt.Errorf("plugin not found: %s", name)
	}

	// Create a plugin client with gRPC support
	logger.Log.Info("[MANAGER] Creating plugin client")
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"plugin": &GRPCPlugin{},
		},
		Cmd:              exec.Command("sh", "-c", pluginPath),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	// Connect to the plugin
	logger.Log.Info("[MANAGER] Connecting to plugin")
	rpcClient, err := client.Client()
	if err != nil {
		logger.Log.Error("[MANAGER] Failed to connect to plugin: %v", err)
		client.Kill()
		return nil, fmt.Errorf("failed to connect to plugin: %v", err)
	}

	// Request the plugin
	logger.Log.Info("[MANAGER] Dispensing plugin")
	raw, err := rpcClient.Dispense("plugin")
	if err != nil {
		logger.Log.Error("[MANAGER] Failed to dispense plugin: %v", err)
		client.Kill()
		return nil, fmt.Errorf("failed to dispense plugin: %v", err)
	}

	// Store the client for cleanup
	m.clients[name] = client

	// Cast to our plugin interface
	logger.Log.Info("[MANAGER] Casting to plugin interface")
	plugin, ok := raw.(PluginInterface)
	if !ok {
		logger.Log.Info("[MANAGER] Plugin does not implement required interface")
		client.Kill()
		return nil, fmt.Errorf("plugin does not implement the required interface")
	}

	logger.Log.Info("[MANAGER] Successfully loaded plugin %s", name)
	return plugin, nil
}

// ExecutePlugin executes a plugin with the given arguments
func (m *Manager) ExecutePlugin(name string, args []string) (string, error) {
	logger.Log.Info("[MANAGER] Executing plugin %s with args: %v", name, args)
	plugin, err := m.LoadPlugin(name)
	if err != nil {
		logger.Log.Error("[MANAGER] Failed to load plugin: %v", err)
		return "", err
	}

	logger.Log.Info("[MANAGER] Calling plugin.Execute()")
	result, err := plugin.Execute(args)
	if err != nil {
		logger.Log.Error("[MANAGER] Plugin execution failed: %v", err)
		return "", err
	}

	logger.Log.Info("[MANAGER] Plugin execution successful")
	return result, nil
}

// Cleanup closes all plugin clients
func (m *Manager) Cleanup() {
	logger.Log.Info("[MANAGER] Cleaning up plugin clients")
	for _, client := range m.clients {
		client.Kill()
	}
}

// GetPluginNames returns a map of plugin names to their paths
func (m *Manager) GetPluginNames() (map[string]string, error) {
	logger.Log.Info("[MANAGER] Getting plugin names from directory: %s", m.pluginDir)
	plugins, err := m.ListPlugins()
	if err != nil {
		logger.Log.Error("[MANAGER] Error listing plugins: %v", err)
		return nil, err
	}
	logger.Log.Info("[MANAGER] Found raw plugin list: %v", plugins)

	pluginMap := make(map[string]string)
	for _, plugin := range plugins {
		// Remove the 'ti-' prefix if it exists
		name := strings.TrimPrefix(plugin, "ti-")
		logger.Log.Info("[MANAGER] Mapping plugin %s to command name %s", plugin, name)
		pluginMap[name] = plugin
	}
	logger.Log.Info("[MANAGER] Final plugin map: %v", pluginMap)
	return pluginMap, nil
}
