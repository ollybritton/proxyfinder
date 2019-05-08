package main

import (
	"fmt"

	"gitlab.com/ollybritton/proxyfinder"
)

func main() {
	proxies := proxyfinder.NewBroker()
	fmt.Println("Downloading proxies...")
	proxies.LoadProvider("proxylistdownload")

	for i := 0; i < 10; i++ {
		proxy := proxies.New()
		fmt.Println(proxy.URL.String())
	}
}
