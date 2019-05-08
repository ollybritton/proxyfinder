package main

import (
	"fmt"

	"gitlab.com/ollybritton/proxyfinder"
)

func main() {
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	locales := []string{"US", "CA", "MX", "AR", "DE", "BE", "NL", "IE", "ES", "IT", "CH"}

	fmt.Println("Proxies before filter:", len(proxies.All()))

	// You need to change the location of the database. (Unless you have the same name as me and keep all your files in the same place)
	proxies.MarkCountries("/Users/olly/Library/_geoip.mmdb")
	proxies.OnlyCountries(locales)

	fmt.Println("Proxies after filter:", len(proxies.All()))
}
