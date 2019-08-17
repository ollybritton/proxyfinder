package proxyfinder_test

import (
	"log"
	"testing"

	"gitlab.com/ollybritton/proxyfinder"
)

// This package aims at testing that the providers for library still work, and does
// not actually test the functionality of the library.

// For example, if proxy-list.download goes down, this test suite will find that out.

func ProviderSuite(t *testing.T, name string, provider func() ([]proxyfinder.Proxy, error)) {
	log.Printf("testing provider %v", name)
	proxies, err := provider()

	if err != nil {
		t.Errorf("unable to reach %v", name)
	}

	if len(proxies) == 0 {
		t.Errorf("provider %v returned no proxies but did not error", name)
	}

	log.Printf("provider %v returned %d proxies", name, len(proxies))
}

func TestProviders(t *testing.T) {
	for name, providerFunc := range realProviders {
		ProviderSuite(t, name, providerFunc)
	}
}
