package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/titan-syndicate/titanium-plugin-sdk/pkg/logger"
	"github.com/titan-syndicate/titanium/internal/cli"
	pluginhost "github.com/titan-syndicate/titanium/internal/pluginhost"
)

// RegisterPluginCommands registers all plugin-related commands
func RegisterPluginCommands(cliApp *cli.CLI) {
	logger.Log.Info("[PLUGIN] Starting plugin command registration")
	pluginManager := pluginhost.NewManager()

	// Register plugin management commands
	logger.Log.Info("[PLUGIN] Registering plugin management commands")
	cliApp.RegisterCommand(cli.Command{
		Name:        "plugin",
		Description: "Manage plugins",
		Subcommands: []cli.Command{
			{
				Name:        "install",
				Description: "Install a plugin",
				RunE: func(cmd *cobra.Command, args []string) error {
					if len(args) == 0 {
						return fmt.Errorf("plugin path is required")
					}
					return pluginManager.InstallPlugin(args[0])
				},
			},
			{
				Name:        "uninstall",
				Description: "Uninstall a plugin",
				RunE: func(cmd *cobra.Command, args []string) error {
					if len(args) == 0 {
						return fmt.Errorf("plugin name is required")
					}
					return pluginManager.UninstallPlugin(args[0])
				},
			},
			{
				Name:        "uninstall-all",
				Description: "Uninstall all plugins",
				RunE: func(cmd *cobra.Command, args []string) error {
					return pluginManager.UninstallAllPlugins()
				},
			},
			{
				Name:        "list",
				Description: "List installed plugins",
				RunE: func(cmd *cobra.Command, args []string) error {
					plugins, err := pluginManager.ListPlugins()
					if err != nil {
						return err
					}
					for _, plugin := range plugins {
						fmt.Println(plugin)
					}
					return nil
				},
			},
		},
	})

	// Register installed plugins as top-level commands
	logger.Log.Info("[PLUGIN] Getting installed plugin names")
	pluginNames, err := pluginManager.GetPluginNames()
	if err != nil {
		logger.Log.Error("[PLUGIN] Error getting plugin names: %v", err)
		return
	}
	logger.Log.Info("[PLUGIN] Found plugins: %v", pluginNames)

	for name, pluginPath := range pluginNames {
		logger.Log.Info("[PLUGIN] Registering plugin command: %s (path: %s)", name, pluginPath)
		// Create a closure to capture the plugin name and path
		pluginName := name
		pluginPath := pluginPath

		cliApp.RegisterCommand(cli.Command{
			Name:               pluginName,
			Description:        fmt.Sprintf("Run the %s plugin", pluginName),
			DisableFlagParsing: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				logger.Log.Info("[PLUGIN] Executing plugin %s with args: %v", pluginName, args)
				result, err := pluginManager.ExecutePlugin(pluginPath, args)
				if err != nil {
					logger.Log.Error("[PLUGIN] Error executing plugin %s: %v", pluginName, err)
					return err
				}
				fmt.Println(result)
				return nil
			},
		})
		logger.Log.Info("[PLUGIN] Successfully registered plugin command: %s", name)
	}
	logger.Log.Info("[PLUGIN] Completed plugin command registration")
}
