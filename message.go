package main

type StateMessage struct {
	messageType     string
	ballPosition    Position
	player1Position Position
	player2Position Position
	playing         bool
}

type InputMessage struct {
	PlayerNumber int
	Direction    string
}
