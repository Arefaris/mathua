package server

// HTTP-level checks for translated lesson bodies: a request that asks for Russian gets the
// Russian lesson when there is one and the English lesson when there is not, and the steps find
// their worked example in whichever body came back.

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/storage"
)

// lessonFixture has three concepts: "a" with a translated lesson (and a figure), "b" with an
// English-only lesson, and "o" whose English lesson carries the OpenStax attribution and so has
// no translation.
func lessonFixture(t *testing.T) http.Handler {
	t.Helper()
	mt := concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}
	dag, err := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "Alpha", Domain: "d", GradingType: "numeric", Prerequisites: []string{}, MasteryThreshold: mt},
		{ID: "b", Label: "Beta", Domain: "d", GradingType: "numeric", Prerequisites: []string{}, MasteryThreshold: mt},
		{ID: "o", Label: "Omega", Domain: "d", GradingType: "numeric", Prerequisites: []string{}, MasteryThreshold: mt},
	})
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	w := func(rel, body string) { writeTestFile(t, filepath.Join(dir, "lessons", filepath.FromSlash(rel)), body) }
	w("lessons.json", `[
		{"concept_id":"a","source":"teaching/a.md"},
		{"concept_id":"b","source":"teaching/b.md"},
		{"concept_id":"o","source":"teaching/o.md"}
	]`)
	w("teaching/a.md", "> Hand-authored lesson for the Mathua project.\n\n# Adding\n\n## First Idea\n\nEnglish idea $x+1$.\n\n![Figure](svg/fig-a.svg)\n")
	w("teaching/ru/a.md", "> Hand-authored lesson for the Mathua project.\n\n# Сложение\n\n## Первая мысль\n\nРусская мысль $x+1$.\n\n![Рисунок](svg/fig-a.svg)\n")
	w("kp/a.json", `[{"label":"One","section":"First Idea","subgoals":["x","y","z"]}]`)
	w("kp/ru/a.json", `[{"label":"Раз","section":"Первая мысль","subgoals":["икс","игрек","зет"]}]`)
	w("teaching/b.md", "# Bee\n\n## Bee Idea\n\nEnglish bee.\n")
	w("kp/b.json", `[{"label":"Bee","section":"Bee Idea","subgoals":["s","t","u"]}]`)
	w("kp/ru/b.json", `[{"label":"Пчела","section":"Идея пчелы","subgoals":["с","т","у"]}]`)
	w("teaching/o.md", "> Content sourced from [OpenStax Prealgebra](https://openstax.org/books/prealgebra) — CC BY 4.0\n\n# Open\n\n## Open Idea\n\nEnglish open.\n")
	w("kp/o.json", `[{"label":"Open","section":"Open Idea","subgoals":["s","t","u"]}]`)
	ll, err := lessons.Load(filepath.Join(dir, "lessons"))
	if err != nil {
		t.Fatal(err)
	}
	reg := generator.NewRegistry()
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	return localeServer(t, engine.New(store, dag, reg, ll, nil), store)
}

type bodyResp struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Locale string `json:"locale"`
	Assets []struct {
		ID      string `json:"id"`
		Section string `json:"section"`
	} `json:"assets"`
}

func lessonBody(t *testing.T, h http.Handler, title, acceptLanguage string) bodyResp {
	t.Helper()
	var out bodyResp
	rec := getAs(h, "/api/lessons/body?title="+url.QueryEscape(title), acceptLanguage)
	decode(t, rec, &out)
	if v := rec.Header().Get("Vary"); !strings.Contains(v, "Accept-Language") {
		t.Errorf("Vary = %q: a cache must not serve one language's lesson to another", v)
	}
	return out
}

