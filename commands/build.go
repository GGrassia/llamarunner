package commands

import (
	"fmt"
	"github/llamarunner/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BuildCommand implements the Command interface for building llama.cpp
type BuildCommand struct {
	*BaseCommand
}

// NewBuildCommand creates a new build command
func NewBuildCommand() *BuildCommand {
	return &BuildCommand{
		BaseCommand: NewBaseCommand(
			"build",
			"Builds llama.cpp with CUDA detection and optimizations",
			"llamarunner build [directory]",
		),
	}
}

// Run executes the build command
func (c *BuildCommand) Run(args []string) {
	var buildDir string

	if len(args) > 0 {
		buildDir = args[0]
	} else {
		// Default to llama.cpp directory from settings
		settings, err := LoadSettings()
		if err != nil || settings.LlamaCppPath == "" {
			fmt.Println("Error: no installation directory specified and no default found in settings")
			return
		}
		buildDir = filepath.Join(settings.LlamaCppPath, "llama.cpp")
	}

	c.buildLlamaCpp(buildDir)
}

// BuildLlamaCpp is an exported function that builds llama.cpp in the specified directory
func BuildLlamaCpp(buildDir string) error {
	fmt.Printf("Building llama.cpp in: %s\n", buildDir)

	// Check if the directory exists
	if _, err := os.Stat(buildDir); err != nil {
		return fmt.Errorf("directory %s does not exist", buildDir)
	}

	// Change to build directory
	err := os.Chdir(buildDir)
	if err != nil {
		return fmt.Errorf("error changing directory: %v", err)
	}

	// Load settings
	settings, err := LoadSettings()
	if err != nil {
		return fmt.Errorf("error loading settings: %v", err)
	}

	// Detect CUDA availability
	hasCUDA, err := utils.HasCUDA()
	if err != nil {
		return fmt.Errorf("error detecting CUDA: %v", err)
	}

	fmt.Printf("CUDA detection: available=%t, forceCPU=%t\n", hasCUDA, settings.ForceCPU)

	// Determine build configuration
	var shouldForceCPU bool

	if !hasCUDA && !settings.ForceCPU {
		// No CUDA and not forcing CPU - prompt user
		fmt.Print("CUDA not detected. Build with CPU only? (y/N) ")

		var input string
		fmt.Scanln(&input)

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" {
			shouldForceCPU = true

			// Ask if they want to update the forcing policy
			fmt.Print("Update settings to force CPU builds for future builds? (Y/n) ")
			fmt.Scanln(&input)

			input = strings.TrimSpace(strings.ToLower(input))
			if input != "n" && input != "" {
				settings.ForceCPU = true

				// Save updated settings
				err = SaveSettings(settings)
				if err != nil {
					return fmt.Errorf("error saving settings: %v", err)
				}

				fmt.Println("Settings updated to force CPU builds.")
			} else {
				fmt.Println("Proceeding with CPU build for this session only.")
			}

			fmt.Println("Building with CPU support only")
		} else {
			return fmt.Errorf("CUDA is required for full functionality. Use --force-cpu or set force_cpu=true in settings to build with CPU only")
		}
	} else if !hasCUDA || settings.ForceCPU {
		// No CUDA available or forcing CPU
		shouldForceCPU = true
		fmt.Println("Building with CPU support only")
	} else {
		// CUDA is available and not forcing CPU
		shouldForceCPU = false
		fmt.Println("Building with CUDA support")
	}

	// Build with cmake
	err = runCMakeBuild(shouldForceCPU)
	if err != nil {
		return fmt.Errorf("cmake build failed: %v", err)
	}

	// Copy binaries
	err = copyBinaries()
	if err != nil {
		return fmt.Errorf("error copying binaries: %v", err)
	}

	fmt.Println("llama.cpp built successfully!")
	return nil
}

// runCMakeBuild runs the cmake build process with appropriate flags
func runCMakeBuild(forceCPU bool) error {
	// Create build directory
	buildDir := "build"
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("error creating build directory: %v", err)
	}

	// Prepare cmake arguments
	cmakeArgs := []string{
		"-B", buildDir,
		"-DBUILD_SHARED_LIBS=OFF",
		"-DLLAMA_CURL=ON",
	}

	if forceCPU {
		cmakeArgs = append(cmakeArgs, "-DGGML_CUDA=OFF")
		fmt.Println("CMAKE args: -DGGML_CUDA=OFF")
	} else {
		cmakeArgs = append(cmakeArgs,
			"-DGGML_CUDA=ON",
			"-DGGML_CUDA_FA_ALL_QUANTS=ON",
		)
		fmt.Println("CMAKE args: -DGGML_CUDA=ON -DGGML_CUDA_FA_ALL_QUANTS=ON")
	}

	// Run cmake
	fmt.Println("Running cmake...")
	cmakeCmd := exec.Command("cmake", cmakeArgs...)
	cmakeCmd.Stdout = os.Stdout
	cmakeCmd.Stderr = os.Stderr

	err := cmakeCmd.Run()
	if err != nil {
		return fmt.Errorf("cmake configuration failed: %v", err)
	}

	// Build the targets
	fmt.Println("Building llama-cli and llama-gguf-split...")
	buildCmd := exec.Command("cmake",
		"--build", buildDir,
		"--config", "Release",
		"-j",
		"--clean-first",
		"--target", "llama-cli",
		"--target", "llama-gguf-split",
		"--target", "llama-server",
	)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	err = buildCmd.Run()
	if err != nil {
		return fmt.Errorf("cmake build failed: %v", err)
	}

	return nil
}

// copyBinaries copies the built binaries to the expected location
func copyBinaries() error {
	// Source directory (build output)
	srcDir := "build/bin"

	// Destination directory (parent of llama.cpp)
	destDir := ".."

	// Check if source directory exists
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("build output directory not found: %s", srcDir)
	}

	// Find all llama-* binaries
	binaries, err := filepath.Glob(filepath.Join(srcDir, "llama-*"))
	if err != nil {
		return fmt.Errorf("error finding binaries: %v", err)
	}

	if len(binaries) == 0 {
		return fmt.Errorf("no llama-* binaries found in %s", srcDir)
	}

	// Copy each binary
	for _, binary := range binaries {
		binaryName := filepath.Base(binary)
		destPath := filepath.Join(destDir, binaryName)

		fmt.Printf("Copying %s to %s\n", binaryName, destPath)

		err = os.Rename(binary, destPath)
		if err != nil {
			return fmt.Errorf("error copying %s: %v", binaryName, err)
		}
	}

	return nil
}

// buildLlamaCpp handles the building of llama.cpp (internal method)
func (c *BuildCommand) buildLlamaCpp(buildDir string) {
	err := BuildLlamaCpp(buildDir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Register the build command automatically
func init() {
	RegisterCommand("build", NewBuildCommand())
}
