package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
}

func TestWithAuth_Disabled(t *testing.T) {
	h := withAuth("", okHandler())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/mcp", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("disabled auth should pass, got %d", rec.Code)
	}
}

func TestWithAuth_Enabled(t *testing.T) {
	h := withAuth("secret", okHandler())

	cases := []struct {
		name   string
		header string
		want   int
	}{
		{"valid token", "secret", http.StatusOK},
		{"wrong token", "nope", http.StatusUnauthorized},
		{"missing token", "", http.StatusUnauthorized},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/mcp", nil)
			if c.header != "" {
				req.Header.Set(authHeader, c.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != c.want {
				t.Errorf("%s: got %d, want %d", c.name, rec.Code, c.want)
			}
		})
	}
}
