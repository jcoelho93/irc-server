package commands

import (
	"fmt"
	"log/slog"
)

type UserCommand struct {
	Username   string
	Hostname   string
	Servername string
	Realname   string
}

func (c UserCommand) Name() string { return "USER" }

func (c UserCommand) Arguments() []string {
	return []string{c.Username, c.Hostname, c.Realname}
}

func (c UserCommand) Validate() error {
	if c.Name() == "" || c.Name() != "USER" {
		return fmt.Errorf("invalid command name: %s", c.Name())
	}
	if len(c.Arguments()) < 3 {
		return fmt.Errorf("USER command requires at least three arguments")
	}
	return nil
}

func (c UserCommand) Execute(ctx *Ctx) error {
	ctx.Server.SetUser(ctx.Connection, c.Username, c.Hostname, c.Realname)

	user, _ := ctx.Server.GetClient(ctx.Connection)

	if user.GetNickname() != "" && user.GetUsername() != "" {
		// Send welcome message
		welcomeMessage := fmt.Sprintf(":%s 001 %s :Welcome to the Internet Relay Chat Network %s!%s@%s",
			ctx.Server.GetHostname(), user.GetNickname(),
			user.GetUsername(), user.GetNickname(), ctx.Server.GetHostname())
		_, err := ctx.Connection.Write([]byte(welcomeMessage + "\r\n"))
		if err != nil {
			return fmt.Errorf("failed to send welcome message: %w", err)
		}
	}

	slog.Info("User set", "username", c.Username, "hostname", c.Hostname, "realname", c.Realname)
	return nil
}
