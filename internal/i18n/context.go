package i18n

import "context"

type localeKey struct{}

// WithLocale returns a context that carries loc.
func WithLocale(ctx context.Context, loc Locale) context.Context {
	return context.WithValue(ctx, localeKey{}, loc)
}

// FromContext returns the locale stored by WithLocale, or Default when the
// context carries none, so code that never saw the middleware still serves English.
func FromContext(ctx context.Context) Locale {
	if loc, ok := ctx.Value(localeKey{}).(Locale); ok && loc.Valid() {
		return loc
	}
	return Default
}
