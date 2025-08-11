package main

import (
	"fmt"
	"os"

	"github/llamarunner/commands"
)

func main() {
	// Initialize commands
	initializeCommands()

	if len(os.Args) < 2 {
		cmd, exists := commands.GetCommand("help")
		if !exists {
			fmt.Println("Error: help command not found")
			return
		}
		cmd.Run(nil)
		return
	}

	// Check for -h flag
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		cmd, exists := commands.GetCommand("help")
		if !exists {
			fmt.Println("Error: help command not found")
			return
		}
		cmd.Run(nil)
		return
	}

	commandName := os.Args[1]
	cmd, exists := commands.GetCommand(commandName)

	if !exists {
		// If it's just a preset name, run it with the run command
		if len(os.Args) >= 2 {
			runCmd, exists := commands.GetCommand("run")
			if exists {
				runCmd.Run([]string{os.Args[1]})
				return
			}
		}
		fmt.Printf("Unknown command: %s\n", commandName)
		return
	}

	// Check if any argument is -h or --help
	for _, arg := range os.Args[2:] {
		if arg == "-h" || arg == "--help" {
			fmt.Printf("%s - %s\n", cmd.Name(), cmd.Description())
			fmt.Println("Usage: " + cmd.Usage())
			return
		}
	}

	cmd.Run(os.Args[2:])
}

// initializeCommands registers all available commands
func initializeCommands() {
	// Register core commands
	commands.RegisterCommand("install", commands.NewInstallCommand())
	commands.RegisterCommand("build", commands.NewBuildCommand())
	commands.RegisterCommand("run", commands.NewRunCommand())
	commands.RegisterCommand("init", commands.NewInitCommand())
	commands.RegisterCommand("list", commands.NewListCommand())
	commands.RegisterCommand("set", commands.NewSetCommand())
	commands.RegisterCommand("help", commands.NewHelpCommand())
}
