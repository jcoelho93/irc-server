package types

type User struct {
	Nickname string
	Username string
	Hostname string
	Realname string
	Password string
}

func (u User) GetNickname() string { return u.Nickname }
func (u User) GetUsername() string { return u.Username }
