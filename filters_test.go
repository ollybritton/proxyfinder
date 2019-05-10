package proxyfinder_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ollybritton/proxyfinder"
)

func TestBasicFilter(t *testing.T) {
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	proxies.Filter(func(p proxyfinder.Proxy) bool {
		return p.URL.String() == "http://41.216.230.154:48705"
	})

	assert.Equal(t, len(proxies.All()), 1)
}

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
