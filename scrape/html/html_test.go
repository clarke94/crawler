package html

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/clarke94/crawler/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHTML_Scrape(t *testing.T) {
	tests := []struct {
		name         string
		givenCloser  io.ReadCloser
		givenRequest *http.Request
		want         []*url.URL
		wantErr      error
	}{
		{
			name:         "expect no links when reader empty",
			givenCloser:  io.NopCloser(strings.NewReader("")),
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "http://localhost", nil),
			want:         nil,
			wantErr:      nil,
		},
		{
			name:         "expect link when provided in reader",
			givenCloser:  io.NopCloser(strings.NewReader(`<a href="https://example.com"></a>`)),
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "https://example.com", nil),
			want: []*url.URL{
				testutil.URLMustParse("https://example.com"),
			},
			wantErr: nil,
		},
		{
			name:         "expect continue when a non anchor provided",
			givenCloser:  io.NopCloser(strings.NewReader(`<div><a href="https://example.com"></a></div>`)),
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "https://example.com", nil),
			want: []*url.URL{
				testutil.URLMustParse("https://example.com"),
			},
			wantErr: nil,
		},
		{
			name:         "expect link when different attribute are present",
			givenCloser:  io.NopCloser(strings.NewReader(`<a data-id="foo" href="https://example.com" class="bar"></a>`)),
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "https://example.com", nil),
			want: []*url.URL{
				testutil.URLMustParse("https://example.com"),
			},
			wantErr: nil,
		},
		{
			name:         "expect continue when invalid URL",
			givenCloser:  io.NopCloser(strings.NewReader(`<a href="%"></a><a href="https://example.com"></a>`)),
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "https://example.com", nil),
			want: []*url.URL{
				testutil.URLMustParse("https://example.com"),
			},
			wantErr: nil,
		},
		{
			name:         "expect relative route URL",
			givenCloser:  io.NopCloser(strings.NewReader(`<a href="%"></a><a href="/"></a>`)),
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "https://example.com", nil),
			want: []*url.URL{
				testutil.URLMustParse("https://example.com/"),
			},
			wantErr: nil,
		},
		{
			name:         "expect relative path URL",
			givenCloser:  io.NopCloser(strings.NewReader(`<a href="/foo"></a>`)),
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "https://example.com", nil),
			want: []*url.URL{
				testutil.URLMustParse("https://example.com/foo"),
			},
			wantErr: nil,
		},
		{
			name:         "expect URL to stay the same given external URL",
			givenCloser:  io.NopCloser(strings.NewReader(`<a href="https://foo.com/foo"></a>`)),
			givenRequest: testutil.HTTPMustRequests(context.Background(), http.MethodGet, "https://example.com", nil),
			want: []*url.URL{
				testutil.URLMustParse("https://foo.com/foo"),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &HTML{}
			got, err := o.Scrape(tt.givenRequest, tt.givenCloser)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(got, tt.want) {
				t.Fatal(cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestNewHTML(t *testing.T) {
	tests := []struct {
		name string
		want *HTML
	}{
		{
			name: "expect HTML Scraper to initialize",
			want: &HTML{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()

			if !cmp.Equal(got, tt.want) {
				t.Errorf(cmp.Diff(got, tt.want))
			}
		})
	}
}
