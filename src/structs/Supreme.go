package structs

// ########################################### STRUCTS
// The JSON response from the requests
type SupremeProduct struct {
	Categories struct {
		Bags []SupremeCategory `json:"Bags"`
		Pants []SupremeCategory `json:"Pants"`
		Accessories []SupremeCategory `json:"Accessories"`
		Skate []SupremeCategory `json:"Skate"`
		Shoes []SupremeCategory `json:"Shoes"`
		Hats []SupremeCategory `json:"Hats"`
		Shirts []SupremeCategory `json:"Shirts"`
		Sweatshirts []SupremeCategory `json:"Sweatshirts"`
		Tops_Sweaters []SupremeCategory `json:"Tops/Sweaters"`
		Jackets []SupremeCategory `json:"Jackets"`
		T_Shirts []SupremeCategory `json:"T-Shirts"`
		New []SupremeCategory `json:"new"`
	} `json:"products_and_categories"`
}

type SupremeCategory struct {
	Name string `json:"name"`
	ID int `json:"id"`
	ImageURL string `json:"image_url_hi"`
	Price float64 `json:"price"`
}

type SupremeProductDetailed struct {
	Styles []struct {
		ID int `json:"id"`
		Name string `json:"name"`
		Currency string `json:"currency"`
		ImageURL string `json:"image_url_hi"`
		Variants []struct {
			Name string `json:"name"`
			ID int `json:"id"`
			Quantity int `json:"stock_level"`
		} `json:"sizes"`
	} `json:"styles"`
	Description string `json:"description"`
}
