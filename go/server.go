package main

import (
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/appengine"
)

func recommendationsHandler(w http.ResponseWriter, r *http.Request) {
	cli := NewClient(r)

	u := r.Header.Get("User-Uri")
	w.Header().Set("User-Uri", u)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	log.Print("Handling request for user: ", u)

	out, err := cli.RecommendationsFor(u)
	if err != nil {
		log.Print("User does not exist: ", u)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if len(out) == 0 {
		log.Print("No recommendations for user: ", u)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonList, err := json.Marshal(out)
	if err != nil {
		log.Print("Could not marshal JSON; error in recs?: ", u)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Defaults to status 200.
	w.Write(jsonList)
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
  case "OPTIONS":
    w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "User-Uri")
  case "GET":
		recommendationsHandler(w, r)
  }
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world."))
}

func main() {
	log.Print("Serving recommendations at :8080.")
	log.Print("API docs at https://gist.github.com/lukasschwab/948817751b4bd1ed4909fd31eb7d9fad.")

	http.HandleFunc("/recommendations", corsHandler)
	appengine.Main()
}
