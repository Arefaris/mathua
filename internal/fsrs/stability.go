package fsrs

import "math"

// InitialStability is the stability assigned after the first review of a
// concept, given its rating.
func (p Parameters) InitialStability(rating Rating) float64 {
	if !rating.Valid() {
		return StabilityMin
	}
	return clampStability(p[rating-1])
}

// NextStability is the stability after a review.
//
// For an incorrect answer (Again) it is the forgetting formula, capped so that
// a lapse can never leave the concept more stable than it was. For a correct
// answer it is the recall formula, where Hard and Easy scale the gain by w15
// and w16 and Good scales by nothing.
//
// There is no same-day branch. See the package documentation: Mathua practice
// is not a scheduled review, so the short-term formula must not be reachable
// from here.
func (p Parameters) NextStability(stability, difficulty, retrievability float64, rating Rating) float64 {
	s := clampStability(stability)
	d := clampDifficulty(difficulty)
	if math.IsNaN(retrievability) {
		retrievability = DefaultRetention
	}
	if retrievability < 0 {
		retrievability = 0
	}
	if retrievability > 1 {
		retrievability = 1
	}
	if !rating.Valid() {
		return s
	}

	hard, easy := 1.0, 1.0
	if rating == Hard {
		hard = p[15]
	}
	if rating == Easy {
		easy = p[16]
	}
	if rating == Again {
		forget := p[11] * math.Pow(d, -p[12]) *
			(math.Pow(s+1, p[13]) - 1) *
			math.Exp((1-retrievability)*p[14]) * hard * easy
		cap := s / math.Exp(p[17]*p[18])
		return clampStability(math.Min(forget, cap))
	}
	recall := s * (1 + math.Exp(p[8])*(11-d)*math.Pow(s, -p[9])*
		(math.Exp((1-retrievability)*p[10])-1)*hard*easy)
	return clampStability(recall)
}

func clampStability(s float64) float64 {
	if math.IsNaN(s) || s < StabilityMin {
		return StabilityMin
	}
	if s > MaxStabilityDays {
		return MaxStabilityDays
	}
	return s
}
