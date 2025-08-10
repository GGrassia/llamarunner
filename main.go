package main

import (
	"fmt"
	"github/llamarunner/commands"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		commands.ShowHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "install":
		commands.InstallLlamaCpp()
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: llamarunner run <preset-name>")
			return
		}
		commands.RunWithPreset(os.Args[2])
	case "init":
		commands.InitPreset()
	case "list":
		commands.ListPresets()
	case "set":
		if len(os.Args) < 3 {
			fmt.Println("Usage: llamarunner set <target>")
			fmt.Println("Targets:")
			fmt.Println("  d    Set default settings")
			fmt.Println("  e    Edit settings file")
			return
		}
		target := os.Args[2]
		switch target {
		case "d":
			commands.SetDefaultSettings()
		case "e":
			commands.EditSettingsFile()
		default:
			fmt.Printf("Unknown target: %s\n", target)
			return
		}
	case "help", "-h", "--help":
		commands.ShowHelp()
	default:
		// If it's just a preset name, run it
		if len(os.Args) < 2 {
			fmt.Println("Usage: llamarunner run <preset-name>")
			return
		}
		commands.RunWithPreset(os.Args[2])
	}
}
