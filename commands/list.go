package commands

import (
	"fmt"
	"os"
	"strings"

	"github/llamarunner/utils"
)

func ListPresets() {
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
