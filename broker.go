package proxyfinder

import (
	"math/rand"
	"sync"
)

// Broker is a type which provides access to proxies.
type Broker struct {
	// A list of all proxies that have been found.
	proxies []Proxy
}

// NewBroker is a simple function to create a new broker.
func NewBroker() Broker {
	return Broker{}
}

/* Loading & Managing proxies */

// Add adds a list of proxies.
func (b *Broker) Add(proxies []Proxy) {
	b.proxies = UniqueProxies(append(b.proxies, proxies...))
}

// FetchProvider fetches all the proxies from a given provider.
func (b *Broker) FetchProvider(provider string) {
	proxies := Providers[provider]()

	b.Add(proxies)
}

// Fetch fetches all the proxies from all providers.
func (b *Broker) Fetch() {
	var wg sync.WaitGroup

	wg.Add(len(Providers))

	for provider := range Providers {
		go func(p string) {
			defer wg.Done()
			b.FetchProvider(p)
		}(provider)
	}

	wg.Wait()
}

// PurgeProvider deletes all proxies from a given provider.
func (b *Broker) PurgeProvider(provider string) {
	var newproxies []Proxy

	for _, proxy := range b.proxies {
		if proxy.Provider != provider {
			newproxies = append(newproxies, proxy)
		}
	}

	b.proxies = newproxies
}

// Purge will delete all proxies from all providers.
func (b *Broker) Purge() {
	b.proxies = []Proxy{}
}

// LoadProvider will load all the proxies from a given provider.
func (b *Broker) LoadProvider(provider string) {
	b.PurgeProvider(provider)
	b.FetchProvider(provider)
}

// Load will load all proxies from all providers.
func (b *Broker) Load() {
	for provider := range Providers {
		b.LoadProvider(provider)
	}
}

/* Getting/Using Proxies */

// Use marks a proxy as used.
func (b *Broker) Use(proxy Proxy) {

	for i, possibleProxy := range b.proxies {
		if possibleProxy == proxy {
			b.proxies[i].Used = true
		}
	}
}

// All returns a list of all proxies.
func (b *Broker) All() []Proxy {
	return b.proxies
}

// Unused returns a list of all unused proxies.
func (b *Broker) Unused() []Proxy {
	var unused []Proxy

	for _, proxy := range b.proxies {
		if proxy.Used == false {
			unused = append(unused, proxy)
		}
	}

	return unused
}

// Random fetches a random proxy.
func (b *Broker) Random() Proxy {
	chosen := b.proxies[rand.Intn(len(b.proxies)-1)]
	b.Use(chosen)

	return chosen
}

// New returns a new, unused proxy.
func (b *Broker) New() Proxy {
	unused := b.Unused()

	switch len(unused) {
	case 0:
		panic("No unused proxies left.")
	case 1:
		chosen := unused[0]
		b.Use(chosen)
		return chosen
	default:
		chosen := unused[rand.Intn(len(unused)-1)]
		b.Use(chosen)
		return chosen
	}
}
