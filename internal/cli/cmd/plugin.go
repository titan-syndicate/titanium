package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/titan-syndicate/titanium/internal/cli"
	pluginhost "github.com/titan-syndicate/titanium/internal/pluginhost"
)

// RegisterPluginCommands registers all plugin-related commands
func RegisterPluginCommands(cliApp *cli.CLI) {
	log.Printf("[PLUGIN] Starting plugin command registration")
	pluginManager := pluginhost.NewManager()

	// Register plugin management commands
	log.Printf("[PLUGIN] Registering plugin management commands")
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
	log.Printf("[PLUGIN] Getting installed plugin names")
	pluginNames, err := pluginManager.GetPluginNames()
	if err != nil {
		log.Printf("[PLUGIN] Error getting plugin names: %v", err)
		return
	}
	log.Printf("[PLUGIN] Found plugins: %v", pluginNames)

	for name, pluginPath := range pluginNames {
		log.Printf("[PLUGIN] Registering plugin command: %s (path: %s)", name, pluginPath)
		// Create a closure to capture the plugin name and path
		pluginName := name
		pluginPath := pluginPath

		cliApp.RegisterCommand(cli.Command{
			Name:               pluginName,
			Description:        fmt.Sprintf("Run the %s plugin", pluginName),
			DisableFlagParsing: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				log.Printf("[PLUGIN] Executing plugin %s with args: %v", pluginName, args)
				result, err := pluginManager.ExecutePlugin(pluginPath, args)
				if err != nil {
					log.Printf("[PLUGIN] Error executing plugin %s: %v", pluginName, err)
					return err
				}
				fmt.Println(result)
				return nil
			},
		})
		log.Printf("[PLUGIN] Successfully registered plugin command: %s", name)
	}
	log.Printf("[PLUGIN] Completed plugin command registration")
}
