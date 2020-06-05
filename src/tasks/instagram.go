package tasks

// ########################################### IMPORTS
import(
  "encoding/json"
  "strings"
	"time"
  "log"
  "context"

  "github.com/parnurzeal/gorequest"

	"go.mongodb.org/mongo-driver/mongo" // MongoDB

  "github.com/pusher/pusher-http-go" // Pusher

  . "github.com/logrusorgru/aurora" // colors

  // ## local shits ##
	"../database"
  "../proxies"
  "../requests"
  "../structs"
	"../utils"
)

// ########################################### VARIABLES
func Instagram(instagramHandle string, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	identifier := "instagram"
	timeout := 864 * time.Millisecond
	for {
    var foundHandle = false
		for i := 0;  i < len(database.DatabaseInstagramHandles); i++ {
			if instagramHandle == database.DatabaseInstagramHandles[i].Handle {
				foundHandle = true
				break
			}
	  }
		if foundHandle {
      scrapeInstagram(instagramHandle, identifier, mongoClient, pusherClient)
  		time.Sleep(timeout)
		} else {
			break
		}
	}
}

func scrapeInstagram(instagramHandle string, identifier string, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  // ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

  // ########################################### VARIABLES
	url := "https://www.instagram.com/" + instagramHandle + "/?__a=1"
  curProxy := proxies.GrabProxy()

	// ########################################### START FIRST REQUEST
	request := gorequest.New()
	// resp, bodyBytes, request_err := request.Proxy(curProxy).Get(url).Set("Referer", "https://www.instagram.com/" + instagramHandle + "/").Set("User-Agent", requests.RandomPhoneUserAgent()).EndBytes()
  // resp, bodyBytes, request_err := request.Proxy(curProxy).Get(url).Set("Accept", "*/*").Set("Accept-Encoding", "deflate, br").Set("Accept-Language", "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7").Set("Cookie", `mid=Xb4YogALAAEX6FaDI4IQRJ8GTrfe; fbm_124024574287414=base_domain=.instagram.com; csrftoken=MuE3947nqnbXh3csoUx3rphiwlhjQmHV; ds_user_id=4189588419; sessionid=4189588419%3AdHmLDynbUKCYkA%3A22; ig_did=F27DF5FA-E378-45EA-86F3-73E92D299F38; shbid=7084; shbts=1580084551.3943954; rur=ATN; urlgen="{\"2601:882:100:7c13:c550:9e34:d5ae:b10c\": 7922\054 \"2601:882:100:7c13:bc5e:e6d1:be77:5d3d\": 7922}:1iwBmh:mVhSKKM8PmA2DjCybY7igVKvbTg"`).Set("Dnt", "1").Set("Referer", "https://www.instagram.com/" + instagramHandle + "/").Set("Sec-Fetch-Mode", "cors").Set("Sec-Fetch-Site", "same-origin").Set("User-Agent", requests.RandomPhoneUserAgent()).Set("X-Csrftoken", "MuE3947nqnbXh3csoUx3rphiwlhjQmHV").Set("X-Ig-App-Id", "936619743392459").Set("X-Ig-Www-Claim", "hmac.AR3hleTo9xJHFmRF6CKzFIm38NV1hVjGUUiJcxQBm81cMOT8").EndBytes()
	resp, bodyBytes, request_err := request.Proxy(curProxy).Get(url).Set("Accept", "*/*").Set("Accept-Encoding", "deflate, br").Set("Accept-Language", "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7").Set("Cookie", `mid=Xb4YogALAAEX6FaDI4IQRJ7GTrfe; fbm_124024574287414=base_domain=.instagram.com; csrftoken=DuE3947nqnbXh3csoUx3rphiwlhjQmHV; ds_user_id=3189588419; sessionid=3189588419%3AdHmLDynbUKCYkA%3A22; ig_did=F17DF5FA-E378-45EA-86F3-73E92D299F38; shbid=6084; shbts=2580084551.3943954; rur=ATN; urlgen="{\"3601:882:100:7c13:c550:9e34:d5ae:b10c\": 3922\054 \"2201:882:100:7c13:bc5e:e6d1:be77:5d3d\": 7923}:2iwBmh:mVhSKKM8PmA2DjCybY7igVKvbTg"`).Set("Dnt", "1").Set("Referer", "https://www.instagram.com/" + instagramHandle + "/").Set("Sec-Fetch-Mode", "cors").Set("Sec-Fetch-Site", "same-origin").Set("User-Agent", requests.RandomPhoneUserAgent()).Set("X-Csrftoken", "MuE3947nqnbXh3csoUx3rphiwlhjQmHV").Set("X-Ig-App-Id", "936619743392459").Set("X-Ig-Www-Claim", "hmac.AR3hleTo9xJHFmRF6CKzFIm38NV1hVjGUUiJcxQBm81cMOT8").EndBytes()
  // resp, bodyBytes, request_err := request.Get(url).Set("Referer", "https://www.instagram.com/" + instagramHandle + "/").Set("User-Agent", requests.RandomUserAgent()).EndBytes()
	if request_err != nil {
		log.Println(request_err)
		return
	}

	// ########################################### HANDLE FIRST ERRORS
	if requests.ParseHTTPErrors(resp, bodyBytes, identifier, instagramHandle, true) {
		// failed_connections++
		return
	}
  // ########################################### HANDLE NOT LOGGED IN ERROR
  if strings.Index(string(bodyBytes), "not-logged-in") > -1 {
    if true {
      log.Println(Red("[" + "Instagram: " + instagramHandle + "] " + "Not logged in."))
      return
    }
    log.Println(Red("[" + "Instagram: " + instagramHandle + "] " + "Not logged in. Attempting to login."))

    BASE_URL := "https://www.instagram.com/accounts/login/"
    LOGIN_URL := BASE_URL + "ajax/"
    USERNAME := "dashpings_monitor"
    PASSWORD := "Kirby320"

    initialLoginRequest := gorequest.New()
    _, initialLoginBodyBytes, initialLoginRequest_err := initialLoginRequest.Proxy(curProxy).Get(BASE_URL).Set("User-Agent", requests.RandomPhoneUserAgent()).EndBytes()
    if initialLoginRequest_err != nil {
      log.Println(Red("[" + "Instagram: " + instagramHandle + "] " + "Error logging in."))
      return
    }

    csrf_token_index := strings.Index(string(initialLoginBodyBytes), `"csrf_token":"`)
    csrf := string(initialLoginBodyBytes)[csrf_token_index:csrf_token_index+100] // strings get from starting index of csrf_token to end of "
    csrf = strings.ReplaceAll(csrf, `"csrf_token":"`, "")
    csrf = csrf[:strings.Index(csrf, `"`)]

    requestBody, err := json.Marshal(map[string]string{
      "username": USERNAME,
      "password": PASSWORD,
    })
    if err != nil {
      log.Println(err)
      return
    }
    loginRequest := gorequest.New()
    _, loginBodyBytes, loginRequest_err := loginRequest.Proxy(curProxy).Post(LOGIN_URL).Send(string(requestBody)).Set("Referer", "https://www.instagram.com/accounts/login/").Set("User-Agent", requests.RandomPhoneUserAgent()).Set("X-CSRFToken", csrf).EndBytes()
    if loginRequest_err != nil {
      log.Println(Red("[" + "Instagram: " + instagramHandle + "] " + "Error logging in."))
      return
    }
    // log.Println(string(loginBodyBytes))
    if strings.Index(string(loginBodyBytes), `"authenticated": true`) > -1 {
			log.Println(Yellow("[" + "Instagram: " + instagramHandle + "] " + "Successfully logged in."))
    } else {
			log.Println(Red("[" + "Instagram: " + instagramHandle + "] " + "Login authentication failed."))
      return
		}


  }

	// ########################################### INITIAL CHECK
	initialChecked := false
	if utils.StringInSlice("instagram_" + instagramHandle, initialCheckedURLs) {
		log.Println(Green("[" + "Instagram: " + instagramHandle + "] " + "Successful connection."))
		initialChecked = true
	} else {
		initialCheckedURLs = append(initialCheckedURLs, "instagram_" + instagramHandle)
		log.Println(Inverse("[" + "Instagram: " + instagramHandle + "] " + "Initial Check Done."))
	}
	successful_connections++

	// ########################################### HANDLE FIRST RESPONSE
	data := &structs.InstagramUserJSON{}
	err := json.Unmarshal(bodyBytes, data)
	if err != nil {
		log.Println(err)
		return
	}

  var instagramUser structs.InstagramUser
  instagramUser.Username = data.GraphQL.User.Username
  instagramUser.ID = data.GraphQL.User.ID
  instagramUser.ImageURL = data.GraphQL.User.ImageURL
  instagramUser.FullName = data.GraphQL.User.FullName
  instagramUser.Biography = data.GraphQL.User.Biography
	instagramUser.ExternalURL = data.GraphQL.User.ExternalURL
  for i := 0;  i < len(data.GraphQL.User.Posts.Edges); i++ {
    caption := ""
    if len(data.GraphQL.User.Posts.Edges[i].Node.Caption.Edges) > 0 {
      caption = data.GraphQL.User.Posts.Edges[i].Node.Caption.Edges[0].Node.Text
    }
    post := structs.InstagramPost{
      data.GraphQL.User.Posts.Edges[i].Node.ID,
      data.GraphQL.User.Posts.Edges[i].Node.ImageURL,
      data.GraphQL.User.Posts.Edges[i].Node.IsVideo,
      caption,
      data.GraphQL.User.Posts.Edges[i].Node.Timestamp,
    }
    instagramUser.Posts = append(instagramUser.Posts, post)
  }

  // if !initialChecked {
  //   utils.ReverseSlice(instagramUser.Posts)
  // }

	database.SendToInstagramDatabase(instagramUser, identifier, initialChecked, mongoClient, pusherClient)

  // ########################################### START SECOND REQUEST
  storiesUrl := "https://www.instagram.com/graphql/query/?query_hash=52a36e788a02a3c612742ed5146f1676&variables=%7B%22reel_ids%22%3A%5B%22" + data.GraphQL.User.ID + "%22%5D%2C%22tag_names%22%3A%5B%5D%2C%22location_ids%22%3A%5B%5D%2C%22highlight_reel_ids%22%3A%5B%5D%2C%22precomposed_overlay%22%3Afalse%2C%22show_story_viewer_list%22%3Atrue%2C%22story_viewer_fetch_count%22%3A50%2C%22story_viewer_cursor%22%3A%22%22%2C%22stories_video_dash_manifest%22%3Afalse%7D"
  storiesRequest := gorequest.New()
	storiesResp, storiesBodyBytes, storiesRequest_err := storiesRequest.Proxy(curProxy).Get(storiesUrl).Set("Accept", "*/*").Set("Accept-Encoding", "deflate, br").Set("Accept-Language", "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7").Set("Cookie", `mid=Xb4YogALAAEX6FaDI4IQRJ7GTrfe; fbm_124024574287414=base_domain=.instagram.com; csrftoken=DuE3947nqnbXh3csoUx3rphiwlhjQmHV; ds_user_id=3189588419; sessionid=3189588419%3AdHmLDynbUKCYkA%3A22; ig_did=F17DF5FA-E378-45EA-86F3-73E92D299F38; shbid=6084; shbts=2580084551.3943954; rur=ATN; urlgen="{\"3601:882:100:7c13:c550:9e34:d5ae:b10c\": 3922\054 \"2201:882:100:7c13:bc5e:e6d1:be77:5d3d\": 7923}:2iwBmh:mVhSKKM8PmA2DjCybY7igVKvbTg"`).Set("Dnt", "1").Set("Referer", "https://www.instagram.com/stories/" + instagramHandle + "/").Set("Sec-Fetch-Mode", "cors").Set("Sec-Fetch-Site", "same-origin").Set("User-Agent", requests.RandomPhoneUserAgent()).Set("X-Csrftoken", "MuE3947nqnbXh3csoUx3rphiwlhjQmHV").Set("X-Ig-App-Id", "936619743392459").Set("X-Ig-Www-Claim", "hmac.AR3hleTo9xJHFmRF6CKzFIm38NV1hVjGUUiJcxQBm81cMOT8").EndBytes()
	// storiesResp, storiesBodyBytes, storiesRequest_err := storiesRequest.Proxy(curProxy).Get(storiesUrl).Set("Accept", "*/*").Set("Accept-Encoding", "deflate, br").Set("Accept-Language", "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7").Set("Cookie", `mid=Xb4YogALAAEX6FaDI4IQRJ8GTrfe; fbm_124024574287414=base_domain=.instagram.com; csrftoken=MuE3947nqnbXh3csoUx3rphiwlhjQmHV; ds_user_id=4189588419; sessionid=4189588419%3AdHmLDynbUKCYkA%3A22; ig_did=F27DF5FA-E378-45EA-86F3-73E92D299F38; shbid=7084; shbts=1580084551.3943954; rur=ATN; urlgen="{\"2601:882:100:7c13:c550:9e34:d5ae:b10c\": 7922\054 \"2601:882:100:7c13:bc5e:e6d1:be77:5d3d\": 7922}:1iwDZM:psuIVVYSEDjWHXg1viICuPo-Jmk"`).Set("Dnt", "1").Set("Referer", "https://www.instagram.com/stories/" + instagramHandle + "/").Set("Sec-Fetch-Mode", "cors").Set("Sec-Fetch-Site", "same-origin").Set("User-Agent", requests.RandomPhoneUserAgent()).Set("X-Csrftoken", "MuE3947nqnbXh3csoUx3rphiwlhjQmHV").Set("X-Ig-App-Id", "936619743392459").Set("X-Ig-Www-Claim", "hmac.AR3hleTo9xJHFmRF6CKzFIm38NV1hVjGUUiJcxQBm81cMOT8").EndBytes()
  // storiesResp, storiesBodyBytes, storiesRequest_err := request.Get(storiesUrl).Set("Accept", "*/*").Set("Accept-Encoding", "deflate, br").Set("Accept-Language", "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7").Set("Cookie", `mid=Xb4YogALAAEX6FaDI4IQRJ8GTrfe; fbm_124024574287414=base_domain=.instagram.com; csrftoken=MuE3947nqnbXh3csoUx3rphiwlhjQmHV; ds_user_id=4189588419; sessionid=4189588419%3AdHmLDynbUKCYkA%3A22; ig_did=F27DF5FA-E378-45EA-86F3-73E92D299F38; shbid=7084; shbts=1580084551.3943954; rur=ATN; urlgen="{\"2601:882:100:7c13:c550:9e34:d5ae:b10c\": 7922\054 \"2601:882:100:7c13:bc5e:e6d1:be77:5d3d\": 7922}:1iwBmh:mVhSKKM8PmA2DjCybY7igVKvbTg"`).Set("Dnt", "1").Set("Referer", "https://www.instagram.com/stories/" + instagramHandle + "/").Set("Sec-Fetch-Mode", "cors").Set("Sec-Fetch-Site", "same-origin").Set("User-Agent", requests.RandomPhoneUserAgent()).Set("X-Csrftoken", "MuE3947nqnbXh3csoUx3rphiwlhjQmHV").Set("X-Ig-App-Id", "936619743392459").Set("X-Ig-Www-Claim", "hmac.AR3hleTo9xJHFmRF6CKzFIm38NV1hVjGUUiJcxQBm81cMOT8").EndBytes()
	if storiesRequest_err != nil {
		log.Println(storiesRequest_err)
		return
	}

  // ########################################### HANDLE SECOND ERRORS
	if requests.ParseHTTPErrors(storiesResp, storiesBodyBytes, identifier, instagramHandle, true) {
		// failed_connections++
		return
	}
  if strings.Index(string(storiesBodyBytes), "not-logged-in") > -1 {
    log.Println(Red("[" + "Instagram: " + instagramHandle + "] " + "Not logged in."))
    return
  }

	// ########################################### INITIAL CHECK
	storiesinitialChecked := false
	if utils.StringInSlice("instagram_" + instagramHandle + "_stories", initialCheckedURLs) {
		log.Println(Green("[" + "Instagram Stories: " + instagramHandle + "] " + "Successful connection."))
		storiesinitialChecked = true
	} else {
		initialCheckedURLs = append(initialCheckedURLs, "instagram_" + instagramHandle + "_stories")
		log.Println(Inverse("[" + "Instagram Stories: " + instagramHandle + "] " + "Initial Check Done."))
	}
	successful_connections++

	// ########################################### HANDLE SECOND RESPONSE
  storiesData := &structs.InstagramStoriesJSON{}
	storiesErr := json.Unmarshal(storiesBodyBytes, storiesData)
	if storiesErr != nil {
		log.Println(storiesErr)
		return
	}

	if len(storiesData.Data.ReelsMedia) > 0 {
		for i := 0;  i < len(storiesData.Data.ReelsMedia[0].Items); i++ {
		  story := structs.InstagramStory{
			    storiesData.Data.ReelsMedia[0].Items[i].ID,
			    storiesData.Data.ReelsMedia[0].Items[i].ImageURL,
			    storiesData.Data.ReelsMedia[0].Items[i].IsVideo,
			    storiesData.Data.ReelsMedia[0].Items[i].ExternalURL,
			  }
			  instagramUser.Stories = append(instagramUser.Stories, story)
		}
		database.SendToInstagramStoryDatabase(instagramUser, identifier, storiesinitialChecked, mongoClient, pusherClient)
	}

  return
}
