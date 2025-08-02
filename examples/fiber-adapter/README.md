# Fiber Geolocation Example

This example demonstrates how to use the go-geolocation package with the Fiber web framework.

## Running the Example

```bash
go run main.go
```

The server will start on port 8083.

## Endpoints

- `GET /` - Get geolocation info from request headers
- `GET /simulate/{country}` - Simulate geolocation for a specific country
- `GET /countries` - Get list of available countries for simulation
- `GET /health` - Health check endpoint

## Testing

```bash
# Basic geolocation (will show local development info)
curl http://localhost:8083/

# Simulate a specific country
curl http://localhost:8083/simulate/JP

# Get available countries
curl http://localhost:8083/countries

# Health check
curl http://localhost:8083/health
```

## Example Response

```json
{
  "location": {
    "ip": "192.168.6.123",
    "country": "JP"
  },
  "user_agent": "Fiber Example Bot/1.0",
  "accept_lang": "ja,en;q=0.9,de;q=0.8",
  "is_local": false
}
```

## Note

Since Fiber uses fasthttp instead of net/http, the client info parsing is slightly different. This example shows how to work with the Fiber-specific context and header extraction.
