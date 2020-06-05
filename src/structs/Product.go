package structs

// ########################################### STRUCTS
type Store struct {
  URL string
	Name string
	Currency string
}

// Our custom Product struct to use for posting to webhooks/saving to DB
type Product struct {
	// == important ==
	Store string
	StoreName string

	// == main info ==
  Name string
	URL string
	Price string
	ImageURL string
  Description string
	Available bool
	Variants []Variant
	Identifier string

	// == extra ==
	LaunchDate string
	Color string
	Collection string // (Brand)
	Keywords string
	Tags []string
	OverrideURL string
  MD5 string
}

// Custom variant struct
type Variant struct {
	Name string
	ID string
	Available bool
	Price string
	Quantity int
}

type Products []Product
