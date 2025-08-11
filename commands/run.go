package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github/llamarunner/utils"
)

// RunCommand implements the Command interface for running presets
type RunCommand struct {
	*BaseCommand
}

// NewRunCommand creates a new run command
func NewRunCommand() *RunCommand {
	return &RunCommand{
		BaseCommand: NewBaseCommand(
			"run",
			"Load model with preset",
			"llamarunner run <preset-name>",
		),
	}
}

// Run executes the run command
func (c *RunCommand) Run(args []string) {
	if len(args) < 1 {
		fmt.Println(c.Usage())
		return
	}
	presetName := args[0]
	// Load the enhanced preset configuration
	preset, err := utils.LoadPresetConfig(presetName)
	if err != nil {
		fmt.Printf("Error loading preset config: %v\n", err)
		return
	}

	// The preset now contains the complete command with binary path, host, port and arguments
	// We need to split it into binary path and arguments
	commandLine := strings.TrimSpace(preset)
	parts := strings.Fields(commandLine)
	if len(parts) < 1 {
		fmt.Println("Error: Empty command line in preset config")
		return
	}

	// Extract the binary path (first part of the command line)
	binaryPath := parts[0]

	// Extract arguments (everything after the binary path)
	runArgs := parts[1:]

	// Build command with direct argument passing
	cmd := exec.Command(binaryPath, runArgs...)

	// Connect stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}

// Register the run command automatically
func init() {
	RegisterCommand("run", NewRunCommand())
}
