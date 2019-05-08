package proxyfinder

import "net/url"

// Proxy represents a URL that can be used.
type Proxy struct {
	URL      url.URL
	Used     bool
	Provider string
	Country  string
}

// NewProxy is a convenience function for generating a new proxy struct.
func NewProxy(url url.URL, provider string) Proxy {
	return Proxy{URL: url, Provider: provider, Used: false}
}

// UniqueProxies returns a list of proxies with all duplicates removed.
func UniqueProxies(proxies []Proxy) []Proxy {
	keys := make(map[Proxy]bool)
	list := []Proxy{}
	for _, entry := range proxies {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
