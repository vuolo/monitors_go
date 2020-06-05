package structs

// ########################################### STRUCTS
// The JSON response from the requests
type SNKRSProduct struct {
	Products []struct {
		// == main info ==
	  Name string `json:"seoTitle"`
		ID string `json:"id"`
		Handle string `json:"seoSlug"`
		ImageURL string `json:"imageUrl"`
		Cards []struct {
			Description string `json:"description"`
		} `json:"cards"`
		SEODescription string `json:"seoDescription"`
		Product struct {
			Available bool `json:"available"`
			SKUS []struct {
				Name string `json:"localizedSize"`
				Available bool `json:"available"`
			}
			Price struct {
				Retail float64 `json:"currentRetailPrice"`
			} `json:"price"`
			// == extra ==
			LaunchDate string `json:"startSellDate"`
		} `json:"product"`

		// == extra ==
		Restricted bool `json:"restricted"`
		Method string `json:"selectionEngine"`
		Tags []string `json:"tags"`
	} `json:"threads"`
}
