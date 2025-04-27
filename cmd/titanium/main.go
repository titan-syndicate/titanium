package main

import (
	"flag"
	"os"

	"github.com/titan-syndicate/titanium/internal/cli"
	ti "github.com/titan-syndicate/titanium/internal/cli/cmd"
	"github.com/titan-syndicate/titanium/internal/server"
)

func main() {
	mode := flag.String("mode", "cli", "Run mode: 'server' or 'cli'")
	flag.Parse()

	switch *mode {
	case "server":
		runServer()
	case "cli":
		runCLI()
	default:
		println("Invalid mode. Use 'server' or 'cli'")
		os.Exit(1)
	}
}

func runServer() {
	srv := server.New()
	if err := srv.Start(":8080"); err != nil {
		println("Error starting server:", err.Error())
		os.Exit(1)
	}
}

func runCLI() {
	c := cli.New()
	ti.RegisterCommands(c)

	if err := c.Run(os.Args[1:]); err != nil {
		println("Error:", err.Error())
		c.PrintUsage()
		os.Exit(1)
	}
}
