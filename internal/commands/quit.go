package commands

import (
	"fmt"
	"log/slog"
)

type QuitCommand struct{}

func (c QuitCommand) Name() string { return "QUIT" }

func (c QuitCommand) Arguments() []string {
	return []string{}
}

func (c QuitCommand) Validate() error {
	if c.Name() == "" || c.Name() != "QUIT" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	return nil
}

func (c QuitCommand) Execute(ctx *Ctx) error {
	slog.Info("Removing client connection from server", "remote_addr", ctx.Connection.RemoteAddr().String())

	// Remove the client from the server's client map
	if _, exists := ctx.Server.GetClient(ctx.Connection); exists {
		delete(ctx.Server.GetClients(), ctx.Connection)
		slog.Info("Client connection removed", "remote_addr", ctx.Connection.RemoteAddr().String())
	} else {
		slog.Warn("Client connection not found", "remote_addr", ctx.Connection.RemoteAddr().String())
	}
	return nil
}
