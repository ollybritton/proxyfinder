package proxyfinder

import (
	"net/http"
	"net/url"
)

// proxyRequest returns a HTTP client and a request.
func proxyClient(proxy url.URL) (*http.Client, *http.Request) {
	client := &http.Client{}
	client.Transport = &http.Transport{Proxy: http.ProxyURL(&proxy)}

	req, _ := http.NewRequest("GET", proxy.String(), nil)

	return client, req
}

func TestProxy() {

}
