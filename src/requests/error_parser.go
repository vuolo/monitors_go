package requests

import (
  "net/http"
  "strings"
  "log"
  "strconv"

  . "github.com/logrusorgru/aurora" // colors

  "github.com/valyala/fasthttp" // fast http
)

// ########################################### UTILITY FUNCTIONS
func ParseHTTPErrors(response *http.Response, responseBody []byte, identifier string, store_url string, displayErrors bool) bool {
  if response.StatusCode == 401 {
    if displayErrors {
      log.Println(Yellow("[" + store_url + "] " + "Password Protected Response."))
    }
    return false
  } else if response.StatusCode == 403 || (identifier == "shopify" && strings.Index(string(responseBody), "large amounts of web requests") > -1) {
    if displayErrors {
      log.Println(Red("[" + store_url + "] " + "Connection failed. (Temp ban occured)"))
    }
    return true
  } else if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
    if displayErrors {
      log.Println(Red("[" + store_url + "] " + "Connection failed. (Invalid site? - Code: " + strconv.Itoa(response.StatusCode) + " - Message: " + string(http.StatusText(response.StatusCode)) + " )"))
    }
    return true
  } else if strings.Index(response.Request.URL.String(), "/password") > -1 {
    if displayErrors {
      log.Println(Yellow("[" + store_url + "] " + "Password Protected Response"))
    }
    return false
	} else if identifier == "shopify" && strings.Index(response.Header.Get("Content-Type"), "json") == -1 {
    if displayErrors {
      log.Println(Red("[" + store_url + "] " + "Connection failed. (Invalid site type? - Not JSON response)"))
    }
    return true
  }
  return false
}

func ParseHTTPErrorsRaw(responseStatusCode int64, responseBody []byte, identifier string, store_url string, displayErrors bool) bool {
  if responseStatusCode == 401 {
    if displayErrors {
      log.Println(Yellow("[" + store_url + "] " + "Password Protected Response."))
    }
    return false
  } else if responseStatusCode == 403 || (identifier == "shopify" && strings.Index(string(responseBody), "large amounts of web requests") > -1) {
    if displayErrors {
      log.Println(Red("[" + store_url + "] " + "Connection failed. (Temp ban occured)"))
    }
    return true
  } else if !(responseStatusCode >= 200 && responseStatusCode <= 299) {
    if displayErrors {
      log.Println(Red("[" + store_url + "] " + "Connection failed. (Invalid site? - Code: " + strconv.FormatInt(responseStatusCode, 10) + " - Message: " + string(http.StatusText(int(responseStatusCode))) + " )"))
    }
    return true
  }
  return false
}

// func ParseFastHTTPErrors(responseStatusCode int, responseBody []byte, identifier string, store_url string) bool {
//   if responseStatusCode == 401 {
//     log.Println(Red("[" + store_url + "] " + "Connection failed. (Unauthorized)"))
//     return true
//   } else if identifier == "shopify" && strings.Index(string(responseBody), "large amounts of web requests") > -1 {
//     log.Println(Red("[" + store_url + "] " + "Connection failed. (Temp ban occured)"))
//     return true
//   } else if identifier == "shopify" && strings.Index(string(responseBody), "/password") > -1 { // TODO: check for redirect URL
//     log.Println(Red("[" + store_url + "] " + "Connection failed. (Temp ban occured)"))
//     return true
//   } else if !(responseStatusCode >= 200 && responseStatusCode <= 299) {
//     log.Println(Red("[" + store_url + "] " + "Connection failed. (Invalid site? - Code: " + strconv.Itoa(responseStatusCode) + " - Message: " + string(http.StatusText(responseStatusCode)) + " )"))
//     return true
//   } else if false {//response.Header.Get("Location") != "" && strings.Index(response.Header.Get("Location"), "/password") > -1 {
// 		log.Println(Red("[" + store_url + "] " + "Password Protected Response."))
//     return false
// 	} else if false {//strings.Index(response.Header.Get("Content-Type"), "json") == -1 {
//     log.Println(Red("[" + store_url + "] " + "Connection failed. (Invalid site type? - Not JSON response)"))
//     return true
//   }
//   return false
// }

func ParseFastHTTPErrorsWithFullResponse(response *fasthttp.Response, responseBody []byte, identifier string, store_url string) bool {
  if response.StatusCode() == 401 {
    log.Println(Yellow("[" + store_url + "] " + "Password Protected Response."))
    return false
  } else if response.StatusCode() == 403 {
    log.Println(Red("[" + store_url + "] " + "Connection failed. (Temp ban occured)"))
    return true
  } else if len(response.Body()) == 0 {
    log.Println(Red("[" + store_url + "] " + "Connection failed. (No body)"))
    return true
  } else if identifier == "shopify" && strings.Index(string(responseBody), "large amounts of web requests") > -1 {
    log.Println(Red("[" + store_url + "] " + "Connection failed. (Temp ban occured)"))
    return true
  } else if identifier == "shopify" && strings.Index(string(responseBody), "/password") > -1 { // TODO: check for redirect URL
    log.Println(Red("[" + store_url + "] " + "Connection failed. (Temp ban occured)"))
    return true
  } else if !(response.StatusCode() >= 200 && response.StatusCode() <= 299) {
    log.Println(Red("[" + store_url + "] " + "Connection failed. (Invalid site? - Code: " + strconv.Itoa(response.StatusCode()) + " - Message: " + string(http.StatusText(response.StatusCode())) + " )"))
    return true
  }
  return false
}
