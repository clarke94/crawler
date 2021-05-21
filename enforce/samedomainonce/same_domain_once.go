package samedomainonce

import (
	"net/url"
	"strings"
)

// SameDomainOnce is an Enforcer that enforces same domain and only visit once.
type SameDomainOnce struct{}

// New initializes a new SameDomainOnce Enforcer.
func New() *SameDomainOnce {
	return &SameDomainOnce{}
}

// Enforce enforces same URL domains and only visit once and checks equality without trailing suffix.
func (s *SameDomainOnce) Enforce(visited map[url.URL]bool, u *url.URL) bool {
	visitedDomain := getDomain(visited)
	if isNothingVisited(visitedDomain) {
		return true
	}

	if !isDomainEqual(*visitedDomain, u.Hostname()) {
		return false
	}

	if isVisited(visited, u) {
		return false
	}

	if isPathEqual(visited, u) {
		return false
	}

	return true
}

func getDomain(urls map[url.URL]bool) *string {
	for k := range urls {
		hn := k.Hostname()
		return &hn
	}

	return nil
}

func isNothingVisited(u *string) bool {
	return u == nil
}

func isDomainEqual(visited, found string) bool {
	return visited == found
}

func isVisited(visited map[url.URL]bool, found *url.URL) bool {
	if _, ok := visited[*found]; ok {
		return true
	}

	return false
}

func isPathEqual(visited map[url.URL]bool, found *url.URL) bool {
	for k := range visited {
		if strings.TrimSuffix(k.Path, "/") == strings.TrimSuffix(found.Path, "/") {
			return true
		}
	}

	return false
}
