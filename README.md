# go-geolocation

[![CI](https://github.com/rumendamyanov/go-geolocation/actions/workflows/ci.yml/badge.svg)](https://github.com/rumendamyanov/go-geolocation/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/rumendamyanov/go-geolocation/branch/master/graph/badge.svg)](https://codecov.io/gh/rumendamyanov/go-geolocation)

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
