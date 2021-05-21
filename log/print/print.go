package print

import (
	"fmt"
	"net/url"
)

// Print is a Logger that prints with the standard format package.
type Print struct{}

// Error prints the given error to the console.
func (p *Print) Error(err error) {
	fmt.Println(err)
}

// Info prints the given parameters to the console.
func (p *Print) Info(visited *url.URL, found []*url.URL) {
	fmt.Printf("Visited %s and found %v \n", visited.String(), found)
}
