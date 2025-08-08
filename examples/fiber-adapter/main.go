package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.rumenx.com/geolocation"
	fiberadapter "go.rumenx.com/geolocation/adapters/fiber"
)

func main() {
	fmt.Println("🌍 Starting Fiber Geolocation Example Server on :8083")
	fmt.Println("📡 Try: curl http://localhost:8083/")
	fmt.Println("🔧 Or simulate a country: curl http://localhost:8083/simulate/JP")
	fmt.Println("🗺️  Available countries: curl http://localhost:8083/countries")

	app := fiber.New(fiber.Config{
		AppName:          "Fiber Geolocation Example",
		DisableKeepalive: true,
	})

	// Use geolocation middleware
	app.Use(fiberadapter.Middleware())

	// Basic geolocation endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		loc := fiberadapter.FromContext(c)
		if loc == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "geolocation not available",
			})
		}

		// For Fiber, we extract info directly from headers since it uses fasthttp
		userAgent := c.Get("User-Agent")
		acceptLang := c.Get("Accept-Language")

		return c.JSON(fiber.Map{
			"location":    loc,
			"user_agent":  userAgent,
			"accept_lang": acceptLang,
			"is_local":    loc.IP == "" || loc.Country == "",
		})
	})

	// Simulation endpoint
	app.Get("/simulate/:country", func(c *fiber.Ctx) error {
		country := c.Params("country")
		if country == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "country parameter required",
			})
		}

		// Create simulated request with custom options
		options := &geolocation.SimulationOptions{
			UserAgent: "Fiber Example Bot/1.0",
			Languages: []string{"ja", "en", "de"},
		}
		simulated := geolocation.SimulateRequest(country, options)

		// Extract geolocation info from simulated request
		loc := geolocation.FromRequest(simulated)
		clientInfo := geolocation.ParseClientInfo(simulated)
		langInfo := geolocation.ParseLanguageInfo(simulated)

		return c.JSON(fiber.Map{
			"simulated":   true,
			"country":     country,
			"location":    loc,
			"client_info": clientInfo,
			"language":    langInfo,
		})
	})

	// Available countries for simulation
	app.Get("/countries", func(c *fiber.Ctx) error {
		countries := geolocation.GetAvailableCountries()
		return c.JSON(fiber.Map{
			"available_countries": countries,
			"random_country":      geolocation.RandomCountry(),
			"total_countries":     len(countries),
		})
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "fiber-geolocation-example",
		})
	})

	fmt.Println("✅ Server started successfully!")
	app.Listen(":8083")
}
