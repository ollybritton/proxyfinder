package proxyfinder

import (
	"math/rand"
	"sync"
)

var p Proxyfinder

// Proxyfinder is a simple Go package for finding proxies.
type Proxyfinder struct {
	// A list of all proxies that have been found.
	proxies []Proxy

	// All unused proxies.
	unused []Proxy
}

/* Loading & Managing proxies */

// Add adds a list of proxies.
func Add(proxies []Proxy) { p.Add(proxies) }
func (p *Proxyfinder) Add(proxies []Proxy) {
	p.proxies = append(p.proxies, proxies...)
	p.unused = append(p.unused, proxies...)
}

// FetchProvider fetches all the proxies from a given provider.
func FetchProvider(provider string) { p.FetchProvider(provider) }
func (p *Proxyfinder) FetchProvider(provider string) {
	proxies := Providers[provider]()
	newUniques := UniqueProxies(append(p.proxies, proxies...))

	p.Add(UniqueProxies(newUniques))
}

// Fetch fetches all the proxies from all providers.
func Fetch() { p.Fetch() }
func (p *Proxyfinder) Fetch() {
	var wg sync.WaitGroup

	wg.Add(len(Providers))

	for provider := range Providers {
		defer wg.Done()
		p.FetchProvider(provider)
	}

	wg.Wait()
}

// PurgeProvider deletes all proxies from a given provider.
func PurgeProvider(provider string) { p.PurgeProvider(provider) }
func (p *Proxyfinder) PurgeProvider(provider string) {
	var newproxies []Proxy
	var newunused []Proxy

	for _, proxy := range p.proxies {
		if proxy.provider != provider {
			newproxies = append(newproxies, proxy)
		}
	}

	for _, proxy := range p.unused {
		if proxy.provider != provider {
			newunused = append(newproxies, proxy)
		}
	}

	p.proxies = newproxies
	p.unused = newunused
}

// Purge will delete all proxies from all providers.
func Purge() { p.Purge() }
func (p *Proxyfinder) Purge() {
	// Since we are deleting everything, then we just reset the proxyfinder.
	p = new(Proxyfinder)
}

// LoadProvider will load all the proxies from a given provider.
func LoadProvider(provider string) { p.LoadProvider(provider) }
func (p *Proxyfinder) LoadProvider(provider string) {
	p.PurgeProvider(provider)
	p.FetchProvider(provider)
}

// Load will load all proxies from all providers.
func Load() { p.Load() }
func (p *Proxyfinder) Load() {
	for provider := range Providers {
		p.LoadProvider(provider)
	}
}

/* Getting/Using Proxies */

// Use marks a proxy as used.
func Use(proxy Proxy) { p.Use(proxy) }
func (p *Proxyfinder) Use(proxy Proxy) {
	var newunused []Proxy

	for _, possibleProxy := range p.unused {
		if possibleProxy != proxy {
			newunused = append(newunused, possibleProxy)
		}
	}

	p.unused = newunused

}

// All returns a list of all proxies.
func All() []Proxy { return p.All() }
func (p *Proxyfinder) All() []Proxy {
	return p.proxies
}

// Unused returns a list of all unused proxies.
func Unused() []Proxy { return p.Unused() }
func (p *Proxyfinder) Unused() []Proxy {
	return p.unused
}

// Random fetches a random proxy.
func Random() Proxy { return p.Random() }
func (p *Proxyfinder) Random() Proxy {
	chosen := p.proxies[rand.Intn(len(p.proxies)-1)]
	p.Use(chosen)

	return chosen
}

// New returns a new, unused proxy.
func New() Proxy { return p.New() }
func (p *Proxyfinder) New() Proxy {
	switch len(p.unused) {
	case 0:
		panic("No proxies unused proxies left.")
	case 1:
		chosen := p.unused[0]
		p.Use(chosen)
		return chosen
	default:
		chosen := p.unused[rand.Intn(len(p.unused)-1)]
		p.Use(chosen)
		return chosen
	}
}

func init() {
	p = Proxyfinder{}
}
