package engine

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

func traceEngine(t *testing.T) *Engine {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{
			ID: "found", Label: "Foundation", Domain: "d", GradingType: "numeric",
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
		{
			ID: "dep", Label: "Dependent", Domain: "d", GradingType: "numeric",
			Prerequisites:    []string{"found"},
			Encompasses:      []string{"found"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("found", &testGen{answer: "42"})
	reg.Register("dep", &testGen{answer: "99"})
	return New(store, d, reg, nil, nil)
}

// TestCombinedTrace reproduces, on one branch carrying 150, 152 and 155
