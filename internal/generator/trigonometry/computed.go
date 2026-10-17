package trigonometry

// The computed trigonometry concepts the domain was missing: a coterminal angle, the quadrant of an
// angle, arc length and sector area as multiples of pi, the sine and cosine of a doubled angle, the
// cofunction of an angle, a triangle's area from its three sides, and the amplitude of a sine wave.
//
// The domain was thirty-two concepts and, before this, its `numeric` surface was built from
// Pythagorean triples read off a table: given a ratio, name the other one. These eight compute an
// answer from parameters the learner has not seen, which is the shift section 4.2 of
// docs/content-expansion.md measures.
//
// Every answer is an integer or an exact reduced fraction by construction: the coterminal angle is
// the remainder the turn count leaves, the quadrant is the count of quarter-turns, the radius is a
// multiple of the angle's denominator so the pi cancels, the doubled angle's fraction reduces by a
// gcd, the cofunction complements to 90 degrees, Heron's formula is used only on triples whose area
// is whole, and the amplitude is the absolute value of a coefficient.
import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerComputed(reg *generator.Registry) {
	reg.Register("trig.basics.coterminal", &coterminalGen{})
	reg.Register("trig.basics.quadrant", &quadrantGen{})
	reg.Register("trig.basics.arc_length", &arcLengthGen{})
	reg.Register("trig.basics.sector_area", &sectorAreaGen{})
	reg.Register("trig.ident.double_angle", &doubleAngleGen{})
	reg.Register("trig.ident.cofunction", &cofunctionGen{})
	reg.Register("trig.adv.heron_area", &heronAreaGen{})
	reg.Register("trig.graph.amplitude", &amplitudeGen{})
}

