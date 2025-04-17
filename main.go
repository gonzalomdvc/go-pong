package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const SPEED = 10 * time.Millisecond

var ballCenter = 250 - BALL_DIAMETER/2

func main() {

	hub := &Hub{}
	hub.reset()

	go hub.run()
	go step(hub)

	defer func() {
		for client := range hub.clients {
			hub.unregister <- client
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\nReceived CTRL+C, shutting down gracefully...")
		os.Exit(0) // Exit the program after handling the signal
	}()

	http.HandleFunc("/pvp-game", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Errorf("%v", err)
		}
		var player int
		if len(hub.clients) >= 2 {
			http.Error(w, "Game is full", http.StatusBadRequest)
		} else if len(hub.clients) == 1 {
			var currentPlayer *Client
			for client := range hub.clients {
				currentPlayer = client
				break
			}
			if currentPlayer.player == 1 {
				player = 2
			} else {
				player = 1
			}
		} else if len(hub.clients) == 0 {
			player = 1
		}

		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), player: player}
		hub.register <- client
		client.sendPlayerNumber()
		go client.readPump()
		go listen(conn, hub)

	})

	http.HandleFunc("/pvai-game", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		if err != nil {
			fmt.Errorf("%v", err)
		}

		if len(hub.clients) > 0 {
			http.Error(w, "Game is full", http.StatusBadRequest)
		}

		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), player: 1}
		hub.register <- client
		hub.ai <- 2
		client.sendPlayerNumber()
		go client.readPump()
		go listen(conn, hub)
	})

	http.HandleFunc("/aivai-game", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		if err != nil {
			fmt.Errorf("%v", err)
		}

		if len(hub.clients) > 0 {
			http.Error(w, "Game is full", http.StatusBadRequest)
		}

		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), player: 0}
		hub.register <- client
		hub.ai <- 1
		hub.ai <- 2
		go client.readPump()
		go listen(conn, hub)
		hub.state.playing = true
		hub.state.ballDirection = Direction{up: 1, down: 0, left: 0, right: 1}

	})

	http.HandleFunc("/pvp", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pong.html")
	})
	http.HandleFunc("/pvai", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pong.html")
	})
	http.HandleFunc("/aivai", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pong.html")
	})
	http.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "menu.html")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/menu", http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)

}

func listen(conn *websocket.Conn, hub *Hub) {
	var message InputMessage
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		message, err = parseMessage(msg)
		if err != nil {
			fmt.Printf("Error parsing message: %v\n", err)
			continue
		}
		playerNumber := message.PlayerNumber
		switch message.Direction {
		case "up":
			hub.state.calculateNewPosition("up", playerNumber)
		case "down":
			hub.state.calculateNewPosition("down", playerNumber)
		case "space":
			centerBallOrPlay(hub)
		case "right":
			hub.state.playing = !hub.state.playing
		}

	}

}

func step(hub *Hub) {
	for {

		time.Sleep(SPEED)

		if !hub.state.playing {
			// Restart the game automatically for AI vs AI
			if hub.aiPlayers[1] && hub.aiPlayers[2] {
				time.Sleep(time.Second)
				centerBallOrPlay(hub)
			}
			hub.broadcast <- getStateMessage(hub)
			continue
		}

		hub.state.ballMovement()

		for num, _ := range hub.aiPlayers {
			// Calculate player movement automatically for AI players
			newDirection := hub.state.calculateNewDirectionForPlayer(num)
			hub.state.calculateNewPosition(newDirection, num)
		}

		hub.broadcast <- getStateMessage(hub)

	}
}

func parseMessage(msg []byte) (InputMessage, error) {
	var inputMsg InputMessage
	err := json.Unmarshal(msg, &inputMsg)
	if err != nil {
		return InputMessage{}, fmt.Errorf("Error parsing input message: %v\n", err)
	}
	//fmt.Printf("Parsed message: %v\n", inputMsg)
	return inputMsg, nil
}

func getStateMessage(hub *Hub) []byte {
	encodedState := map[string]interface{}{
		"messageType":     "gameState",
		"ballPosition":    hub.state.ballPosition,
		"player1Position": hub.state.player1Position,
		"player2Position": hub.state.player2Position,
		"playing":         hub.state.playing,
		"score":           hub.state.score,
	}
	stateMessage, err := json.Marshal(encodedState)

	if err != nil {
		fmt.Errorf("Could not parse game state")
		panic("An error occurred")
	}

	return stateMessage
}

func centerBallOrPlay(hub *Hub) {
	if (!hub.state.playing && hub.state.ballPosition == Position{X: ballCenter, Y: ballCenter}) {
		hub.state.playing = true
		hub.state.ballDirection = Direction{up: 0, down: 0, left: 0, right: 1}
	} else {
		if !hub.state.playing {
			hub.state.ballPosition = Position{X: ballCenter, Y: ballCenter}
			hub.state.ballDirection = Direction{up: 0, down: 0, left: 0, right: 0}
		}
	}
}
