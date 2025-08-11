package commands

// BaseCommand provides common functionality for all commands
type BaseCommand struct {
	name        string
	description string
	usage       string
}

// NewBaseCommand creates a new base command with metadata
func NewBaseCommand(name, description, usage string) *BaseCommand {
	return &BaseCommand{
		name:        name,
		description: description,
		usage:       usage,
	}
}

// Name returns the command name
func (c *BaseCommand) Name() string {
	return c.name
}

// Description returns the command description
func (c *BaseCommand) Description() string {
	return c.description
}

// Usage returns the command usage information
func (c *BaseCommand) Usage() string {
	return c.usage
}
