package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// corsConfig holds the CORS configuration.
type corsConfig struct {
	allowOrigins    []string
	allowMethods    []string
	allowHeaders    []string
	allowCredentials bool
	maxAge          int
}

// CORSOption configures the CORS middleware.
type CORSOption func(*corsConfig)

// AllowOrigins sets the allowed origins. An empty list or a list containing "*"
// allows all origins.
func AllowOrigins(origins ...string) CORSOption {
	return func(c *corsConfig) {
		c.allowOrigins = origins
	}
}

// AllowMethods sets the allowed HTTP methods for CORS requests.
func AllowMethods(methods ...string) CORSOption {
	return func(c *corsConfig) {
		c.allowMethods = methods
	}
}

// AllowHeaders sets the allowed request headers for CORS requests.
func AllowHeaders(headers ...string) CORSOption {
	return func(c *corsConfig) {
		c.allowHeaders = headers
	}
}

// AllowCredentials enables the Access-Control-Allow-Credentials header.
func AllowCredentials() CORSOption {
	return func(c *corsConfig) {
		c.allowCredentials = true
	}
}

// MaxAge sets the maximum time (in seconds) that preflight results can be cached.
func MaxAge(seconds int) CORSOption {
	return func(c *corsConfig) {
		c.maxAge = seconds
	}
}

// CORS returns middleware that handles Cross-Origin Resource Sharing.
// By default it allows all origins and the methods GET, POST, HEAD, PUT, DELETE, PATCH.
func CORS(opts ...CORSOption) Middleware {
	cfg := &corsConfig{
		allowOrigins: []string{"*"},
		allowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			allowed := isOriginAllowed(origin, cfg.allowOrigins)
			if !allowed {
				next.ServeHTTP(w, r)
				return
			}

			if len(cfg.allowOrigins) == 1 && cfg.allowOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			}

			if cfg.allowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.allowMethods, ", "))
				if len(cfg.allowHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.allowHeaders, ", "))
				}
				if cfg.maxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.maxAge))
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}
	return false
}
