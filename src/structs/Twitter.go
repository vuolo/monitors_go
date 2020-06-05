package structs

import "github.com/dghubble/go-twitter/twitter"

// ########################################### STRUCTS
type TwitterUser struct {
  Username string `bson:"username"`
  ID string `bson:"id"`
  ImageURL string `bson:"imageurl"`
  FullName string `bson:"fullname"`
  Biography string `bson:"biography"`
  ExternalURL string `bson:"externalurl"`
  Tweets []twitter.Tweet `bson:"tweets"`
}
