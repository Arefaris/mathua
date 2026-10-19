package numerical

import "github.com/chuma-beep/mathua/internal/generator"

// Register adds the numerical-analysis generators to reg.
func Register(reg *generator.Registry) {
	reg.Register("numa.errors.absolute", &absRelErrorGen{})
	reg.Register("numa.errors.roundoff", &roundoffGen{})
	reg.Register("numa.errors.cancellation", &cancellationGen{})
	reg.Register("numa.errors.conditioning", &conditioningGen{})
	reg.Register("numa.errors.stability", &stabilityGen{})
	reg.Register("numa.errors.taylor", &taylorErrorGen{})
	reg.Register("numa.roots.bisection", &bisectionGen{})
	reg.Register("numa.roots.fixed_point", &fixedPointGen{})
	reg.Register("numa.roots.newton", &newtonGen{})
	reg.Register("numa.roots.secant", &secantGen{})
	reg.Register("numa.roots.false_position", &falsePositionGen{})
	reg.Register("numa.roots.order", &orderGen{})
	reg.Register("numa.roots.multiple", &multipleRootGen{})
	reg.Register("numa.roots.bracketing", &bracketingGen{})
	reg.Register("numa.roots.horner", &hornerGen{})
	reg.Register("numa.interp.lagrange", &lagrangeGen{})
	reg.Register("numa.interp.divided_differences", &dividedDiffGen{})
	reg.Register("numa.interp.error", &interpErrorGen{})
	reg.Register("numa.interp.chebyshev", &chebyshevGen{})
	reg.Register("numa.interp.splines", &splinesGen{})
}
