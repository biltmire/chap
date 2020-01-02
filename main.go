package main

import (
	"html/template"
	"net/http"
	"os"
	"fmt"
	"log"
	"github.com/gorilla/websocket"
)

// Define our message object
type Message struct {
	Message  string `json:"message"`
}

//Package level template definition
var index_templ = template.Must(template.ParseFiles("templates/index.html"))
var queue_templ = template.Must(template.ParseFiles("templates/queue.html"))
var game_templ = template.Must(template.ParseFiles("templates/game.html"))

//List of hosts
var host_queue []string

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message) // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Writes the index template to the response writer
func indexHandler(w http.ResponseWriter, r *http.Request) {
	index_templ.Execute(w, nil)
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
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
		log.Printf(msg.Message)
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
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

//Handler for the queue which adds players to the queue until two eligible
//players are found and then removes them from the queue
func queueHandler(w http.ResponseWriter, r *http.Request) {
	//Establish the websocket with the client
	if len(host_queue) > 0 {
		for host := range host_queue {
			if(host_queue[host] != r.Host) {
				host_queue = host_queue[1:]
				http.Redirect(w,r,"/game", http.StatusFound)
			}
		}
	} else {
		host_queue = append(host_queue,r.Host)
		queue_templ.Execute(w,nil)
	}
	fmt.Println(host_queue)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	game_templ.Execute(w,nil)
}

func main() {
  //Check to see if the PORT env variable is avaialbe and if so set it
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go handleMessages()

  //Standard http Multiplexer for routing
	mux := http.NewServeMux()

  //Tell the Multiplexer to route all calls to static files to file_server
  file_server := http.FileServer(http.Dir("static"))
  mux.Handle("/static/", http.StripPrefix("/static/", file_server))
	//Route the mux from requests to handlers
	mux.HandleFunc("/queue",queueHandler)
	mux.HandleFunc("/game",gameHandler)
	mux.HandleFunc("/ws",handleConnections)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
