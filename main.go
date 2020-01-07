package main

import (
	"html/template"
	"net/http"
	"os"
	"fmt"
)

//Package level template definition
var index_templ = template.Must(template.ParseFiles("templates/index.html"))
var queue_templ = template.Must(template.ParseFiles("templates/queue.html"))
var game_templ = template.Must(template.ParseFiles("templates/game.html"))

//List of hosts waiting for game
var host_queue []string
//var game_list []Game
var game_list = make(map[string]*Game)

//Writes the index template to the response writer
func indexHandler(w http.ResponseWriter, r *http.Request) {
	index_templ.Execute(w, nil)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()["id"]
	player := r.URL.Query()["player"]
	game_templ.Execute(w,hostLookup(key[0],player[0]))
}

//Handler for the queue which adds players to the queue until two eligible
//players are found and then removes them from the queue
func queueHandler(w http.ResponseWriter, r *http.Request) {
	//Get name of player
	player_name := getPlayerName(r)
	if len(host_queue) > 0 {
		for host := range host_queue {
			//Found suitable opponent, start new game
			if(host_queue[host] != player_name) {
				id := gameManager(host_queue[host],player_name)
				//Delete the host from the host_queue
				host_queue = append(host_queue[:host],host_queue[host+1:]...)
				http.Redirect(w,r,"/game?id="+id+"&player="+player_name, http.StatusFound)
				fmt.Println("Go routine spining up")
				var p *Game = game_list[id]
				go handleMoves(p)
			}
		}
	} else {
		host_queue = append(host_queue,player_name)
	}
	fmt.Println(r.RemoteAddr)
  queue_templ.Execute(w,nil)
}

//Returns name of player
func getPlayerName(r *http.Request) string {
	q := r.URL.Query()
	playerName := q["player"]
	return playerName[0]
}

func main() {
  //Check to see if the PORT env variable is avaialbe and if so set it
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
	mux.HandleFunc("/gamesock",gameConnections)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
