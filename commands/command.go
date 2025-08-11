package commands

// Command interface defines the standard interface for all commands
type Command interface {
	// Execute the command with given arguments
	Run(args []string)

	// Get the name of the command
	Name() string

	// Get a brief description of the command
	Description() string

	// Get detailed usage information
	Usage() string
}

// commandRegistry maps command names to their implementations
var commandRegistry = map[string]Command{}

// RegisterCommand registers a command with the registry
func RegisterCommand(name string, cmd Command) {
	commandRegistry[name] = cmd
}

// GetCommand retrieves a command by name from the registry
func GetCommand(name string) (Command, bool) {
	cmd, exists := commandRegistry[name]
	return cmd, exists
}
