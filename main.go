package main

import (
	"html/template"
	"net/http"
	"os"
)

//Package level template definition
var index_templ = template.Must(template.ParseFiles("templates/index.html"))

//Writes the index template to the response writer
func indexHandler(w http.ResponseWriter, r *http.Request) {
	index_templ.Execute(w, nil)
}

func main() {
  //Check to see if the PORT env variable is avaialbe and if so set it
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

  //Standard http Multiplexer for routing
	mux := http.NewServeMux()

  //Tell the Multiplexer to route all calls to static files to file_server
  file_server := http.FileServer(http.Dir("static"))
  mux.Handle("/static/", http.StripPrefix("/static/", file_server))

	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
