package commands

import "fmt"

type PrivMsgCommand struct {
	Target  string
	Message string
}

func (c PrivMsgCommand) Name() string { return "PRIVMSG" }

func (c PrivMsgCommand) Arguments() []string {
	return []string{c.Target, c.Message}
}

func (c PrivMsgCommand) Validate() error {
	if c.Name() == "" || c.Name() != "PRIVMSG" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	if c.Target == "" || c.Message == "" {
		return fmt.Errorf("PRIVMSG command requires a target and a message")
	}
	return nil
}

func (c PrivMsgCommand) Execute(ctx *Ctx) error {
	response := "hun?"

	_, err := ctx.Connection.Write([]byte(response + "\r\n"))
	if err != nil {
		return fmt.Errorf("failed to send PRIVMSG response: %w", err)
	}
	return nil
}
