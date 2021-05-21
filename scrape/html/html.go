package html

import (
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// HTML is a Scraper that scrapes urls from HTML.
type HTML struct{}

// New initializes a new HTML Scraper.
func New() *HTML {
	return &HTML{}
}

// Scrape extracts all HTML URLs from a reader.
func (o *HTML) Scrape(req *http.Request, closer io.ReadCloser) ([]*url.URL, error) {
	defer closer.Close()

	var links []*url.URL

	z := html.NewTokenizer(closer)

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			return links, nil
		}

		if tt == html.StartTagToken {
			token := z.Token()
			if token.Data != "a" {
				continue
			}

			for _, attr := range token.Attr {
				if attr.Key != "href" {
					continue
				}

				v, err := url.Parse(attr.Val)
				if err != nil {
					continue
				}

				u := req.URL.ResolveReference(v)

				links = append(links, u)
			}
		}
	}
}
