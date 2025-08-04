package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Settings struct {
	LlamaCppPath string `toml:"llama_cpp_path"`
	ModelPath    string `toml:"model_path"`
	ConfigPath   string `toml:"config_path"`
	Host         string `toml:"host"`
	Port         string `toml:"port"`
}

const (
	SETTINGS_FILE      = "/usr/local/llama-presets/settings.toml"
	USER_SETTINGS_FILE = "/home/yourusername/.llama-presets/settings.toml"
)

func LoadSettings() (*Settings, error) {
	// Try user settings first
	settings, err := loadSettingsFromFile(USER_SETTINGS_FILE)
	if err == nil && settings != nil {
		return settings, nil
	}

	// Try system settings
	settings, err = loadSettingsFromFile(SETTINGS_FILE)
	if err == nil && settings != nil {
		return settings, nil
	}

	// Create default settings
	return createDefaultSettings()
}

func loadSettingsFromFile(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var settings Settings
	err = toml.Unmarshal(data, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func createDefaultSettings() (*Settings, error) {
	settings := &Settings{
		LlamaCppPath: "/usr/local/llama.cpp",
		ModelPath:    "/usr/local/llama-presets/models",
		ConfigPath:   "/usr/local/llama-presets",
		Host:         "localhost",
		Port:         "8080",
	}

	// Create settings directory if needed
	dir := filepath.Dir(SETTINGS_FILE)
	os.MkdirAll(dir, 0755)

	// Save default settings
	err := SaveSettings(settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func SaveSettings(settings *Settings) error {
	// Save to user settings
	data, err := toml.Marshal(settings)
	if err != nil {
		return err
	}

	err = os.WriteFile(USER_SETTINGS_FILE, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SetDefaultSettings() {
	settings := &Settings{
		LlamaCppPath: "/usr/local/llama.cpp",
		ModelPath:    "/usr/local/llama-presets/models",
		ConfigPath:   "/usr/local/llama-presets",
		Host:         "localhost",
		Port:         "8080",
	}

	err := SaveSettings(settings)
	if err != nil {
		fmt.Printf("Error setting default settings: %v\n", err)
	} else {
		fmt.Println("Default settings applied successfully!")
	}
}

func EditSettingsFile() {
	settings, err := LoadSettings()
	if err != nil {
		fmt.Printf("Error loading settings: %v\n", err)
		return
	}

	fmt.Printf("Editing settings file: %s\n", USER_SETTINGS_FILE)
	fmt.Println("Current settings:")
	fmt.Printf("LlamaCppPath: %s\n", settings.LlamaCppPath)
	fmt.Printf("ModelPath: %s\n", settings.ModelPath)
	fmt.Printf("ConfigPath: %s\n", settings.ConfigPath)
	fmt.Printf("Host: %s\n", settings.Host)
	fmt.Printf("Port: %s\n", settings.Port)

	// You can add logic here to actually open nano or other editor
	// For now, just show that it would edit the file
	fmt.Println("Note: This would normally open your editor to modify the file.")
	fmt.Println("You can manually edit: " + USER_SETTINGS_FILE)
}
