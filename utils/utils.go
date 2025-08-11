package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	// Check if llama.cpp is in build/bin directory
	if FileExists(filepath.Join(llamaDir, "build", "bin", "llama-server")) {
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

	// Construct the full path to settings.toml
	configPath := filepath.Join(configDir, "settings.toml")

	// Check if the file exists
	if !FileExists(configPath) {
		return "", "", fmt.Errorf("settings.toml not found: %s", configPath)
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

	// Parse host and port from settings.toml
	host, port, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	// Get the llama.cpp directory
	llamaDir := FindLlamaCppDir()
	binaryPath := filepath.Join(llamaDir, "build", "bin", "llama-server")

	// Read preset content and build enhanced command
	presetContent := strings.TrimSpace(string(data))

	// Build the enhanced command string
	// Format: llama-server --host <host> --port <port> [preset arguments]
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

// HasCUDA checks if CUDA is available on the system
func HasCUDA() (bool, error) {
	// Check if nvcc is available in PATH
	if !isCommandAvailable("nvcc") {
		return false, nil
	}

	// Additional CUDA checks could be added here if needed
	// For now, just checking for nvcc availability

	return true, nil
}

// isCommandAvailable checks if a command is available in the system PATH
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// DownloadFile downloads a file from a URL to a local path
func DownloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// GetLatestGitHubRelease fetches the latest release from GitHub
func GetLatestGitHubRelease(owner, repo string) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release GitHubRelease
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return nil, err
	}

	return &release, nil
}

// GetBinaryName returns the appropriate binary name for Linux
func GetBinaryName() string {
	goarch := runtime.GOARCH

	if goarch == "arm64" {
		return "llamarunner-linux-arm64"
	}
	return "llamarunner-linux-amd64"
}
