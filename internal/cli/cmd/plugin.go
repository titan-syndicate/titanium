package ti

import (
	"fmt"

	"github.com/titan-syndicate/titanium/internal/cli"
	"github.com/titan-syndicate/titanium/internal/plugin"
)

// RegisterPluginCommands registers all plugin-related commands
func RegisterPluginCommands(cliApp *cli.CLI) {
	pluginManager := plugin.NewManager()

	cliApp.RegisterCommand(cli.Command{
		Name:        "plugin",
		Description: "Manage plugins",
		Subcommands: []cli.Command{
			{
				Name:        "install",
				Description: "Install a plugin",
				Run: func(cli *cli.CLI, args []string) error {
					if len(args) == 0 {
						return fmt.Errorf("plugin path is required")
					}
					return pluginManager.InstallPlugin(args[0])
				},
			},
			{
				Name:        "list",
				Description: "List installed plugins",
				Run: func(cli *cli.CLI, args []string) error {
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
}
