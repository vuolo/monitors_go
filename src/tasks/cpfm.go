package tasks

// ########################################### IMPORTS
import (
	"log"
	"time"
	"context"
	// "net"
	"encoding/json"

	"github.com/chouandy/shopify-graphql"

	"go.mongodb.org/mongo-driver/mongo" // MongoDB

	"github.com/pusher/pusher-http-go" // Pusher

	. "github.com/logrusorgru/aurora" // colors

	// ## local shits ##
	"../database"
	"../structs"
	"../utils"
)

// ########################################### VARIABLES
var shopifyStore = "cactusplantfleamarket.myshopify.com"
var shopifyStorefrontAccessToken = "574d05e81d915e3c13a16c514c678649"
var shopifyClient = shopifygraphql.NewStorefrontClient(shopifyStore, shopifyStorefrontAccessToken)

func CPFM(store structs.Store, mongoClient *mongo.Client, pusherClient *pusher.Client) {
	identifier := "cpfm"
	// timeout := 1 * time.Second
	// timeout := 500 * time.Millisecond
	for {
		scrapeCPFM(store, identifier, mongoClient, pusherClient)
		// time.Sleep(timeout)
	}
}

func scrapeCPFM(store structs.Store, identifier string, mongoClient *mongo.Client, pusherClient *pusher.Client) {

	// ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// ########################################### START REQUEST
	var productQuery structs.ProductQuery
  queryErr := shopifyClient.Query(context.Background(), &productQuery, nil)
  if queryErr != nil {
    // log.Println(Red("ERROR:"), queryErr)
		log.Println(Red("[" + "CACTUS PLANT FLEA MARKET" + "] " + "Connection failed (Temp ban occured)"))
		failed_connections++
		time.Sleep(5 * time.Minute)
    return
  }

	// ########################################### INITIAL CHECK
	initialChecked := false
	if utils.StringInSlice(store.URL, initialCheckedURLs) {
		log.Println(Green("[" + store.Name + "] " + "Successful connection."))
		initialChecked = true
	} else {
		initialCheckedURLs = append(initialCheckedURLs, store.URL)
		log.Println(Inverse("[" + store.Name + "] " + "Initial Check Done."))
	}
	successful_connections++

	// ########################################### HANDLE RESPONSE
	var products structs.Products
	for i := 0;  i < len(productQuery.Products.Edges); i++ {
		var product structs.Product
		// == important ==
		product.Store = store.URL
		product.StoreName = store.Name

		// == main info ==
		product.Name = string(productQuery.Products.Edges[i].Node.Title)
		product.URL = "https://" + shopifyStore + "/products/" + string(productQuery.Products.Edges[i].Node.Handle)
		product.Price = store.Currency + string(productQuery.Products.Edges[i].Node.Variants.Edges[0].Node.PriceV2.Amount)
		if len(productQuery.Products.Edges[i].Node.Images.Edges) == 0 {
			product.ImageURL = "https://i.imgur.com/fip3nw5.png";
		} else {
			product.ImageURL = string(productQuery.Products.Edges[i].Node.Images.Edges[0].Node.OriginalSrc)
		}
		product.Description = string(productQuery.Products.Edges[i].Node.Description)
		productAvailable := false
		for j := 0;  j < len(productQuery.Products.Edges[i].Node.Variants.Edges); j++ {
			variant := structs.Variant{
				string(productQuery.Products.Edges[i].Node.Options[0].Values[j]),
				string(productQuery.Products.Edges[i].Node.Variants.Edges[j].Node.ID),
				bool(productQuery.Products.Edges[i].Node.Variants.Edges[j].Node.AvailableForSale),
				string(productQuery.Products.Edges[i].Node.Variants.Edges[j].Node.PriceV2.Amount),
				-420,
			}
			product.Variants = append(product.Variants, variant)
			if !productAvailable && bool(productQuery.Products.Edges[i].Node.Variants.Edges[j].Node.AvailableForSale) {
				productAvailable = true
				// break // break is only necessary if variants are NOT setup in this nested for loop
			}
		}
		product.Available = productAvailable
		product.Identifier = identifier

		// == extra ==
		product.LaunchDate = string(productQuery.Products.Edges[i].Node.PublishedAt)
		product.OverrideURL = "https://" + store.URL
		convertedVariants, _ := json.Marshal(product.Variants)
		product.MD5 = utils.GetMD5Hash(string(convertedVariants))
		products = append(products, product)
	}

	database.SendToDatabase(products, identifier, store.URL, store.Name, initialChecked, mongoClient, pusherClient)

}
