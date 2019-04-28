package proxyfinder

import (
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// Providers maps the colloqial names of proxy providers to the function that returns their proxies.
var Providers = map[string]func() []Proxy{
	"freeproxylists": FreeProxyLists,
}

// FreeProxyLists returns all the HTTP proxies that it can find on the http://www.freeproxylists.com/ website.
func FreeProxyLists() (proxies []Proxy) {

	initialLinks := FindLinks("http://www.freeproxylists.com/elite.html", `^elite #\d+`)

	links := []string{}

	// Vomit, proxy regex
	proxyRegex := regexp.MustCompile(`&lt;td&gt;\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}&lt;\/td&gt;&lt;td&gt;\d+&lt;\/td&gt;`)
	// End Vomit

	for _, link := range initialLinks {
		parsedLink := "http://freeproxylists.com/load_elite_" + link[6:]
		links = append(links, parsedLink)
	}

	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()

			doc, err := GetURL(l)
			if err != nil {
				return
			}

			matches := proxyRegex.FindAllString(doc, -1)

			for i := range matches {
				matches[i] = strings.Replace(matches[i], "&lt;td&gt;", "", -1)
				matches[i] = strings.Replace(matches[i], "&lt;/td&gt;", ":", 1)
				matches[i] = strings.Replace(matches[i], "&lt;/td&gt;", "", 1)

				parsedMatch, err := url.Parse("http://" + matches[i])
				if err != nil {
					log.Printf("error parsing url %q: %v", matches[i], err.Error())
					continue
				}

				proxy := NewProxy(*parsedMatch, "freeproxylists")

				proxies = append(proxies, proxy)
			}
		}(link)
	}

	wg.Wait()

	return proxies
}
