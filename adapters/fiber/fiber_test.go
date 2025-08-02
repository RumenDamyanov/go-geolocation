package fiber

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rumendamyanov/go-geolocation"
)

func TestMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		loc := FromContext(c)
		if loc == nil {
			t.Error("location not found in context")
		}
		if loc.Country != "BG" {
			t.Errorf("expected country 'BG', got '%s'", loc.Country)
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("CF-Connecting-IP", "1.2.3.4")
	req.Header.Set("CF-IPCountry", "BG")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fiber app test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestMiddleware_MissingHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		loc := FromContext(c)
		if loc == nil {
			t.Error("location not found in context")
		}
		if loc.IP != "" || loc.Country != "" {
			t.Errorf("expected empty IP and Country, got %+v", loc)
		}
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("fiber app test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestFromContext_Helper(t *testing.T) {
	app := fiber.New()
	// No value set
	app.Get("/no-value", func(c *fiber.Ctx) error {
		if FromContext(c) != nil {
			t.Error("expected nil when no value is set")
		}
		return nil
	})
	req := httptest.NewRequest("GET", "/no-value", nil)
	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected error or status: %v, %d", err, resp.StatusCode)
	}

	// Set non-location value
	app.Get("/wrong-type", func(c *fiber.Ctx) error {
		c.Locals("geolocation", "not_a_location")
		if FromContext(c) != nil {
			t.Error("expected nil for non-location value")
		}
		return nil
	})
	req = httptest.NewRequest("GET", "/wrong-type", nil)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected error or status: %v, %d", err, resp.StatusCode)
	}

	// Set valid value
	app.Get("/valid", func(c *fiber.Ctx) error {
		loc := &geolocation.Location{Country: "BG"}
		c.Locals("geolocation", loc)
		got := FromContext(c)
		if got == nil || got.Country != "BG" {
			t.Error("expected to retrieve the set location value")
		}
		return nil
	})
	req = httptest.NewRequest("GET", "/valid", nil)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected error or status: %v, %d", err, resp.StatusCode)
	}
}
