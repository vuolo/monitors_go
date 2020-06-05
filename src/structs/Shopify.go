package structs

// ########################################### STRUCTS
// The JSON response from the requests
type ShopifyProduct struct {
	Products []struct {
		// == main info ==
	  Name string `json:"title"`
		Handle string `json:"handle"`
		Images []struct {
			URL string `json:"src"`
		} `json:"images"`
	  Description string `json:"body_html"`
		Variants []struct {
			Name string `json:"title"`
			ID int `json:"id"`
			Available bool `json:"available"`
			Price string `json:"price"`
		} `json:"variants"`

		// == extra ==
		PublishedAt string `json:"published_at"`
		Collection string `json:"vendor"` // (Brand)
		Tags []string `json:"tags"`
	} `json:"products"`
}
