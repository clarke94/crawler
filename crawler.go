package crawler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/clarke94/crawler/enforce/samedomainonce"
	print2 "github.com/clarke94/crawler/log/print"
	"github.com/clarke94/crawler/request/get"
	"github.com/clarke94/crawler/scrape/html"
	"github.com/clarke94/crawler/storage/memory"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidURL is the returned error for a URL that
	// cannot be parse or is nil.
	ErrInvalidURL = errors.New("invalid URL")
	// ErrRequester is the annotated error that is wrapped with
	// the returned error from the Requester.
	ErrRequester = errors.New("requester error")
	// ErrScraper is the annotated error that is wrapped with
	// the returned error from the Scraper.
	ErrScraper = errors.New("scraper error")
	// ErrStorer is the annotated error that is wrapped with
	// the returned error from the Storer.
	ErrStorer = errors.New("storage error")
)

// Storer provides an interface to the storage layer.
type Storer interface {
	Read() (map[url.URL]bool, error)
	Write(visited *url.URL) error
}

// Scraper provides an interface to extract data and return urls.
type Scraper interface {
	Scrape(req *http.Request, closer io.ReadCloser) ([]*url.URL, error)
}

// Requester provides the interface for a HTTP Request.
type Requester interface {
	Request(ctx context.Context, rawURL string, body io.Reader) (*http.Request, error)
	Do(req *http.Request) (io.ReadCloser, error)
}

// Logger provides the interface to log output from the Crawler.
type Logger interface {
	Error(err error)
	Info(visited *url.URL, found []*url.URL)
}

// Enforcer provides an interface to enforce logic before scraping.
type Enforcer interface {
	Enforce(data map[url.URL]bool, url *url.URL) bool
}

// Option is a functional option to modify the default Crawler instance.
type Option func(crawler *Crawler)

// Crawler provides a web crawler.
type Crawler struct {
	Context context.Context

	requester Requester
	enforcer  Enforcer
	scraper   Scraper
	storer    Storer
	logger    Logger
	wg        *sync.WaitGroup
	mu        *sync.RWMutex
	errMu     *sync.RWMutex
	err       error
}

// New initializes a new default Crawler.
func New(options ...Option) *Crawler {
	c := &Crawler{
		Context: context.Background(),

		requester: get.New(),
		enforcer:  samedomainonce.New(),
		scraper:   html.New(),
		storer:    memory.New(),
		logger:    &print2.Print{},
		wg:        &sync.WaitGroup{},
		mu:        &sync.RWMutex{},
		errMu:     &sync.RWMutex{},
		err:       nil,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// WithEnforcer replaces the default enforcer with the provided one.
func WithEnforcer(enforcer Enforcer) Option {
	return func(c *Crawler) {
		c.enforcer = enforcer
	}
}

// WithLogger replaces the default logger with the provided one.
func WithLogger(logger Logger) Option {
	return func(c *Crawler) {
		c.logger = logger
	}
}

// WithRequester replaces the default requester with the provided one.
func WithRequester(requester Requester) Option {
	return func(c *Crawler) {
		c.requester = requester
	}
}

// WithScraper replaces the default scraper with the provided one.
func WithScraper(scraper Scraper) Option {
	return func(c *Crawler) {
		c.scraper = scraper
	}
}

// WithStorer replaces the default storer with the provided one.
func WithStorer(storer Storer) Option {
	return func(c *Crawler) {
		c.storer = storer
	}
}

// Crawl sends a request to a given URL and scrapes data from the response.
func (c *Crawler) Crawl(u *url.URL) error {
	if u == nil {
		return ErrInvalidURL
	}

	c.goCrawl(u)
	c.wg.Wait()

	c.errMu.RLock()
	defer c.errMu.RUnlock()

	return c.err
}

// goCrawl is a concurrent wrapper around the crawl method.
func (c *Crawler) goCrawl(u *url.URL) {
	c.wg.Add(1)

	go func(c *Crawler) {
		if err := c.crawl(u); err != nil {
			c.errMu.Lock()
			defer c.errMu.Unlock()
			c.err = err
		}
	}(c)
}

// crawl checks the URL with the enforcer to see if the conditions are met
// and then invokes the requester to create and send the request.
// The request is stored in the storer and the response is passed
// to the scraper to extract the data and return found URLs, the
// request is passed to the logger and any found urls are recursively crawled.
func (c *Crawler) crawl(u *url.URL) error {
	defer c.wg.Done()

	c.mu.Lock()
	ok, err := c.check(u)
	c.mu.Unlock()

	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	req, err := c.requester.Request(c.Context, u.String(), nil)
	if err != nil {
		return errors.Wrap(ErrRequester, err.Error())
	}

	resp, err := c.requester.Do(req)
	if err != nil {
		return errors.Wrap(ErrRequester, err.Error())
	}

	urls, err := c.scraper.Scrape(req, resp)
	if err != nil {
		return errors.Wrap(ErrScraper, err.Error())
	}

	c.logger.Info(req.URL, urls)

	for _, foundURL := range urls {
		c.goCrawl(foundURL)
	}

	return nil
}

func (c *Crawler) check(u *url.URL) (bool, error) {
	visitedURLs, err := c.storer.Read()
	if err != nil {
		return false, errors.Wrap(ErrStorer, err.Error())
	}

	if ok := c.enforcer.Enforce(visitedURLs, u); !ok {
		return false, nil
	}

	if writeErr := c.storer.Write(u); writeErr != nil {
		return false, errors.Wrap(ErrStorer, writeErr.Error())
	}

	return true, nil
}
