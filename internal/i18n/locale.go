// Package i18n carries the learner's language through the server.
//
// English is the source language and the default: every lookup that can miss
// (a label, a course description, a lesson step) falls back to it item by item,
// so a translation can ship in pieces and an untranslated item is never blank.
package i18n

import (
	"sort"
	"strconv"
	"strings"
)

// Locale is a language the server can answer in. It is the bare language subtag
// ("ru"), not a full BCP 47 tag, because nothing here varies by region.
type Locale string

const (
	EN Locale = "en"
	RU Locale = "ru"
	UK Locale = "uk"
)

// Default is the locale of every untranslated item and of any request that does
// not name a supported language.
const Default = EN

// Supported lists the locales the server has translations for, English first.
var Supported = []Locale{EN, RU, UK}

// aliases maps tags that people and some browsers send for a language but that are not its
// ISO 639-1 code. "ua" is the country code of Ukraine, often sent for Ukrainian by mistake.
var aliases = map[string]Locale{"ua": UK}

// Overlays are the locales that ship translation files next to the English
// originals. English is the original, so it never has an overlay.
func Overlays() []Locale {
	out := make([]Locale, 0, len(Supported))
	for _, l := range Supported {
		if l != Default {
			out = append(out, l)
		}
	}
	return out
}

// Valid reports whether l is one of the supported locales.
func (l Locale) Valid() bool {
	for _, s := range Supported {
		if s == l {
			return true
		}
	}
	return false
}

// String implements fmt.Stringer.
func (l Locale) String() string { return string(l) }

// Parse maps a single language tag to a supported locale. "ru", "RU" and "ru-RU" are
// Russian; "uk", "uk-UA" and "ua" are Ukrainian; anything else, including the empty string,
// is Default.
func Parse(tag string) Locale {
	if l, ok := match(tag); ok {
		return l
	}
	return Default
}

// match resolves one language tag to a supported locale.
func match(tag string) (Locale, bool) {
	tag = strings.ToLower(strings.TrimSpace(tag))
	if i := strings.IndexAny(tag, "-_"); i >= 0 {
		tag = tag[:i]
	}
	for _, s := range Supported {
		if string(s) == tag {
			return s, true
		}
	}
	if l, ok := aliases[tag]; ok {
		return l, true
	}
	return "", false
}

// ParseAcceptLanguage picks the locale for an Accept-Language header value.
//
// It is deliberately lenient, because the header is client-supplied and a bad
// one must never fail a request: a malformed entry is skipped, a missing or
// wholly unsupported header is Default. Among the supported languages the
// highest q-value wins and the header's own order breaks a tie, so
// "ru,en;q=0.8" is Russian and "en,ru;q=0.5" is English. A language with q=0 is
// "not acceptable" and is ignored.
func ParseAcceptLanguage(header string) Locale {
	type choice struct {
		loc   Locale
		q     float64
		order int
	}
	var picks []choice
	for i, part := range strings.Split(header, ",") {
		tag, params, _ := strings.Cut(part, ";")
		loc, ok := match(tag)
		if !ok {
			continue
		}
		q := 1.0
		for _, p := range strings.Split(params, ";") {
			k, v, found := strings.Cut(strings.TrimSpace(p), "=")
			if !found || !strings.EqualFold(strings.TrimSpace(k), "q") {
				continue
			}
			f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err != nil {
				q = 0
				continue
			}
			q = f
		}
		if q <= 0 {
			continue
		}
		picks = append(picks, choice{loc, q, i})
	}
	if len(picks) == 0 {
		return Default
	}
	sort.SliceStable(picks, func(a, b int) bool {
		if picks[a].q != picks[b].q {
			return picks[a].q > picks[b].q
		}
		return picks[a].order < picks[b].order
	})
	return picks[0].loc
}
