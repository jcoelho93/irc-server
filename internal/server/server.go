package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/jcoelho93/irc/internal/commands"
)

type User struct {
	Nickname string
	Username string
	Hostname string
	Realname string
	Password string
}

func (u User) GetNickname() string { return u.Nickname }
func (u User) GetUsername() string { return u.Username }

type InternetRelayChatServer struct {
	Port    string
	Clients map[net.Conn]User
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
	irc.Clients = make(map[net.Conn]User)
	l, err := net.Listen("tcp4", irc.Port)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return err
		}
		go irc.handleConnection(c)
	}
}

func (irc *InternetRelayChatServer) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	defer c.Close()

	for {
		command, err := irc.readCommand(c)
		if err != nil {
			panic(err)
		}

		err = command.Validate()
		if err != nil {
			fmt.Printf("Invalid command: %s\n", err)
			panic(err)
		}
		fmt.Printf("Received command: %s\n", command.Name())

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
		fmt.Printf("Connection %s is not registered\n", conn.RemoteAddr().String())
		return false
	}
	fmt.Printf("Connection %s is registered\n", conn.RemoteAddr().String())
	return true
}

func (irc *InternetRelayChatServer) SetNick(conn net.Conn, nick string) {
	fmt.Printf("Setting nickname to %s for connection %s\n", nick, conn.RemoteAddr().String())
	if user, exists := irc.Clients[conn]; exists {
		user.Nickname = nick
		irc.Clients[conn] = user
	} else {
		irc.Clients[conn] = User{Nickname: nick}
	}
}

func (irc *InternetRelayChatServer) SetUser(conn net.Conn, username, hostname, realname string) error {
	fmt.Printf("Setting user %s (%s@%s) for connection %s\n", username, hostname, realname, conn.RemoteAddr().String())
	if user, exists := irc.Clients[conn]; exists {
		user.Username = username
		user.Hostname = hostname
		user.Realname = realname
		irc.Clients[conn] = user
	} else {
		irc.Clients[conn] = User{
			Username: username,
			Hostname: hostname,
			Realname: realname,
		}
	}
	return nil
}

func (irc *InternetRelayChatServer) SetPassword(conn net.Conn, password string) error {
	fmt.Printf("Setting password for connection %s\n", conn.RemoteAddr().String())
	if user, exists := irc.Clients[conn]; exists {
		user.Password = password
		irc.Clients[conn] = user
	} else {
		irc.Clients[conn] = User{Password: password}
	}
	return nil
}

func (irc *InternetRelayChatServer) GetClient(connection net.Conn) commands.User {
	return irc.Clients[connection]
}
