package database

import (
  "context"
  "log"
  "strconv"
  // "strings"
  "math"

  "go.mongodb.org/mongo-driver/mongo" // MongoDB
	"go.mongodb.org/mongo-driver/mongo/options" // MongoDB Options
	"go.mongodb.org/mongo-driver/bson" // MongoDB BSON
	"go.mongodb.org/mongo-driver/x/bsonx" // MongoDB BSONx

  "github.com/pusher/pusher-http-go" // Pusher

  . "github.com/logrusorgru/aurora" // colors

  "github.com/dghubble/go-twitter/twitter" // Twitter

  // ## local shits ##
	"../structs"
  "../requests"
)

// ########################################### UTILITY FUNCTIONS
func SendToDatabase(products structs.Products, identifier string, store_url string, store_name string, initialChecked bool, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  collection := mongoClient.Database("monitors").Collection(identifier + "_products")

  // // ########################################### INITIAL CHECK (cached_products approach)
  // if !utils.StringInSlice(store_name, initialCheckedURLs) {
  // 	products = utils.QuicksortProducts(products) // Method for local variables from initial check
  // 	// for _, product := range products {
  // 	// 	log.Println(product)
  // 	// }
  // 	cached_products = products
  // 	initialCheckedURLs = append(initialCheckedURLs, store_url)
  // 	log.Println(Inverse("[" + store_name + "] " + "Initial Check Done."))
  // 	return
  // }
  //
  // // ########################################### HANDLE FOUND PRODUCTS (cached_products approach)
  // for i := 0;  i < len(cached_products); i++ { // Method for local variables from initial check
  // 	// log.Println(results[i])
  // 	if cached_products[i].Store != products[i].Store {
  // 		continue
  // 	}
  // 	if cached_products[i].URL != products[i].URL {
  // 		indifferent_products = append(indifferent_products, products[i])
  // 		log.Println(Inverse("[" + store_name + "] " + "New Item Found: " + products[i].Name))
  // 		i++ // offset to prevent update for all subsequent products
  // 	} else if cached_products[i].URL == products[i].URL && (cached_products[i].Available != products[i].Available || cached_products[i].LaunchDate != products[i].LaunchDate) {
  // 		stock_update_products = append(stock_update_products, products[i])
  // 		log.Println(Inverse("[" + store_name + "] " + "New Stock Update Found: " + products[i].Name))
  // 	}
  // }

  // ########################################### FIND PRODUCTS FROM DB
  // specify the Sort option to sort the returned documents by url in ascending order
  // find_opts := options.Find().SetSort(bson.D{{"url", 1}}) // Method for local variables from initial check
  cursor, find_error := collection.Find(context.TODO(), bson.D{{"store", store_url}})//, find_opts)
  if find_error != nil {
    log.Println(find_error)
    return
  }

  // get a list of all returned documents
  var results []bson.M
  if find_error = cursor.All(context.TODO(), &results); find_error != nil {
    log.Println(find_error)
    return
  }

  // ########################################### HANDLE FOUND PRODUCTS
  var indifferent_products structs.Products // all new products
  var stock_update_products structs.Products // all products that availability has changed or published_at is different
  var stock_update_products_availabilities structs.Products // all products that md5 has changed
  for i := 0;  i < len(products); i++ {
    // log.Println(products[i])
    foundProduct := false
    foundProduct_obj := bson.M{}
    for j := 0;  j < len(results); j++ {
      if results[j]["url"] == products[i].URL {
        foundProduct = true
        foundProduct_obj = results[j]
        break
      }
    }
    if !foundProduct {
      indifferent_products = append(indifferent_products, products[i])
      if initialChecked {
        log.Println(Inverse("[" + store_name + "] " + "New Product Found: " + products[i].Name))
      }
    } else if foundProduct_obj["available"] != products[i].Available || foundProduct_obj["launchdate"] != products[i].LaunchDate {
      stock_update_products = append(stock_update_products, products[i])
      if initialChecked {
        log.Println(Inverse("[" + store_name + "] " + "New Stock Update Found: " + products[i].Name))
      }
    } else if foundProduct_obj["md5"] != products[i].MD5 {
      stock_update_products_availabilities = append(stock_update_products, products[i])
      if initialChecked {
        log.Println(Inverse("[" + store_name + "] " + "Availability Update Found: " + products[i].Name))
      }
    }
  }

  if len(indifferent_products) > 0 {

    if initialChecked {
      for i := 0;  i < len(indifferent_products); i++ {
        go requests.PostProduct(indifferent_products[i], true, pusherClient)
      }
    }

    // ########################################### SETUP DB INDEX
    // Set index (unique variable for url)
    _, index_err := collection.Indexes().CreateOne(
      context.Background(),
      mongo.IndexModel{
        Keys   : bsonx.Doc{{"url", bsonx.Int32(1)}},
        Options: options.Index().SetUnique(true),
      },
    )

    if index_err != nil {
      log.Println(index_err)
      return
    }

    // ########################################### INSERT MANY TO DB
    output_products := []interface{}{} // TODO: figure out a one-line Products to []interface converter?
    for i := 0;  i < len(indifferent_products); i++ {
      output_products = append(output_products, indifferent_products[i])
    }
    _, insert_err := collection.InsertMany(context.TODO(), output_products, options.InsertMany().SetOrdered(false))
    if insert_err != nil {
        // log.Println(insert_err)
        // return
    } else {
      // fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
    }
  }

  // ########################################### UPDATE PRODUCT IN DB
  for i := 0;  i < len(stock_update_products); i++ {
    if initialChecked {
      go requests.PostProduct(stock_update_products[i], false, pusherClient)
    }
    update_opts := options.FindOneAndUpdate().SetUpsert(true)
    filter := bson.D{{"url", stock_update_products[i].URL}}
    update := bson.D{{"$set", stock_update_products[i]}}
    var updatedDocument bson.M
    update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
    if update_err != nil {
        // ErrNoDocuments means that the filter did not match any documents in the collection
        if update_err == mongo.ErrNoDocuments {
            // return
            log.Println("no documents to update found!")
        } else {
          log.Println(update_err)
          return
        }
    } else {
      // log.Printf("updated document %v", updatedDocument)
    }
  }

  // ########################################### UPDATE PRODUCT IN DB (AVAILABILITIES)
  for i := 0;  i < len(stock_update_products_availabilities); i++ {
    if initialChecked {
      go requests.UpdateProduct(stock_update_products_availabilities[i], pusherClient)
    }
    update_opts := options.FindOneAndUpdate().SetUpsert(true)
    filter := bson.D{{"url", stock_update_products_availabilities[i].URL}}
    update := bson.D{{"$set", stock_update_products_availabilities[i]}}
    var updatedDocument bson.M
    update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
    if update_err != nil {
        // ErrNoDocuments means that the filter did not match any documents in the collection
        if update_err == mongo.ErrNoDocuments {
            // return
            log.Println("no documents to update found!")
        } else {
          log.Println(update_err)
          return
        }
    } else {
      // log.Printf("updated document %v", updatedDocument)
    }
  }

	return
}

