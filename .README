Based on the Chat example implementation of [gorilla/websocket](https://github.com/gorilla/websocket)
To start a server, go build and go run . from root. 
## Player vs Player
Up to two client connections are handled separately. State is sent to both and each controls a player, sending input messages via socket. 

## Player vs AI
One client is connected and the other is an AI (just moves towards the ball).

## AI vs AI 
Two AIs play each other, but a client is created to visualize the game. 

--- 

Hub coordinates the client connections and the State. When a player joins a match, it's registered in the Hub. 
Endpoints serve HTML files, and socket endpoints (ending in "-game") upgrade HTTP to Socket connections and create clients. 
Hub broadcasts the state to clients, which pong.html receives on its end and draws out the state in each step. 
I couldn't manage to put game constants into one single source, so they're set both in pong.html and state.go. Keep in mind if you make changes. 
AI vs AI games run in a loop, as they get stuck on the same initial positions and make the same decisions always. 
Would be cool to implement an actual AI. 