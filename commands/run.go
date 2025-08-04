package commands

import (
	"fmt"
	"github/llamarunner/utils"
	"os"
	"os/exec"
	"path/filepath"
)

func RunWithPreset() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: llamarunner run <preset>")
		return
	}

	presetName := os.Args[2]
	runWithPresetName(presetName)
}

func runWithPresetName(presetName string) {
	// Find llama.cpp binary and config directory
	llamaDir := utils.FindLlamaCppDir()
	configDir := utils.FindConfigDir()

	// Check if preset config exists
	configPath := filepath.Join(configDir, presetName+".cfg")
	if !utils.FileExists(configPath) {
		fmt.Printf("Error: Preset config not found: %s\n", configPath)
		return
	}

	// Load config
	config := utils.LoadPresetConfig(configPath)

	// Build command
	cmd := exec.Command(
		filepath.Join(llamaDir, "llama-cli"),
		"-m", config["model"],
		"-t", config["threads"],
		"-n", config["n_predict"],
		"--ctx_size", config["ctx_size"],
	)

	// Connect stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running llama.cpp: %v\n", err)
	}
}
