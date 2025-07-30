// Package geolocation provides framework-agnostic geolocation extraction from Cloudflare headers
// and user agent parsing for browser, OS, device, and language information.
package geolocation

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mssola/user_agent"
	"gopkg.in/yaml.v3"
)

// Location represents a geolocation result, typically extracted from Cloudflare headers.
type Location struct {
	IP      string // The user's public IP address (from CF-Connecting-IP)
	Country string // The user's country code (from CF-IPCountry)
}

// ClientInfo holds browser, OS, and device information parsed from the User-Agent header.
type ClientInfo struct {
	BrowserName    string // e.g., Chrome, Firefox
	BrowserVersion string // e.g., 123.0.0.0
	OS             string // e.g., Windows NT 10.0
	Device         string // e.g., Mobile, Desktop, Tablet
}

// LanguageInfo holds the user's preferred and supported languages from Accept-Language.
type LanguageInfo struct {
	Default   string   // The default (first) language
	Supported []string // All languages in order of preference
}

// Resolution holds screen resolution information.
type Resolution struct {
	Width  int // Screen width in pixels
	Height int // Screen height in pixels
}

// GeoInfo holds all geolocation and client information.
type GeoInfo struct {
	CountryCode       string     `json:"country_code"`
	IP                string     `json:"ip"`
	PreferredLanguage string     `json:"preferred_language"`
	AllLanguages      []string   `json:"all_languages"`
	OS                string     `json:"os"`
	Browser           string     `json:"browser"`
	BrowserVersion    string     `json:"browser_version"`
	Device            string     `json:"device"`
	Resolution        Resolution `json:"resolution"`
}

// Config holds module configuration, including country-to-language mapping, defaults, and cookie name.
type Config struct {
	DefaultLanguage      string              `json:"default_language" yaml:"default_language"`
	CountryToLanguageMap map[string][]string `json:"country_to_language_map" yaml:"country_to_language_map"`
	CookieName           string              `json:"cookie_name" yaml:"cookie_name"`
}

// FromRequest extracts geolocation info from Cloudflare headers in the request.
//
// Example:
//
//	loc := geolocation.FromRequest(r)
//	fmt.Println(loc.IP, loc.Country)
func FromRequest(r *http.Request) *Location {
	return &Location{
		IP:      r.Header.Get("CF-Connecting-IP"),
		Country: r.Header.Get("CF-IPCountry"),
	}
}

// ParseClientInfo parses the User-Agent header for browser, OS, and device info.
//
// Example:
//
//	info := geolocation.ParseClientInfo(r)
//	fmt.Println(info.BrowserName, info.OS, info.Device)
func ParseClientInfo(r *http.Request) *ClientInfo {
	ua := user_agent.New(r.UserAgent())
	name, version := ua.Browser()
	device := "Desktop"
	userAgent := r.UserAgent()
	if ua.Mobile() {
		device = "Mobile"
	} else if strings.Contains(strings.ToLower(userAgent), "ipad") ||
		strings.Contains(strings.ToLower(userAgent), "tablet") {
		device = "Tablet"
	}
	return &ClientInfo{
		BrowserName:    name,
		BrowserVersion: version,
		OS:             ua.OS(),
		Device:         device,
	}
}

// ParseLanguageInfo parses the Accept-Language header for language preferences.
//
// Example:
//
//	lang := geolocation.ParseLanguageInfo(r)
//	fmt.Println(lang.Default, lang.Supported)
func ParseLanguageInfo(r *http.Request) *LanguageInfo {
	header := r.Header.Get("Accept-Language")
	if header == "" {
		return &LanguageInfo{}
	}
	langs := parseAcceptLanguage(header)
	defaultLang := ""
	if len(langs) > 0 {
		defaultLang = langs[0]
	}
	return &LanguageInfo{
		Default:   defaultLang,
		Supported: langs,
	}
}

// parseAcceptLanguage parses the Accept-Language header into a slice of language codes.
func parseAcceptLanguage(header string) []string {
	var langs []string
	for _, part := range strings.Split(header, ",") {
		lang := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		if lang != "" {
			langs = append(langs, lang)
		}
	}
	return langs
}

// LookupIP is deprecated: use FromRequest for Cloudflare geolocation.
func LookupIP(ip string) (*Location, error) {
	return &Location{IP: ip}, nil
}

// LoadConfig loads configuration from a JSON or YAML file.
// The format is determined by the file extension (.json or .yaml/.yml).
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	switch {
	case strings.HasSuffix(path, ".json"):
		err = json.Unmarshal(data, cfg)
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		err = yaml.Unmarshal(data, cfg)
	default:
		return nil, errors.New("unsupported config file format")
	}
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// ActiveLanguages returns the list of languages for a given country code, or the default if not mapped.
func (c *Config) ActiveLanguages(country string) []string {
	if langs, ok := c.CountryToLanguageMap[country]; ok && len(langs) > 0 {
		return langs
	}
	return []string{c.DefaultLanguage}
}

// ActiveLanguage returns the first language for a given country code, or the default if not mapped.
func (c *Config) ActiveLanguage(country string) string {
	langs := c.ActiveLanguages(country)
	if len(langs) > 0 {
		return langs[0]
	}
	return c.DefaultLanguage
}

