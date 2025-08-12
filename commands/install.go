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
	var installDir string
	var input string
	defaultDir := utils.FindLlamaCppDir()
	// Logic is buggy, TODO fix user prompt for folder
	if defaultDir != "" {
		installDir = defaultDir
	} else {
		fmt.Printf("Enter installation directory [default: %s]: ", defaultDir)
		fmt.Scanln(&input)

		if input != "" {
			installDir = input
		}
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

	// Check if build flag is present
	buildFlag := false
	for _, arg := range args {
		if arg == "-b" || arg == "--build" {
			buildFlag = true
			break
		}
	}

	// Build llama.cpp only if -b or --build flag is present
	if buildFlag {
		buildCmd := NewBuildCommand()
		buildCmd.buildLlamaCpp(filepath.Join(installDir, "llama.cpp"))
	} else {
		fmt.Println("Skipping build. Use -b or --build flag to build llama.cpp after installation.")
	}

	// Save the installation path to settings
	settings, err := utils.LoadSettings()
	if err != nil {
		settings = &utils.Settings{}
	}

	settings.LlamaCppPath = installDir
	err = utils.SaveSettings(settings)
	if err != nil {
		fmt.Printf("Error saving settings: %v\n", err)
	} else {
		fmt.Println("llama.cpp installed and configured successfully!")
		fmt.Printf("Installation path: %s\n", installDir)
	}
	returnCmd := exec.Command("cd", "..")
	returnCmd.Stderr = os.Stderr
	returnCmd.Stdin = os.Stdin

	retErr := returnCmd.Run()
	if err != nil {
		fmt.Printf("No idea how but we could not return to the parent folder -> %v", retErr)
		return
	}
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Register the install command automatically
func init() {
	RegisterCommand("install", NewInstallCommand())
}
