package ti

import (
	"fmt"

	"github.com/titan-syndicate/titanium/internal/cli"
)

// RegisterCommands registers all Titanium CLI commands
func RegisterCommands(cliApp *cli.CLI) {
	cliApp.RegisterCommand(cli.Command{
		Name:        "pack",
		Description: "Pack a directory into an archive",
		Run:         runPack,
	})

	cliApp.RegisterCommand(cli.Command{
		Name:        "build",
		Description: "Build an archive",
		Run:         runBuild,
	})

	cliApp.RegisterCommand(cli.Command{
		Name:        "goreleaser",
		Description: "Run goreleaser",
		Run:         runGoreleaser,
	})
}

func runPack(args []string) error {
	// TODO: Implement pack command
	fmt.Println("Pack command not implemented yet")
	return nil
}

func runBuild(args []string) error {
	// TODO: Implement build command
	fmt.Println("Build command not implemented yet")
	return nil
}

func runGoreleaser(args []string) error {
	// TODO: Implement goreleaser command
	fmt.Println("Goreleaser command not implemented yet")
	return nil
}
