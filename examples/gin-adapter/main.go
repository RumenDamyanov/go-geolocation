package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.rumenx.com/geolocation"
	ginadapter "go.rumenx.com/geolocation/adapters/gin"
)

func main() {
	fmt.Println("🌍 Starting Gin Geolocation Example Server on :8081")
	fmt.Println("📡 Try: curl http://localhost:8081/")
	fmt.Println("🔧 Or simulate a country: curl http://localhost:8081/simulate/DE")
	fmt.Println("🗺️  Available countries: curl http://localhost:8081/countries")

	// Set Gin to release mode for production-like output
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Use geolocation middleware
	r.Use(ginadapter.Middleware())

	// Basic geolocation endpoint
	r.GET("/", func(c *gin.Context) {
		loc := ginadapter.FromContext(c)
		if loc == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "geolocation not available"})
			return
		}

		// Get additional client info
		clientInfo := geolocation.ParseClientInfo(c.Request)
		langInfo := geolocation.ParseLanguageInfo(c.Request)

		c.JSON(http.StatusOK, gin.H{
			"location":    loc,
			"client_info": clientInfo,
			"language":    langInfo,
			"is_local":    geolocation.IsLocalDevelopment(c.Request),
		})
	})

	// Simulation endpoint
	r.GET("/simulate/:country", func(c *gin.Context) {
		country := c.Param("country")
		if country == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "country parameter required"})
			return
		}

		// Create simulated request with custom options
		options := &geolocation.SimulationOptions{
			UserAgent: "Gin Example Bot/1.0",
			Languages: []string{"en", "de", "fr"},
		}
		simulated := geolocation.SimulateRequest(country, options)

		// Extract geolocation info from simulated request
		loc := geolocation.FromRequest(simulated)
		clientInfo := geolocation.ParseClientInfo(simulated)
		langInfo := geolocation.ParseLanguageInfo(simulated)

		c.JSON(http.StatusOK, gin.H{
			"simulated":   true,
			"country":     country,
			"location":    loc,
			"client_info": clientInfo,
			"language":    langInfo,
		})
	})

	// Available countries for simulation
	r.GET("/countries", func(c *gin.Context) {
		countries := geolocation.GetAvailableCountries()
		c.JSON(http.StatusOK, gin.H{
			"available_countries": countries,
			"random_country":      geolocation.RandomCountry(),
			"total_countries":     len(countries),
		})
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "gin-geolocation-example",
		})
	})

	fmt.Println("✅ Server started successfully!")
	r.Run(":8081")
}
