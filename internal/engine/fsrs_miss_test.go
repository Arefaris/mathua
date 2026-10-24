package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/fsrs"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
)

// reviewStateRow builds a row that already has a schedule to lose: a stability of 500
// days, a review timestamp from 480 days ago and a due date 400 days out. Every assertion
// about a lapse has to start from a row in that state, because the defect was invisible
// on a row that had never been scheduled.
func reviewStateRow(t *testing.T, studentID, status string, streak int, now time.Time) *storage.ConceptProgress {
	t.Helper()
	stability := 500.0
	difficulty := 5.0
	lastReview := now.AddDate(0, 0, -480)
	due := now.AddDate(0, 0, 400)
	return &storage.ConceptProgress{
		StudentID:       studentID,
		ConceptID:       "a",
		Status:          status,
		Streak:          streak,
		BestStreak:      streak,
		AvgResponseTime: 6,
		Attempts:        6,
		LastReviewed:    &lastReview,
		NextReviewDue:   &due,
		SM2Repetitions:  4,
		SM2Interval:     400,
		SM2EFactor:      2.5,
		WeaknessScore:   0.1,
		FSRSStability:   &stability,
		FSRSDifficulty:  &difficulty,
		FSRSLastReview:  &lastReview,
	}
}

// TestFSRSMissSchedulesImmediately pins the defect this repair exists for. A missed
// review used to collapse the stability but leave FSRSLastReview and NextReviewDue
// holding the values from the last successful review, so retrievability computed from a
// stale clock while the shelf displayed a due date months away.
func TestFSRSMissSchedulesImmediately(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("fsrs-miss")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	p := reviewStateRow(t, st.ID, string(mastery.StatusMastered), 4, now)
	stabilityBefore := *p.FSRSStability
	dueBefore := *p.NextReviewDue

	if _, _, err := e.applyGradedAttempt(p, gradedAttempt{
		studentID: st.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 40, correct: false,
	}); err != nil {
		t.Fatalf("apply miss: %v", err)
	}

	// The lapse policy still collapses stability rather than nudging it.
	if p.FSRSStability == nil || *p.FSRSStability >= stabilityBefore/2 {
		t.Fatalf("stability after miss = %v, want a collapse from %v", p.FSRSStability, stabilityBefore)
	}
	// The review clock moves to this attempt, not the one from months ago.
	if p.FSRSLastReview == nil {
		t.Fatal("FSRSLastReview is nil after a miss")
	}
	if age := time.Since(*p.FSRSLastReview); age > time.Minute || age < -time.Minute {
		t.Errorf("FSRSLastReview after miss is %v old, want this attempt", age)
	}
	// The concept stays immediately eligible instead of being parked months out.
	if p.NextReviewDue == nil {
		t.Fatal("NextReviewDue is nil after a miss")
	}
	if p.NextReviewDue.After(time.Now().UTC().Add(time.Minute)) {
		t.Errorf("NextReviewDue after miss = %v, want no later than now (was %v)", p.NextReviewDue, dueBefore)
	}
	// SM-2 keeps its own lapse behaviour: interval resets to a single day.
	if p.SM2Interval != 1 {
		t.Errorf("SM2Interval after miss = %d, want 1", p.SM2Interval)
	}

	// Retrievability now reads from this attempt. With the stale timestamp it was
	// effectively zero, which is what kept the row permanently due.
	params := fsrs.DefaultParameters()
	elapsed := time.Since(*p.FSRSLastReview).Hours() / 24
	if r := params.Retrievability(elapsed, *p.FSRSStability); r < fsrs.DefaultRetention {
		t.Errorf("retrievability %.3f after miss, want at least %.2f", r, fsrs.DefaultRetention)
	}

	// The queue agrees with the record, and the persisted row is the one asserted above.
	if due, err := e.ReviewsDue(st.ID); err != nil {
		t.Fatalf("reviews due: %v", err)
	} else if due < 1 {
		t.Errorf("reviews due after miss = %d, want the concept counted", due)
	}
	stored, err := e.repo.GetProgress(st.ID, "a")
	if err != nil || stored == nil {
		t.Fatalf("read stored progress: %v", err)
	}
	if stored.FSRSLastReview == nil || stored.FSRSLastReview.Before(now.Add(-time.Minute)) {
		t.Errorf("stored FSRSLastReview = %v, want this attempt", stored.FSRSLastReview)
	}
	if stored.NextReviewDue == nil || stored.NextReviewDue.After(time.Now().UTC().Add(time.Minute)) {
		t.Errorf("stored NextReviewDue = %v, want no later than now", stored.NextReviewDue)
	}

	// A miss schedules work; it never creates mastery. The row above was already
	// MASTERED and the ladder's hysteresis keeps it there, which is the same rule 155
	// applied to placements, so the assertion belongs on a row that is not yet mastered.
	st2, err := e.CreateStudent("fsrs-miss-new")
	if err != nil {
		t.Fatalf("create second student: %v", err)
	}
	fresh := reviewStateRow(t, st2.ID, string(mastery.StatusPracticing), 1, now)
	if _, _, err := e.applyGradedAttempt(fresh, gradedAttempt{
		studentID: st2.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 40, correct: false,
	}); err != nil {
		t.Fatalf("apply miss to practicing row: %v", err)
	}
	if fresh.Status == string(mastery.StatusMastered) {
		t.Errorf("status after a miss on a practicing row = %q, want not mastered", fresh.Status)
	}
}

