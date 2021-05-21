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

	c := crawler.New(
		crawler.WithStorer(New()),
	)

	err = c.Crawl(u)
	if err != nil {
		log.Println(err)
	}
}

// CustomStorer provides a custom storer.
type CustomStorer struct {
	example map[url.URL]bool
}

// New initializes a new Storer.
func New() *CustomStorer {
	return &CustomStorer{
		example: map[url.URL]bool{},
	}
}

// Read from a custom database.
func (c *CustomStorer) Read() (map[url.URL]bool, error) {
	return c.example, nil
}

// Write to a custom database.
func (c *CustomStorer) Write(visited *url.URL) error {
	c.example[*visited] = true

	return nil
}
