package proxyfinder_test

import (
	"fmt"
	"time"

	"gitlab.com/ollybritton/proxyfinder"
)

func ExampleBroker_Filter() {
	// Create a new proxy broker.
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	// Filter all proxies by checking their connections:
	proxies.Filter(proxyfinder.CheckConnection)
}

func ExampleCheckConnection() {
	// Create new proxy broker.
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	// Get a new proxy and check if it is valid.
	proxy := proxies.New()
	fmt.Printf("Checking proxy %q:\n", proxy.URL.String())
	fmt.Printf("Valid: %t\n", proxyfinder.CheckConnection(proxy))
}

func ExampleCheckConnection_second() {
	// Create a new proxy broker.
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	// Filter all proxies by checking their connections:
	proxies.Filter(proxyfinder.CheckConnection)
}

func ExampleCheckConnection_third() {
	// Create new proxy broker.
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	// Change the timeout to 2 seconds.
	proxyfinder.TimeoutLength = 2 * time.Second

	// Get a new proxy and check if it is valid.
	proxy := proxies.New()
	fmt.Printf("Checking proxy %q:\n", proxy.URL.String())
	fmt.Printf("Valid: %t\n", proxyfinder.CheckConnection(proxy))
}
