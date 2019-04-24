# proxyfinder
proxyfinder is a simple Go package for scraping proxies off the web. It is at its early stages, so it does not have much functionality yet.

## Usage
The only really useful function at the moment is `Proxies()`, which will return all the proxies it can find. This is rather slow as it has to find all the proxies first.

```go
func main() {
	fmt.Println(proxyfinder.Proxies()) // Has type []url.URL
}
```