func CheckIfProductInDatabase(product_url string, identifier string, mongoClient *mongo.Client) bool {

  collection := mongoClient.Database("monitors").Collection(identifier + "_products")

  // ########################################### FIND PRODUCTS FROM DB
  cursor, find_error := collection.Find(context.TODO(), bson.D{{"url", bson.D{{"$regex", ".*" + product_url + ".*"}}}})
  if find_error != nil {
    log.Println(find_error)
    return true
  }

  // get a list of all returned documents
  var results []bson.M
  if find_error = cursor.All(context.TODO(), &results); find_error != nil {
    log.Println(find_error)
    return true
  }

	return len(results) > 0
}

func TryUpdatePasswordBook(store_name string, store_url string, passwordEnabled bool, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  collection := mongoClient.Database("monitors").Collection("password_book")

  // ########################################### FIND SITE IN DB
  var result bson.M
  find_err := collection.FindOne(context.TODO(), bson.D{{"store", store_url}}).Decode(&result)
  if find_err != nil {
    // ErrNoDocuments means that the filter did not match any documents in the collection
    if find_err == mongo.ErrNoDocuments {
      // ########################################### UPDATE SITE IN DB
      update_opts := options.FindOneAndUpdate().SetUpsert(true)
      filter := bson.D{{"store", store_url}}
      update := bson.D{{"$set", bson.M{"password_enabled": passwordEnabled}}}
      var updatedDocument bson.M
      update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
      if update_err != nil {
          // ErrNoDocuments means that the filter did not match any documents in the collection
          if update_err == mongo.ErrNoDocuments {
              // return
              log.Println(Cyan("**ADDING " + store_url + " TO PASSWORD BOOK DATABASE**"))
          } else {
            log.Println(update_err)
            return
          }
      } else {
        // log.Printf("updated document %v", updatedDocument)
      }
      return
    }
    log.Println(find_err)
    return
  }

  // ########################################### HANDLE SITE IN DB
  if result["password_enabled"] == passwordEnabled {
    return
  }

  // ########################################### UPDATE SITE IN DB
  update_opts := options.FindOneAndUpdate().SetUpsert(true)
  filter := bson.D{{"store", store_url}}
  update := bson.D{{"$set", bson.M{"password_enabled": passwordEnabled}}}
  var updatedDocument bson.M
  update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
  if update_err != nil {
      // ErrNoDocuments means that the filter did not match any documents in the collection
      if update_err == mongo.ErrNoDocuments {
          // return
          log.Println(Cyan("**ADDING " + store_url + " TO PASSWORD BOOK DATABASE**"))
      } else {
        log.Println(update_err)
        return
      }
  } else {
    // log.Printf("updated document %v", updatedDocument)
    go requests.PostPassword(store_name, store_url, passwordEnabled)
  }
}

