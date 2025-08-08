# go-geolocation

[![CI](https://github.com/rumendamyanov/go-geolocation/actions/workflows/ci.yml/badge.svg)](https://github.com/rumendamyanov/go-geolocation/actions/workflows/ci.yml)
![CodeQL](https://github.com/rumendamyanov/go-geolocation/actions/workflows/github-code-scanning/codeql/badge.svg)
![Dependabot](https://github.com/rumendamyanov/go-geolocation/actions/workflows/dependabot/dependabot-updates/badge.svg)
[![codecov](https://codecov.io/gh/rumendamyanov/go-geolocation/branch/master/graph/badge.svg)](https://codecov.io/gh/rumendamyanov/go-geolocation)
[![Go Report Card](https://goreportcard.com/badge/go.rumenx.com/geolocation?)](https://goreportcard.com/report/go.rumenx.com/geolocation)
[![Go Reference](https://pkg.go.dev/badge/go.rumenx.com/geolocation.svg)](https://pkg.go.dev/go.rumenx.com/geolocation)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/rumendamyanov/go-geolocation/blob/master/LICENSE.md)

A framework-agnostic Go module for geolocation, inspired by [php-geolocation](https://github.com/RumenDamyanov/php-geolocation). Provides core geolocation features and adapters for popular Go web frameworks.

## Features

- Extracts country from Cloudflare headers
- Parses browser, OS, device, and language from standard headers
- **Local development simulation** - Fake Cloudflare headers for testing without production setup
- **Auto-detection** of local environments (localhost, local IPs, missing Cloudflare headers)
- **Advanced language negotiation** - matches browser and available site languages for multi-language countries
- **Comprehensive client info** - browser, OS, device type (including tablet), screen resolution
- **Built-in country data** for 8 countries (US, CA, GB, DE, FR, JP, AU, BR)
- Middleware/adapters for net/http, Gin, Echo, Fiber
- Testable, modular design
- High test coverage and CI integration

## Installation

```sh
go get go.rumenx.com/geolocation
```

## Usage

### Core Usage

```go
import "go.rumenx.com/geolocation"

// In your handler:
loc := geolocation.FromRequest(r)
client := geolocation.ParseClientInfo(r)
languages := geolocation.ParseLanguageInfo(r)
```

## Framework Adapters

Ready-to-use examples for popular Go web frameworks are available in the [examples/](examples/) directory:

| Framework | Port | Adapter Path | Description |
|-----------|------|--------------|-------------|
| **net/http** | 8080 | [adapters/nethttp](adapters/nethttp) | Standard library HTTP server integration |
| **Gin** | 8081 | [adapters/gin](adapters/gin) | Gin web framework integration |
| **Echo** | 8082 | [adapters/echo](adapters/echo) | Echo web framework integration |
| **Fiber** | 8083 | [adapters/fiber](adapters/fiber) | Fiber web framework integration |

### Two Integration Approaches

1. **Import Adapter Packages** — For clean middleware integration:

```go
go get go.rumenx.com/geolocation/adapters/gin    # or echo, fiber, nethttp
```

1. **Copy Example Applications** — For quick start with full applications:

```bash
# Clone and run complete example servers
git clone https://github.com/RumenDamyanov/go-geolocation.git
cd go-geolocation/examples/gin-adapter && go run main.go
```

### Quick Start with Framework Examples

```bash
# Clone the repository
git clone https://github.com/RumenDamyanov/go-geolocation.git
cd go-geolocation

# Run net/http example (port 8080)
cd examples/nethttp-adapter && go run main.go

# Run Gin example (port 8081)
cd examples/gin-adapter && go run main.go

# Run Echo example (port 8082)
cd examples/echo-adapter && go run main.go

# Run Fiber example (port 8083)
cd examples/fiber-adapter && go run main.go
```

### Test the Examples

```bash
# Download geolocation data
curl "http://localhost:8080/" | jq

# Simulate a specific country
curl "http://localhost:8081/simulate/DE" | jq

# Get available countries
curl "http://localhost:8082/countries" | jq
```

### Integration Pattern Example

Each framework adapter follows the same pattern:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go.rumenx.com/geolocation"
    ginadapter "go.rumenx.com/geolocation/adapters/gin"
)

func main() {
    r := gin.Default()
    r.Use(ginadapter.Middleware())

    r.GET("/location", func(c *gin.Context) {
        loc := ginadapter.FromContext(c)
        clientInfo := geolocation.ParseClientInfo(c.Request)

        c.JSON(200, gin.H{
            "location": loc,
            "client":   clientInfo,
        })
    })

    r.Run(":8080")
}
```

### Gin Example

```go
package main

import (
    "github.com/gin-gonic/gin"
    "go.rumenx.com/geolocation"
    ginadapter "go.rumenx.com/geolocation/adapters/gin"
)

func main() {
    r := gin.Default()
    r.Use(ginadapter.Middleware())

    r.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        loc := ginadapter.FromContext(c)

        c.JSON(200, gin.H{
            "user":     name,
            "location": loc,
            "local":    geolocation.IsLocalDevelopment(c.Request),
        })
    })

    r.Run(":8080")
}
```

### Echo Example

```go
package main

import (
    "github.com/labstack/echo/v4"
    "go.rumenx.com/geolocation"
    echoadapter "go.rumenx.com/geolocation/adapters/echo"
)

func main() {
    e := echo.New()
    e.Use(echoadapter.Middleware())

    e.GET("/user/:name", func(c echo.Context) error {
        name := c.Param("name")
        loc := echoadapter.FromContext(c)

        return c.JSON(200, map[string]interface{}{
            "user":     name,
            "location": loc,
            "local":    geolocation.IsLocalDevelopment(c.Request()),
        })
    })

    e.Start(":8080")
}
```

### Fiber Example

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    fiberadapter "go.rumenx.com/geolocation/adapters/fiber"
)

func main() {
    app := fiber.New()
    app.Use(fiberadapter.Middleware())

    app.Get("/user/:name", func(c *fiber.Ctx) error {
        name := c.Params("name")
        loc := fiberadapter.FromContext(c)

        return c.JSON(fiber.Map{
            "user":     name,
            "location": loc,
            "local":    loc.IP == "",
        })
    })

    app.Listen(":8080")
}
```

### net/http Example

```go
package main

import (
    "encoding/json"
    "net/http"
    "go.rumenx.com/geolocation"
    httpadapter "go.rumenx.com/geolocation/adapters/nethttp"
)

func main() {
    mux := http.NewServeMux()

    mux.Handle("/user", httpadapter.HTTPMiddleware(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            loc := httpadapter.FromContext(r.Context())
            clientInfo := geolocation.ParseClientInfo(r)

            response := map[string]interface{}{
                "location": loc,
                "client":   clientInfo,
                "local":    geolocation.IsLocalDevelopment(r),
            }

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(response)
        })))

    http.ListenAndServe(":8080", mux)
}
```

## Example Output

```text
IP: 1.2.3.4
Country: BG
Browser: Chrome 123.0.0.0
OS: Windows NT 10.0
Device: Desktop
DefaultLang: en-US
AllLangs: [en-US en bg de]
```

## Local Development Simulation

When developing locally where Cloudflare is not available, you can simulate its functionality:

### Quick Simulation

```go
// Create a simulated request for a specific country
req := geolocation.Simulate("DE", nil)
info := geolocation.GetGeoInfo(req)
fmt.Printf("Country: %s, IP: %s\n", info.CountryCode, info.IP)
```

### Advanced Simulation

```go
// Custom simulation options
req := geolocation.Simulate("JP", &geolocation.SimulationOptions{
    UserAgent: "Custom User Agent",
    Languages: []string{"ja", "en"},
})
```

### Auto-Detection of Local Environment

```go
req := getRequest() // your HTTP request
if geolocation.IsLocalDevelopment(req) {
    fmt.Println("Running in local development mode")
}
```

### Available Countries for Simulation

```go
countries := geolocation.GetAvailableCountries()
// Returns: ["US", "CA", "GB", "DE", "FR", "JP", "AU", "BR"]

randomCountry := geolocation.RandomCountry()
```

## Advanced Features

### Comprehensive Client Information

```go
info := geolocation.GetGeoInfo(req)
fmt.Printf("Country: %s\n", info.CountryCode)
fmt.Printf("IP: %s\n", info.IP)
fmt.Printf("Browser: %s %s\n", info.Browser, info.BrowserVersion)
fmt.Printf("OS: %s\n", info.OS)
fmt.Printf("Device: %s\n", info.Device) // Desktop, Mobile, Tablet
fmt.Printf("Languages: %v\n", info.AllLanguages)
fmt.Printf("Resolution: %dx%d\n", info.Resolution.Width, info.Resolution.Height)
```

### Advanced Language Negotiation

```go
cfg := &geolocation.Config{
    DefaultLanguage: "en",
    CountryToLanguageMap: map[string][]string{
        "CA": {"en", "fr"}, // Canada: English (default), French
        "CH": {"de", "fr", "it"}, // Switzerland: German, French, Italian
    },
}

// Sophisticated language selection based on:
// 1. Browser preferred language matching country languages and available site languages
// 2. All browser languages for matches with available languages
// 3. First country language as fallback
availableSiteLanguages := []string{"en", "fr", "de", "es"}
bestLang := geolocation.GetLanguageForCountry(req, cfg, "CH", availableSiteLanguages)

// Check if language cookie should be set
if geolocation.ShouldSetLanguage(req, "lang") {
    // Set language in your application
    geolocation.SetCookie(w, "lang", bestLang, &http.Cookie{MaxAge: 86400 * 30})
}
```

### Screen Resolution Detection

```go
// Frontend JavaScript would set these headers:
// req.Header.Set("X-Screen-Width", "1920")
// req.Header.Set("X-Screen-Height", "1080")

resolution := geolocation.GetResolution(req)
fmt.Printf("Screen: %dx%d\n", resolution.Width, resolution.Height)
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

## Documentation

For comprehensive documentation and additional examples:

- **[Contributing Guidelines](CONTRIBUTING.md)** - How to contribute to this project
- **[Security Policy](SECURITY.md)** - Security guidelines and vulnerability reporting
- **[Code of Conduct](CODE_OF_CONDUCT.md)** - Community guidelines and behavior expectations
- **[Funding](FUNDING.md)** - Support and sponsorship information
- **[Wiki](https://github.com/RumenDamyanov/go-geolocation/wiki)** - Extended examples and tutorials

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on how to submit pull requests, report issues, and contribute to the project.

## Support

If you find this project helpful, please consider:

- ⭐ Starring the repository
- 📝 Reporting issues or suggesting features
- 💝 Supporting via [GitHub Sponsors](FUNDING.md)

For detailed support information, see [FUNDING.md](FUNDING.md).
