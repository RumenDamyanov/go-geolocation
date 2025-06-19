# go-geolocation

[![CI](https://github.com/rumendamyanov/go-geolocation/actions/workflows/ci.yml/badge.svg)](https://github.com/rumendamyanov/go-geolocation/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/rumendamyanov/go-geolocation/branch/main/graph/badge.svg)](https://codecov.io/gh/rumendamyanov/go-geolocation)

A framework-agnostic Go module for geolocation, inspired by php-geolocation. Provides core geolocation features and adapters for popular Go web frameworks.

## Features

- Extracts country from Cloudflare headers
- Parses browser, OS, device, and language from standard headers
- Middleware/adapters for net/http, Gin, Echo, Fiber
- Testable, modular design
- High test coverage and CI integration

## Installation

```sh
go get github.com/rumendamyanov/go-geolocation
```

## Usage

### Core Usage

```go
import "github.com/rumendamyanov/go-geolocation"

// In your handler:
loc := geolocation.FromRequest(r)
client := geolocation.ParseClientInfo(r)
languages := geolocation.ParseLanguageInfo(r)
```

### net/http Example

See: [examples/nethttp.go](examples/nethttp.go)

### Gin Example

See: [examples/gin.go](examples/gin.go)

### Echo Example

See: [examples/echo.go](examples/echo.go)

### Fiber Example

See: [examples/fiber.go](examples/fiber.go)

## Example Output

```
IP: 1.2.3.4
Country: BG
Browser: Chrome 123.0.0.0
OS: Windows NT 10.0
Device: Desktop
DefaultLang: en-US
AllLangs: [en-US en bg de]
```

## Testing

Run all tests:

```sh
go test ./...
```

Check coverage:

```sh
go test -cover ./...
```

> **Note on Coverage:**
>
> All error branches and edge cases in the core package are thoroughly tested. Due to Go's coverage tool behavior, a few lines in `LoadConfig` may not be counted as covered, even though all real error paths (file not found, invalid JSON/YAML, unsupported extension) are exercised in tests. The code is idiomatic and robust; further refactoring for the sake of 100% coverage is not recommended.

## CI & Coverage

- GitHub Actions for CI
- Codecov integration for test coverage

## License

[MIT](LICENSE.md)

## Advanced Usage

### Combining Geolocation, Client Info, and Language

You can extract all available information in a single handler:

```go
import (
    "fmt"
    "net/http"
    "github.com/rumendamyanov/go-geolocation"
)

func handler(w http.ResponseWriter, r *http.Request) {
    loc := geolocation.FromRequest(r)
    info := geolocation.ParseClientInfo(r)
    lang := geolocation.ParseLanguageInfo(r)
    fmt.Fprintf(w, "IP: %s\nCountry: %s\nBrowser: %s %s\nOS: %s\nDevice: %s\nDefaultLang: %s\nAllLangs: %v\n",
        loc.IP, loc.Country, info.BrowserName, info.BrowserVersion, info.OS, info.Device, lang.Default, lang.Supported)
}
```

### Custom Middleware Example (net/http)

You can create your own middleware to attach all info to the request context:

```go
import (
    "context"
    "net/http"
    "github.com/rumendamyanov/go-geolocation"
)

type contextKey struct{}

func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        loc := geolocation.FromRequest(r)
        info := geolocation.ParseClientInfo(r)
        lang := geolocation.ParseLanguageInfo(r)
        ctx := context.WithValue(r.Context(), contextKey{}, struct {
            *geolocation.Location
            *geolocation.ClientInfo
            *geolocation.LanguageInfo
        }{loc, info, lang})
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Using in API Responses

You can return all extracted info as JSON in your API endpoints for debugging or analytics:

```go
import (
    "encoding/json"
    "net/http"
    "github.com/rumendamyanov/go-geolocation"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
    resp := struct {
        *geolocation.Location
        *geolocation.ClientInfo
        *geolocation.LanguageInfo
    }{
        geolocation.FromRequest(r),
        geolocation.ParseClientInfo(r),
        geolocation.ParseLanguageInfo(r),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
```
