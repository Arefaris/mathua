package storage

import (
	"testing"
	"time"
)

// TestRecordXPEvent_AddsToTheLedgerOnly pins the contract the weekly replay
// depends on: restoring history must not award XP a second time. The counters
// already include the original award, so this writes the event and leaves
// xp_total and xp_today untouched, and a repeat insert is refused by the key.
func TestRecordXPEvent_AddsToTheLedgerOnly(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateStudent("ada")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	since := time.Now().UTC().AddDate(0, 0, -1)

	inserted, err := store.RecordXPEvent(st.ID, 7, "replay", "attempt:1001")
	if err != nil || !inserted {
		t.Fatalf("first insert = %v, %v; want true, nil", inserted, err)
	}
	again, err := store.RecordXPEvent(st.ID, 7, "replay", "attempt:1001")
	if err != nil {
		t.Fatalf("repeat insert: %v", err)
	}
	if again {
		t.Error("repeat insert was accepted; the event key is not unique")
	}

	weekly, err := store.WeeklyXPByStudent(since)
	if err != nil {
		t.Fatalf("weekly: %v", err)
	}
	if weekly[st.ID] != 7 {
		t.Errorf("weekly XP = %d, want 7 (a duplicate must not double it)", weekly[st.ID])
	}
	total, today, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 0 || today != 0 {
		t.Errorf("counters moved to total=%d today=%d; a replay must leave them alone", total, today)
	}
}

// TestAttemptsSince_ReturnsOnlyTheWindow pins the query the replay uses: an
// attempt from last week must not be replayed into this week's board, and an
// RFC3339 timestamp with a non-UTC offset must still be compared correctly.
func TestAttemptsSince_ReturnsOnlyTheWindow(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateStudent("ada")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := store.EnsureSession("s1", st.ID); err != nil {
		t.Fatalf("session: %v", err)
	}
	now := time.Now().UTC()
	entry := AttemptEntry{
		SessionID: "s1", StudentID: st.ID, ConceptID: "arith.add.single",
		Answer: "1", Expected: "1", Correct: true, ElapsedSeconds: 3,
		Source: "practice",
	}
	old := entry
	old.Timestamp = now.AddDate(0, 0, -8)
	if _, err := store.RecordAttempt(old); err != nil {
		t.Fatalf("old attempt: %v", err)
	}
	fresh := entry
	fresh.Timestamp = now.Add(-time.Hour)
	freshID, err := store.RecordAttempt(fresh)
	if err != nil {
		t.Fatalf("fresh attempt: %v", err)
	}

	rows, err := store.AttemptsSince(now.AddDate(0, 0, -7))
	if err != nil {
		t.Fatalf("attempts since: %v", err)
	}
	if len(rows) != 1 || rows[0].ID != freshID {
		t.Fatalf("window rows = %+v, want only attempt %d", rows, freshID)
	}
	if rows[0].StudentID != st.ID || !rows[0].Correct || rows[0].Source != "practice" {
		t.Errorf("row fields = %+v", rows[0])
	}
}
