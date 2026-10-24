package complex

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

// batchGradingTypes states, per converted concept, how the corpus grades it.
// The answers here are complex strings ("49/50+7/50i"), fractions and signed
// integers, so "the generator produced an answer" is not the same as "the
// grader can read it".
var batchGradingTypes = map[string]grader.GradingType{
	"complex.basics.parts":          grader.GradingNumeric,
	"complex.basics.powers":         grader.GradingComplex,
	"complex.basics.reciprocal":     grader.GradingComplex,
	"complex.basics.geometry":       grader.GradingNumeric,
	"complex.basics.arg_arithmetic": grader.GradingNumeric,
	"complex.basics.equations":      grader.GradingComplex,
	"complex.adv.geometric_series":  grader.GradingComplex,
	"complex.adv.mobius_fixed":      grader.GradingComplex,
	"complex.adv.residue_simple":    grader.GradingComplex,
	"complex.adv.cross_ratio":       grader.GradingComplex,
}

// TestComplexBatchAnswersGradeAsThemselves runs each published answer back
// through the router. An answer form the grader cannot parse marks a correct
// learner wrong, and no test in this package could see it before.
func TestComplexBatchAnswersGradeAsThemselves(t *testing.T) {
	reg := batchRegistry(t)
	router := grader.NewRouter()
	checked := 0
	for id, gt := range batchGradingTypes {
		for i := 0; i < 100; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       rand.Int63(),
			})
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			if res := router.Grade(gt, p.Answer, p.Answer); !res.Correct {
				t.Fatalf("%s: answer %q does not grade as itself: %s", id, p.Answer, res.Feedback)
			}
			checked++
		}
	}
	if checked < 900 {
		t.Fatalf("checked only %d answers, expected 1000", checked)
	}
}

// TestComplexBatchEquationsAreNotTrivial rejects the zero right-hand side. The
// root was drawn without excluding 0 + 0i, so one draw in a few hundred asked
// the learner to solve (2 + 3i)z = 0, whose answer is 0 + 0i whatever the
// coefficient is.
func TestComplexBatchEquationsAreNotTrivial(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 200; i++ {
		p, err := reg.GenerateContext("complex.basics.equations", generator.GeneratorContext{
			Difficulty: rand.Float64(),
			Seed:       rand.Int63(),
		})
		if err != nil {
			t.Fatal(err)
		}
		c, _ := strconv.Atoi(p.Facts["c"])
		d, _ := strconv.Atoi(p.Facts["d"])
		if c == 0 && d == 0 {
			t.Fatalf("equation has a zero right-hand side: %s (answer %s)", p.Question, p.Answer)
		}
	}
}

// TestComplexBatchSeriesAnswersVary pins the widened ratio range. The series
// answer is a function of the ratio's denominator alone, so the old range of
// 2..8 meant seven possible answers no matter how many times the concept was
// served.
func TestComplexBatchSeriesAnswersVary(t *testing.T) {
	reg := batchRegistry(t)
	answers := map[string]bool{}
	for i := 0; i < 200; i++ {
		p, err := reg.GenerateContext("complex.adv.geometric_series", generator.GeneratorContext{
			Difficulty: 1,
			Seed:       rand.Int63(),
		})
		if err != nil {
			t.Fatal(err)
		}
		answers[p.Answer] = true
	}
	if len(answers) < 10 {
		t.Fatalf("series serves only %d distinct answers at full difficulty, want at least 10", len(answers))
	}
}
