package commands

import (
	"fmt"
	"github/llamarunner/utils"
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
	ForceCPU     bool   `toml:"force_cpu"`
}

const (
	SETTINGS_FILE = "$HOME/.llama-presets/settings.toml"
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
		SetDefaultSettings()
	case "e":
		EditSettingsFile()
	default:
		fmt.Printf("Unknown target: %s\n", target)
	}
}

func LoadSettings() (*Settings, error) {
	// Try user settings first
	userSettingsFile := getUserSettingsFile()
	settings, err := loadSettingsFromFile(userSettingsFile)
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

// getUserSettingsFile returns the path to the user's settings file
func getUserSettingsFile() string {
	homeDir := os.Getenv("HOME")
	return filepath.Join(homeDir, ".llama-presets", "settings.toml")
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
		ModelPath:    getUserSettingsFile(),
		ConfigPath:   filepath.Dir(getUserSettingsFile()),
		Host:         "localhost",
		Port:         "8080",
		ForceCPU:     false,
	}

	// Create settings directory if needed
	userDir := filepath.Dir(getUserSettingsFile())
	os.MkdirAll(userDir, 0755)

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

	// Use utils function to get the config directory
	configDir := utils.GetDefaultConfigDir()
	userFile := filepath.Join(configDir, "settings.toml")

	err = os.WriteFile(userFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func SetDefaultSettings() {
	settings := &Settings{
		LlamaCppPath: "/usr/local/llama.cpp",
		ModelPath:    filepath.Join(filepath.Dir(getUserSettingsFile()), "models"),
		ConfigPath:   filepath.Dir(getUserSettingsFile()),
		Host:         "localhost",
		Port:         "8080",
		ForceCPU:     false,
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

	// Get the actual expanded path
	actualPath := utils.GetDefaultConfigDir()
	settingsFile := filepath.Join(actualPath, "settings.toml")

	fmt.Printf("Editing settings file: %s\n", settingsFile)
	fmt.Println("Current settings:")
	fmt.Printf("LlamaCppPath: %s\n", settings.LlamaCppPath)
	fmt.Printf("ModelPath: %s\n", settings.ModelPath)
	fmt.Printf("ConfigPath: %s\n", settings.ConfigPath)
	fmt.Printf("Host: %s\n", settings.Host)
	fmt.Printf("Port: %s\n", settings.Port)
	fmt.Printf("ForceCPU: %t\n", settings.ForceCPU)

	// You can add logic here to actually open nano or other editor
	// For now, just show that it would edit the file
	fmt.Println("Note: This would normally open your editor to modify the file.")
	fmt.Println("You can manually edit: " + settingsFile)
}
