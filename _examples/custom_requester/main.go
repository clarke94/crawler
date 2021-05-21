package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/clarke94/crawler"
)

func main() {
	u, err := url.Parse("http://example.com")
	if err != nil {
		log.Fatalln(err)
	}

	c := crawler.New(
		crawler.WithRequester(&CustomRequester{}),
	)

	err = c.Crawl(u)
	if err != nil {
		log.Println(err)
	}
}

// CustomRequester provides a custom requester.
type CustomRequester struct{}

// Request attaches a custom context with timeout of 10 seconds to a request.
func (c CustomRequester) Request(_ context.Context, rawURL string, _ io.Reader) (*http.Request, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
}

// Do sends the request with a custom client.
func (c *CustomRequester) Do(req *http.Request) (io.ReadCloser, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
