package fractions

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 100; i++ {
		d := rand.Float64()
		p := gen.Generate(generator.GeneratorContext{Difficulty: d})
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
		}
	}
}

func TestFracConceptGen(t *testing.T)      { fuzzGen(t, &fracConceptGen{}) }
func TestFracPartsGen(t *testing.T)        { fuzzGen(t, &fracPartsGen{}) }
func TestFracNumberLineGen(t *testing.T)   { fuzzGen(t, &fracNumberLineGen{}) }
func TestFracEquivalentGen(t *testing.T)   { fuzzGen(t, &fracEquivalentGen{}) }
func TestFracSimplifyGen(t *testing.T)     { fuzzGen(t, &fracSimplifyGen{}) }
func TestFracCompareGen(t *testing.T)      { fuzzGen(t, &fracCompareGen{}) }
func TestFracBenchmarkGen(t *testing.T)    { fuzzGen(t, &fracBenchmarkGen{}) }
func TestFracToDecimalGen(t *testing.T)    { fuzzGen(t, &fracToDecimalGen{}) }
func TestFracOpSameDenGen(t *testing.T)    { fuzzGen(t, &fracOpSameDenGen{}) }
func TestFracOpDiffDenGen(t *testing.T)    { fuzzGen(t, &fracOpDiffDenGen{}) }
func TestFracWordGen(t *testing.T)         { fuzzGen(t, &fracWordGen{}) }
func TestFracMultGen(t *testing.T)         { fuzzGen(t, &fracMultGen{}) }
func TestFracDivGen(t *testing.T)          { fuzzGen(t, &fracDivGen{}) }
func TestFracMixedConvertGen(t *testing.T) { fuzzGen(t, &fracMixedConvertGen{}) }
func TestFracMixedOpGen(t *testing.T)      { fuzzGen(t, &fracMixedOpGen{}) }
func TestFracMixedMultGen(t *testing.T)    { fuzzGen(t, &fracMixedMultGen{}) }

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

// The bar was drawn as "X X _ _ _" inside math delimiters, where the spaces
// collapse and the underscores become subscripts, so the learner was asked what
// fraction of a bar they could not see. The question now says the counts.
func TestFractionConceptAsksInWords(t *testing.T) {
	g := &fracConceptGen{}
	for i := 0; i < 200; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64(), Seed: rand.Int63()})
		if strings.Contains(p.Question, "X X") || strings.Contains(p.Question, "\\(") {
			t.Fatalf("question still draws a bar: %q", p.Question)
		}
		if !strings.Contains(p.Question, "equal parts") {
			t.Fatalf("question does not describe the parts: %q", p.Question)
		}
		if p.Answer != p.Facts["filled"]+"/"+p.Facts["parts"] {
			t.Fatalf("answer %q does not match facts %v", p.Answer, p.Facts)
		}
	}
}
