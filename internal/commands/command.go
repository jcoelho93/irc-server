package commands

import (
	"net"
)

type Command interface {
	Name() string
	Arguments() []string
	Validate() error
	Execute(ctx *Ctx) error
}

type User interface {
	GetNickname() string
	GetUsername() string
}

type Server interface {
	Start() error
	IsUsernameTaken(nick string) bool
	IsConnectionRegistered(conn net.Conn) bool
	SetNick(conn net.Conn, nick string)
	SetUser(conn net.Conn, username, hostname, realname string) error
	SetPassword(conn net.Conn, password string) error
	GetClient(conn net.Conn) User
	GetHostname() string
}

type Ctx struct {
	Server     Server
	Connection net.Conn
}
