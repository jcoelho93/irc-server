package commands

import (
	"fmt"
	"log/slog"
)

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
	slog.Info("Received PING command", "connection", ctx.Connection.RemoteAddr().String())

	response := "PONG"
	_, err := ctx.Connection.Write([]byte(response + "\r\n"))
	if err != nil {
		return fmt.Errorf("failed to send PONG response: %w", err)
	}
	slog.Info("Sent PONG response", "connection", ctx.Connection.RemoteAddr().String())
	return nil
}
