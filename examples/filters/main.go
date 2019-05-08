package main

import (
	"fmt"

	"gitlab.com/ollybritton/proxyfinder"
)

func main() {
	// Create a new proxy broker and load all proxies from all providers into it.
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	// Print the length of all the countries.
	fmt.Printf("Length of all proxies: %d\n", len(proxies.All()))

	// Check the connection on all the proxies, defaults to a timeout of 10 seconds.
	proxies.Filter(proxyfinder.CheckConnection)

	// Print the new length of proxies.
	fmt.Printf("New Proxy Length: %d\n", len(proxies.All()))

	// Example Output:
	// Length of all proxies: 13228
	// New Proxy Length: 9255
}
