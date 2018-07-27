package main

import (
  "strconv"
  "log"
)

// TODO: replace with spanner.

// TODO: add mutexes.

var userLikesCache map[user][]video

var videoFansCache map[video][]user

type user struct {
  id int
}

type video struct {
  id int
}

func (c Client) getUserLikes(userID int) []int {
  u := user{id: userID}
  videos, ok := userLikesCache[u]
  if ok {
    log.Print("Cache hit.")
  } else {
    log.Print("Cache miss.")
    // TODO: get videos from Vimeo.

    returned, _, err := c.Users.ListLikedVideo(strconv.Itoa(userID))
    log.Print(returned)
    log.Print(err) // if it's a 404, report.

    // Process them into []video.
    // userLikesCache[u] = videos
  }
  return toIDs(videos)
}

func toIDs(vids []video) []int {
  var out []int
  if vids == nil {
    return []int{}
  }

  for _, vid := range vids {
    out = append(out, vid.id)
  }
  return out
}
