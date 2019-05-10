package proxyfinder_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ollybritton/proxyfinder"
)

func TestUniques(t *testing.T) {
	url1, _ := url.Parse("http://192.168.0.1")
	url2, _ := url.Parse("http://192.168.0.2")
	url3, _ := url.Parse("http://192.168.0.1")

	proxy1 := proxyfinder.NewProxy(*url1, "test1")
	proxy3 := proxyfinder.NewProxy(*url2, "test3")
	proxy2 := proxyfinder.NewProxy(*url3, "test2")

	filtered := proxyfinder.UniqueProxies([]proxyfinder.Proxy{
		proxy1, proxy2, proxy3,
	})

	assert.Equal(t, len(filtered), 2)
}
