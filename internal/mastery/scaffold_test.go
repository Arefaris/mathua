package mastery

import "testing"

func TestKPStateMachine_TransitionsAndStreaks(t *testing.T) {
	sm := NewKPStateMachine(3)
	if sm.CurrentOrdinal != 1 {
		t.Fatalf("expected initial ordinal 1, got %d", sm.CurrentOrdinal)
	}
	if sm.CurrentStreak != 0 {
		t.Fatalf("expected initial streak 0, got %d", sm.CurrentStreak)
	}

	// Attempt 1: correct. Streak becomes 1. Does not advance yet.
	adv, comp := sm.Transition(true)
	if adv || comp {
		t.Fatalf("expected adv=false comp=false, got adv=%v comp=%v", adv, comp)
	}
	if sm.CurrentOrdinal != 1 || sm.CurrentStreak != 1 {
		t.Fatalf("expected ordinal=1 streak=1, got ordinal=%d streak=%d", sm.CurrentOrdinal, sm.CurrentStreak)
	}

	// Attempt 2: correct. Streak reaches 2. Advances to ordinal 2.
	adv, comp = sm.Transition(true)
	if !adv || comp {
		t.Fatalf("expected adv=true comp=false, got adv=%v comp=%v", adv, comp)
	}
	if sm.CurrentOrdinal != 2 || sm.CurrentStreak != 0 {
		t.Fatalf("expected ordinal=2 streak=0, got ordinal=%d streak=%d", sm.CurrentOrdinal, sm.CurrentStreak)
	}
}

func TestKPStateMachine_StreakResetOnMiss(t *testing.T) {
	sm := NewKPStateMachine(2)

	// 1 correct answer sets streak to 1.
	sm.Transition(true)
	if sm.CurrentStreak != 1 {
		t.Fatalf("expected streak 1, got %d", sm.CurrentStreak)
	}

	// Incorrect answer resets streak to 0 and stays on current ordinal.
	adv, comp := sm.Transition(false)
	if adv || comp {
		t.Fatalf("expected adv=false comp=false on miss, got adv=%v comp=%v", adv, comp)
	}
	if sm.CurrentOrdinal != 1 || sm.CurrentStreak != 0 {
		t.Fatalf("expected ordinal=1 streak=0 after miss, got ordinal=%d streak=%d", sm.CurrentOrdinal, sm.CurrentStreak)
	}

	// Then two consecutive correct answers advance to ordinal 2.
	sm.Transition(true)
	adv, comp = sm.Transition(true)
	if !adv || comp {
		t.Fatalf("expected advance to KP 2, got adv=%v comp=%v", adv, comp)
	}
	if sm.CurrentOrdinal != 2 {
		t.Fatalf("expected ordinal 2, got %d", sm.CurrentOrdinal)
	}
}

func TestKPStateMachine_CompletionAndMastery(t *testing.T) {
	sm := NewKPStateMachine(2)

	// Clear KP 1.
	sm.Transition(true)
	sm.Transition(true)
	if sm.CurrentOrdinal != 2 {
		t.Fatalf("expected ordinal 2, got %d", sm.CurrentOrdinal)
	}

	// 1 correct on KP 2.
	adv, comp := sm.Transition(true)
	if adv || comp {
		t.Fatalf("expected not completed yet")
	}

	// 2nd consecutive correct on KP 2 completes all Knowledge Points.
	adv, comp = sm.Transition(true)
	if !adv || !comp {
		t.Fatalf("expected adv=true comp=true, got adv=%v comp=%v", adv, comp)
	}
	if !sm.Completed {
		t.Fatalf("expected sm.Completed=true")
	}

	// Additional attempts keep completed state.
	adv, comp = sm.Transition(true)
	if adv || !comp {
		t.Fatalf("expected completed stays true, got adv=%v comp=%v", adv, comp)
	}
}

