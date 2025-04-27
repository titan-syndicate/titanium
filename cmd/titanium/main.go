package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"
	"github.com/titan-syndicate/titanium/internal/cli"
	ti "github.com/titan-syndicate/titanium/internal/cli/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "titanium",
		Short: "Titanium CLI",
		Long:  `Titanium CLI for building and managing applications`,
	}
	cliInstance *cli.CLI
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build an application",
	Long:  `Build an application using various build tools`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		// Default to pack if no subcommand is specified
		return runPack(cmd, args)
	},
}

var buildGoCmd = &cobra.Command{
	Use:   "go",
	Short: "Build a Go application",
	Long:  `Build a Go application using goreleaser`,
	RunE:  ti.RunGoreleaser,
}

var buildPackCmd = &cobra.Command{
	Use:   "pack",
	Short: "Build using Cloud Native Buildpacks",
	Long:  `Build an application using Cloud Native Buildpacks`,
	RunE:  runPack,
}

func init() {
	// Add build subcommands
	buildCmd.AddCommand(buildGoCmd)
	buildCmd.AddCommand(buildPackCmd)

	// Add build command to root
	rootCmd.AddCommand(buildCmd)
}

func main() {
	// Initialize Docker client with specific API version
	dockerCli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.48"),
	)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer dockerCli.Close()

	// Create CLI instance
	cliInstance = cli.New(dockerCli)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runPack(cmd *cobra.Command, args []string) error {
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

	dockerCli := cliInstance.GetDockerClient()

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

	// Set up a channel to receive container logs
	logs, err := dockerCli.ContainerLogs(context.Background(), resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return fmt.Errorf("failed to get container logs: %v", err)
	}
	defer logs.Close()

	// Set up a channel to receive container completion
	statusCh, errCh := dockerCli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)

	// Stream logs in a goroutine
	go func() {
		stdcopy.StdCopy(os.Stdout, os.Stderr, logs)
	}()

	// Wait for container completion
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for container: %v", err)
		}
	case <-statusCh:
		// Container completed successfully
	}

	// Remove the container
	if err := dockerCli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{}); err != nil {
		log.Printf("Warning: Failed to remove container: %v", err)
	}

	return nil
}
