package commands

import (
	"fmt"
)

// HelpCommand implements the Command interface for showing help
type HelpCommand struct {
	*BaseCommand
}

// NewHelpCommand creates a new help command
func NewHelpCommand() *HelpCommand {
	return &HelpCommand{
		BaseCommand: NewBaseCommand(
			"help",
			"Show this help message",
			"llamarunner help",
		),
	}
}

// Run executes the help command
func (c *HelpCommand) Run(args []string) {
	fmt.Println("Available commands:")

	// Get all registered commands
	commands := GetAllCommands()

	for name, cmd := range commands {
		fmt.Printf("  %-12s %s\n", name, cmd.Description())
	}

	fmt.Println("\nUse 'llamarunner <command> --help' for more information about a command.")
}

// Register the help command automatically
func init() {
	RegisterCommand("help", NewHelpCommand())
}
