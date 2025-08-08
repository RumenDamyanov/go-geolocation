//go:build ignore
// +build ignore

// This file is for documentation/example purposes and is not built with the main module.

package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.rumenx.com/geolocation"
	echoadapter "go.rumenx.com/geolocation/adapters/echo"
)

func main() {
	fmt.Println("🌍 Starting Echo Geolocation Example Server on :8082")
	fmt.Println("📡 Try: curl http://localhost:8082/")
	fmt.Println("🔧 Or simulate a country: curl http://localhost:8082/simulate/FR")

	e := echo.New()

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

		// Create simulated request
		simulated := geolocation.SimulateRequest(country, nil)

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
		})
	})

	e.Logger.Fatal(e.Start(":8082"))
}
