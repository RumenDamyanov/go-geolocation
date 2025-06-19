package ginadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rumendamyanov/go-geolocation"
)

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware())
	r.GET("/", func(c *gin.Context) {
		loc := FromContext(c)
		if loc == nil {
			t.Error("location not found in context")
		}
		if loc.Country != "BG" {
			t.Errorf("expected country 'BG', got '%s'", loc.Country)
		}
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("CF-Connecting-IP", "1.2.3.4")
	req.Header.Set("CF-IPCountry", "BG")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestMiddleware_MissingHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Middleware())
	r.GET("/", func(c *gin.Context) {
		loc := FromContext(c)
		if loc == nil {
			t.Error("location not found in context")
		}
		if loc.IP != "" || loc.Country != "" {
			t.Errorf("expected empty IP and Country, got %+v", loc)
		}
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestFromContext(t *testing.T) {
	c := &gin.Context{}
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
