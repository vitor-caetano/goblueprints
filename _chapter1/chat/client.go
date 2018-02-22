package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// client represents a single chat user
type client struct {
	// socket is the websocket for this client
	socket *websocket.Conn

	// send is a channel which messages are sent
	send chan *message

	// room is the room this client is chatting in
	room *room

	// userData holds information about the user
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarURL, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}
		c.room.foward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
