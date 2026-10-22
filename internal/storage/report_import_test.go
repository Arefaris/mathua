package storage

import (
	"testing"
	"time"
)

// The import is the boundary between production and a local workbench, so the
// two properties that matter are that it keeps the evidence intact and that
// running the same export twice does not double the queue.
func TestImportReport_IsIdempotentAndKeepsTheEvidence(t *testing.T) {
	store := newTestStore(t)
	filed := time.Date(2026, 10, 20, 9, 30, 0, 0, time.UTC)

	r := QuestionReport{
		ReporterID:  "rep-1",
		ConceptID:   "arith.div.remainder",
		Kind:        "question",
		Question:    "What is the remainder when 17 is divided by 5?",
		Expected:    "2",
		AnswerGiven: "3",
		Reason:      "wrong_expected",
		Detail:      "The expected answer looks wrong.",
		Status:      "open",
		CreatedAt:   filed,
	}

	inserted, err := store.ImportReport("73", r)
	if err != nil {
		t.Fatalf("first import: %v", err)
	}
	if !inserted {
		t.Fatal("first import reported the report as already present")
	}

	// A second run of the same export must not add a row, even if the payload
	// carries a later local edit: the file is a copy of one report, not a feed.
	r.Question = "edited after export"
	inserted, err = store.ImportReport("73", r)
	if err != nil {
		t.Fatalf("second import: %v", err)
	}
	if inserted {
		t.Error("a repeated import inserted a duplicate")
	}

	reports, err := store.ListReports("open", 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("queue has %d reports, want 1", len(reports))
	}
	got := reports[0]
	if got.ConceptID != "arith.div.remainder" || got.Expected != "2" || got.AnswerGiven != "3" {
		t.Errorf("evidence not preserved: %+v", got)
	}
	if !got.CreatedAt.UTC().Equal(filed) {
		t.Errorf("created_at = %s, want %s", got.CreatedAt.UTC(), filed)
	}
	if got.Question != "What is the remainder when 17 is divided by 5?" {
		t.Errorf("the second payload overwrote the first: %q", got.Question)
	}
}

func TestImportReport_RequiresAnOriginID(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.ImportReport("   ", QuestionReport{ConceptID: "a"}); err == nil {
		t.Fatal("import accepted an empty origin id")
	}
}

// Locally filed reports have no origin id, and the uniqueness rule must not
// treat "none" as a value: two reports from the same concept are ordinary.
func TestImportReport_DoesNotConstrainLocallyFiledReports(t *testing.T) {
	store := newTestStore(t)
	r := QuestionReport{ConceptID: "arith.div.remainder", Kind: "question", Status: "open", CreatedAt: time.Now().UTC()}
	if _, err := store.CreateReport(r); err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err := store.CreateReport(r); err != nil {
		t.Fatalf("create second: %v", err)
	}
	reports, err := store.ListReports("open", 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(reports) != 2 {
		t.Fatalf("queue has %d reports, want 2", len(reports))
	}
}
