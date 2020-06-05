package requests

import (
  "encoding/json"
  "log"
  "time"
  "strings"
  "strconv"
  "context"
  "math"
  "net/http"
  // "os"

  "github.com/parnurzeal/gorequest"

  "github.com/pusher/pusher-http-go" // Pusher

  . "github.com/logrusorgru/aurora" // colors

  "github.com/valyala/fasthttp" // fast http
  "github.com/mvdan/xurls" // urls

  "github.com/dghubble/go-twitter/twitter" // twitter

  // ## local shits ##
  "../proxies"
  "../structs"
)

// ########################################### UTILITY VARIABLES
var quicktask_link = "https://resell.monster/quicktask"
var atc_link = "https://resell.monster/atc"
var mlc_link = "https://resell.monster/mlc"

var keywords = []string{
  "+og, +daddy, +yankee",
  "+yeezy,-700,-500,-pant",
  "+air,+fear,+god,-pant,-women,-skylon,-vapormax,-parka,-warm,-shorts",
  "+sacai,-coat,-pant,-bag,-t-shirt,-tee,-shirt,-scarf,-sock",
  "+stranger things",
  "+off,+white,+kiger",
  "+off,+white,+chuck,-women",
  "+cactus,+jack,-women",
  "+jordan,+1,+high,+og,-hyper,-guava,-women,-wheat,-skyhigh",
  "+zx,+4000,+4D,-women",
  "+yeezy,+utility,-500,-pant",
  "+air,+max,+sketch",
  "+jordan,+1,+neutral,+grey,-women,-td",
  "+jordan,+11,+denim",
  "+air,+fear,+god,+1,-skylon,-spruce,-women,-vapormax,-flyknit",
  "+stranger,+things,+blazer",
  "+patta,+jordan,-hat,-tee,-shirt",
  "+jordan,+12,+retro",
  "+aj,+12,+retro",
  "+jordan,+6,+retro",
  "+aj,+6,+retro",
  "+air,+max,+270,+sketch",
  "+air,+jordan,+1,+retro,+high,+og,-strap,-guava,-phantom,-royal,-coral",
  "+travis",
  "+golf,+le,+fleur",
  "+air,+force,+clot,+low,-women,-wmn",
  "+jordan,+1,+fearless",
  "+born,+again,+hooded,-chinatown",
  "+air,+fear,+god,+1,-moc,-raid,-skyl,-pant,-women,-parka,-warm,-short,-top,-sleeve,-hat,-zip,-tee,-jack,-cap,-hat,-hood,-dress",
  "+jordan,+1,+fearless,+blue,+great,-shirt,-short,-pant,-metallic,-womens,-parka,-warm,-short,-top,-sleeve,-hat,-zip,-tee,-jack,-cap,-hat,-hood,-dress",
  "+off,+vapor",
  "+balko",
  "+eve,+aio",
  "+kaws",
  "+off,+white,+dunk",
  "+yeezy,+700,+v3,+azael",
  "+nike,+raygun",
  "+raygun",
  "+air,+force,+clot",
  "+air,+force,+silk",
  "+yeezy, +350, +yeshaya",
  "+nike, +dunk, +plum",
  "+nike, +dunk, +CU1726-500",
  "+nike, +sb, +strangelove",
  "+nike, +sb, +CT2552-800",
  "+air, +jordan, +hi, +85",
  "+air, +jordan, +hi, +BQ4422-600",
  "+yeezy, +700, +mnvn",
  "+yeezy, +700, +FV4440",
}

type Webhook struct {
  Username string `json:"username,omitempty"`
  AvatarURL string `json:"avatar_url,omitempty"`
  Content string `json:"content,omitempty"`
  Embeds []Embed `json:"embeds,omitempty"`
}

type Embed struct {
  Title string `json:"title,omitempty"`
  URL string `json:"url,omitempty"`
  Description string `json:"description,omitempty"`
  Color int `json:"color,omitempty"`
  Footer Footer `json:"footer,omitempty"`
  Image Image `json:"image,omitempty"`
  Thumbnail Thumbnail `json:"thumbnail,omitempty"`
  Author Author `json:"author,omitempty"`
  Fields []Field `json:"fields,omitempty"`
}

type Footer struct {
  Text string `json:"text,omitempty"`
  IconURL string `json:"icon_url,omitempty"`
}

type Image struct {
  URL string `json:"url,omitempty"`
}

type Thumbnail struct {
  URL string `json:"url,omitempty"`
}

type Author struct {
  Name string `json:"name,omitempty"`
  URL string `json:"url,omitempty"`
  IconURL string `json:"icon_url,omitempty"`
}

type Field struct {
  Name string `json:"name,omitempty"`
  Value string `json:"value,omitempty"`
  Inline bool `json:"inline,omitempty"`
}

type WebhookError struct {
  Global bool `json:"global"`
  Message string `json:"message"`
  RetryAfter int `json:"retry_after"`
}

