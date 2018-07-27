package main

import (
  "log"
  "strings"

  "github.com/silentsokolov/go-vimeo/vimeo"
)

// TODO: replace with spanner.

// TODO: add mutexes.

var userLikesCache map[string][]string

var videoFansCache map[string][]string

func (c Client) getUserLikes(userID string) ([]string, error) {
  videos, ok := userLikesCache[userID]
  if ok {
    log.Print("Cache hit.")
    return videos, nil
  } else {
    log.Print("Cache miss.")
    returned, _, err := c.Users.ListLikedVideo(userID)
    if err != nil {
      log.Print("Error from Vimeo API: ", err)
      return nil, err
    }
    out := APIFilterVideos(returned)
    userLikesCache[userID] = out
    return out, nil
  }
}

func APIFilterVideos(vids []*vimeo.Video) []string {
  var out []string
  if vids == nil {
    return []string{}
  }

  for _, vid := range vids {
    sliced := strings.Split(vid.URI, "/")
    out = append(out, sliced[len(sliced) - 1])
  }
  return out
}
