package commands

import "fmt"

type PassCommand struct {
	Password string
}

func (c PassCommand) Name() string { return "PASS" }
func (c PassCommand) Arguments() []string {
	return []string{c.Password}
}

func (c PassCommand) Validate() error {
	if c.Name() == "" || c.Name() != "PASS" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	if len(c.Arguments()) != 1 {
		return fmt.Errorf("PASS command requires exactly one argument")
	}
	return nil
}

func (c PassCommand) Execute(ctx *Ctx) error {
	if ctx.Server.IsConnectionRegistered(ctx.Connection) {
		return fmt.Errorf("connection already registered, cannot set password")
	}
	if c.Password == "" {
		panic("password cannot be empty")
	}
	ctx.Server.SetPassword(ctx.Connection, c.Password)
	return nil
}
