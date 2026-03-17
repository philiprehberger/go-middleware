package middleware

import "net/http"

// BearerAuth validates the Authorization: Bearer <token> header using the provided validate function.
// If valid, the request is passed downstream.
// If invalid, responds with 401 Unauthorized.
func BearerAuth(validate func(token string) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if len(header) < 8 || header[:7] != "Bearer " {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token := header[7:]
			if err := validate(token); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