func TryUpdateCheckpointBook(store_name string, store_url string, checkpointEnabled bool, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  collection := mongoClient.Database("monitors").Collection("checkpoint_book")

  // ########################################### FIND SITE IN DB
  var result bson.M
  find_err := collection.FindOne(context.TODO(), bson.D{{"store", store_url}}).Decode(&result)
  if find_err != nil {
    // ErrNoDocuments means that the filter did not match any documents in the collection
    if find_err == mongo.ErrNoDocuments {
      // ########################################### UPDATE SITE IN DB
      update_opts := options.FindOneAndUpdate().SetUpsert(true)
      filter := bson.D{{"store", store_url}}
      update := bson.D{{"$set", bson.M{"checkpoint_enabled": checkpointEnabled}}}
      var updatedDocument bson.M
      update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
      if update_err != nil {
          // ErrNoDocuments means that the filter did not match any documents in the collection
          if update_err == mongo.ErrNoDocuments {
              // return
              log.Println(Cyan("**ADDING " + store_url + " TO CHECKPOINT BOOK DATABASE**"))
          } else {
            log.Println(update_err)
            return
          }
      } else {
        // log.Printf("updated document %v", updatedDocument)
      }
      return
    }
    log.Println(find_err)
    return
  }

  // ########################################### HANDLE SITE IN DB
  if result["checkpoint_enabled"] == checkpointEnabled {
    return
  }

  // ########################################### UPDATE SITE IN DB
  update_opts := options.FindOneAndUpdate().SetUpsert(true)
  filter := bson.D{{"store", store_url}}
  update := bson.D{{"$set", bson.M{"checkpoint_enabled": checkpointEnabled}}}
  var updatedDocument bson.M
  update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
  if update_err != nil {
      // ErrNoDocuments means that the filter did not match any documents in the collection
      if update_err == mongo.ErrNoDocuments {
          // return
          log.Println(Cyan("**ADDING " + store_url + " TO CHECKPOINT BOOK DATABASE**"))
      } else {
        log.Println(update_err)
        return
      }
  } else {
    // log.Printf("updated document %v", updatedDocument)
    go requests.PostCheckpoint(store_name, store_url, checkpointEnabled)
  }
}

