package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const APIIDKey contextKey = "api_id"

// APIIDMiddleware generates a unique API ID for each request
// and adds it to the request context for tracing
func APIIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate unique API ID
		apiID := uuid.New().String()

		// Add to context
		ctx := context.WithValue(r.Context(), APIIDKey, apiID)

		// Add to response header for easy debugging
		w.Header().Set("X-API-ID", apiID)

		// Continue with the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAPIID retrieves the API ID from the request context
func GetAPIID(ctx context.Context) string {
	if apiID, ok := ctx.Value(APIIDKey).(string); ok {
		return apiID
	}
	return "unknown"
}
