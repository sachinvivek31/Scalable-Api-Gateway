package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Define a custom type for context keys to avoid collisions
type contextKey string

const RequestIDKey contextKey = "requestID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add it to the request context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		
		// Add it to the response header so the client sees it
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}