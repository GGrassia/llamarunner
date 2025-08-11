// install.go
package commands

import (
	"fmt"
	"github/llamarunner/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// InstallCommand implements the Command interface for installing llama.cpp
type InstallCommand struct {
	*BaseCommand
}

// NewInstallCommand creates a new install command
func NewInstallCommand() *InstallCommand {
	return &InstallCommand{
		BaseCommand: NewBaseCommand(
			"install",
			"Downloads and builds llama.cpp with optimizations",
			"llamarunner install",
		),
	}
}

// Run executes the install command
func (c *InstallCommand) Run(args []string) {
	fmt.Println("Installing llama.cpp...")

	// Check if git is available
	if !isCommandAvailable("git") {
		fmt.Println("Error: git is required but not found")
		return
	}

	// Default installation directory
	defaultDir := utils.FindLlamaCppDir()
	// Ask user for installation path
	fmt.Print("Enter installation directory [default: /usr/local/share/llama.cpp]: ")
	var input string
	fmt.Scanln(&input)

	installDir := defaultDir
	if input != "" {
		installDir = input
	}

	fmt.Printf("Installing to: %s\n", installDir)

	// Check if directory exists
	if _, err := os.Stat(installDir); err == nil {
		fmt.Printf("Directory %s already exists. Overwrite? (y/n): ", installDir)
		fmt.Scanln(&input)
		if strings.ToLower(input) != "y" {
			fmt.Println("Installation cancelled")
			return
		}
	}

	// Clone the repository
	fmt.Println("Cloning llama.cpp repository...")
	cloneCmd := exec.Command("git", "clone", "https://github.com/ggerganov/llama.cpp.git", installDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	err := cloneCmd.Run()
	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		return
	}

	// Build llama.cpp
	fmt.Println("Building llama.cpp...")
	buildDir := filepath.Join(installDir, "llama.cpp")

	// Change to build directory
	err = os.Chdir(buildDir)
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return
	}

	// Run make to build
	buildCmd := exec.Command("make")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	err = buildCmd.Run()
	if err != nil {
		fmt.Printf("Error building llama.cpp: %v\n", err)
		return
	}

	// Save the installation path to settings
	settings, err := LoadSettings()
	if err != nil {
		settings = &Settings{}
	}

	settings.LlamaCppPath = installDir
	err = SaveSettings(settings)
	if err != nil {
		fmt.Printf("Error saving settings: %v\n", err)
	} else {
		fmt.Println("llama.cpp installed and configured successfully!")
		fmt.Printf("Installation path: %s\n", installDir)
	}
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
