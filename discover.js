var Vimeo = require('vimeo').Vimeo;
var keys = require('./keys');
var CUR_USER = '/users/' + '3164416';
var client = new Vimeo(keys.CLIENT_ID, keys.CLIENT_SECRET, keys.ACCESS_TOKEN);
const express = require('express');
const app = express();

app.get('/', async (req, res) => {
  const results = await main();
  console.log(results);
  res.send(results);
});

app.listen(3000, () => {
  console.log('Listening on port 3000');
});


function getLikes(path, numResults) {
  return new Promise(function(resolve, reject) {
    client.request({
      path: path + '/likes',
      query: {
        page: 1,
        per_page: numResults,
        fields: 'uri'
      }
    }, function(error, body, status_code, headers) {
      if (error) {
        reject(error);
      } else {
        resolve(body.data);
      }
    });
  });
}

function getDistinct(myLikesObj, otherLikes) {
  var distinct = [];

  otherLikes.forEach(function(video) {
    if (!(video.uri in myLikesObj)) {
      distinct.push(video.uri);
    }
  });
  return distinct;
}

async function main() {
  var videoList = {};
  // get all videos the current user has liked
  var myLikedVids = await getLikes(CUR_USER, 100);
  // create object from myLikedVids
  myLikesObj = {}
  myLikedVids.forEach(function(video) {
    myLikesObj[video.uri] = true;
  });
  // get 10 videos the current user has liked
  var myTenLikedVids = myLikedVids.slice(0,10);

  if (myTenLikedVids != null) {
    for (var i = 0; i < myTenLikedVids.length; i++) {
      // Create a list of 10 users who have liked the same video
      var userList = await getLikes(myTenLikedVids[i].uri, 10);
      if (userList != null) {
        // Get list of 50 liked videos for user
        for (var j = 0; j < userList.length; j++) {
          var userLikedVids = await getLikes(userList[j].uri, 50);
          if (userLikedVids != null) {
            // Create a recommended videos list
            var distinctVids = await getDistinct(myLikesObj, userLikedVids);
            var priority = userLikedVids.length - distinctVids.length;
            distinctVids.forEach(function(video) {
              if (!(video in videoList)) {
                videoList[video] = priority;
              } else {
                videoList[video] += priority;
              }
            });
          }
        }
      }
    }
  }
  // Sort videoList by its user's mutual likes
  videoList = Object.keys(videoList).sort(function(a, b) {
    return videoList[b] - videoList[a];
  })

  // Get only the ten videos with most mutual likes. If fewer than ten, don't send any data.
  if (videoList.length < 10) {
    return [];
  } else {
    videoList = videoList.slice(0,10);
    for (var i = 0; i < videoList.length; i++) {
      videoList[i] = videoList[i].substring(videoList[i].lastIndexOf('/'));
    }
    return videoList;
  }
}
