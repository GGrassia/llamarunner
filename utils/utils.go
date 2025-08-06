package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_CONFIG_DIR = "/usr/local/share/llama-presets"
	DEFAULT_LLAMA_DIR  = "/usr/local/share/llama.cpp"
)

func FindLlamaCppDir() string {
	// Try to find llama.cpp installation
	llamaDir := DEFAULT_LLAMA_DIR

	// Check if llama.cpp is in standard location
	if FileExists(filepath.Join(llamaDir, "llama-server")) {
		return llamaDir
	}

	// If not found, ask user
	fmt.Printf("llama.cpp not found in %s\n", llamaDir)
	fmt.Print("Enter llama.cpp installation path (or press Enter for default): ")

	var input string
	fmt.Scanln(&input)

	if input != "" {
		return input
	}

	return llamaDir // fallback
}

func FindConfigDir() string {
	// Try user config directory first
	userConfigDir := filepath.Join(os.Getenv("HOME"), ".llama-presets")
	if _, err := os.Stat(userConfigDir); err == nil {
		return userConfigDir
	}

	configDir := DEFAULT_CONFIG_DIR

	// If directory exists, use it
	if _, err := os.Stat(configDir); err == nil {
		return configDir
	}

	// Ask user or create default
	fmt.Printf("Config directory not found: %s\n", configDir)
	fmt.Print("Enter config directory path (or press Enter to create it): ")

	var input string
	fmt.Scanln(&input)

	if input != "" {
		return input
	}

	// Create the default directory
	os.MkdirAll(configDir, 0755)
	return configDir
}

func LoadPresetConfig(configPath string) map[string]string {
	config := make(map[string]string)

	// Simple config parser
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				config[parts[0]] = strings.TrimSpace(parts[1])
			}
		}
	}

	return config
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
