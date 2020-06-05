package tasks

// ########################################### IMPORTS
import (
	"encoding/json"
	"log"
	"time"
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
func StockX(productHandle string, mongoClient *mongo.Client, pusherClient *pusher.Client) {
	identifier := "stockx"
	timeout := 25 * time.Second
	for {
		scrapeStockX(productHandle, identifier, mongoClient, pusherClient)
		time.Sleep(timeout)
	}
}

func scrapeStockX(productHandle string, identifier string, mongoClient *mongo.Client, pusherClient *pusher.Client) {

	// ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// ########################################### VARIABLES
	url := "https://stockx.com/api/products/" + productHandle + "?includes=market&currency=USD"

	// ########################################### START REQUEST
	request := gorequest.New()
	resp, bodyBytes, request_err := request.Proxy(proxies.GrabProxy()).Get(url).Set("User-Agent", requests.RandomUserAgent()).EndBytes()
	if request_err != nil {
		log.Println(request_err)
		return
	}

	// ########################################### HANDLE ERRORS
	if requests.ParseHTTPErrors(resp, bodyBytes, identifier, productHandle, true) {
		failed_connections++
		return
	}

	// ########################################### INITIAL CHECK
	initialChecked := false
	if utils.StringInSlice(url, initialCheckedURLs) {
		log.Println(Green("[" + "StockX: " + productHandle + "] " + "Successful connection."))
		initialChecked = true
	} else {
		initialCheckedURLs = append(initialCheckedURLs, url)
		log.Println(Inverse("[" + "StockX: " + productHandle + "] " + "Initial Check Done."))
	}
	successful_connections++

	// ########################################### HANDLE RESPONSE
	data := &structs.StockXProductJSON{}
	err := json.Unmarshal(bodyBytes, data)
	if err != nil {
		log.Println(err)
		return
	}

	if err := json.Unmarshal(bodyBytes, &data.Product.Variants.Variant); err != nil {
		log.Println(err)
		return
  }

	var product structs.StockXProduct
	product.Name = data.Product.Name
	product.Handle = data.Product.Handle
	product.ImageURL = data.Product.Image.URL

	for _, record := range data.Product.Variants.Variant {
    if rec, ok := record.(map[string]interface{}); ok {
      for key, val := range rec {
				if key == "children" {

					variant, _ := val.(map[string]interface{})
					for _, variant_value := range variant {

						child, _ := variant_value.(map[string]interface{})
						shoeSize := child["shoeSize"]
						retailPrice := child["retailPrice"] // DOES NOT WORK ON NON-SHOES (non-shoes require to look at traits)
						if retailPrice == nil {
							for _, traits_value := range child["traits"].([]interface{}) {
								traits_value2, _ := traits_value.(map[string]interface{})
								if traits_value2["format"] != nil && traits_value2["format"] == "currency" {
									retailPrice = traits_value2["value"] // NON-SHOE RETAIL PRICE
								}
							}
						}

						child_market, _ := child["market"].(map[string]interface{})
						lowestAsk := child_market["lowestAsk"]
						highestBid := child_market["highestBid"]
						lastSale := child_market["lastSale"]

						if retailPrice == nil {
							retailPrice = 0.00
						}

						variant := structs.StockXVariant{
							shoeSize.(string),
							retailPrice.(float64),
							lowestAsk.(float64),
							highestBid.(float64),
							lastSale.(float64),
						}
						product.Variants = append(product.Variants, variant)

					}

				}
      }
    }
	}

	product.Collection = data.Product.Collection
	product.DeadstockSold = data.Product.Market.DeadstockSold

	database.SendToStockXDatabase(product, identifier, initialChecked, mongoClient, pusherClient)

}
