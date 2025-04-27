package ti

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/spf13/cobra"
	"github.com/titan-syndicate/titanium/internal/cli"
)

// RegisterGoreleaserCommands registers all goreleaser-related commands
func RegisterGoreleaserCommands(cliApp *cli.CLI) {
	cliApp.RegisterCommand(cli.Command{
		Name:        "goreleaser",
		Description: "Run goreleaser",
		Run: func(cli *cli.CLI, args []string) error {
			return RunGoreleaser(nil, args)
		},
	})
}

// RunGoreleaser runs the goreleaser command
func RunGoreleaser(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("version argument is required")
	}

	version := args[0]

	// Get path from args or use current directory
	path := "."
	if len(args) > 1 {
		path = args[1]
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Initialize Dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return fmt.Errorf("failed to connect to Dagger: %v", err)
	}
	defer client.Close()

	// Get the source directory
	sourceDir := client.Host().Directory(absPath)

	// Create a container with goreleaser
	container := client.Container().
		From("goreleaser/goreleaser:latest").
		WithDirectory("/src", sourceDir).
		WithWorkdir("/src").
		WithEnvVariable("GORELEASER_CURRENT_TAG", version).
		WithExec([]string{"goreleaser", "release", "--snapshot", "--skip", "docker,homebrew", "--verbose"})

	// Execute goreleaser
	output, err := container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to run goreleaser: %v", err)
	}

	fmt.Println(output)

	// Export the dist directory
	distDir := container.Directory("/src/dist")
	if _, err := distDir.Export(ctx, "./dist"); err != nil {
		return fmt.Errorf("failed to export dist directory: %v", err)
	}

	return nil
}
