package middleware

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

// Recover returns middleware that recovers from panics in downstream handlers.
// When a panic occurs, it logs the panic value and stack trace to stderr and
// responds with a 500 Internal Server Error.
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Fprintf(os.Stderr, "panic: %v\n%s\n", err, debug.Stack())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
