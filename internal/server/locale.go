package server

import (
	"net/http"

	"github.com/chuma-beep/mathua/internal/i18n"
)

// LocaleMiddleware resolves the request's language from Accept-Language and stores it in
// the request context, where handlers read it with localeOf.
//
// The client sends the header on every API call and the server trusts nothing else: there
// is no cookie and no ?lang= parameter, so a response is a function of the URL, the body and
// this one header. That is also why Vary is set. A cache or proxy between the browser and the
// server that ignored it would serve one learner's Russian labels to the next learner's
// English request.
func LocaleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc := i18n.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		w.Header().Add("Vary", "Accept-Language")
		next.ServeHTTP(w, r.WithContext(i18n.WithLocale(r.Context(), loc)))
	})
}

// localeOf is the language to answer the request in; English when the request did not pass
// through LocaleMiddleware.
func localeOf(r *http.Request) i18n.Locale {
	return i18n.FromContext(r.Context())
}
