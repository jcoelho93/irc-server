package commands

import (
	"fmt"
	"log/slog"
)

type NickCommand struct {
	NewNick string
}

func (c NickCommand) Name() string { return "NICK" }
func (c NickCommand) Arguments() []string {
	return []string{c.NewNick}
}

func (c NickCommand) Validate() error {
	if c.Name() == "" || c.Name() != "NICK" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	if len(c.Arguments()) != 1 {
		return fmt.Errorf("NICK command requires exactly one argument")
	}
	return nil
}

func (c NickCommand) Execute(ctx *Ctx) error {
	if ctx.Server.IsUsernameTaken(c.NewNick) {
		return fmt.Errorf("nickname %s is already in use", c.NewNick)
	}
	ctx.Server.SetNick(ctx.Connection, c.NewNick)
	slog.Info("Nickname changed", "new_nick", c.NewNick)

	user, _ := ctx.Server.GetClient(ctx.Connection)

	if user.GetNickname() != "" && user.GetUsername() != "" {
		// Send welcome message
		welcomeMessage := fmt.Sprintf(":%s NOTICE %s :Welcome to the Internet Relay Chat Network %s!%s@%s",
			ctx.Server.GetHostname(), user.GetNickname(),
			user.GetUsername(), user.GetNickname(), ctx.Server.GetHostname())
		_, err := ctx.Connection.Write([]byte(welcomeMessage + "\r\n"))
		if err != nil {
			return fmt.Errorf("failed to send welcome message: %w", err)
		}
	}

	return nil
}
