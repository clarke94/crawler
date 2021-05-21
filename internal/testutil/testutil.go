package testutil

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type Handler struct {
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewTestServer(handlers ...Handler) *httptest.Server {
	mux := http.NewServeMux()

	for _, h := range handlers {
		mux.HandleFunc(h.Pattern, h.HandlerFunc)
	}

	l, _ := net.Listen("tcp", "127.0.0.1:8080")

	ts := httptest.NewUnstartedServer(mux)

	_ = ts.Listener.Close()

	ts.Listener = l

	ts.Start()

	return ts
}

func URLMustParse(rawURL string) *url.URL {
	u, _ := url.Parse(rawURL)
	return u
}

func HTTPMustRequests(ctx context.Context, method, rawURL string, body io.Reader) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, method, rawURL, body)
	return req
}
