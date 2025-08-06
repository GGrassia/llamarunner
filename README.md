# llamarunner

A Linux CLI tool to manage the installation and launching of [llama.cpp](https://github.com/ggerganov/llama.cpp) with configurable presets.

## Features

- Download and build `llama.cpp` automatically with optimizations
- Manage multiple models with custom presets
- Define flags like context size, threads, quantization, and prompt files
- Quickly switch between models and configurations

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
