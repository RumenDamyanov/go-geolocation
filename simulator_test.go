package geolocation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeolocationSimulator(t *testing.T) {
	// Test FakeCloudflareHeaders
	headers := FakeCloudflareHeaders("DE", nil)
	if headers["CF-IPCountry"] != "DE" {
		t.Errorf("expected CF-IPCountry to be 'DE', got '%s'", headers["CF-IPCountry"])
	}
	if headers["CF-Connecting-IP"] == "" {
		t.Error("expected CF-Connecting-IP to be set")
	}
	if !strings.Contains(headers["Accept-Language"], "de") {
		t.Error("expected Accept-Language to contain 'de'")
	}

	// Test with custom options
	options := &SimulationOptions{
		UserAgent: "Custom Test Agent",
		Languages: []string{"en", "fr"},
	}
	headers = FakeCloudflareHeaders("CA", options)
	if headers["User-Agent"] != "Custom Test Agent" {
		t.Errorf("expected custom user agent, got '%s'", headers["User-Agent"])
	}
	if !strings.Contains(headers["Accept-Language"], "en") {
		t.Error("expected Accept-Language to contain 'en'")
	}
}

func TestGetAvailableCountries(t *testing.T) {
	countries := GetAvailableCountries()
	expectedCountries := []string{"US", "CA", "GB", "DE", "FR", "JP", "AU", "BR"}

	if len(countries) != len(expectedCountries) {
		t.Errorf("expected %d countries, got %d", len(expectedCountries), len(countries))
	}

	for _, expected := range expectedCountries {
		found := false
		for _, country := range countries {
			if country == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected country '%s' not found in available countries", expected)
		}
	}
}

