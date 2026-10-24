package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/repetition"
	"github.com/chuma-beep/mathua/internal/storage"
)

// TestPlacementCreditCycleOnlyAdvancesOnce settles the invariant the three
// engine changes share: a provisional placement is not evidence, and implicit
// credit earned downstream cannot turn ten answers in one sitting into ten
// spaced reviews of the prerequisite.
//
// The scenario is the one reported as a 30,595-day interval: a learner places
// out of "found", then answers ten questions on "dep" correctly. Under the
// cycle guard the ten answers bank one ordinal, projecting it once moves the
// schedule by one step, a retry moves nothing, and "found" stays provisional
// because no direct evidence exists for it.
func TestPlacementCreditCycleOnlyAdvancesOnce(t *testing.T) {
	e := traceEngine(t)
	diff := 0.5

	st, err := e.CreateStudent("credit-cycle")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := e.ApplyDiagnosticReport(st.ID, &diagnostic.DiagnosticReport{
		ConditionallyCompleted: []string{"found"},
		FrontierConditional:    true,
	}); err != nil {
		t.Fatalf("apply diagnostic report: %v", err)
	}

	// Ten correct answers on the concept that encompasses "found".
	for i := 1; i <= 10; i++ {
		dp, err := e.repo.GetProgress(st.ID, "dep")
		if err != nil {
			t.Fatalf("read dep progress: %v", err)
		}
		if dp == nil {
			dp = &storage.ConceptProgress{StudentID: st.ID, ConceptID: "dep", Status: "LEARNING"}
		}
		if _, _, err := e.applyGradedAttempt(dp, gradedAttempt{
			studentID: st.ID, conceptID: "dep", requiredStreak: 3,
			timeThreshold: 60, elapsed: 30, correct: true, difficulty: &diff,
		}); err != nil {
			t.Fatalf("downstream attempt %d: %v", i, err)
		}
		if err := e.repo.UpsertProgress(dp); err != nil {
			t.Fatalf("save dep after attempt %d: %v", i, err)
		}
	}

	cs, ok := e.repo.(storage.CreditStore)
	if !ok {
		t.Fatal("repository does not implement CreditStore")
	}
	banked, err := cs.GetBankedOrdinals(st.ID)
	if err != nil {
		t.Fatalf("banked ordinals: %v", err)
	}
	if got := banked["found"]; got != 1 {
		t.Fatalf("ten answers in one cycle banked %d ordinals for found, want 1", got)
	}

	// A review of the prerequisite falls due, which is the point at which
	// implicit credit may stand in for a direct answer.
	found, err := e.repo.GetProgress(st.ID, "found")
	if err != nil || found == nil {
		t.Fatalf("read found progress: %v", err)
	}
	due := time.Now().UTC().AddDate(0, 0, -1)
	found.NextReviewDue = &due
	if err := e.repo.UpsertProgress(found); err != nil {
		t.Fatalf("save found: %v", err)
	}
	epochBefore := epochOfCredit(t, cs, st.ID, "found")

	const stepDays = 25
	projected, err := cs.ProjectCredit(st.ID, "found", 0, stepDays, time.Now().UTC())
	if err != nil {
		t.Fatalf("project credit: %v", err)
	}
	if len(projected) != 1 {
		t.Fatalf("projection consumed %d ordinals, want 1", len(projected))
	}
	epochAfter := epochOfCredit(t, cs, st.ID, "found")
	if epochAfter != epochBefore+1 {
		t.Fatalf("epoch after projection = %d, want %d", epochAfter, epochBefore+1)
	}

	// One step, bounded: the projected due date must sit one step away, not
	// decades. MaxStepDays is the ceiling the projection itself enforces.
	if stepDays > repetition.MaxStepDays {
		t.Fatalf("test step of %d days exceeds MaxStepDays %d", stepDays, repetition.MaxStepDays)
	}
	projectedDue, err := time.Parse(time.RFC3339, projected[0].CounterfactualDue)
	if err != nil {
		t.Fatalf("parse projected due %q: %v", projected[0].CounterfactualDue, err)
	}
	shift := projectedDue.Sub(due)
	if shift < 0 || shift > time.Duration(stepDays*2)*24*time.Hour {
		t.Fatalf("projection shifted the schedule by %v, want one step near %d days", shift, stepDays)
	}
	if got := epochOfCredit(t, cs, st.ID, "found"); got != epochAfter {
		t.Fatalf("epoch moved to %d during the interval check, want %d", got, epochAfter)
	}

	// Retrying must consume nothing and leave the cycle exactly where it was.
	again, err := cs.ProjectCredit(st.ID, "found", 0, stepDays, time.Now().UTC())
	if err != nil {
		t.Fatalf("retry projection: %v", err)
	}
	if len(again) != 0 {
		t.Fatalf("retry consumed %d ordinals, want 0", len(again))
	}
	if got := epochOfCredit(t, cs, st.ID, "found"); got != epochAfter {
		t.Fatalf("retry advanced the epoch to %d, want %d", got, epochAfter)
	}

	// Implicit credit is not evidence: the placement stays provisional.
	found, err = e.repo.GetProgress(st.ID, "found")
	if err != nil || found == nil {
		t.Fatalf("re-read found: %v", err)
	}
	if !found.IsProvisional || !found.PlacementSeeded {
		t.Fatalf("placement was certified by implicit credit: provisional=%v seeded=%v",
			found.IsProvisional, found.PlacementSeeded)
	}
	if found.Status == "MASTERED" {
		t.Fatalf("implicit credit promoted a placement to MASTERED")
	}
}

// epochOfCredit resolves the concept's current credit cycle.
func epochOfCredit(t *testing.T, cs storage.CreditStore, student, concept string) int64 {
	t.Helper()
	acc, err := cs.AccountOf(student, concept, 0)
	if err != nil {
		t.Fatalf("account of %s: %v", concept, err)
	}
	return acc.Epoch
}
