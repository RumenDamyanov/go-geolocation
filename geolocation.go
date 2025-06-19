// Package geolocation provides framework-agnostic geolocation extraction from Cloudflare headers
// and user agent parsing for browser, OS, device, and language information.
package geolocation

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
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
	Device         string // e.g., Mobile, Desktop
}

// LanguageInfo holds the user's preferred and supported languages from Accept-Language.
type LanguageInfo struct {
	Default   string   // The default (first) language
	Supported []string // All languages in order of preference
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
	if ua.Mobile() {
		device = "Mobile"
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
