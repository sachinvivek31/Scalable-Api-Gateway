package middleware

import "net/http"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		
		// In a real app, use a JWT library to validate this
		if token != "Bearer my-secret-pro-token" {
			http.Error(w, "401 Unauthorized - Valid token required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}