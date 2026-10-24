package all

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// TestNoConceptIsRegisteredTwice fails when two generator packages claim the same concept.
//
// Every generator package ignores the error from Register, which is why a collision produced no
// symptom: the first registration won, the second implementation was unreachable and its tests
// still passed because they built the type directly. The six basic combinatorics counts were
// registered by both `discrete` and `combinatorics` on 23 October 2026, and only a probe of the
// registry showed which one a learner actually got.
func TestNoConceptIsRegisteredTwice(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)

	dups := reg.Duplicates()
	if len(dups) == 0 {
		return
	}
	t.Fatalf("two generators claim the same concept: %s\nKeep one registration and delete the other implementation, or rename the concept.",
		strings.Join(dups, ", "))
}