// TestFSRSMissThenSuccessReschedules checks that the repair does not freeze the row in
// the lapsed state: repeated correct answers have to advance it, move the due date into
// the future and keep the review clock current.
func TestFSRSMissThenSuccessReschedules(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("fsrs-recover")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	p := reviewStateRow(t, st.ID, string(mastery.StatusPracticing), 0, now)

	if _, _, err := e.applyGradedAttempt(p, gradedAttempt{
		studentID: st.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 40, correct: false,
	}); err != nil {
		t.Fatalf("apply miss: %v", err)
	}
	lapsed := *p.FSRSStability

	previous := lapsed
	for i := 1; i <= 3; i++ {
		if _, _, err := e.applyGradedAttempt(p, gradedAttempt{
			studentID: st.ID, conceptID: "a", requiredStreak: 3,
			timeThreshold: 60, elapsed: 6, correct: true,
		}); err != nil {
			t.Fatalf("apply recovery %d: %v", i, err)
		}
		if *p.FSRSStability < previous {
			t.Errorf("stability fell on recovery %d: %f after %f", i, *p.FSRSStability, previous)
		}
		previous = *p.FSRSStability
		if age := time.Since(*p.FSRSLastReview); age > time.Minute || age < -time.Minute {
			t.Errorf("FSRSLastReview after recovery %d is %v old, want this attempt", i, age)
		}
	}

	if p.NextReviewDue == nil || !p.NextReviewDue.After(time.Now().UTC().Add(24*time.Hour)) {
		t.Errorf("NextReviewDue after three correct answers = %v, want a future date", p.NextReviewDue)
	}
}

// TestFSRSCorrectAnswerMovesReviewClock pins the second half of the repair: a correct
// answer that does not advance a tier moves the review clock while leaving the due date
// alone. Without it, stability is recomputed from an elapsed time that never resets, so
// repeating a question inflates the estimate with no time passing.
func TestFSRSCorrectAnswerMovesReviewClock(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("fsrs-clock")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	p := reviewStateRow(t, st.ID, string(mastery.StatusPracticing), 0, now)
	staleReview := now.AddDate(0, 0, -60)
	p.FSRSLastReview = &staleReview
	dueBefore := *p.NextReviewDue
	stabilityBefore := *p.FSRSStability

	// One correct answer with the streak still below the tier threshold, so no advance.
	if _, _, err := e.applyGradedAttempt(p, gradedAttempt{
		studentID: st.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 6, correct: true,
	}); err != nil {
		t.Fatalf("apply correct answer: %v", err)
	}

	if age := time.Since(*p.FSRSLastReview); age > time.Minute || age < -time.Minute {
		t.Errorf("FSRSLastReview after a correct answer is %v old, want this attempt", age)
	}
	if *p.FSRSStability <= stabilityBefore {
		t.Errorf("stability after a correct answer = %f, want growth from %f", *p.FSRSStability, stabilityBefore)
	}
	// The due date is gated on the tier, and this answer did not reach it.
	if !p.NextReviewDue.Equal(dueBefore) {
		t.Errorf("NextReviewDue moved on a non-advancing answer: %v, want %v", p.NextReviewDue, dueBefore)
	}
}

// TestFSRSImplicitCreditAdvancesOncePerCycle pins that implicit credit moves FSRS state
// once per cycle and never again. PR #152 made the ordinal side of that rule explicit;
// this asserts the stability side follows it, because the FSRS advance sits behind the
// same ProjectCredit guard.
func TestFSRSImplicitCreditAdvancesOncePerCycle(t *testing.T) {
	e := traceEngine(t)
	diff := 0.5
	st, err := e.CreateStudent("fsrs-credit")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := e.ApplyDiagnosticReport(st.ID, &diagnostic.DiagnosticReport{
		ConditionallyCompleted: []string{"found"},
		FrontierConditional:    true,
	}); err != nil {
		t.Fatalf("apply diagnostic report: %v", err)
	}

	// Ten answers in one sitting: the cycle guard banks exactly one ordinal.
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

	// Give the prerequisite FSRS state and make its review due, which is the only moment
	// implicit credit may stand in for a direct answer.
	found, err := e.repo.GetProgress(st.ID, "found")
	if err != nil || found == nil {
		t.Fatalf("read found progress: %v", err)
	}
	stability := 40.0
	difficulty := 5.0
	lastReview := time.Now().UTC().AddDate(0, 0, -60)
	due := time.Now().UTC().AddDate(0, 0, -1)
	found.FSRSStability = &stability
	found.FSRSDifficulty = &difficulty
	found.FSRSLastReview = &lastReview
	found.NextReviewDue = &due
	if err := e.repo.UpsertProgress(found); err != nil {
		t.Fatalf("save found: %v", err)
	}

	if err := e.SatisfyImplicitReviews(st.ID); err != nil {
		t.Fatalf("project credit: %v", err)
	}
	after, err := e.repo.GetProgress(st.ID, "found")
	if err != nil || after == nil || after.FSRSStability == nil {
		t.Fatalf("read found after projection: %v", err)
	}
	first := *after.FSRSStability
	if first <= stability {
		t.Fatalf("stability after projection = %f, want above the seeded %f", first, stability)
	}

	// Further calls must find nothing left to project and leave stability untouched.
	for i := 0; i < 2; i++ {
		if err := e.SatisfyImplicitReviews(st.ID); err != nil {
			t.Fatalf("reproject %d: %v", i, err)
		}
		again, err := e.repo.GetProgress(st.ID, "found")
		if err != nil || again == nil || again.FSRSStability == nil {
			t.Fatalf("read found after reproject %d: %v", i, err)
		}
		if *again.FSRSStability != first {
			t.Fatalf("reprojection %d moved stability from %f to %f, want it unchanged", i, first, *again.FSRSStability)
		}
	}
}
