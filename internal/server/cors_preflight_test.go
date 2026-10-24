package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCORS_PreflightAllowsTheRequestedHeaders pins the failure that would have
// taken the whole cross-origin API down: the browser only honours a preflight
// when every header the request carries is listed in Access-Control-Allow-
// Headers, and a missing one fails silently, the fetch never leaves the page.
// X-Request-ID is now attached to every client request, so the preflight must
// allow it.
func TestCORS_PreflightAllowsTheRequestedHeaders(t *testing.T) {
	handler := cors(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("a preflight must not reach the handler")
	})
	req := httptest.NewRequest(http.MethodOptions, "/api/config", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "authorization,content-type,x-request-id")
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d, want 204", rec.Code)
	}
	allowed := strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers"))
	for _, want := range []string{"authorization", "content-type", "x-request-id"} {
		if !strings.Contains(allowed, want) {
			t.Errorf("Access-Control-Allow-Headers = %q, missing %q", allowed, want)
		}
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("Allow-Origin = %q, want the calling origin", got)
	}
}

// TestCORS_PreflightDoesNotReflectHeadersForAForeignOrigin keeps the reflection
// tied to the origin allowlist. Reflecting for anyone would let a page we do not
// trust ask the browser for a preflight it should not get.
func TestCORS_PreflightDoesNotReflectHeadersForAForeignOrigin(t *testing.T) {
	handler := cors(func(w http.ResponseWriter, r *http.Request) {})
	req := httptest.NewRequest(http.MethodOptions, "/api/config", nil)
	req.Header.Set("Origin", "https://not-mathua.example")
	req.Header.Set("Access-Control-Request-Headers", "x-not-ours")
	rec := httptest.NewRecorder()

	handler(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin = %q for a foreign origin, want empty", got)
	}
	if strings.Contains(strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers")), "x-not-ours") {
		t.Error("a foreign origin's requested headers were reflected")
	}
}
