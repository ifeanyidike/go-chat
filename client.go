package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user
type client struct {
	// socket is the websocket for this client.
	socket *websocket.Conn
	// send is a channel on which messages are sent.
	send chan *message
	// room is the room this client is chatting in.
	room *room
	// userData holds information about the user
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err != nil {
			break
		} else {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			// if avatarUrl, ok := c.userData["avatar_url"]; ok {
			// 	msg.AvatarURL = avatarUrl.(string)
			// }
			// use Gravatar instead

			// msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(user)

			if avatarUrl, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarUrl.(string)
			}

			c.room.forward <- msg
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(&msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
