package main

import (
	"flag"
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/titan-syndicate/titanium/internal/cli"
	ti "github.com/titan-syndicate/titanium/internal/cli/cmd"
	"github.com/titan-syndicate/titanium/internal/server"
)

func main() {
	mode := flag.String("mode", "cli", "Run mode: 'server' or 'cli'")
	flag.Parse()

	// Initialize Docker client with specific API version
	dockerCli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.48"),
	)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer dockerCli.Close()

	switch *mode {
	case "server":
		runServer(dockerCli)
	case "cli":
		runCLI(dockerCli)
	default:
		println("Invalid mode. Use 'server' or 'cli'")
		os.Exit(1)
	}
}

func runServer(dockerCli *client.Client) {
	srv := server.New(dockerCli)
	if err := srv.Start(":8080"); err != nil {
		println("Error starting server:", err.Error())
		os.Exit(1)
	}
}

func runCLI(dockerCli *client.Client) {
	c := cli.New(dockerCli)
	ti.RegisterCommands(c)

	if err := c.Run(os.Args[1:]); err != nil {
		println("Error:", err.Error())
		c.PrintUsage()
		os.Exit(1)
	}
}