func SendToStockXDatabase(product structs.StockXProduct, identifier string, initialChecked bool, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  collection := mongoClient.Database("monitors").Collection(identifier + "_products")

  // ########################################### SETUP DB INDEX
  // Set index (unique variable for handle)
  _, index_err := collection.Indexes().CreateOne(
    context.Background(),
    mongo.IndexModel{
      Keys   : bsonx.Doc{{"handle", bsonx.Int32(1)}},
      Options: options.Index().SetUnique(true),
    },
  )

  if index_err != nil {
    log.Println(index_err)
    return
  }

  // ########################################### FIND PRODUCT FROM DB
  var result structs.StockXProductBSON
  find_error := collection.FindOne(context.TODO(), bson.D{{"handle", product.Handle}}).Decode(&result)
  if find_error != nil {
    // ErrNoDocuments means that the filter did not match any documents in the collection
    if find_error == mongo.ErrNoDocuments {
      // log.Println("Product not found in database... Adding now.")
      // ########################################### INSERT PRODUCT TO DB
      _, insert_err := collection.InsertOne(context.TODO(), product)
      if insert_err != nil {
          // log.Println(insert_err)
          // return
      } else {
        // fmt.Println("Inserted document: ", res)
      }
      return
    }
    log.Println(find_error)
    return
  }

  // ########################################### HANDLE FOUND PRODUCT
  needUpdate := false
  for i := 0;  i < len(product.Variants); i++ {
    for j := 0;  j < len(result.Variants); j++ { // slow method of finding indifferences. Better method is sorting both variants arrays and comparing by a single index
      if product.Variants[i].Name != result.Variants[j].Name {
        continue
      }

      notifyDescription := ""
      if product.Variants[i].LowestAsk != result.Variants[j].LowestAsk {
        needUpdate = true // force update
        // 5% checker under previous lowest ask
        if product.Variants[i].LowestAsk < result.Variants[j].LowestAsk {
          percentageDiff_lowest_ask := math.Round((((result.Variants[j].LowestAsk - product.Variants[i].LowestAsk) / result.Variants[j].LowestAsk) * 100) * 10) / 10 // convert to whole number
          if percentageDiff_lowest_ask > 5 {
            // add to description
            notifyDescription += "> Lowest Ask is **" + strconv.Itoa(int(percentageDiff_lowest_ask)) + "%** ($" + strconv.FormatFloat(result.Variants[j].LowestAsk - product.Variants[i].LowestAsk, 'f', 2, 64) + ") lower than previous Lowest Ask. ($" + strconv.FormatFloat(result.Variants[j].LowestAsk, 'f', 2, 64) + ")\n"
          }
        }
        // 5% checker under retail price
        percentageDiff_retail := math.Round((((product.Variants[i].RetailPrice - product.Variants[i].LowestAsk) / product.Variants[i].RetailPrice) * 100) * 10) / 10 // convert to whole number
        if percentageDiff_retail > 5 {
          // add to description
          notifyDescription += "> Lowest Ask is **" + strconv.Itoa(int(percentageDiff_retail)) + "%** ($" + strconv.FormatFloat(product.Variants[i].RetailPrice - product.Variants[i].LowestAsk, 'f', 2, 64) + ") lower than Retail Price. ($" + strconv.FormatFloat(product.Variants[i].RetailPrice, 'f', 2, 64) + ")\n"
        }
      }
      if product.Variants[i].LastSale != result.Variants[j].LastSale {
        needUpdate = true // force update
        // 5% checker under last sale
        if product.Variants[i].LastSale < result.Variants[j].LastSale {
          percentageDiff_last_sale := math.Round((((result.Variants[j].LastSale - product.Variants[i].LastSale) / result.Variants[j].LastSale) * 100) * 10) / 10 // convert to whole number
          if percentageDiff_last_sale > 5 {
            // add to description
            notifyDescription += "> Last Sale is **" + strconv.Itoa(int(percentageDiff_last_sale)) + "%** ($" + strconv.FormatFloat(result.Variants[j].LastSale - product.Variants[i].LastSale, 'f', 2, 64) + ") lower than previous Last Sale. ($" + strconv.FormatFloat(result.Variants[j].LastSale, 'f', 2, 64) + ")\n"
          }
        }
      }
      if notifyDescription != "" {
        notifyDescription = "**Size: " + product.Variants[i].Name + "**\n**BUY NOW 📲 __$" + strconv.FormatFloat(product.Variants[i].LowestAsk, 'f', 2, 64) + "__**\n" + notifyDescription // exactly what to put in description
        if initialChecked {
          go requests.PostStockX(product, product.Variants[i], notifyDescription)
        }
        needUpdate = true
      }

    }
  }

  // ########################################### UPDATE PRODUCT IN DB
  if !needUpdate {
    return
  }
  update_opts := options.FindOneAndUpdate().SetUpsert(true)
  filter := bson.D{{"handle", product.Handle}}
  update := bson.D{{"$set", product}}
  var updatedDocument bson.M
  update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
  if update_err != nil {
      // ErrNoDocuments means that the filter did not match any documents in the collection
      if update_err == mongo.ErrNoDocuments {
          // return
          log.Println("no documents to update found!")
      } else {
        log.Println(update_err)
        return
      }
  } else {
    // log.Printf("updated document %v", updatedDocument)
  }

	return
}

