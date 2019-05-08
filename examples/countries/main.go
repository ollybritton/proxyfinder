package main

import (
	"fmt"

	"gitlab.com/ollybritton/proxyfinder"
)

func main() {
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	// Two-digit iso codes for that country. Other names will not work.
	locales := []string{"US", "CA", "MX", "AR", "DE", "BE", "NL", "IE", "ES", "IT", "CH"}

	fmt.Println("Proxies before filter:", len(proxies.All()))

	// Some proxies (such as the ones from the `proxylistdownload` provider), come bundled with country information already.
	// For this reason, some proxies which do come from the locale you specified will be discarded as they do not yet have locale information.
	proxies.FilterCountries(locales)

	// In order to add locale information to all proxies, you can use the following:
	// proxies.MarkCountries("path/to/maxmind/geoip/database")

	fmt.Println("Proxies after filter:", len(proxies.All()))
}
