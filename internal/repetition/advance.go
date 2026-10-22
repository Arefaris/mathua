package repetition

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/chuma-beep/mathua/internal/scheduler"
)

// PolicyImplicitAdvance is the name of the counterfactual scheduling rule used
// by the shadow projection. It is deliberately not called SM-2.
//
// Ordinary SM-2 advances a repetition count and an ease factor together with
// the interval. This policy advances only the interval, because the implicit
// credit is a claim about rehearsal, not a recorded recall: the repetition
// count is history and the ease factor is a property of the learner's answers,
// so neither may move. A reader who sees "SM-2" would assume both moved.
const PolicyImplicitAdvance = "implicit-advance"

const (
	MinStepDays = 1
	MaxStepDays = 36500
)

var (
	ErrStepRange = errors.New("repetition: step outside the permitted range")
	ErrDueRange  = errors.New("repetition: due date outside the representable range")
)

// NormalizeSpeed returns a usable learning speed for the projection.
//
// The production path leaves the speed at 1.0 when GetTopicSpeed fails, and
// neither ComputeSM2WithSpeed nor UpdateLearningSpeed guards against NaN: a NaN
// speed would compare false against both clamp bounds and flow into
// math.Round, producing an implementation-defined interval. Normalising here
// keeps that fault out of the ledger.
func NormalizeSpeed(speed float64) float64 {
	if math.IsNaN(speed) || math.IsInf(speed, 0) || speed <= 0 {
		return 1.0
	}
	if speed < 0.5 {
		return 0.5
	}
	if speed > 2.0 {
		return 2.0
	}
	return speed
}

// StepDays is the number of days one implicit repetition would add under
// PolicyImplicitAdvance, derived from the same function the explicit review
// path uses with quality 4 (a correct answer at the threshold).
func StepDays(prev scheduler.SM2, speed float64) (int, error) {
	next := scheduler.ComputeSM2WithSpeed(prev, 4, NormalizeSpeed(speed))
	if next.Interval < MinStepDays {
		return MinStepDays, nil
	}
	if next.Interval > MaxStepDays {
		return 0, fmt.Errorf("%w: %d", ErrStepRange, next.Interval)
	}
	return next.Interval, nil
}

// ChainDue projects ordinals in order from a single baseline.
//
// The first ordinal starts at the real due date; each later ordinal starts at
// the previous counterfactual, so three ordinals compose rather than each
// measuring from the same origin. Every date keeps the baseline's time of day
// and is never earlier than one day after now.
func ChainDue(baseline, now time.Time, steps []int) ([]time.Time, error) {
	out := make([]time.Time, 0, len(steps))
	due := baseline
	floor := now.UTC().AddDate(0, 0, MinStepDays)
	for _, step := range steps {
		if step < MinStepDays || step > MaxStepDays {
			return nil, fmt.Errorf("%w: %d", ErrStepRange, step)
		}
		if due.Year() > 9000 {
			return nil, fmt.Errorf("%w: %s", ErrDueRange, due.Format(time.RFC3339))
		}
		due = due.AddDate(0, 0, step)
		if due.Before(floor) {
			due = floor
		}
		out = append(out, due)
	}
	return out, nil
}
