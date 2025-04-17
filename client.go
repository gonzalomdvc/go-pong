package main

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	player int
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message := <-c.send:
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendPlayerNumber() {
	if c == nil {
		fmt.Printf("Client is null")
	} else {
		c.conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(c.player)))
	}
}
