package main

import (
	"fmt"
	"time"

	"gitlab.com/ollybritton/proxyfinder"
)

func main() {
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	proxyfinder.TimeoutLength = 5 * time.Second

	fmt.Printf("Unfiltered Proxy Length: %d\n", len(proxies.All()))
	proxies.Filter(proxyfinder.CheckConnection)

	fmt.Printf("New Proxy Length: %d\n", len(proxies.All()))
}
