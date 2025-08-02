# Gin Geolocation Example

This example demonstrates how to use the go-geolocation package with the Gin web framework.

## Running the Example

```bash
go run main.go
```

The server will start on port 8081.

## Endpoints

- `GET /` - Get geolocation info from request headers
- `GET /simulate/{country}` - Simulate geolocation for a specific country
- `GET /countries` - Get list of available countries for simulation
- `GET /health` - Health check endpoint

## Testing

```bash
# Basic geolocation (will show local development info)
curl http://localhost:8081/

# Simulate a specific country
curl http://localhost:8081/simulate/DE

# Get available countries
curl http://localhost:8081/countries

# Health check
curl http://localhost:8081/health
```

## Example Response

```json
{
  "location": {
    "ip": "192.168.4.123",
    "country": "DE"
  },
  "client_info": {
    "browser_name": "Chrome",
    "browser_version": "91.0.4472.124",
    "os": "Linux",
    "device": "Desktop"
  },
  "language": {
    "default": "de-DE",
    "supported": ["de-DE", "de", "en"]
  },
  "is_local": false
}
```
