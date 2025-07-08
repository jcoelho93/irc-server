package commands

import "fmt"

type PingCommand struct{}

func (c PingCommand) Name() string { return "PING" }

func (c PingCommand) Arguments() []string {
	return []string{}
}

func (c PingCommand) Validate() error {
	if c.Name() == "" || c.Name() != "PING" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	return nil
}

func (c PingCommand) Execute(ctx *Ctx) error {
	fmt.Println("Received PING command")

	response := "PONG"
	_, err := ctx.Connection.Write([]byte(response + "\r\n"))
	if err != nil {
		return fmt.Errorf("failed to send PONG response: %w", err)
	}
	fmt.Println("Sent PONG response")
	return nil
}
