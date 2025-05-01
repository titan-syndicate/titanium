package ti

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/titan-syndicate/titanium/internal/cli"
	pluginhost "github.com/titan-syndicate/titanium/internal/pluginhost"
)

// RegisterPluginCommands registers all plugin-related commands
func RegisterPluginCommands(cliApp *cli.CLI) {
	pluginManager := pluginhost.NewManager()

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
			{
				Name:        "exec",
				Description: "Execute a plugin",
				RunE: func(cmd *cobra.Command, args []string) error {
					if len(args) < 1 {
						return fmt.Errorf("plugin name is required")
					}
					pluginName := args[0]
					pluginArgs := args[1:]
					result, err := pluginManager.ExecutePlugin(pluginName, pluginArgs)
					if err != nil {
						return err
					}
					fmt.Println(result)
					return nil
				},
			},
		},
	})
}