func SendToTwitterDatabase(twitterUser structs.TwitterUser, identifier string, initialChecked bool, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  collection := mongoClient.Database("monitors").Collection(identifier + "_users")

  // ########################################### SETUP DB INDEX
  // Set index (unique variable for handle)
  _, index_err := collection.Indexes().CreateOne(
    context.Background(),
    mongo.IndexModel{
      Keys   : bsonx.Doc{{"username", bsonx.Int32(1)}},
      Options: options.Index().SetUnique(true),
    },
  )

  if index_err != nil {
    log.Println(index_err)
    return
  }

  // ########################################### FIND TWITTER USER FROM DB
  var resultUser structs.TwitterUser
  find_error := collection.FindOne(context.TODO(), bson.D{{"username", twitterUser.Username}}).Decode(&resultUser)
  if find_error != nil {
    // ErrNoDocuments means that the filter did not match any documents in the collection
    if find_error == mongo.ErrNoDocuments {
      // log.Println("Twitter user not found in database... Adding now.")
      // ########################################### INSERT TWITTER USER TO DB
      _, insert_err := collection.InsertOne(context.TODO(), twitterUser)
      if insert_err != nil {
          // log.Println(insert_err)
          // return
      } else {
        // fmt.Println("Inserted document: ", res)
        if initialChecked {
          log.Println(Inverse("[" + "Twitter: " + twitterUser.Username + "] " + "New tweet found."))
          go requests.PostTwitterTweet(twitterUser, twitterUser.Tweets[0], pusherClient)
        }
      }
      return
    }
    log.Println(find_error)
    return
  }

  // ########################################### HANDLE FOUND TWITTER USER
  var indifferentTweets []twitter.Tweet
  for i := 0;  i < len(twitterUser.Tweets); i++ {
    foundTweet := false
    for j := 0;  j < len(resultUser.Tweets); j++ {
      if resultUser.Tweets[j].ID == twitterUser.Tweets[i].ID {
        foundTweet = true
        break
      }
    }
    if !foundTweet {
      indifferentTweets = append(indifferentTweets, twitterUser.Tweets[i])
    }
  }

  needUpdate := false
  if resultUser.ImageURL != twitterUser.ImageURL {
    if initialChecked {
      go requests.PostTwitterProfilePicture(twitterUser, pusherClient)
    }
    log.Println(Inverse("[" + "Twitter: " + twitterUser.Username + "] " + "Profile picture updated."))
    needUpdate = true
  }
  if resultUser.FullName != twitterUser.FullName {
    if initialChecked {
      go requests.PostTwitterFullName(twitterUser, pusherClient)
    }
    log.Println(Inverse("[" + "Twitter: " + twitterUser.Username + "] " + "Full name updated."))
    needUpdate = true
  }
  if resultUser.Biography != twitterUser.Biography {
    if initialChecked {
      go requests.PostTwitterBiography(twitterUser, pusherClient)
    }
    log.Println(Inverse("[" + "Twitter: " + twitterUser.Username + "] " + "Biography updated."))
    needUpdate = true
  }
  if resultUser.ExternalURL != twitterUser.ExternalURL {
    if initialChecked {
      go requests.PostTwitterExternalURL(twitterUser, pusherClient)
    }
    log.Println(Inverse("[" + "Twitter: " + twitterUser.Username + "] " + "External URL updated."))
    needUpdate = true
  }
  for i := 0;  i < len(indifferentTweets); i++ {
    // resultUser.Tweets = append(resultUser.Tweets, indifferentTweets[i]) // keep extending DB stored tweets...
    resultUser.Tweets = append([]twitter.Tweet{indifferentTweets[i]}, resultUser.Tweets...) // keep extending DB stored tweets...
    if initialChecked {
      log.Println(Inverse("[" + "Twitter: " + twitterUser.Username + "] " + "New tweet found."))
      go requests.PostTwitterTweet(twitterUser, indifferentTweets[i], pusherClient)
    }
    needUpdate = true
    if i == len(indifferentTweets)-1 { // keep extending DB stored tweets...
      twitterUser.Tweets = resultUser.Tweets // keep extending DB stored tweets...
    } // keep extending DB stored tweets...
  }

  // ########################################### UPDATE TWITTER USER IN DB
  if !needUpdate {
    return
  }
  update_opts := options.FindOneAndUpdate().SetUpsert(true)
  filter := bson.D{{"username", twitterUser.Username}}
  update := bson.D{{"$set", twitterUser}}
  var updatedDocument bson.M
  update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
  if update_err != nil {
      // ErrNoDocuments means that the filter did not match any documents in the collection
      if update_err == mongo.ErrNoDocuments {
          // return
          log.Println("no documents to update found!")
      } else {
        log.Println(update_err)
        return
      }
  } else {
    // log.Printf("updated document %v", updatedDocument)
  }

	return
}

