package calculus

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var convertedNumericIDs = []string{
	"calc.seq.limit_value",
	"calc.seq.recursive_value",
	"calc.seq.arithmetic_sum",
	"calc.multi.mass_lamina",
	"calc.integral.substitution_value",
	"calc.integral.mass_density",
	"calc.deriv.higher_order_value",
}

func testCalculusRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func TestCalculusBatchConversions_CognitiveScaffolding(t *testing.T) {
	reg := testCalculusRegistry(t)
	for _, id := range convertedNumericIDs {
		for i := 0; i < 30; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 9.0,
				Seed:       int64(10000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			if !strings.Contains(p.Explanation, "Subgoal 1: State governing formula") {
				t.Fatalf("%s: missing Subgoal 1 in explanation: %s", id, p.Explanation)
			}
			if !strings.Contains(p.Explanation, "Subgoal 2: Substitute parameters") {
				t.Fatalf("%s: missing Subgoal 2 in explanation: %s", id, p.Explanation)
			}
			if !strings.Contains(p.Explanation, "Subgoal 3: Evaluate and simplify") {
				t.Fatalf("%s: missing Subgoal 3 in explanation: %s", id, p.Explanation)
			}
		}
	}
}

func TestCalculusBatchConversions_DifficultyScaling(t *testing.T) {
	reg := testCalculusRegistry(t)
	for _, id := range convertedNumericIDs {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("generator not found: %s: %v", id, err)
		}
		for i := 0; i < 50; i++ {
			d := rand.Float64()
			prob := gen.Generate(generator.GeneratorContext{Difficulty: d})
			if prob.Question == "" || prob.Answer == "" || prob.Explanation == "" {
				t.Fatalf("%s: empty problem output at difficulty %.2f", id, d)
			}
			if prob.Facts == nil || prob.Facts["answer"] != prob.Answer {
				t.Fatalf("%s: fact answer missing or mismatch: got %q, want %q", id, prob.Facts["answer"], prob.Answer)
			}
			if _, err := strconv.Atoi(prob.Answer); err != nil {
				t.Fatalf("%s: answer %q is not an integer", id, prob.Answer)
			}
		}
	}
}

func TestCalculusBatchConversions_ParameterTransparency(t *testing.T) {
	reg := testCalculusRegistry(t)
	checked := 0
	for _, id := range convertedNumericIDs {
		for i := 0; i < 60; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 9.0,
				Seed:       int64(20000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			for key, raw := range p.Facts {
				if key == "answer" {
					continue
				}
				value, err := strconv.Atoi(raw)
				if err != nil {
					continue
				}
				absVal := value
				if absVal < 0 {
					absVal = -absVal
				}
				pattern := regexp.MustCompile(`(^|\D)` + regexp.QuoteMeta(strconv.Itoa(absVal)) + `(\D|$)`)
				if !pattern.MatchString(p.Question) {
					t.Fatalf("%s: names parameter %s = %d, but question does not print it: %s",
						id, key, value, p.Question)
				}
				checked++
			}
		}
	}
	if checked < 400 {
		t.Fatalf("only %d facts checked", checked)
	}
}

func TestCalculusBatchConversions_ZeroOxfordComma(t *testing.T) {
	reg := testCalculusRegistry(t)
	oxfordComma := regexp.MustCompile(`,\s+and\s+`)
	for _, id := range convertedNumericIDs {
		for i := 0; i < 20; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 9.0,
				Seed:       int64(30000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			if oxfordComma.MatchString(p.Explanation) {
				t.Fatalf("%s: explanation contains Oxford comma: %s", id, p.Explanation)
			}
			if oxfordComma.MatchString(p.Question) {
				t.Fatalf("%s: question contains Oxford comma: %s", id, p.Question)
			}
		}
	}
}

// TestCalculusBatchAnswersAreIntegers pins the battery to numeric answers. All
// seven concepts are graded by the numeric grader, so a decimal or a symbolic
// string would mark a correct learner wrong.
func TestCalculusBatchAnswersAreIntegers(t *testing.T) {
	reg := testCalculusRegistry(t)
	for _, id := range convertedNumericIDs {
		for i := 0; i < 100; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i%10) / 9.0,
				Seed:       int64(40000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			if _, err := strconv.Atoi(p.Answer); err != nil {
				t.Fatalf("%s: answer %q is not an integer: %s", id, p.Answer, p.Question)
			}
		}
	}
}

// TestCalculusBatchStaysVaried guards against a generator that degenerates to one
// question or one answer. Two batch conversions on 23 October 2026 published
// constant answers in most modes, which a package suite that only checked
// registration and scaffolding could not see.
func TestCalculusBatchStaysVaried(t *testing.T) {
	reg := testCalculusRegistry(t)
	const draws = 300
	const minDistinct = 6
	for _, id := range convertedNumericIDs {
		questions := map[string]bool{}
		answers := map[string]bool{}
		for i := 0; i < draws; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       rand.Int63(),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			questions[p.Question] = true
			answers[p.Answer] = true
		}
		if len(questions) < minDistinct || len(answers) < minDistinct {
			t.Fatalf("%s: only %d distinct questions and %d distinct answers in %d draws, want at least %d of each",
				id, len(questions), len(answers), draws, minDistinct)
		}
	}
}
