package engine

import (
	"sync"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/i18n"
	"github.com/chuma-beep/mathua/internal/storage"
)

type localeProbe struct {
	mu   sync.Mutex
	seen []i18n.Locale
}

func (g *localeProbe) Generate(ctx generator.GeneratorContext) generator.Problem {
	g.mu.Lock()
	g.seen = append(g.seen, ctx.Locale)
	g.mu.Unlock()
	return generator.Problem{Question: "q", Answer: "1", Explanation: "e"}
}

func localeEngine(t *testing.T) (*Engine, *localeProbe) {
	t.Helper()
	mt := concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}
	d, err := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "Concept A", Domain: "d", GradingType: "numeric", Prerequisites: []string{}, MasteryThreshold: mt},
	})
	if err != nil {
		t.Fatal(err)
	}
	d.SetLabels(i18n.RU, map[string]string{"a": "Тема А"})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	probe := &localeProbe{}
	reg := generator.NewRegistry()
	reg.Register("a", probe)
	return New(store, d, reg, nil, nil), probe
}

func TestNextQuestionInLocalizesTheNameAndTellsTheGenerator(t *testing.T) {
	e, probe := localeEngine(t)
	st, _ := e.CreateStudent("kid")
	sess, _ := e.repo.CreateSession(st.ID)

	q, err := e.NextQuestionIn(i18n.RU, sess.ID, st.ID)
	if err != nil || q == nil {
		t.Fatalf("question: %v %v", q, err)
	}
	if q.ConceptName != "Тема А" {
		t.Errorf("ConceptName = %q", q.ConceptName)
	}
	if got := probe.seen[len(probe.seen)-1]; got != i18n.RU {
		t.Errorf("generator saw %q", got)
	}

	// The same pending question read back in English is renamed from its id, not served
	// from the name that was stored when it was generated.
	cur, err := e.GetCurrentQuestionIn(i18n.EN, sess.ID, st.ID)
	if err != nil || cur == nil {
		t.Fatalf("current: %v %v", cur, err)
	}
	if cur.ConceptName != "Concept A" {
		t.Errorf("current name in en = %q", cur.ConceptName)
	}
	cur, _ = e.GetCurrentQuestionIn(i18n.RU, sess.ID, st.ID)
	if cur.ConceptName != "Тема А" {
		t.Errorf("current name in ru = %q", cur.ConceptName)
	}

	// The old entry points are English and tell the generator so.
	sess2, _ := e.repo.CreateSession(st.ID)
	q2, _ := e.NextQuestion(sess2.ID, st.ID)
	if q2 == nil || q2.ConceptName != "Concept A" {
		t.Fatalf("default question = %+v", q2)
	}
	if got := probe.seen[len(probe.seen)-1]; got != i18n.EN {
		t.Errorf("default generator locale = %q", got)
	}
}

func TestDiagnosticQuestionAndReportAreLocalized(t *testing.T) {
	e, probe := localeEngine(t)
	// One concept is enough for a diagnostic run.
	sess := e.diag.Start()
	prob, cid, err := e.NextDiagnosticQuestionIn(i18n.RU, sess)
	if err != nil || prob == nil || cid != "a" {
		t.Fatalf("diagnostic question: %v %v %q", prob, err, cid)
	}
	if got := probe.seen[len(probe.seen)-1]; got != i18n.RU {
		t.Errorf("diagnostic generator saw %q", got)
	}
	rep := e.DiagnosticReportIn(i18n.RU, sess)
	for _, labels := range rep.GapsByDomain {
		for _, l := range labels {
			if l != "Тема А" {
				t.Errorf("gap label = %q", l)
			}
		}
	}
	if len(rep.GapsByDomain["d"]) == 0 {
		t.Error("expected the unprobed concept in the gaps")
	}
}
