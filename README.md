![logo](https://github.com/ollybritton/Assets/raw/master/proxyfinder.jpg)
[![GoDoc](https://godoc.org/gitlab.com/ollybritton/proxyfinder?status.svg)](https://godoc.org/gitlab.com/ollybritton/proxyfinder)

`proxyfinder` is a simple Go package for anonymous open proxy servers. It works by congregating proxies from several different providers on the internet which can be accessed through a simple API. On my machine, it can fetch ~2500 proxies in roughly 2.6 seconds.

- [Usage](#usage)
  - [Loading proxies](#loading-proxies)
  - [Accessing proxies](#accessing-proxies)
  - [Filtering proxies](#filtering-proxies)
  - [The `Proxy` type](#the-proxy-type)
- [More examples](#more-examples)
- [Providers](#providers)

## Usage
Proxies can be accessed through a `Broker`.

```go
import (
  "gitlab.com/ollybritton/proxyfinder"
  "fmt"
)
proxies := proxyfinder.NewBroker()
proxies.Load()
fmt.Println(proxies.New())
```

### Loading proxies
All examples assume you have defined a broker named `proxies`, through a statement such as `proxies := proxyfinder.NewBroker()`.

```go
proxies.Load() // Load all proxies - this means deleting all old proxies and getting the new ones
proxies.LoadProvider("freeproxylists") // Load only proxies from the `freeproxylists` provider.

proxies.Fetch() // Fetch all proxies - get all proxies, but don't delete the old ones. This function respects duplicates.
proxies.FetchProvider("freeproxylists") // Fetch only proxies from the `freeproxylists` provider.

proxies.Purge() // Delete all proxies.
proxies.PurgeProvider("freeproxylists") // Delete only proxies from the `freeproxylists` provider.
```

`.Load()` is equivalent to a call to `.Purge()`, followed by a `.Fetch()`.

### Accessing proxies
After proxies have been loaded, you can use the following to access proxies.
```go
proxies.All() // Get a list of all proxies.
proxies.Unused() // Get a list of all unused proxies.

proxies.Random() // Pick a random proxy. This proxy may have been used already.
proxies.New() // Pick a proxy that has not been used yet. 
```

### Filtering proxies
You can also filter proxies using the `Filter` function.

```go
proxies.Load()

proxies.Filter(proxyfinder.CheckConnection) // Finds all proxies that don't timeout (default 10 seconds).

proxyfinder.TimeoutLength = 5 * time.Second
proxies.Filter(proxyfinder.CheckConnection) // Finds all proxies that don't timeout after 5 seconds.

proxies.OnlyCountries([]string{"GB"}) // Remove all proxies that don't come from the UK.

proxies.Filter(func(p proxyfinder.Proxy) bool {
  // Your custom filter here...
})
```

### The `Proxy` type
Functions return a `Proxy` type. `Proxy` is just a wrapper around `url.URL`, additionally containing information about the provider.

For example:

```go
proxies.Load()
proxy := proxies.New()

proxy.URL // Access the url.URL representation of the URL.
proxy.Provider // Get the provider of the proxy.
```

## More examples
More examples can be found in the `examples` folder.  

## Providers
| URL                                                              | Name             |
| ---------------------------------------------------------------- | ---------------- |
| [http://www.freeproxylists.com/](http://www.freeproxylists.com/) | `freeproxylists` |
| [http://spys.me/proxy.txt](http://spys.me/proxy.txt)             | `spysme`         |

