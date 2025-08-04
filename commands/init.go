package commands

import (
	"fmt"
	"github/llamarunner/utils"
	"os"
	"path/filepath"
)

func InitPreset() {
	fmt.Println("Initializing new preset...")

	// Get preset name
	fmt.Print("Enter preset name: ")
	var presetName string
	fmt.Scanln(&presetName)

	if presetName == "" {
		fmt.Println("Invalid preset name")
		return
	}

	// Get model path
	fmt.Print("Enter model path: ")
	var modelPath string
	fmt.Scanln(&modelPath)

	// Get threads
	fmt.Print("Enter thread count (default 8): ")
	var threads string
	fmt.Scanln(&threads)
	if threads == "" {
		threads = "8"
	}

	// Get n_predict
	fmt.Print("Enter n_predict value (default 200): ")
	var nPredict string
	fmt.Scanln(&nPredict)
	if nPredict == "" {
		nPredict = "200"
	}

	// Get ctx_size
	fmt.Print("Enter context size (default 2048): ")
	var ctxSize string
	fmt.Scanln(&ctxSize)
	if ctxSize == "" {
		ctxSize = "2048"
	}

	// Create config file
	configDir := utils.FindConfigDir()
	configPath := filepath.Join(configDir, presetName+".cfg")

	configContent := fmt.Sprintf("model=%s\nthreads=%s\nn_predict=%s\nctx_size=%s\n",
		modelPath, threads, nPredict, ctxSize)

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("Error creating config file: %v\n", err)
	} else {
		fmt.Printf("Preset '%s' created successfully at %s\n",
			presetName, configPath)
	}
}
