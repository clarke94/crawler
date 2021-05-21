package samedomainonce

import (
	"net/url"
	"testing"

	"github.com/clarke94/crawler/internal/testutil"
	"github.com/clarke94/crawler/storage/memory"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewSameDomain(t *testing.T) {
	tests := []struct {
		name string
		want *SameDomainOnce
	}{
		{
			name: "expect SameDomainOnce Enforcer to initialize",
			want: &SameDomainOnce{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(SameDomainOnce{}, memory.Memory{})) {
				t.Error(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(SameDomainOnce{}, memory.Memory{})))
			}
		})
	}
}

func TestSameDomain_Enforce(t *testing.T) {
	tests := []struct {
		name      string
		givenData map[url.URL]bool
		givenURL  url.URL
		want      bool
	}{
		{
			name:      "expect enforcer true given no URL has been visited",
			givenData: nil,
			givenURL:  *testutil.URLMustParse("https://example.com"),
			want:      true,
		},
		{
			name: "expect enforcer false given a URL already visited",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://example.com"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com"),
			want:     false,
		},
		{
			name: "expect enforcer false given a URL with a different domain",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("http://localhost"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com"),
			want:     false,
		},
		{
			name: "expect enforcer false given a URL with a different sub-domain",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://sub.example.com"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com"),
			want:     false,
		},
		{
			name: "expect enforcer true given URL with a different protocol",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("http://localhost/some/url"): true,
			},
			givenURL: *testutil.URLMustParse("https://localhost"),
			want:     true,
		},
		{
			name: "expect enforcer false given the same URL with a trailing slash stored",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://example.com/"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com"),
			want:     false,
		},
		{
			name: "expect enforcer false given the same URL with a trailing slash found",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://example.com"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com/"),
			want:     false,
		},
		{
			name: "expect enforcer false given the same URL with a path and trailing slash stored",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://example.com/some/path/"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com/some/path"),
			want:     false,
		},
		{
			name: "expect enforcer false given the same URL with a path and trailing slash found",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://example.com/some/path"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com/some/path/"),
			want:     false,
		},
		{
			name: "expect enforcer false given the same URL with a trailing reference",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://example.com/foo"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com/foo#something"),
			want:     false,
		},
		{
			name: "expect enforcer false given the same URL with a trailing reference after slash",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://example.com/foo"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com/foo/#something"),
			want:     false,
		},
		{
			name: "expect enforcer false given the same URL with a query param",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("https://example.com/foo"): true,
			},
			givenURL: *testutil.URLMustParse("https://example.com/foo?foo=bar&bar=foo"),
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()

			got := s.Enforce(tt.givenData, &tt.givenURL)

			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}
		})
	}
}
