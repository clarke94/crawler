package print

import (
	"errors"
	"net/url"
	"testing"

	"github.com/clarke94/crawler/internal/testutil"
)

func TestPrint_Error(t *testing.T) {
	tests := []struct {
		name       string
		givenError error
	}{
		{
			name:       "expect log to print without panic",
			givenError: errors.New("foo"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Print{}

			p.Error(tt.givenError)
		})
	}
}

func TestPrint_Info(t *testing.T) {
	tests := []struct {
		name         string
		givenVisited url.URL
		givenFound   []*url.URL
	}{
		{
			name:         "expect log to print without panic",
			givenVisited: *testutil.URLMustParse("http://localhost"),
			givenFound:   []*url.URL{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Print{}

			p.Info(&tt.givenVisited, tt.givenFound)
		})
	}
}
