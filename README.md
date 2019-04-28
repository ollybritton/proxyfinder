# proxyfinder
`proxyfinder` is a simple Go package for finding proxies. It works by congregating proxies from several different providers on the internet which can be accessed through a simple API.

- [proxyfinder](#proxyfinder)
  - [Usage](#usage)
    - [Basic usage](#basic-usage)
    - [Loading proxies](#loading-proxies)
    - [Getting proxies](#getting-proxies)
  - [Providers](#providers)

## Usage
### Basic usage
If you just want to get proxies and use them, use this:
```go
proxyfinder.Load() // Loads proxies into the program.

proxy1 := proxyfinder.New() // Gets a new proxy, one that has not been used.
proxy2 := proxyfinder.Random() // Gets a random proxy, might have been used previously.

fmt.Println(proxy1.URL()) // Get url.URL representation of url.
fmt.Println(proxy1.Provider()) // Get the proxy provider for that url.
```

### Loading proxies
In order to load proxies into the program, two functions are avaliable: `Load` and `LoadProvider`.

`Load` will load all proxies from all the providers. It is used like so:
```go
proxyfinder.Load() // Loads proxies from all providers into the program.
```

`LoadProvider` will only load proxies from one provider into the program. A list of providers is avaliable [here](#providers)

```go
proxyfinder.LoadProvider("freeproxylists")
```

Loading will delete existing proxies from that provider before fetching new ones. If you do not want this, use the functions `Fetch` and `FetchProvider`.

```go
proxyfinder.Fetch() // Fetch all proxies
proxyfinder.FetchProvider("freeproxylists") // Fetch proxies from http://www.freeproxylists.com/
```

Please note `Fetch` respects duplicates, so it will not add multiple proxies with the same URL. If you want to delete proxies, then use `Purge` and `PurgeProvider`.

```go
proxyfinder.Fetch() // Fetch all proxies

proxyfinder.PurgeProvider("foo") // Delete all proxies from the provider 'foo' (non-existent)
proxyfinder.Purge() // Delete all proxies from all providers.
```

### Getting proxies
To get all the proxies as a list:

```go
proxyfinder.All()
```

To get all unused proxies as a list:

```go
proxyfinder.Unused()
```

To get a single, unused proxy:

```go
proxyfinder.New()
```

To get a random proxy that might of been used:

```go
proxyfinder.Random()
```

Proxies are represented by the `Proxy` struct.

```go
proxyfinder.Load() // Load all proxies from all providers.

proxy1 := proxyfinder.New() // Get a new proxy.
proxy1.URL() // Get proxy1 as a url.URL type.
proxy1.URL().String() // Get the string representation of the proxy.

proxy1.Provider() // Get the proxy's provider.
```

## Providers
| URL                                                              | Name             |
| ---------------------------------------------------------------- | ---------------- |
| [http://www.freeproxylists.com/](http://www.freeproxylists.com/) | `freeproxylists` |
