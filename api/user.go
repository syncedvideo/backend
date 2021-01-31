package syncedvideo

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	ID      uuid.UUID `db:"id" json:"id"`
	Name    string    `db:"name" json:"username"`
	Color   string    `db:"color" json:"chatColor"`
	IsAdmin bool      `db:"is_admin" json:"isAdmin"`

	conn *websocket.Conn
	// send chan []byte
}

func (u *User) SetConnection(conn *websocket.Conn) {
	u.conn = conn
}

func (u *User) SetUsername(name string) *User {
	u.Name = name
	return u
}

func (u *User) SetChatColor(color string) *User {
	u.Color = color
	return u
}

func (u *User) CanUpdateRoom() bool {
	return true
}
