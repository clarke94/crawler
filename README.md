# Crawler

Crawler is a Go package that crawls a given URL, extracts all the links from the response and logs the request. Any extracted link that has not yet been visited and has the same domain will be crawled recursively.

## Example

A basic example that will crawl the given URL and visit each link that is under the same domain, logging where it
visited and what links it found.

```go
package main

import (
	"net/url"

	"github.com/clarke94/crawler"
)

func main() {
	u, _ := url.Parse("http://example.com")

	c := crawler.New()

	_ = c.Crawl(u)
}
```

> More examples of extending and customising Crawler in `_examples/` 
