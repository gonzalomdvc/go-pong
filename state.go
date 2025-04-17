package main

type Position struct {
	X int
	Y int
}

type Direction struct {
	up    float32
	down  float32
	left  float32
	right float32
}

type State struct {
	playing         bool
	ballPosition    Position
	ballDirection   Direction
	player1Position int
	player2Position int
	leftBound       int
	rightBound      int
	topBound        int
	bottomBound     int
	score           [2]int
}

const (
	PLAYER_BAR_WIDTH       = 10
	PLAYER_BAR_HEIGHT      = 100
	BALL_DIAMETER          = 10
	PLAYER_MOVEMENT_FACTOR = 40
	BALL_MOVEMENT_FACTOR   = 4
	STEP_TIME              = 20
	GAME_MODE              = "live"
)

func (s *State) updateBallPosition(newPosition Position) {
	s.ballPosition = newPosition
}

func (s *State) updatePlayerPosition(newPosition int, player int) {
	if newPosition < s.topBound {
		newPosition = s.topBound
	} else if newPosition > s.bottomBound {
		newPosition = s.bottomBound
	}

	if player == 1 {
		s.player1Position = newPosition
	} else {
		s.player2Position = newPosition
	}

}

func (s *State) calculateNewPosition(direction string, player int) {
	var newPosition int
	if player == 1 {
		newPosition = s.player1Position
	} else if player == 2 {
		newPosition = s.player2Position
	}
	//fmt.Printf("Previous position: %d\n", newPosition)
	switch direction {
	case "up":
		newPosition = newPosition - PLAYER_MOVEMENT_FACTOR
		if newPosition < s.topBound {
			newPosition = s.topBound
		}
	case "down":
		newPosition = newPosition + PLAYER_MOVEMENT_FACTOR
		if newPosition > s.bottomBound-PLAYER_BAR_HEIGHT {
			newPosition = s.bottomBound - PLAYER_BAR_HEIGHT
		}
	case "none":
		// Do nothing
	}
	//fmt.Printf("New position: %d\n", newPosition)
	s.updatePlayerPosition(newPosition, player)
	//fmt.Printf("Player %d position: %d\n", player, newPosition)
}

func (s *State) calculateBallPosition() {
	newPosition := s.ballPosition
	newPosition.X = newPosition.X + int(float32(BALL_MOVEMENT_FACTOR)*s.ballDirection.right) - int(float32(BALL_MOVEMENT_FACTOR)*s.ballDirection.left)
	newPosition.Y = newPosition.Y + int(float32(BALL_MOVEMENT_FACTOR)*s.ballDirection.down) - int(float32(BALL_MOVEMENT_FACTOR)*s.ballDirection.up)
	s.updateBallPosition(newPosition)
}

func (s *State) ballMovement() {
	if s.ballPosition.X < s.leftBound+BALL_DIAMETER/2 {
		s.ballDirection.right = 1
		s.ballDirection.left = 0
	} else if s.ballPosition.X > s.rightBound-BALL_DIAMETER {
		s.ballDirection.right = 0
		s.ballDirection.left = 1
	}

	if s.ballPosition.Y < s.topBound+BALL_DIAMETER/2 {
		s.ballDirection.down = 1
		s.ballDirection.up = 0
	} else if s.ballPosition.Y > s.bottomBound-BALL_DIAMETER {
		s.ballDirection.down = 0
		s.ballDirection.up = 1
	}

	s.calculateBallPosition()
	s.calculateCollision(1)
	s.calculateCollision(2)

	if GAME_MODE == "live" {
		// Collide with bounds and update score
		if s.ballPosition.X <= s.leftBound+BALL_DIAMETER/2 {
			s.playing = false
			s.score[1]++
		} else if s.ballPosition.X >= s.rightBound-BALL_DIAMETER {
			s.playing = false
			s.score[0]++
		}
	}

}

func (s *State) calculateCollision(player int) {
	var position int
	var direction float32
	var collision bool
	if player == 1 {
		position = s.player1Position
		direction = s.ballDirection.left
		collision = s.ballPosition.X < PLAYER_BAR_WIDTH
	} else if player == 2 {
		position = s.player2Position
		direction = s.ballDirection.right
		collision = s.ballPosition.X > s.rightBound-PLAYER_BAR_WIDTH-BALL_DIAMETER
	}
	//fmt.Printf("Player %d position: %d\n", player, position)
	//fmt.Printf("Ball position: X: %d Y: %d\n", s.ballPosition.X, s.ballPosition.Y)
	//fmt.Printf("Ball direction: up: %f down: %f left: %f right: %f\n", s.ballDirection.up, s.ballDirection.down, s.ballDirection.left, s.ballDirection.right)
	barMiddle := position + PLAYER_BAR_HEIGHT/2

	// Check for collision with player
	if collision && s.ballPosition.Y > position && s.ballPosition.Y < position+PLAYER_BAR_HEIGHT && direction > 0 {
		if s.ballPosition.Y < barMiddle {
			s.ballDirection.up = normalizeToRange(barMiddle - s.ballPosition.Y)
		} else {
			s.ballDirection.up = 0
		}
		if s.ballPosition.Y > barMiddle {
			s.ballDirection.down = normalizeToRange(s.ballPosition.Y - barMiddle)
		} else {
			s.ballDirection.down = 0
		}

		if player == 1 {
			s.ballDirection.left = 0
			s.ballDirection.right = 1
		} else {
			s.ballDirection.right = 0
			s.ballDirection.left = 1
		}
	}
}

func (s *State) calculateNewDirectionForPlayer(player int) (direction string) {
	var position int
	var currentDirection float32
	if player == 1 {
		position = s.player1Position
		currentDirection = s.ballDirection.left
	} else if player == 2 {
		position = s.player2Position
		currentDirection = s.ballDirection.right
	}
	if (s.ballPosition.Y) > position+PLAYER_BAR_HEIGHT && currentDirection > 0 {
		return "down"
	} else if (s.ballPosition.Y) < position && currentDirection > 0 {
		return "up"
	} else {
		return "none"
	}
}

func normalizeToRange(value int) float32 {
	oldMin := float32(0)
	oldMax := float32(50)

	// New range is 0 to 2
	newMin := float32(0)
	newMax := float32(2)

	return ((float32(value)-oldMin)/(oldMax-oldMin))*(newMax-newMin) + newMin

}
