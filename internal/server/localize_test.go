package server

import (
	"os"
	"regexp"
	"testing"

	"github.com/chuma-beep/mathua/internal/scheduler"
)

// Every badge the scheduler and engine can put on a card must have a Russian wording, or a
// Russian learner sees one English card among Russian ones. The badges are read from source so
// a new card added upstream fails here instead of shipping half-translated.
func TestEveryRecommendationBadgeHasRussianWording(t *testing.T) {
	badge := regexp.MustCompile(`Badge\s*(?::|=)\s*"([^"]+)"`)
	seen := 0
	for _, f := range []string{"../scheduler/recommend.go", "../engine/recommend.go"} {
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range badge.FindAllStringSubmatch(string(src), -1) {
			seen++
			if _, ok := recommendationRU[m[1]]; !ok {
				t.Errorf("%s: badge %q has no Russian wording in recommendationRU", f, m[1])
			}
		}
	}
	if seen == 0 {
		t.Fatal("found no badges; the pattern no longer matches the source")
	}
}

func TestWordRecommendationRU(t *testing.T) {
	rec := scheduler.Recommendation{ID: "diagnostic", Badge: "Recommended", Detail: "x", CTA: "y"}
	wordRecommendationRU(&rec)
	if rec.Badge != "Рекомендуем" || rec.CTA != "Пройти диагностику →" || rec.ConceptTitle != "Диагностика" {
		t.Fatalf("diagnostic card not worded in Russian: %+v", rec)
	}

	concept := scheduler.Recommendation{ID: "frac.add", ConceptTitle: "Сложение дробей", Badge: "New"}
	wordRecommendationRU(&concept)
	if concept.ConceptTitle != "Сложение дробей" || concept.Badge != "Новая тема" {
		t.Fatalf("concept card: %+v", concept)
	}

	unknown := scheduler.Recommendation{ID: "x", Badge: "Something new", Detail: "d"}
	wordRecommendationRU(&unknown)
	if unknown.Badge != "Something new" || unknown.Detail != "d" {
		t.Fatalf("unknown badge must keep its English wording: %+v", unknown)
	}
}
