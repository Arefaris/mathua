package storage

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/repetition"
)

func newCreditStore(t *testing.T) *SQLiteStore {
	t.Helper()
	store, err := NewSQLiteStore(t.TempDir() + "/credit.db")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func seedDue(t *testing.T, store *SQLiteStore, studentID, conceptID string, due time.Time) {
	t.Helper()
	if _, err := store.db.Exec(`
		INSERT INTO concept_progress (student_id, concept_id, status, next_review_due)
		VALUES (?, ?, 'PRACTICING', ?)
		ON CONFLICT(student_id, concept_id) DO UPDATE SET next_review_due = excluded.next_review_due
	`, studentID, conceptID, due.UTC().Format("2006-01-02 15:04:05")); err != nil {
		t.Fatalf("seed progress: %v", err)
	}
}

func creditLearner(t *testing.T, store *SQLiteStore) string {
	t.Helper()
	st, err := store.CreateStudent("credit learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	return st.ID
}

func account(t *testing.T, store *SQLiteStore, studentID, conceptID string, epoch int64) CreditAccount {
	t.Helper()
	acc, err := store.AccountOf(studentID, conceptID, epoch)
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	return acc
}

func TestAllocateCredit_EarnsFractionalCreditWithoutBanking(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)

	alloc, fresh, err := store.AllocateCredit(student, "attempt:1", 1, []CreditCandidate{{ConceptID: "a", RawUnits: 800_000}})
	if err != nil || !fresh {
		t.Fatalf("allocate: %v fresh=%v", err, fresh)
	}
	if alloc.AllocatedUnits != 800_000 || len(alloc.Events) != 1 {
		t.Fatalf("allocation = %+v", alloc)
	}
	acc := account(t, store, student, "a", 1)
	if acc.RetainedUnits != 800_000 || acc.Banked != 0 || acc.Balance() != 800_000 {
		t.Fatalf("account = %+v balance=%d", acc, acc.Balance())
	}
}

func TestAllocateCredit_ScalesTheFanOutToThePerAttemptCap(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)

	alloc, _, err := store.AllocateCredit(student, "attempt:2", 1, []CreditCandidate{
		{ConceptID: "a", RawUnits: 800_000},
		{ConceptID: "b", RawUnits: 800_000},
	})
	if err != nil {
		t.Fatalf("allocate: %v", err)
	}
	if alloc.AllocatedUnits != repetition.Micro {
		t.Fatalf("allocated = %d, want one repetition", alloc.AllocatedUnits)
	}
	for _, ev := range alloc.Events {
		if ev.ScaledUnits != 500_000 {
			t.Errorf("target %s scaled = %d, want 500000", ev.ConceptID, ev.ScaledUnits)
		}
	}
}

func TestAllocateCredit_ReplaysWithoutBankingTwice(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)

	first, fresh, err := store.AllocateCredit(student, "attempt:3", 1, []CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}})
	if err != nil || !fresh {
		t.Fatalf("first: %v fresh=%v", err, fresh)
	}
	if got := account(t, store, student, "a", 1).Banked; got != 1 {
		t.Fatalf("banked after first = %d, want 1", got)
	}
	second, fresh, err := store.AllocateCredit(student, "attempt:3", 1, []CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}})
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if fresh {
		t.Error("replay reported a fresh allocation")
	}
	if second.AllocatedUnits != first.AllocatedUnits || len(second.Events) != len(first.Events) {
		t.Errorf("replay = %+v, want %+v", second, first)
	}
	if got := account(t, store, student, "a", 1).Banked; got != 1 {
		t.Errorf("replay banked again: ordinals = %d, want 1", got)
	}
}

func TestAllocateCredit_BanksOncePerWholeRepetitionAndCapsLifetime(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)

	// 900k three times is 2.7 repetitions: two ordinals, 0.7 left over.
	for i := 0; i < 3; i++ {
		if _, _, err := store.AllocateCredit(student, fmt.Sprintf("attempt:%d", 10+i), 1,
			[]CreditCandidate{{ConceptID: "a", RawUnits: 900_000}}); err != nil {
			t.Fatalf("allocate %d: %v", i, err)
		}
	}
	acc := account(t, store, student, "a", 1)
	if acc.Banked != 2 || acc.Balance() != 700_000 {
		t.Fatalf("account = %+v balance=%d, want 2 ordinals and 700000", acc, acc.Balance())
	}

	// A fourth 900k crosses the cap: 300k is retained, 600k discarded.
	alloc, _, err := store.AllocateCredit(student, "attempt:20", 1, []CreditCandidate{{ConceptID: "a", RawUnits: 900_000}})
	if err != nil {
		t.Fatalf("allocate over cap: %v", err)
	}
	ev := alloc.Events[0]
	if ev.RetainedUnits != 300_000 || ev.DiscardedUnits != 600_000 || ev.DiscardReason != "epoch_bank_limit" {
		t.Fatalf("over-cap event = %+v", ev)
	}
	if got := account(t, store, student, "a", 1).RetainedUnits; got != 3*repetition.Micro {
		t.Errorf("retained = %d, want the lifetime cap", got)
	}
}

