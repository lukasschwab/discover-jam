package main

import (
	// "time"

	"github.com/silentsokolov/go-vimeo/vimeo"
	"golang.org/x/oauth2"
)

type Client struct {
	*vimeo.Client
}

func NewClient() Client {
	userLikesCache = make(map[string][]string)
	videoFansCache = make(map[string][]string)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: TOKEN},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return Client{vimeo.NewClient(tc, nil)}
}

func (c Client) RecommendationsFor(userID string) ([]string, error) {
	out, err := c.getUserLikes(userID)
	return out, err
}

var defaultValues = []int{96431363, 270062970, 277328815, 276246978, 276103410}
