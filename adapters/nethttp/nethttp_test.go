package nethttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.rumenx.com/geolocation"
)

func TestHTTPMiddleware(t *testing.T) {
	h := HTTPMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc := FromContext(r.Context())
		if loc == nil {
			t.Error("location not found in context")
		}
	}))
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
}

func TestHTTPMiddleware_WithHeaders(t *testing.T) {
	h := HTTPMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc := FromContext(r.Context())
		if loc == nil || loc.IP != "2.2.2.2" || loc.Country != "DE" {
			t.Errorf("expected IP 2.2.2.2 and country DE, got %+v", loc)
		}
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("CF-Connecting-IP", "2.2.2.2")
	req.Header.Set("CF-IPCountry", "DE")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	loc := &geolocation.Location{IP: "1.2.3.4"}
	ctx = context.WithValue(ctx, contextKey{}, loc)
	got := FromContext(ctx)
	if got == nil || got.IP != "1.2.3.4" {
		t.Error("FromContext did not return correct location")
	}
}

func TestFromContext_Nil(t *testing.T) {
	ctx := context.Background()
	if loc := FromContext(ctx); loc != nil {
		t.Error("expected nil location from empty context")
	}
}
