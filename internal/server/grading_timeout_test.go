package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestGradingTimeout_ReturnsJSONBeforeTheClientGivesUp is the D3 regression:
// a grade slower than the deadline must produce a 503 the browser can read and
// retry, rather than a request the browser aborts while the server keeps
// working on it.
func TestGradingTimeout_ReturnsJSONBeforeTheClientGivesUp(t *testing.T) {
	slow := func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
		case <-r.Context().Done():
			return
		}
		io.WriteString(w, `{"late":true}`)
	}
	srv := httptest.NewServer(gradingTimeout(30*time.Millisecond, slow))
	defer srv.Close()

	start := time.Now()
	res, err := http.Post(srv.URL, "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	elapsed := time.Since(start)

	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", res.StatusCode)
	}
	if !strings.Contains(string(body), "grading took too long") {
		t.Errorf("body = %q, want the retryable grading message", body)
	}
	if elapsed > time.Second {
		t.Errorf("timeout took %s, deadline was 30ms", elapsed)
	}
}

// TestGradingTimeout_LeavesFastHandlersAlone guards against the wrapper
// truncating a normal answer.
func TestGradingTimeout_LeavesFastHandlersAlone(t *testing.T) {
	fast := func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, `{"ok":true}`) }
	srv := httptest.NewServer(gradingTimeout(time.Second, fast))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 || !strings.Contains(string(body), `"ok":true`) {
		t.Errorf("fast handler = %d %q", res.StatusCode, body)
	}
}
