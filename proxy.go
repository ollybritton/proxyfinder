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
// Instead of just using the plain Proxy type, it only checks if the URL matches.
// This is because a proxy that is identical to another from a different provider will not match because the provider will be different, although fundementally it is the same proxy.
func UniqueProxies(proxies []Proxy) []Proxy {
	keys := make(map[url.URL]bool)
	list := []Proxy{}

	for _, entry := range proxies {
		if _, value := keys[entry.URL]; !value {
			keys[entry.URL] = true
			list = append(list, entry)
		}
	}
	return list
}
