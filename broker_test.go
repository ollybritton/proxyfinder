package proxyfinder_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ollybritton/proxyfinder"
)

var realProviders map[string]func() ([]proxyfinder.Proxy, error)

// Here we change the proxyfinder.Providers variable so that methods such as `.Load()` load proxies from our dummy provider.
// The problem is, if a proxy provider goes down (as it did for me), then the tests will fail, which is a problem since the issue isn't with the package, but with the provider.
func init() {
	realProviders = make(map[string]func() ([]proxyfinder.Proxy, error))

	for key := range proxyfinder.Providers {
		realProviders[key] = proxyfinder.Providers[key]
		delete(proxyfinder.Providers, key)
	}

	proxyfinder.Providers["dummy"] = DummyProvider
}

func DummyProvider() ([]proxyfinder.Proxy, error) {
	var proxies []proxyfinder.Proxy

	proxyURLs := []string{"http://117.191.11.105:80",
		"http://139.255.123.194:4550",
		"http://41.216.230.154:48705",
		"http://43.228.221.171:61267",
		"http://158.58.130.222:53281",
		"http://159.138.22.112:80",
		"http://118.97.180.132:30793",
		"http://176.35.51.2:53281",
		"http://103.250.152.20:61354",
		"http://90.182.159.82:53803",
	}

	for _, proxyURL := range proxyURLs {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			return []proxyfinder.Proxy{}, err
		}

		proxy := proxyfinder.NewProxy(*parsed, "dummy")
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

// TestLoad checks that the .Load method of a broker works correctly.
// Since `dummyProvider` returns 10 proxies, we should expect 10 proxies back.
func TestLoad(t *testing.T) {
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	assert.Equal(t, len(proxies.All()), 10)
}

// TestFetch checks the .Fetch and by extension .FetchProvider methods are working correctly.
// If we fetch identical sets of proxies twice, we shouldn't have duplicates.
func TestFetch(t *testing.T) {
	proxies := proxyfinder.NewBroker()

	proxies.Fetch()
	initialLength := len(proxies.All())

	proxies.Fetch()
	newLength := len(proxies.All())

	assert.Equal(t, initialLength, newLength)
}

// TestPurge checks that the .Purge and .PurgeProvider method works correctly.
// If we fetch proxies, and then purge that broker, we should be left with 0 proxies.
func TestPurge(t *testing.T) {
	proxies1 := proxyfinder.NewBroker()
	proxies1.Fetch()
	proxies1.Purge()
	assert.Equal(t, len(proxies1.All()), 0)

	proxies2 := proxyfinder.NewBroker()
	proxies2.FetchProvider("dummy")
	proxies2.PurgeProvider("dummy")
	assert.Equal(t, len(proxies2.All()), 0)
}

// BenchmarkLoad tests how fast proxies can be loaded.
// This is kind of useless, as this is completely dependent on internet speed.
func BenchmarkLoad(b *testing.B) {
	proxies := proxyfinder.NewBroker()

	for n := 0; n < b.N; n++ {
		proxies.Load()
	}
}

// BenchmarkRandom tests how fast a random proxy can be selected.
func BenchmarkRandom(b *testing.B) {
	proxies := proxyfinder.NewBroker()
	proxies.Load()

	for n := 0; n < b.N; n++ {
		proxies.Random()
	}
}
