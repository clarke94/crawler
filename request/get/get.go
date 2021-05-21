package get

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	// ErrRequest is returned when a HTTP request cannot be
	// parsed with the given URL.
	ErrRequest = errors.New("unable to parse request")
	// ErrDo is returned when a given request is unable to
	// be sent by the client.
	ErrDo = errors.New("unable to send request")
)

const timeout = 10

// Get is a Requester for HTTP GET requests.
type Get struct {
	client *http.Client
}

// New initializes a new Get Requester.
func New() *Get {
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	return &Get{
		client: client,
	}
}

// Request parses the given URL and returns a HTTP request.
func (r *Get) Request(ctx context.Context, rawURL string, _ io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, ErrRequest
	}

	return req, nil
}

// Do sends a HTTP request and returns the response.
func (r *Get) Do(req *http.Request) (io.ReadCloser, error) {
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, ErrDo
	}

	return resp.Body, nil
}
