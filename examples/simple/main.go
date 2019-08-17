package main

import (
	"fmt"

	"gitlab.com/ollybritton/proxyfinder"
)

func main() {
	// Create a new proxy broker
	proxies := proxyfinder.NewBroker()

	// Load proxies from all providers
	fmt.Println("Downloading proxies...")
	proxies.Load()

	for i := 0; i < 10; i++ {
		// Print the urls of 10 proxies.
		proxy, err := proxies.New()
		if err != nil {
			panic(err)
		}

		fmt.Println(proxy.URL.String())
	}
}
