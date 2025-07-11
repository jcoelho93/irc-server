package server

import (
	"errors"
	"fmt"
	"io"
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
	defer func() {
		defer c.Close()
		slog.Info("Connection closed", "remote_addr", c.RemoteAddr().String())
	}()

	for {
		commandList, err := irc.readCommand(c)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Connection closed by client", "remote_addr", c.RemoteAddr().String())
				return
			}
			slog.Error("Failed to read command", "error", err, "remote_addr", c.RemoteAddr().String())
			return
		}
		for _, cmd := range commandList {
			if err := cmd.Validate(); err != nil {
				slog.Error("Command validation failed", "error", err, "remote_addr", c.RemoteAddr().String())
				continue
			}

			slog.Info("Executing command", "command", cmd.Name(), "remote_addr", c.RemoteAddr().String())

			if err := cmd.Execute(&commands.Ctx{
				Server:     irc,
				Connection: c,
			}); err != nil {
				slog.Error("Failed to execute command", "error", err, "command", cmd.Name(), "remote_addr", c.RemoteAddr().String())
				continue
			}
		}
	}

}

func (irc *InternetRelayChatServer) readCommand(c net.Conn) ([]commands.Command, error) {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}

	rawInput := string(buf[:n])
	lines := strings.Split(rawInput, "\r\n")

	var parsedCommands []commands.Command

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		arguments := strings.Split(line, " ")
		command := arguments[0]

		slog.Info("Command received", "remote_addr", c.RemoteAddr().String(), "command", command)
		slog.Info("", "command", command, "arguments", strings.Join(arguments, " "))

		cmd, err := irc.parseCommand(command, arguments)
		if err != nil {
			slog.Warn("Failed to parse command", "error", err, "remote_addr", c.RemoteAddr().String())
			continue
		}
		parsedCommands = append(parsedCommands, cmd)
	}

	return parsedCommands, nil
}

func (irc *InternetRelayChatServer) parseCommand(command string, arguments []string) (commands.Command, error) {
	command = strings.ToUpper(command)

	switch command {
	case "PASS":
		if len(arguments) < 2 {
			return nil, fmt.Errorf("PASS requires a password")
		}
		return commands.PassCommand{Password: arguments[1]}, nil
	case "NICK":
		if len(arguments) < 2 {
			return nil, fmt.Errorf("NICK requires a nickname")
		}
		return commands.NickCommand{NewNick: arguments[1]}, nil
	case "ECHO":
		return commands.EchoCommand{Message: strings.Join(arguments[1:], " ")}, nil
	case "USER":
		if len(arguments) < 5 {
			return nil, fmt.Errorf("USER requires at least 5 arguments")
		}
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
		return commands.JoinCommand{Channels: arguments[1:]}, nil
	case "QUIT":
		return commands.QuitCommand{}, nil
	case "PRIVMSG":
		if len(arguments) < 3 {
			return nil, fmt.Errorf("PRIVMSG requires a target and a message")
		}
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
