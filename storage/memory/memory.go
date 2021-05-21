package memory

import (
	"errors"
	"net/url"
)

var ErrInvalidURL = errors.New("invalid url")

// Memory is a Storer that stores urls in memory.
type Memory struct {
	data map[url.URL]bool
}

// New initializes a new Memory Storer.
func New() *Memory {
	return &Memory{
		data: map[url.URL]bool{},
	}
}

// Read returns all visited urls in memory.
func (m *Memory) Read() (map[url.URL]bool, error) {
	return m.data, nil
}

// Write stores the provided urls in memory.
func (m *Memory) Write(visited *url.URL) error {
	if visited == nil {
		return ErrInvalidURL
	}

	m.data[*visited] = true

	return nil
}
