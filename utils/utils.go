package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
)

const (
	DEFAULT_LLAMA_DIR = "/usr/local/share/llama.cpp"
)

// GetDefaultConfigDir returns the default config directory path
func GetDefaultConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".llama-presets")
}

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
	userConfigDir := GetDefaultConfigDir()
	if _, err := os.Stat(userConfigDir); err == nil {
		return userConfigDir
	}

	configDir := "/usr/local/share/llama-presets"

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

func LoadConfig() (string, string, error) {
	// Find the config directory
	configDir := FindConfigDir()

	// Construct the full path to config.toml
	configPath := filepath.Join(configDir, "config.toml")

	// Check if the file exists
	if !FileExists(configPath) {
		return "", "", fmt.Errorf("config.toml not found: %s", configPath)
	}

	// Read and parse the TOML file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", "", err
	}

	// Parse TOML
	tomlData, err := toml.LoadBytes(data)
	if err != nil {
		return "", "", err
	}

	// Extract host and port values with defaults
	host := tomlData.Get("host").(string)
	port := tomlData.Get("port").(string)

	return host, port, nil
}

func LoadPresetConfig(presetName string) (string, error) {
	// Find the config directory
	configDir := FindConfigDir()

	// Construct the full path to the preset config file
	configPath := filepath.Join(configDir, presetName+".cfg")

	// Check if the file exists
	if !FileExists(configPath) {
		return "", fmt.Errorf("preset config file not found: %s", configPath)
	}

	// Read the content of the cfg file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	// Parse host and port from config.toml
	host, port, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	// Get the llama.cpp directory
	llamaDir := FindLlamaCppDir()
	binaryPath := filepath.Join(llamaDir, "llama-server")

	// Read preset content and build enhanced command
	presetContent := strings.TrimSpace(string(data))

	// Build the enhanced command string
	// Format: ./bin/build/llama-server --host <host> --port <port> [preset arguments]
	var builder strings.Builder
	builder.WriteString(binaryPath)
	builder.WriteString(" --host ")
	builder.WriteString(host)
	builder.WriteString(" --port ")
	builder.WriteString(port)
	builder.WriteString(" ")
	builder.WriteString(presetContent)

	return builder.String(), nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