func SendToInstagramDatabase(instagramUser structs.InstagramUser, identifier string, initialChecked bool, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  collection := mongoClient.Database("monitors").Collection(identifier + "_users")

  // ########################################### SETUP DB INDEX
  // Set index (unique variable for handle)
  _, index_err := collection.Indexes().CreateOne(
    context.Background(),
    mongo.IndexModel{
      Keys   : bsonx.Doc{{"username", bsonx.Int32(1)}},
      Options: options.Index().SetUnique(true),
    },
  )

  if index_err != nil {
    log.Println(index_err)
    return
  }

  // ########################################### FIND INSTAGRAM USER FROM DB
  var resultUser structs.InstagramUser
  find_error := collection.FindOne(context.TODO(), bson.D{{"username", instagramUser.Username}}).Decode(&resultUser)
  if find_error != nil {
    // ErrNoDocuments means that the filter did not match any documents in the collection
    if find_error == mongo.ErrNoDocuments {
      // log.Println("Instagram user not found in database... Adding now.")
      // ########################################### INSERT INSTAGRAM USER TO DB
      _, insert_err := collection.InsertOne(context.TODO(), instagramUser)
      if insert_err != nil {
          // log.Println(insert_err)
          // return
      } else {
        // fmt.Println("Inserted document: ", res)
        if initialChecked {
          log.Println(Inverse("[" + "Instagram: " + instagramUser.Username + "] " + "New post found."))
          go requests.PostInstagramPost(instagramUser, instagramUser.Posts[0], pusherClient)
        }
      }
      return
    }
    log.Println(find_error)
    return
  }

  // ########################################### HANDLE FOUND INSTAGRAM USER
  var indifferentPosts []structs.InstagramPost
  for i := 0;  i < len(instagramUser.Posts); i++ {
    foundPost := false
    for j := 0;  j < len(resultUser.Posts); j++ {
      if resultUser.Posts[j].ID == instagramUser.Posts[i].ID {
        foundPost = true
        break
      }
    }
    if !foundPost {
      indifferentPosts = append(indifferentPosts, instagramUser.Posts[i])
    }
  }

  needUpdate := false
  // resultUser_shortenedImageUrl := resultUser.ImageURL
  // instagramUser_shortenedImageUrl := instagramUser.ImageURL
  // if strings.Index(resultUser_shortenedImageUrl, "?") > -1 {
  //   resultUser_shortenedImageUrl = resultUser_shortenedImageUrl[:strings.Index(resultUser_shortenedImageUrl, "?")]
  // }
  // if strings.Index(instagramUser_shortenedImageUrl, "?") > -1 {
  //   instagramUser_shortenedImageUrl = instagramUser_shortenedImageUrl[:strings.Index(instagramUser_shortenedImageUrl, "?")]
  // }
  // if resultUser_shortenedImageUrl != instagramUser_shortenedImageUrl {
  //   if initialChecked {
  //     go requests.PostInstagramProfilePicture(instagramUser)
  //   }
  //   log.Println(Inverse("[" + "Instagram: " + instagramUser.Username + "] " + "Profile picture updated."))
  //   needUpdate = true
  // }
  if resultUser.FullName != instagramUser.FullName {
    if initialChecked {
      go requests.PostInstagramFullName(instagramUser, pusherClient)
    }
    log.Println(Inverse("[" + "Instagram: " + instagramUser.Username + "] " + "Full name updated."))
    needUpdate = true
  }
  if resultUser.Biography != instagramUser.Biography {
    if initialChecked {
      go requests.PostInstagramBiography(instagramUser, pusherClient)
    }
    log.Println(Inverse("[" + "Instagram: " + instagramUser.Username + "] " + "Biography updated."))
    needUpdate = true
  }
  if resultUser.ExternalURL != instagramUser.ExternalURL {
    if initialChecked {
      go requests.PostInstagramExternalURL(instagramUser, pusherClient)
    }
    log.Println(Inverse("[" + "Instagram: " + instagramUser.Username + "] " + "External URL updated."))
    needUpdate = true
  }
  for i := 0;  i < len(indifferentPosts); i++ {
    // resultUser.Posts = append(resultUser.Posts, indifferentPosts[i]) // keep extending DB stored posts...
    resultUser.Posts = append([]structs.InstagramPost{indifferentPosts[i]}, resultUser.Posts...) // keep extending DB stored tweets...
    if initialChecked {
      log.Println(Inverse("[" + "Instagram: " + instagramUser.Username + "] " + "New post found."))
      go requests.PostInstagramPost(instagramUser, indifferentPosts[i], pusherClient)
    }
    needUpdate = true
    if i == len(indifferentPosts)-1 { // keep extending DB stored posts...
      instagramUser.Posts = resultUser.Posts // keep extending DB stored posts...
    } // keep extending DB stored posts...
  }

  // ########################################### UPDATE INSTAGRAM USER IN DB
  if !needUpdate {
    return
  }
  update_opts := options.FindOneAndUpdate().SetUpsert(true)
  filter := bson.D{{"username", instagramUser.Username}}
  update := bson.D{{"$set", instagramUser}}
  var updatedDocument bson.M
  update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
  if update_err != nil {
      // ErrNoDocuments means that the filter did not match any documents in the collection
      if update_err == mongo.ErrNoDocuments {
          // return
          log.Println("no documents to update found!")
      } else {
        log.Println(update_err)
        return
      }
  } else {
    // log.Printf("updated document %v", updatedDocument)
  }

	return
}

