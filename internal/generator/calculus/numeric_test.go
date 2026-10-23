package calculus

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestSequenceLimitPrintsNoStrayCoefficient(t *testing.T) {
	g := &sequenceLimitGen{}
	for i := 0; i < 100; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		if strings.Contains(p.Question, "1n^2") {
			t.Fatalf("question prints a coefficient of 1: %q", p.Question)
		}
		if !strings.Contains(p.Question, "n^2") {
			t.Fatalf("question lost its squared term: %q", p.Question)
		}
	}
}
