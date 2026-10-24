package numtheory

// Guards for the eight numeric generators this batch made parametric.
//
// Their own suite checks registration, scaffolding and that the answer matches
// Facts["answer"]. None of that would notice a fact the question never prints,
// which is the shape prob.dist.order_stats shipped (its key held n = 3 while
// the stem said "let n observations") and the shape the parametric-k reports
// were. The checks here compare what a learner can read with what the grader
// expects.

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var batchNumericIDs = []string{
	"nt.adv.jacobi",
	"nt.adv.discrete_log",
	"nt.adv.modular_sqrt",
	"nt.adv.two_squares",
	"nt.adv.pell",
	"nt.adv.pisano",
	"nt.adv.carmichael_function",
	"nt.adv.prime_gap",
}

func batchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

// TestNumTheoryBatchAnswersAreIntegers pins the grading contract: the eight
// concepts are registered as numeric, so anything else would never compare
// equal to a typed answer.
func TestNumTheoryBatchAnswersAreIntegers(t *testing.T) {
	reg := batchRegistry(t)
	for _, id := range batchNumericIDs {
		for i := 0; i < 100; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 10, Seed: int64(i),
			})
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			if _, err := strconv.Atoi(p.Answer); err != nil {
				t.Fatalf("%s answers %q, which is not an integer", id, p.Answer)
			}
			if !regexp.MustCompile(`Subgoal 1`).MatchString(p.Explanation) ||
				!regexp.MustCompile(`Subgoal 3`).MatchString(p.Explanation) {
				t.Fatalf("%s explanation is missing its subgoals: %q", id, p.Explanation)
			}
		}
	}
}

// TestNumTheoryBatchFactsAppearInTheStem requires every number the answer is
// derived from to be readable in the question. A modulus, exponent or bound
// that exists only in Facts is unrecoverable by the learner.
func TestNumTheoryBatchFactsAppearInTheStem(t *testing.T) {
	reg := batchRegistry(t)
	checked := 0
	for _, id := range batchNumericIDs {
		for i := 0; i < 100; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 10, Seed: int64(i) * 7,
			})
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			for name, value := range p.Facts {
				if name == "answer" {
					continue
				}
				num, err := strconv.Atoi(value)
				if err != nil {
					continue // a word like kind = "gap"
				}
				if num < 0 {
					num = -num
				}
				// Digit-bounded so 5 does not match inside 25.
				pattern := regexp.MustCompile(`(^|\D)` + strconv.Itoa(num) + `(\D|$)`)
				if !pattern.MatchString(p.Question) {
					t.Fatalf("%s: names %s = %s in Facts but the question does not print it: %q",
						id, name, value, p.Question)
				}
				checked++
			}
		}
	}
	if checked < 1000 {
		t.Fatalf("only %d facts checked, so the sweep is not seeing the generators", checked)
	}
}
