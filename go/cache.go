package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/silentsokolov/go-vimeo/vimeo"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func (c Client) getUserLikes(userID string) ([]string, error) {
	likes := c.cacheGet(userID)
	if likes != nil {
		return likes, nil
	}
	returned, _, err := c.vc.Users.ListLikedVideo(userID)
	if err != nil {
		log.Errorf(c.ctx, "Error from the Vimeo API: "+err.Error())
		return nil, err
	}
	out := APIFilterVideos(returned)
	c.cacheSet(userID, out)
	return out, nil
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
	fans := c.cacheGet(videoID)
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
	c.cacheSet(videoID, out)
	return out, nil
}

func (c Client) cacheGet(id string) []string {
	item, err := memcache.Get(c.ctx, id)
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

func (c Client) cacheSet(id string, value []string) {
	valueJSON, _ := json.Marshal(value)
	item := &memcache.Item{
		Key:   id,
		Value: valueJSON,
	}
	if err := memcache.Set(c.ctx, item); err != nil {
		log.Errorf(c.ctx, "Error writing to memcache: ", err.Error())
	}
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
