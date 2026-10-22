package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/repetition"
)

// moveDueInsideProjection installs a hook that moves the real due date from
// inside the projection transaction, immediately before the conflict re-read.
//
// This is the only way the conflict path can be reached: the guard is defensive,
// not concurrent. SQLite holds the write reservation for the whole transaction
// (BEGIN IMMEDIATE) and PostgreSQL holds the concept_progress row FOR UPDATE, so
// a competing writer does not race the projection - it waits. What the test
// proves is therefore the part that matters operationally: when the baseline the
// chain was built from no longer matches, every ordinal write in that
// transaction is rolled back rather than half-committed.
func moveDueInsideProjection(t *testing.T, repo Repository, studentID, conceptID, newDue string) {
	t.Helper()
	_, postgres := repo.(*PostgresStore)
	creditProjectTestHook = func(exec func(query string, args ...any) error) {
		var err error
		if postgres {
			err = exec(
				`UPDATE concept_progress SET next_review_due = $1 WHERE student_id = $2 AND concept_id = $3`,
				newDue, studentID, conceptID)
		} else {
			err = exec(
				`UPDATE concept_progress SET next_review_due = ? WHERE student_id = ? AND concept_id = ?`,
				newDue, studentID, conceptID)
		}
		if err != nil {
			t.Fatalf("hook could not move the due date: %v", err)
		}
	}
	t.Cleanup(func() { creditProjectTestHook = nil })
}

// TestCreditParity_ProjectionRollsBackWhenTheDueDateMoves covers the
// ErrCreditDueConflict path on both backends: the chain is abandoned, the
// ordinals stay banked, and the real schedule is untouched.
func TestCreditParity_ProjectionRollsBackWhenTheDueDateMoves(t *testing.T) {
	for tag, repo := range creditStores(t) {
		repo := repo
		t.Run(tag, func(t *testing.T) {
			credit := repo.(CreditStore)
			stu, err := repo.CreateStudent("credit conflict " + tag)
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			due := "2026-11-01T09:30:00Z"
			seedDueRepo(t, repo, stu.ID, "a", mustTime(t, due))

			if _, fresh, err := credit.AllocateCredit(stu.ID, tag+":conflict", 1,
				[]CreditCandidate{{ConceptID: "a", RawUnits: repetition.Micro}}); err != nil || !fresh {
				t.Fatalf("allocate: %v fresh=%v", err, fresh)
			}
			moveDueInsideProjection(t, repo, stu.ID, "a", "2026-12-01T09:30:00Z")

			now := time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)
			if _, err := credit.ProjectCredit(stu.ID, "a", 1, 25, now); !errors.Is(err, ErrCreditDueConflict) {
				t.Fatalf("projection error = %v, want ErrCreditDueConflict", err)
			}

			// Rollback undoes the hook's write as well: the transaction's own
			// change cannot survive the failure that change caused.
			if got := readDueRepo(t, repo, stu.ID, "a"); got != "2026-11-01 09:30:00" {
				t.Errorf("due after rollback = %q, want the original schedule", got)
			}
			acc, err := credit.AccountOf(stu.ID, "a", 1)
			if err != nil {
				t.Fatalf("account: %v", err)
			}
			if acc.Banked != 1 || acc.Projected != 0 || acc.Stale != 0 {
				t.Errorf("account after rollback = %+v, want one unprojected ordinal", acc)
			}
			if acc.Balance() != 0 || acc.RetainedUnits != repetition.Micro {
				t.Errorf("credit after rollback = retained %d balance %d", acc.RetainedUnits, acc.Balance())
			}
			// Nothing is left to project: a later honest run still has the ordinal.
			creditProjectTestHook = nil
			ords, err := credit.ProjectCredit(stu.ID, "a", 1, 25, now)
			if err != nil || len(ords) != 1 || ords[0].State != "projected" {
				t.Fatalf("projection after recovery = %+v err=%v", ords, err)
			}
		})
	}
}

// TestCreditParity_ResetBeforeProjectionIsNotUndone covers the other ordering:
// the reset lands first, stales the old banked ordinal and moves the epoch, and a
// projection still holding the old epoch must refuse rather than project a stale
// ordinal or write into the new epoch.
func TestCreditParity_ResetBeforeProjectionIsNotUndone(t *testing.T) {
	for tag, repo := range creditStores(t) {
		repo := repo
		t.Run(tag, func(t *testing.T) {
			credit := repo.(CreditStore)
			stu, err := repo.CreateStudent("credit reset first " + tag)
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			seedDueRepo(t, repo, stu.ID, "a", mustTime(t, "2026-11-01T09:30:00Z"))

			if _, fresh, err := credit.AllocateCredit(stu.ID, tag+":reset", 1,
				[]CreditCandidate{{ConceptID: "a", RawUnits: repetition.Micro}}); err != nil || !fresh {
				t.Fatalf("allocate: %v fresh=%v", err, fresh)
			}
			next, err := credit.ResetCreditEpoch(stu.ID, "a")
			if err != nil || next != 2 {
				t.Fatalf("reset = %d err=%v", next, err)
			}
			now := time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)
			if _, err := credit.ProjectCredit(stu.ID, "a", 1, 25, now); !errors.Is(err, ErrCreditOrdinalStale) {
				t.Fatalf("projection of the closed epoch = %v, want ErrCreditOrdinalStale", err)
			}
			old, err := credit.AccountOf(stu.ID, "a", 1)
			if err != nil {
				t.Fatalf("old account: %v", err)
			}
			if old.Stale != 1 || old.Projected != 0 {
				t.Errorf("old epoch = %+v, want the ordinal stale and unprojected", old)
			}
			fresh, err := credit.AccountOf(stu.ID, "a", 2)
			if err != nil {
				t.Fatalf("new account: %v", err)
			}
			if fresh.RetainedUnits != 0 || fresh.Banked != 0 || fresh.Projected != 0 || fresh.Stale != 0 {
				t.Errorf("new epoch inherited %+v", fresh)
			}
			if got := readDueRepo(t, repo, stu.ID, "a"); got != "2026-11-01 09:30:00" {
				t.Errorf("due moved to %q during a refused projection", got)
			}
		})
	}
}

func mustTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("parse %q: %v", value, err)
	}
	return parsed
}

// TestProjectionConflictGuardIsDefensive documents, next to the code it
// describes, why the conflict cannot be produced by two goroutines. If a future
// change drops the row lock this assertion is what fails, rather than the
// guarantee silently disappearing.
func TestProjectionConflictGuardIsDefensive(t *testing.T) {
	store := newCreditStore(t)
	stu := creditLearner(t, store)
	due := time.Date(2026, 11, 1, 9, 30, 0, 0, time.UTC)
	seedDue(t, store, stu, "a", due)
	if _, _, err := store.AllocateCredit(stu, "defensive:1", 1,
		[]CreditCandidate{{ConceptID: "a", RawUnits: repetition.Micro}}); err != nil {
		t.Fatalf("allocate: %v", err)
	}
	// No hook: under the write reservation the baseline cannot move, so the
	// projection succeeds and the guard is never triggered.
	now := time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)
	if _, err := store.ProjectCredit(stu, "a", 1, 25, now); err != nil {
		t.Fatalf("clean projection: %v", err)
	}
	if creditProjectTestHook != nil {
		t.Fatal("a test hook leaked into the production path")
	}
}
