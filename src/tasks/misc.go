package tasks

import (
  "fmt"
  "time"
  "log"

  "go.mongodb.org/mongo-driver/mongo" // MongoDB
  "go.mongodb.org/mongo-driver/bson/primitive" // MongoDB primitive

  "github.com/pusher/pusher-http-go" // Pusher

  . "github.com/logrusorgru/aurora" // colors

  // ## local shits ##
	"../database"
  "../structs"
)

// ########################################### VARIABLES
var successful_connections = 0
var failed_connections = 0
var initialCheckedURLs = []string{}

// ########################################### UTILITY FUNCTIONS

func SetupSocialPlusActiveHandles(mongoClient *mongo.Client, pusherClient *pusher.Client) {
  results := database.GatherSocialPlusActiveHandles(mongoClient)

  var foundTwitterHandles = []structs.SocialMediaHandle{}
  var foundInstagramHandles = []structs.SocialMediaHandle{}

  // ########################################### GATHER ACTIVE TWITTER AND INSTAGRAM HANDLES TO BE MONITORED/ALREADY MONITORED
  for i := 0;  i < len(results); i++ {
    var socialMediaHandle structs.SocialMediaHandle
    socialMediaHandle.Handle = fmt.Sprintf("%v", results[i]["handle"])
    socialMediaHandle.Platform = fmt.Sprintf("%v", results[i]["platform"])
    socialMediaHandle.IsFollowingHandle = results[i]["following_handle"].(bool)
    var twitterApp structs.TwitterApp
    twitterApp.ConsumerKey = fmt.Sprintf("%v", results[i]["twitter_app"].(primitive.M)["consumer_key"])
  	twitterApp.ConsumerSecret = fmt.Sprintf("%v", results[i]["twitter_app"].(primitive.M)["consumer_secret"])
    twitterApp.AccessTokenKey = fmt.Sprintf("%v", results[i]["twitter_app"].(primitive.M)["access_token"])
    twitterApp.AccessTokenSecret = fmt.Sprintf("%v", results[i]["twitter_app"].(primitive.M)["access_token_secret"])
    socialMediaHandle.TwitterApp = twitterApp
    var reserveTwitterApp structs.TwitterApp
    reserveTwitterApp.ConsumerKey = fmt.Sprintf("%v", results[i]["reserve_twitter_app"].(primitive.M)["consumer_key"])
  	reserveTwitterApp.ConsumerSecret = fmt.Sprintf("%v", results[i]["reserve_twitter_app"].(primitive.M)["consumer_secret"])
    reserveTwitterApp.AccessTokenKey = fmt.Sprintf("%v", results[i]["reserve_twitter_app"].(primitive.M)["access_token"])
    reserveTwitterApp.AccessTokenSecret = fmt.Sprintf("%v", results[i]["reserve_twitter_app"].(primitive.M)["access_token_secret"])
    socialMediaHandle.ReserveTwitterApp = reserveTwitterApp
    if results[i]["platform"] == "twitter" {
      foundTwitterHandles = append(foundTwitterHandles, socialMediaHandle)
    } else if results[i]["platform"] == "instagram" {
      foundInstagramHandles = append(foundInstagramHandles, socialMediaHandle)
    }
  }

  // ########################################### LOOK FOR TWITTER HANDLES TO STOP MONITORING
  for i := 0;  i < len(database.DatabaseTwitterHandles); i++ {
    var foundHandle = false
    for j := 0;  j < len(foundTwitterHandles); j++ {
      if foundTwitterHandles[j].Handle == database.DatabaseTwitterHandles[i].Handle {
        foundHandle = true
        break
      }
    }
    if !foundHandle {
      log.Println(Inverse("[" + "Twitter (Social+)" + "] " + "No-longer Monitoring Handle: " + database.DatabaseTwitterHandles[i].Handle))
      // stop monitoring this Twitter account
      go pusherClient.Trigger("socialPlusChannel", "App\\Events\\handleInactive", "handle=" + database.DatabaseTwitterHandles[i].Handle + ";platform=twitter")
      // remove from initialCheckedURLs
      for j := 0;  j < len(initialCheckedURLs); j++ {
        if initialCheckedURLs[j] == "twitter_" + database.DatabaseTwitterHandles[i].Handle {
          // Remove the element at index i from slice
          initialCheckedURLs[j] = initialCheckedURLs[len(initialCheckedURLs)-1] // Copy last element to index i.
          initialCheckedURLs = initialCheckedURLs[:len(initialCheckedURLs)-1]   // Truncate slice.
          break
        }
      }

      // Remove the element at index i from slice
      database.DatabaseTwitterHandles[i] = database.DatabaseTwitterHandles[len(database.DatabaseTwitterHandles)-1] // Copy last element to index i.
      database.DatabaseTwitterHandles = database.DatabaseTwitterHandles[:len(database.DatabaseTwitterHandles)-1]   // Truncate slice.
      break
    }
  }

  // ########################################### LOOK FOR INSTAGRAM HANDLES TO STOP MONITORING
  for i := 0;  i < len(database.DatabaseInstagramHandles); i++ {
    var foundHandle = false
    for j := 0;  j < len(foundInstagramHandles); j++ {
      if foundInstagramHandles[j].Handle == database.DatabaseInstagramHandles[i].Handle {
        foundHandle = true
        break
      }
    }
    if !foundHandle {
      log.Println(Inverse("[" + "Instagram (Social+)" + "] " + "No-longer Monitoring Handle: " + database.DatabaseInstagramHandles[i].Handle))
      // stop monitoring this Instagram account
      go pusherClient.Trigger("socialPlusChannel", "App\\Events\\handleInactive", "handle=" + database.DatabaseInstagramHandles[i].Handle + ";platform=instagram")
      // remove from initialCheckedURLs
      var removeCount = 0
      for j := 0;  j < len(initialCheckedURLs); j++ {
        if removeCount == 2 {
          break
        }
        if initialCheckedURLs[j] == "instagram_" + database.DatabaseInstagramHandles[i].Handle {
          // Remove the element at index i from slice
          initialCheckedURLs[j] = initialCheckedURLs[len(initialCheckedURLs)-1] // Copy last element to index i.
          initialCheckedURLs = initialCheckedURLs[:len(initialCheckedURLs)-1]   // Truncate slice.
          removeCount++
        } else if initialCheckedURLs[j] == "instagram_" + database.DatabaseInstagramHandles[i].Handle + "_stories" {
          // Remove the element at index i from slice
          initialCheckedURLs[j] = initialCheckedURLs[len(initialCheckedURLs)-1] // Copy last element to index i.
          initialCheckedURLs = initialCheckedURLs[:len(initialCheckedURLs)-1]   // Truncate slice.
          removeCount++
        }
      }

      // Remove the element at index i from slice
      database.DatabaseInstagramHandles[i] = database.DatabaseInstagramHandles[len(database.DatabaseInstagramHandles)-1] // Copy last element to index i.
      database.DatabaseInstagramHandles = database.DatabaseInstagramHandles[:len(database.DatabaseInstagramHandles)-1]   // Truncate slice.
      break
    }
  }

  // ########################################### LOOK FOR TWITTER HANDLES TO BEGIN MONITORING
  for i := 0;  i < len(foundTwitterHandles); i++ {
    var foundHandle = false
    for j := 0;  j < len(database.DatabaseTwitterHandles); j++ {
      if database.DatabaseTwitterHandles[j].Handle == foundTwitterHandles[i].Handle {
        foundHandle = true
        break
      }
    }
    if !foundHandle {
      log.Println(Inverse("[" + "Twitter (Social+)" + "] " + "Now Monitoring Handle: " + foundTwitterHandles[i].Handle))
      database.DatabaseTwitterHandles = append(database.DatabaseTwitterHandles, foundTwitterHandles[i])
      // begin monitoring this Twitter account
      go Twitter(foundTwitterHandles[i].Handle, foundTwitterHandles[i].TwitterApp, foundTwitterHandles[i].ReserveTwitterApp, 864 * time.Millisecond, mongoClient, pusherClient)
    }
  }

  // ########################################### LOOK FOR INSTAGRAM HANDLES TO BEGIN MONITORING
  for i := 0;  i < len(foundInstagramHandles); i++ {
    var foundHandle = false
    for j := 0;  j < len(database.DatabaseInstagramHandles); j++ {
      if database.DatabaseInstagramHandles[j].Handle == foundInstagramHandles[i].Handle {
        foundHandle = true
        break
      }
    }
    if !foundHandle {
      log.Println(Inverse("[" + "Instagram (Social+)" + "] " + "Now Monitoring Handle: " + foundInstagramHandles[i].Handle))
      database.DatabaseInstagramHandles = append(database.DatabaseInstagramHandles, foundInstagramHandles[i])
      // begin monitoring this Instagram account
      go Instagram(foundInstagramHandles[i].Handle, mongoClient, pusherClient)
    }
  }

}
