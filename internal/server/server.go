package server

import (
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/jcoelho93/irc/internal/commands"
	"github.com/jcoelho93/irc/internal/types"
)

type InternetRelayChatServer struct {
	Port    string
	Clients map[net.Conn]types.User
}

func (irc *InternetRelayChatServer) GetHostname() string {
	return "irc.example.com"
}

func NewInternetRelayChatServer(port string) *InternetRelayChatServer {
	return &InternetRelayChatServer{
		Port: port,
	}
}

func (irc *InternetRelayChatServer) Start() error {
	slog.Info("Starting IRC server", "port", irc.Port)
	irc.Clients = make(map[net.Conn]types.User)
	l, err := net.Listen("tcp4", irc.Port)
	if err != nil {
		slog.Error("Failed to start IRC server", "error", err)
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			return err
		}
		go irc.handleConnection(c)
	}
}

func (irc *InternetRelayChatServer) handleConnection(c net.Conn) {
	slog.Info("New connection established", "remote_addr", c.RemoteAddr().String())
	defer c.Close()

	for {
		command, err := irc.readCommand(c)
		if err != nil {
			panic(err)
		}

		err = command.Validate()
		if err != nil {
			slog.Error("Invalid command", "error", err)
			panic(err)
		}
		slog.Info("command received", "command", command.Name(), "remote_addr", c.RemoteAddr().String())

		err = command.Execute(&commands.Ctx{Server: irc, Connection: c})
		if err != nil {
			panic(err)
		}
	}
}

func (irc *InternetRelayChatServer) readCommand(c net.Conn) (commands.Command, error) {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}

	rawCommand := string(buf[:n])
	rawCommand = strings.TrimSpace(rawCommand)
	arguments := strings.Split(rawCommand, " ")
	command := arguments[0]
	slog.Info("Command received: ", "remote_addr", c.RemoteAddr().String(), "command", command)
	slog.Info("", "command", command, "arguments", strings.Join(arguments, " "))
	switch command {
	case "PASS":
		return commands.PassCommand{
			Password: arguments[1],
		}, nil
	case "NICK":
		return commands.NickCommand{
			NewNick: arguments[1],
		}, nil
	case "ECHO":
		return commands.EchoCommand{
			Message: strings.Join(arguments[1:], " "),
		}, nil
	case "USER":
		return commands.UserCommand{
			Username:   arguments[1],
			Hostname:   arguments[2],
			Servername: arguments[3],
			Realname:   strings.Join(arguments[4:], " "),
		}, nil
	case "PING":
		return commands.PingCommand{}, nil
	case "CAP":
		return commands.CapCommand{}, nil
	case "JOIN":
		return commands.JoinCommand{
			Channels: arguments[1:],
		}, nil
	case "QUIT":
		return commands.QuitCommand{}, nil
	case "PRIVMSG":
		return commands.PrivMsgCommand{
			Target:  arguments[1],
			Message: strings.Join(arguments[2:], " "),
		}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}

func (irc *InternetRelayChatServer) IsUsernameTaken(username string) bool {
	for _, existingUsername := range irc.Clients {
		if existingUsername.Username == username {
			return true
		}
	}
	return false
}

func (irc *InternetRelayChatServer) IsConnectionRegistered(conn net.Conn) bool {
	_, exists := irc.Clients[conn]
	if !exists {
		slog.Warn("Connection is not registered", "remote_addr", conn.RemoteAddr().String())
		return false
	}
	slog.Info("Connection is registered", "remote_addr", conn.RemoteAddr().String())
	return true
}

func (irc *InternetRelayChatServer) SetNick(conn net.Conn, nick string) {
	slog.Info("Setting nickname", "nickname", nick, "remote_addr", conn.RemoteAddr().String())
	if user, exists := irc.Clients[conn]; exists {
		user.Nickname = nick
		irc.Clients[conn] = user
	} else {
		irc.Clients[conn] = types.User{Nickname: nick}
	}
}

func (irc *InternetRelayChatServer) SetUser(conn net.Conn, username, hostname, realname string) error {
	slog.Info("Setting user details", "username", username, "hostname", hostname, "realname", realname, "remote_addr", conn.RemoteAddr().String())
	if user, exists := irc.Clients[conn]; exists {
		user.Username = username
		user.Hostname = hostname
		user.Realname = realname
		irc.Clients[conn] = user
	} else {
		irc.Clients[conn] = types.User{
			Username: username,
			Hostname: hostname,
			Realname: realname,
		}
	}
	return nil
}

func (irc *InternetRelayChatServer) SetPassword(conn net.Conn, password string) error {
	slog.Info("Setting password for connection", "remote_addr", conn.RemoteAddr().String())
	if user, exists := irc.Clients[conn]; exists {
		user.Password = password
		irc.Clients[conn] = user
	} else {
		irc.Clients[conn] = types.User{Password: password}
	}
	return nil
}

func (irc *InternetRelayChatServer) GetClient(connection net.Conn) (types.User, bool) {
	user := irc.Clients[connection]
	if connection == nil {
		slog.Warn("Client connection not found", "remote_addr", connection.RemoteAddr().String())
		return types.User{}, false
	}
	slog.Info("Client connection found", "remote_addr", connection.RemoteAddr().String())
	return user, true
}

func (irc *InternetRelayChatServer) GetClients() map[net.Conn]types.User {
	slog.Info("Retrieving all clients")
	if irc.Clients == nil {
		slog.Warn("No clients found")
		return make(map[net.Conn]types.User)
	}
	slog.Info("Clients retrieved", "count", len(irc.Clients))
	return irc.Clients
}
