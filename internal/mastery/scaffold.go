package mastery

// RequiredKPStreak is the count of consecutive correct answers required to clear
// a Knowledge Point before advancing to the next.
// Justin Skycak, The Math Academy Way, Chapter 4 and Chapter 14.
const RequiredKPStreak = 2

// KPStateMachine tracks progress through a sequence of Knowledge Points.
type KPStateMachine struct {
	CurrentOrdinal int  `json:"current_ordinal"`
	CurrentStreak  int  `json:"current_streak"`
	TotalKPs       int  `json:"total_kps"`
	Completed      bool `json:"completed"`
}

// NewKPStateMachine initializes a state machine for a concept with totalKPs Knowledge Points.
// NewKPStateMachineAt resumes a scaffold at a known point rather than at the
// first one. A session is rebuilt from storage after a deploy, and its
// knowledge-point position is persisted with it, so restarting every learner at
// KP 1 would show them work they had already done. An ordinal past the last
// point means the sequence was finished before the restart.
func NewKPStateMachineAt(totalKPs, ordinal int) *KPStateMachine {
	if totalKPs <= 0 {
		return NewKPStateMachine(totalKPs)
	}
	if ordinal < 1 {
		ordinal = 1
	}
	sm := NewKPStateMachine(totalKPs)
	sm.CurrentOrdinal = ordinal
	if ordinal > totalKPs {
		sm.CurrentOrdinal = totalKPs
		sm.Completed = true
	}
	return sm
}

func NewKPStateMachine(totalKPs int) *KPStateMachine {
	if totalKPs <= 0 {
		return &KPStateMachine{
			CurrentOrdinal: 0,
			CurrentStreak:  0,
			TotalKPs:       0,
			Completed:      true,
		}
	}
	return &KPStateMachine{
		CurrentOrdinal: 1,
		CurrentStreak:  0,
		TotalKPs:       totalKPs,
		Completed:      false,
	}
}

// Transition evaluates an attempt.
// It returns advanced (true if this attempt advanced to the next Knowledge Point or completed)
// and completed (true if all Knowledge Points are now completed).
func (sm *KPStateMachine) Transition(correct bool) (bool, bool) {
	if sm == nil || sm.TotalKPs <= 0 {
		return false, sm != nil && sm.Completed
	}
	if sm.Completed && correct {
		return false, true
	}
	if !correct {
		// A miss invalidates the point being worked on, and a completed sequence
		// reopens its last point: otherwise finishing the last Knowledge Point
		// once would be permanent, and a learner who has since forgotten the
		// material would keep the completion forever. The ordinal itself is not
		// decremented: a miss on a point already started is what the streak
		// reset records, and stepping back a point on every slip would push a
		// learner away from the material they just reached.
		sm.CurrentStreak = 0
		if sm.Completed {
			sm.Completed = false
			sm.CurrentOrdinal = sm.TotalKPs
		}
		return false, false
	}
	sm.CurrentStreak++
	if sm.CurrentStreak >= RequiredKPStreak {
		if sm.CurrentOrdinal >= sm.TotalKPs {
			sm.Completed = true
			return true, true
		}
		sm.CurrentOrdinal++
		sm.CurrentStreak = 0
		return true, false
	}
	return false, false
}
