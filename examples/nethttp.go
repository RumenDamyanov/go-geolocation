//go:build ignore
// +build ignore

// This file is for documentation/example purposes and is not built with the main module.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rumendamyanov/go-geolocation"
)

var cfg *geolocation.Config

func handler(w http.ResponseWriter, r *http.Request) {
	loc := geolocation.FromRequest(r)
	info := geolocation.ParseClientInfo(r)
	lang := geolocation.ParseLanguageInfo(r)
	activeLangs := cfg.ActiveLanguages(loc.Country)

	// Check for existing cookie
	cookieVal := geolocation.GetCookie(r, cfg.CookieName)
	if cookieVal == "" {
		// Set cookie to first active language if not present
		geolocation.SetCookie(w, cfg.CookieName, activeLangs[0], &http.Cookie{
			Path:     "/",
			MaxAge:   86400 * 30, // 30 days
			HttpOnly: true,
			Expires:  time.Now().Add(30 * 24 * time.Hour),
		})
		cookieVal = activeLangs[0]
	}

	fmt.Fprintf(w, "IP: %s\nCountry: %s\nActiveLangs: %v\nCookie: %s=%s\nBrowser: %s %s\nOS: %s\nDevice: %s\nDefaultLang: %s\nAllLangs: %v\n",
		loc.IP, loc.Country, activeLangs, cfg.CookieName, cookieVal, info.BrowserName, info.BrowserVersion, info.OS, info.Device, lang.Default, lang.Supported)
}

func main() {
	var err error
	cfg, err = geolocation.LoadConfig("examples/config.yaml") // or config.json
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	http.HandleFunc("/", handler)
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
