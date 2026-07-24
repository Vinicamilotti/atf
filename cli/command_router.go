package cli

import "fmt"

type CliFunc interface {
	Exec(args ...[]string) error
}

type CommandRouter struct {
	Commands map[string]CliFunc
}

func (c *CommandRouter) AddCommand(command string, exec CliFunc) {
	c.Commands[command] = exec
}

func (c *CommandRouter) GetCommand(command string) (CliFunc, error) {
	exec, ok := c.Commands[command]

	if !ok {
		return nil, fmt.Errorf("Command not found")
	}

	return exec, nil
}