func TestKPStateMachine_MachineNextIntegration(t *testing.T) {
	m := &Machine{}
	sm := NewKPStateMachine(2)

	// Initial wrong answer leaves status UNSEEN.
	sm.Transition(false)
	ctx := TransitionCtx{KPState: sm}
	if status := m.Next(StatusUnseen, ctx); status != StatusUnseen {
		t.Fatalf("expected UNSEEN after wrong answer, got %s", status)
	}

	// 1st correct answer moves UNSEEN to LEARNING.
	sm.Transition(true)
	ctx = TransitionCtx{KPState: sm}
	if status := m.Next(StatusUnseen, ctx); status != StatusLearning {
		t.Fatalf("expected LEARNING after 1st correct, got %s", status)
	}

	// 2nd correct answer clears KP 1 and advances to KP 2 (PRACTICING).
	sm.Transition(true)
	ctx = TransitionCtx{KPState: sm}
	if status := m.Next(StatusLearning, ctx); status != StatusPracticing {
		t.Fatalf("expected PRACTICING after advancing to KP 2, got %s", status)
	}

	// Still on KP 2 with 1 correct: stays PRACTICING.
	sm.Transition(true)
	ctx = TransitionCtx{KPState: sm}
	if status := m.Next(StatusPracticing, ctx); status != StatusPracticing {
		t.Fatalf("expected PRACTICING, got %s", status)
	}

	// 2nd correct on KP 2 completes every knowledge point. That is progress, not proof:
	// the ladder still refuses MASTERED without demonstrated evidence, so the concept
	// stays PRACTICING. This assertion previously expected MASTERED, which is the defect
	// this test now guards: scaffolding could certify a concept on evidence the ladder
	// rejects.
	sm.Transition(true)
	ctx = TransitionCtx{KPState: sm}
	if status := m.Next(StatusPracticing, ctx); status != StatusPracticing {
		t.Fatalf("expected PRACTICING after completing all KPs without evidence, got %s", status)
	}

	// The same completed state does reach MASTERED once the evidence rung is met, so the
	// ladder is authoritative rather than blocked by the scaffold.
	if status := m.Next(StatusPracticing, TransitionCtx{Evidence: demonstrated(), KPState: sm}); status != StatusMastered {
		t.Fatalf("expected MASTERED when the ladder is satisfied, got %s", status)
	}
}

func TestKPStateMachine_BackwardsCompatibility(t *testing.T) {
	m := &Machine{}

	// When KPState is nil, Machine uses evidence rungs.
	ctx := TransitionCtx{
		Evidence: started(),
		KPState:  nil,
	}
	if got := m.Next(StatusUnseen, ctx); got != StatusLearning {
		t.Fatalf("expected fallback to evidence model LEARNING, got %s", got)
	}

	ctx = TransitionCtx{
		Evidence: consistent(),
		KPState:  nil,
	}
	if got := m.Next(StatusLearning, ctx); got != StatusPracticing {
		t.Fatalf("expected fallback to evidence model PRACTICING, got %s", got)
	}
}

// A completed knowledge-point sequence must be able to regress. The scaffold is
// pacing, not proof, so a learner who misses after finishing the last point has
// to re-earn it rather than keeping the completion for the life of the row.
func TestKPStateMachine_MissReopensCompletionAndRetreats(t *testing.T) {
	sm := NewKPStateMachine(2)
	for i := 0; i < RequiredKPStreak; i++ {
		sm.Transition(true)
	}
	for i := 0; i < RequiredKPStreak; i++ {
		sm.Transition(true)
	}
	if !sm.Completed {
		t.Fatalf("expected the sequence to complete first, got %+v", sm)
	}

	sm.Transition(false)
	if sm.Completed {
		t.Fatalf("a miss must reopen the completed sequence, got %+v", sm)
	}
	if sm.CurrentOrdinal != 2 {
		t.Fatalf("expected the last point to reopen, got ordinal %d", sm.CurrentOrdinal)
	}

	sm.Transition(false)
	if sm.CurrentOrdinal != 2 {
		t.Fatalf("repeated misses must stay on the point being re-earned, got %d", sm.CurrentOrdinal)
	}
	if sm.CurrentStreak != 0 {
		t.Fatalf("expected the streak to stay zero, got %d", sm.CurrentStreak)
	}

	for i := 0; i < 2*RequiredKPStreak; i++ {
		sm.Transition(true)
	}
	if !sm.Completed {
		t.Fatalf("expected recovery to complete the sequence again, got %+v", sm)
	}
}

func TestKPStateMachine_ResumesWhereTheSessionLeftOff(t *testing.T) {
	// A deploy empties the in-memory map, and the session row carries the
	// ordinal. Rebuilding at 1 made a learner repeat points they had finished.
	sm := NewKPStateMachineAt(3, 2)
	if sm.CurrentOrdinal != 2 || sm.Completed {
		t.Fatalf("resumed at ordinal %d completed=%v, want 2 and not completed", sm.CurrentOrdinal, sm.Completed)
	}
	if sm.TotalKPs != 3 {
		t.Fatalf("total = %d, want 3", sm.TotalKPs)
	}

	// An ordinal past the last point means the sequence finished before the
	// restart, so it comes back completed rather than reopening the last point.
	done := NewKPStateMachineAt(3, 4)
	if !done.Completed || done.CurrentOrdinal != 3 {
		t.Fatalf("completed=%v ordinal=%d, want completed at 3", done.Completed, done.CurrentOrdinal)
	}

	// Zero or negative means no position was stored, which is a fresh learner.
	if fresh := NewKPStateMachineAt(3, 0); fresh.CurrentOrdinal != 1 || fresh.Completed {
		t.Fatalf("fresh ordinal=%d completed=%v, want 1 and not completed", fresh.CurrentOrdinal, fresh.Completed)
	}
}
