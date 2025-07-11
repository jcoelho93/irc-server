package commands

import "fmt"

type CapCommand struct{}

func (c CapCommand) Name() string { return "CAP" }

func (c CapCommand) Arguments() []string {
	return []string{}
}

func (c CapCommand) Validate() error {
	if c.Name() == "" || c.Name() != "CAP" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	return nil
}

func (c CapCommand) Execute(ctx *Ctx) error {
	response := "CAP * ACK :multi-prefix"
	_, err := ctx.Connection.Write([]byte(response + "\r\n"))
	if err != nil {
		return fmt.Errorf("failed to send CAP response: %w", err)
	}
	return nil
}
