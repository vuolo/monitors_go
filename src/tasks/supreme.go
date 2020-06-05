package tasks

// ########################################### IMPORTS
import (
	"encoding/json"
	"log"
	"strings"
	"time"
	"context"
	"strconv"
	"fmt"
	"math/rand"

	"github.com/parnurzeal/gorequest"

	"github.com/chromedp/chromedp/device"
	"github.com/chromedp/cdproto/network"
  "github.com/chromedp/chromedp"

	"go.mongodb.org/mongo-driver/mongo" // MongoDB

	"github.com/pusher/pusher-http-go" // Pusher

	. "github.com/logrusorgru/aurora" // colors

	// ## local shits ##
	"../database"
	// "../proxies"
	"../requests"
	"../structs"
	"../utils"
)

// ########################################### VARIABLES
func Supreme(store structs.Store, mongoClient *mongo.Client, pusherClient *pusher.Client) {
	var restockModeEnabled = true
	var useChrome = true
	identifier := "supreme_" + strings.ToLower(strings.Split(store.Name, " ")[1])
	timeout := 1 * time.Second
	// if restockModeEnabled {
	// 	timeout = 5 * time.Second
	// }
	for {
		scrapeSupreme(store, identifier, mongoClient, pusherClient, restockModeEnabled, useChrome)
		time.Sleep(timeout)
	}
}

