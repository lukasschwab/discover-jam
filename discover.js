var Vimeo = require('vimeo').Vimeo;
var CUR_USER = '/users/' + '3164416';
var CLIENT_ID = "b72457b81d83ea94cb92fa281ca45026e3bc0785";
var CLIENT_SECRET = "ACYd0wJDHHpFP8UBRutpiG6DILC7K27QdQQn/tvz9o5i5oH9zWbsglVEKkomFcwrt2DIPxpezcfNXoAmzILDYU5d5B8juO3TrXMQrVfpgCK6vk1uUro/DU7qwTgGQ1K1";
var ACCESS_TOKEN = "cb6b2e8f14964cee9fa23f53166e9c3a";
var client = new Vimeo(CLIENT_ID, CLIENT_SECRET, ACCESS_TOKEN);

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
// // https://stackoverflow.com/questions/20696527/get-list-of-items-in-one-list-that-are-not-in-another
// function getDistinctArr(myLikes,otherLikes) {
//   var distinct=otherLikes.filter(function(item){
//     return myLikes.indexOf(item)==-1;
//   });
//   return distinct;
// }

function getDistinct(myLikesObj, otherLikes) {
  var distinct = [];

  otherLikes.forEach(function(video) {
    // console.log(video.uri);
    // console.log(myLikesObj);
    if (!(video.uri in myLikesObj)) {
      distinct.push(video.uri);
    }
  });
  return distinct;
}

async function main() {
  var videoList = {};

  // get 10 videos the current user has liked
  var myLikedVids = await getLikes(CUR_USER, 2); //10
  // console.log("my liked vids:")
  // console.log(myLikedVids)
  // create object from myLikedVids
  myLikesObj = {}
  myLikedVids.forEach(function(video) {
    myLikesObj[video.uri] = true;
  });

  if (myLikedVids != null) {
    for (var i = 0; i < myLikedVids.length; i++) {
      // Create a list of 50 users who have liked the same video
      var userList = await getLikes(myLikedVids[i].uri, 2); //50
      if (userList != null) {
        // Get list of 100 liked videos for user
        for (var j = 0; j < userList.length; j++) {
          var userLikedVids = await getLikes(userList[j].uri, 2); //100
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
  console.log(myLikedVids)
  console.log(videoList)

  // Sort videoList by its user's mutual likes
  videoList = Object.keys(videoList).sort(function(a, b) {
    return videoList[b] - videoList[a];
  })

}

main();
