package main

import (
  "log"
  "strings"
  "strconv"

  "github.com/silentsokolov/go-vimeo/vimeo"
  "google.golang.org/appengine/memcache"
)

// TODO: add mutexes.

var userLikesCache map[string][]string

var videoFansCache map[string][]string

func (c Client) getUserLikes(userID string) ([]string, error) {
  likes := getLikesFromCache(userID)
  if likes != nil {
    return likes, nil
  }
  returned, _, err := c.vc.Users.ListLikedVideo(userID)
  if err != nil {
    log.Print("Error from Vimeo API: ", err)
    return nil, err
  }
  out := APIFilterVideos(returned)
  setLikesInCache(userID, out)
  return out, nil
}

func getLikesFromCache(userID string) []string {
  videos, ok := userLikesCache[userID]
  if ok {
    log.Print("User cache hit.")
    return videos
  }
  log.Print("User cache miss.")
  return nil
}

func setLikesInCache(userID string, likes []string) {
  userLikesCache[userID] = likes
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
  fans := getFansFromCache(videoID)
  if fans != nil {
    return fans, nil
  }
  idAsInt, _ := strconv.Atoi(videoID)
  returned, _, err := c.vc.Videos.LikeList(idAsInt)
  if err != nil {
    log.Print("Error from Vimeo API: ", err)
    return nil, err
  }
  out := APIFilterUsers(returned)
  setFansInCache(videoID, out)
  return out, nil
}

func getFansFromCache(videoID string) []string {
  fans, ok := userLikesCache[videoID]
  if ok {
    log.Print("Video cache hit.")
    return fans
  }
  log.Print("Video cache miss.")
  return nil
}

func setFansInCache(videoID string, fans []string) {
  videoFansCache[videoID] = fans
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
