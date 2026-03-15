// Package middleware provides composable HTTP middleware for net/http.
//
// Middleware functions wrap http.Handler to add cross-cutting concerns like
// logging, panic recovery, CORS, security headers, timeouts, compression,
// and ETag generation. Use Chain to compose multiple middleware together.
package middleware

import "net/http"

// Middleware is a function that wraps an http.Handler to add behavior.
type Middleware = func(http.Handler) http.Handler

// Chain composes multiple middleware into a single Middleware.
// Middleware are applied in the order given: the first middleware in the list
// wraps the outermost layer, so it executes first on the way in and last on
// the way out.
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
