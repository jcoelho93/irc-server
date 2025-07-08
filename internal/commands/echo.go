package commands

import (
	"fmt"
)

type EchoCommand struct {
	Message string
}

func (c EchoCommand) Name() string { return "ECHO" }
func (c EchoCommand) Arguments() []string {
	return []string{c.Message}
}
func (c EchoCommand) Validate() error {
	if c.Name() == "" || c.Name() != "ECHO" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	if len(c.Arguments()) != 1 {
		return fmt.Errorf("ECHO command requires one arguments")
	}
	return nil
}
func (c EchoCommand) Execute(ctx *Ctx) error {
	if c.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}
	fmt.Printf("Echoing message: %s\n", c.Message)
	_, err := ctx.Connection.Write([]byte(c.Message + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
