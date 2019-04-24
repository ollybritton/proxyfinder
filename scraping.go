package proxyfinder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gocolly/colly"
)

// FindLinks will return the URLs on page at a given URL if the text matches the regex specified.
func FindLinks(url, regex string) (links []string) {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting URL '%v'\n", r.URL)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		regex := regexp.MustCompile(regex)

		if regex.Match([]byte(e.Text)) {
			log.Println("Found matching URL:", e.Attr("href"))
			links = append(links, e.Attr("href"))
		}

	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)
	return links
}

// GetURL will return the string of the response at the URL specified.
func GetURL(url string) (string, error) {
	var client http.Client

	response, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("invalid status code")
	}

	body, err := ioutil.ReadAll(response.Body)
	doc := string(body)

	return doc, nil
}
