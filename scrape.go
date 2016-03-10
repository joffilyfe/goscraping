package main

import (
  "fmt"
  "log"
  "strings"
  "github.com/PuerkitoBio/goquery"
  "gopkg.in/olivere/elastic.v3"
)

type News struct  {
  Title string `json:"title"`
  Text string `json:"text"`
  Url string `json:url`    
}

var count = 0

func Indexer(new *News) {
  client, err := elastic.NewClient()
  if err != nil {
    panic(err)
  }

  _, err = client.Index().
    Index("scrape").
    Type("news").
    BodyJson(new).
    Do()
}


func ParseMain(mainUrl string) {
	doc, err := goquery.NewDocument(mainUrl) 
	if err != nil {
		log.Fatal(err)
	}

  doc.Find(".vevent.contenttype-news-item").Each(func(i int, s *goquery.Selection) {
    url, _ := s.Find("a").Attr("href")
    ParseNews(url)
  })
}

func ParseNews(url string) {
	document, err := goquery.NewDocument(url)
  defer count += 1

	if err != nil {
		log.Fatal(err)
	}

	title := document.Find("h1 span").Text()
  text := ""
    
  document.Find("#parent-fieldname-text p").Each(func(i int, s *goquery.Selection) {   
    text = text + s.Text()  
  })
    
  news := News{
    Title: title,
    Text: text,
    Url: url,
  }

  Indexer(&news)
  fmt.Printf("> %d) %s.\n> URL: %s\n", count, strings.Trim(title, " \r \n"), url)
}

func main() {
  for i := 0; i < 25; i += 5 {
    root := "http://www.ifpb.edu.br/"
    path := "campi/noticias-campi/folder_summary_view?b_start:int="
    fullUrl := root + path + fmt.Sprintf("%d", i)
    ParseMain(fullUrl)
  }
}