package commands

import (
	"net"

	"github.com/jcoelho93/irc/internal/types"
)

type Command interface {
	Name() string
	Arguments() []string
	Validate() error
	Execute(ctx *Ctx) error
}

type Server interface {
	Start() error
	IsUsernameTaken(nick string) bool
	IsConnectionRegistered(conn net.Conn) bool
	SetNick(conn net.Conn, nick string)
	SetUser(conn net.Conn, username, hostname, realname string) error
	SetPassword(conn net.Conn, password string) error
	GetClient(conn net.Conn) (types.User, bool)
	GetClients() map[net.Conn]types.User
	GetHostname() string
}

type Ctx struct {
	Server     Server
	Connection net.Conn
}
