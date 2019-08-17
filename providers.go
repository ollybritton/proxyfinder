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

	"github.com/imroc/req"
)

// Providers maps the colloqial names of proxy providers to the function that returns their proxies.
var Providers = map[string]func() ([]Proxy, error){
	"freeproxylists":    FreeProxyLists,
	"spysme":            SpysMe,
	"proxylistdownload": ProxyListDownload,
	"proxyscrape":       ProxyScrape,
	"static":            Static,
}

// FreeProxyLists returns all the HTTP proxies that it can find on the http://www.freeproxylists.com/ website.
func FreeProxyLists() ([]Proxy, error) {
	var proxies []Proxy

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
	var proxyFetchError error

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()

			doc, err := GetURL(l)
			if err != nil {
				proxyFetchError = fmt.Errorf("unable to request freeproxylists: %v", err)
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
	if proxyFetchError != nil {
		return []Proxy{}, proxyFetchError
	}

	return proxies, nil
}

// SpysMe returns all the HTTP proxies that it can find on the website http://spys.me/proxy.txt
func SpysMe() ([]Proxy, error) {
	var proxies []Proxy

	resp, err := http.Get("http://spys.me/proxy.txt")
	if err != nil {
		return []Proxy{}, fmt.Errorf("unable to request spysme proxies: %v", err)
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

	return proxies, nil

}

// ProxyListDownload returns all the HTTP proxies that it can find on the website https://www.proxy-list.download/HTTP.
func ProxyListDownload() ([]Proxy, error) {
	var proxies []Proxy

	response, err := http.Get("https://www.proxy-list.download/api/v0/get?l=en&t=http")

	if err != nil {
		return []Proxy{}, fmt.Errorf("unable to request proxylistdownload: %v", err)
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []Proxy{}, fmt.Errorf("unable to read proxylistdownload response: %v", err)
	}

	type ProxyJSON struct {
		IP   string `json:"IP"`
		Port string `json:"PORT"`
		ISO  string `json:"ISO"`
	}

	type ProxyList struct {
		Proxies []ProxyJSON `json:"LISTA"`
	}

	var result []ProxyList

	err = json.Unmarshal(contents, &result)
	if err != nil {
		return []Proxy{}, fmt.Errorf("error unmarshling response from proxylist-download site: %v", err)
	}

	if len(result) == 0 {
		return []Proxy{}, fmt.Errorf("empty response from proxylist-download: %v", err)
	}

	for _, proxy := range result[0].Proxies {
		ip, err := url.Parse("http://" + proxy.IP + ":" + proxy.Port)

		if err != nil {
			return []Proxy{}, fmt.Errorf("unable to parse proxy url: %v", err)
		}

		parsedProxy := NewProxy(*ip, "proxylistdownload")
		parsedProxy.Country = proxy.ISO

		proxies = append(proxies, parsedProxy)
	}

	return proxies, nil

}

// ProxyScrape will get proxies using the service http://proxyscrape.com.
func ProxyScrape() ([]Proxy, error) {
	var formatURL = "https://api.proxyscrape.com/?request=getproxies&proxytype=http&timeout=10000&country=%v&ssl=all&anonymity=all"
	var countries = []string{"AF", "AL", "AM", "AR", "AT", "AU", "BA", "BD", "BG", "BO", "BR", "BY", "CA", "CL", "CM", "CN", "CO", "CZ", "DE", "EC", "EG", "ES", "FR", "GB", "GE", "GN", "GR", "GT", "HK", "HN", "HU", "ID", "IN", "IQ", "IR", "IT", "JP", "KE", "KG", "KH", "KR", "KZ", "LB", "LT", "LV", "LY", "MD", "MM", "MN", "MU", "MW", "MX", "MY", "NG", "NL", "NO", "NP", "PE", "PH", "PK", "PL", "PS", "PY", "RO", "RS", "RU", "SC", "SE", "SG", "SK", "SY", "TH", "TR", "TW", "TZ", "UA", "UG", "US", "VE", "VN", "ZA"}

	var proxies []Proxy
	var wg sync.WaitGroup

	for _, country := range countries {
		go func(country string) {
			wg.Add(1)
			defer wg.Done()

			countryURL := fmt.Sprintf(formatURL, country)
			r, err := req.Get(countryURL)

			if err != nil {
				return
			}

			list, err := r.ToString()
			if err != nil {
				return
			}

			proxyList := strings.Split(list, "\n")

			for _, stringProxy := range proxyList {
				parsedURL, err := url.Parse(strings.TrimSpace("http://" + stringProxy))
				if err != nil {
					return
				}

				proxy := NewProxy(*parsedURL, "proxyscrape")
				proxy.Country = country

				proxies = append(proxies, proxy)
			}
		}(country)

	}

	wg.Wait()

	if len(proxies) == 0 {
		return []Proxy{}, fmt.Errorf("no proxies gathered from proxyscrape provider")
	}

	return proxies, nil
}

// Static gets a saved offline static proxy list. Although this can be accessed when offline, it means
// that the proxies can go offline while still appearing in this list.
func Static() ([]Proxy, error) {
	var proxies []Proxy

	for _, rawURL := range staticProxies {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return []Proxy{}, err
		}

		proxy := NewProxy(*parsed, "static")
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}
