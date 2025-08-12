package commands

import (
	"fmt"
	"github/llamarunner/utils"
	"os"
	"os/exec"
)

// UpdateCommand implements the Command interface for updating llamarunner
type UpdateCommand struct {
	*BaseCommand
}

// NewUpdateCommand creates a new update command
func NewUpdateCommand() *UpdateCommand {
	return &UpdateCommand{
		BaseCommand: NewBaseCommand(
			"update",
			"Updates llamarunner to the latest version",
			"llamarunner update [options]\nOptions:\n  --check    Check for updates without installing\n  --force    Force update even if already latest",
		),
	}
}

// Run executes the update command
func (c *UpdateCommand) Run(args []string) {
	checkOnly := false
	forceUpdate := false

	// Parse arguments
	for _, arg := range args {
		switch arg {
		case "--check":
			checkOnly = true
		case "--force":
			forceUpdate = true
		case "-h", "--help":
			fmt.Println(c.Usage())
			return
		default:
			fmt.Printf("Unknown option: %s\n", arg)
			fmt.Println(c.Usage())
			return
		}
	}

	// Check for updates first
	release, err := utils.GetLatestGitHubRelease("GGrassia", "llamarunner")
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		return
	}

	// Load current settings to get version
	settings, err := utils.LoadSettings()
	if err != nil {
		fmt.Printf("Error loading settings: %v\n", err)
		return
	}

	currentVersion := settings.Version
	latestVersion := release.TagName

	fmt.Printf("Current version: %s\n", currentVersion)
	fmt.Printf("Latest version: %s\n", latestVersion)

	if currentVersion == latestVersion && !forceUpdate {
		fmt.Println("You are already using the latest version.")
		return
	}

	if checkOnly {
		fmt.Println("A new update is available!")
		fmt.Printf("Release: %s\n", release.Name)
		fmt.Printf("Changes:\n%s\n", release.Body)
		return
	}

	// Perform the update by re-running install.sh
	fmt.Println("Updating llamarunner...")
	err = c.runInstallScript()
	if err != nil {
		fmt.Printf("Error updating llamarunner: %v\n", err)
		return
	}

	// Update version in settings
	settings.Version = latestVersion
	err = utils.SaveSettings(settings)
	if err != nil {
		fmt.Printf("Warning: Could not update version in settings: %v\n", err)
	} else {
		fmt.Println("Version updated successfully!")
	}

	fmt.Println("Update completed successfully!")
}

// runInstallScript executes the install.sh script
func (c *UpdateCommand) runInstallScript() error {
	// Get the directory where llamarunner is installed
	scriptPath := "install.sh"

	// Check if install.sh exists in current directory
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("install.sh not found. Please run this command from the llamarunner installation directory")
	}

	// Make sure the script is executable
	err := os.Chmod(scriptPath, 0755)
	if err != nil {
		return fmt.Errorf("error making install.sh executable: %v", err)
	}

	// Run the install script
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("install script failed: %v", err)
	}

	return nil
}

// Register the update command automatically
func init() {
	RegisterCommand("update", NewUpdateCommand())
}
