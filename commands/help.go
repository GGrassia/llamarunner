package commands

import (
	"fmt"
)

func ShowHelp() {
	fmt.Println("llamarunner - CLI manager for llama.cpp")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  llamarunner install	Downloads and builds llama.cpp with optimizations")
	fmt.Println("  llamarunner <preset>    Shorthand to load model with preset")
	fmt.Println("  llamarunner run <preset>    Load model with preset")
	fmt.Println("  llamarunner init             Initialize new preset")
	fmt.Println("  llamarunner set d          Reset llamarunner to default configuration")
	fmt.Println("  llamarunner set e          Edit configuration file")
	fmt.Println("  llamarunner help              Show this help")
	fmt.Println("  llamarunner -h                Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  llamarunner qwen3-flash")
	fmt.Println("  llamarunner run llama3-8b")
	fmt.Println("  llamarunner init")
	fmt.Println("  llamarunner set e")
	fmt.Println()
	fmt.Println("Preset files should be located in /usr/local/llama-presets/")
	fmt.Println("or ~/.llama-presets/")
}
