package server

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	burst    int
	interval time.Duration
	// ipHeader is the header that carries the real client IP, or "" when no
	// header is trusted and the socket address is used. See
	// resolveTrustedIPHeader.
	ipHeader string
	// now is the clock, replaced in tests so refill can be measured instead of
	// slept through.
	now func() time.Time
}

type visitor struct {
	// tokens is fractional: a bucket that only credited whole intervals let a
	// learner answering every second earn nothing at all.
	tokens    float64
	lastCheck time.Time
}

func newRateLimiter(rate, burst int, interval time.Duration, ipHeader string) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		interval: interval,
		ipHeader: ipHeader,
		now:      time.Now,
	}
	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-30 * time.Minute)
		for ip, v := range rl.visitors {
			if v.lastCheck.Before(cutoff) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, ok := rl.visitors[key]
	now := rl.now()

	if !ok {
		rl.visitors[key] = &visitor{tokens: float64(rl.burst - 1), lastCheck: now}
		return true
	}

	// Refill continuously. The old arithmetic added int(elapsed/interval)
	// tokens, which is zero for every request less than an interval after the
	// previous one, so a normal answering pace never refilled the bucket and
	// the twentieth answer of a diagnostic was rejected.
	if elapsed := now.Sub(v.lastCheck).Seconds(); elapsed > 0 {
		v.tokens += elapsed * float64(rl.rate) / rl.interval.Seconds()
		if max := float64(rl.burst); v.tokens > max {
			v.tokens = max
		}
	}
	v.lastCheck = now
	if v.tokens < 1 {
		return false
	}
	v.tokens--
	return true
}

// key names the bucket for a request: the signed-in learner when there is one,
// otherwise the client address.
//
// IP keying alone rationed an account by everyone who shares its address, which
// on a school network or a carrier NAT is a whole class. The identity is only
// visible to this middleware when the auth wrapper runs first, so the practice
// routes are registered that way.
func (rl *rateLimiter) key(r *http.Request) string {
	if id, _ := r.Context().Value(authStudentKey{}).(string); id != "" {
		return "student:" + id
	}
	return "ip:" + clientIP(r, rl.ipHeader)
}

func (rl *rateLimiter) middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(rl.key(r)) {
			w.Header().Set("Retry-After", "2")
			http.Error(w, `{"error":"too many answers at once — retrying in a moment"}`, 429)
			return
		}
		next(w, r)
	}
}

// resolveTrustedIPHeader decides which request header, if any, carries the true
// client IP.
//
// A header is only trustworthy when the edge that sets it also overwrites any
// value the client supplied. Otherwise it is just another spoofable input, and
// both the rate limiter and the guest AI quota are bypassed by setting it: a
// caller sends `Fly-Client-IP: 1.2.3.4` and is limited under a key of their own
// choosing, or, worse, spends somebody else's AI budget.
//
// Fly's edge sets Fly-Client-IP and overwrites it, and Fly puts FLY_APP_NAME in
// the machine environment, so that header is trusted exactly when FLY_APP_NAME
// is present. A deployment behind a different proxy sets AI_TRUSTED_IP_HEADER to
// the header that proxy controls. With neither, the socket address is the only
// untamperable signal, which is the correct answer for a directly exposed
// process.
func resolveTrustedIPHeader() string {
	if h := strings.TrimSpace(os.Getenv("AI_TRUSTED_IP_HEADER")); h != "" {
		return h
	}
	if os.Getenv("FLY_APP_NAME") != "" {
		return "Fly-Client-IP"
	}
	return ""
}

// clientIP returns the caller's address for rate limiting and guest AI quotas.
//
// It reads trustedHeader only, never X-Forwarded-For: that header's first entry
// is client-controlled, so a parse that trusted it was a rate-limit bypass for
// anyone willing to set a header. When no header is trusted, or the trusted one
// is absent, the socket address is used, so a request that never traversed the
// edge is keyed by the address that actually connected.
func clientIP(r *http.Request, trustedHeader string) string {
	if trustedHeader != "" {
		if ip := strings.TrimSpace(r.Header.Get(trustedHeader)); ip != "" {
			return ip
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
