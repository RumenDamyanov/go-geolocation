package geolocation

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// GeolocationSimulator provides methods to simulate Cloudflare geolocation headers
// for local development when actual Cloudflare infrastructure is not available.
type GeolocationSimulator struct{}

// CountryData holds simulation data for a specific country.
type CountryData struct {
	Country   string   `json:"country"`
	IPRanges  []string `json:"ip_ranges"`
	Languages []string `json:"languages"`
	Timezone  string   `json:"timezone"`
}

// Built-in country data for simulation
var countryData = map[string]CountryData{
	"US": {
		Country:   "US",
		IPRanges:  []string{"192.168.1.", "10.0.0.", "172.16.0."},
		Languages: []string{"en-US", "en", "es"},
		Timezone:  "America/New_York",
	},
	"CA": {
		Country:   "CA",
		IPRanges:  []string{"192.168.2.", "10.0.1.", "172.16.1."},
		Languages: []string{"en-CA", "en", "fr-CA", "fr"},
		Timezone:  "America/Toronto",
	},
	"GB": {
		Country:   "GB",
		IPRanges:  []string{"192.168.3.", "10.0.2.", "172.16.2."},
		Languages: []string{"en-GB", "en"},
		Timezone:  "Europe/London",
	},
	"DE": {
		Country:   "DE",
		IPRanges:  []string{"192.168.4.", "10.0.3.", "172.16.3."},
		Languages: []string{"de-DE", "de", "en"},
		Timezone:  "Europe/Berlin",
	},
	"FR": {
		Country:   "FR",
		IPRanges:  []string{"192.168.5.", "10.0.4.", "172.16.4."},
		Languages: []string{"fr-FR", "fr", "en"},
		Timezone:  "Europe/Paris",
	},
	"JP": {
		Country:   "JP",
		IPRanges:  []string{"192.168.6.", "10.0.5.", "172.16.5."},
		Languages: []string{"ja-JP", "ja", "en"},
		Timezone:  "Asia/Tokyo",
	},
	"AU": {
		Country:   "AU",
		IPRanges:  []string{"192.168.7.", "10.0.6.", "172.16.6."},
		Languages: []string{"en-AU", "en"},
		Timezone:  "Australia/Sydney",
	},
	"BR": {
		Country:   "BR",
		IPRanges:  []string{"192.168.8.", "10.0.7.", "172.16.7."},
		Languages: []string{"pt-BR", "pt", "en"},
		Timezone:  "America/Sao_Paulo",
	},
}

// SimulationOptions holds additional options for customizing simulation.
type SimulationOptions struct {
	UserAgent  string   `json:"user_agent"`
	ServerName string   `json:"server_name"`
	IPRange    string   `json:"ip_range"`
	Languages  []string `json:"languages"`
}

// FakeCloudflareHeaders generates fake Cloudflare headers for a specific country.
func FakeCloudflareHeaders(countryCode string, options *SimulationOptions) map[string]string {
	countryCode = strings.ToUpper(countryCode)
	data, exists := countryData[countryCode]
	if !exists {
		data = countryData["US"] // fallback to US
	}

	// Generate fake IP
	ipRange := data.IPRanges[0]
	if options != nil && options.IPRange != "" {
		ipRange = options.IPRange
	}
	fakeIP := fmt.Sprintf("%s%d", ipRange, rand.Intn(254)+1)

	// Select languages
	languages := data.Languages
	if options != nil && len(options.Languages) > 0 {
		languages = options.Languages
	}

	// Build Accept-Language header
	var acceptLang strings.Builder
	for i, lang := range languages {
		if i > 0 {
			acceptLang.WriteString(",")
		}
		if i == 0 {
			acceptLang.WriteString(lang)
		} else {
			weight := 1.0 - (float64(i) * 0.1)
			acceptLang.WriteString(fmt.Sprintf("%s;q=%.1f", lang, weight))
		}
	}

	// Generate user agent
	userAgent := generateFakeUserAgent()
	if options != nil && options.UserAgent != "" {
		userAgent = options.UserAgent
	}

	// Generate CF-Ray header (fake)
	cfRay := fmt.Sprintf("%016x-%s", rand.Int63(), strings.ToLower(countryCode))

	headers := map[string]string{
		"CF-IPCountry":     countryCode,
		"CF-Connecting-IP": fakeIP,
		"CF-Ray":           cfRay,
		"Accept-Language":  acceptLang.String(),
		"User-Agent":       userAgent,
		"X-Forwarded-For":  fakeIP,
	}

	if options != nil && options.ServerName != "" {
		headers["Server-Name"] = options.ServerName
		headers["HTTP_HOST"] = options.ServerName
	}

	return headers
}

// generateFakeUserAgent returns a random realistic user agent string.
func generateFakeUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
	}
	return userAgents[rand.Intn(len(userAgents))]
}

// GetAvailableCountries returns a list of available country codes for simulation.
func GetAvailableCountries() []string {
	countries := make([]string, 0, len(countryData))
	for country := range countryData {
		countries = append(countries, country)
	}
	return countries
}

// RandomCountry returns a random country code for simulation.
func RandomCountry() string {
	countries := GetAvailableCountries()
	return countries[rand.Intn(len(countries))]
}

// AddCountryData adds custom country data for simulation.
func AddCountryData(countryCode string, data CountryData) {
	countryData[strings.ToUpper(countryCode)] = data
}

// SimulateRequest creates a new HTTP request with simulated Cloudflare headers.
func SimulateRequest(countryCode string, options *SimulationOptions) *http.Request {
	req, _ := http.NewRequest("GET", "/", nil)
	headers := FakeCloudflareHeaders(countryCode, options)

	for key, value := range headers {
		// Convert header names to standard HTTP format
		switch key {
		case "CF-IPCountry":
			req.Header.Set("CF-IPCountry", value)
		case "CF-Connecting-IP":
			req.Header.Set("CF-Connecting-IP", value)
		case "CF-Ray":
			req.Header.Set("CF-Ray", value)
		case "Accept-Language":
			req.Header.Set("Accept-Language", value)
		case "User-Agent":
			req.Header.Set("User-Agent", value)
		case "X-Forwarded-For":
			req.Header.Set("X-Forwarded-For", value)
		default:
			req.Header.Set(key, value)
		}
	}

	return req
}

// init initializes the random seed for simulation
func init() {
	rand.Seed(time.Now().UnixNano())
}
