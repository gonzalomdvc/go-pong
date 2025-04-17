package main

import "fmt"

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// AI players
	aiPlayers map[int]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Adds an AI to play
	ai chan int

	//Game state
	state *State
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// Add the client to the hub
			h.clients[client] = true
			fmt.Println("Client registered\n")
		case client := <-h.unregister:
			// Remove the client from the hub
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				fmt.Println("Client unregistered\n")
			}
			if len(h.clients) == 0 {
				h.reset()
			}
		case message := <-h.broadcast:
			// Broadcast the message to all clients
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case aiPlayer := <-h.ai:
			h.aiPlayers[aiPlayer] = true
			fmt.Println("AI player registered")
		}
	}
}

func (h *Hub) reset() {
	h.state = &State{
		playing:         false,
		ballPosition:    Position{X: ballCenter, Y: ballCenter},
		ballDirection:   Direction{up: 0, down: 0, left: 0, right: 1},
		player1Position: 200,
		player2Position: 200,
		leftBound:       0,
		rightBound:      500,
		topBound:        0,
		bottomBound:     500,
	}
	h.aiPlayers = make(map[int]bool)
	h.clients = make(map[*Client]bool)
	h.broadcast = make(chan []byte)
	h.register = make(chan *Client)
	h.unregister = make(chan *Client)
	h.ai = make(chan int)
}
