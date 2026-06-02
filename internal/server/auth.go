package server

import (
	"crypto/subtle"
	"net/http"
)

// authHeader is the header carrying the access token for HTTP/SSE transports.
const authHeader = "X-MCP-AUTH"

// withAuth wraps h so that, when token is non-empty, every request must present
// a matching X-MCP-AUTH header. An empty token disables auth (handler passes
// through unchanged). Comparison is constant-time.
func withAuth(token string, h http.Handler) http.Handler {
	if token == "" {
		return h
	}
	want := []byte(token)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := []byte(r.Header.Get(authHeader))
		if subtle.ConstantTimeCompare(got, want) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
