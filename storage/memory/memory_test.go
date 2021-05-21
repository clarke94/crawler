package memory

import (
	"net/url"
	"testing"

	"github.com/clarke94/crawler/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMemory_Read(t *testing.T) {
	tests := []struct {
		name      string
		givenData map[url.URL]bool
		want      map[url.URL]bool
		wantErr error
	}{
		{
			name: "expect all data on Read",
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("http://localhost"):    true,
				*testutil.URLMustParse("https://example.com"): true,
			},
			want: map[url.URL]bool{
				*testutil.URLMustParse("http://localhost"):    true,
				*testutil.URLMustParse("https://example.com"): true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				data: tt.givenData,
			}
			got, err := m.Read()
			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestMemory_Write(t *testing.T) {
	tests := []struct {
		name         string
		givenVisited *url.URL
		givenData    map[url.URL]bool
		wantErr error
	}{
		{
			name:         "expect write when data is empty",
			givenVisited: testutil.URLMustParse("http://localhost"),
			givenData:    map[url.URL]bool{},
		},
		{
			name:         "expect write when data is already in",
			givenVisited: testutil.URLMustParse("http://localhost"),
			givenData: map[url.URL]bool{
				*testutil.URLMustParse("http://localhost"): true,
			},
		},
		{
			name:         "expect error when nil URL provided",
			givenVisited: nil,
			givenData:    map[url.URL]bool{},
			wantErr:      ErrInvalidURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				data: tt.givenData,
			}

			err := m.Write(tt.givenVisited)

			if tt.givenVisited != nil {
				if _, ok := m.data[*tt.givenVisited]; !ok {
					t.Error(cmp.Diff(tt.givenVisited, m.data))
				}
			}

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestNewMemory(t *testing.T) {
	tests := []struct {
		name string
		want *Memory
	}{
		{
			name: "expect Memory Storer to initialize",
			want: &Memory{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreUnexported(Memory{})) {
				t.Error(cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(Memory{})))
			}
		})
	}
}