func SendToInstagramStoryDatabase(instagramUser structs.InstagramUser, identifier string, initialChecked bool, mongoClient *mongo.Client, pusherClient *pusher.Client) {

  collection := mongoClient.Database("monitors").Collection(identifier + "_users")

  // ########################################### SETUP DB INDEX
  // Set index (unique variable for handle)
  _, index_err := collection.Indexes().CreateOne(
    context.Background(),
    mongo.IndexModel{
      Keys   : bsonx.Doc{{"username", bsonx.Int32(1)}},
      Options: options.Index().SetUnique(true),
    },
  )

  if index_err != nil {
    log.Println(index_err)
    return
  }

  // ########################################### FIND INSTAGRAM USER FROM DB
  var resultUser structs.InstagramUser
  find_error := collection.FindOne(context.TODO(), bson.D{{"username", instagramUser.Username}}).Decode(&resultUser)
  if find_error != nil {
    // ErrNoDocuments means that the filter did not match any documents in the collection
    if find_error == mongo.ErrNoDocuments {
      // log.Println("Instagram user not found in database... Adding now.")
      // ########################################### INSERT INSTAGRAM USER TO DB
      _, insert_err := collection.InsertOne(context.TODO(), instagramUser)
      if insert_err != nil {
          // log.Println(insert_err)
          // return
      } else {
        // fmt.Println("Inserted document: ", res)
        if initialChecked {
          go requests.PostInstagramStory(instagramUser, instagramUser.Stories[0], pusherClient)
        }
      }
      return
    }
    log.Println(find_error)
    return
  }

  // ########################################### HANDLE FOUND INSTAGRAM USER
  var indifferentStories []structs.InstagramStory
  for i := 0;  i < len(instagramUser.Stories); i++ {
    foundStory := false
    for j := 0;  j < len(resultUser.Stories); j++ {
      if resultUser.Stories[j].ID == instagramUser.Stories[i].ID {
        foundStory = true
        break
      }
    }
    if !foundStory {
      indifferentStories = append(indifferentStories, instagramUser.Stories[i])
    }
  }

  needUpdate := false
  for i := 0;  i < len(indifferentStories); i++ {
    resultUser.Stories = append(resultUser.Stories, indifferentStories[i]) // keep extending DB stored stories...
    if initialChecked {
      go requests.PostInstagramStory(instagramUser, indifferentStories[i], pusherClient)
    }
    needUpdate = true
    if i == len(indifferentStories)-1 { // keep extending DB stored stories...
      instagramUser.Stories = resultUser.Stories // keep extending DB stored stories...
    } // keep extending DB stored stories...
  }

  // ########################################### UPDATE INSTAGRAM USER IN DB
  if !needUpdate {
    return
  }
  update_opts := options.FindOneAndUpdate().SetUpsert(true)
  filter := bson.D{{"username", instagramUser.Username}}
  update := bson.D{{"$set", instagramUser}}
  var updatedDocument bson.M
  update_err := collection.FindOneAndUpdate(context.TODO(), filter, update, update_opts).Decode(&updatedDocument)
  if update_err != nil {
      // ErrNoDocuments means that the filter did not match any documents in the collection
      if update_err == mongo.ErrNoDocuments {
          // return
          log.Println("no documents to update found!")
      } else {
        log.Println(update_err)
        return
      }
  } else {
    // log.Printf("updated document %v", updatedDocument)
  }

	return
}

func GatherSocialPlusActiveHandles(mongoClient *mongo.Client) []bson.M {
  collection := mongoClient.Database("monitors").Collection("social_media_handles")

  // ########################################### FIND PRODUCTS FROM DB
  cursor, find_error := collection.Find(context.TODO(), bson.D{{"concurrent_connections_length", bson.D{{"$gt", 0}}}})
  if find_error != nil {
    log.Println(find_error)
    return []bson.M{}
  }

  // get a list of all returned documents
  var results []bson.M
  if find_error = cursor.All(context.TODO(), &results); find_error != nil {
    log.Println(find_error)
    return []bson.M{}
  }

	return results
}
