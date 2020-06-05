package structs

import(
  "github.com/dghubble/go-twitter/twitter" // twitter
)

// ########################################### STRUCTS
type SocialMediaHandle struct {
  Handle string
  Platform string
  IsFollowingHandle bool
  TwitterApp TwitterApp
  ReserveTwitterApp TwitterApp
	// ConcurrentConnections []ConcurrentConnection
}

type TwitterApp struct {
  ConsumerKey string
  ConsumerSecret string
  AccessTokenKey string
  AccessTokenSecret string
}

type ConcurrentConnection struct {
	DiscordID string
	Timestamp int
}

type TwitterHandleObj struct {
  User TwitterUser
  Tweet twitter.Tweet
}

type InstagramHandleObj struct {
  User InstagramUser
  Post InstagramPost
}

type InstagramHandleObjStory struct {
  User InstagramUser
  Story InstagramStory
}
