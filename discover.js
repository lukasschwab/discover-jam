var Vimeo = require('vimeo').Vimeo;
var CUR_USER = '/users/'+'3164416';
var CLIENT_ID = "b72457b81d83ea94cb92fa281ca45026e3bc0785";
var CLIENT_SECRET = "ACYd0wJDHHpFP8UBRutpiG6DILC7K27QdQQn/tvz9o5i5oH9zWbsglVEKkomFcwrt2DIPxpezcfNXoAmzILDYU5d5B8juO3TrXMQrVfpgCK6vk1uUro/DU7qwTgGQ1K1";
var ACCESS_TOKEN = "cb6b2e8f14964cee9fa23f53166e9c3a";
var client = new Vimeo(CLIENT_ID, CLIENT_SECRET, ACCESS_TOKEN);


function getLikes(path, numResults, callback) {
  client.request({
    path: path + '/likes',
    query: {
      page: 1,
      per_page: numResults,
      fields: 'uri'
    }
  }, function (error, body, status_code, headers) {
    if (error) {
      return null;
    } else {
      return callback(body.data);
    }
  });
}

// // https://stackoverflow.com/questions/20696527/get-list-of-items-in-one-list-that-are-not-in-another
// function getDistinctArr(myLikes,otherLikes) {
//   var distinct=otherLikes.filter(function(item){
//     return myLikes.indexOf(item)==-1;
//   });
//   return distinct;
// }

function getDistinct(myLikesObj,otherLikes) {
  var distinct = [];

  otherLikes.forEach(function(video) {
    if (!(video in myLikesObj)) {
      distinct.push(video.uri);
    }
  });
  return distinct;
}

function main() {
  var bigUsrLst = [];
  var newVideoList = {};

  // get 10 videos the current user has liked
  getLikes(CUR_USER,2,function(myLikedVids){ //10
    // create object from myLikedVids
    myLikesObj = {}
    myLikedVids.forEach(function(video) {
      myLikesObj[video.uri] = true;
    });

    if (myLikedVids != null) {
      for (var i = 0; i < myLikedVids.length; i++) {
        // Create a list of 50 users who have liked the same video
        getLikes(myLikedVids[i].uri,2,function(userList){ //50
          if (userList != null) {
            for (var i = 0; i < userList.length; i++) {
              // Get list of 100 liked videos for user
              getLikes(userList[i].uri,2,function(userLikedVids){ //100
                if (userLikedVids != null) {
                  // Create a recommended videos list
                  var distinctVids = getDistinct(myLikesObj,userLikedVids);
                  var numMutualLikes = myLikedVids.length - distinctVids.length;
                  if (!(numMutualLikes in newVideoList)) {
                    newVideoList[numMutualLikes] = distinctVids;
                  } else {
                    newVideoList[numMutualLikes] = newVideoList[numMutualLikes].concat(distinctVids)
                  }
                  console.log(newVideoList);
                  // associate(userList[i],numMutualLikes);
                }
              });
            }
            // newVideoList = sort(newVideoList,mutualLikes);
            // bigUsrLst.append(newVideoList)
          }
        });
      }
    }
    // var recList = bigUsrLst[:10].values();
  });
}

main();
