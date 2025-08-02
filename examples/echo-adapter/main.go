package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rumendamyanov/go-geolocation"
	echoadapter "github.com/rumendamyanov/go-geolocation/adapters/echo"
)

func main() {
	fmt.Println("🌍 Starting Echo Geolocation Example Server on :8082")
	fmt.Println("📡 Try: curl http://localhost:8082/")
	fmt.Println("🔧 Or simulate a country: curl http://localhost:8082/simulate/FR")
	fmt.Println("🗺️  Available countries: curl http://localhost:8082/countries")

	e := echo.New()

	// Hide Echo banner for cleaner output
	e.HideBanner = true

	// Use geolocation middleware
	e.Use(echoadapter.Middleware())

	// Basic geolocation endpoint
	e.GET("/", func(c echo.Context) error {
		loc := echoadapter.FromContext(c)
		if loc == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "geolocation not available",
			})
		}

		// Get additional client info
		clientInfo := geolocation.ParseClientInfo(c.Request())
		langInfo := geolocation.ParseLanguageInfo(c.Request())

		return c.JSON(http.StatusOK, map[string]interface{}{
			"location":    loc,
			"client_info": clientInfo,
			"language":    langInfo,
			"is_local":    geolocation.IsLocalDevelopment(c.Request()),
		})
	})

	// Simulation endpoint
	e.GET("/simulate/:country", func(c echo.Context) error {
		country := c.Param("country")
		if country == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "country parameter required",
			})
		}

		// Create simulated request with custom options
		options := &geolocation.SimulationOptions{
			UserAgent: "Echo Example Bot/1.0",
			Languages: []string{"fr", "en", "de"},
		}
		simulated := geolocation.SimulateRequest(country, options)

		// Extract geolocation info from simulated request
		loc := geolocation.FromRequest(simulated)
		clientInfo := geolocation.ParseClientInfo(simulated)
		langInfo := geolocation.ParseLanguageInfo(simulated)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"simulated":   true,
			"country":     country,
			"location":    loc,
			"client_info": clientInfo,
			"language":    langInfo,
		})
	})

	// Available countries for simulation
	e.GET("/countries", func(c echo.Context) error {
		countries := geolocation.GetAvailableCountries()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"available_countries": countries,
			"random_country":      geolocation.RandomCountry(),
			"total_countries":     len(countries),
		})
	})

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "echo-geolocation-example",
		})
	})

	fmt.Println("✅ Server started successfully!")
	e.Logger.Fatal(e.Start(":8082"))
}
