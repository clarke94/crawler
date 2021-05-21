package main

import (
	"log"
	"net/url"

	"github.com/clarke94/crawler"
	"github.com/clarke94/crawler/enforce/samedomainonce"
)

const maxDepth = 10

func main() {
	u, err := url.Parse("http://example.com")
	if err != nil {
		log.Fatalln(err)
	}

	c := crawler.New(
		crawler.WithEnforcer(NewEnforcer()),
	)

	err = c.Crawl(u)
	if err != nil {
		log.Println(err)
	}
}

// CustomEnforcer provides a custom enforcer that decorates the default Crawler enforcer.
type CustomEnforcer struct {
	MaxDepth int
	enforcer crawler.Enforcer
}

// NewEnforcer initializes a new custom enforcer.
func NewEnforcer() *CustomEnforcer {
	return &CustomEnforcer{
		MaxDepth: 0,
		enforcer: samedomainonce.New(),
	}
}

// Enforce enforces a max depth less than 10 for the same domain that is only visited once.
func (c *CustomEnforcer) Enforce(data map[url.URL]bool, u *url.URL) bool {
	c.MaxDepth++
	if c.MaxDepth >= maxDepth {
		return false
	}

	return c.enforcer.Enforce(data, u)
}
