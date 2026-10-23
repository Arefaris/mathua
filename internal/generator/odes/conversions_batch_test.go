package odes

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var convertedBatchIDs = []string{
	"ode.ivp.evaluate",
	"ode.characteristic.sum",
	"ode.equilibrium.nonzero",
	"ode.wronskian",
	"ode.numerical.steps",
	"ode.laplace.value",
	"ode.pde.discriminant",
}

func testRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func TestODESBatchConversions_CognitiveScaffolding(t *testing.T) {
	reg := testRegistry(t)
	for _, id := range convertedBatchIDs {
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

func TestODESBatchConversions_DifficultyScaling(t *testing.T) {
	reg := testRegistry(t)
	for _, id := range convertedBatchIDs {
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

func TestODESBatchConversions_ParameterTransparency(t *testing.T) {
	reg := testRegistry(t)
	checked := 0
	for _, id := range convertedBatchIDs {
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

func TestODESBatchConversions_ZeroOxfordComma(t *testing.T) {
	reg := testRegistry(t)
	oxfordComma := regexp.MustCompile(`,\s+and\s+`)
	for _, id := range convertedBatchIDs {
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

// TestODEsBatchAnswersAreIntegers pins the battery to numeric answers. All seven
// concepts are graded by the numeric grader, so a decimal or a symbolic string
// would mark a correct learner wrong.
func TestODEsBatchAnswersAreIntegers(t *testing.T) {
	reg := testRegistry(t)
	for _, id := range convertedBatchIDs {
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

// TestODEsBatchDoesNotPrintItsAnswer covers the read-off class. The equilibrium
// generator used to draw the capacity a and ask for a, so every one of 500
// sampled draws printed its own answer in the stem and the question tested
// copying rather than solving. It now writes the model as y' = a y - b y^2 and
// asks for a/b. A share above this threshold means a question can be answered
// without doing the work, which is the defect this guards.
func TestODEsBatchDoesNotPrintItsAnswer(t *testing.T) {
	reg := testRegistry(t)
	const draws = 500
	const maxShare = 0.6
	for _, id := range convertedBatchIDs {
		printed := 0
		for i := 0; i < draws; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       rand.Int63(),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			answer := p.Answer
			if strings.HasPrefix(answer, "-") {
				answer = answer[1:]
			}
			pattern := regexp.MustCompile(`(^|\D)` + regexp.QuoteMeta(answer) + `(\D|$)`)
			if pattern.MatchString(p.Question) {
				printed++
			}
		}
		if share := float64(printed) / draws; share > maxShare {
			t.Fatalf("%s: %d of %d draws print the answer in the stem (%.0f%%), limit %.0f%%",
				id, printed, draws, share*100, maxShare*100)
		}
	}
}
