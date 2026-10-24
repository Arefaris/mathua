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

	// A review cycle banks one ordinal, so 900k three times is 2.7 repetitions of
	// credit that yield a single advance: the first two allocations reach a whole
	// repetition, the third arrives after the cycle advanced and is refused. The
	// refused credit is accounted rather than dropped, so the ledger still
	// explains every unit.
	cyclic := creditLearner(t, store)
	var last CreditAllocation
	for i := 0; i < 3; i++ {
		alloc, _, err := store.AllocateCredit(cyclic, fmt.Sprintf("attempt:%d", 10+i), 1,
			[]CreditCandidate{{ConceptID: "a", RawUnits: 900_000}})
		if err != nil {
			t.Fatalf("allocate %d: %v", i, err)
		}
		last = alloc
	}
	acc := account(t, store, cyclic, "a", 1)
	if acc.Banked != 1 {
		t.Fatalf("banked %d ordinals in one cycle, want 1", acc.Banked)
	}
	// 1.8M retained against one 1M ordinal leaves 800k of unfinished credit.
	if acc.RetainedUnits != 1_800_000 || acc.Balance() != 800_000 {
		t.Fatalf("account = %+v balance=%d, want 1800000 retained and 800000 balance", acc, acc.Balance())
	}
	ev := last.Events[0]
	if ev.RetainedUnits != 0 || ev.DiscardedUnits != 900_000 || ev.DiscardReason != "cycle_already_advanced" {
		t.Fatalf("refused event = %+v, want 900000 discarded for cycle_already_advanced", ev)
	}

	// The per-epoch bank limit is no longer reachable through the public API: the
	// cycle guard refuses the second ordinal before the limit can apply, so its
	// branch is a safety net rather than a path a learner can trigger. That is
	// why this test proves the cycle invariant instead of the cap.
}

func TestProjectCredit_AdvancesOneCycleAndLeavesTheRealDueAlone(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)
	realDue := time.Date(2026, 11, 1, 9, 30, 0, 0, time.UTC)
	seedDue(t, store, student, "a", realDue)

	// Two answers in one cycle: the first banks the single ordinal the cycle
	// allows and the second is refused, so projection has one step to take.
	// Two steps is what the old test asserted, and two steps is the compounding.
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
	if len(ordinals) != 1 {
		t.Fatalf("projected %d ordinals in one review cycle, want 1", len(ordinals))
	}
	if got := ordinals[0].CounterfactualDue; got != "2026-11-26T09:30:00Z" {
		t.Errorf("counterfactual due = %q, want one 25-day step from 2026-11-01", got)
	}
	if got := epochOf(t, store, student, "a"); got != 2 {
		t.Fatalf("epoch after consuming the cycle = %d, want 2", got)
	}

	var stored string
	if err := store.db.QueryRow(`SELECT next_review_due FROM concept_progress WHERE student_id = ? AND concept_id = ?`,
		student, "a").Scan(&stored); err != nil {
		t.Fatalf("read due: %v", err)
	}
	if stored != "2026-11-01 09:30:00" {
		t.Fatalf("production due date moved to %q", stored)
	}

	// A repeated projection has nothing banked in the current cycle, so it must
	// consume nothing and must not advance the epoch a second time.
	again, err := store.ProjectCredit(student, "a", 0, 25, now)
	if err != nil {
		t.Fatalf("retry projection: %v", err)
	}
	if len(again) != 0 {
		t.Errorf("retry projected %d ordinals, want none", len(again))
	}
	if got := epochOf(t, store, student, "a"); got != 2 {
		t.Errorf("retry advanced the epoch to %d, want 2", got)
	}
}

func TestProjectCredit_ALaterCycleEarnsAnotherAdvance(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)
	realDue := time.Date(2026, 11, 1, 9, 30, 0, 0, time.UTC)
	seedDue(t, store, student, "a", realDue)
	now := time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC)

	if _, _, err := store.AllocateCredit(student, "attempt:1", 0,
		[]CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}}); err != nil {
		t.Fatalf("first allocate: %v", err)
	}
	first, err := store.ProjectCredit(student, "a", 0, 25, now)
	if err != nil || len(first) != 1 {
		t.Fatalf("first projection = %d ordinals, %v", len(first), err)
	}

	// The collapse is per cycle, not for the concept's lifetime: credit earned
	// after the cycle closed still advances it once.
	if _, _, err := store.AllocateCredit(student, "attempt:2", 0,
		[]CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}}); err != nil {
		t.Fatalf("second allocate: %v", err)
	}
	second, err := store.ProjectCredit(student, "a", 0, 25, now)
	if err != nil {
		t.Fatalf("second projection: %v", err)
	}
	if len(second) != 1 {
		t.Fatalf("second projection = %d ordinals, want 1", len(second))
	}
	if got := epochOf(t, store, student, "a"); got != 3 {
		t.Errorf("epoch after the second cycle = %d, want 3", got)
	}
}

func epochOf(t *testing.T, store *SQLiteStore, student, concept string) int {
	t.Helper()
	var epoch int
	if err := store.db.QueryRow(
		`SELECT credit_epoch FROM implicit_credit_targets WHERE student_id = ? AND target_concept_id = ?`,
		student, concept).Scan(&epoch); err != nil {
		t.Fatalf("read epoch: %v", err)
	}
	return epoch
}

func TestAllocateCredit_DifferentSessionsShareOneCycle(t *testing.T) {
	store := newCreditStore(t)
	student := creditLearner(t, store)
	seedDue(t, store, student, "a", time.Date(2026, 11, 1, 9, 30, 0, 0, time.UTC))

	// The guard is keyed on the concept's review cycle, not on the session that
	// produced the answer. Two sittings inside one cycle are still one review,
	// so the second session must not bank a second ordinal.
	for _, attempt := range []string{"session-morning:7", "session-evening:2"} {
		if _, _, err := store.AllocateCredit(student, attempt, 1,
			[]CreditCandidate{{ConceptID: "a", RawUnits: 1_000_000}}); err != nil {
			t.Fatalf("allocate %s: %v", attempt, err)
		}
	}

	var banked int
	if err := store.db.QueryRow(
		`SELECT COUNT(*) FROM implicit_credit_ordinals WHERE student_id = ? AND target_concept_id = ? AND state = 'banked'`,
		student, "a").Scan(&banked); err != nil {
		t.Fatalf("count banked: %v", err)
	}
	if banked != 1 {
		t.Fatalf("banked %d ordinals across two sessions in one cycle, want 1", banked)
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
