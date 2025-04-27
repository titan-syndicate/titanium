package cli

import (
	"fmt"
)

// CLI represents the command-line interface
type CLI struct {
	commands map[string]Command
}

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Run         func(args []string) error
}

// New creates a new CLI instance
func New() *CLI {
	return &CLI{
		commands: make(map[string]Command),
	}
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

	return cmd.Run(args[1:])
}

// PrintUsage prints the usage information
func (c *CLI) PrintUsage() {
	fmt.Println("Usage: ti [command] [options]")
	fmt.Println("\nCommands:")
	for _, cmd := range c.commands {
		fmt.Printf("  %-15s %s\n", cmd.Name, cmd.Description)
	}
}
