package main

import (
	"log"
	"sync"
	"sort"

	"github.com/silentsokolov/go-vimeo/vimeo"
	"golang.org/x/oauth2"
	"github.com/ekzhu/counter"
	// "github.com/gomodule/redigo/redis"
)

type Client struct {
	vc *vimeo.Client
	// redisPool *redis.Pool
}

func NewClient() Client {
	userLikesCache = make(map[string][]string)
	videoFansCache = make(map[string][]string)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: TOKEN},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return Client{
		vc: vimeo.NewClient(tc, nil),
	// 	redisPool: &redis.Pool{
  //     Dial: func() (redis.Conn, error) {
  //       conn, err := redis.Dial("tcp", REDIS_ADDR)
  //       if REDIS_PASS == "" {
  //         return conn, err
  //       }
  //       if err != nil {
  //         return nil, err
  //       }
  //       if _, err := conn.Do("AUTH", REDIS_PASS); err != nil {
  //         conn.Close()
  //         return nil, err
  //       }
  //       return conn, nil
  //     },
  //   },
	}
}

func (c Client) RecommendationsFor(userID string) ([]string, error) {
	vids, err := c.getUserLikes(userID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	out := make(chan string, 1000000)
	for x, vid := range vids {
		wg.Add(1)
		go c.compileRecs(vid, out, &wg)
		if x > 5 { // FIXME
			break
		}
	}
	wg.Wait()
	close(out)

	ctr := counter.NewCounter()
	for candidate := range out {
		ctr.Update(candidate)
	}
	uqs, freqs := ctr.Freqs()
	sort.Slice(uqs, func (i, j int) bool {
		return freqs[i] > freqs[j]
	})

	var recs []string
	i := 0
	for len(recs) < 10 {
		s, _ := uqs[i].(string)
		// Don't include videos that have already been liked.
		if !contains(vids, s) {
			recs = append(recs, s)
		}
		i++
	}

	return recs, err
}

func (c Client) compileRecs(vid string, out chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Print("Goroutine started.")
	fans, err := c.getVideoFans(vid)
	if err != nil {
		return
	}
	for x, fan := range(fans) {
		log.Print("Iterating over a video's fans.")
		recCandidates, err := c.getUserLikes(fan)
		if err != nil {
			return
		}
		for _, r := range recCandidates {
			out <- r
		}
		if x > 5 { // FIXME
			break
		}
	}
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
