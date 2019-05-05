package proxyfinder

import (
	"net/http"
	"sync"
	"time"
)

var (
	// TimeoutLength configures how long to wait when checking a proxy's connection until it is considered not working.
	TimeoutLength = 10 * time.Second
)

// Filter is a method that allows a user to filter their list of proxies.
// You pass a filterFunc which, given a proxy, will return true or false depending if the proxy should stay or not.
func (b *Broker) Filter(filterFunc func(Proxy) bool) {
	var newproxies []Proxy

	var wg sync.WaitGroup

	for _, proxy := range b.proxies {
		go func(p Proxy) {
			wg.Add(1)
			defer wg.Done()

			if filterFunc(p) == true {
				newproxies = append(newproxies, p)
			}
		}(proxy)
	}

	wg.Wait()

	b.proxies = newproxies
}

// CheckConnection checks that a proxy is working.
// Because all these proxies are just random ones found on the internet, they are quite intermittent or slow.
// This function (by default) will mark a proxy as not working if there is no response after 10 seconds. This can be changed by redefining proxyfinder.TimeoutLength, for example:
//  proxyfinder.TimeoutLength = 5 * time.Seconds
func CheckConnection(proxy Proxy) bool {
	timeout := time.Duration(TimeoutLength)

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(&proxy.URL)},
		Timeout:   timeout,
	}

	// TODO: there's probably a much better way to do this.
	req, _ := http.NewRequest("GET", "https://ismyinternetworking.com/", nil)
	_, err := client.Do(req)

	if err != nil {
		return false
	}

	return true
}