func trigSmallInt(lo, hi int) int {
	if hi < lo {
		hi = lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// piFrac renders k*pi/m in the compact form the corpus uses, with the coefficient only when it is
// not one.
func piFrac(k, m int) string {
	if m == 1 {
		if k == 1 {
			return "\\pi"
		}
		return fmt.Sprintf("%d\\pi", k)
	}
	if k == 1 {
		return fmt.Sprintf("\\pi/%d", m)
	}
	return fmt.Sprintf("%d\\pi/%d", k, m)
}

// coterminalGen reduces an angle the learner is handed outside the first turn.
type coterminalGen struct{}

func (g *coterminalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	turns := 1 + rand.Intn(3+int(ctx.Difficulty*3))
	base := 1 + rand.Intn(359)
	angle := base + 360*turns
	dir := "positive"
	if rand.Intn(2) == 0 {
		angle = base - 360*turns
		dir = "negative"
	}
	return generator.Problem{
		Question: fmt.Sprintf("Find the angle in \\([0^{\\circ}, 360^{\\circ})\\) that is coterminal with \\(%d^{\\circ}\\).", angle),
		Answer:   strconv.Itoa(base),
		Explanation: fmt.Sprintf("A %s angle of %d degrees and %d degrees differ by %d full turns, and coterminal angles differ by whole turns. Adding or subtracting \\(360^{\\circ}\\) repeatedly leaves %d\\(^{\\circ}\\) in the first turn, so that is the coterminal angle.",
			dir, angle, base, turns, base),
		Facts: map[string]string{
			"value": strconv.Itoa(base),
			"given": strconv.Itoa(angle),
			"turns": strconv.Itoa(turns),
		},
	}
}

// quadrantGen counts the quarter-turns an angle crosses.
type quadrantGen struct{}

func (g *quadrantGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	angle := 1 + rand.Intn(359)
	for angle%90 == 0 {
		angle = 1 + rand.Intn(359)
	}
	q := angle/90 + 1
	return generator.Problem{
		Question: fmt.Sprintf("In which quadrant does an angle of \\(%d^{\\circ}\\) terminate? Answer 1, 2, 3 or 4.", angle),
		Answer:   strconv.Itoa(q),
		Explanation: fmt.Sprintf("The quadrant boundaries are the multiples of 90 degrees, and %d degrees lies between %d and %d, which is quadrant %d. The angle is measured counterclockwise from the positive x-axis, so the first quadrant spans 0 to 90 degrees and each next quadrant adds another 90.",
			angle, (q-1)*90, q*90, q),
		Facts: map[string]string{"value": strconv.Itoa(q), "angle": strconv.Itoa(angle)},
	}
}

// arcLengthGen uses s = r * theta with theta a fraction of pi a radius is chosen to cancel.
type arcLengthGen struct{}

func (g *arcLengthGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	dens := []int{2, 3, 4, 6}
	m := dens[rand.Intn(len(dens))]
	k := 1 + rand.Intn(m-1)
	for gcd(k, m) != 1 {
		k = 1 + rand.Intn(m-1)
	}
	t := 1 + rand.Intn(4+int(ctx.Difficulty*4))
	r := m * t // chosen so r*k/m is whole and the pi stays symbolic
	n := k * t
	frac := piFrac(k, m)
	return generator.Problem{
		Question: fmt.Sprintf("A circle has radius %d. A central angle of \\(%s\\) subtends an arc. What is the arc length, written as a multiple of \\(\\pi\\)?", r, frac),
		Answer:   strconv.Itoa(n),
		Explanation: fmt.Sprintf("Arc length is \\(s = r\\theta\\), and in radians that is the radius times the angle. Here \\(s = %d \\times %s = %d\\pi\\), so the arc length is \\(%d\\pi\\) and the multiplier is %d. Radians are the unit that makes this formula free of a conversion factor, which is why arc length is stated in radians and not degrees.",
			r, frac, n, n, n),
		Facts: map[string]string{"value": strconv.Itoa(n), "r": strconv.Itoa(r), "num": strconv.Itoa(k), "den": strconv.Itoa(m), "frac": frac},
	}
}

// sectorAreaGen uses A = r^2 theta / 2 with the same radius trick.
type sectorAreaGen struct{}

func (g *sectorAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	dens := []int{2, 4, 6}
	m := dens[rand.Intn(len(dens))]
	k := 1 + rand.Intn(m-1)
	for gcd(k, m) != 1 {
		k = 1 + rand.Intn(m-1)
	}
	t := 1 + rand.Intn(3+int(ctx.Difficulty*3))
	r := m * t
	n := (m / 2) * t * t * k // A = (1/2) r^2 (k pi / m) = m t^2 k pi / 2
	frac := piFrac(k, m)
	return generator.Problem{
		Question: fmt.Sprintf("A circle has radius %d. A sector has central angle \\(%s\\). What is the sector's area, written as a multiple of \\(\\pi\\)?", r, frac),
		Answer:   strconv.Itoa(n),
		Explanation: fmt.Sprintf("Sector area is half the radius squared times the angle in radians: \\(A = \\tfrac{1}{2} r^2 \\theta\\). Here \\(A = \\tfrac{1}{2}(%d)^2 \\cdot %s = %d\\pi\\), so the multiplier is %d. The factor of one half is the same one that makes a triangle half its bounding rectangle, and the whole circle is the special case where the angle is \\(2\\pi\\) and the area is \\(\\pi r^2\\).",
			r, frac, n, n),
		Facts: map[string]string{"value": strconv.Itoa(n), "r": strconv.Itoa(r), "num": strconv.Itoa(k), "den": strconv.Itoa(m), "frac": frac},
	}
}

// doubleAngleGen applies sin(2t) = 2 sin t cos t and cos(2t) = cos^2 t - sin^2 t to a triple.
type doubleAngleGen struct{}

func (g *doubleAngleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type triple struct{ opp, adj, hyp int }
	triples := []triple{{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {20, 21, 29}, {7, 24, 25}, {9, 40, 41}}
	tr := triples[rand.Intn(len(triples))]
	if rand.Intn(2) == 0 { // sin(2t) = 2 sin t cos t
		num := 2 * tr.opp * tr.adj
		den := tr.hyp * tr.hyp
		d := gcd(num, den)
		return generator.Problem{
			Question: fmt.Sprintf("If \\(\\sin\\theta = %d/%d\\) and \\(\\theta\\) is acute, what is \\(\\sin(2\\theta)\\), as a fraction?", tr.opp, tr.hyp),
			Answer:   fmt.Sprintf("%d/%d", num/d, den/d),
			Explanation: fmt.Sprintf("The double-angle formula is \\(\\sin(2\\theta) = 2\\sin\\theta\\cos\\theta\\). The opposite side is %d and the hypotenuse is %d, so the adjacent side is %d and \\(\\cos\\theta = %d/%d\\). Then \\(\\sin(2\\theta) = 2 \\cdot %d/%d \\cdot %d/%d = %d/%d\\). Because \\(\\theta\\) is acute both ratios are positive, which fixes the sign of the product.",
				tr.opp, tr.hyp, tr.adj, tr.adj, tr.hyp, tr.opp, tr.hyp, tr.adj, tr.hyp, num/d, den/d),
			Facts: map[string]string{"num": strconv.Itoa(num / d), "den": strconv.Itoa(den / d), "sin_num": strconv.Itoa(tr.opp), "sin_den": strconv.Itoa(tr.hyp)},
		}
	}
	// cos(2t) = cos^2 t - sin^2 t
	cn := tr.adj*tr.adj - tr.opp*tr.opp
	sign := ""
	if cn < 0 {
		sign = "-"
		cn = -cn
	}
	den := tr.hyp * tr.hyp
	d := gcd(cn, den)
	cn, den = cn/d, den/d
	return generator.Problem{
		Question: fmt.Sprintf("If \\(\\sin\\theta = %d/%d\\) and \\(\\theta\\) is acute, what is \\(\\cos(2\\theta)\\), as a fraction?", tr.opp, tr.hyp),
		Answer:   fmt.Sprintf("%s%d/%d", sign, cn, den),
		Explanation: fmt.Sprintf("The double-angle formula is \\(\\cos(2\\theta) = \\cos^2\\theta - \\sin^2\\theta\\). With opposite %d and hypotenuse %d the adjacent side is %d, so \\(\\cos\\theta = %d/%d\\). Then \\(\\cos(2\\theta) = %d/%d - %d/%d = %s%d/%d\\). The sign follows from the formula rather than from the acute angle: doubling can carry \\(2\\theta\\) past 90 degrees, where the cosine turns negative.",
			tr.opp, tr.hyp, tr.adj, tr.adj, tr.hyp, tr.adj*tr.adj, tr.hyp*tr.hyp, tr.opp*tr.opp, tr.hyp*tr.hyp, sign, cn, den),
		Facts: map[string]string{"num": sign + strconv.Itoa(cn), "den": strconv.Itoa(den), "sin_num": strconv.Itoa(tr.opp), "sin_den": strconv.Itoa(tr.hyp)},
	}
}

// cofunctionGen uses the fact that the sine of an angle is the cosine of its complement.
type cofunctionGen struct{}

func (g *cofunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	theta := 1 + rand.Intn(88) // 1..88 keeps the complement strictly inside (0, 90)
	phi := 90 - theta
	return generator.Problem{
		Question: fmt.Sprintf("Given that \\(\\sin(%d^{\\circ}) = \\cos(\\varphi)\\) with \\(0^{\\circ} < \\varphi < 90^{\\circ}\\), what is \\(\\varphi\\)?", theta),
		Answer:   strconv.Itoa(phi),
		Explanation: fmt.Sprintf("The cofunction identity says a trigonometric function of an angle equals the cofunction of its complement: \\(\\sin\\theta = \\cos(90^{\\circ} - \\theta)\\). Here \\(90 - %d = %d\\), so \\(\\varphi = %d^{\\circ}\\). That complement is why cosine is the co-sine of sine, and it is a statement about the two acute angles of a right triangle summing to 90 degrees.",
			theta, phi, phi),
		Facts: map[string]string{"value": strconv.Itoa(phi), "theta": strconv.Itoa(theta)},
	}
}

// heronAreaGen uses triples whose area by Heron's formula is whole.
type heronAreaGen struct{}

func (g *heronAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type tri struct{ a, b, c, area int }
	// Heronian triples: integer sides and integer area.
	table := []tri{
		{3, 4, 5, 6},
		{5, 12, 13, 30},
		{6, 8, 10, 24},
		{8, 15, 17, 60},
		{7, 24, 25, 84},
		{13, 14, 15, 84},
		{10, 24, 26, 120},
		{15, 20, 25, 150},
		{9, 40, 41, 180},
	}
	tr := table[rand.Intn(len(table))]
	semi := (tr.a + tr.b + tr.c) / 2
	return generator.Problem{
		Question: fmt.Sprintf("A triangle has side lengths %d, %d and %d. What is its area?", tr.a, tr.b, tr.c),
		Answer:   strconv.Itoa(tr.area),
		Explanation: fmt.Sprintf("Heron's formula is \\(A = \\sqrt{s(s-a)(s-b)(s-c)}\\) with \\(s\\) the semiperimeter. Here the semiperimeter is %d, so \\(A = \\sqrt{%d \\cdot %d \\cdot %d \\cdot %d} = %d\\). Three sides determine the triangle, so an area is recoverable from them even though no height was given; the law of cosines is what lets Heron's formula be derived from the sides.",
			semi, semi, semi-tr.a, semi-tr.b, semi-tr.c, tr.area),
		Facts: map[string]string{"value": strconv.Itoa(tr.area), "a": strconv.Itoa(tr.a), "b": strconv.Itoa(tr.b), "c": strconv.Itoa(tr.c)},
	}
}

// amplitudeGen reads the amplitude off the coefficient, sign and all.
type amplitudeGen struct{}

func (g *amplitudeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := trigSmallInt(2, 8+int(ctx.Difficulty*4))
	sign := ""
	if rand.Intn(2) == 0 {
		a = -a
		sign = "-"
	}
	fn := "sin"
	if rand.Intn(2) == 0 {
		fn = "cos"
	}
	amp := a
	if amp < 0 {
		amp = -amp
	}
	return generator.Problem{
		Question: fmt.Sprintf("What is the amplitude of \\(y = %s%d\\,%s(x)\\)?", sign, amp, fn),
		Answer:   strconv.Itoa(amp),
		Explanation: fmt.Sprintf("The amplitude is the distance from the midline to a peak, which is the absolute value of the coefficient: \\(|%d| = %d\\). The sign is not part of the amplitude because a reflection over the axis changes where the wave starts, not how far it swings. The coefficient never changes the period either, which is set by the multiplier on \\(x\\) alone.",
			a, amp),
		Facts: map[string]string{"value": strconv.Itoa(amp), "a": strconv.Itoa(a), "fn": fn, "fnmath": "\\" + fn},
	}
}
