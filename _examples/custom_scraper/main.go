package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/clarke94/crawler"
	"golang.org/x/net/html"
)

func main() {
	u, err := url.Parse("http://example.com")
	if err != nil {
		log.Fatalln(err)
	}

	c := crawler.New(
		crawler.WithScraper(&CustomScraper{}),
	)

	err = c.Crawl(u)
	if err != nil {
		log.Println(err)
	}
}

// CustomScraper provides a custom scraper.
type CustomScraper struct{}

// Scrape scrapes all text and prints it to the console.
func (c CustomScraper) Scrape(_ *http.Request, closer io.ReadCloser) ([]*url.URL, error) {
	defer closer.Close()

	var links []*url.URL

	z := html.NewTokenizer(closer)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			return links, nil
		}

		if tt == html.TextToken {
			token := z.Token()

			fmt.Println(strings.TrimSpace(token.Data))
		}
	}
}
