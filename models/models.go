package models

import "github.com/gorilla/websocket"

type value int

const (
	VIEWER value = iota
	ADMIN
)

type User struct {
	Conn   *websocket.Conn
	UserId int
	room   *Room
	Role   int
}

type Room struct {
	Users        []*User
	roomPassword string
	defaultRole  int
}

type ClientInput struct {
	UserID  string `json:"userId"`
	RoomID  string `json:"roomId"`
	Command string `json:"command"`
}
