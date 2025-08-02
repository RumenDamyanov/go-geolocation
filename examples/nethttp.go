//go:build ignore
// +build ignore

// This file is for documentation/example purposes and is not built with the main module.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rumendamyanov/go-geolocation"
	httpadapter "github.com/rumendamyanov/go-geolocation/adapters/nethttp"
)

func main() {
	fmt.Println("🌍 Starting net/http Geolocation Example Server on :8080")
	fmt.Println("📡 Try: curl http://localhost:8080/")
	fmt.Println("🔧 Or simulate a country: curl http://localhost:8080/simulate/US")

	mux := http.NewServeMux()

	// Basic geolocation endpoint with middleware
	mux.Handle("/", httpadapter.HTTPMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc := httpadapter.FromContext(r.Context())
		if loc == nil {
			http.Error(w, `{"error": "geolocation not available"}`, http.StatusInternalServerError)
			return
		}

		// Get additional client info
		clientInfo := geolocation.ParseClientInfo(r)
		langInfo := geolocation.ParseLanguageInfo(r)

		response := map[string]interface{}{
			"location":    loc,
			"client_info": clientInfo,
			"language":    langInfo,
			"is_local":    geolocation.IsLocalDevelopment(r),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})))

	// Simulation endpoint
	mux.HandleFunc("/simulate/", func(w http.ResponseWriter, r *http.Request) {
		country := strings.TrimPrefix(r.URL.Path, "/simulate/")
		if country == "" {
			http.Error(w, `{"error": "country parameter required"}`, http.StatusBadRequest)
			return
		}

		// Create simulated request
		simulated := geolocation.SimulateRequest(country, nil)

		// Extract geolocation info from simulated request
		loc := geolocation.FromRequest(simulated)
		clientInfo := geolocation.ParseClientInfo(simulated)
		langInfo := geolocation.ParseLanguageInfo(simulated)

		response := map[string]interface{}{
			"simulated":   true,
			"country":     country,
			"location":    loc,
			"client_info": clientInfo,
			"language":    langInfo,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Available countries for simulation
	mux.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
		countries := geolocation.GetAvailableCountries()
		response := map[string]interface{}{
			"available_countries": countries,
			"random_country":      geolocation.RandomCountry(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
