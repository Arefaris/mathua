package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/i18n"
)

func TestLocaleMiddleware(t *testing.T) {
	var seen i18n.Locale
	h := LocaleMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = localeOf(r)
	}))
	cases := []struct {
		header string
		want   i18n.Locale
	}{
		{"", i18n.EN},
		{"en", i18n.EN},
		{"ru", i18n.RU},
		{"ru-RU", i18n.RU},
		{"ru,en;q=0.8", i18n.RU},
		{"de,fr;q=0.5", i18n.EN},
	}
	for _, c := range cases {
		req := httptest.NewRequest("GET", "/api/graph", nil)
		if c.header != "" {
			req.Header.Set("Accept-Language", c.header)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if seen != c.want {
			t.Errorf("Accept-Language %q: handler saw %q, want %q", c.header, seen, c.want)
		}
		if got := rec.Header().Get("Vary"); got != "Accept-Language" {
			t.Errorf("Accept-Language %q: Vary = %q, want Accept-Language", c.header, got)
		}
	}
}

// The static frontend is the same bytes in every language, so it must not be told to vary.
func TestLocaleMiddlewareLeavesStaticAssetsCacheable(t *testing.T) {
	h := LocaleMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/_next/static/app.js", nil)
	req.Header.Set("Accept-Language", "ru")
	h.ServeHTTP(rec, req)
	if got := rec.Header().Get("Vary"); got != "" {
		t.Errorf("Vary = %q on a static asset", got)
	}
}

// Without the middleware the server answers in English rather than failing.
func TestLocaleOfDefaultsToEnglish(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/graph", nil)
	req.Header.Set("Accept-Language", "ru")
	if got := localeOf(req); got != i18n.EN {
		t.Fatalf("localeOf without middleware = %q, want en", got)
	}
}

func TestCORSAllowsAcceptLanguage(t *testing.T) {
	h := cors(func(w http.ResponseWriter, r *http.Request) {})
	req := httptest.NewRequest("OPTIONS", "/api/graph", nil)
	req.Header.Set("Origin", "https://example.test")
	rec := httptest.NewRecorder()
	h(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Headers"); !containsToken(got, "Accept-Language") {
		t.Fatalf("Access-Control-Allow-Headers = %q, want it to list Accept-Language", got)
	}
}

func containsToken(list, tok string) bool {
	for _, t := range strings.Split(list, ",") {
		if strings.TrimSpace(t) == tok {
			return true
		}
	}
	return false
}
