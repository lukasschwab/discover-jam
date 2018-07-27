package main

import (
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/appengine"
)

func getRecommendationsHandler(cli Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Header.Get("User-Uri")
		w.Header().Set("User-Uri", u)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world."))
}

func main() {
	log.Print("Serving recommendations at :8080.")
	log.Print("API docs at https://gist.github.com/lukasschwab/948817751b4bd1ed4909fd31eb7d9fad.")

	http.HandleFunc("/recommendations", getRecommendationsHandler(NewClient()))
	// http.HandleFunc("/", helloWorldHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))
	appengine.Main()
}
