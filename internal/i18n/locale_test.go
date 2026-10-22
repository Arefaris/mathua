package i18n

import (
	"context"
	"testing"
)

func TestParseAcceptLanguage(t *testing.T) {
	cases := []struct {
		header string
		want   Locale
	}{
		{"", EN},
		{"en", EN},
		{"ru", RU},
		{"RU", RU},
		{"ru-RU", RU},
		{"ru_RU", RU},
		{"ru,en;q=0.8", RU},
		{"en,ru;q=0.5", EN},
		{"en;q=0.4, ru;q=0.9", RU},
		{"fr, ru;q=0.5", RU},
		{"fr-CA, de;q=0.8", EN},
		{"*", EN},
		{"ru;q=0", EN},
		{"ru;q=0, en;q=0.5", EN},
		{"ru;q=oops", EN},
		{"  ru  ;  q=0.7 ,en;q=0.7", RU},
		{"ru;q=0.5,en;q=0.5", RU},
		{";;;,,,", EN},
		{"russian", EN},
	}
	for _, c := range cases {
		if got := ParseAcceptLanguage(c.header); got != c.want {
			t.Errorf("ParseAcceptLanguage(%q) = %q, want %q", c.header, got, c.want)
		}
	}
}

func TestParse(t *testing.T) {
	for tag, want := range map[string]Locale{"ru": RU, "ru-RU": RU, "en-GB": EN, "": EN, "de": EN} {
		if got := Parse(tag); got != want {
			t.Errorf("Parse(%q) = %q, want %q", tag, got, want)
		}
	}
}

func TestContextRoundTrip(t *testing.T) {
	if got := FromContext(context.Background()); got != Default {
		t.Fatalf("bare context = %q, want default %q", got, Default)
	}
	if got := FromContext(WithLocale(context.Background(), RU)); got != RU {
		t.Fatalf("WithLocale(ru) = %q", got)
	}
	// A value that is not a supported locale must not leak through.
	if got := FromContext(WithLocale(context.Background(), Locale("xx"))); got != Default {
		t.Fatalf("unsupported locale = %q, want default", got)
	}
}

func TestOverlaysExcludeDefault(t *testing.T) {
	for _, l := range Overlays() {
		if l == Default {
			t.Fatalf("Overlays() contains the default locale %q", l)
		}
	}
	if len(Overlays()) != 1 || Overlays()[0] != RU {
		t.Fatalf("Overlays() = %v, want [ru]", Overlays())
	}
}
