package structs

// ########################################### STRUCTS
type StockXProduct struct {
	// == main info ==
  Name string
	Handle string
	ImageURL string
	Variants []StockXVariant

	// == extra ==
	Collection string // (Brand)
	DeadstockSold int
}

type StockXVariant struct {
	Name string
	RetailPrice float64
	LowestAsk float64
	HighestBid float64
	LastSale float64
}

// The JSON response from the requests
type StockXProductJSON struct {
	Product struct {
		// == main info ==
	  Name string `json:"title"`
		Handle string `json:"urlKey"`
		Image struct {
			URL string `json:"imageUrl"`
		} `json:"media"`
		Variants struct {
      Variant map[string]interface{} `json:"-"`
      // Variant struct {
      //   Name string `json:"shoeSize"`
      //   RetailPrice int `json:"retailPrice"`
      //   Market struct {
      //     LowestAsk int `json:"lowestAsk"`
      //     HighestBid int `json:"highestBid"`
      //     LastSale int `json:"lastSale"`
      //   } `json:"market"`
      // } `json:"-"`
		} `json:"children"`

		// == extra ==
		Collection string `json:"brand"` // (Brand)
    Market struct {
      DeadstockSold int `json:"deadstockSold"`
    } `json:"market"`
	} `json:"product"`
}

type StockXProductBSON struct {
  // == main info ==
  Name string `bson:"name"`
	Handle string `bson:"handle"`
	ImageURL string `bson:"imageurl"`
	Variants []StockXVariantBSON `bson:"variants"`

	// == extra ==
	Collection string `bson:"collection"` // (Brand)
	DeadstockSold int `bson:"deadstocksold"`
}

type StockXVariantBSON struct {
  Name string `bson:"name"`
  RetailPrice float64 `bson:"retailprice"`
  LowestAsk float64 `bson:"lowestask"`
  HighestBid float64 `bson:"highestbid"`
  LastSale float64 `bson:"lastsale"`
}
