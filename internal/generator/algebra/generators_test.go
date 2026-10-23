package algebra

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	gr := grader.NewRouter()
	for i := 0; i < 100; i++ {
		d := rand.Float64()
		p := gen.Generate(generator.GeneratorContext{Difficulty: d})
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
		}
		res := gr.Grade(grader.GradingNumeric, p.Answer, p.Answer)
		if !res.Correct {
			res2 := gr.Grade(grader.GradingMultipleChoice, p.Answer, p.Answer)
			if !res2.Correct {
				t.Errorf("self-grade failed for %q (answer=%q)", p.Question, p.Answer)
			}
		}
	}
}

func TestSlopeGen(t *testing.T)             { fuzzGen(t, &slopeGen{}) }
func TestSlopeInterceptGen(t *testing.T)    { fuzzGen(t, &slopeInterceptGen{}) }
func TestMultiStepEqGen(t *testing.T)       { fuzzGen(t, &multiStepEqGen{}) }
func TestVarsBothSidesGen(t *testing.T)     { fuzzGen(t, &varsBothSidesGen{}) }
func TestVarsBothSidesWordGen(t *testing.T) { fuzzGen(t, &varsBothSidesWordGen{}) }
func TestIneqWordGen(t *testing.T)          { fuzzGen(t, &ineqWordGen{}) }
func TestStdFormGen(t *testing.T)           { fuzzGen(t, &stdFormGen{}) }
func TestParallelPerpGen(t *testing.T)      { fuzzGen(t, &parallelPerpGen{}) }

func TestDifficultyScaling(t *testing.T) {
	gen := &multiStepEqGen{}
	// Low difficulty should produce smaller numbers
	low := gen.Generate(generator.GeneratorContext{Difficulty: 0.1})
	high := gen.Generate(generator.GeneratorContext{Difficulty: 1.0})
	if low.Question == "" || high.Question == "" {
		t.Fatal("empty questions")
	}
	t.Logf("low=%.2f: %s", 0.1, low.Question)
	t.Logf("high=%.2f: %s", 1.0, high.Question)
}

func TestFuzz(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)
	for _, id := range reg.Concepts() {
		gen, _ := reg.Get(id)
		t.Run(id, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				d := rand.Float64()
				p := gen.Generate(generator.GeneratorContext{Difficulty: d})
				if p.Question == "" || p.Answer == "" || p.Explanation == "" {
					t.Fatalf("empty field at difficulty=%.2f for %s: q=%q a=%q", d, id, p.Question, p.Answer)
				}
			}
		})
	}
}
func TestInjectivityGen(t *testing.T)    { fuzzGen(t, &injectivityGen{}) }
func TestPowerFuncGen(t *testing.T)      { fuzzGen(t, &powerFuncGen{}) }
func TestIrrationalFuncGen(t *testing.T) { fuzzGen(t, &irrationalFuncGen{}) }

// Nine reports on 23 October 2026 said the same thing: the question asked for k
// and the equation had no k in it. The answer was unrelated to the printed
// quadratic as well. This pins the equation a learner reads against the answer.
func TestQuadParametricAsksForAParameterItPrints(t *testing.T) {
	for i := 0; i < 100; i++ {
		p := (&quadParametricGen{}).Generate(generator.GeneratorContext{})
		if !strings.Contains(p.Question, "kx") {
			t.Fatalf("question does not print the parameter it asks for: %q", p.Question)
		}
		k, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("answer %q is not an integer", p.Answer)
		}
		c, err := strconv.Atoi(p.Facts["const"])
		if err != nil {
			t.Fatalf("const fact %q is not an integer", p.Facts["const"])
		}
		if k*k != 4*c {
			t.Fatalf("k=%d gives discriminant %d for x^2 + kx + %d, want zero", k, k*k-4*c, c)
		}
		if p.Facts["k"] != p.Answer || p.Facts["coef"] != p.Answer {
			t.Fatalf("facts disagree with the answer: %v", p.Facts)
		}
		if !strings.Contains(p.Question, strconv.Itoa(c)) {
			t.Fatalf("constant %d is not printed in %q", c, p.Question)
		}
	}
}
