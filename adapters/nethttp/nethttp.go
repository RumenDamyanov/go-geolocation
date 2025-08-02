package nethttp

import (
	"context"
	"net/http"

	"github.com/rumendamyanov/go-geolocation"
)

// contextKey is used for storing location in context.
type contextKey struct{}

// HTTPMiddleware adds geolocation information to the request context.
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc := geolocation.FromRequest(r)
		ctx := context.WithValue(r.Context(), contextKey{}, loc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext retrieves the geolocation info from the context.
func FromContext(ctx context.Context) *geolocation.Location {
	loc, _ := ctx.Value(contextKey{}).(*geolocation.Location)
	return loc
}
