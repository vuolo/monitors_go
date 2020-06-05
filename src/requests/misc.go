package requests

import(
  "fmt"
  "bufio"
  "net"
  "time"
  "net/http"
  "strings"

  "github.com/valyala/fasthttp" // fast http
)

// ########################################### UTILITY FUNCTIONS
func ConvertCookies(cookies_string string) []*http.Cookie {
  var cookies []*http.Cookie
  split_cookies := strings.Split(strings.ReplaceAll(cookies_string, "\"", "'"), "; ")
  for i := 0;  i < len(split_cookies); i++ {
    cookies = append(cookies, &http.Cookie{Name: strings.Split(split_cookies[i], "=")[0], Value:strings.Split(split_cookies[i], "=")[1]})
  }
  return cookies
}

func FasthttpHTTPDialer(proxyAddr string) fasthttp.DialFunc {
  return func(addr string) (net.Conn, error) {
    conn, err := fasthttp.DialDualStackTimeout(proxyAddr, 15 * time.Second) //fasthttp.DialTimeout(proxyAddr, 15 * time.Second) //fasthttp.Dial(proxyAddr)
    if err != nil {
      return nil, err
    }

    req := "CONNECT " + addr + " HTTP/1.1\r\n"
    // req += "Proxy-Authorization: xxx\r\n"
    req += "\r\n"

    if _, err := conn.Write([]byte(req)); err != nil {
      return nil, err
    }

    res := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(res)

    res.SkipBody = true

    if err := res.Read(bufio.NewReader(conn)); err != nil {
      conn.Close()
      return nil, err
    }
    if res.Header.StatusCode() != 200 {
      conn.Close()
      return nil, fmt.Errorf("could not connect to proxy")
    }
    return conn, nil
  }
}
