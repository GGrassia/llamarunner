# llamarunner

A Linux CLI tool to manage the installation and launching of [llama.cpp](https://github.com/ggerganov/llama.cpp) with configurable presets.

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
