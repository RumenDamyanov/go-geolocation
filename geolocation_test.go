package geolocation

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLookupIP(t *testing.T) {
	ip := "8.8.8.8"
	loc, err := LookupIP(ip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.IP != ip {
		t.Errorf("expected IP %s, got %s", ip, loc.IP)
	}
}

func TestFromRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("CF-Connecting-IP", "1.2.3.4")
	r.Header.Set("CF-IPCountry", "BG")
	loc := FromRequest(r)
	if loc.IP != "1.2.3.4" || loc.Country != "BG" {
		t.Errorf("unexpected location: %+v", loc)
	}
}

func TestFromRequest_MissingHeaders(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	loc := FromRequest(r)
	if loc.IP != "" || loc.Country != "" {
		t.Errorf("expected empty IP and Country, got %+v", loc)
	}
}

func TestParseClientInfo(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	info := ParseClientInfo(r)
	if info.BrowserName == "" || info.OS == "" || info.Device == "" {
		t.Errorf("unexpected client info: %+v", info)
	}
}

func TestParseClientInfo_EmptyUserAgent(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	info := ParseClientInfo(r)
	if info.BrowserName != "" || info.OS != "" || info.Device != "Desktop" {
		t.Errorf("unexpected client info for empty UA: %+v", info)
	}
}

func TestParseClientInfo_Mobile(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	info := ParseClientInfo(r)
	if info.Device != "Mobile" {
		t.Errorf("expected device to be 'Mobile', got '%s'", info.Device)
	}
}

func TestParseLanguageInfo(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "en-US,en;q=0.9,bg;q=0.8,de;q=0.7")
	lang := ParseLanguageInfo(r)
	if lang.Default != "en-US" {
		t.Errorf("expected default 'en-US', got '%s'", lang.Default)
	}
	if len(lang.Supported) != 4 || lang.Supported[2] != "bg" {
		t.Errorf("unexpected supported languages: %+v", lang.Supported)
	}
}

func TestParseLanguageInfo_EmptyHeader(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	lang := ParseLanguageInfo(r)
	if lang.Default != "" || len(lang.Supported) != 0 {
		t.Errorf("unexpected language info for empty header: %+v", lang)
	}
}

func TestParseLanguageInfo_WeirdHeader(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", ",,,en-US;q=0.9,,bg;q=0.8,,")
	lang := ParseLanguageInfo(r)
	if lang.Default != "en-US" || len(lang.Supported) != 2 || lang.Supported[1] != "bg" {
		t.Errorf("unexpected language info for weird header: %+v", lang)
	}
}

func TestConfig_ActiveLanguages(t *testing.T) {
	cfg := &Config{
		DefaultLanguage: "en",
		CountryToLanguageMap: map[string][]string{
			"BG": {"bg"},
			"CA": {"en", "fr"},
		},
	}
	langs := cfg.ActiveLanguages("CA")
	if len(langs) != 2 || langs[0] != "en" || langs[1] != "fr" {
		t.Errorf("unexpected active languages: %+v", langs)
	}
	langs = cfg.ActiveLanguages("XX")
	if len(langs) != 1 || langs[0] != "en" {
		t.Errorf("unexpected fallback language: %+v", langs)
	}
}

func TestConfig_ActiveLanguage(t *testing.T) {
	cfg := &Config{
		DefaultLanguage: "en",
		CountryToLanguageMap: map[string][]string{
			"BG": {"bg"},
			"CA": {"en", "fr"},
		},
	}
	if lang := cfg.ActiveLanguage("CA"); lang != "en" {
		t.Errorf("expected 'en', got '%s'", lang)
	}
	if lang := cfg.ActiveLanguage("XX"); lang != "en" {
		t.Errorf("expected fallback 'en', got '%s'", lang)
	}
}

func TestConfig_ActiveLanguage_EmptySlice(t *testing.T) {
	cfg := &Config{
		DefaultLanguage: "en",
		CountryToLanguageMap: map[string][]string{
			"ZZ": {},
		},
	}
	lang := cfg.ActiveLanguage("ZZ")
	if lang != "en" {
		t.Errorf("expected fallback 'en', got '%s'", lang)
	}
}

func TestConfig_ActiveLanguage_NilConfig(t *testing.T) {
	var cfg *Config
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on nil Config")
		}
	}()
	_ = cfg.ActiveLanguage("US")
}

func TestConfig_ActiveLanguages_NilMap(t *testing.T) {
	cfg := &Config{DefaultLanguage: "en"}
	langs := cfg.ActiveLanguages("US")
	if len(langs) != 1 || langs[0] != "en" {
		t.Errorf("expected fallback 'en', got %+v", langs)
	}
}

func TestGetSetCookie(t *testing.T) {
	r, w := httptest.NewRequest("GET", "/", nil), httptest.NewRecorder()
	SetCookie(w, "test_cookie", "test_value", &http.Cookie{Path: "/", MaxAge: 60})
	resp := w.Result()
	for _, c := range resp.Cookies() {
		r.AddCookie(c)
	}
	val := GetCookie(r, "test_cookie")
	if val != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", val)
	}
	if GetCookie(r, "missing_cookie") != "" {
		t.Error("expected empty string for missing cookie")
	}
}

func TestLoadConfig_JSON(t *testing.T) {
	jsonData := `{"default_language":"en","country_to_language_map":{"BG":["bg"],"CA":["en","fr"]},"cookie_name":"geo_lang"}`
	f, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString(jsonData)
	f.Close()
	cfg, err := LoadConfig(f.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.DefaultLanguage != "en" || cfg.CookieName != "geo_lang" || len(cfg.CountryToLanguageMap["CA"]) != 2 {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfig_YAML(t *testing.T) {
	yamlData := `default_language: en
country_to_language_map:
  BG: [bg]
  CA: [en, fr]
cookie_name: geo_lang
`
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString(yamlData)
	f.Close()
	cfg, err := LoadConfig(f.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.DefaultLanguage != "en" || cfg.CookieName != "geo_lang" || len(cfg.CountryToLanguageMap["CA"]) != 2 {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString("{invalid json}")
	f.Close()
	_, err = LoadConfig(f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString(": invalid yaml")
	f.Close()
	_, err = LoadConfig(f.Name())
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadConfig_InvalidYML(t *testing.T) {
	f, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.WriteString(": invalid yml")
	f.Close()
	_, err = LoadConfig(f.Name())
	if err == nil {
		t.Error("expected error for invalid YML")
	}
}
