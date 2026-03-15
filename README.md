# go-middleware

Composable HTTP middleware collection for Go's `net/http`. No framework required.

## Installation

```sh
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
		mw.Logger(logger),
		mw.Recover(),
		mw.CORS(),
		mw.SecureHeaders(),
		mw.Timeout(10*time.Second),
		mw.Compress(),
		mw.ETag(),
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

## License

MIT
