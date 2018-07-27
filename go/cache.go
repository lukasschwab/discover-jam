package main

import (
  "log"
  "strings"
  "strconv"

  "github.com/silentsokolov/go-vimeo/vimeo"
)

// TODO: replace with spanner.

// TODO: add mutexes.

var userLikesCache map[string][]string

var videoFansCache map[string][]string

func (c Client) getUserLikes(userID string) ([]string, error) {
  videos, ok := userLikesCache[userID]
  if ok {
    log.Print("User cache hit.")
    return videos, nil
  } else {
    log.Print("User cache miss.")
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
  if vids == nil {
    return []string{}
  }

  var out []string
  for _, vid := range vids {
    sliced := strings.Split(vid.URI, "/")
    out = append(out, sliced[len(sliced) - 1])
  }
  return out
}

func (c Client) getVideoFans(videoID string) ([]string, error) {
  fans, ok := videoFansCache[videoID]
  if ok {
    log.Print("Video cache hit.")
    return fans, nil
  } else {
    log.Print("Video cache miss.")
    idAsInt, _ := strconv.Atoi(videoID)
    returned, _, err := c.Videos.LikeList(idAsInt)
    if err != nil {
      log.Print("Error from Vimeo API: ", err)
      return nil, err
    }
    out := APIFilterUsers(returned)
    videoFansCache[videoID] = out
    return out, nil
  }
}

func APIFilterUsers(fans []*vimeo.User) []string {
  if fans == nil {
    return []string{}
  }

  var out []string
  for _, fan := range fans {
    sliced := strings.Split(fan.URI, "/")
    out = append(out, sliced[len(sliced) - 1])
  }
  return out
}
