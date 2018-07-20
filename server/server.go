package main

import (
	"net/http"
	"time"

	"./secrets"
	"github.com/silentsokolov/go-vimeo/vimeo"
	"golang.org/x/oauth2"
)

// type videoList struct {
// 	data []video `json:data,`
// }

type video struct {
	uri string
}

type cacheEntry struct {
	likes       []video
	lastUpdated time.Time
}

var cache map[string]cacheEntry

func getUserLikes(cli *vimeo.Client, uid string) []*vimeo.Video {
	vids, _, err := cli.Users.ListVideo(uid)
	if err != nil {
		println(err.Error())
	}
	return vids
}

func getVideoLikers(cli *vimeo.Client, uid string) {

}

func getRecommendationsHandler(cli *vimeo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		getUserLikes(cli, "colintoupe")
		// if len(vids) == 0 {
		//   println("oops")
		// }
		// for _, vid := range(vids) {
		//   println(vid.Name)
		// }
		w.Write([]byte("hello world"))
	}
}

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: secrets.TOKEN},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := vimeo.NewClient(tc, nil)

	http.HandleFunc("/", getRecommendationsHandler(client))
	http.ListenAndServe(":8999", nil)
}
