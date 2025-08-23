# llamarunner

A Linux CLI tool to manage the installation and launching of [llama.cpp](https://github.com/ggerganov/llama.cpp) with configurable presets.

## Table of Contents

- [llamarunner](#llamarunner)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Commands](#commands)
    - [Preset Configuration](#preset-configuration)
    - [Settings Management](#settings-management)
  - [System Functionalities](#system-functionalities)
    - [✅ Working](#-working)
    - [⚠️ Known Limitations / Areas for Future Enhancement](#️-known-limitations--areas-for-future-enhancement)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
  - [Examples](#examples)
    - [Basic Preset Example](#basic-preset-example)
    - [Performance Optimization Examples](#performance-optimization-examples)
  - [Troubleshooting](#troubleshooting)
    - [Common Issues](#common-issues)
    - [Getting Help](#getting-help)
  - [Contributing](#contributing)
  - [License](#license)
  - [Uninstallation](#uninstallation)
  - [Foreseeable upgrades](#foreseeable-upgrades)

## Features

- **Automated Installation**: Download and build `llama.cpp` automatically with system-specific optimizations (CUDA support detection).
- **Preset Management**: Create, list, and manage multiple model configurations with custom settings.
- **Flexible Configuration**: Define flags like context size, threads, quantization, prompt files, and server parameters (host/port).
- **Quick Execution**: Run models directly by preset name or execute specific `llama.cpp` binaries.
- **Settings Management**: Configure default paths, force CPU builds, and manage global settings.
- **Self-Updating**: Update llamarunner itself to the latest version from GitHub releases.

## Installation

You can install llamarunner using the following command:

```bash
curl -fsSL https://raw.githubusercontent.com/GGrassia/llamarunner/main/install.sh | sh
```

This will download and install the appropriate binary for your system to `~/.local/bin/`.

Alternatively, you can build from source:

```bash
git clone https://github.com/GGrassia/llamarunner.git
cd llamarunner
go build -o llamarunner
```

## Usage

### Commands

- `help`: Show this help message.
- `install`: Downloads and builds llama.cpp with optimizations. Optionally build immediately with `-b` or `--build`.
- `build [directory]`: Builds llama.cpp in the specified directory (or default from settings) with CUDA detection.
- `init`: Initialize a new preset configuration interactively.
- `list`: List all available presets.
- `run <preset-name>`: Load and run a model using the specified preset. Can also directly execute `llama.cpp` binaries if the preset name matches a binary path.
- `set <target>`: Manage configuration settings.
  - `set d`: Set default settings (paths, host, port, etc.).
  - `set e`: Edit the settings file (shows current settings, manual edit required).
- `update [options]`: Updates llamarunner to the latest version from GitHub.
  - `--check`: Check for updates without installing.
  - `--force`: Force update even if already on the latest version.

### Preset Configuration

Presets are stored as `.cfg` files in `~/.llama-presets/`. Each preset file contains command-line arguments for `llama-server`.

Example `my-model.cfg`:
```
model=/path/to/model.gguf
threads=8
n_predict=200
ctx_size=2048
```

When run, llamarunner automatically enhances this with:
```
llama-server --host localhost --port 8080 --model /path/to/model.gguf -threads 8 -n_predict 200 -ctx_size 2048
```

### Settings Management

Global settings are stored in `~/.llama-presets/settings.toml` and include:
- `llama_cpp_path`: Default directory for llama.cpp installation.
- `model_path`: Default directory for model files.
- `config_path`: Directory for preset configurations.
- `host`: Default server host (default: "localhost").
- `port`: Default server port (default: "8080").
- `force_cpu`: Force CPU builds even if CUDA is available (default: false).
- `version`: Current llamarunner version.

## System Functionalities

### ✅ Working
- **Automated llama.cpp Installation**: Downloads and compiles `llama.cpp` from source with optimizations.
- **CUDA Detection**: Automatically detects CUDA availability and prompts for CPU-only build if CUDA is not found. Can force CPU builds via settings.
- **Preset Creation & Management**: Interactive `init` command to create presets, `list` command to view available presets.
- **Model Execution**: `run` command loads presets and executes `llama-server` with all specified parameters (model path, threads, context size, predictions, host, port).
- **Settings Persistence**: Saves and loads global settings (paths, build preferences) in `~/.llama-presets/settings.toml`.
- **Binary Management**: Builds `llama-cli`, `llama-gguf-split`, and `llama-server` binaries. Keeps them in the `build/bin` directory within the llama.cpp installation.
- **Self-Updating**: `update` command checks GitHub for new releases and re-runs the installation script to update llamarunner itself.
- **Command-Line Interface**: Supports `-h`/`--help` for individual commands, argument parsing, and direct execution of preset names.

### ⚠️ Known Limitations / Areas for Future Enhancement
- **Preset Editing**: The `set e` command currently displays current settings but does not open an editor for inline editing. Manual file editing is required.
- **Advanced Server Features**: While `llama-server` is executed, advanced server configurations (e.g., different API endpoints beyond basic host/port) would require manual preset editing or direct binary execution.
- **Model Management**: No built-in model downloading or management features beyond specifying paths in presets. Users must handle model file acquisition and placement.
- **Error Handling for Missing Binaries**: If `llama-server` (or other binaries) fail to build, the error message is generic. More specific feedback on build failures could be added.
- **Cross-Platform Testing**: While designed for Linux, broader platform compatibility (beyond the provided Linux binary names) would require additional testing and potentially conditional logic.
- **Configuration Validation**: Limited validation of preset configurations beyond file existence. Invalid parameter combinations might lead to runtime errors from `llama.cpp` itself.

## Getting Started

1. **Install llamarunner**:
   ```bash
   curl -fsSL https://raw.githubusercontent.com/GGrassia/llamarunner/main/install.sh | sh
   ```

2. **Install llama.cpp** (if not already present):
   ```bash
   llamarunner install
   ```
   Follow the prompts to set the installation directory and optionally build immediately.

3. **Create a Preset**:
   ```bash
   llamarunner init
   ```
   Enter a name, model path, and desired parameters (threads, context size, etc.).

4. **Run Your Model**:
   ```bash
   llamarunner run <your-preset-name>
   ```

5. **List Available Presets**:
   ```bash
   llamarunner list
   ```

6. **Update llamarunner**:
   ```bash
   llamarunner update
   ```

## Configuration

Default paths and settings are managed in `~/.llama-presets/settings.toml`. You can edit this file manually or use `llamarunner set d` to reset defaults.

The tool automatically detects CUDA availability during the build process. If CUDA is not found, it will prompt you to build with CPU support only and optionally update your settings to force CPU builds in the future.

## Examples

### Basic Preset Example
```
model=/path/to/model.gguf
threads=8
n_predict=200
ctx_size=2048
```

Note: Host and port are automatically added by llamarunner when running the preset. To change these values, modify your global settings with `llamarunner set d`.

### Performance Optimization Examples

For CPU-only systems:
```
model=/path/to/model.gguf
threads= -1
ctx_size=2048
```

For GPU-accelerated systems (with CUDA):
```
model=/path/to/model.gguf
ngl=35  # Number of layers to offload to GPU
threads=4  # Fewer threads when using GPU
ctx_size=4096
```

## Troubleshooting

### Common Issues

**CUDA not detected during build:**
- If CUDA is installed but not detected, ensure you have the NVIDIA drivers and CUDA toolkit properly installed.
- You can force a CPU-only build by running `llamarunner set d` and setting `force_cpu = true`.

**Model loading fails:**
- Verify the model path in your preset is correct.
- Ensure the model file exists and has read permissions.
- Check that the model format is compatible with llama.cpp.

**Port already in use:**
- Change the port in your preset or settings to a different value.
- Kill any process using the desired port (e.g., `sudo fuser -k 8080/tcp`).

### Getting Help

If you encounter issues not covered here:
1. Check the [llama.cpp repository](https://github.com/ggerganov/llama.cpp) for model-specific issues.
2. Open an issue in the [llamarunner GitHub repository](https://github.com/GGrassia/llamarunner/issues).
3. Include your system information, llamarunner version (`llamarunner set d`), and preset configuration when reporting issues.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Uninstallation

To uninstall llamarunner:
1. Remove the binary from your system:
   ```bash
   rm ~/.local/bin/llamarunner
   ```
2. Optionally, remove all configurations and presets:
   ```bash
   rm -rf ~/.llama-presets
   ```

## Foreseeable upgrades

- I'd like to make a preset repo with different model families and the relative optimal parameters, to further simplify the process of running or testing a new model.
- ik_llama implementation: while it *technically* works if you download ik_llama, build it with llama-server, set the default folder to the ik_llama one in the config.toml and then try to run models, it would be nice to implement a command to pull and build ik_llama just as normal llama.cpp does.
