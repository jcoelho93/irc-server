package commands

import (
	"fmt"
)

type JoinCommand struct {
	Channels []string
}

func (c JoinCommand) Name() string { return "JOIN" }
func (c JoinCommand) Arguments() []string {
	return c.Channels
}

func (c JoinCommand) Validate() error {
	if c.Name() == "" || c.Name() != "JOIN" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	if len(c.Arguments()) < 1 {
		return fmt.Errorf("JOIN command requires at least one argument")
	}
	return nil
}

func (c JoinCommand) Execute(ctx *Ctx) error {
	return nil
}
