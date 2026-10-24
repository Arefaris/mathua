package server

// HTTP-level checks that a request's Accept-Language reaches the response: the concept
// labels, course text and lesson steps come back in Russian when asked for, and in English
// otherwise or for any item the overlay does not cover.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/i18n"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/storage"
)

// localeServer routes through LocaleMiddleware exactly as cmd/mathua does.
func localeServer(t *testing.T, eng *engine.Engine, store storage.Repository) http.Handler {
	t.Helper()
	s := New(eng, store, nil)
	mux := http.NewServeMux()
	s.Register(mux)
	return LocaleMiddleware(mux)
}

func getAs(h http.Handler, path, acceptLanguage string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", path, nil)
	if acceptLanguage != "" {
		req.Header.Set("Accept-Language", acceptLanguage)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}
}

type graphResp struct {
	Count int `json:"count"`
	Nodes []struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"nodes"`
}

func graphLabels(t *testing.T, h http.Handler, acceptLanguage string) (map[string]string, int) {
	t.Helper()
	var g graphResp
	decode(t, getAs(h, "/api/graph", acceptLanguage), &g)
	out := make(map[string]string, len(g.Nodes))
	for _, n := range g.Nodes {
		out[n.ID] = n.Label
	}
	return out, g.Count
}

// The real corpus and the shipped overlay: a Russian request gets Russian labels for all
// of it, and a request without the header gets the English originals.
func TestGraphLabelsFollowAcceptLanguageWithShippedData(t *testing.T) {
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load corpus: %v", err)
	}
	if err := dag.LoadOverlays("../../data/i18n"); err != nil {
		t.Fatalf("load overlays: %v", err)
	}
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	h := localeServer(t, engine.New(store, dag, generator.NewRegistry(), nil, nil), store)

	en, n := graphLabels(t, h, "")
	if n != 1250 {
		t.Fatalf("graph has %d concepts, want 1250", n)
	}
	if got := en["arith.add.single"]; got != "Single-digit addition" {
		t.Errorf("no header: %q, want the English label", got)
	}
	for _, header := range []string{"en", "en-US,en;q=0.9", "de", "fr,de;q=0.5"} {
		got, _ := graphLabels(t, h, header)
		if got["arith.add.single"] != "Single-digit addition" {
			t.Errorf("Accept-Language %q: %q, want English", header, got["arith.add.single"])
		}
	}

	for _, header := range []string{"ru", "ru-RU", "ru,en;q=0.8"} {
		ru, _ := graphLabels(t, h, header)
		if got := ru["arith.add.single"]; got != "Сложение однозначных чисел" {
			t.Errorf("Accept-Language %q: %q, want the Russian label", header, got)
		}
		for id, label := range ru {
			if label == en[id] && label != "" {
				// A handful of labels are identical in both languages only if they are
				// formula-only; none are today.
				t.Errorf("Accept-Language %q: %s was served untranslated (%q)", header, id, label)
				break
			}
		}
	}
}

// localeGen records the locale its generator is called with.
type localeGen struct {
	mu   sync.Mutex
	seen []i18n.Locale
}

func (g *localeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	g.mu.Lock()
	g.seen = append(g.seen, ctx.Locale)
	n := len(g.seen)
	g.mu.Unlock()
	return generator.Problem{Question: "q" + string(rune('a'+n%26)), Answer: "1", Explanation: "e"}
}

func (g *localeGen) last() i18n.Locale {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.seen) == 0 {
		return "none"
	}
	return g.seen[len(g.seen)-1]
}

func writeTestFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// fixture builds a two-concept corpus where only "a" is translated, one course and one
// lesson with one translated step.
func fixture(t *testing.T) (http.Handler, *localeGen) {
	t.Helper()
	mt := concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}
	dag, err := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "Alpha", Domain: "d", GradingType: "numeric", Prerequisites: []string{}, MasteryThreshold: mt},
		{ID: "b", Label: "Beta", Domain: "d", GradingType: "numeric", Prerequisites: []string{"a"}, MasteryThreshold: mt},
	})
	if err != nil {
		t.Fatal(err)
	}
	dag.SetLabels(i18n.RU, map[string]string{"a": "Альфа"})

	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "courses.json"), `[{"id":"c1","name":"Course","grade":"1","description":"About it","targets":["b"]}]`)
	writeTestFile(t, filepath.Join(dir, "i18n", "ru", "courses.json"), `{"c1":{"name":"Курс","description":"Об этом"}}`)
	planner, err := planning.Load(filepath.Join(dir, "courses.json"), dag)
	if err != nil {
		t.Fatal(err)
	}
	if err := planner.LoadOverlays(filepath.Join(dir, "i18n")); err != nil {
		t.Fatal(err)
	}

	writeTestFile(t, filepath.Join(dir, "lessons", "lessons.json"), `[{"concept_id":"a","source":"a.md"}]`)
	writeTestFile(t, filepath.Join(dir, "lessons", "a.md"), "# A\n\n## Section One\n\nbody")
	writeTestFile(t, filepath.Join(dir, "lessons", "kp", "a.json"), `[{"label":"One","section":"Section One","subgoals":["x","y","z"]}]`)
	writeTestFile(t, filepath.Join(dir, "lessons", "kp", "ru", "a.json"), `[{"label":"Раз","section":"Раздел один","subgoals":["икс","игрек","зет"]}]`)
	ll, err := lessons.Load(filepath.Join(dir, "lessons"))
	if err != nil {
		t.Fatal(err)
	}

	gen := &localeGen{}
	reg := generator.NewRegistry()
	reg.Register("a", gen)
	reg.Register("b", gen)
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	return localeServer(t, engine.New(store, dag, reg, ll, planner), store), gen
}

func TestPerItemFallbackOverHTTP(t *testing.T) {
	h, _ := fixture(t)

	ru, _ := graphLabels(t, h, "ru")
	if ru["a"] != "Альфа" {
		t.Errorf("a = %q", ru["a"])
	}
	if ru["b"] != "Beta" {
		t.Errorf("b has no translation and must stay English, got %q", ru["b"])
	}
	en, _ := graphLabels(t, h, "")
	if en["a"] != "Alpha" || en["b"] != "Beta" {
		t.Errorf("English = %v", en)
	}

	var cur struct {
		Topics []struct {
			ConceptID string `json:"conceptId"`
			Title     string `json:"title"`
		} `json:"topics"`
	}
	decode(t, getAs(h, "/api/curriculum/domains/d", "ru"), &cur)
	titles := map[string]string{}
	for _, tp := range cur.Topics {
		titles[tp.ConceptID] = tp.Title
	}
	if titles["a"] != "Альфа" || titles["b"] != "Beta" {
		t.Errorf("curriculum titles = %v", titles)
	}

	// b is locked behind a: the prerequisite's title is localized, b's own is not translated.
	var ready struct {
		Prerequisites []struct {
			ConceptID string `json:"conceptId"`
			Title     string `json:"title"`
		} `json:"prerequisites"`
		Missing []struct {
			ID    string `json:"id"`
			Label string `json:"label"`
		} `json:"missing"`
	}
	decode(t, getAs(h, "/api/lessons/b/readiness", "ru"), &ready)
	if len(ready.Prerequisites) != 1 || ready.Prerequisites[0].Title != "Альфа" {
		t.Errorf("eligibility prerequisites = %+v", ready.Prerequisites)
	}
	if len(ready.Missing) != 1 || ready.Missing[0].Label != "Альфа" {
		t.Errorf("missing = %+v", ready.Missing)
	}
	decode(t, getAs(h, "/api/lessons/b/readiness", ""), &ready)
	if ready.Prerequisites[0].Title != "Alpha" {
		t.Errorf("English prerequisite title = %q", ready.Prerequisites[0].Title)
	}
}

func TestLessonStepsAndCoursesOverHTTP(t *testing.T) {
	h, _ := fixture(t)

	type kpResp struct {
		KPs []struct {
			Label    string   `json:"label"`
			Section  string   `json:"section"`
			Subgoals []string `json:"subgoals"`
		} `json:"kps"`
	}
	var ru, en kpResp
	rec := getAs(h, "/api/lessons/a/kp", "ru")
	decode(t, rec, &ru)
	decode(t, getAs(h, "/api/lessons/a/kp", ""), &en)
	if len(ru.KPs) != 1 || len(en.KPs) != 1 {
		t.Fatalf("steps: ru %d, en %d; body %s", len(ru.KPs), len(en.KPs), rec.Body.String())
	}
	if ru.KPs[0].Label != "Раз" || ru.KPs[0].Section != "Раздел один" || ru.KPs[0].Subgoals[0] != "икс" {
		t.Errorf("ru step = %+v", ru.KPs[0])
	}
	if en.KPs[0].Label != "One" || en.KPs[0].Section != "Section One" || en.KPs[0].Subgoals[0] != "x" {
		t.Errorf("en step = %+v", en.KPs[0])
	}

	type coursesResp struct {
		Courses []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"courses"`
	}
	var cru, cen coursesResp
	decode(t, getAs(h, "/api/courses", "ru"), &cru)
	decode(t, getAs(h, "/api/courses", ""), &cen)
	if cru.Courses[0].Name != "Курс" || cru.Courses[0].Description != "Об этом" {
		t.Errorf("ru course = %+v", cru.Courses[0])
	}
	// The same shared course must still be English for the next request.
	if cen.Courses[0].Name != "Course" || cen.Courses[0].Description != "About it" {
		t.Errorf("en course after a ru request = %+v", cen.Courses[0])
	}
	var detail struct {
		Name string `json:"name"`
	}
	decode(t, getAs(h, "/api/courses/c1", "ru"), &detail)
	if detail.Name != "Курс" {
		t.Errorf("course detail name = %q", detail.Name)
	}
	decode(t, getAs(h, "/api/courses/c1", ""), &detail)
	if detail.Name != "Course" {
		t.Errorf("course detail name (en) = %q", detail.Name)
	}
}

