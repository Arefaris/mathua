package probability

import "github.com/chuma-beep/mathua/internal/generator"

// Register adds the probability-theory generators to reg.
func Register(reg *generator.Registry) {
	reg.Register("prob.space.axioms", &axiomsGen{})
	reg.Register("prob.space.sigma_algebra", &sigmaAlgebraGen{})
	reg.Register("prob.space.inclusion_exclusion", &inclusionExclusionGen{})
	reg.Register("prob.space.continuity", &continuityGen{})
	reg.Register("prob.cond.multiplication", &multiplicationGen{})
	reg.Register("prob.cond.total_probability", &totalProbabilityGen{})
	reg.Register("prob.bayes.prior", &bayesGen{})
	reg.Register("prob.independence.pairwise", &independenceGen{})
	reg.Register("prob.rv.definition", &rvDefinitionGen{})
	reg.Register("prob.rv.cdf", &cdfGen{})
	reg.Register("prob.rv.pmf", &pmfGen{})
	reg.Register("prob.rv.pdf", &pdfGen{})
	reg.Register("prob.rv.joint", &jointPMFGen{})
	reg.Register("prob.rv.expectation", &expectationGen{})
	reg.Register("prob.rv.variance", &varianceGen{})
	reg.Register("prob.rv.mgf", &mgfGen{})
	reg.Register("prob.rv.transform", &transformGen{})
	reg.Register("prob.dist.order_stats", &orderStatsGen{})
	reg.Register("prob.dist.memoryless", &memorylessGen{})
	reg.Register("prob.dist.poisson_process", &poissonProcessGen{})
	reg.Register("prob.dist.exponential_family", &exponentialFamilyGen{})
	reg.Register("prob.limit.markov", &markovGen{})
	reg.Register("prob.limit.chebyshev", &chebyshevGen{})
	reg.Register("prob.limit.weak_law", &weakLawGen{})
	reg.Register("prob.limit.clt", &cltGen{})
}