// ########################################### UTILITY FUNCTIONS
func PostProduct(product structs.Product, isNew bool, pusherClient *pusher.Client) {

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // // ################### DISABLE OOS FOR STOCK UPDATE
  // if !isNew && !product.Available {
  //   return
  // }

  // ########################################### VARIABLES
  product.Keywords = grabKeywords(product.Name)

  productTitle := product.Name
  if productTitle == "" {
    productTitle = "Unspecified Product Name"
  }

  productURL := product.URL
  if product.OverrideURL != "" {
    productURL = product.OverrideURL
  }

  if isNew {
    go pusherClient.Trigger("productsChannel", "App\\Events\\addProduct", product)
  } else {
    // if !(product.Store == "juicestore.com" || product.Store == "culturekings.com" || product.Store == "culturekings.co.nz" || product.Store == "culturekings.com.au") {
      go pusherClient.Trigger("productsChannel", "App\\Events\\updateProduct", product)
    // }
  }

  return; // disables product monitors

  webhookURLs := grabWebhookURLs(product.Store, product.Identifier, product.Keywords != "", false, false)

  description := "🆕 NEW"
  if !isNew {
    description = "🔄 STOCK UPDATE"
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  groupColor := 0xd8b97b

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + product.Store,
    IconURL: groupIconURL,
  }

  webhookImage := Image{
    URL: "",
  }

  webhookThumbnail := Thumbnail{
    URL: product.ImageURL,
  }

  webhookAuthor := Author{
    Name: product.StoreName,
    URL: "https://" + product.Store + "#resellmonster",
    // IconURL: product.StoreImageURL,
  }

  webhookFields := setupFields(product)

  webhookOutput := Webhook{
    Username: product.StoreName,
    Embeds: []Embed{
      {
        Title: productTitle,
        URL: productURL,
        Description: description,
        Color: groupColor,
        Footer: webhookFooter,
        Image: webhookImage,
        Thumbnail: webhookThumbnail,
        Author: webhookAuthor,
        Fields: webhookFields,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }
  // log.Println(string(requestBody))

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func UpdateProduct(product structs.Product, pusherClient *pusher.Client) {

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  product.Keywords = grabKeywords(product.Name)

  if !(product.Store == "juicestore.com" || product.Store == "culturekings.com" || product.Store == "culturekings.co.nz" || product.Store == "culturekings.com.au") {
    go pusherClient.Trigger("productsChannel", "App\\Events\\updateProductAvailabilities", product)
  }

}

func PostPassword(store_name string, store_url string, passwordEnabled bool) {

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(store_url, "shopify", false, true, false)

  avatar_url := "https://resell.monster/images/lock.png"
  color := 0xff0000
  description := "🔒 **Password page is UP.**"

  if (!passwordEnabled) {
    avatar_url = "https://resell.monster/images/unlock.png"
    color = 0x98fb98
    description = "🔓 **Password page is DOWN.**"
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + store_url,
    IconURL: groupIconURL,
  }

  webhookOutput := Webhook{
    Username: "PASSWORD MONITOR",
    AvatarURL: avatar_url,
    Embeds: []Embed{
      {
        Title: store_name,
        URL: "https://" + store_url + "#resellmonster",
        Description: description,
        Color: color,
        Footer: webhookFooter,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }
  // log.Println(string(requestBody))

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostCheckpoint(store_name string, store_url string, checkpointEnabled bool) {

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(store_url, "shopify", false, false, true)

  avatar_url := "https://resell.monster/images/checkmark.png"
  color := 0x98fb98;
  description := "✅ **Checkpoint is UP.**"

  if (!checkpointEnabled) {
    avatar_url = "https://resell.monster/images/x-mark.png"
    color = 0xff0000;
    description = "❌ **Checkpoint is DOWN.**"
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + store_url + "/checkpoint",
    IconURL: groupIconURL,
  }

  webhookOutput := Webhook{
    Username: "CHECKPOINT MONITOR",
    AvatarURL: avatar_url,
    Embeds: []Embed{
      {
        Title: store_name,
        URL: "https://" + store_url + "/checkpoint#resellmonster",
        Description: description,
        Color: color,
        Footer: webhookFooter,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }
  // log.Println(string(requestBody))

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostStockX(product structs.StockXProduct, variant structs.StockXVariant, description string) {

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  identifier := "stockx_sneakers"
  if strings.Index(strings.ToLower(product.Collection), "supreme") > -1 {
    identifier = "stockx_supreme"
  }
  webhookURLs := grabWebhookURLs(product.Handle, identifier, false, false, false)

  productTitle := product.Name

  productURL := "https://stockx.com/" + product.Handle + "?size=" + variant.Name

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  groupColor := 0xd8b97b

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + "stockx.com",
    IconURL: groupIconURL,
  }

  webhookThumbnail := Thumbnail{
    URL: product.ImageURL,
  }

  webhookAuthor := Author{
    Name: "StockX",
    URL: "https://stockx.com#resellmonster",
  }

  webhookOutput := Webhook{
    Username: "StockX",
    Embeds: []Embed{
      {
        Title: productTitle,
        URL: productURL,
        Description: description,
        Color: groupColor,
        Footer: webhookFooter,
        Thumbnail: webhookThumbnail,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostTwitterProfilePicture(twitterUser structs.TwitterUser, pusherClient *pusher.Client) {

  var twitterUserNoStatuses = structs.TwitterUser{twitterUser.Username, twitterUser.ID, twitterUser.ImageURL, twitterUser.FullName, twitterUser.Biography, twitterUser.ExternalURL, []twitter.Tweet{}}
  go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newTwitterProfilePicture", twitterUserNoStatuses)

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(twitterUser.Username, "twitter", false, false, false)

  profileURL := "https://twitter.com/" + twitterUser.Username

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0x55acef // Twitter Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookImage := Image{
    URL: twitterUser.ImageURL,
  }

  webhookAuthor := Author{
    Name: twitterUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: twitterUser.ImageURL,
  }

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "PROFILE PICTURE UPDATE",
        Color: groupColor,
        Footer: webhookFooter,
        Image: webhookImage,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostTwitterFullName(twitterUser structs.TwitterUser, pusherClient *pusher.Client) {

  var twitterUserNoStatuses = structs.TwitterUser{twitterUser.Username, twitterUser.ID, twitterUser.ImageURL, twitterUser.FullName, twitterUser.Biography, twitterUser.ExternalURL, []twitter.Tweet{}}
  go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newTwitterFullName", twitterUserNoStatuses)

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(twitterUser.Username, "twitter", false, false, false)

  profileURL := "https://twitter.com/" + twitterUser.Username
  webhookDescription := twitterUser.FullName
  if webhookDescription == "" {
    webhookDescription = "User has removed their display name."
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0x55acef // Twitter Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookAuthor := Author{
    Name: twitterUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: twitterUser.ImageURL,
  }

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "DISPLAY NAME UPDATE",
        Description: webhookDescription,
        Color: groupColor,
        Footer: webhookFooter,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostTwitterBiography(twitterUser structs.TwitterUser, pusherClient *pusher.Client) {

  var twitterUserNoStatuses = structs.TwitterUser{twitterUser.Username, twitterUser.ID, twitterUser.ImageURL, twitterUser.FullName, twitterUser.Biography, twitterUser.ExternalURL, []twitter.Tweet{}}
  go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newTwitterBiography", twitterUserNoStatuses)

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(twitterUser.Username, "twitter", false, false, false)

  profileURL := "https://twitter.com/" + twitterUser.Username
  webhookDescription := twitterUser.Biography
  if webhookDescription == "" {
    webhookDescription = "User has removed their biography."
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0x55acef // Twitter Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookAuthor := Author{
    Name: twitterUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: twitterUser.ImageURL,
  }

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "BIOGRAPHY UPDATE",
        Description: webhookDescription,
        Color: groupColor,
        Footer: webhookFooter,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostTwitterExternalURL(twitterUser structs.TwitterUser, pusherClient *pusher.Client) {

  var twitterUserNoStatuses = structs.TwitterUser{twitterUser.Username, twitterUser.ID, twitterUser.ImageURL, twitterUser.FullName, twitterUser.Biography, twitterUser.ExternalURL, []twitter.Tweet{}}
  go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newTwitterExternalURL", twitterUserNoStatuses)

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(twitterUser.Username, "twitter", false, false, false)

  profileURL := "https://twitter.com/" + twitterUser.Username
  webhookDescription := twitterUser.ExternalURL
  if webhookDescription == "" {
    webhookDescription = "User has removed their displayed url."
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0x55acef // Twitter Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookAuthor := Author{
    Name: twitterUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: twitterUser.ImageURL,
  }

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "DISPLAYED URL UPDATE",
        Description: webhookDescription,
        Color: groupColor,
        Footer: webhookFooter,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostTwitterTweet(twitterUser structs.TwitterUser, tweet twitter.Tweet, pusherClient *pusher.Client) {

  var twitterUserNoStatuses = structs.TwitterUser{twitterUser.Username, twitterUser.ID, twitterUser.ImageURL, twitterUser.FullName, twitterUser.Biography, twitterUser.ExternalURL, []twitter.Tweet{}}
  jsonResult, jsonErr := json.Marshal(&structs.TwitterHandleObj{twitterUserNoStatuses, tweet})
  if jsonErr != nil {
    log.Fatal(jsonErr)
  } else {
    go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newTwitterStatus", jsonResult)
  }

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(twitterUser.Username, "twitter", false, false, false)

  profileURL := "https://twitter.com/" + twitterUser.Username
  tweetURL := profileURL + "/status/" + tweet.IDStr
  tweetText := ""
  if tweet.FullText != "" {
    tweetText = tweet.FullText
    if tweet.Entities != nil {
      for i := 0;  i < len(tweet.Entities.Urls); i++ {
        tweetText = strings.ReplaceAll(tweetText, tweet.Entities.Urls[i].URL, tweet.Entities.Urls[i].ExpandedURL)
      }
      for i := 0;  i < len(tweet.Entities.Media); i++ {
        tweetText = strings.ReplaceAll(tweetText, tweet.Entities.Media[i].URLEntity.URL, "")
      }
    } else if tweet.ExtendedTweet != nil && tweet.ExtendedTweet.Entities != nil {
      for i := 0;  i < len(tweet.ExtendedTweet.Entities.Urls); i++ {
        tweetText = strings.ReplaceAll(tweetText, tweet.Entities.Urls[i].URL, tweet.Entities.Urls[i].ExpandedURL)
      }
      for i := 0;  i < len(tweet.ExtendedTweet.Entities.Media); i++ {
        tweetText = strings.ReplaceAll(tweetText, tweet.Entities.Media[i].URLEntity.URL, "")
      }
    }
    if tweet.QuotedStatus != nil {
      quotedProfileURL := "https://twitter.com/" + tweet.QuotedStatus.User.ScreenName
      quotedTweetURL := quotedProfileURL + "/status/" + tweet.QuotedStatusIDStr
      tweetText = strings.ReplaceAll(tweetText, quotedTweetURL, "\n\n**[QUOTED TWEET](" + quotedTweetURL + ")**:\n```" + tweet.QuotedStatus.User.Name + " (@" + tweet.QuotedStatus.User.ScreenName + "):\n" + tweet.QuotedStatus.FullText + "```")
    }
  }

  tweetImageURL := ""
  if tweet.Entities != nil && len(tweet.Entities.Media) > 0 {
    tweetImageURL = tweet.Entities.Media[0].MediaURL
  } else if tweet.ExtendedTweet != nil && tweet.ExtendedTweet.Entities != nil && len(tweet.ExtendedTweet.Entities.Media) > 0 {
    tweetImageURL = tweet.ExtendedTweet.Entities.Media[0].MediaURL
  }

  webhookContent := ""
  lowercaseContent := strings.ToLower(tweetText)
  if strings.Index(lowercaseContent, "giveaway") > -1 || (strings.Index(lowercaseContent, "giv") > -1 && strings.Index(lowercaseContent, "away") > -1) || strings.Index(lowercaseContent, "free") > -1 {
    webhookContent = "<@&661297580789071873> Potential Giveaway Detected! 🎁";
  } else if (strings.Index(lowercaseContent, "restock") > -1 || strings.Index(lowercaseContent, "fcfs") > -1) && !(strings.Index(lowercaseContent, "restock world") > -1 || strings.Index(lowercaseContent, "restock world") > -1) {
    webhookContent = "<@&661297590435971093> Potential Restock Detected! 💸";
  }
  rawURLs := xurls.Relaxed().FindAllString(tweetText, -1)
  for i := 0;  i < len(rawURLs); i++ {
    if strings.Index(rawURLs[i], "discordapp.com") > -1 || strings.Index(rawURLs[i], "discord.gg") > -1 {
      if webhookContent == "" {
        webhookContent += rawURLs[i] + " "
      } else {
          webhookContent += "• " + rawURLs[i]
      }
    }
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0x55acef // Twitter Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookImage := Image{
    URL: tweetImageURL,
  }

  webhookAuthor := Author{
    Name: twitterUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: twitterUser.ImageURL,
  }

  webhookFields := setupTwitterFields(tweet)

  webhookOutput := Webhook{
    Content: webhookContent,
    Embeds: []Embed{
      {
        Title: "LINK TO ORIGINAL TWEET",
        URL: tweetURL,
        Description: tweetText,
        Color: groupColor,
        Footer: webhookFooter,
        Image: webhookImage,
        Author: webhookAuthor,
        Fields: webhookFields,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostInstagramProfilePicture(instagramUser structs.InstagramUser, pusherClient *pusher.Client) {

  var instagramUserNoStatuses = structs.InstagramUser{instagramUser.Username, instagramUser.ID, instagramUser.ImageURL, instagramUser.FullName, instagramUser.Biography, instagramUser.ExternalURL, []structs.InstagramPost{}, []structs.InstagramStory{}}
  go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newInstagramProfilePicture", instagramUserNoStatuses)

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(instagramUser.Username, "instagram", false, false, false)

  profileURL := "https://instagram.com/" + instagramUser.Username

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0xC13584 // Instagram Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookImage := Image{
    URL: instagramUser.ImageURL,
  }

  webhookAuthor := Author{
    Name: instagramUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: instagramUser.ImageURL,
  }

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "PROFILE PICTURE UPDATE",
        Color: groupColor,
        Footer: webhookFooter,
        Image: webhookImage,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostInstagramFullName(instagramUser structs.InstagramUser, pusherClient *pusher.Client) {

  var instagramUserNoStatuses = structs.InstagramUser{instagramUser.Username, instagramUser.ID, instagramUser.ImageURL, instagramUser.FullName, instagramUser.Biography, instagramUser.ExternalURL, []structs.InstagramPost{}, []structs.InstagramStory{}}
  go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newInstagramFullName", instagramUserNoStatuses)

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(instagramUser.Username, "instagram", false, false, false)

  profileURL := "https://instagram.com/" + instagramUser.Username
  webhookDescription := instagramUser.FullName
  if webhookDescription == "" {
    webhookDescription = "User has removed their display name."
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0xC13584 // Instagram Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookAuthor := Author{
    Name: instagramUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: instagramUser.ImageURL,
  }

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "DISPLAY NAME UPDATE",
        Description: webhookDescription,
        Color: groupColor,
        Footer: webhookFooter,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostInstagramBiography(instagramUser structs.InstagramUser, pusherClient *pusher.Client) {

  var instagramUserNoStatuses = structs.InstagramUser{instagramUser.Username, instagramUser.ID, instagramUser.ImageURL, instagramUser.FullName, instagramUser.Biography, instagramUser.ExternalURL, []structs.InstagramPost{}, []structs.InstagramStory{}}
  go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newInstagramBiography", instagramUserNoStatuses)

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(instagramUser.Username, "instagram", false, false, false)

  profileURL := "https://instagram.com/" + instagramUser.Username
  webhookDescription := instagramUser.Biography
  if webhookDescription == "" {
    webhookDescription = "User has removed their biography."
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0xC13584 // Instagram Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookAuthor := Author{
    Name: instagramUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: instagramUser.ImageURL,
  }

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "BIOGRAPHY UPDATE",
        Description: webhookDescription,
        Color: groupColor,
        Footer: webhookFooter,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostInstagramExternalURL(instagramUser structs.InstagramUser, pusherClient *pusher.Client) {

  var instagramUserNoStatuses = structs.InstagramUser{instagramUser.Username, instagramUser.ID, instagramUser.ImageURL, instagramUser.FullName, instagramUser.Biography, instagramUser.ExternalURL, []structs.InstagramPost{}, []structs.InstagramStory{}}
  go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newInstagramExternalURL", instagramUserNoStatuses)

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(instagramUser.Username, "instagram", false, false, false)

  profileURL := "https://instagram.com/" + instagramUser.Username
  webhookDescription := instagramUser.ExternalURL
  if webhookDescription == "" {
    webhookDescription = "User has removed their displayed url."
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0xC13584 // Instagram Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookAuthor := Author{
    Name: instagramUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: instagramUser.ImageURL,
  }

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "DISPLAYED URL UPDATE",
        Description: webhookDescription,
        Color: groupColor,
        Footer: webhookFooter,
        Author: webhookAuthor,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostInstagramPost(instagramUser structs.InstagramUser, post structs.InstagramPost, pusherClient *pusher.Client) {

  var instagramUserNoStatuses = structs.InstagramUser{instagramUser.Username, instagramUser.ID, instagramUser.ImageURL, instagramUser.FullName, instagramUser.Biography, instagramUser.ExternalURL, []structs.InstagramPost{}, []structs.InstagramStory{}}
  jsonResult, jsonErr := json.Marshal(&structs.InstagramHandleObj{instagramUserNoStatuses, post})
  if jsonErr != nil {
    log.Fatal(jsonErr)
  } else {
    go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newInstagramStatus", jsonResult)
  }

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(instagramUser.Username, "instagram", false, false, false)

  profileURL := "https://instagram.com/" + instagramUser.Username
  postURL := "https://instagram.com/p/" + post.ID

  webhookContent := ""
  lowercaseContent := strings.ToLower(post.Caption)
  if strings.Index(lowercaseContent, "giveaway") > -1 || (strings.Index(lowercaseContent, "giv") > -1 && strings.Index(lowercaseContent, "away") > -1) || strings.Index(lowercaseContent, "free") > -1 {
    webhookContent = "<@&661297580789071873> Potential Giveaway Detected! 🎁";
  } else if (strings.Index(lowercaseContent, "restock") > -1 || strings.Index(lowercaseContent, "fcfs") > -1) && !(strings.Index(lowercaseContent, "restock world") > -1 || strings.Index(lowercaseContent, "restock world") > -1) {
    webhookContent = "<@&661297590435971093> Potential Restock Detected! 💸";
  }
  rawURLs := xurls.Relaxed().FindAllString(post.Caption, -1)
  for i := 0;  i < len(rawURLs); i++ {
    if strings.Index(rawURLs[i], "discordapp.com") > -1 || strings.Index(rawURLs[i], "discord.gg") > -1 {
      if webhookContent == "" {
        webhookContent += rawURLs[i] + " "
      } else {
          webhookContent += " • " + rawURLs[i]
      }
    }
  }

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0xC13584 // Instagram Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookImage := Image{
    URL: post.ImageURL,
  }

  webhookAuthor := Author{
    Name: instagramUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: instagramUser.ImageURL,
  }

  webhookFields := setupInstagramFields(post)

  webhookOutput := Webhook{
    Content: webhookContent,
    Embeds: []Embed{
      {
        Title: "LINK TO ORIGINAL POST",
        URL: postURL,
        Description: post.Caption,
        Color: groupColor,
        Footer: webhookFooter,
        Image: webhookImage,
        Author: webhookAuthor,
        Fields: webhookFields,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostInstagramStory(instagramUser structs.InstagramUser, story structs.InstagramStory, pusherClient *pusher.Client) {

  var instagramUserNoStatuses = structs.InstagramUser{instagramUser.Username, instagramUser.ID, instagramUser.ImageURL, instagramUser.FullName, instagramUser.Biography, instagramUser.ExternalURL, []structs.InstagramPost{}, []structs.InstagramStory{}}
  jsonResult, jsonErr := json.Marshal(&structs.InstagramHandleObjStory{instagramUserNoStatuses, story})
  if jsonErr != nil {
    log.Fatal(jsonErr)
  } else {
    go pusherClient.Trigger("socialPlusChannel", "App\\Events\\newInstagramStory", jsonResult)
  }

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
  webhookURLs := grabWebhookURLs(instagramUser.Username, "instagram", false, false, false)

  profileURL := "https://instagram.com/" + instagramUser.Username
  storyURL := "https://instagram.com/stories/" + instagramUser.Username + "/" + story.ID

  // groupIconURL := "https://dashpings.com/DashPings.png"
  // groupText := "Powered by Dash Pings"
  // groupColor := 0xd8b97b
  groupIconURL := "https://resell.monster/images/buddy-transparent-compressed.png"
  groupText := "Resell Monster"
  // groupColor := 0xd8b97b

  groupColor := 0xC13584 // Instagram Color

  webhookFooter := Footer{
    Text: groupText + " \u2022 " + time.Now().Format("January 2, 2006 (3:04:05 PM) [MST]") + " \u2022 " + strings.Replace(profileURL, "https://", "", 1),
    IconURL: groupIconURL,
  }

  webhookImage := Image{
    URL: story.ImageURL,
  }

  webhookAuthor := Author{
    Name: instagramUser.FullName,
    URL: profileURL + "#resellmonster",
    IconURL: instagramUser.ImageURL,
  }

  webhookFields := setupInstagramStoryFields(story)

  webhookOutput := Webhook{
    Embeds: []Embed{
      {
        Title: "LINK TO ORIGINAL STORY",
        URL: storyURL,
        Color: groupColor,
        Footer: webhookFooter,
        Image: webhookImage,
        Author: webhookAuthor,
        Fields: webhookFields,
      },
    },
  }

  requestBody, err := json.Marshal(webhookOutput)
  if err != nil {
    log.Println(err)
    return
  }

  for i := 0;  i < len(webhookURLs); i++ {
    go tryPostWebhook(requestBody, webhookURLs[i])
  }

}

func PostTest(webhookURL string) {

  // ########################################### VARIABLES
  requestBody, err := json.Marshal(map[string]string{
    "content": "testing testing 123...",
  })
  if err != nil {
    log.Println(err)
    return
  }

  go tryPostWebhook(requestBody, webhookURL)

}

func postWebhook(requestBody []byte, webhookURL string) bool {
  // return true // disables ALL discord webhooks from being sent
  var strPost = []byte("POST")
  var strRequestURI = []byte(webhookURL)

  req := fasthttp.AcquireRequest()
  req.SetBody(requestBody)
  req.Header.SetMethodBytes(strPost)
  req.Header.SetContentType("application/json")
  req.SetRequestURIBytes(strRequestURI)
  res := fasthttp.AcquireResponse()
  if err := fasthttp.Do(req, res); err != nil {
    log.Println(err)
    log.Println(req)
    log.Println(res)
    time.Sleep(5 * time.Second)
    return false
  }
  defer fasthttp.ReleaseRequest(req)

  body := res.Body()

  if len(string(body)) > 0 {
    if strings.Index(string(body), "error code: 1015") > -1 {
      log.Println(Yellow("POST ERROR: Code: " + strconv.Itoa(res.StatusCode()) + " - Message: " + string(http.StatusText(res.StatusCode())) + " )"))
      log.Println(Yellow(string(body)))
      request := gorequest.New()
      goreq_req, goreq_body, gorequest_err := request.Proxy(proxies.GrabProxy()).Post(webhookURL).Send(string(requestBody)).EndBytes()
      if gorequest_err != nil {
        log.Println(Yellow("ATTEMPT 2 POST with proxy ERRORED"))
      } else {
        log.Println(Yellow("ATTEMPT 2 POST with proxy SUCCEEDED."))
        log.Println(Yellow(goreq_req.StatusCode))
        log.Println(Yellow(string(goreq_body)))
        return true
      }
      time.Sleep(5 * time.Second)
      return false
    } else {
      log.Println(Yellow(string(body)))
    }
  }

  defer fasthttp.ReleaseResponse(res) // Only when you are done with body!

  if strings.Index(string(body), "You are being rate limited.") > -1 {
    log.Println(Yellow(string(body)))
    data := &WebhookError{}
    err := json.Unmarshal(body, data)
    if err != nil {
      log.Println(err)
      return false
    }
    time.Sleep(time.Duration(data.RetryAfter) * time.Millisecond)
    return false
  } else if strings.Index(string(body), "400 Bad Request") > -1 {
    log.Println(Yellow(string(body)))
    time.Sleep(5 * time.Second)
    return false
  }
  return true
}

func tryPostWebhook(requestBody []byte, webhookURL string) {

  // keep retrying post until it works
  for {
    if postWebhook(requestBody, webhookURL) {
      break
    } else {
      log.Println(string(requestBody))
      log.Println(webhookURL)
    }
  }

}

func setupFields(product structs.Product) []Field {
  var fields []Field
  if (product.Identifier != "shopify" && product.Identifier != "cpfm") && product.LaunchDate != "" {
    launchDate, _ := time.Parse(time.RFC3339, product.LaunchDate)
    fields = append(fields, Field{
      Name: "**Launch Date**",
      Value: launchDate.Format("January 2, 2006 (3:04:05 PM) [MST]"),
      Inline: true,
    })
  }
  if product.Price != "" {
    fields = append(fields, Field{
      Name: "**Price**",
      Value: product.Price,
      Inline: true,
    })
  }
  if len(product.Variants) > 0 {
    if product.Available {
      if product.Variants[0].Quantity == -420 { // Quantity not available
        totalStock := 0
        for i := 0;  i < len(product.Variants); i++ {
          if product.Variants[i].Available {
            totalStock++
          }
        }
        fields = append(fields, Field{
          Name: "**Total Stock**",
          Value: strconv.Itoa(totalStock) + "+",
          Inline: true,
        })
      } else { // Quantity available: aggregate all stock levels
        totalStock := 0
        for i := 0;  i < len(product.Variants); i++ {
          if product.Variants[i].Quantity != -420 {
            totalStock += product.Variants[i].Quantity
          }
        }
        fields = append(fields, Field{
          Name: "**Total Stock**",
          Value: strconv.Itoa(totalStock),
          Inline: true,
        })
      }
    } else { // Not available: display "SOLD OUT"
      fields = append(fields, Field{
        Name: "**Total Stock**",
        Value: "~~**SOLD OUT**~~",
        Inline: true,
      })
    }
  }
  if product.Color != "" {
    fields = append(fields, Field{
      Name: "**Color**",
      Value: product.Color,
      Inline: true,
    })
  }
  if product.Collection != "" {
    fields = append(fields, Field{
      Name: "**Collection**",
      Value: product.Collection,
      Inline: true,
    })
  }
  if len(product.Variants) > 0 {
    if product.Available {
      sizesString := ""
      for i := 0;  i < len(product.Variants); i++ {
        if !product.Variants[i].Available {
          continue
        }
        if product.Identifier == "supreme" { // supreme atc: http://localhost:420/atc?store=supreme&pid=172919&id=74706
          choppedURL := product.OverrideURL[strings.Index(product.OverrideURL, "new/"):]
          atcLink := atc_link + "?store=supreme" + "&pid=" + choppedURL[strings.Index(choppedURL, "/")+1:strings.LastIndex(choppedURL, "/")] + "&id=" + product.Variants[i].ID
          sizesString += "[**" + product.Variants[i].Name + "** (ATC)](" + atcLink + ")" + " [" + strconv.Itoa(product.Variants[i].Quantity) + "]" + "\n"
        } else if product.Identifier == "offwhite" { // offwhite atc: http://localhost:420/atc?store=off-white&id=118025
          atcLink := atc_link + "?store=off-white" + "&id=" + product.Variants[i].ID
          sizesString += "[**" + product.Variants[i].Name + "** (ATC)](" + atcLink + ")" + " [" + strconv.Itoa(product.Variants[i].Quantity) + "]" + "\n"
        } else if product.Identifier == "offspring" { // offspring atc: http://localhost:420/atc?store=offspring&id=3760396114060
          atcLink := atc_link + "?store=offspring" + "&id=" + product.Variants[i].ID
          sizesString += "[**" + product.Variants[i].Name + "** (ATC)](" + atcLink + ")" + " [" + strconv.Itoa(product.Variants[i].Quantity) + "]" + "\n"
        } else if product.Identifier == "solebox" { // solebox atc: https://www.solebox.com/index.php?aproducts[0][am]=1&fnc=changebasket&cl=basket&action=atc&aproducts[0][aid]=44053
          atcLink := "https://www.solebox.com/index.php?aproducts[0][am]=1&fnc=changebasket&cl=basket&action=atc&aproducts[0][aid]=" + product.Variants[i].ID
          sizesString += "[**" + product.Variants[i].Name + "** (ATC)](" + atcLink + ")" + " [" + strconv.Itoa(product.Variants[i].Quantity) + "]" + "\n"
        } else if product.Identifier == "cpfm" { // cpfm atc: http://localhost:420/atc?store=cpfm&id=Z2lkOi8vc2hvcGlmeS9Qcm9kdWN0VmFyaWFudC8yODIxNjE5MjQ5OTc5Mg==
          atcLink := atc_link + "?store=cpfm" + "&id=" + product.Variants[i].ID
          sizesString += "[**" + product.Variants[i].Name + "** (ATC)](" + atcLink + ")" + " [" + strconv.Itoa(product.Variants[i].Quantity) + "]" + "\n"
        } else if product.Identifier == "shopify" || product.Identifier == "cpfm" || product.Identifier == "dsm-eflash" { // product (shopify) atc: http://localhost:420/atc?url=https://solestop.com/cart/31141503238186:1
          atcLink := "https://" + product.Store + "/cart/" + product.Variants[i].ID + ":1"
          qtLink := atc_link + "?url=" + atcLink
          sizesString += "**[" + product.Variants[i].Name + "](" + atcLink + ")** [(QT)](" + qtLink + ")" + " [" + strconv.Itoa(product.Variants[i].Quantity) + "]" + "\n"
        } else {
          sizesString += "**" + product.Variants[i].Name + "**" + " (" + product.Variants[i].ID + ")" + " [" + strconv.Itoa(product.Variants[i].Quantity) + "]" + "\n"
        }
      }
      sizesString = strings.ReplaceAll(sizesString, "[-420]", "[1+]")
      fieldsNeeded := int(math.Floor(float64(len(sizesString)) / 1024) + 1)
      sizesArrayRaw := strings.Split(sizesString, "\n")
      sizesArray := chunkIt(sizesArrayRaw[:len(sizesArrayRaw)-1], fieldsNeeded, true) // slice sizesArray to remove last empty element
      useInline := false
      if len(fields) % 3 == 0 {
        useInline = true
      }
      for i := 0; i < len(sizesArray); i++ {
        if i == 0 { // first size chunk
          fields = append(fields, Field{
            Name: "**Sizes**",
            Value: strings.Join(sizesArray[i], "\n"),
            Inline: useInline,
          })
        } else { // all extra size chunks
          fields = append(fields, Field{
            Name: "**⠀**",
            Value: strings.Join(sizesArray[i], "\n"),
            Inline: useInline,
          })
        }
      }
    }
  }
  if product.Keywords != "" {
    fields = append(fields, Field{
      Name: "**Keywords**",
      Value: "`" + product.Keywords + "`",
      Inline: false,
    })
  }
  fields = append(fields, Field{
    Name: "**Links**",
    Value: "[StockX](https://stockx.com/search?s=" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(product.Name, " ", "+"), "(", "%28"), ")", "%29") + ") | [QuickTask Setup](" + quicktask_link + ") | [Mass Link Change](" + mlc_link + "?url=" + product.URL + ") | **[Execute QT](" + atc_link + "?url=" + product.URL + ")**",
    Inline: false,
  })

  return fields
}

func setupTwitterFields(tweet twitter.Tweet) []Field {
  var fields []Field
  URLs := []string{}
  if tweet.Entities != nil {
    for i := 0;  i < len(tweet.Entities.Urls); i++ {
      URLs = append(URLs, "• " + tweet.Entities.Urls[i].ExpandedURL)
    }
  } else if tweet.ExtendedTweet != nil && tweet.ExtendedTweet.Entities != nil {
    for i := 0;  i < len(tweet.ExtendedTweet.Entities.Urls); i++ {
      URLs = append(URLs, "• " + tweet.ExtendedTweet.Entities.Urls[i].ExpandedURL)
    }
  }

  if tweet.QuotedStatus != nil {
    // quote url test
    for i := 0;  i < len(URLs); i++ {
      if strings.Contains(URLs[i], tweet.QuotedStatus.IDStr) {
        URLs = append(URLs[:i], URLs[i+1:]...)
        break
      }
    }
  }

  // self url test
  for i := 0;  i < len(URLs); i++ {
    if strings.Contains(URLs[i], tweet.IDStr) {
      URLs = append(URLs[:i], URLs[i+1:]...)
      break
    }
  }

  if len(URLs) > 0 {
    fields = append(fields, Field{
      Name: "**URLs**",
      Value: strings.Join(URLs, "\n"),
      Inline: false,
    })
  }

  return fields
}

func setupInstagramFields(post structs.InstagramPost) []Field {
  var fields []Field
  rawURLs := xurls.Relaxed().FindAllString(post.Caption, -1)
  URLs := []string{}
  for i := 0;  i < len(rawURLs); i++ {
    URLs = append(URLs, "• " + rawURLs[i])
  }

  if post.IsVideo {
    fields = append(fields, Field{
      Name: "**Video?**",
      Value: "✅",
      Inline: true,
    })
  }

  if len(URLs) > 0 {
    fields = append(fields, Field{
      Name: "**URLs**",
      Value: strings.Join(URLs, "\n"),
      Inline: false,
    })
  }

  return fields
}

func setupInstagramStoryFields(story structs.InstagramStory) []Field {
  var fields []Field

  if story.IsVideo {
    fields = append(fields, Field{
      Name: "**Video?**",
      Value: "✅",
      Inline: true,
    })
  }

  if story.ExternalURL != "" {
    fields = append(fields, Field{
      Name: "**External URL**",
      Value: story.ExternalURL,
      Inline: true,
    })
  }

  return fields
}

func chunkIt(sizesArray []string, numChunks int, balanced bool) [][]string {
  len := len(sizesArray)
  var out [][]string
  i := 0
  size := 0

  if numChunks < 2 {
    out = append(out, sizesArray)
    return out
  }

  if len % numChunks == 0 {
    size = int(math.Floor(float64(len / numChunks)))
    for i < len {
      out = append(out, sizesArray[i:i+size])
      i+=size
    }
  } else if balanced {
    for i < len {
      numChunks--
      size = int(math.Ceil(float64((len - i) / numChunks)))
      out = append(out, sizesArray[i:i+size])
      i+=size
    }
  } else {
    numChunks--
    size = int(math.Floor(float64(len / numChunks)))
    if len % size == 0 {
      size--
    }
    for i < size * numChunks {
      out = append(out, sizesArray[i:i+size])
      i+=size
    }
    out = append(out, sizesArray[size*numChunks:])
  }

  return out
}

func grabWebhookURLs(store_url string, identifier string, isFiltered bool, isPassword bool, isCheckpoint bool) []string {
  webhookURLs := []string{}
  if identifier == "shopify" {
    if store_url == "travis-scott-secure.myshopify.com" {
      webhookURL := "https://discordapp.com/api/webhooks/645146366103912448/PPTHVvA5pnK15GKZ9h0wdBEdA4r-58MzHUKtK3r5ndYxNt1Wjy3xNIhT-xt28VtJlD8b"
      webhookURLs = append(webhookURLs, webhookURL) // Travis Scott
    } else if store_url == "eflash-jp.doverstreetmarket.com" ||
              store_url == "eflash-sg.doverstreetmarket.com" ||
              store_url == "eflash-us.doverstreetmarket.com" ||
              store_url == "eflash.doverstreetmarket.com" {
      webhookURL := "https://discordapp.com/api/webhooks/609528780251332610/268vHAw6g1EZPWU_s7W4pHDwI8Za-j0NIEU5YtwlXtpXj7fbgb0D6ZLKmW5Gg8HIqoqm"
      webhookURLs = append(webhookURLs, webhookURL) // Dover Street Market (E-FLASH)
    } else if store_url == "shop-usa.palaceskateboards.com" {
      webhookURL := "https://discordapp.com/api/webhooks/668146855263207474/EIhRhh_2Yy5m8SU2R99gjiLX-4QfJWqpoXGWu9OtDdKsRdnm7hiOeYffNR5zbZy2uGnF"
      webhookURLs = append(webhookURLs, webhookURL) // Palace
    } else if store_url == "undefeated.com" {
      webhookURL := "https://discordapp.com/api/webhooks/605553457616781312/1VIGSRQ_cZE0s9Yf9lmZJKcn2uPGUecI9K26Zy-Idm4Y5FYnw-39hJoZ2Fu56Bb2yk5L"
      webhookURLs = append(webhookURLs, webhookURL) // Undefeated
    } else if store_url == "bdgastore.com" {
      webhookURL := "https://discordapp.com/api/webhooks/605553561115295784/UbTOy1e1gw8Bnqk535bx926tkme5hX5gwQabW8FiI1BLRa-sIO-76VoflYuHPU7NSnNR"
      webhookURLs = append(webhookURLs, webhookURL) // Bodega
    } else if store_url == "us.bape.com" {
      webhookURL := "https://discordapp.com/api/webhooks/605518906110640171/TGBsjqvHd4OV6GSosGwocvqHSsjZL_9smWu0WsJEPANlRuiRMxgt4-Ys-KGaiYkT32Zy"
      webhookURLs = append(webhookURLs, webhookURL) // BAPE
    } else if store_url == "kith.com" {
      webhookURL := "https://discordapp.com/api/webhooks/605551872576716801/syBU1DMEVYrqI6ATSGbfecK6JGxDo50fvlF4NTodbjMmo_U-0jhsh3oNarobPgaKCax2"
      webhookURLs = append(webhookURLs, webhookURL) // KITH
    } else if store_url == "everoboticsinc.com" ||
              store_url == "f3ather.io" ||
              store_url == "purchase.spectrebots.com" ||
              store_url == "shop.balkobot.com" ||
              store_url == "shop.destroyerbots.com" ||
              store_url == "shop.ghostaio.com" ||
              store_url == "soleaio.com" {
      webhookURL := "https://discordapp.com/api/webhooks/651216896292552715/yH6AGUVB2A90Y7Qqg5PsB8GV0PUz2kKFfMlMzXoh5rs3bRj_MoFBo9qIli1sRZ84ulc-"
      webhookURLs = append(webhookURLs, webhookURL) // bots
    }
    if isCheckpoint {
      webhookURL := "https://discordapp.com/api/webhooks/659514744394481681/rViouzYT_MrADVIcYiNh7WxOHpLy4kELwd8xvBfxKfSrL_u4RMLoO3e22cJq8zMluazB"
      webhookURLs = append(webhookURLs, webhookURL) // Shopify Checkpoint
    }
    if isPassword {
      webhookURL := "https://discordapp.com/api/webhooks/609573294282113024/mM5YW-Hy7GpYvhO7I7OIORCmmmv-WKpESzCUzeQXNy_Dpu96ynH4r5B5IhFf-xT3DwQX"
      webhookURLs = append(webhookURLs, webhookURL) // Shopify Password
    }
    if isFiltered {
      webhookURL := "https://discordapp.com/api/webhooks/609295924639825930/cN1-mvN4v6C-pjkXqIuAcqOtFclPZ90piKt6lrQAc6Jh3d323vksh3H2zzyzQgmR6K8Z"
      webhookURLs = append(webhookURLs, webhookURL) // Shopify Filtered
    }
    webhookURL := "https://discordapp.com/api/webhooks/609569084065185792/KS4GWVNDYeq1JEDtWf4TVG2drYcjPKtLyGKa64M7cXgkP65PGmgUDWEu0swVL1QJIIwA"
    webhookURLs = append(webhookURLs, webhookURL) // Shopify Unfiltered
  } else if identifier == "cpfm" {
    webhookURL := "https://discordapp.com/api/webhooks/644303769173229578/EDKg9WWp_fB79RbC4GxysGeDpRRDgyg6gDYc_rpyjKk-GGbIkC4gQ8_Wd1gRiKt9nlNA"
    webhookURLs = append(webhookURLs, webhookURL) // CPFM
  } else if identifier == "stockx_sneakers" {
    webhookURL := "https://discordapp.com/api/webhooks/661352850907070496/YhXFh_jf-KQyFZ2xIVAW-kASHZyFkXAIjKho9J5RBDZFd8V7AAlCG-ilohVsNOgcF592"
    webhookURLs = append(webhookURLs, webhookURL) // StockX Sneakers
  } else if identifier == "stockx_supreme" {
    webhookURL := "https://discordapp.com/api/webhooks/661352615107362843/zQ4XNOIafVeSoUuS06sT1_sZOFd5TuvYDuKkuvFXFTwcUs1dExPNBefPQl8pHaPqOr6p"
    webhookURLs = append(webhookURLs, webhookURL) // StockX Supreme
  } else if identifier == "twitter" {
    if strings.ToLower(store_url) == "cybersole" || strings.ToLower(store_url) == "offline" {
      webhookURL := "https://discordapp.com/api/webhooks/651643248141533206/textITlu86932zh1zjDsHG17eZ4JBPqzyjTr3mBKfTjgsojIvCMoHuph9zcQ2aZT_OJx"
      webhookURLs = append(webhookURLs, webhookURL) // Twitter Cybersole
    } else if strings.ToLower(store_url) == "balkobot" {
      webhookURL := "https://discordapp.com/api/webhooks/651643763961364501/KYur4389bxW_kD0s8Rc4gv7HN-sapLAos2bb3hpkEBH_Y4ZnEMjQavh3dKj_mMfdAsyV"
      webhookURLs = append(webhookURLs, webhookURL) // Twitter balkobot
    } else if strings.ToLower(store_url) == "adeptbots" {
      webhookURL := "https://discordapp.com/api/webhooks/651644071856701450/aFHaNiuPgWryyhVCv6sU1lFpwxkQGebSXjM8Pj-czJNi4zJzPv6emdP3gjRkhBN42ClH"
      webhookURLs = append(webhookURLs, webhookURL) // Twitter Adept
    } else if strings.ToLower(store_url) == "nova_aio" {
      webhookURL := "https://discordapp.com/api/webhooks/651644277469741057/sb8oR11rWgtqQ4xjFI6ih0ubj5kSp6IIPCv9dsrAsrq9dXjFC4eU0Yru7nmGzJErI2b2"
      webhookURLs = append(webhookURLs, webhookURL) // Twitter Nova
    } else if strings.ToLower(store_url) == "amnotify" {
      webhookURL := "https://discordapp.com/api/webhooks/651645407797051392/TqG5kBHqqDStFxZOMPhFicZxbCr5PKiliRs15OoeLa5JzLy-VlHFEbyB2QK6rOtXcFXI"
      webhookURLs = append(webhookURLs, webhookURL) // Twitter AMNotify
    } else if strings.ToLower(store_url) == "restockworld" || strings.ToLower(store_url) == "impersonated" || strings.ToLower(store_url) == "ryxnSZN" {
      webhookURL := "https://discordapp.com/api/webhooks/651645777810161681/T7_P77o1NEU64SpBtKLXuYm_UB4iTtWHBjUfIwm3TWb6xHrGtLbol_NisnMnbmxgNN24"
      webhookURLs = append(webhookURLs, webhookURL) // Twitter Restock World
    }
    webhookURL := "https://discordapp.com/api/webhooks/609146658051194887/QSaA0XjLCJ7P2yx1kdrCGaELSrTvqTENnYwGKl7fMiS9ymHTsQ_0aRz4Dev0mamtQ_t7"
    webhookURLs = append(webhookURLs, webhookURL) // Twitter Unfiltered
  } else if identifier == "instagram" {
    if strings.ToLower(store_url) == "offspringhq" {
      webhookURL := "https://discordapp.com/api/webhooks/652608247316086795/BCMjd-EJKog2v8IY7EGQWWIzTUYtVvvU0cTzPXpHRJ8jj2VciNNQDCRZL7-Kyx15nY8t"
      webhookURLs = append(webhookURLs, webhookURL) // Instagram Offspring
    } else if strings.ToLower(store_url) == "cybersole"{
      webhookURL := "https://discordapp.com/api/webhooks/651643248141533206/textITlu86932zh1zjDsHG17eZ4JBPqzyjTr3mBKfTjgsojIvCMoHuph9zcQ2aZT_OJx"
      webhookURLs = append(webhookURLs, webhookURL) // Instagram Cybersole
    } else if strings.ToLower(store_url) == "balkobot" {
      webhookURL := "https://discordapp.com/api/webhooks/651643763961364501/KYur4389bxW_kD0s8Rc4gv7HN-sapLAos2bb3hpkEBH_Y4ZnEMjQavh3dKj_mMfdAsyV"
      webhookURLs = append(webhookURLs, webhookURL) // Instagram balkobot
    } else if strings.ToLower(store_url) == "adeptbots" {
      webhookURL := "https://discordapp.com/api/webhooks/651644071856701450/aFHaNiuPgWryyhVCv6sU1lFpwxkQGebSXjM8Pj-czJNi4zJzPv6emdP3gjRkhBN42ClH"
      webhookURLs = append(webhookURLs, webhookURL) // Instagram Adept
    } else if strings.ToLower(store_url) == "nova_aio" {
      webhookURL := "https://discordapp.com/api/webhooks/651644277469741057/sb8oR11rWgtqQ4xjFI6ih0ubj5kSp6IIPCv9dsrAsrq9dXjFC4eU0Yru7nmGzJErI2b2"
      webhookURLs = append(webhookURLs, webhookURL) // Instagram Nova
    } else if strings.ToLower(store_url) == "amnotifyus" {
      webhookURL := "https://discordapp.com/api/webhooks/651645407797051392/TqG5kBHqqDStFxZOMPhFicZxbCr5PKiliRs15OoeLa5JzLy-VlHFEbyB2QK6rOtXcFXI"
      webhookURLs = append(webhookURLs, webhookURL) // Instagram AMNotify
    } else if strings.ToLower(store_url) == "restockworld" {
      webhookURL := "https://discordapp.com/api/webhooks/651645777810161681/T7_P77o1NEU64SpBtKLXuYm_UB4iTtWHBjUfIwm3TWb6xHrGtLbol_NisnMnbmxgNN24"
      webhookURLs = append(webhookURLs, webhookURL) // Instagram Restock World
    }
    webhookURL := "https://discordapp.com/api/webhooks/609148205539524639/MAG95h05xMcJGkEcVQbwWnBsEe_Xx8GTCtocf5-7nJFJ81WMJDvyZI9B4vph6HA8m38Q"
    webhookURLs = append(webhookURLs, webhookURL) // Instagram Unfiltered
  }
  // return []string{"https://discordapp.com/api/webhooks/614963312555327491/Qx1-49umYNt3OlcLMI7y_Al2S5fcTJVY_NB4dwC1zSKUBCj5JqL33f9SpdqhAEhUWxT6"} // TEST WEBHOOK
  return webhookURLs
}

func grabKeywords(product_name string) string {
  product_name = strings.ToLower(product_name)
  for i := 0;  i < len(keywords); i++ {
    cur_keyword := strings.ReplaceAll(strings.ToLower(keywords[i]), ", ", ",") // remove whitespace after comma
    split_keyword := strings.Split(cur_keyword, ",")
    for j := 0;  j <= len(split_keyword); j++ {
      if j == len(split_keyword) { // reached end of current keyword search. found exact match!
        return strings.ReplaceAll(cur_keyword, ",", ", ")
      }
      keyword_prefix := split_keyword[j][0:1]
      keyword := split_keyword[j][1:]
      if keyword_prefix == "-" {
        log.Println(keyword_prefix + keyword)
        if strings.Index(product_name, keyword) > -1 { // found BAD keyword
          break
        }
      } else if keyword_prefix == "+" {
        if strings.Index(product_name, keyword) > -1 { // found GOOD keyword
          continue
        } else if strings.Index(product_name, keyword) == -1 { // found GOOD keyword
          break
        }
      }
    }
  }
	return ""
}
