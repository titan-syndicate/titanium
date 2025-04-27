package ti

import (
	"fmt"

	"github.com/titan-syndicate/titanium/internal/cli"
)

// RegisterBuildCommands registers all build-related commands
func RegisterBuildCommands(cliApp *cli.CLI) {
	cliApp.RegisterCommand(cli.Command{
		Name:        "build",
		Description: "Build an application using various build tools",
		Run:         runBuildCmd,
		Subcommands: []cli.Command{
			{
				Name:        "go",
				Description: "Build a Go application",
				Run: func(cli *cli.CLI, args []string) error {
					return RunGoreleaser(nil, args)
				},
			},
			{
				Name:        "pack",
				Description: "Build an application using Cloud Native Buildpacks",
				Run:         runPack,
			},
		},
	})
}

func runBuildCmd(cli *cli.CLI, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("path argument is required")
	}

	// If no subcommand is specified, use pack
	return runPack(cli, args)
}