func TestProjectCredit_ChainsFromTheRealDueAndLeavesItAlone(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)
	realDue := time.Date(2026, 11, 1, 9, 30, 0, 0, time.UTC)
	seedDue(t, store, student, "a", realDue)

	for i := 0; i < 2; i++ {
		if _, _, err := store.AllocateCredit(student, fmt.Sprintf("attempt:%d", 30+i), 1,
			[]CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}}); err != nil {
			t.Fatalf("allocate: %v", err)
		}
	}

	now := time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)
	ordinals, err := store.ProjectCredit(student, "a", 1, 25, now)
	if err != nil {
		t.Fatalf("project: %v", err)
	}
	if len(ordinals) != 2 {
		t.Fatalf("projected %d ordinals, want 2", len(ordinals))
	}
	if ordinals[0].CounterfactualDue != "2026-11-26T09:30:00Z" || ordinals[1].CounterfactualDue != "2026-12-21T09:30:00Z" {
		t.Errorf("chain = %q then %q", ordinals[0].CounterfactualDue, ordinals[1].CounterfactualDue)
	}

	var stored string
	if err := store.db.QueryRow(`SELECT next_review_due FROM concept_progress WHERE student_id = ? AND concept_id = ?`,
		student, "a").Scan(&stored); err != nil {
		t.Fatalf("read due: %v", err)
	}
	if stored != "2026-11-01 09:30:00" {
		t.Fatalf("production due date moved to %q", stored)
	}

	again, err := store.ProjectCredit(student, "a", 1, 25, now)
	if err != nil {
		t.Fatalf("retry projection: %v", err)
	}
	for i := range again {
		if again[i].CounterfactualDue != ordinals[i].CounterfactualDue {
			t.Errorf("retry changed ordinal %d: %q vs %q", i, again[i].CounterfactualDue, ordinals[i].CounterfactualDue)
		}
	}
}

func TestProjectCredit_WithoutADueDateStaysBanked(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)
	if _, _, err := store.AllocateCredit(student, "attempt:40", 1,
		[]CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}}); err != nil {
		t.Fatalf("allocate: %v", err)
	}
	out, err := store.ProjectCredit(student, "a", 1, 25, time.Now().UTC())
	if err != nil {
		t.Fatalf("project without due: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("projected %d ordinals without a schedule", len(out))
	}
	if got := account(t, store, student, "a", 1).Banked; got != 1 {
		t.Errorf("banked = %d, want it left alone", got)
	}
}

func TestResetCreditEpoch_StalesBankedOrdinalsAndStartsFresh(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)
	seedDue(t, store, student, "a", time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC))
	if _, _, err := store.AllocateCredit(student, "attempt:50", 1,
		[]CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}}); err != nil {
		t.Fatalf("allocate: %v", err)
	}

	next, err := store.ResetCreditEpoch(student, "a")
	if err != nil || next != 2 {
		t.Fatalf("reset: epoch=%d err=%v", next, err)
	}
	old := account(t, store, student, "a", 1)
	// The old epoch keeps its history: one ordinal, now stale, and the credit
	// that funded it. What it must not keep is a usable banked repetition.
	if old.Banked != 1 || old.Stale != 1 || old.Balance() != 0 {
		t.Errorf("old epoch = %+v balance=%d, want one stale ordinal and zero balance", old, old.Balance())
	}
	fresh := account(t, store, student, "a", 2)
	if fresh.RetainedUnits != 0 || fresh.Banked != 0 || fresh.Stale != 0 {
		t.Errorf("new epoch = %+v, want empty", fresh)
	}
	if _, _, err := store.AllocateCredit(student, "attempt:51", 1,
		[]CreditCandidate{{ConceptID: "a", RawUnits: 100_000}}); err == nil {
		t.Error("allocation into a closed epoch was accepted")
	}
}

func TestVoidCreditEvent_RespectsFundingAndEpoch(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)
	seedDue(t, store, student, "a", time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC))

	if _, _, err := store.AllocateCredit(student, "attempt:60", 1,
		[]CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}}); err != nil {
		t.Fatalf("allocate: %v", err)
	}
	// One banked ordinal is funded by this event; voiding it would leave the
	// ordinal unfunded, so the ledger must refuse rather than restate history.
	if err := store.VoidCreditEvent(student, "attempt:60", "a", "test"); err == nil {
		t.Fatal("void left a banked ordinal unfunded")
	}

	if _, _, err := store.AllocateCredit(student, "attempt:61", 1,
		[]CreditCandidate{{ConceptID: "a", RawUnits: 500_000}}); err != nil {
		t.Fatalf("allocate: %v", err)
	}
	// The 500k event is not funding the one banked ordinal, so removing it
	// leaves every ordinal funded and must be allowed.
	if err := store.VoidCreditEvent(student, "attempt:61", "a", "credited in error"); err != nil {
		t.Fatalf("void funded event: %v", err)
	}
	if got := account(t, store, student, "a", 1).RetainedUnits; got != 1_000_000 {
		t.Errorf("retained after void = %d, want the other event", got)
	}
}

func TestAllocateCredit_ConcurrentDuplicatesInsertOnce(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)

	const n = 8
	var wg sync.WaitGroup
	errs := make([]error, n)
	fresh := make([]bool, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, fresh[i], errs[i] = store.AllocateCredit(student, "attempt:70", 1,
				[]CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}})
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("worker %d: %v", i, err)
		}
	}
	winners := 0
	for _, f := range fresh {
		if f {
			winners++
		}
	}
	if winners != 1 {
		t.Errorf("%d workers reported a fresh allocation, want 1", winners)
	}
	if got := account(t, store, student, "a", 1).Banked; got != 1 {
		t.Errorf("banked = %d, want exactly one ordinal", got)
	}
}
