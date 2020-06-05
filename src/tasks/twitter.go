package tasks

// ########################################### IMPORTS
import (
	"log"
	"time"
	"context"
	"strings"
	// "math/rand"

	"github.com/dghubble/go-twitter/twitter"
  "github.com/dghubble/oauth1"

	"go.mongodb.org/mongo-driver/mongo" // MongoDB

	"github.com/pusher/pusher-http-go" // Pusher

	. "github.com/logrusorgru/aurora" // colors

	// ## local shits ##
	"../database"
	"../structs"
	"../utils"
)

// ########################################### VARIABLES
func Twitter(twitterHandle string, twitterApp structs.TwitterApp, reserveTwitterApp structs.TwitterApp, timeoutMS time.Duration, mongoClient *mongo.Client, pusherClient *pusher.Client) {

	// ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	identifier := "twitter"
	// timeout := 2 * 864 * time.Millisecond
	timeout := timeoutMS

	// ########################################### VARIABLES

	// Twitter client auth
	config := oauth1.NewConfig(twitterApp.ConsumerKey, twitterApp.ConsumerSecret)
  token := oauth1.NewToken(twitterApp.AccessTokenKey, twitterApp.AccessTokenSecret)
  httpClient := config.Client(oauth1.NoContext, token)

  // Twitter client
  twitterClient := twitter.NewClient(httpClient)

	useReserveApp := false
	// Reserve Twitter client auth
	config_reserve := oauth1.NewConfig(reserveTwitterApp.ConsumerKey, reserveTwitterApp.ConsumerSecret)
  token_reserve := oauth1.NewToken(reserveTwitterApp.AccessTokenKey, reserveTwitterApp.AccessTokenSecret)
  httpClient_reserve := config_reserve.Client(oauth1.NoContext, token_reserve)

  // Reserve Twitter client
  reserveTwitterClient := twitter.NewClient(httpClient_reserve)

  // User Show
  params := &twitter.UserTimelineParams{
    ScreenName: twitterHandle,
    Count: 20,
    IncludeRetweets: twitter.Bool(false),
    ExcludeReplies: twitter.Bool(true),
    TweetMode: "extended",
  }

	for {
		var foundHandle = false
		for i := 0;  i < len(database.DatabaseTwitterHandles); i++ {
			if twitterHandle == database.DatabaseTwitterHandles[i].Handle {
				foundHandle = true
				break
			}
	  }
		if foundHandle {
			scrapeTwitter(twitterHandle, identifier, mongoClient, twitterClient, reserveTwitterClient, params, &useReserveApp, pusherClient)
			time.Sleep(timeout)
		} else {
			break
		}
	}
}

func scrapeTwitter(twitterHandle string, identifier string, mongoClient *mongo.Client, twitterClient *twitter.Client, reserveTwitterClient *twitter.Client, params *twitter.UserTimelineParams, useReserveApp *bool, pusherClient *pusher.Client) {

	// ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// ########################################### VARIABLES
	var tweets []twitter.Tweet

	// ########################################### START REQUEST
	if !*useReserveApp {
		var tl_err error
		tweets, _, tl_err = twitterClient.Timelines.UserTimeline(params)
	  if tl_err != nil {
			log.Println(Yellow("[" + "Twitter: " + twitterHandle + "] " + "Failed fetching user timeline. Using reserve instead. (" + tl_err.Error() + ")"))
			*useReserveApp = true
	    // return
	  }
	} else {
		var tl_err_reserve error
		tweets, _, tl_err_reserve = reserveTwitterClient.Timelines.UserTimeline(params)
		if tl_err_reserve != nil {
			log.Println(Yellow("[" + "Twitter: " + twitterHandle + "] " + "Failed fetching user timeline with reserve. Using assigned twitter app. (" + tl_err_reserve.Error() + ")"))
			*useReserveApp = false
			// time.Sleep(30*time.Second)
			var tl_err error
			tweets, _, tl_err = twitterClient.Timelines.UserTimeline(params)
		  if tl_err != nil {
				log.Println(Red("[" + "Twitter: " + twitterHandle + "] " + "Failed fetching user timeline after switching to assigned twitter app. (" + tl_err.Error() + ")"))
				return
		  }
		}
	}

  if len(tweets) == 0 {
		log.Println(Yellow("[" + "Twitter: " + twitterHandle + "] " + "No tweets found for " + twitterHandle +  "."))
    return
	}

	// ########################################### HANDLE RESPONSE
	var twitterUser structs.TwitterUser
  twitterUser.Username = tweets[0].User.ScreenName
  twitterUser.ID = tweets[0].User.IDStr
  twitterUser.ImageURL = strings.Replace(tweets[0].User.ProfileImageURL, "_normal", "", 1)
  twitterUser.FullName = tweets[0].User.Name
  twitterUser.Biography = tweets[0].User.Description
	if len(tweets[0].User.Entities.URL.Urls) > 0 {
		twitterUser.ExternalURL = tweets[0].User.Entities.URL.Urls[0].ExpandedURL // tweets[0].User.URL
	}
  twitterUser.Tweets = tweets

	// ########################################### INITIAL CHECK
	initialChecked := false
	if utils.StringInSlice("twitter_" + twitterHandle, initialCheckedURLs) {
		log.Println(Green("[" + "Twitter: " + twitterHandle + "] " + "Successful connection."))
		initialChecked = true
	} else {
		// utils.ReverseSlice(twitterUser.Tweets)
		initialCheckedURLs = append(initialCheckedURLs, "twitter_" + twitterHandle)
		log.Println(Inverse("[" + "Twitter: " + twitterHandle + "] " + "Initial Check Done."))
	}
	successful_connections++

	database.SendToTwitterDatabase(twitterUser, identifier, initialChecked, mongoClient, pusherClient)

}
