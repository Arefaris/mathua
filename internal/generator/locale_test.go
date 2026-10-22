package generator

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/i18n"
)

type localeCapture struct{ seen []i18n.Locale }

func (g *localeCapture) Generate(ctx GeneratorContext) Problem {
	g.seen = append(g.seen, ctx.Locale)
	return Problem{Question: "q", Answer: "1"}
}

// A generator can compare against i18n.RU without guarding for the empty string, and an
// unsupported value is treated as English.
func TestRegistryNormalizesTheLocale(t *testing.T) {
	cases := []struct {
		name string
		run  func(r *Registry)
		want i18n.Locale
	}{
		{"zero context", func(r *Registry) { r.GenerateContext("c", GeneratorContext{Difficulty: 0.5}) }, i18n.EN},
		{"unsupported", func(r *Registry) { r.GenerateContext("c", GeneratorContext{Difficulty: 0.5, Locale: "xx"}) }, i18n.EN},
		{"ru context", func(r *Registry) { r.GenerateContext("c", GeneratorContext{Difficulty: 0.5, Locale: i18n.RU}) }, i18n.RU},
		{"Generate", func(r *Registry) { r.Generate("c", 0.5) }, i18n.EN},
		{"GenerateIn", func(r *Registry) { r.GenerateIn(i18n.RU, "c", 0.5) }, i18n.RU},
		{"batch ru", func(r *Registry) { r.BatchGenerateContext("c", 1, GeneratorContext{Difficulty: 0.5, Locale: i18n.RU}) }, i18n.RU},
		{"batch zero", func(r *Registry) { r.BatchGenerateContext("c", 1, GeneratorContext{Difficulty: 0.5}) }, i18n.EN},
	}
	for _, c := range cases {
		r := NewRegistry()
		g := &localeCapture{}
		if err := r.Register("c", g); err != nil {
			t.Fatal(err)
		}
		c.run(r)
		if len(g.seen) == 0 {
			t.Errorf("%s: the generator was not called", c.name)
		}
		// A batch calls the generator repeatedly to find distinct questions.
		for _, got := range g.seen {
			if got != c.want {
				t.Errorf("%s: generator saw %q, want %q", c.name, got, c.want)
			}
		}
	}
}