func TestLessonBodyFollowsAcceptLanguage(t *testing.T) {
	h := lessonFixture(t)

	for _, header := range []string{"ru", "ru-RU", "ru,en;q=0.8"} {
		got := lessonBody(t, h, "Adding", header)
		if got.Locale != "ru" || !strings.Contains(got.Body, "Русская мысль") || got.Title != "Сложение" {
			t.Errorf("Accept-Language %q: %+v", header, got)
		}
		// The attribution survives the translation byte for byte.
		if !strings.Contains(got.Body, "> Hand-authored lesson for the Mathua project.") {
			t.Errorf("Accept-Language %q: attribution line missing from %q", header, got.Body)
		}
		if len(got.Assets) != 1 || got.Assets[0].ID != "fig-a.svg" || got.Assets[0].Section != "Первая мысль" {
			t.Errorf("Accept-Language %q: assets = %+v", header, got.Assets)
		}
	}
	for _, header := range []string{"", "en", "en-US,en;q=0.9", "de", "uk"} {
		got := lessonBody(t, h, "Adding", header)
		if got.Locale != "en" || !strings.Contains(got.Body, "English idea") || got.Title != "Adding" {
			t.Errorf("Accept-Language %q: want the English lesson, got %+v", header, got)
		}
		if len(got.Assets) != 1 || got.Assets[0].Section != "First Idea" {
			t.Errorf("Accept-Language %q: assets = %+v", header, got.Assets)
		}
	}
}

// A lesson with no translation stays English under ru, and so does one from OpenStax, which is
// never translated.
func TestUntranslatedLessonStaysEnglishOverHTTP(t *testing.T) {
	h := lessonFixture(t)
	for title, want := range map[string]string{"Bee": "English bee", "Open": "English open"} {
		got := lessonBody(t, h, title, "ru")
		if got.Locale != "en" || !strings.Contains(got.Body, want) || got.Title != title {
			t.Errorf("%s under ru = %+v", title, got)
		}
	}
	// An unknown title is still a 404, in any language.
	if rec := getAs(h, "/api/lessons/body?title=Nope", "ru"); rec.Code != 404 {
		t.Errorf("unknown title: status %d", rec.Code)
	}
}

type kpsResp struct {
	KPs []struct {
		Label         string   `json:"label"`
		Section       string   `json:"section"`
		Subgoals      []string `json:"subgoals"`
		WorkedExample string   `json:"worked_example"`
		Assets        []struct {
			ID      string `json:"id"`
			Section string `json:"section"`
		} `json:"assets"`
	} `json:"kps"`
}

// The step's section finds its worked example in the body that is served: the Russian heading in
// the Russian lesson, the English heading in the English one, and a lesson with no translation
// keeps looking in its English body with the English heading while the step text is Russian.
func TestStepSectionsResolveInTheServedBody(t *testing.T) {
	h := lessonFixture(t)

	var ru, en kpsResp
	decode(t, getAs(h, "/api/lessons/a/kp", "ru"), &ru)
	decode(t, getAs(h, "/api/lessons/a/kp", ""), &en)
	if len(ru.KPs) != 1 || len(en.KPs) != 1 {
		t.Fatalf("steps: ru %d, en %d", len(ru.KPs), len(en.KPs))
	}
	k := ru.KPs[0]
	if k.Label != "Раз" || k.Section != "Первая мысль" || k.Subgoals[0] != "икс" {
		t.Errorf("ru step = %+v", k)
	}
	if !strings.HasPrefix(k.WorkedExample, "## Первая мысль") || !strings.Contains(k.WorkedExample, "Русская мысль $x+1$") {
		t.Errorf("ru worked example = %q", k.WorkedExample)
	}
	if len(k.Assets) != 1 || k.Assets[0].ID != "fig-a.svg" {
		t.Errorf("ru figure for the step = %+v", k.Assets)
	}
	e := en.KPs[0]
	if e.Label != "One" || e.Section != "First Idea" || !strings.Contains(e.WorkedExample, "English idea") {
		t.Errorf("en step = %+v", e)
	}
	if len(e.Assets) != 1 || e.Assets[0].ID != "fig-a.svg" {
		t.Errorf("en figure for the step = %+v", e.Assets)
	}

	// b: the step is translated, the lesson is not. The worked example is the English section,
	// which a Russian heading could not have found.
	var b kpsResp
	decode(t, getAs(h, "/api/lessons/b/kp", "ru"), &b)
	if b.KPs[0].Label != "Пчела" || !strings.Contains(b.KPs[0].WorkedExample, "English bee") {
		t.Errorf("untranslated lesson, translated step = %+v", b.KPs[0])
	}
	// o: nothing is translated.
	var o kpsResp
	decode(t, getAs(h, "/api/lessons/o/kp", "ru"), &o)
	if o.KPs[0].Label != "Open" || !strings.Contains(o.KPs[0].WorkedExample, "English open") {
		t.Errorf("OpenStax lesson = %+v", o.KPs[0])
	}
}