func scrapeSupreme(store structs.Store, identifier string, mongoClient *mongo.Client, pusherClient *pusher.Client, restockModeEnabled bool, useChrome bool) {

	// ########################################### GARBAGE COLLECTOR
	_, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// ########################################### VARIABLES
	url := "https://www." + store.URL + "/mobile_stock.json"
	if store.Name == "" {
		store.Name = store.URL
	}

	var bodyBytes []byte

	// ########################################### START REQUEST
	if useChrome {
	  headers := map[string]interface{}{
	    "user-agent": requests.RandomPhoneUserAgent(),
	    "connection": "keep-alive",
	    "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	    "accept-language": "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7",
	    "cache": "max-age=0",
	    "cookie": "__utmz=74692624.1588174094.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utma=74692624.1052763426.1588174094.1588648181.1588839568.7; __utmc=74692624; __utmt=1; shoppingSessionId=1588839580822; lastVisitedFragment=products/173300; __utmb=74692624." + "18" + ".10.1588839568; _ticket=" + utils.GetMD5Hash(strconv.Itoa(rand.Int())) + "b04664143fa6476c428fa226190cac5e579b1fe58ca40baf9fc147f02735043845d3116c503a897cbcc1683840f31ba" + "41588839595",
	    "dnt": "1",
	    "referer": "https://www.supremenewyork.com/mobile/",
	    "sec-fetch-dest": "empty",
	    "sec-fetch-mode": "cors",
	    "sec-fetch-site": "same-origin",
	    "sec-fetch-user": "?1",
	    "x-requested-with": "XMLHttpRequest",
	    "upgrade-insecure-requests": "1",
	  }

		scrapeSupremeWithChrome(url, headers, &bodyBytes)
	} else {
		request := gorequest.New().Timeout(10*time.Second)
		// resp, bodyBytes, request_err := request.Proxy(proxies.GrabResidentialProxy()).Get(url).
		resp, bodyBytes, request_err := request.Get(url).
		// Set("user-agent", requests.RandomUserAgent()).
		Set("user-agent", requests.RandomPhoneUserAgent()).
		// Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.122 Safari/537.36").
		Set("connection", "keep-alive").
		Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		// Set("accept-encoding", "gzip, deflate, br").
		Set("accept-language", "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7").
		Set("cache-control", "max-age=0").
		// Set("cookie", "__utmz=74692624.1588174094.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utma=74692624.1052763426.1588174094.1588648181.1588839568.7; __utmc=74692624; __utmt=1; shoppingSessionId=1588839580822; lastVisitedFragment=products/173300; __utmb=74692624.7.10.1588839568; _ticket=d0e32ac6c9e4c1edb732ec8a3bf7040bb04664143fa6476c428fa226190cac5e579b1fe58ca40baf9fc147f02735043845d3116c503a897cbcc1683840f31ba81588839595").
		Set("cookie", "__utmz=74692624.1588174094.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utma=74692624.1052763426.1588174094.1588648181.1588839568.7; __utmc=74692624; __utmt=1; shoppingSessionId=1588839580822; lastVisitedFragment=products/173300; __utmb=74692624." + "18" + ".10.1588839568; _ticket=" + utils.GetMD5Hash(strconv.Itoa(rand.Int())) + "b04664143fa6476c428fa226190cac5e579b1fe58ca40baf9fc147f02735043845d3116c503a897cbcc1683840f31ba" + "41588839595").
		Set("dnt", "1").
		Set("referer", "https://www.supremenewyork.com/mobile/").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-origin").
		Set("sec-fetch-user", "?1").
		Set("x-requested-with", "XMLHttpRequest").
		Set("upgrade-insecure-requests", "1").
		EndBytes()
		if request_err != nil {
			log.Println(request_err)
			return
		}

		// ########################################### HANDLE ERRORS
		if requests.ParseHTTPErrors(resp, bodyBytes, identifier, store.URL, true) {
			failed_connections++
			return
		}
	}

	// ########################################### INITIAL CHECK
	initialChecked := restockModeEnabled
	if utils.StringInSlice(identifier, initialCheckedURLs) {
		if !restockModeEnabled {
			log.Println(Green("[" + store.Name + "] " + "Successful connection."))
		} else {
			log.Println(Green("[" + store.Name + " Restocks" + "] " + "Successful connection."))
		}
		initialChecked = true
	} else {
		initialCheckedURLs = append(initialCheckedURLs, identifier)
		if !restockModeEnabled {
			log.Println(Inverse("[" + store.Name + "] " + "Initial Check Done."))
		} else {
			log.Println(Inverse("[" + store.Name + " Restocks" + "] " + "Initial Check Done."))
		}
	}
	successful_connections++

	// ########################################### HANDLE RESPONSE
	data := &structs.SupremeProduct{}
	err := json.Unmarshal(bodyBytes, data)
	if err != nil {
		log.Println(err)
		return
	}

	// fmt.Printf("%+v\n", data.Categories.New[0]) // print first product

	// LOAD NEW ITEMS FIRST
	fetchProducts(data, "New", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome) // only allow restocks for new items
	if restockModeEnabled {
		return;
	}
	fetchProducts(data, "Bags", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Pants", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Accessories", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Skate", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Shoes", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Hats", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Shirts", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Sweatshirts", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Tops_Sweaters", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "Jackets", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
	fetchProducts(data, "T_Shirts", identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)

}

func fetchProducts(data *structs.SupremeProduct, category string, identifier string, store structs.Store, mongoClient *mongo.Client, pusherClient *pusher.Client, initialChecked bool, restockModeEnabled bool, useChrome bool) {
	cur_category := data.Categories.Bags
	if category == "Bags" {
		cur_category = data.Categories.Bags
	} else if category == "Pants" {
		cur_category = data.Categories.Pants
	} else if category == "Accessories" {
		cur_category = data.Categories.Accessories
	} else if category == "Skate" {
		cur_category = data.Categories.Skate
	} else if category == "Shoes" {
		cur_category = data.Categories.Shoes
	} else if category == "Hats" {
		cur_category = data.Categories.Hats
	} else if category == "Shirts" {
		cur_category = data.Categories.Shirts
	} else if category == "Sweatshirts" {
		cur_category = data.Categories.Sweatshirts
	} else if category == "Tops_Sweaters" {
		cur_category = data.Categories.Tops_Sweaters
	} else if category == "Jackets" {
		cur_category = data.Categories.Jackets
	} else if category == "T_Shirts" {
		cur_category = data.Categories.T_Shirts
	} else if category == "New" {
		cur_category = data.Categories.New
	}
	// var products structs.Products
	for i := 0;  i < len(cur_category); i++ {
		product_url := "https://www." + store.URL + "/shop/new/" + strconv.Itoa(cur_category[i].ID)

		// DATABASE DETECTION TO SEE WHETHER OR NOT THE ITEM IS ALREADY IN THE DATABASE. IF NOT: scrape full product.
		if database.CheckIfProductInDatabase(product_url, identifier, mongoClient) {
			if !restockModeEnabled {
				continue
			} else {
				// time.Sleep(135 * time.Millisecond)
			}
		}

		// if category == "New" {
		// 	fetchProduct(data, product_url, cur_category[i], identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
		// } else {
			go fetchProduct(data, product_url, cur_category[i], identifier, store, mongoClient, pusherClient, initialChecked, restockModeEnabled, useChrome)
		// }
	}
}

func fetchProduct(data *structs.SupremeProduct, product_url string, cur_product structs.SupremeCategory, identifier string, store structs.Store, mongoClient *mongo.Client, pusherClient *pusher.Client, initialChecked bool, restockModeEnabled bool, useChrome bool) {

	var bodyBytes []byte

	// ########################################### START REQUEST
	if useChrome {
	  headers := map[string]interface{}{
	    "user-agent": requests.RandomPhoneUserAgent(),
	    "connection": "keep-alive",
	    "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	    "accept-language": "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7",
	    "cache": "max-age=0",
	    "cookie": "__utmz=74692624.1588174094.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utma=74692624.1052763426.1588174094.1588648181.1588839568.7; __utmc=74692624; __utmt=1; shoppingSessionId=1588839580822; lastVisitedFragment=products/173300; __utmb=74692624." + "18" + ".10.1588839568; _ticket=" + utils.GetMD5Hash(strconv.Itoa(rand.Int())) + "b04664143fa6476c428fa226190cac5e579b1fe58ca40baf9fc147f02735043845d3116c503a897cbcc1683840f31ba" + "41588839595",
	    "dnt": "1",
	    "referer": "https://www.supremenewyork.com/mobile/",
	    "sec-fetch-dest": "empty",
	    "sec-fetch-mode": "cors",
	    "sec-fetch-site": "same-origin",
	    "sec-fetch-user": "?1",
	    "x-requested-with": "XMLHttpRequest",
	    "upgrade-insecure-requests": "1",
	  }

		scrapeSupremeWithChrome(product_url + ".json", headers, &bodyBytes)
	} else {
		request := gorequest.New().Timeout(5*time.Second)
		// resp, bodyBytes, request_err := request.Proxy(proxies.GrabResidentialProxy()).Get(product_url + ".json").
		resp, bodyBytes, request_err := request.Get(product_url + ".json").
		// Set("User-Agent", requests.RandomUserAgent()).
		Set("User-Agent", requests.RandomUserAgent()).
		Set("connection", "keep-alive").
		Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		// Set("accept-encoding", "gzip, deflate, br").
		Set("accept-language", "es-US,es;q=0.9,en;q=0.8,en-US;q=0.7").
		Set("cache-control", "max-age=0").
		// Set("cookie", "__utmz=74692624.1588174094.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utma=74692624.1052763426.1588174094.1588174094.1588258943.2; __utmc=74692624; __utmt=1; __utmb=74692624.10.9.1588260542137; _ticket=28517f2a0e86d63bcf273829b8be317694885a8904de8d96ce83f180fafe15cbd480ec1758e855e0d46c22a6ebd49cd981ca340eff12e3eb4ad968a09645911a1588260558").
		Set("dnt", "1").
		Set("referer", "https://www.supremenewyork.com/mobile/").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-origin").
		Set("sec-fetch-user", "?1").
		Set("x-requested-with", "XMLHttpRequest").
		Set("upgrade-insecure-requests", "1").
		EndBytes()
		if request_err != nil {
			log.Println(request_err)
			return /*products*/
		}

		// ########################################### HANDLE ERRORS
		if requests.ParseHTTPErrors(resp, bodyBytes, identifier, store.URL, true) {
			failed_connections++
			return /*products*/
		}
	}

	// ########################################### HANDLE RESPONSE
	product_data := &structs.SupremeProductDetailed{}
	err := json.Unmarshal(bodyBytes, product_data)
	if err != nil {
		log.Println(err)
		return /*products*/
	}

	var products structs.Products
	for j := 0;  j < len(product_data.Styles); j++ {
		var product structs.Product
		// == important ==
		product.Store = store.URL
		product.StoreName = store.Name

		// == main info ==
		product.Name = cur_product.Name + " - " + product_data.Styles[j].Name
		product.URL = product_url + "/" + strconv.Itoa(product_data.Styles[j].ID)
		product.Price = fmt.Sprintf("%.2f", cur_product.Price/100) + " " + product_data.Styles[j].Currency
		product.ImageURL = "https:" + product_data.Styles[j].ImageURL
		product.Description = product_data.Description
		productAvailable := false
		for k := 0;  k < len(product_data.Styles[j].Variants); k++ {
			variant := structs.Variant{
				product_data.Styles[j].Variants[k].Name,
				strconv.Itoa(product_data.Styles[j].Variants[k].ID),
				product_data.Styles[j].Variants[k].Quantity > 0,
				product.Price,
				product_data.Styles[j].Variants[k].Quantity,
			}
			product.Variants = append(product.Variants, variant)
			if !productAvailable && product_data.Styles[j].Variants[k].Quantity > 0 {
				productAvailable = true
				// break // break is only necessary if variants are NOT setup in this nested for loop
			}
		}
		product.Available = productAvailable
		product.Identifier = identifier

		// == extra ==
		product.Color = product_data.Styles[j].Name
		convertedVariants, _ := json.Marshal(product.Variants)
		product.MD5 = utils.GetMD5Hash(string(convertedVariants))
		products = append(products, product)
	}
	// send after each product
	database.SendToDatabase(products, identifier, store.URL, store.Name, initialChecked, mongoClient, pusherClient)
}

func scrapeSupremeWithChrome(url string, headers map[string]interface{}, bodyBytes *[]byte) {
	chromeContext, cancelContext := chromedp.NewContext(context.Background())
	defer cancelContext()

	var response string
	var statusCode int64
	var responseHeaders map[string]interface{}

	runError := chromedp.Run(
		chromeContext,
		runWithTimeout(&chromeContext, 5, chromeTask(
			chromeContext, url, headers,
			&response, &statusCode, &responseHeaders)))

	if runError != nil {
		log.Println(runError)
		return
	}

	// log.Printf(
	// 	"\n\n{%s}\n\n > %s\n status: %d\nheaders: %s\n\n",
	// 	response, url, statusCode, responseHeaders)

	*bodyBytes = []byte(response)

	// ########################################### HANDLE ERRORS
	if requests.ParseHTTPErrorsRaw(statusCode, *bodyBytes, "supreme", "supremenewyork.com", true) {
		failed_connections++
		return
	}
}

func chromeTask(chromeContext context.Context, url string, requestHeaders map[string]interface{}, response *string, statusCode *int64, responseHeaders *map[string]interface{}) chromedp.Tasks {
  chromedp.ListenTarget(chromeContext, func(event interface{}) {
    switch responseReceivedEvent := event.(type) {
    case *network.EventResponseReceived:
      response := responseReceivedEvent.Response
      if response.URL == url {
        *statusCode = response.Status
        *responseHeaders = response.Headers
      }
    }
  })

  return chromeScrape(url, requestHeaders, response)
}

func chromeScrape(url string, headers map[string]interface{}, str *string) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(headers)),
		// emulate iPhone 7 landscape
		chromedp.Emulate(device.IPhone7landscape),
		chromedp.Navigate(url),

		chromedp.InnerHTML("pre", str), // get contents
	}
}

func runWithTimeout(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout * time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}
