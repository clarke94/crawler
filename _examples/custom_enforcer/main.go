package main

import (
	"log"
	"net/url"

	"github.com/clarke94/crawler"
)

const maxDepth = 5

func main() {
	u, err := url.Parse("http://example.com")
	if err != nil {
		log.Fatalln(err)
	}

	c := crawler.New(
		crawler.WithEnforcer(&CustomEnforcer{}),
	)

	err = c.Crawl(u)
	if err != nil {
		log.Println(err)
	}
}

// CustomEnforcer provides a custom enforcer with max depth.
type CustomEnforcer struct {
	MaxDepth int
}

// Enforce enforces a max depth less than 5.
func (c *CustomEnforcer) Enforce(_ map[url.URL]bool, _ *url.URL) bool {
	c.MaxDepth++
	return c.MaxDepth < maxDepth
}
