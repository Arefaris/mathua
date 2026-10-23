package fsrs

import "math"

// Retrievability is the probability that the learner still recalls a concept
// after elapsedDays, given a stability that equals the interval at which recall
// is 90 percent.
//
// R(t, S) = (1 + FACTOR * t / S) ^ DECAY, with DECAY negative, so R falls
// towards zero as t grows and R(S, S) is exactly 0.9.
func (p Parameters) Retrievability(elapsedDays, stability float64) float64 {
	if stability <= 0 || math.IsNaN(stability) {
		stability = StabilityMin
	}
	if elapsedDays <= 0 {
		// A review cannot be more than perfectly recalled, and a negative
		// elapsed time is a clock anomaly rather than evidence of forgetting.
		return 1
	}
	return math.Pow(1+p.Factor()*elapsedDays/stability, p.Decay())
}

// Interval is the gap at which the learner is predicted to recall the concept
// with DefaultRetention (0.9).
func (p Parameters) Interval(stability float64) int {
	return p.IntervalForRetention(stability, DefaultRetention)
}

// IntervalForRetention is the inverse of Retrievability: the gap after which
// recall is predicted to be retention, rounded to whole days and clamped to
// [1, MaxIntervalDays].
//
// The clamp is Mathua's, not the reference implementation's, and it is on
// purpose. A computed gap of less than a day is not a schedule, and a gap
// beyond a century is not a plan. Retention outside (0, 1) is not meaningful
// and falls back to the default rather than producing a negative or infinite
// interval.
func (p Parameters) IntervalForRetention(stability, retention float64) int {
	if stability <= 0 || math.IsNaN(stability) {
		return 1
	}
	if !(retention > 0 && retention < 1) {
		retention = DefaultRetention
	}
	raw := (stability / p.Factor()) * (math.Pow(retention, 1/p.Decay()) - 1)
	// RoundToEven matches the reference implementation's Python round, which
	// rounds halves to even. math.Round would differ on exact halves and make
	// a boundary test fight the oracle.
	days := int(math.RoundToEven(raw))
	if days < 1 {
		return 1
	}
	if days > MaxIntervalDays {
		return MaxIntervalDays
	}
	return days
}
