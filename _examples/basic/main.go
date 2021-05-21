package main

import (
	"log"
	"net/url"

	"github.com/clarke94/crawler"
)

func main() {
	u, err := url.Parse("http://example.com")
	if err != nil {
		log.Fatalln(err)
	}

	c := crawler.New()

	err = c.Crawl(u)
	if err != nil {
		log.Println(err)
	}
}
