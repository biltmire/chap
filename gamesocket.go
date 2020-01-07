package main
import (	"log"
          "github.com/gorilla/websocket"
          "net/http"
          "fmt")


//Struct that contains information on the game
type Game struct {
	PlayerMap map[string]string
	ConnectionList [2]*websocket.Conn
  MovesChannel chan Message
}

// Define our message object
type Message struct {
	Color string `json:"color"`
  From string `json:"from"`
  To string `json:"to"`
  Flags string `json:"flags"`
  Piece string `json:"piece"`
  San string `json:"san"`
  Fen string `json:"fen"`
}

//Handle connections and the receiving of messagess
func gameConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
  q := r.URL.Query()
  color := q["color"]
  fmt.Println(color)
  key := q["id"]
  p := game_list[key[0]]
  if(color[0] == "white") {
    (*p).ConnectionList[0] = ws
  } else {
    (*p).ConnectionList[1] = ws
  }

	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
      fmt.Println(len(game_list))
      delete(game_list,key[0])
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		(*p).MovesChannel <- msg
	}
}

//Sends messages to clients
func handleMoves(game_obj *Game) {
  var client *websocket.Conn
	for {
		// Grab the next message from the broadcast channel
		msg := <-game_obj.MovesChannel
		// Send it out to the opponent of the received message sender
	  if(msg.Color == "w"){
      client = game_obj.ConnectionList[1]
		} else {
      client = game_obj.ConnectionList[0]
    }
    err := client.WriteJSON(msg)
    if err != nil {
      log.Println("error: %v", err)
      client.Close()
      delete(clients, client)
    }
	}
}
