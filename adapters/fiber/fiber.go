package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rumendamyanov/go-geolocation"
)

// Middleware attaches geolocation info to Fiber context.
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		loc := &geolocation.Location{
			IP:      c.Get("CF-Connecting-IP"),
			Country: c.Get("CF-IPCountry"),
		}
		c.Locals("geolocation", loc)
		return c.Next()
	}
}

// FromContext retrieves the Location from Fiber context.
func FromContext(c *fiber.Ctx) *geolocation.Location {
	loc := c.Locals("geolocation")
	if l, ok := loc.(*geolocation.Location); ok {
		return l
	}
	return nil
}
