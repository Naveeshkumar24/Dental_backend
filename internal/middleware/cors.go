package middleware

import "net/http"

// CORS handles Cross-Origin Resource Sharing
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow all origins (SAFE for APIs; restrict in production if needed)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allowed HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allowed headers
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Authorization, Content-Type, X-Requested-With",
		)

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
