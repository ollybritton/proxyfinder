# proxyfinder
`proxyfinder` is a simple Go package for finding proxies. It works by congregating proxies from several different providers on the internet which can be accessed through a simple API. On my machine, it can fetch ~2500 proxies in roughly 2.6 seconds.

- [proxyfinder](#proxyfinder)
  - [Usage](#usage)
    - [Loading proxies](#loading-proxies)
    - [Accessing proxies](#accessing-proxies)
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
