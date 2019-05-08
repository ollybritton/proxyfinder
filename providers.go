package proxyfinder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// Providers maps the colloqial names of proxy providers to the function that returns their proxies.
var Providers = map[string]func() []Proxy{
	"freeproxylists":    FreeProxyLists,
	"spysme":            SpysMe,
	"proxylistdownload": ProxyListDownload,
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

// SpysMe returns all the HTTP proxies that it can find on the website http://spys.me/proxy.txt
func SpysMe() (proxies []Proxy) {
	resp, err := http.Get("http://spys.me/proxy.txt")
	if err != nil {
		return []Proxy{}
	}
	defer resp.Body.Close()

	// Response is 4 lines of metadata, follwed by a list of proxies, followed by 2 lines of more metadata.
	// Each proxy is of the form:
	// <IP>:<PORT> <INFO>
	// Therefore, if we take the first item from the string splitted on a space, it will be the IP & PORT.

	body, err := ioutil.ReadAll(resp.Body)

	lines := strings.Split(string(body), "\n")
	lines = lines[4 : len(lines)-2]

	for _, line := range lines {
		rawURL := "http://" + strings.Split(line, " ")[0]
		proxyURL, err := url.Parse(rawURL)

		if err != nil {
			continue
		}

		proxies = append(proxies, NewProxy(*proxyURL, "spysme"))
	}

	return proxies

}

// ProxyListDownload returns all the HTTP proxies that it can find on the website https://www.proxy-list.download/HTTP.
func ProxyListDownload() (proxies []Proxy) {
	fmt.Println("Yoyoyo")
	response, err := http.Get("https://www.proxy-list.download/api/v0/get?l=en&t=http")
	fmt.Println("Yoyoyo")

	if err != nil {
		return []Proxy{}
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []Proxy{}
	}

	type ProxyResponse struct {
		proxies []map[string]interface{} `json:"LISTA"`
	}

	var result []ProxyResponse

	json.Unmarshal(contents, &result)

	fmt.Println(result)

	return

}