// The locale reaches the generator on the practice path, and defaults to English without
// the header. Generators write English today, so this is the plumbing the next round uses.
func TestPracticeGenerationReceivesTheRequestLocale(t *testing.T) {
	h, gen := fixture(t)

	for header, want := range map[string]i18n.Locale{"ru": i18n.RU, "ru-RU,ru;q=0.9": i18n.RU, "": i18n.EN, "en": i18n.EN, "de": i18n.EN} {
		rec := getAs(h, "/api/lessons/a/practice?count=1", header)
		if rec.Code != 200 {
			t.Fatalf("%q: status %d: %s", header, rec.Code, rec.Body.String())
		}
		if got := gen.last(); got != want {
			t.Errorf("Accept-Language %q: generator saw %q, want %q", header, got, want)
		}
	}
}

// The goal diagnostic names the concept it is asking about in the request's language, and
// passes that language to the generator that wrote the question.
func TestGoalDiagnosticNamesTheConceptInTheRequestLocale(t *testing.T) {
	h, gen := fixture(t)

	start := func(header string) string {
		t.Helper()
		req := httptest.NewRequest("POST", "/api/goal/diagnostic", strings.NewReader(`{"name":"kid","concept_ids":["a"]}`))
		req.Header.Set("Content-Type", "application/json")
		if header != "" {
			req.Header.Set("Accept-Language", header)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		var out struct {
			ConceptID   string `json:"concept_id"`
			ConceptName string `json:"concept_name"`
		}
		decode(t, rec, &out)
		if out.ConceptID != "a" {
			t.Fatalf("first concept = %q, want a", out.ConceptID)
		}
		return out.ConceptName
	}

	if got := start("ru"); got != "Альфа" {
		t.Errorf("ru name = %q", got)
	}
	if got := gen.last(); got != i18n.RU {
		t.Errorf("generator locale = %q", got)
	}
	if got := start(""); got != "Alpha" {
		t.Errorf("en name = %q", got)
	}
}
