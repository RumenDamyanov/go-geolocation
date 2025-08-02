package echo

import (
	"github.com/labstack/echo/v4"
	"github.com/rumendamyanov/go-geolocation"
)

// Middleware attaches geolocation info to Echo context.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			loc := geolocation.FromRequest(c.Request())
			c.Set("geolocation", loc)
			return next(c)
		}
	}
}

// FromContext retrieves the Location from Echo context.
func FromContext(c echo.Context) *geolocation.Location {
	loc := c.Get("geolocation")
	if l, ok := loc.(*geolocation.Location); ok {
		return l
	}
	return nil
}
