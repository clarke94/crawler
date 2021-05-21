package get

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/clarke94/crawler/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewGet(t *testing.T) {
	tests := []struct {
		name string
		want *Get
	}{
		{
			name: "expect Get Requester to initialize",
			want: &Get{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(Get{})) {
				t.Error(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(Get{})))
			}
		})
	}
}

func TestGet_Request(t *testing.T) {
	tests := []struct {
		name     string
		givenURL string
		want     *http.Request
		wantErr  error
	}{
		{
			name:     "expect a GET request for the given URL",
			givenURL: "http://localhost",
			want: &http.Request{
				Method:     "GET",
				URL:        testutil.URLMustParse("http://localhost"),
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{},
				Host:       "localhost",
			},
			wantErr: nil,
		},
		{
			name:     "expect error given an invalid URL",
			givenURL: "%",
			want:     nil,
			wantErr:  ErrRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()

			got, err := r.Request(context.Background(), tt.givenURL, nil)
			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(http.Request{})) {
				t.Fatal(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(http.Request{})))
			}
		})
	}
}

func TestGet_Do(t *testing.T) {
	tests := []struct {
		name         string
		givenRequest *http.Request
		givenHandler testutil.Handler
		want         []byte
		wantErr      error
	}{
		{
			name:         "expect response from valid request",
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "http://localhost:8080/", nil),
			givenHandler: testutil.Handler{
				Pattern: "/",
				HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
					_, err := rw.Write([]byte("foo"))
					if err != nil {
						return
					}
				},
			},
			want:    []byte("foo"),
			wantErr: nil,
		},
		{
			name:         "expect empty body given status response",
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "http://localhost:8080", nil),
			givenHandler: testutil.Handler{
				Pattern: "/",
				HandlerFunc: func(rw http.ResponseWriter, rr *http.Request) {
					rw.WriteHeader(http.StatusInternalServerError)
				},
			},
			want:    []byte{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := testutil.NewTestServer(tt.givenHandler)
			defer ts.Close()

			r := New()

			got, err := r.Do(tt.givenRequest)
			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			defer got.Close()

			body, err := ioutil.ReadAll(got)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(body, tt.want, cmp.AllowUnexported(http.Request{})) {
				t.Fatal(cmp.Diff(body, tt.want, cmp.AllowUnexported(http.Request{})))
			}
		})
	}
}

func TestGet_Do_Fail(t *testing.T) {
	tests := []struct {
		name         string
		givenRequest *http.Request
		givenHandler testutil.Handler
		want         []byte
		wantErr      error
	}{
		{
			name:         "expect error when context canceled",
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "http://localhost:8080", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			req := tt.givenRequest.WithContext(ctx)

			cancel()

			_, err := r.Do(req)
			if err == nil {
				t.Error("Expected error")
			}
		})
	}
}
