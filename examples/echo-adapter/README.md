# Echo Geolocation Example

This example demonstrates how to use the go-geolocation package with the Echo web framework.

## Running the Example

```bash
go run main.go
```

The server will start on port 8082.

## Endpoints

- `GET /` - Get geolocation info from request headers
- `GET /simulate/{country}` - Simulate geolocation for a specific country
- `GET /countries` - Get list of available countries for simulation
- `GET /health` - Health check endpoint

## Testing

```bash
# Basic geolocation (will show local development info)
curl http://localhost:8082/

# Simulate a specific country
curl http://localhost:8082/simulate/FR

# Get available countries
curl http://localhost:8082/countries

# Health check
curl http://localhost:8082/health
```

## Example Response

```json
{
  "location": {
    "ip": "192.168.5.123",
    "country": "FR"
  },
  "client_info": {
    "browser_name": "Firefox",
    "browser_version": "89.0",
    "os": "Windows NT 10.0",
    "device": "Desktop"
  },
  "language": {
    "default": "fr-FR",
    "supported": ["fr-FR", "fr", "en"]
  },
  "is_local": false
}
```
