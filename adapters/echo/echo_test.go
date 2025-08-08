package echo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"go.rumenx.com/geolocation"
)

func TestMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(Middleware())
	e.GET("/", func(c echo.Context) error {
		loc := FromContext(c)
		if loc == nil {
			t.Error("location not found in context")
		}
		if loc.Country != "BG" {
			t.Errorf("expected country 'BG', got '%s'", loc.Country)
		}
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("CF-Connecting-IP", "1.2.3.4")
	req.Header.Set("CF-IPCountry", "BG")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestMiddleware_MissingHeaders(t *testing.T) {
	e := echo.New()
	e.Use(Middleware())
	e.GET("/", func(c echo.Context) error {
		loc := FromContext(c)
		if loc == nil {
			t.Error("location not found in context")
		}
		if loc.IP != "" || loc.Country != "" {
			t.Errorf("expected empty IP and Country, got %+v", loc)
		}
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestFromContext_Helper(t *testing.T) {
	e := echo.New()
	c := e.NewContext(nil, nil)
	// No value set
	if FromContext(c) != nil {
		t.Error("expected nil when no value is set")
	}
	// Set non-location value
	c.Set("geolocation", "not_a_location")
	if FromContext(c) != nil {
		t.Error("expected nil for non-location value")
	}
	// Set valid value
	loc := &geolocation.Location{Country: "BG"}
	c.Set("geolocation", loc)
	if got := FromContext(c); got == nil || got.Country != "BG" {
		t.Error("expected to retrieve the set location value")
	}
}
