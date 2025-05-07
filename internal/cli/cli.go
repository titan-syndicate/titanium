package cli

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// Command represents a CLI command
type Command struct {
	Name               string
	Description        string
	Subcommands        []Command
	Run                func(*CLI, []string) error
	RunE               func(*cobra.Command, []string) error
	DisableFlagParsing bool
}

// CLI represents the command-line interface
type CLI struct {
	commands  map[string]Command
	dockerCli *client.Client
}

// NewCLI creates a new CLI instance
func NewCLI() *CLI {
	return &CLI{
		commands: make(map[string]Command),
	}
}

// GetCommand returns a command by name
func (c *CLI) GetCommand(name string) *cobra.Command {
	cmd, ok := c.commands[name]
	if !ok {
		return nil
	}

	cobraCmd := &cobra.Command{
		Use:                cmd.Name,
		Short:              cmd.Description,
		RunE:               cmd.RunE,
		DisableFlagParsing: cmd.DisableFlagParsing,
	}

	// Add subcommands
	for _, subcmd := range cmd.Subcommands {
		subcobraCmd := &cobra.Command{
			Use:                subcmd.Name,
			Short:              subcmd.Description,
			RunE:               subcmd.RunE,
			DisableFlagParsing: subcmd.DisableFlagParsing,
		}
		cobraCmd.AddCommand(subcobraCmd)
	}

	return cobraCmd
}

// RegisterCommand registers a new command
func (c *CLI) RegisterCommand(cmd Command) {
	c.commands[cmd.Name] = cmd
}

// Run runs the CLI with the given arguments
func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no command specified")
	}

	cmd, ok := c.commands[args[0]]
	if !ok {
		return fmt.Errorf("unknown command: %s", args[0])
	}

	// If there are subcommands and we have more args, check for subcommand
	if len(cmd.Subcommands) > 0 && len(args) > 1 {
		for _, subcmd := range cmd.Subcommands {
			if subcmd.Name == args[1] {
				return subcmd.Run(c, args[2:])
			}
		}
	}

	// If no subcommand matched or no subcommands exist, run the main command
	return cmd.Run(c, args[1:])
}

// PrintUsage prints the usage information
func (c *CLI) PrintUsage() {
	fmt.Println("Usage: ti [command] [options]")
	fmt.Println("\nCommands:")
	for _, cmd := range c.commands {
		fmt.Printf("  %-15s %s\n", cmd.Name, cmd.Description)
		if len(cmd.Subcommands) > 0 {
			for _, subcmd := range cmd.Subcommands {
				fmt.Printf("    %-13s %s\n", subcmd.Name, subcmd.Description)
			}
		}
	}
}

// GetDockerClient returns the Docker client instance
func (c *CLI) GetDockerClient() *client.Client {
	return c.dockerCli
}

// GetCommands returns a map of all registered commands
func (c *CLI) GetCommands() map[string]Command {
	return c.commands
}
