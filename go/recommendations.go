package main

import (
	"context"
	"net/http"
	"sort"
	"sync"

	"github.com/ekzhu/counter"
	"github.com/silentsokolov/go-vimeo/vimeo"
	"golang.org/x/oauth2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const iterLimit = 5 // FIXME

type Client struct {
	vc  *vimeo.Client
	ctx context.Context
}

func NewClient(r *http.Request) Client {
	ctx := appengine.NewContext(r)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)

	return Client{
		vc:  vimeo.NewClient(tc, nil),
		ctx: ctx,
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
		if x > iterLimit {
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
	sort.Slice(uqs, func(i, j int) bool {
		return freqs[i] > freqs[j]
	})

	var recs []string
	i := 0
	for len(recs) < 10 && i < len(uqs) {
		s, _ := uqs[i].(string)
		// Don't include videos that have already been liked.
		if !contains(vids, s) {
			recs = append(recs, s)
		}
		i++
	}

	log.Infof(c.ctx, "Done computing recommendations.")
	return recs, err
}

func (c Client) compileRecs(vid string, out chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	fans, err := c.getVideoFans(vid)
	if err != nil {
		return
	}
	for x, fan := range fans {
		recCandidates, err := c.getUserLikes(fan)
		if err != nil {
			return
		}
		for _, r := range recCandidates {
			out <- r
		}
		if x > iterLimit {
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
