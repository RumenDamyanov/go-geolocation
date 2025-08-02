# go-geolocation Examples

This directory contains complete working examples for integrating go-geolocation with popular Go web frameworks.

## Available Examples

| Framework | Port | Directory | Description |
|-----------|------|-----------|-------------|
| **net/http** | 8080 | [nethttp-adapter/](nethttp-adapter/) | Standard library HTTP server |
| **Gin** | 8081 | [gin-adapter/](gin-adapter/) | Gin web framework |
| **Echo** | 8082 | [echo-adapter/](echo-adapter/) | Echo web framework |
| **Fiber** | 8083 | [fiber-adapter/](fiber-adapter/) | Fiber web framework |

## Quick Start

1. **Clone the repository:**
```bash
git clone https://github.com/rumendamyanov/go-geolocation.git
cd go-geolocation
```

2. **Run any example:**
```bash
# net/http example (port 8080)
cd examples/nethttp-adapter && go run main.go

# Gin example (port 8081)
cd examples/gin-adapter && go run main.go

# Echo example (port 8082)
cd examples/echo-adapter && go run main.go

# Fiber example (port 8083)
cd examples/fiber-adapter && go run main.go
```

3. **Test the examples:**
```bash
# Basic geolocation
curl http://localhost:8080/

# Simulate different countries
curl http://localhost:8081/simulate/DE
curl http://localhost:8082/simulate/FR
curl http://localhost:8083/simulate/JP

# Get available countries
curl http://localhost:8080/countries
```

## Common Features

All examples include:

- **Basic geolocation** - Extract location info from Cloudflare headers
- **Local development simulation** - Fake geolocation data for testing
- **Client info parsing** - Browser, OS, device detection
- **Language negotiation** - Parse Accept-Language headers
- **Country simulation** - Test with different countries
- **Health checks** - Basic monitoring endpoints

## Example Response Format

```json
{
  "location": {
    "ip": "192.168.1.123",
    "country": "US"
  },
  "client_info": {
    "browser_name": "Chrome",
    "browser_version": "91.0.4472.124",
    "os": "macOS",
    "device": "Desktop"
  },
  "language": {
    "default": "en-US",
    "supported": ["en-US", "en", "fr"]
  },
  "is_local": true
}
```

## Legacy Examples

For backward compatibility, simplified single-file examples are available:

- [gin.go](gin.go) - Basic Gin example
- [echo.go](echo.go) - Basic Echo example
- [fiber.go](fiber.go) - Basic Fiber example
- [nethttp.go](nethttp.go) - Basic net/http example
- [advanced.go](advanced.go) - Advanced features demo

## Integration Tips

1. **Choose your framework** and navigate to the corresponding example directory
2. **Copy the example code** as a starting point for your application
3. **Customize the endpoints** according to your needs
4. **Add authentication, logging, and other middleware** as required
5. **Deploy to production** with proper Cloudflare configuration

For detailed integration guides, see the [project wiki](https://github.com/RumenDamyanov/go-geolocation/wiki).
