//go:build ignore
// +build ignore

// Example demonstrating advanced geolocation features including simulation,
// language negotiation, and comprehensive client information extraction.
package main

import (
	"fmt"
	"strings"

	"go.rumenx.com/geolocation"
)

func main() {
	fmt.Println("🌍 Advanced Go Geolocation Demo")
	fmt.Println(strings.Repeat("=", 50))

	// Create configuration with country-to-language mapping
	cfg := &geolocation.Config{
		DefaultLanguage: "en",
		CountryToLanguageMap: map[string][]string{
			"US": {"en"},
			"CA": {"en", "fr"},
			"DE": {"de"},
			"FR": {"fr"},
			"JP": {"ja", "en"},
			"CH": {"de", "fr", "it"}, // Switzerland with multiple languages
		},
		CookieName: "site_lang",
	}

	// Demonstrate simulation for local development
	fmt.Println("\n🔧 Local Development Simulation")
	fmt.Println("Available countries:", geolocation.GetAvailableCountries())

	// Test different countries
	countries := []string{"US", "CA", "DE", "JP", "CH"}
	availableSiteLanguages := []string{"en", "fr", "de", "es", "ja"}

	for _, country := range countries {
		fmt.Printf("\n🌍 Simulating visitor from %s:\n", country)

		// Create simulated request with custom options
		req := geolocation.Simulate(country, &geolocation.SimulationOptions{
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
		})

		// Get comprehensive geolocation info
		info := geolocation.GetGeoInfo(req)
		fmt.Printf("  📍 Country: %s, IP: %s\n", info.CountryCode, info.IP)
		fmt.Printf("  🗣️  Languages: %v (preferred: %s)\n", info.AllLanguages, info.PreferredLanguage)
		fmt.Printf("  💻 Device: %s, Browser: %s %s, OS: %s\n",
			info.Device, info.Browser, info.BrowserVersion, info.OS)

		// Advanced language negotiation
		bestLang := geolocation.GetLanguageForCountry(req, cfg, country, availableSiteLanguages)
		fmt.Printf("  🎯 Best language for site: %s\n", bestLang)

		// Check if we should set language cookie
		shouldSet := geolocation.ShouldSetLanguage(req, cfg.CookieName)
		fmt.Printf("  🍪 Should set language cookie: %t\n", shouldSet)

		// Check if local development
		isLocal := geolocation.IsLocalDevelopment(req)
		fmt.Printf("  🏠 Local development: %t\n", isLocal)
	}

	// Demonstrate tablet detection
	fmt.Println("\n📱 Device Type Detection")
	tabletReq := geolocation.Simulate("US", &geolocation.SimulationOptions{
		UserAgent: "Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
	})
	tabletInfo := geolocation.ParseClientInfo(tabletReq)
	fmt.Printf("Tablet Detection: %s device detected\n", tabletInfo.Device)

	// Demonstrate screen resolution (would be set by frontend JavaScript)
	fmt.Println("\n🖥️  Screen Resolution")
	req := geolocation.Simulate("US", nil)
	req.Header.Set("X-Screen-Width", "1920")
	req.Header.Set("X-Screen-Height", "1080")
	resolution := geolocation.GetResolution(req)
	fmt.Printf("Resolution: %dx%d\n", resolution.Width, resolution.Height)

	// Demonstrate complex language negotiation for Switzerland
	fmt.Println("\n🇨🇭 Switzerland Language Negotiation Example")
	chReq := geolocation.Simulate("CH", &geolocation.SimulationOptions{
		Languages: []string{"fr-CH", "fr", "de", "en"}, // French speaker in Switzerland
	})

	// Site supports: English, German, French
	siteLangs := []string{"en", "de", "fr"}
	bestLangCH := geolocation.GetLanguageForCountry(chReq, cfg, "CH", siteLangs)
	fmt.Printf("French speaker in Switzerland → Best language: %s\n", bestLangCH)

	// Add custom country data
	fmt.Println("\n🔧 Custom Country Data")
	geolocation.AddCountryData("XX", geolocation.CountryData{
		Country:   "XX",
		IPRanges:  []string{"192.168.99."},
		Languages: []string{"xx-XX", "xx"},
		Timezone:  "UTC",
	})

	customReq := geolocation.Simulate("XX", nil)
	customInfo := geolocation.GetGeoInfo(customReq)
	fmt.Printf("Custom country XX: IP=%s, Languages=%v\n",
		customInfo.IP, customInfo.AllLanguages)

	fmt.Println("\n✅ Advanced geolocation demo complete!")
}
