package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/clarke94/crawler"
)

func main() {
	u, err := url.Parse("http://example.com")
	if err != nil {
		log.Fatalln(err)
	}

	c := crawler.New(
		crawler.WithLogger(&CustomLogger{}),
	)

	err = c.Crawl(u)
	if err != nil {
		log.Println(err)
	}
}

// CustomLogger provides a custom logger.
type CustomLogger struct{}

// Error prints a custom error.
func (c CustomLogger) Error(err error) {
	fmt.Printf("Custom: %v", err)
}

// Info prints custom info.
func (c CustomLogger) Info(visited *url.URL, found []*url.URL) {
	fmt.Printf("Custom: %s %v", visited.String(), found)
}
