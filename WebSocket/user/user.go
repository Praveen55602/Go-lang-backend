package user

import "github.com/gorilla/websocket"

type User struct {
	Name string
	Conn *websocket.Conn
}

//in first bracket we are telling go that this method can be called upon a user object only and will have access to it's properties also
func (user *User) RecieveMessage(message string) error {
	return user.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}
