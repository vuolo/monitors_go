package tasks

// ########################################### IMPORTS
import (
	"encoding/json"
	"log"
	"strings"
	"time"
	"context"
	"strconv"

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
func Shopify(store structs.Store, mongoClient *mongo.Client, pusherClient *pusher.Client) {
	identifier := "shopify"
	timeout := 2 * time.Second
	for {
		scrapeShopify(store, identifier, mongoClient, pusherClient)
		go checkCheckpoint(store, mongoClient, pusherClient)
		time.Sleep(timeout)
	}
}

func scrapeShopify(store structs.Store, identifier string, mongoClient *mongo.Client, pusherClient *pusher.Client) {

	// ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// ########################################### VARIABLES
	url := "https://" + store.URL + "/products.json"
	if store.Name == "" {
		store.Name = store.URL
	}

	// ########################################### START REQUEST
	request := gorequest.New()
	resp, bodyBytes, request_err := request.Proxy(proxies.GrabProxy()).Get(url).Set("User-Agent", requests.RandomUserAgent()).EndBytes()
	if request_err != nil {
		log.Println(request_err)
		return
	}

	// ########################################### HANDLE ERRORS
	if requests.ParseHTTPErrors(resp, bodyBytes, identifier, store.URL, true) {
		failed_connections++
		return
	}

	// ########################################### INITIAL CHECK
	initialChecked := false
	if (resp.StatusCode != 401 && strings.Index(resp.Request.URL.String(), "/password") == -1) && utils.StringInSlice(store.URL, initialCheckedURLs) {
		log.Println(Green("[" + "Shopify: " + store.Name + "] " + "Successful connection."))
		initialChecked = true
		go database.TryUpdatePasswordBook(store.Name, store.URL, false, mongoClient, pusherClient)
	} else if resp.StatusCode == 401 || strings.Index(resp.Request.URL.String(), "/password") > -1 { // Detect if Password Enabled
		if !utils.StringInSlice(store.URL, initialCheckedURLs) {
			initialCheckedURLs = append(initialCheckedURLs, store.URL)
			log.Println(Inverse("[" + "Shopify: " + store.Name + "] " + "Initial Check Done."))
			return
		} else {
			go database.TryUpdatePasswordBook(store.Name, store.URL, true, mongoClient, pusherClient)
			return
		}
	}
	successful_connections++

	// ########################################### HANDLE RESPONSE
	data := &structs.ShopifyProduct{}
	err := json.Unmarshal(bodyBytes, data)
	if err != nil {
		log.Println(err)
		return
	}

	// fmt.Printf("%+v\n", data.Products[0]) // print first product

	var products structs.Products
	for i := 0;  i < len(data.Products); i++ {

		// // ############## FUNKO VALIDATOR START ##############
		// // +pop, +exclusive, -dorbz, -plushies, -hikari, -vynl, -plush, -mystery, -mini, -pocket, -lanyard, -tee
		// if store.URL == "shop.funko.com" || store.URL == "bigpopshop.com" || store.URL == "bungiestore.com" || store.URL == "galactictoys.com" {
		// 	if !(strings.Contains(data.Products[i].Name, "pop") && strings.Contains(data.Products[i].Name, "exclusive") && !strings.Contains(data.Products[i].Name, "dorbz") && !strings.Contains(data.Products[i].Name, "plushies") && !strings.Contains(data.Products[i].Name, "hikari") && !strings.Contains(data.Products[i].Name, "vynl") && !strings.Contains(data.Products[i].Name, "plush") && !strings.Contains(data.Products[i].Name, "mystery") && !strings.Contains(data.Products[i].Name, "mini") && !strings.Contains(data.Products[i].Name, "pocket") && !strings.Contains(data.Products[i].Name, "lanyard") && !strings.Contains(data.Products[i].Name, "tee")) {
		// 		continue
		// 	}
		// }
		// // ############## FUNKO VALIDATOR END ##############

		var product structs.Product
		// == important ==
		product.Store = store.URL
		product.StoreName = store.Name

		// == main info ==
		product.Name = data.Products[i].Name
		product.URL = "https://" + store.URL + "/products/" + data.Products[i].Handle
		product.Price = store.Currency + data.Products[i].Variants[0].Price
		if data.Products[i].Images == nil || len(data.Products[i].Images) < 1  {
			product.ImageURL = "https://i.imgur.com/fip3nw5.png";
		} else {
			product.ImageURL = data.Products[i].Images[0].URL
		}
		product.Description = data.Products[i].Description
		productAvailable := false
		for j := 0;  j < len(data.Products[i].Variants); j++ {
			variant := structs.Variant{
				strings.Replace(data.Products[i].Variants[j].Name, "Default Title", "One Size", 1),
				strconv.Itoa(data.Products[i].Variants[j].ID),
				data.Products[i].Variants[j].Available,
				data.Products[i].Variants[j].Price,
				-420,
			}
			product.Variants = append(product.Variants, variant)
			if !productAvailable && data.Products[i].Variants[j].Available {
				productAvailable = true
				// break // break is only necessary if variants are NOT setup in this nested for loop
			}
		}
		product.Available = productAvailable
		product.Identifier = identifier

		// == extra ==
		product.LaunchDate = data.Products[i].PublishedAt
		product.Collection = data.Products[i].Collection // (Brand)
		product.Tags = data.Products[i].Tags
		convertedVariants, _ := json.Marshal(product.Variants)
		product.MD5 = utils.GetMD5Hash(string(convertedVariants))
		products = append(products, product)
	}

	database.SendToDatabase(products, identifier, store.URL, store.Name, initialChecked, mongoClient, pusherClient)

	if !initialChecked {
		initialCheckedURLs = append(initialCheckedURLs, store.URL)
		log.Println(Inverse("[" + "Shopify: " + store.Name + "] " + "Initial Check Done."))
	}

}

func checkCheckpoint(store structs.Store, mongoClient *mongo.Client, pusherClient *pusher.Client) {
	// ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// ########################################### VARIABLES
	url := "https://" + store.URL + "/checkpoint"
	if store.Name == "" {
		store.Name = store.URL
	}

	// ########################################### START REQUEST
	request := gorequest.New()
	resp, _, request_err := request.Proxy(proxies.GrabProxy()).Get(url).Set("User-Agent", requests.RandomUserAgent()).EndBytes()
	if request_err != nil {
		log.Println(request_err)
		return
	}

	// ########################################### HANDLE ERRORS
	// Detect if Checkpoint Enabled
	if resp.StatusCode == 200 && strings.Index(resp.Request.URL.String(), "/checkpoint") > -1 {
		go database.TryUpdateCheckpointBook(store.Name, store.URL, true, mongoClient, pusherClient)
	} else if resp.StatusCode == 404 || resp.StatusCode == 401 {
		go database.TryUpdateCheckpointBook(store.Name, store.URL, false, mongoClient, pusherClient)
	}
	successful_connections++
	return
}
