package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/silentsokolov/go-vimeo/vimeo"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

// TODO: add mutexes.

var userLikesCache map[string][]string

var videoFansCache map[string][]string

func (c Client) getUserLikes(userID string) ([]string, error) {
	likes := c.getLikesFromCache(userID)
	if likes != nil {
		return likes, nil
	}
	returned, _, err := c.vc.Users.ListLikedVideo(userID)
	if err != nil {
		log.Errorf(c.ctx, "Error from the Vimeo API: "+err.Error())
		return nil, err
	}
	out := APIFilterVideos(returned)
	c.setLikesInCache(userID, out)
	return out, nil
}

func (c Client) getLikesFromCache(userID string) []string {
	item, err := memcache.Get(c.ctx, userID)
	if err == nil {
		var likes []string
		err := json.Unmarshal(item.Value, &likes)
		if err != nil {
			log.Errorf(c.ctx, "Error unmarshaling memcache value JSON: "+err.Error())
			return nil
		}
		return likes
	}
	log.Warningf(c.ctx, "Could not get data from memcache: "+err.Error())
	return nil
}

func (c Client) setLikesInCache(userID string, likes []string) {
	// userLikesCache[userID] = likes
	likesJSON, _ := json.Marshal(likes)
	item := &memcache.Item{
		Key:   userID,
		Value: likesJSON,
	}
	if err := memcache.Set(c.ctx, item); err != nil {
		log.Errorf(c.ctx, "Error writing to memcache: ", err.Error())
	}
}

func APIFilterVideos(vids []*vimeo.Video) []string {
	if vids == nil {
		return []string{}
	}

	var out []string
	for _, vid := range vids {
		sliced := strings.Split(vid.URI, "/")
		out = append(out, sliced[len(sliced)-1])
	}
	return out
}

func (c Client) getVideoFans(videoID string) ([]string, error) {
	fans := c.getFansFromCache(videoID)
	if fans != nil {
		return fans, nil
	}
	idAsInt, _ := strconv.Atoi(videoID)
	returned, _, err := c.vc.Videos.LikeList(idAsInt)
	if err != nil {
		log.Errorf(c.ctx, "Error from the Vimeo API: "+err.Error())
		return nil, err
	}
	out := APIFilterUsers(returned)
	c.setFansInCache(videoID, out)
	return out, nil
}

func (c Client) getFansFromCache(videoID string) []string {
	fans, ok := userLikesCache[videoID]
	if ok {
		log.Infof(c.ctx, "Video cache hit.")
		return fans
	}
	log.Infof(c.ctx, "Video cache miss.")
	return nil
}

func (c Client) setFansInCache(videoID string, fans []string) {
	videoFansCache[videoID] = fans
}

func APIFilterUsers(fans []*vimeo.User) []string {
	if fans == nil {
		return []string{}
	}

	var out []string
	for _, fan := range fans {
		sliced := strings.Split(fan.URI, "/")
		out = append(out, sliced[len(sliced)-1])
	}
	return out
}
