package ti

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stdcopy"
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
		Description: "Build an application using Cloud Native Buildpacks",
		Run:         runBuild,
	})

	cliApp.RegisterCommand(cli.Command{
		Name:        "goreleaser",
		Description: "Run goreleaser",
		Run:         runGoreleaser,
	})
}

func runPack(cli *cli.CLI, args []string) error {
	// TODO: Implement pack command
	fmt.Println("Pack command not implemented yet")
	return nil
}

func runBuild(cli *cli.CLI, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("path argument is required")
	}

	sourcePath := args[0]

	// Validate path exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", sourcePath)
	}

	// Get absolute path of source directory
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	dockerCli := cli.GetDockerClient()

	// Pull the pack image if not available
	imageName := "buildpacksio/pack:latest"
	_, _, err = dockerCli.ImageInspectWithRaw(context.Background(), imageName)
	if err != nil {
		log.Printf("Pulling image %s...", imageName)
		out, err := dockerCli.ImagePull(context.Background(), imageName, image.PullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image: %v", err)
		}
		defer out.Close()
		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}

	// Prepare command arguments
	buildArgs := []string{"build", "test-app"}
	buildArgs = append(buildArgs, "--path", "/workspace")
	buildArgs = append(buildArgs, "--builder", "paketobuildpacks/builder-jammy-base")
	buildArgs = append(buildArgs, "--creation-time", "now")

	// Create container config
	config := &container.Config{
		Image: imageName,
		Cmd:   buildArgs,
		User:  "root", // Run as root to ensure access to Docker socket
	}

	// Create host config with volume mount
	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/workspace", absPath),
			"/var/run/docker.sock:/var/run/docker.sock",
		},
		// Ensure the container has access to the Docker socket
		SecurityOpt: []string{"label:disable"},
	}

	// Create the container
	resp, err := dockerCli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	if err := dockerCli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	// Wait for the container to finish
	statusCh, errCh := dockerCli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for container: %v", err)
		}
	case <-statusCh:
	}

	// Get the container logs
	out, err := dockerCli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return fmt.Errorf("failed to get container logs: %v", err)
	}

	// Copy the logs to stdout and stderr
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// Remove the container
	if err := dockerCli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{}); err != nil {
		log.Printf("Warning: Failed to remove container: %v", err)
	}

	return nil
}

func runGoreleaser(cli *cli.CLI, args []string) error {
	// TODO: Implement goreleaser command
	fmt.Println("Goreleaser command not implemented yet")
	return nil
}
