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

	cmd.Run(os.Args[2:])
}

// initializeCommands registers all available commands
func initializeCommands() {
	// Register core commands
	commands.RegisterCommand("install", commands.NewInstallCommand())
	commands.RegisterCommand("run", commands.NewRunCommand())
	commands.RegisterCommand("init", commands.NewInitCommand())
	commands.RegisterCommand("list", commands.NewListCommand())
	commands.RegisterCommand("set", commands.NewSetCommand())
	commands.RegisterCommand("help", commands.NewHelpCommand())
}
