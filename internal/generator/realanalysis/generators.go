package realanalysis

import "github.com/chuma-beep/mathua/internal/generator"

// Register adds the real-analysis generators to reg.
func Register(reg *generator.Registry) {
	// The order and the completeness of the real line.
	reg.Register("ra.order.supremum", &supremumGen{})
	reg.Register("ra.order.completeness", &completenessGen{})
	reg.Register("ra.order.archimedean", &archimedeanGen{})
	reg.Register("ra.order.nested_intervals", &nestedIntervalsGen{})
	reg.Register("ra.order.density", &densityGen{})
	reg.Register("ra.order.countable", &countableGen{})
	reg.Register("ra.order.uncountable", &uncountableGen{})
	reg.Register("ra.order.extended", &extendedGen{})
	// Sequences and series.
	reg.Register("ra.seq.convergence", &sequenceGen{})
	reg.Register("ra.seq.monotone", &monotoneGen{})
	reg.Register("ra.seq.subsequence", &subsequenceGen{})
	reg.Register("ra.seq.bolzano_weierstrass", &bolzanoWeierstrassGen{})
	reg.Register("ra.seq.cauchy", &cauchyGen{})
	reg.Register("ra.seq.series", &seriesGen{})
	reg.Register("ra.seq.tests", &ratioTestGen{})
	reg.Register("ra.seq.absolute", &absoluteGen{})
	reg.Register("ra.seq.rearrangement", &rearrangementGen{})
	reg.Register("ra.seq.cauchy_criterion", &cauchyCriterionGen{})
	// Limits and continuity.
	reg.Register("ra.limit.function", &functionLimitGen{})
	reg.Register("ra.limit.sequential", &sequentialGen{})
	reg.Register("ra.limit.continuity", &continuityGen{})
	reg.Register("ra.limit.ivt", &ivtGen{})
	reg.Register("ra.limit.evt", &evtGen{})
	reg.Register("ra.limit.uniform", &uniformGen{})
	reg.Register("ra.limit.compact", &compactGen{})
}