// GetCookie retrieves a named cookie value from the request. Returns empty string if not found.
func GetCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// SetCookie sets a cookie with the given name and value on the response writer.
// You can pass additional options via the http.Cookie struct.
func SetCookie(w http.ResponseWriter, name, value string, opts *http.Cookie) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	if opts != nil {
		if opts.Path != "" {
			cookie.Path = opts.Path
		}
		cookie.Domain = opts.Domain
		cookie.Expires = opts.Expires
		cookie.MaxAge = opts.MaxAge
		cookie.Secure = opts.Secure
		cookie.HttpOnly = opts.HttpOnly
		cookie.SameSite = opts.SameSite
	}
	http.SetCookie(w, cookie)
}

// GetResolution retrieves screen resolution from custom headers (if set by frontend JS).
func GetResolution(r *http.Request) Resolution {
	width := 0
	height := 0
	if w := r.Header.Get("X-Screen-Width"); w != "" {
		if parsed, err := strconv.Atoi(w); err == nil {
			width = parsed
		}
	}
	if h := r.Header.Get("X-Screen-Height"); h != "" {
		if parsed, err := strconv.Atoi(h); err == nil {
			height = parsed
		}
	}
	return Resolution{Width: width, Height: height}
}

// GetGeoInfo returns all geolocation and client information in a single struct.
//
// Example:
//
//	info := geolocation.GetGeoInfo(r)
//	fmt.Printf("Country: %s, Browser: %s, Device: %s", info.CountryCode, info.Browser, info.Device)
func GetGeoInfo(r *http.Request) *GeoInfo {
	loc := FromRequest(r)
	client := ParseClientInfo(r)
	lang := ParseLanguageInfo(r)
	resolution := GetResolution(r)

	return &GeoInfo{
		CountryCode:       loc.Country,
		IP:                loc.IP,
		PreferredLanguage: lang.Default,
		AllLanguages:      lang.Supported,
		OS:                client.OS,
		Browser:           client.BrowserName,
		BrowserVersion:    client.BrowserVersion,
		Device:            client.Device,
		Resolution:        resolution,
	}
}

// GetLanguageForCountry returns the best language for a given country code,
// considering browser languages and available site languages.
//
// Logic:
// 1. If browser preferred language matches a country language and is available, use it
// 2. Check all browser languages for a match with available languages
// 3. Use the first country language as fallback
// 4. Returns empty string if no match found
func GetLanguageForCountry(r *http.Request, cfg *Config, countryCode string, availableSiteLanguages []string) string {
	if cfg == nil || countryCode == "" {
		return ""
	}

	// Check if country is actually mapped
	countryKey := strings.ToUpper(countryCode)
	_, exists := cfg.CountryToLanguageMap[countryKey]
	if !exists {
		return ""
	}

	langs := cfg.ActiveLanguages(countryKey)
	if len(langs) == 0 {
		return ""
	}

	langInfo := ParseLanguageInfo(r)

	if len(availableSiteLanguages) > 0 {
		// 1. Check preferred language
		if langInfo.Default != "" {
			preferredShort := getLanguageCode(langInfo.Default)
			if contains(langs, preferredShort) && contains(availableSiteLanguages, preferredShort) {
				return preferredShort
			}
		}

		// 2. Check all browser languages
		for _, browserLang := range langInfo.Supported {
			langCode := getLanguageCode(browserLang)
			if contains(langs, langCode) && contains(availableSiteLanguages, langCode) {
				return langCode
			}
		}

		// 3. Fallback: first country language that is available
		for _, lang := range langs {
			if contains(availableSiteLanguages, lang) {
				return lang
			}
		}
		return ""
	}

	// If no availableSiteLanguages provided, return first country language
	if len(langs) > 0 {
		return langs[0]
	}
	return ""
}

// ShouldSetLanguage returns true if the language cookie should be set (i.e., if no language cookie exists).
func ShouldSetLanguage(r *http.Request, cookieName string) bool {
	cookie := GetCookie(r, cookieName)
	return cookie == ""
}

// IsLocalDevelopment checks if we're in a local development environment.
// Returns true for localhost, local IPs, or missing Cloudflare headers.
func IsLocalDevelopment(r *http.Request) bool {
	loc := FromRequest(r)
	ip := loc.IP
	host := r.Host

	// Check for localhost, local IPs, or missing Cloudflare headers
	return ip == "" ||
		ip == "127.0.0.1" ||
		ip == "::1" ||
		strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "172.16.") ||
		strings.Contains(host, "localhost") ||
		strings.Contains(host, ".local") ||
		r.Header.Get("CF-IPCountry") == ""
}

// Simulate creates a geolocation-enabled HTTP request with simulated Cloudflare headers
// for local development and testing.
//
// Example:
//
//	req := geolocation.Simulate("DE", &geolocation.SimulationOptions{
//		UserAgent: "Custom Test Agent",
//	})
//	info := geolocation.GetGeoInfo(req)
func Simulate(countryCode string, options *SimulationOptions) *http.Request {
	return SimulateRequest(countryCode, options)
}

// Helper functions

// getLanguageCode extracts the language code from a locale string (e.g., "en-US" -> "en")
func getLanguageCode(locale string) string {
	parts := strings.Split(locale, "-")
	if len(parts) > 0 {
		return strings.ToLower(parts[0])
	}
	return locale
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
