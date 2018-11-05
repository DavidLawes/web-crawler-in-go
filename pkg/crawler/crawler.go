package crawler

import (
  "net/http"
  "fmt"
  "web-crawler-in-go/pkg/getBody"
  "web-crawler-in-go/pkg/hrefExtractor"
)

type HttpClient interface {
  Get(string) (*http.Response, error)
}

type UrlMap map[string][]string

func Crawl(seedUrl string, client HttpClient) (urlMap UrlMap) {
  urlMap = make(UrlMap)
  urlQueue := make(chan string)
  urlCrawled := make(chan bool)
  crawlComplete := false
  i := 0

  go func() { urlQueue <- seedUrl }()

  for crawlComplete == false {
    select {
    case url := <- urlQueue:
      go getLinks(url, client, urlMap, urlCrawled, urlQueue)
    case <- urlCrawled:
      fmt.Println("urlCrawled")
      i++
      fmt.Println(i)
      fmt.Println(len(urlQueue))
      fmt.Println(urlMap)
      fmt.Println(len(urlMap))
      if i == len(urlMap) {
        fmt.Println("all urls crawled")
        crawlComplete = true
      }
    }
  }

  return
}

func getLinks(url string, client HttpClient, urlMap UrlMap, urlCrawled chan bool, urlQueue chan string) {
  body := getBody.GetBody(client, url)
  links := hrefExtractor.Extract(body)
  addToMap(url, urlMap, links, urlQueue)
  urlCrawled <- true
}

func addToMap(url string, urlMap UrlMap, links []string, urlQueue chan string) {
  fmt.Println("addToMap")
  urlMap[url] = links
  for _, url := range links {
    if _, ok := urlMap[url]; !ok {
      urlMap[url] = []string{}
      fmt.Println("url added to queue: ", url)
      // only add url to queue if it's not already in map
      urlQueue <- url
    }
  }
}
