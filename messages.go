package main

type CommandType string

const (
	CommandQuit CommandType = "QUIT"
	CommandPass CommandType = "PASS"
	CommandNick CommandType = "NICK"
)

type Message struct {
	Prefix  string
	Command CommandType
	Params  []string
}

func (m *Message) String() string {
	result := m.Prefix
	if result != "" {
		result += " "
	}
	result += string(m.Command)
	for _, param := range m.Params {
		result += " " + param
	}
	return result
}

// NewMessage creates a new Message with a valid command
func NewMessage(prefix string, command CommandType, params ...string) *Message {
	return &Message{
		Prefix:  prefix,
		Command: command,
		Params:  params,
	}
}
