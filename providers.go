package proxyfinder

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

// Provider represents a website's proxies.
type Provider struct {
	name    string    // The url of the proxy website
	proxies []url.URL // The list of proxies found on that website.
}

// FreeProxyListsHTTP returns all the HTTP proxies that it can find on the http://www.freeproxylists.com/ website.
func FreeProxyListsHTTP() (proxies []url.URL) {

	/*
		Either in order to prevent scraping or to make fetching the proxies easier, this website will not load the proxies straight away. Instead, it will make a request to another URL which can be accessed using the "key" found in the original URL.

		For example, if the original URL is "...elite/1556123424.html", the proxies can be found at ".../load_elite_1556123424.html".

		The resource it points to is not actually a HTML resource. It is a XML document. This code does not do any fancy parsing of the XML document, and instead just uses regular expressions to match the proxies, then removing the XML open and close tags in between (I know, I'm very sorry).

	*/

	initialLinks := FindLinks("http://www.freeproxylists.com/elite.html", `^elite #\d+`)

	links := []string{}

	// Vomit, proxy regex
	proxyRegex := regexp.MustCompile(`&lt;td&gt;\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}&lt;\/td&gt;&lt;td&gt;\d+&lt;\/td&gt;`)
	// End Vomit

	for _, link := range initialLinks {
		parsedLink := "http://freeproxylists.com/load_elite_" + link[6:]
		links = append(links, parsedLink)
	}

	for _, link := range links {
		log.Println("Scanning", link, "for proxies...")

		doc, err := GetURL(link)
		if err != nil {
			return []url.URL{}
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

			proxies = append(proxies, *parsedMatch)
		}
	}

	return proxies
}

// Proxies will return a list of all the proxies it can find.
func Proxies() (proxies []url.URL) {
	proxies = append(proxies, FreeProxyListsHTTP()...)
	return proxies
}