func TestRandomCountry(t *testing.T) {
	country := RandomCountry()
	countries := GetAvailableCountries()

	found := false
	for _, c := range countries {
		if c == country {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("random country '%s' not in available countries", country)
	}
}

func TestFakeCloudflareHeaders_EdgeCases(t *testing.T) {
	// Test with empty country
	headers := FakeCloudflareHeaders("", nil)

	// Function sets the countryCode as provided, but uses US data
	if headers["CF-IPCountry"] != "" {
		t.Errorf("expected empty string for CF-IPCountry with empty country, got %q", headers["CF-IPCountry"])
	}
	// Should use US IP ranges as fallback
	if !strings.HasPrefix(headers["CF-Connecting-IP"], "192.168.1.") {
		t.Errorf("expected US IP range fallback, got %q", headers["CF-Connecting-IP"])
	}

	// Test with whitespace-only country
	headers = FakeCloudflareHeaders("  ", nil)

	// Function converts to uppercase but retains whitespace
	if headers["CF-IPCountry"] != "  " {
		t.Errorf("expected whitespace preserved for CF-IPCountry, got %q", headers["CF-IPCountry"])
	}

	// Test with non-existent country
	headers = FakeCloudflareHeaders("XX", nil)

	// Function returns XX as country but uses US data as fallback
	if headers["CF-IPCountry"] != "XX" {
		t.Errorf("expected 'XX' for CF-IPCountry, got %q", headers["CF-IPCountry"])
	}
	// Should use US IP ranges as fallback
	if !strings.HasPrefix(headers["CF-Connecting-IP"], "192.168.1.") {
		t.Errorf("expected US IP range fallback for non-existent country, got %q", headers["CF-Connecting-IP"])
	}

	// Test with custom simulation options
	options := &SimulationOptions{
		UserAgent:  "Custom Agent",
		ServerName: "test.example.com",
		IPRange:    "203.0.113.",
		Languages:  []string{"es", "en"},
	}
	headers = FakeCloudflareHeaders("US", options)

	if headers["User-Agent"] != "Custom Agent" {
		t.Errorf("expected custom user agent, got %q", headers["User-Agent"])
	}
	if headers["Server-Name"] != "test.example.com" {
		t.Errorf("expected custom server name, got %q", headers["Server-Name"])
	}
	if !strings.HasPrefix(headers["CF-Connecting-IP"], "203.0.113.") {
		t.Errorf("expected IP with custom range, got %q", headers["CF-Connecting-IP"])
	}
	if !strings.Contains(headers["Accept-Language"], "es") {
		t.Errorf("expected custom languages in Accept-Language, got %q", headers["Accept-Language"])
	}
}

func TestSimulateRequest_EdgeCases(t *testing.T) {
	// Test with empty country (should fallback to US)
	req := SimulateRequest("", nil)

	if req.Header.Get("CF-IPCountry") != "" {
		t.Errorf("expected empty string for CF-IPCountry with empty country, got %q", req.Header.Get("CF-IPCountry"))
	}

	// Test with non-existent country code
	req = SimulateRequest("XX", nil)

	// Should return XX as country but use US data as fallback
	if req.Header.Get("CF-IPCountry") != "XX" {
		t.Errorf("expected 'XX' for CF-IPCountry, got %q", req.Header.Get("CF-IPCountry"))
	}

	// Test with custom options that include server name (triggers default case)
	options := &SimulationOptions{
		UserAgent:  "Test Agent",
		ServerName: "test.example.com",
		Languages:  []string{"de", "en"},
	}
	req = SimulateRequest("DE", options)

	if req.Header.Get("User-Agent") != "Test Agent" {
		t.Errorf("expected custom user agent, got %q", req.Header.Get("User-Agent"))
	}
	if req.Header.Get("CF-IPCountry") != "DE" {
		t.Errorf("expected 'DE', got %q", req.Header.Get("CF-IPCountry"))
	}
	// This should trigger the default case in the switch statement
	if req.Header.Get("Server-Name") != "test.example.com" {
		t.Errorf("expected Server-Name header to be set via default case, got %q", req.Header.Get("Server-Name"))
	}
	if req.Header.Get("HTTP_HOST") != "test.example.com" {
		t.Errorf("expected HTTP_HOST header to be set via default case, got %q", req.Header.Get("HTTP_HOST"))
	}

	// Verify all expected headers are set
	expectedHeaders := []string{"CF-IPCountry", "CF-Connecting-IP", "CF-Ray", "Accept-Language", "User-Agent", "X-Forwarded-For"}
	for _, header := range expectedHeaders {
		if req.Header.Get(header) == "" {
			t.Errorf("expected header %q to be set", header)
		}
	}
}

func TestGetResolution(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Screen-Width", "1920")
	req.Header.Set("X-Screen-Height", "1080")

	resolution := GetResolution(req)
	if resolution.Width != 1920 {
		t.Errorf("expected width 1920, got %d", resolution.Width)
	}
	if resolution.Height != 1080 {
		t.Errorf("expected height 1080, got %d", resolution.Height)
	}

	// Test with missing headers
	req2 := httptest.NewRequest("GET", "/", nil)
	resolution2 := GetResolution(req2)
	if resolution2.Width != 0 || resolution2.Height != 0 {
		t.Errorf("expected zero resolution for missing headers, got %+v", resolution2)
	}
}

func TestGetGeoInfo(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("CF-IPCountry", "US")
	req.Header.Set("CF-Connecting-IP", "192.168.1.1")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("X-Screen-Width", "1920")
	req.Header.Set("X-Screen-Height", "1080")

	info := GetGeoInfo(req)
	if info.CountryCode != "US" {
		t.Errorf("expected country code 'US', got '%s'", info.CountryCode)
	}
	if info.IP != "192.168.1.1" {
		t.Errorf("expected IP '192.168.1.1', got '%s'", info.IP)
	}
	if info.PreferredLanguage != "en-US" {
		t.Errorf("expected preferred language 'en-US', got '%s'", info.PreferredLanguage)
	}
	if info.Resolution.Width != 1920 {
		t.Errorf("expected resolution width 1920, got %d", info.Resolution.Width)
	}
}

func TestGetLanguageForCountry(t *testing.T) {
	cfg := &Config{
		DefaultLanguage: "en",
		CountryToLanguageMap: map[string][]string{
			"CA": {"en", "fr"},
			"DE": {"de"},
		},
	}

	// Test with browser language matching country and available languages
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "fr-CA,fr;q=0.9,en;q=0.8")

	lang := GetLanguageForCountry(req, cfg, "CA", []string{"en", "fr", "de"})
	if lang != "fr" {
		t.Errorf("expected language 'fr', got '%s'", lang)
	}

	// Test fallback to first country language
	lang2 := GetLanguageForCountry(req, cfg, "CA", []string{"en", "de"})
	if lang2 != "en" {
		t.Errorf("expected fallback language 'en', got '%s'", lang2)
	}

	// Test without available site languages
	lang3 := GetLanguageForCountry(req, cfg, "DE", nil)
	if lang3 != "de" {
		t.Errorf("expected language 'de', got '%s'", lang3)
	}

	// Test unmapped country
	lang4 := GetLanguageForCountry(req, cfg, "XX", []string{"en"})
	if lang4 != "" {
		t.Errorf("expected empty string for unmapped country, got '%s'", lang4)
	}
}

