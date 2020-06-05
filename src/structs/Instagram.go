package structs

// ########################################### STRUCTS
type InstagramUser struct {
  Username string `bson:"username"`
  ID string `bson:"id"`
  ImageURL string `bson:"imageurl"`
  FullName string `bson:"fullname"`
  Biography string `bson:"biography"`
  ExternalURL string `bson:"externalurl"`
  Posts []InstagramPost `bson:"posts"`
  Stories []InstagramStory `bson:"stories"`
}

type InstagramPost struct {
  ID string `bson:"id"`
  ImageURL string `bson:"imageurl"`
  IsVideo bool `bson:"isvideo"`
  Caption string `bson:"caption"`
  Timestamp int `bson:"timestamp"`
}

type InstagramStory struct {
  ID string `bson:"id"`
  ImageURL string `bson:"imageurl"`
  IsVideo bool `bson:"isvideo"`
  ExternalURL string `bson:"externalurl"`
}

type InstagramUserJSON struct {
  GraphQL struct{
    User struct {
      Username string `json:"username"`
      ID string `json:"id"`
      ImageURL string `json:"profile_pic_url_hd"`
      FullName string `json:"full_name"`
      Biography string `json:"biography"`
      ExternalURL string `json:"external_url"`
      Posts struct{
        Edges []struct {
          Node struct {
            ID string `json:"shortcode"` // `json:"id"`
            ImageURL string `json:"display_url"`
            IsVideo bool `json:"is_video"`
            Timestamp int `json:"taken_at_timestamp"`
            Caption struct{
              Edges []struct {
                Node struct {
                  Text string `json:"text"`
                } `json:"node,omitempty"`
              } `json:"edges"`
            } `json:"edge_media_to_caption"`
          } `json:"node"`
        } `json:"edges"`
      } `json:"edge_owner_to_timeline_media"`
    } `json:"user"`
  } `json:"graphql"`
}

type InstagramStoriesJSON struct {
  Data struct{
    ReelsMedia []struct {
      Items []struct {
        ID string `json:"id"`
        ImageURL string `json:"display_url"`
        IsVideo bool `json:"is_video"`
        ExternalURL string `json:"story_cta_url"`
      } `json:"items"`
    } `json:"reels_media"`
  } `json:"data"`
}
