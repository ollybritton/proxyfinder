package proxyfinder

import "net/url"

// Proxy represents a URL that can be used.
type Proxy struct {
	url      url.URL
	provider string
}

// URL returns the url.URL representation of the proxy.
func (p *Proxy) URL() url.URL {
	return p.url
}

// Provider returns the name of the provider.
func (p *Proxy) Provider() string {
	return p.provider
}

// NewProxy is a convenience function for generating a new proxy struct.
func NewProxy(url url.URL, provider string) Proxy {
	return Proxy{url: url, provider: provider}
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
