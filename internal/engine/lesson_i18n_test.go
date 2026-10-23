package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/i18n"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/storage"
)

// The lesson that rides along with a question is in the request's language when that lesson has
// a translation, and English when it has none.
func TestQuestionCarriesTheLessonInTheRequestLanguage(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("lessons.json", `[{"concept_id":"a","source":"teaching/a.md"}]`)
	write("teaching/a.md", "# Adding\n\n## First\n\nEnglish text.\n")
	write("teaching/ru/a.md", "# Сложение\n\n## Первое\n\nРусский текст.\n")
	write("kp/a.json", `[{"label":"One","section":"First","subgoals":["x","y","z"]}]`)
	write("kp/ru/a.json", `[{"label":"Раз","section":"Первое","subgoals":["икс","игрек","зет"]}]`)
	ll, err := lessons.Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	mt := concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}
	d, err := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "Concept A", Domain: "d", GradingType: "numeric", Prerequisites: []string{}, MasteryThreshold: mt},
	})
	if err != nil {
		t.Fatal(err)
	}
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &localeProbe{})
	e := New(store, d, reg, ll, nil)
	st, _ := e.CreateStudent("kid")

	for _, tc := range []struct {
		loc  i18n.Locale
		want string
	}{{i18n.RU, "Русский текст"}, {i18n.EN, "English text"}, {i18n.UK, "English text"}} {
		sess, _ := e.repo.CreateSession(st.ID)
		q, err := e.NextQuestionIn(tc.loc, sess.ID, st.ID)
		if err != nil || q == nil || q.Lesson == nil {
			t.Fatalf("%s: question %+v, err %v", tc.loc, q, err)
		}
		if !strings.Contains(q.Lesson.Body, tc.want) {
			t.Errorf("%s: lesson body = %q, want it to contain %q", tc.loc, q.Lesson.Body, tc.want)
		}
		// Reading the pending question back follows the language of the reading request.
		cur, _ := e.GetCurrentQuestionIn(tc.loc, sess.ID, st.ID)
		if cur == nil || cur.Lesson == nil || !strings.Contains(cur.Lesson.Body, tc.want) {
			t.Errorf("%s: current question lesson = %+v", tc.loc, cur)
		}
	}
}
