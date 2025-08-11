package commands

import (
	"fmt"
	"os"
	"strings"

	"github/llamarunner/utils"
)

// ListCommand implements the Command interface for listing presets
type ListCommand struct {
	*BaseCommand
}

// NewListCommand creates a new list command
func NewListCommand() *ListCommand {
	return &ListCommand{
		BaseCommand: NewBaseCommand(
			"list",
			"List available presets",
			"llamarunner list",
		),
	}
}

// Run executes the list command
func (c *ListCommand) Run(args []string) {
	// Find the config directory using the existing utility function
	configDir := utils.FindConfigDir()

	// Read all files in the config directory
	files, err := os.ReadDir(configDir)
	if err != nil {
		fmt.Printf("Error reading presets directory: %v\n", err)
		return
	}

	fmt.Println("Available presets:")

	// Filter and print preset files (ending with .cfg)
	presetCount := 0
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".cfg") {
			// Remove the .cfg extension for display
			presetName := strings.TrimSuffix(file.Name(), ".cfg")
			fmt.Printf("  - %s\n", presetName)
			presetCount++
		}
	}

	if presetCount == 0 {
		fmt.Println("  No presets found")
	} else {
		fmt.Printf("\nTotal presets: %d\n", presetCount)
	}
}

// Register the list command automatically
func init() {
	RegisterCommand("list", NewListCommand())
}
