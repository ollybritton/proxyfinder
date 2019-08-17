package proxyfinder

import (
	"fmt"
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
func (b *Broker) FetchProvider(provider string) error {
	proxies, err := Providers[provider]()

	if err != nil {
		return err
	}

	b.Add(proxies)

	return nil
}

// Fetch fetches all the proxies from all providers.
func (b *Broker) Fetch() error {
	var wg sync.WaitGroup

	wg.Add(len(Providers))
	var err error

	for provider := range Providers {
		go func(p string) {
			defer wg.Done()
			err = b.FetchProvider(p)
		}(provider)
	}

	wg.Wait()

	if err != nil {
		return err
	}

	return nil
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
func (b *Broker) LoadProvider(provider string) error {
	b.PurgeProvider(provider)
	err := b.FetchProvider(provider)

	return err
}

// Load will load all proxies from all providers.
func (b *Broker) Load() error {
	for provider := range Providers {
		err := b.LoadProvider(provider)

		if err != nil {
			return err
		}
	}

	return nil
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
func (b *Broker) New() (Proxy, error) {
	unused := b.Unused()

	switch len(unused) {
	case 0:
		return Proxy{}, fmt.Errorf("no unused proxies")

	case 1:
		chosen := unused[0]
		b.Use(chosen)
		return chosen, nil

	default:
		chosen := unused[rand.Intn(len(unused)-1)]
		b.Use(chosen)
		return chosen, nil
	}
}
