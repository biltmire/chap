package main
import (	"log"
          "github.com/gorilla/websocket"
          "net/http")

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan QueueMessage) // broadcast channel

type QueueMessage struct {
    Flag string `json:"flag"`
    Url string `json:"url"`
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Handle connections and the receiving of messagess
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	// Register our new client
	clients[ws] = true
	for {
		var msg QueueMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("handleConnections error: %v from client %s", err,r.Host)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

//Sends messages to clients
func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("handleMessages error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
