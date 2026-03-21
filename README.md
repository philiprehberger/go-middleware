# go-middleware

[![CI](https://github.com/philiprehberger/go-middleware/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-middleware/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-middleware.svg)](https://pkg.go.dev/github.com/philiprehberger/go-middleware) [![License](https://img.shields.io/github/license/philiprehberger/go-middleware)](LICENSE)

Composable HTTP middleware collection for Go's `net/http`. No framework required

## Installation

```bash
go get github.com/philiprehberger/go-middleware
```

## Usage

```go
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	mw "github.com/philiprehberger/go-middleware"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	stack := mw.Chain(
		mw.RequestID,
		mw.Logger(logger),
		mw.Recover(),
		mw.CORS(),
		mw.SecureHeaders(),
		mw.BearerAuth(func(token string) error { return nil }),
		mw.Timeout(10*time.Second),
		mw.Compress(),
		mw.ETag(),
		mw.Metrics(func(method, path string, status int, duration time.Duration) {
			logger.Info("metrics", "method", method, "path", path, "status", status, "duration", duration)
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":8080", stack(mux))
}
```

### Available Middleware

**Chain** — Compose multiple middleware into a single wrapper.

```go
handler := mw.Chain(m1, m2, m3)(finalHandler)
```

**Logger** — Log each request with method, path, status, and duration.

```go
mw.Logger(slog.Default())
```

**Recover** — Catch panics and return 500 Internal Server Error.

```go
mw.Recover()
```

**CORS** — Handle Cross-Origin Resource Sharing with functional options.

```go
mw.CORS(
    mw.AllowOrigins("https://example.com"),
    mw.AllowMethods("GET", "POST"),
    mw.AllowHeaders("Authorization"),
    mw.AllowCredentials(),
    mw.MaxAge(3600),
)
```

**SecureHeaders** — Set common security headers (X-Content-Type-Options, X-Frame-Options, Referrer-Policy, X-XSS-Protection).

```go
mw.SecureHeaders()
```

**Timeout** — Enforce a request timeout, returning 503 if exceeded.

```go
mw.Timeout(10 * time.Second)
```

**Compress** — Gzip-compress responses when the client supports it.

```go
mw.Compress()
```

**ETag** — Generate ETag headers and return 304 Not Modified when appropriate.

```go
mw.ETag()
```

**Request ID** — Generate or propagate an `X-Request-ID` header. If the incoming request already has one, it is preserved; otherwise a new random ID is generated. The ID is also stored in the request context.

```go
mw.RequestID
// Extract from context inside a handler:
id := mw.RequestIDFromContext(r.Context())
```

**Bearer Auth** — Validate `Authorization: Bearer <token>` headers using a custom validation function. Returns 401 if missing or invalid.

```go
mw.BearerAuth(func(token string) error {
    if token != "expected" {
        return errors.New("invalid token")
    }
    return nil
})
```

**Metrics** — Call a function after each request with method, path, status code, and duration.

```go
mw.Metrics(func(method, path string, status int, duration time.Duration) {
    log.Printf("%s %s -> %d (%s)", method, path, status, duration)
})
```

### Middleware Ordering

```go
import mw "github.com/philiprehberger/go-middleware"

// Recommended order: recovery → logging → security → business logic
handler := mw.Chain(
    mw.Recover(),
    mw.RequestID,
    mw.Logger(slog.Default()),
    mw.CORS(mw.AllowOrigins("*")),
    mw.BearerAuth(validateToken),
)(yourHandler)
```

### Preset Chains

```go
import mw "github.com/philiprehberger/go-middleware"

// Common API preset
apiChain := mw.Chain(
    mw.Recover(),
    mw.RequestID,
    mw.Logger(slog.Default()),
    mw.Timeout(30 * time.Second),
    mw.CORS(mw.AllowOrigins("*")),
)
http.Handle("/api/", apiChain(apiHandler))
```

## API

| Function | Signature | Description |
|---|---|---|
| `Chain` | `Chain(middlewares ...Middleware) Middleware` | Compose middleware in order |
| `Logger` | `Logger(logger *slog.Logger) Middleware` | Request logging |
| `Recover` | `Recover() Middleware` | Panic recovery |
| `CORS` | `CORS(opts ...CORSOption) Middleware` | CORS handling |
| `AllowOrigins` | `AllowOrigins(origins ...string) CORSOption` | Set allowed origins |
| `AllowMethods` | `AllowMethods(methods ...string) CORSOption` | Set allowed methods |
| `AllowHeaders` | `AllowHeaders(headers ...string) CORSOption` | Set allowed headers |
| `AllowCredentials` | `AllowCredentials() CORSOption` | Enable credentials |
| `MaxAge` | `MaxAge(seconds int) CORSOption` | Set preflight cache duration |
| `SecureHeaders` | `SecureHeaders() Middleware` | Security headers |
| `Timeout` | `Timeout(d time.Duration) Middleware` | Request timeout |
| `Compress` | `Compress() Middleware` | Gzip compression |
| `ETag` | `ETag() Middleware` | ETag generation |
| `RequestID` | `RequestID(next http.Handler) http.Handler` | Request ID generation/propagation |
| `RequestIDFromContext` | `RequestIDFromContext(ctx context.Context) string` | Extract request ID from context |
| `BearerAuth` | `BearerAuth(validate func(string) error) func(http.Handler) http.Handler` | Bearer token authentication |
| `Metrics` | `Metrics(onRequest func(string, string, int, time.Duration)) func(http.Handler) http.Handler` | Request metrics collection |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
