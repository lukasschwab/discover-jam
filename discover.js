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
      // console.log(body.data)
      return callback(body.data);
    }
  });
}

// function getDistinct(myLikes, otherLikes) {
//   for (var i = 0; i < myLikes.length; i++) {
//     if myLikes[i]
//   }
//   return distinctVids;
// }

function getDistinct(myLikes,otherLikes) {

}

function main() {
  var bigUsrLst = [];
  var newUserList = {};

  // get 10 videos the current user has liked
  getLikes(CUR_USER,2,function(myLikedVids){ //10
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
                  newUserList.uri = Object.values(getDistinct(myLikedVids,userLikedVids));
                  var numMutualLikes = myLikedVids.length - newUserList.uri.length;
                  console.log(newUserList.uri);
                  console.log(numMutualLikes);
                  // associate(userList[i],numMutualLikes);
                }
              });
            }
            // userList = sort(userList,mutualLikes);
            // bigUsrLst.append(userList)
          }
        });
      }
    }
    // var recList = bigUsrLst[:10].values();
  });
}

main();
