package geolocation

import (
	"context"
	"net/http"
)

// contextKey is used for storing location in context.
type contextKey struct{}

// HTTPMiddleware attaches geolocation info to the request context.
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc := FromRequest(r)
		ctx := context.WithValue(r.Context(), contextKey{}, loc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext retrieves the Location from context.
func FromContext(ctx context.Context) *Location {
	loc, _ := ctx.Value(contextKey{}).(*Location)
	return loc
}
