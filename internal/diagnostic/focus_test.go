package diagnostic

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
)

// pathOf returns the whole prerequisite closure in the DAG's own order.
func pathOf(t *testing.T, d *concepts.DAG) []*concepts.Concept {
	t.Helper()
	return d.Order()
}

// Report 117 was a learner who chose linear algebra, analysis and abstract
// algebra and was asked "what is the opposite of 1" first. The closure of that
// selection is 246 concepts and its highest-influence node is a root, so the
// first probe was arithmetic. The focus does not drop the foundations: it makes
// the chosen concepts outrank them while both remain equally uncertain.
func TestFocusProbesTheChosenConceptFirst(t *testing.T) {
	d := testDAG(t)
	eng := NewEngine(d, mustRegistry(t))
	path := pathOf(t, d)

	plain := eng.StartWithPath(path)
	if _, cid, err := eng.NextQuestion(plain); err != nil {
		t.Fatalf("next question: %v", err)
	} else if cid == "c" {
		t.Fatalf("without a focus the first probe was the chosen concept, so this test proves nothing")
	}

	focused := eng.StartWithPathFocus(path, []string{"c"})
	if _, cid, err := eng.NextQuestion(focused); err != nil {
		t.Fatalf("next question: %v", err)
	} else if cid != "c" {
		t.Fatalf("first probe with focus c = %q, want %q", cid, "c")
	}
}

func TestFocusStillAsksTheFoundations(t *testing.T) {
	d := testDAG(t)
	eng := NewEngine(d, mustRegistry(t))
	path := pathOf(t, d)

	s := eng.StartWithPathFocus(path, []string{"d"})
	asked := map[string]bool{}
	for {
		_, cid, err := eng.NextQuestion(s)
		if err != nil {
			t.Fatalf("next question: %v", err)
		}
		if cid == "" {
			break
		}
		asked[cid] = true
		eng.SettleCurrent(s)
	}
	for _, want := range []string{"a", "b", "c", "d"} {
		if !asked[want] {
			t.Errorf("concept %q was never asked; focus must not drop the path", want)
		}
	}
}

func TestFocusSurvivesARestart(t *testing.T) {
	d := testDAG(t)
	eng := NewEngine(d, mustRegistry(t))
	path := pathOf(t, d)

	s := eng.StartWithPathFocus(path, []string{"c"})
	s.ID = "focus-restart"
	data, err := s.Encode()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	restored, err := eng.Decode(data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, cid, err := eng.NextQuestion(restored); err != nil {
		t.Fatalf("next question: %v", err)
	} else if cid != "c" {
		t.Fatalf("first probe after a restart = %q, want %q", cid, "c")
	}
}
