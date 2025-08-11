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
	fmt.Println("llamarunner - CLI manager for llama.cpp")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  llamarunner install    Downloads and builds llama.cpp with optimizations")
	fmt.Println("  llamarunner run <preset>    Load model with preset")
	fmt.Println("  llamarunner list         List available presets")
	fmt.Println("  llamarunner init             Initialize new preset")
	fmt.Println("  llamarunner set d          Reset llamarunner to default configuration")
	fmt.Println("  llamarunner set e          Edit configuration file")
	fmt.Println("  llamarunner help              Show this help")
	fmt.Println("  llamarunner -h                Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  llamarunner qwen3-flash")
	fmt.Println("  llamarunner run llama3-8b")
	fmt.Println("  llamarunner init")
	fmt.Println("  llamarunner set e")
	fmt.Println()
	fmt.Println("Preset files can be located in:")
	fmt.Println("  /usr/local/llama-presets/")
	fmt.Println("  ~/.llama-presets/")
}
