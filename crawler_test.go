package crawler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/clarke94/crawler/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Crawler
	}{
		{
			name: "expect crawler to initialize",
			want: &Crawler{
				Context: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(Crawler{})) {
				t.Error(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(Crawler{})))
			}
		})
	}
}

func TestCrawler_Crawl_Success(t *testing.T) {
	tests := []struct {
		name          string
		givenURL      *url.URL
		givenHandlers []testutil.Handler
		wantErr       error
	}{
		{
			name:     "given a URL with no links, expect no output",
			givenURL: testutil.URLMustParse("http://127.0.0.1:8080/"),
			givenHandlers: []testutil.Handler{
				{
					Pattern: "/",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<div>hello world</div>`))
					},
				},
			},
			wantErr: nil,
		},
		{
			name:     "given a URL with a link that is a different domain, expect URL not visited",
			givenURL: testutil.URLMustParse("http://127.0.0.1:8080/"),
			givenHandlers: []testutil.Handler{
				{
					Pattern: "/",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<a href="https://example.com">Example</a>`))
					},
				},
			},
			wantErr: nil,
		},
		{
			name:     "given a URL with a link that is the same domain, expect URL visited",
			givenURL: testutil.URLMustParse("http://127.0.0.1:8080/"),
			givenHandlers: []testutil.Handler{
				{
					Pattern: "/",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<a href="http://127.0.0.1:8080/second">Example</a>`))
					},
				},
				{
					Pattern: "/second",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<div>hello world</div>`))
					},
				},
			},
			wantErr: nil,
		},
		{
			name:     "given a URL with a link that is the same domain, expect URL visited",
			givenURL: testutil.URLMustParse("http://127.0.0.1:8080/"),
			givenHandlers: []testutil.Handler{
				{
					Pattern: "/",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<a href="http://127.0.0.1:8080/second">Example</a>`))
					},
				},
				{
					Pattern: "/second",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<a href="http://127.0.0.1:8080/second">Example</a>`))
					},
				},
			},
			wantErr: nil,
		},
		{
			name:     "given a URL with multiple links, expect all domain URLs visited",
			givenURL: testutil.URLMustParse("http://127.0.0.1:8080/"),
			givenHandlers: []testutil.Handler{
				{
					Pattern: "/",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`
							<a href="http://127.0.0.1:8080/second">Example</a>
							<a href="http://127.0.0.1:8080/">Example</a>
							<a href="http://127.0.0.1:8080/three">Example</a>
							<a href="https://example.com">Example</a>
						`))
					},
				},
				{
					Pattern: "/second",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<a href="http://127.0.0.1:8080/three">Example</a>`))
					},
				},
				{
					Pattern: "/three",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<a href="http://127.0.0.1:8080/">Example</a>`))
					},
				},
			},
			wantErr: nil,
		},
		{
			name:     "given a URL with a broken link, expect continue",
			givenURL: testutil.URLMustParse("http://127.0.0.1:8080/"),
			givenHandlers: []testutil.Handler{
				{
					Pattern: "/",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`
							<a href="http://127.0.0.1:8080/second">Example</a>
							<a href="http://127.0.0.1:8080/not-found">Not found</a>
							<a href="http://127.0.0.1:8080/third">Example</a>
						`))
					},
				},
				{
					Pattern: "/second",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<a href="http://127.0.0.1:8080/second">Example</a>`))
					},
				},
				{
					Pattern: "/third",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`<a href="http://127.0.0.1:8080/not-found">Not found</a>`))
					},
				},
				{
					Pattern: "/not-found",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.WriteHeader(http.StatusNotFound)
					},
				},
			},
			wantErr: nil,
		},
		{
			name:     "given a URL with relative link and query param link, expect continue",
			givenURL: testutil.URLMustParse("http://127.0.0.1:8080/"),
			givenHandlers: []testutil.Handler{
				{
					Pattern: "/",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`
							<a href="http://127.0.0.1:8080/">same</a>
							<a href="http://127.0.0.1:8080/#foo">ref</a>
							<a href="http://127.0.0.1:8080#foo">ref</a>
							<a href="http://127.0.0.1:8080?foo=bar">query</a>
						`))
					},
				},
			},
			wantErr: nil,
		},
		{
			name:     "given a URL with relative link and query param link, expect continue",
			givenURL: testutil.URLMustParse("http://127.0.0.1:8080/contact"),
			givenHandlers: []testutil.Handler{
				{
					Pattern: "/contact",
					HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
						rw.Header().Set("Content-Type", "text/html")
						_, _ = rw.Write([]byte(`
							<a href="http://127.0.0.1:8080/contact#reference">ref</a>
							<a href="http://127.0.0.1:8080/contact/#foo">foo</a>
							<a href="http://127.0.0.1:8080/contact/">Example</a>
							<a href="http://127.0.0.1:8080/contact">Query</a>
						`))
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := testutil.NewTestServer(tt.givenHandlers...)
			defer ts.Close()

			t.Log(tt.name)

			c := New()

			err := c.Crawl(tt.givenURL)
			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			t.Log("\n")
		})
	}
}

func TestCrawler_Crawl_Fail(t *testing.T) {
	testRequest := testutil.HTTPMustRequests(context.Background(), http.MethodGet, "http://localhost", nil)
	tests := []struct {
		name           string
		givenURL       *url.URL
		givenRequester Requester
		givenScraper   Scraper
		givenStorer    Storer
		givenLogger    Logger
		givenEnforcer  Enforcer
		wantErr        error
	}{
		{
			name:     "given a nil pointer URL, expect error",
			givenURL: nil,
			wantErr:  ErrInvalidURL,
		},
		{
			name:     "expect error given fatal requester error",
			givenURL: testutil.URLMustParse("http://localhost"),
			givenRequester: mockRequester{
				GivenRequestError: ErrRequester,
			},
			givenScraper:  mockScraper{},
			givenStorer:   mockStorer{},
			givenLogger:   mockLogger{},
			givenEnforcer: mockEnforcer{
				GivenBool: true,
			},
			wantErr:       ErrRequester,
		},
		{
			name:     "expect error given fatal requester doer error",
			givenURL: testutil.URLMustParse("http://localhost"),
			givenRequester: mockRequester{
				GivenRequest: testRequest,
				GivenDoError: ErrRequester,
			},
			givenScraper:  mockScraper{},
			givenStorer:   mockStorer{},
			givenLogger:   mockLogger{},
			givenEnforcer: mockEnforcer{
				GivenBool: true,
			},
			wantErr:       ErrRequester,
		},
		{
			name:     "expect error given fatal scraper error",
			givenURL: testutil.URLMustParse("http://localhost"),
			givenRequester: mockRequester{
				GivenRequest: testRequest,
			},
			givenScraper: mockScraper{
				GivenError: ErrRequester,
			},
			givenStorer:   mockStorer{},
			givenLogger:   mockLogger{},
			givenEnforcer: mockEnforcer{
				GivenBool: true,
			},
			wantErr:       ErrScraper,
		},
		{
			name:     "expect error given fatal write storage error",
			givenURL: testutil.URLMustParse("http://localhost"),
			givenRequester: mockRequester{
				GivenRequest: testRequest,
			},
			givenScraper: mockScraper{},
			givenStorer: mockStorer{
				GivenWriteError: ErrStorer,
			},
			givenLogger:   mockLogger{},
			givenEnforcer: mockEnforcer{
				GivenBool: true,
			},
			wantErr:       ErrStorer,
		},
		{
			name:     "expect error given fatal read storage error",
			givenURL: testutil.URLMustParse("http://localhost"),
			givenRequester: mockRequester{
				GivenRequest: testRequest,
			},
			givenScraper: mockScraper{
				GivenURLs: []*url.URL{
					testutil.URLMustParse("http://localhost"),
				},
			},
			givenStorer: mockStorer{
				GivenReadError: ErrStorer,
			},
			givenLogger:   mockLogger{},
			givenEnforcer: mockEnforcer{},
			wantErr:       ErrStorer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Crawler{
				requester: tt.givenRequester,
				scraper:   tt.givenScraper,
				storer:    tt.givenStorer,
				logger:    tt.givenLogger,
				enforcer:  tt.givenEnforcer,
				wg:        &sync.WaitGroup{},
				mu:        &sync.RWMutex{},
				errMu:     &sync.RWMutex{},
			}

			err := c.Crawl(tt.givenURL)
			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestCrawler_WithEnforcer(t *testing.T) {
	tests := []struct {
		name          string
		givenEnforcer Enforcer
		want          *Crawler
	}{
		{
			name:          "expect custom enforcer initialize",
			givenEnforcer: mockEnforcer{},
			want: &Crawler{
				enforcer: mockEnforcer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithEnforcer(tt.givenEnforcer))
			if !cmp.Equal(got.enforcer, tt.want.enforcer) {
				t.Error(cmp.Diff(got.enforcer, tt.want.enforcer))
			}
		})
	}
}

func TestCrawler_Withlogger(t *testing.T) {
	tests := []struct {
		name        string
		givenLogger Logger
		want        *Crawler
	}{
		{
			name:        "expect custom logger initialize",
			givenLogger: mockLogger{},
			want: &Crawler{
				logger: mockLogger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithLogger(tt.givenLogger))
			if !cmp.Equal(got.logger, tt.want.logger) {
				t.Error(cmp.Diff(got.logger, tt.want.logger))
			}
		})
	}
}

func TestCrawler_WithRequester(t *testing.T) {
	tests := []struct {
		name           string
		givenRequester Requester
		want           *Crawler
	}{
		{
			name:           "expect custom Requester initialize",
			givenRequester: mockRequester{},
			want: &Crawler{
				requester: mockRequester{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithRequester(tt.givenRequester))
			if !cmp.Equal(got.requester, tt.want.requester) {
				t.Error(cmp.Diff(got.requester, tt.want.requester))
			}
		})
	}
}

func TestCrawler_WithScraper(t *testing.T) {
	tests := []struct {
		name         string
		givenScraper Scraper
		want         *Crawler
	}{
		{
			name:         "expect custom Scraper initialize",
			givenScraper: mockScraper{},
			want: &Crawler{
				scraper: mockScraper{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithScraper(tt.givenScraper))
			if !cmp.Equal(got.scraper, tt.want.scraper) {
				t.Error(cmp.Diff(got.scraper, tt.want.scraper))
			}
		})
	}
}

func TestCrawler_WithStorer(t *testing.T) {
	tests := []struct {
		name        string
		givenStorer Storer
		want        *Crawler
	}{
		{
			name:        "expect custom Storer initialize",
			givenStorer: mockStorer{},
			want: &Crawler{
				storer: mockStorer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithStorer(tt.givenStorer))
			if !cmp.Equal(got.storer, tt.want.storer) {
				t.Error(cmp.Diff(got.storer, tt.want.storer))
			}
		})
	}
}

type mockRequester struct {
	GivenRequestError error
	GivenDoError      error
	GivenRequest      *http.Request
	GivenReadCloser   io.ReadCloser
}

func (m mockRequester) Request(_ context.Context, _ string, _ io.Reader) (*http.Request, error) {
	return m.GivenRequest, m.GivenRequestError
}

func (m mockRequester) Do(_ *http.Request) (io.ReadCloser, error) {
	return m.GivenReadCloser, m.GivenDoError
}

type mockScraper struct {
	GivenError error
	GivenURLs  []*url.URL
}

func (m mockScraper) Scrape(_ *http.Request, _ io.ReadCloser) ([]*url.URL, error) {
	return m.GivenURLs, m.GivenError
}

type mockStorer struct {
	GivenReadError  error
	GivenWriteError error
	GivenNilError   error
	GivenStore      map[url.URL]bool
}

func (m mockStorer) Read() (map[url.URL]bool, error) {
	return m.GivenStore, m.GivenReadError
}

func (m mockStorer) Write(u *url.URL) error {
	if u == nil {
		return m.GivenNilError
	}

	return m.GivenWriteError
}

type mockLogger struct {
}

func (m mockLogger) Error(_ error) {}

func (m mockLogger) Info(_ *url.URL, _ []*url.URL) {}

type mockEnforcer struct {
	GivenBool bool
}

func (m mockEnforcer) Enforce(_ map[url.URL]bool, _ *url.URL) bool {
	return m.GivenBool
}
