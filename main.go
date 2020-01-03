package main

import (
	"html/template"
	"net/http"
	"os"
)

//Package level template definition
var index_templ = template.Must(template.ParseFiles("templates/index.html"))
var queue_templ = template.Must(template.ParseFiles("templates/queue.html"))
var game_templ = template.Must(template.ParseFiles("templates/game.html"))

//List of hosts
var host_queue []string
var game_list []Game

//Writes the index template to the response writer
func indexHandler(w http.ResponseWriter, r *http.Request) {
	index_templ.Execute(w, nil)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	game_templ.Execute(w,hostLookup(r.Host))
}

//Handler for the queue which adds players to the queue until two eligible
//players are found and then removes them from the queue
func queueHandler(w http.ResponseWriter, r *http.Request) {
	//Establish the websocket with the client
	if len(host_queue) > 0 {
		for host := range host_queue {
			if(host_queue[host] != r.Host) {
				gameManager(host_queue[host],r.Host)
				host_queue = host_queue[1:]
				http.Redirect(w,r,"/game", http.StatusFound)
			}
		}
	} else {
		host_queue = append(host_queue,r.Host)
	}
  queue_templ.Execute(w,nil)
}

func main() {
  //Check to see if the PORT env variable is avaialbe and if so set it
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	//Concurrently handle incoming messages from clients
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
