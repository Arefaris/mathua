package generator_test

import (
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/all"
)

// A multi-word answer cannot be typed into the math field. The field strips
// spaces, so "line segment" arrives as "linesegment" or worse, and the grader
// compares that to the authored phrase. Five reports on 23 October 2026 were
// exactly this: three asked a yes/no question, one asked for a geometric term
// and one asked for a sentence.
//
// So: a concept graded multiple_choice whose answer contains a space must serve
// an option list, and the answer must be one of the options. Numeric answers
// are exempt, because they stay typeable and most of the corpus relies on that
// fallback today.
func TestMultipleChoiceProseAnswersServeOptions(t *testing.T) {
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load concepts: %v", err)
	}
	reg := generator.NewRegistry()
	all.Register(reg)

	checked, offenders := 0, 0
	for id, c := range dag.Concepts() {
		if c == nil || c.GradingType != "multiple_choice" || !learnerDomains[c.Domain] {
			continue
		}
		gen, err := reg.Get(id)
		if err != nil {
			continue
		}
		prose := false
		bad := ""
		for i := 0; i < 60; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: float64(i%5) / 4, Seed: int64(i)})
			if !looksLikeWords(p.Answer) {
				continue
			}
			prose = true
			options := generator.ChoiceOptions(p.Question)
			if len(options) < 2 {
				bad = "answer is prose and the question carries no options"
				break
			}
			found := false
			for _, text := range options {
				if text == p.Answer {
					found = true
					break
				}
			}
			if !found {
				bad = "the answer is not one of the options"
				break
			}
		}
		if !prose {
			continue
		}
		checked++
		if bad != "" {
			offenders++
			sample := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: 1})
			t.Errorf("%s: %s\n  question: %s\n  answer: %s", id, bad, sample.Question, sample.Answer)
		}
	}
	if checked == 0 {
		t.Fatal("no multiple_choice concept produced a prose answer: the sweep checked nothing")
	}
	t.Logf("checked %d prose-answer concepts, %d offenders", checked, offenders)
}

// mathWords are words that belong inside an expression rather than in the
// answer. "pi/2 - 1" is typeable, so it is not prose; "two real" is.
var latexCommand = regexp.MustCompile(`\\[A-Za-z]+`)

var mathWords = map[string]bool{
	"sin": true, "cos": true, "tan": true, "cot": true, "sec": true, "csc": true,
	"arc": true, "log": true, "ln": true, "exp": true, "sqrt": true, "inf": true,
}

// looksLikeWords reports whether the answer contains a word a learner could not
// type into a math field. Splitting on non-letters keeps expressions out.
func looksLikeWords(answer string) bool {
	answer = latexCommand.ReplaceAllString(answer, " ")
	for _, field := range strings.FieldsFunc(answer, func(r rune) bool { return !unicode.IsLetter(r) }) {
		w := strings.ToLower(field)
		if len(w) >= 3 && !mathWords[w] {
			return true
		}
	}
	return false
}
