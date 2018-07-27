package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/appengine"
)

func getRecommendationsHandler(cli Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Header.Get("User-Uri")
		w.Header().Set("User-Uri", u)
		w.Header().Set("Content-Type", "application/json")
		log.Print("Handling request for user: ", u)

		user, err := strconv.Atoi(u)
		if err != nil {
			log.Print("Could not convert user UID: ", u)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		out, err := cli.RecommendationsFor(user)
		if err != nil {
			log.Print("User does not exist: ", user)
			w.WriteHeader(http.StatusBadRequest)
		}
		if len(out) == 0 {
			log.Print("No recommendations for user: ", user)
			w.WriteHeader(http.StatusNotFound)
		}

		jsonList, err := json.Marshal(out)
		if err != nil {
			log.Print("Could not marshal JSON; error in recs?: ", user)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Defaults to status 200.
		w.Write(jsonList)
	}
}

func main() {
	log.Print("Serving recommendations at :8080.")
	log.Print("API docs at https://gist.github.com/lukasschwab/948817751b4bd1ed4909fd31eb7d9fad.")

	http.HandleFunc("/recommendations", getRecommendationsHandler(NewClient()))
	// log.Fatal(http.ListenAndServe(":8080", nil))
	appengine.Main()
}
