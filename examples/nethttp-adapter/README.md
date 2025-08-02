# net/http Geolocation Example

This example demonstrates how to use the go-geolocation package with the standard net/http library.

## Running the Example

```bash
go run main.go
```

The server will start on port 8080.

## Endpoints

- `GET /` - Get geolocation info from request headers
- `GET /simulate/{country}` - Simulate geolocation for a specific country
- `GET /countries` - Get list of available countries for simulation
- `GET /health` - Health check endpoint

## Testing

```bash
# Basic geolocation (will show local development info)
curl http://localhost:8080/

# Simulate a specific country
curl http://localhost:8080/simulate/US

# Get available countries
curl http://localhost:8080/countries

# Health check
curl http://localhost:8080/health
```

## Example Response

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
