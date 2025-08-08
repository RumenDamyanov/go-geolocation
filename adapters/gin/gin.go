package gin

import (
	"github.com/gin-gonic/gin"
	"go.rumenx.com/geolocation"
)

type contextKey struct{}

// Middleware attaches geolocation info to Gin context.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loc := geolocation.FromRequest(c.Request)
		c.Set("geolocation", loc)
		c.Next()
	}
}

// FromContext retrieves the Location from Gin context.
func FromContext(c *gin.Context) *geolocation.Location {
	loc, _ := c.Get("geolocation")
	if l, ok := loc.(*geolocation.Location); ok {
		return l
	}
	return nil
}
