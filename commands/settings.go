package commands

import (
	"fmt"
	"github/llamarunner/utils"
)

// SetCommand implements the Command interface for settings management
type SetCommand struct {
	*BaseCommand
}

// NewSetCommand creates a new set command
func NewSetCommand() *SetCommand {
	return &SetCommand{
		BaseCommand: NewBaseCommand(
			"set",
			"Manage configuration settings",
			"llamarunner set <target>\nTargets:\n  d    Set default settings\n  e    Edit settings file",
		),
	}
}

// Run executes the set command
func (c *SetCommand) Run(args []string) {
	if len(args) < 1 {
		fmt.Println(c.Usage())
		return
	}
	target := args[0]
	switch target {
	case "d":
		utils.SetDefaultSettings()
	case "e":
		utils.EditSettingsFile()
	default:
		fmt.Printf("Unknown target: %s\n", target)
	}
}

// Register the set command automatically
func init() {
	RegisterCommand("set", NewSetCommand())
}