func TestShouldSetLanguage(t *testing.T) {
	// Test without cookie
	req := httptest.NewRequest("GET", "/", nil)
	if !ShouldSetLanguage(req, "lang") {
		t.Error("expected ShouldSetLanguage to return true when no cookie is set")
	}

	// Test with cookie
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.AddCookie(&http.Cookie{Name: "lang", Value: "en"})
	if ShouldSetLanguage(req2, "lang") {
		t.Error("expected ShouldSetLanguage to return false when cookie is set")
	}
}

func TestIsLocalDevelopment(t *testing.T) {
	// Test localhost
	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	if !IsLocalDevelopment(req) {
		t.Error("expected IsLocalDevelopment to return true for localhost")
	}

	// Test missing CF headers
	req2 := httptest.NewRequest("GET", "/", nil)
	if !IsLocalDevelopment(req2) {
		t.Error("expected IsLocalDevelopment to return true when CF headers are missing")
	}

	// Test with CF headers (production)
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.Header.Set("CF-IPCountry", "US")
	req3.Header.Set("CF-Connecting-IP", "203.0.113.1")
	req3.Host = "example.com"
	if IsLocalDevelopment(req3) {
		t.Error("expected IsLocalDevelopment to return false for production with CF headers")
	}
}

func TestSimulate(t *testing.T) {
	req := Simulate("FR", &SimulationOptions{
		UserAgent: "Test Agent",
	})

	if req.Header.Get("CF-IPCountry") != "FR" {
		t.Errorf("expected simulated country 'FR', got '%s'", req.Header.Get("CF-IPCountry"))
	}
	if req.Header.Get("User-Agent") != "Test Agent" {
		t.Errorf("expected custom user agent, got '%s'", req.Header.Get("User-Agent"))
	}
}

func TestParseClientInfo_Tablet(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15")

	info := ParseClientInfo(req)
	if info.Device != "Tablet" {
		t.Errorf("expected device to be 'Tablet', got '%s'", info.Device)
	}
}

func TestAddCountryData(t *testing.T) {
	customData := CountryData{
		Country:   "XX",
		IPRanges:  []string{"192.168.99."},
		Languages: []string{"xx-XX", "xx"},
		Timezone:  "UTC",
	}

	AddCountryData("XX", customData)

	headers := FakeCloudflareHeaders("XX", nil)
	if headers["CF-IPCountry"] != "XX" {
		t.Errorf("expected custom country 'XX', got '%s'", headers["CF-IPCountry"])
	}
	if !strings.Contains(headers["Accept-Language"], "xx") {
		t.Error("expected Accept-Language to contain custom language 'xx'")
	}
}
