package statistics

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("stat.data.read_table", &readTableGen{})
	reg.Register("stat.data.bar_graph", &barGraphGen{})
	reg.Register("stat.data.line_plot", &linePlotGen{})
	reg.Register("stat.data.mean", &meanGen{})
	reg.Register("stat.data.median", &medianGen{})
	reg.Register("stat.data.mode", &modeGen{})
	reg.Register("stat.data.range", &rangeGen{})
	reg.Register("stat.prob.sample_space", &sampleSpaceGen{})
	reg.Register("stat.prob.basic", &basicProbGen{})
	reg.Register("stat.prob.complement", &complementProbGen{})
	reg.Register("stat.prob.compound", &compoundProbGen{})
	reg.Register("stat.prob.counting", &countingGen{})

	reg.Register("stat.adv.covariance", &covarianceGen{})
	reg.Register("stat.adv.geometric_mean", &geometricMeanGen{})
	reg.Register("stat.adv.harmonic_mean", &harmonicMeanGen{})
	reg.Register("stat.adv.rms", &rmsGen{})
	reg.Register("stat.adv.variance", &varianceGen{})
	reg.Register("stat.dist.bernoulli", &bernoulliGen{})
	reg.Register("stat.dist.beta", &betaGen{})
	reg.Register("stat.dist.binomial", &binomialDistGen{})
	reg.Register("stat.dist.chi_square", &chiSquareGen{})
	reg.Register("stat.dist.exponential", &exponentialDistGen{})
	reg.Register("stat.dist.gamma", &gammaDistGen{})
	reg.Register("stat.dist.geometric", &geometricDistGen{})
	reg.Register("stat.dist.hypergeometric", &hypergeometricGen{})
	reg.Register("stat.dist.normal", &normalDistGen{})
	reg.Register("stat.func.gaussian", &gaussianFuncGen{})
	reg.Register("stat.dist.poisson", &poissonGen{})
	reg.Register("stat.dist.student_t", &studentTGen{})
	reg.Register("stat.dist.uniform", &uniformDistGen{})
	reg.Register("stat.dist.z_table", &zTableGen{})
	reg.Register("stat.infer.confidence", &confidenceGen{})
	reg.Register("stat.infer.sampling", &samplingGen{})
	reg.Register("stat.prob.bayes", &bayesGen{})
	reg.Register("stat.prob.continuous_rv", &continuousRVGen{})
	reg.Register("stat.prob.discrete_rv", &discreteRVGen{})
	reg.Register("stat.prob.expected_value", &expectedValueGen{})
	registerNumeric(reg)
	registerBatch(reg)
}

type readTableGen struct{}

func (g *readTableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	numRows := rand.Intn(max(1, scale/2)) + 3
	numCols := rand.Intn(max(1, scale/2)) + 2
	rowLabels := make([]string, numRows)
	for i := 0; i < numRows; i++ {
		items := []string{"Apples", "Oranges", "Bananas", "Grapes", "Mangoes", "Cherries", "Pears", "Plums"}
		rowLabels[i] = items[rand.Intn(len(items))]
	}
	colLabels := make([]string, numCols)
	for i := 0; i < numCols; i++ {
		units := []string{"Week 1", "Week 2", "Week 3", "Jan", "Feb", "Mar", "Mon", "Tue", "Wed"}
		colLabels[i] = units[rand.Intn(len(units))]
	}
	table := make([][]int, numRows)
	for i := 0; i < numRows; i++ {
		table[i] = make([]int, numCols)
		for j := 0; j < numCols; j++ {
			table[i][j] = rand.Intn(max(1, scale*4)) + 1
		}
	}
	pickRow := rand.Intn(numRows)
	sum := 0
	for j := 0; j < numCols; j++ {
		sum += table[pickRow][j]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-12s", ""))
	for j := 0; j < numCols; j++ {
		sb.WriteString(fmt.Sprintf("%-10s", colLabels[j]))
	}
	sb.WriteString("\n")
	for i := 0; i < numRows; i++ {
		sb.WriteString(fmt.Sprintf("%-12s", rowLabels[i]))
		for j := 0; j < numCols; j++ {
			sb.WriteString(fmt.Sprintf("%-10d", table[i][j]))
		}
		sb.WriteString("\n")
	}
	return generator.Problem{
		Question: fmt.Sprintf("According to the table, how many %s are there in total?\n\n%s", rowLabels[pickRow], sb.String()),
		Answer:   strconv.Itoa(sum),
		Explanation: fmt.Sprintf("Sum of %s: %s = %d", rowLabels[pickRow], strings.Join(func() []string {
			vals := make([]string, numCols)
			for j, v := range table[pickRow] {
				vals[j] = strconv.Itoa(v)
			}
			return vals
		}(), " + "), sum),
	}
}

type barGraphGen struct{}

func (g *barGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(3) + 3
	labels := make([]string, n)
	values := make([]int, n)
	items := []string{"Cats", "Dogs", "Birds", "Fish", "Hamsters", "Rabbits", "Turtles", "Snakes"}
	chosen := make(map[string]bool)
	for i := 0; i < n; i++ {
		var label string
		for {
			label = items[rand.Intn(len(items))]
			if !chosen[label] {
				chosen[label] = true
				break
			}
		}
		labels[i] = label
		values[i] = rand.Intn(max(1, scale*6)) + 5
	}
	maxIdx := 0
	for i := 1; i < n; i++ {
		if values[i] > values[maxIdx] {
			maxIdx = i
		}
	}
	dataParts := make([]string, n)
	for i := 0; i < n; i++ {
		dataParts[i] = fmt.Sprintf("%s: %d", labels[i], values[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("In a bar graph showing pet popularity (%s), which pet has the highest value?", strings.Join(dataParts, ", ")),
		Answer:      labels[maxIdx],
		Explanation: fmt.Sprintf("%s has the highest value of %d.", labels[maxIdx], values[maxIdx]),
	}
}

type linePlotGen struct{}

func (g *linePlotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(3) + 4
	xs := make([]int, n)
	ys := make([]int, n)
	startX := rand.Intn(5) + 1
	startY := rand.Intn(max(1, scale*4)) + 5
	for i := 0; i < n; i++ {
		xs[i] = startX + i
		if i == 0 {
			ys[i] = startY
		} else {
			delta := rand.Intn(max(1, scale)) - 3
			ys[i] = ys[i-1] + delta
			if ys[i] < 0 {
				ys[i] = 0
			}
		}
	}
	pick := rand.Intn(n)
	points := make([]string, n)
	for i := 0; i < n; i++ {
		points[i] = fmt.Sprintf("(%d,%d)", xs[i], ys[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("In the line plot with points %s, what is the y-value when x = %d?", strings.Join(points, ", "), xs[pick]),
		Answer:      strconv.Itoa(ys[pick]),
		Explanation: fmt.Sprintf("When x = %d, the corresponding y-value is %d.", xs[pick], ys[pick]),
	}
}

type meanGen struct{}

func (g *meanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(4) + 4
	vals := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(max(1, scale*8)) + 1
		sum += vals[i]
	}
	mean := float64(sum) / float64(n)
	valStrs := make([]string, n)
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the mean of: %s", strings.Join(valStrs, ", ")),
		Answer:      fmt.Sprintf("%.1f", mean),
		Explanation: fmt.Sprintf("Mean = (%s) / %d = %d / %d = %.1f", strings.Join(valStrs, " + "), n, sum, n, mean),
	}
}

type medianGen struct{}

func (g *medianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(5) + 5
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(max(1, scale*10)) + 1
	}
	sorted := make([]int, n)
	copy(sorted, vals)
	sort.Ints(sorted)
	var median string
	if n%2 == 1 {
		median = strconv.Itoa(sorted[n/2])
	} else {
		mid1, mid2 := sorted[n/2-1], sorted[n/2]
		avg := float64(mid1+mid2) / 2.0
		median = fmt.Sprintf("%.1f", avg)
	}
	valStrs := make([]string, n)
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	sortedStrs := make([]string, n)
	for i, v := range sorted {
		sortedStrs[i] = strconv.Itoa(v)
	}
	exp := fmt.Sprintf("Sorted: %s. ", strings.Join(sortedStrs, ", "))
	if n%2 == 1 {
		exp += fmt.Sprintf("Middle value (position %d) is %s.", n/2+1, median)
	} else {
		exp += fmt.Sprintf("Average of two middle values (%d and %d) = %s.", sorted[n/2-1], sorted[n/2], median)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the median of: %s", strings.Join(valStrs, ", ")),
		Answer:      median,
		Explanation: exp,
	}
}

type modeGen struct{}

func (g *modeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(5) + 5
	base := rand.Intn(max(1, scale*4)) + 1
	modeVal := base + rand.Intn(10)
	vals := make([]int, 0, n)
	modeCount := rand.Intn(3) + 2
	for i := 0; i < modeCount; i++ {
		vals = append(vals, modeVal)
	}
	otherCount := n - modeCount
	for i := 0; i < otherCount; i++ {
		v := rand.Intn(max(1, scale*6)) + 1
		for v == modeVal {
			v = rand.Intn(max(1, scale*6)) + 1
		}
		vals = append(vals, v)
	}
	rand.Shuffle(len(vals), func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	valStrs := make([]string, n)
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the mode of: %s", strings.Join(valStrs, ", ")),
		Answer:      strconv.Itoa(modeVal),
		Explanation: fmt.Sprintf("%d appears %d times, more than any other value.", modeVal, modeCount),
	}
}

type rangeGen struct{}

func (g *rangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(5) + 5
	vals := make([]int, n)
	minVal := rand.Intn(max(1, scale*4)) + 1
	maxVal := minVal + rand.Intn(max(1, scale*6)) + 10
	for i := 0; i < n; i++ {
		vals[i] = rand.Intn(maxVal-minVal+1) + minVal
	}
	actualMin, actualMax := vals[0], vals[0]
	for _, v := range vals {
		if v < actualMin {
			actualMin = v
		}
		if v > actualMax {
			actualMax = v
		}
	}
	r := actualMax - actualMin
	valStrs := make([]string, n)
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the range of: %s", strings.Join(valStrs, ", ")),
		Answer:      strconv.Itoa(r),
		Explanation: fmt.Sprintf("Range = max - min = %d - %d = %d", actualMax, actualMin, r),
	}
}

type sampleSpaceGen struct{}

func (g *sampleSpaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scenarios := []struct {
		question string
		count    int
		expl     string
	}{
		{"rolling a standard six-sided die", 6, "A die has 6 faces: 1, 2, 3, 4, 5, 6"},
		{"flipping a coin", 2, "A coin has 2 outcomes: Heads, Tails"},
		{"rolling a die and flipping a coin", 12, "6 die outcomes x 2 coin outcomes = 12"},
		{"selecting a card from a standard 52-card deck", 52, "There are 52 unique cards in a standard deck"},
		{"spinning a spinner with 4 equal sections labeled A-D", 4, "4 sections: A, B, C, D"},
	}
	pick := rand.Intn(len(scenarios))
	return generator.Problem{
		Question:    fmt.Sprintf("List all possible outcomes when %s. How many total outcomes are there?", scenarios[pick].question),
		Answer:      strconv.Itoa(scenarios[pick].count),
		Explanation: scenarios[pick].expl,
	}
}

// distinctColors returns n different colour names.
//
// The bag questions used to draw each colour independently, so a bag could read
// "2 blue, 11 blue, 11 yellow". That was not only odd prose: pick selected one
// of the two blue groups, so the stored answer was the probability of that group
// while a reader naturally sums both, and a correct answer was marked wrong.
func distinctColors(n int) []string {
	pool := []string{"red", "blue", "green", "yellow", "orange", "purple"}
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	if n > len(pool) {
		n = len(pool)
	}
	return pool[:n]
}

type basicProbGen struct{}

func (g *basicProbGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	colorNames := distinctColors(3)
	counts := make([]int, 3)
	total := 0
	for i := 0; i < 3; i++ {
		counts[i] = rand.Intn(max(4, scale*4)) + 1
		total += counts[i]
	}
	pick := rand.Intn(3)
	fav := counts[pick]
	gcd := mathutil.GCD(fav, total)
	bagParts := make([]string, 3)
	for i := 0; i < 3; i++ {
		bagParts[i] = fmt.Sprintf("%d %s", counts[i], colorNames[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A bag contains %s. What is \\(P(%s)\\)?", strings.Join(bagParts, ", "), colorNames[pick]),
		Answer:      fmt.Sprintf("%d/%d", fav/gcd, total/gcd),
		Explanation: fmt.Sprintf("\\(P(%s) = \\frac{%d}{%d} = \\frac{%d}{%d}\\)", colorNames[pick], fav, total, fav/gcd, total/gcd),
	}
}

type complementProbGen struct{}

func (g *complementProbGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	colorNames := distinctColors(3)
	counts := make([]int, 3)
	total := 0
	for i := 0; i < 3; i++ {
		counts[i] = rand.Intn(max(4, scale*4)) + 1
		total += counts[i]
	}
	pick := rand.Intn(3)
	other := total - counts[pick]
	gcd := mathutil.GCD(other, total)
	bagParts := make([]string, 3)
	for i := 0; i < 3; i++ {
		bagParts[i] = fmt.Sprintf("%d %s", counts[i], colorNames[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A bag contains %s. What is \\(P(\\text{not }%s)\\)?", strings.Join(bagParts, ", "), colorNames[pick]),
		Answer:      fmt.Sprintf("%d/%d", other/gcd, total/gcd),
		Explanation: fmt.Sprintf("\\(P(\\text{not }%s) = 1 - P(%s) = 1 - \\frac{%d}{%d} = \\frac{%d}{%d}\\)", colorNames[pick], colorNames[pick], counts[pick], total, other/gcd, total/gcd),
	}
}

type compoundProbGen struct{}

func (g *compoundProbGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	colorNames := distinctColors(3)
	counts := make([]int, 3)
	total := 0
	for i := 0; i < 3; i++ {
		counts[i] = rand.Intn(max(4, scale*4)) + 2
		total += counts[i]
	}
	pick1 := rand.Intn(3)
	pick2 := (pick1 + 1 + rand.Intn(2)) % 3
	c1, c2 := counts[pick1], counts[pick2]
	n1, n2 := colorNames[pick1], colorNames[pick2]
	num := c1 * c2
	den := total * (total - 1)
	gcd := mathutil.GCD(num, den)
	bagParts := make([]string, 3)
	for i := 0; i < 3; i++ {
		bagParts[i] = fmt.Sprintf("%d %s", counts[i], colorNames[i])
	}
	return generator.Problem{
		Question:    fmt.Sprintf("A bag contains %s. You draw one marble, then another without replacement. What is \\(P(%s \\text{ and then } %s)\\)?", strings.Join(bagParts, ", "), n1, n2),
		Answer:      fmt.Sprintf("%d/%d", num/gcd, den/gcd),
		Explanation: fmt.Sprintf("\\(P(%s \\text{ then } %s) = \\frac{%d}{%d} \\times \\frac{%d}{%d} = \\frac{%d}{%d} = \\frac{%d}{%d}\\)", n1, n2, c1, total, c2, total-1, num, den, num/gcd, den/gcd),
	}
}

type countingGen struct{}

func (g *countingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	categories := []struct {
		name    string
		options []string
	}{
		{"shirts", []string{"3", "4", "5", "6", "7"}},
		{"pants", []string{"2", "3", "4", "5"}},
		{"hats", []string{"2", "3", "4"}},
		{"shoes", []string{"2", "3"}},
	}
	n := rand.Intn(2) + 2
	items := make([]struct{ name, count string }, n)
	chosen := make(map[int]bool)
	total := 1
	for i := 0; i < n; i++ {
		var idx int
		for {
			idx = rand.Intn(len(categories))
			if !chosen[idx] {
				chosen[idx] = true
				break
			}
		}
		opt := categories[idx].options[rand.Intn(len(categories[idx].options))]
		items[i] = struct{ name, count string }{categories[idx].name, opt}
		v, _ := strconv.Atoi(opt)
		total *= v
	}
	parts := make([]string, n)
	for i, item := range items {
		parts[i] = fmt.Sprintf("%s %s", item.count, item.name)
	}
	multParts := make([]string, n)
	for i, item := range items {
		multParts[i] = item.count
	}
	return generator.Problem{
		Question:    fmt.Sprintf("You have %s. How many different outfit combinations can you make?", strings.Join(parts, ", ")),
		Answer:      strconv.Itoa(total),
		Explanation: fmt.Sprintf("Total combinations = %s = %d", strings.Join(multParts, " x "), total),
	}
}

// ----- stat.adv.* (advanced measures) -----

type covarianceGen struct{}

func (g *covarianceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Do two variables that tend to increase together have positive covariance? (yes/no)", "yes", "Positive covariance means the variables tend to increase together."},
		{"Do two variables that move in opposite directions have positive covariance? (yes/no)", "no", "Negative covariance means one tends to increase while the other decreases."},
		{"What does covariance near zero indicate?", "no linear relationship", "Zero or near-zero covariance suggests no linear relationship between variables."},
		{"Is covariance scale-dependent or scale-invariant?", "scale-dependent", "Covariance depends on the units of measurement; it is not standardized."},
		{"For a sample of size 11, what is the denominator of the sample covariance?", "10", "Sample covariance divides by (n-1) = 10, while population covariance divides by n."},
		{"Is correlation equal to covariance divided by the product of the standard deviations? (yes/no)", "yes", "Pearson correlation standardizes covariance by dividing by the product of standard deviations."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type geometricMeanGen struct{}

func (g *geometricMeanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	p := float64(a * b)
	gm := int(math.Sqrt(p))
	for gm*gm != a*b {
		a = rand.Intn(10) + 1
		b = rand.Intn(10) + 1
		p = float64(a * b)
		gm = int(math.Sqrt(p))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the geometric mean of %d and %d.", a, b),
		Answer:      fmt.Sprintf("%d", gm),
		Explanation: fmt.Sprintf("\\(\\text{Geometric mean} = \\sqrt{%d \\times %d} = \\sqrt{%d} = %d\\)", a, b, a*b, gm),
	}
}

type harmonicMeanGen struct{}

func (g *harmonicMeanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := (rand.Intn(max(6, scale*4)) + 2) * 2
	b := (rand.Intn(max(6, scale*4)) + 2) * 2
	numer := 2 * a * b
	denom := a + b
	gcd := mathutil.GCD(numer, denom)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the harmonic mean of %d and %d.", a, b),
		Answer:      fmt.Sprintf("%d/%d", numer/gcd, denom/gcd),
		Explanation: fmt.Sprintf("\\(H = \\frac{2}{\\frac{1}{%d} + \\frac{1}{%d}} = \\frac{2}{\\frac{%d}{%d}} = \\frac{%d}{%d} = \\frac{%d}{%d}\\)", a, b, a+b, a*b, 2*a*b, a+b, numer/gcd, denom/gcd),
	}
}

type rmsGen struct{}

func (g *rmsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	sq := float64(a*a + b*b)
	p := int(sq / 2)
	rms := int(math.Sqrt(float64(p)))
	for rms*rms != p {
		a = rand.Intn(8) + 1
		b = a
		sq = float64(a*a + b*b)
		p = int(sq / 2)
		rms = int(math.Sqrt(float64(p)))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the root mean square (RMS) of %d and %d.", a, b),
		Answer:      fmt.Sprintf("%d", rms),
		Explanation: fmt.Sprintf("\\(\\text{RMS} = \\sqrt{\\frac{%d^{2}+%d^{2}}{2}} = \\sqrt{\\frac{%d+%d}{2}} = \\sqrt{%d} = %d\\)", a, b, a*a, b*b, p, rms),
	}
}

type varianceGen struct{}

func (g *varianceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	vals := make([]int, 5)
	sum := 0
	for i := range vals {
		vals[i] = (rand.Intn(max(4, scale*4)) + 1) * 2
		sum += vals[i]
	}
	mean := float64(sum) / float64(len(vals))
	var ss float64
	for _, v := range vals {
		d := float64(v) - mean
		ss += d * d
	}
	popVar := ss / float64(len(vals))
	valStrs := make([]string, len(vals))
	for i, v := range vals {
		valStrs[i] = strconv.Itoa(v)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the population variance of: %s", strings.Join(valStrs, ", ")),
		Answer:      fmt.Sprintf("%.0f", popVar),
		Explanation: fmt.Sprintf("\\(\\text{Mean} = \\frac{%d}{%d} = %.0f\\). \\(\\text{Variance} = \\frac{\\Sigma(x_i - \\bar{x})^{2}}{n} = \\frac{%.0f}{%d} = %.0f\\)", sum, len(vals), mean, ss, len(vals), popVar),
	}
}

// ----- stat.dist.* (probability distributions) -----

type bernoulliGen struct{}

func (g *bernoulliGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type param struct {
			n int
			p float64
		}
		params := []param{
			{20, 0.3},
			{50, 0.4},
			{40, 0.5},
			{100, 0.7},
			{80, 0.25},
			{60, 0.6},
		}
		item := params[rand.Intn(len(params))]
		expected := int(math.Round(float64(item.n) * item.p))
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1, X_2, \\dots, X_{%d}\\) be independent Bernoulli random variables each with success probability \\(p = %.2f\\). Compute the expected total number of successes \\(E[\\sum_{i=1}^{%d} X_i]\\).", item.n, item.p, item.n),
			Answer:   fmt.Sprintf("%d", expected),
			Explanation: fmt.Sprintf("By linearity of expectation, the expected value of a sum of independent Bernoulli variables is \\(E[\\sum X_i] = %d \\times %.2f = %d\\).",
				item.n, item.p, expected),
		}
	case 1:
		type varParam struct {
			n int
			p float64
			v int
		}
		varParams := []varParam{
			{40, 0.5, 10},
			{100, 0.5, 25},
			{100, 0.2, 16},
			{100, 0.1, 9},
			{50, 0.2, 8},
			{25, 0.2, 4},
		}
		vp := varParams[rand.Intn(len(varParams))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1, X_2, \\dots, X_{%d}\\) be independent Bernoulli random variables with parameter \\(p = %.1f\\). Find the variance \\(\\operatorname{Var}(\\sum_{i=1}^{%d} X_i)\\).", vp.n, vp.p, vp.n),
			Answer:   fmt.Sprintf("%d", vp.v),
			Explanation: fmt.Sprintf("For independent Bernoulli trials, the variance of the sum is \\(\\operatorname{Var}(\\sum X_i) = %d \\times p(1-p) = %d \\times %.1f \\times %.1f = %d\\).",
				vp.n, vp.n, vp.p, 1-vp.p, vp.v),
		}
	case 2:
		pPct := 10 * (1 + rand.Intn(8))
		p := float64(pPct) / 100.0
		failPct := 100 - pPct
		return generator.Problem{
			Question: fmt.Sprintf("A trial modeled as a Bernoulli random variable succeeds with probability \\(p = %.2f\\). What is the probability (expressed as an integer percentage from 0 to 100) that the trial fails?", p),
			Answer:   fmt.Sprintf("%d", failPct),
			Explanation: fmt.Sprintf("The probability of failure is \\(q = 1 - p = 1 - %.2f = %.2f\\), which equals %d percent.",
				p, 1-p, failPct),
		}
	case 3:
		// The probability is printed as a fraction because the odds it implies
		// are then exact. Rounded to two decimals, 2/3 reads as 0.67, and a
		// learner who divides the printed numbers gets 2.03 for an answer of 2.
		oddsList := []int{1, 2, 3, 4, 9}
		k := oddsList[rand.Intn(len(oddsList))]
		return generator.Problem{
			Question: fmt.Sprintf("A Bernoulli trial has success probability \\(p = \\frac{%d}{%d}\\). Compute the odds in favor of success, defined as the ratio \\(\\frac{p}{1-p}\\).", k, k+1),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("\\(1-p = \\frac{1}{%d}\\), so \\(\\frac{p}{1-p} = \\frac{%d}{%d} \\div \\frac{1}{%d} = %d\\).",
				k+1, k, k+1, k+1, k),
		}
	default:
		type scaleParam struct {
			a   int
			p   float64
			ans int
		}
		scales := []scaleParam{
			{5, 0.4, 10},
			{10, 0.3, 30},
			{4, 0.5, 8},
			{6, 0.5, 18},
			{5, 0.2, 5},
			{10, 0.8, 80},
		}
		sp := scales[rand.Intn(len(scales))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Bernoulli}(%.1f)\\) and define \\(Y = %d X\\). Find the second moment \\(E[Y^2]\\).", sp.p, sp.a),
			Answer:   fmt.Sprintf("%d", sp.ans),
			Explanation: fmt.Sprintf("Because \\(X \\in \\{0, 1\\}\\), \\(X^2 = X\\) and \\(E[X^2] = p = %.1f\\). Therefore \\(E[Y^2] = %d^2 E[X^2] = %d \\times %.1f = %d\\).",
				sp.p, sp.a, sp.a*sp.a, sp.p, sp.ans),
		}
	}
}

type betaGen struct{}

func (g *betaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type meanParam struct {
			a, b, pct int
		}
		params := []meanParam{
			{3, 7, 30},
			{1, 3, 25},
			{7, 13, 35},
			{3, 2, 60},
			{9, 16, 36},
			{1, 1, 50},
			{4, 1, 80},
		}
		p := params[rand.Intn(len(params))]
		return generator.Problem{
			Question: fmt.Sprintf("For a random variable \\(X \\sim \\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\), compute the expected value \\(E[X]\\) expressed as an integer percentage from 0 to 100.", p.a, p.b),
			Answer:   fmt.Sprintf("%d", p.pct),
			Explanation: fmt.Sprintf("The mean of a Beta distribution is \\(E[X] = \\frac{\\alpha}{\\alpha + \\beta} = \\frac{%d}{%d} = %.2f\\), which corresponds to %d percent.",
				p.a, p.a+p.b, float64(p.pct)/100.0, p.pct),
		}
	case 1:
		type modeParam struct {
			a, b, pct int
		}
		params := []modeParam{
			{4, 8, 30},
			{2, 4, 25},
			{3, 4, 40},
			{8, 14, 35},
			{5, 2, 80},
			{3, 3, 50},
		}
		p := params[rand.Intn(len(params))]
		num := p.a - 1
		denom := p.a + p.b - 2
		return generator.Problem{
			Question: fmt.Sprintf("For a random variable \\(X \\sim \\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\) with \\(\\alpha, \\beta > 1\\), compute the mode of the distribution expressed as an integer percentage from 0 to 100.", p.a, p.b),
			Answer:   fmt.Sprintf("%d", p.pct),
			Explanation: fmt.Sprintf("The mode of a Beta distribution with parameters greater than 1 is \\(\\operatorname{Mode} = \\frac{\\alpha - 1}{\\alpha + \\beta - 2} = \\frac{%d}{%d} = %.2f\\), which equals %d percent.",
				num, denom, float64(p.pct)/100.0, p.pct),
		}
	case 2:
		a0 := 1 + rand.Intn(6)
		b0 := 1 + rand.Intn(6)
		s := 2 + rand.Intn(10)
		f := 2 + rand.Intn(10)
		aPost := a0 + s
		return generator.Problem{
			Question: fmt.Sprintf("A Bayesian prior is modeled as \\(\\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\). An experiment observes %d successes and %d failures in independent Bernoulli trials. What is the updated parameter \\(\\alpha_{\\text{post}}\\) of the posterior Beta distribution?", a0, b0, s, f),
			Answer:   fmt.Sprintf("%d", aPost),
			Explanation: fmt.Sprintf("Because the Beta prior is conjugate to the Bernoulli likelihood, the posterior parameter adds the observed successes: \\(\\alpha_{\\text{post}} = \\alpha_0 + s = %d + %d = %d\\).",
				a0, s, aPost),
		}
	case 3:
		a := 2 + rand.Intn(8)
		b := 2 + rand.Intn(8)
		total := a + b
		return generator.Problem{
			Question: fmt.Sprintf("In Bayesian statistics with a \\(\\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\) prior, compute the effective prior pseudo-sample size \\(\\alpha + \\beta\\).", a, b),
			Answer:   fmt.Sprintf("%d", total),
			Explanation: fmt.Sprintf("The sum of the shape parameters represents the effective prior sample size: \\(\\alpha + \\beta = %d + %d = %d\\).",
				a, b, total),
		}
	default:
		k := 2 + rand.Intn(6)
		b := 1 + rand.Intn(5)
		a := k * b
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Beta}(\\alpha = %d, \\beta = %d)\\). Find the ratio of expected values \\(\\frac{E[X]}{E[1 - X]}\\).", a, b),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("Because \\(E[X] = \\frac{\\alpha}{\\alpha + \\beta}\\) and \\(E[1 - X] = \\frac{\\beta}{\\alpha + \\beta}\\), their ratio is \\(\\frac{E[X]}{E[1 - X]} = \\frac{\\alpha}{\\beta} = \\frac{%d}{%d} = %d\\).",
				a, b, k),
		}
	}
}

type binomialDistGen struct{}

func (g *binomialDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type item struct {
			n int
			p float64
		}
		items := []item{
			{20, 0.25},
			{50, 0.40},
			{60, 0.15},
			{80, 0.35},
			{45, 0.20},
			{30, 0.40},
			{75, 0.60},
			{120, 0.05},
		}
		it := items[rand.Intn(len(items))]
		mean := int(math.Round(float64(it.n) * it.p))
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Binomial}(%d, %.2f)\\). Compute the expected value \\(E[X]\\).", it.n, it.p),
			Answer:   fmt.Sprintf("%d", mean),
			Explanation: fmt.Sprintf("The expected value of a binomial random variable is \\(E[X] = n p = %d \\times %.2f = %d\\).",
				it.n, it.p, mean),
		}
	case 1:
		type varItem struct {
			n int
			p float64
			v int
		}
		items := []varItem{
			{36, 0.50, 9},
			{64, 0.50, 16},
			{100, 0.50, 25},
			{144, 0.50, 36},
			{25, 0.20, 4},
			{50, 0.20, 8},
			{75, 0.20, 12},
			{100, 0.20, 16},
			{100, 0.10, 9},
			{200, 0.10, 18},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Binomial}(%d, %.2f)\\). Calculate the variance \\(\\operatorname{Var}(X)\\).", it.n, it.p),
			Answer:   fmt.Sprintf("%d", it.v),
			Explanation: fmt.Sprintf("The variance of a binomial random variable is \\(\\operatorname{Var}(X) = n p (1 - p) = %d \\times %.2f \\times %.2f = %d\\).",
				it.n, it.p, 1-it.p, it.v),
		}
	case 2:
		type sdItem struct {
			n  int
			p  float64
			v  int
			sd int
		}
		items := []sdItem{
			{36, 0.50, 9, 3},
			{64, 0.50, 16, 4},
			{100, 0.50, 25, 5},
			{144, 0.50, 36, 6},
			{25, 0.20, 4, 2},
			{100, 0.20, 16, 4},
			{225, 0.20, 36, 6},
			{100, 0.10, 9, 3},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) follows a binomial distribution \\(\\operatorname{Binomial}(%d, %.2f)\\). Find the standard deviation \\(\\sigma\\).", it.n, it.p),
			Answer:   fmt.Sprintf("%d", it.sd),
			Explanation: fmt.Sprintf("The standard deviation is \\(\\sigma = \\sqrt{n p (1 - p)} = \\sqrt{%d \\times %.2f \\times %.2f} = \\sqrt{%d} = %d\\).",
				it.n, it.p, 1-it.p, it.v, it.sd),
		}
	case 3:
		type combItem struct {
			n int
			k int
			c int
		}
		items := []combItem{
			{5, 2, 10},
			{6, 2, 15},
			{6, 3, 20},
			{7, 2, 21},
			{7, 3, 35},
			{8, 2, 28},
			{8, 3, 56},
			{9, 2, 36},
			{10, 2, 45},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("In a binomial experiment with \\(n = %d\\) trials, how many distinct sequences contain exactly \\(k = %d\\) successes?", it.n, it.k),
			Answer:   fmt.Sprintf("%d", it.c),
			Explanation: fmt.Sprintf("The number of sequences is given by the binomial coefficient \\(\\binom{n}{k} = \\binom{%d}{%d} = \\frac{%d!}{%d!(%d - %d)!} = %d\\).",
				it.n, it.k, it.n, it.k, it.n, it.k, it.c),
		}
	default:
		type sumItem struct {
			n1 int
			n2 int
			p  float64
			e  int
		}
		items := []sumItem{
			{30, 20, 0.40, 20},
			{45, 35, 0.25, 20},
			{25, 15, 0.50, 20},
			{40, 60, 0.30, 30},
			{70, 50, 0.15, 18},
			{50, 50, 0.10, 10},
		}
		it := items[rand.Intn(len(items))]
		nTotal := it.n1 + it.n2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1 \\sim \\operatorname{Binomial}(%d, %.2f)\\) and \\(X_2 \\sim \\operatorname{Binomial}(%d, %.2f)\\) be independent random variables. Find the expected value of their sum \\(E[X_1 + X_2]\\).", it.n1, it.p, it.n2, it.p),
			Answer:   fmt.Sprintf("%d", it.e),
			Explanation: fmt.Sprintf("By linearity of expectation, \\(E[X_1 + X_2] = E[X_1] + E[X_2] = (n_1 + n_2)p = (%d + %d) \\times %.2f = %d \\times %.2f = %d\\).",
				it.n1, it.n2, it.p, nTotal, it.p, it.e),
		}
	}
}

type chiSquareGen struct{}

func (g *chiSquareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		kValues := []int{3, 5, 7, 9, 11, 13, 15, 20}
		k := kValues[rand.Intn(len(kValues))]
		v := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\chi^2(%d)\\) follow a chi-square distribution with %d degrees of freedom. Calculate the variance \\(\\operatorname{Var}(X)\\).", k, k),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf("The variance of a chi-square distribution with \\(k\\) degrees of freedom is \\(\\operatorname{Var}(X) = 2k = 2 \\times %d = %d\\).",
				k, v),
		}
	case 1:
		type sdItem struct {
			k  int
			v  int
			sd int
		}
		items := []sdItem{
			{2, 4, 2},
			{8, 16, 4},
			{18, 36, 6},
			{32, 64, 8},
			{50, 100, 10},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) follows a chi-square distribution with %d degrees of freedom: \\(X \\sim \\chi^2(%d)\\). Find the standard deviation \\(\\sigma\\).", it.k, it.k),
			Answer:   fmt.Sprintf("%d", it.sd),
			Explanation: fmt.Sprintf("For a chi-square distribution, the variance is \\(\\operatorname{Var}(X) = 2k = 2 \\times %d = %d\\). The standard deviation is \\(\\sigma = \\sqrt{2k} = \\sqrt{%d} = %d\\).",
				it.k, it.v, it.v, it.sd),
		}
	case 2:
		k1List := []int{4, 6, 8, 10, 12}
		k2List := []int{5, 7, 9, 11, 15}
		k1 := k1List[rand.Intn(len(k1List))]
		k2 := k2List[rand.Intn(len(k2List))]
		sumK := k1 + k2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1 \\sim \\chi^2(%d)\\) and \\(X_2 \\sim \\chi^2(%d)\\) be independent chi-square random variables. Find the expected value of their sum \\(E[X_1 + X_2]\\).", k1, k2),
			Answer:   fmt.Sprintf("%d", sumK),
			Explanation: fmt.Sprintf("By independence and linearity of expectation, \\(E[X_1 + X_2] = E[X_1] + E[X_2] = k_1 + k_2 = %d + %d = %d\\).",
				k1, k2, sumK),
		}
	case 3:
		kValues := []int{4, 6, 8, 10, 12, 16, 25}
		k := kValues[rand.Intn(len(kValues))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(Z_1, Z_2, \\dots, Z_{%d}\\) be independent standard normal random variables. What is the degrees of freedom parameter \\(k\\) of the variable \\(W = \\sum_{i=1}^{%d} Z_i^2\\)?", k, k),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("The sum of \\(k\\) independent squared standard normal variables follows a chi-square distribution \\(\\chi^2(k)\\). Here %d variables are summed, so \\(k = %d\\).",
				k, k),
		}
	default:
		kValues := []int{3, 4, 5, 6, 7, 8, 10}
		k := kValues[rand.Intn(len(kValues))]
		m2 := k*k + 2*k
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\chi^2(%d)\\). Calculate the second moment \\(E[X^2]\\).", k),
			Answer:   fmt.Sprintf("%d", m2),
			Explanation: fmt.Sprintf("Using the variance formula \\(\\operatorname{Var}(X) = E[X^2] - (E[X])^2\\), we find \\(E[X^2] = \\operatorname{Var}(X) + (E[X])^2 = 2(%d) + (%d)^2 = %d + %d = %d\\).",
				k, k, 2*k, k*k, m2),
		}
	}
}

type exponentialDistGen struct{}

func (g *exponentialDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type rateItem struct {
			m      int
			lambda int
		}
		items := []rateItem{
			{5, 12},
			{6, 10},
			{10, 6},
			{12, 5},
			{15, 4},
			{20, 3},
			{30, 2},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Customer arrivals follow an exponential distribution with a mean waiting time of %d minutes. What is the arrival rate \\(\\lambda\\) in events per hour?", it.m),
			Answer:   fmt.Sprintf("%d", it.lambda),
			Explanation: fmt.Sprintf("With 60 minutes per hour and mean \\(E[X] = %d\\) minutes, the hourly rate parameter is \\(\\lambda = \\frac{60}{%d} = %d\\).",
				it.m, it.m, it.lambda),
		}
	case 1:
		type memItem struct {
			mu int
			s  int
		}
		items := []memItem{
			{4, 10},
			{5, 8},
			{6, 12},
			{8, 20},
			{10, 15},
			{12, 24},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A component lifetime follows an exponential distribution with mean \\(E[X] = %d\\) hours. Given that the component has operated without failure for %d hours, what is the expected remaining lifetime until failure?", it.mu, it.s),
			Answer:   fmt.Sprintf("%d", it.mu),
			Explanation: fmt.Sprintf("By the memoryless property of the exponential distribution, prior operating time does not affect remaining lifetime: \\(E[X - %d \\mid X > %d] = E[X] = %d\\).",
				it.s, it.s, it.mu),
		}
	case 2:
		l1List := []int{2, 3, 4, 5, 6}
		l2List := []int{3, 4, 5, 7, 8}
		l1 := l1List[rand.Intn(len(l1List))]
		l2 := l2List[rand.Intn(len(l2List))]
		sumL := l1 + l2
		return generator.Problem{
			Question: fmt.Sprintf("Two independent sensors monitor a machine. The time until sensor 1 triggers follows \\(\\operatorname{Exp}(\\lambda_1)\\) with rate \\(\\lambda_1 = %d\\) per day, and sensor 2 follows \\(\\operatorname{Exp}(\\lambda_2)\\) with rate \\(\\lambda_2 = %d\\) per day. What is the rate parameter \\(\\lambda_{\\min}\\) of the first sensor trigger time \\(T = \\min(X_1, X_2)\\)?", l1, l2),
			Answer:   fmt.Sprintf("%d", sumL),
			Explanation: fmt.Sprintf("The minimum of independent exponential variables is exponentially distributed with rate equal to the sum of the rates: \\(\\lambda_{\\min} = \\lambda_1 + \\lambda_2 = %d + %d = %d\\).",
				l1, l2, sumL),
		}
	case 3:
		muList := []int{3, 4, 5, 6, 7, 8, 10}
		mu := muList[rand.Intn(len(muList))]
		m2 := 2 * mu * mu
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Exp}(\\lambda)\\) have mean \\(E[X] = %d\\). Calculate the second moment \\(E[X^2]\\).", mu),
			Answer:   fmt.Sprintf("%d", m2),
			Explanation: fmt.Sprintf("For an exponential random variable, \\(\\operatorname{Var}(X) = \\mu^2 = %d^2 = %d\\). Therefore, \\(E[X^2] = \\operatorname{Var}(X) + (E[X])^2 = %d + %d = %d\\).",
				mu, mu*mu, mu*mu, mu*mu, m2),
		}
	default:
		type sumItem struct {
			n  int
			mu int
		}
		items := []sumItem{
			{4, 5},
			{5, 6},
			{6, 7},
			{8, 5},
			{10, 4},
			{12, 3},
		}
		it := items[rand.Intn(len(items))]
		total := it.n * it.mu
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1, X_2, \\dots, X_{%d}\\) be independent exponential random variables each having mean \\(E[X_i] = %d\\). Find the expected value of their sum \\(E[\\sum_{i=1}^{%d} X_i]\\).", it.n, it.mu, it.n),
			Answer:   fmt.Sprintf("%d", total),
			Explanation: fmt.Sprintf("By linearity of expectation, the expected value of the sum is \\(E[\\sum_{i=1}^{%d} X_i] = %d \\times E[X_1] = %d \\times %d = %d\\).",
				it.n, it.n, it.n, it.mu, total),
		}
	}
}

type gammaDistGen struct{}

func (g *gammaDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type meanItem struct {
			k     int
			theta int
		}
		items := []meanItem{
			{3, 4},
			{4, 5},
			{5, 2},
			{6, 3},
			{7, 4},
			{8, 5},
		}
		it := items[rand.Intn(len(items))]
		mean := it.k * it.theta
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Gamma}(k, \\theta)\\) have shape parameter \\(k = %d\\) and scale parameter \\(\\theta = %d\\). Calculate the expected value \\(E[X]\\).", it.k, it.theta),
			Answer:   fmt.Sprintf("%d", mean),
			Explanation: fmt.Sprintf("The expected value of a Gamma distribution with shape \\(k\\) and scale \\(\\theta\\) is \\(E[X] = k \\theta = %d \\times %d = %d\\).",
				it.k, it.theta, mean),
		}
	case 1:
		type varItem struct {
			k     int
			theta int
		}
		items := []varItem{
			{2, 3},
			{3, 4},
			{4, 3},
			{5, 2},
			{6, 4},
			{2, 5},
		}
		it := items[rand.Intn(len(items))]
		v := it.k * it.theta * it.theta
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Gamma}(k, \\theta)\\) have shape \\(k = %d\\) and scale \\(\\theta = %d\\). Calculate the variance \\(\\operatorname{Var}(X)\\).", it.k, it.theta),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf("The variance of a Gamma distribution with shape \\(k\\) and scale \\(\\theta\\) is \\(\\operatorname{Var}(X) = k \\theta^2 = %d \\times %d^2 = %d \\times %d = %d\\).",
				it.k, it.theta, it.k, it.theta*it.theta, v),
		}
	case 2:
		type sdItem struct {
			k     int
			theta int
			sd    int
		}
		items := []sdItem{
			{4, 3, 6},
			{4, 5, 10},
			{9, 2, 6},
			{9, 4, 12},
			{16, 2, 8},
			{16, 3, 12},
		}
		it := items[rand.Intn(len(items))]
		v := it.k * it.theta * it.theta
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) follows a Gamma distribution with shape \\(k = %d\\) and scale \\(\\theta = %d\\). Find the standard deviation \\(\\sigma\\).", it.k, it.theta),
			Answer:   fmt.Sprintf("%d", it.sd),
			Explanation: fmt.Sprintf("The variance is \\(\\operatorname{Var}(X) = k \\theta^2 = %d \\times %d^2 = %d\\). Therefore, the standard deviation is \\(\\sigma = \\sqrt{k \\theta^2} = \\sqrt{%d} = %d\\).",
				it.k, it.theta, v, v, it.sd),
		}
	case 3:
		k1List := []int{2, 3, 4, 5}
		k2List := []int{3, 4, 5, 6}
		thetaList := []int{2, 3, 4}
		k1 := k1List[rand.Intn(len(k1List))]
		k2 := k2List[rand.Intn(len(k2List))]
		th := thetaList[rand.Intn(len(thetaList))]
		sumK := k1 + k2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1 \\sim \\operatorname{Gamma}(%d, %d)\\) and \\(X_2 \\sim \\operatorname{Gamma}(%d, %d)\\) be independent random variables sharing scale parameter \\(\\theta = %d\\). What is the shape parameter \\(k_{\\text{sum}}\\) of their sum \\(S = X_1 + X_2\\)?", k1, th, k2, th, th),
			Answer:   fmt.Sprintf("%d", sumK),
			Explanation: fmt.Sprintf("The sum of independent Gamma variables sharing the same scale parameter has shape parameter equal to the sum of the individual shapes: \\(k_{\\text{sum}} = k_1 + k_2 = %d + %d = %d\\).",
				k1, k2, sumK),
		}
	default:
		type modeItem struct {
			k     int
			theta int
		}
		items := []modeItem{
			{4, 3},
			{5, 2},
			{6, 4},
			{3, 5},
			{7, 3},
			{8, 2},
		}
		it := items[rand.Intn(len(items))]
		mode := (it.k - 1) * it.theta
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X\\) follow a Gamma distribution with shape parameter \\(k = %d\\) and scale parameter \\(\\theta = %d\\). Find the mode of \\(X\\).", it.k, it.theta),
			Answer:   fmt.Sprintf("%d", mode),
			Explanation: fmt.Sprintf("For shape parameter \\(k \\ge 1\\), the mode of \\(\\operatorname{Gamma}(k, \\theta)\\) is \\((k - 1)\\theta = (%d - 1) \\times %d = %d \\times %d = %d\\).",
				it.k, it.theta, it.k-1, it.theta, mode),
		}
	}
}

type geometricDistGen struct{}

func (g *geometricDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type item struct {
			p   float64
			exp int
		}
		items := []item{
			{0.50, 2},
			{0.25, 4},
			{0.20, 5},
			{0.10, 10},
			{0.05, 20},
			{0.04, 25},
			{0.02, 50},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X\\) follow a geometric distribution with success probability \\(p = %.2f\\), counting trials until the first success. Find the expected value \\(E[X]\\).", it.p),
			Answer:   fmt.Sprintf("%d", it.exp),
			Explanation: fmt.Sprintf("The expected value of a geometric random variable is \\(E[X] = \\frac{1}{p} = \\frac{1}{%.2f} = %d\\).",
				it.p, it.exp),
		}
	case 1:
		mList := []int{2, 3, 4, 5, 6, 10}
		m := mList[rand.Intn(len(mList))]
		v := m * (m - 1)
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Geometric}(p)\\) with success probability \\(p = \\frac{1}{%d}\\). Compute the variance \\(\\operatorname{Var}(X)\\).", m),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf("The variance of a geometric distribution is \\(\\operatorname{Var}(X) = \\frac{1-p}{p^2} = \\frac{1 - 1/%d}{(1/%d)^2} = \\frac{%d/%d}{1/%d} = %d \\times %d = %d\\).",
				m, m, m-1, m, m*m, m-1, m, v),
		}
	case 2:
		type cdfItem struct {
			p   float64
			k   int
			pct int
		}
		items := []cdfItem{
			{0.50, 1, 50},
			{0.50, 2, 75},
			{0.20, 1, 20},
			{0.20, 2, 36},
			{0.10, 1, 10},
			{0.10, 2, 19},
			{0.60, 2, 84},
			{0.70, 2, 91},
			{0.80, 2, 96},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("In independent Bernoulli trials each having success probability \\(p = %.2f\\), what is the probability (expressed as an integer percentage from 0 to 100) that the first success occurs on or before trial \\(k = %d\\)?", it.p, it.k),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("The cumulative probability is \\(P(X \\le %d) = 1 - (1-p)^{%d} = 1 - (%.2f)^{%d} = 1 - %.2f = %.2f\\), which equals %d percent.",
				it.k, it.k, 1-it.p, it.k, float64(100-it.pct)/100.0, float64(it.pct)/100.0, it.pct),
		}
	case 3:
		type ratioItem struct {
			p float64
			r int
		}
		items := []ratioItem{
			{0.50, 2},
			{0.75, 4},
			{0.80, 5},
			{0.90, 10},
			{0.95, 20},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Geometric}(p)\\) with success probability \\(p = %.2f\\). Compute the probability ratio \\(\\frac{P(X = 1)}{P(X = 2)}\\).", it.p),
			Answer:   fmt.Sprintf("%d", it.r),
			Explanation: fmt.Sprintf("For a geometric distribution, \\(P(X = 1) = p\\) and \\(P(X = 2) = (1-p)p\\). The ratio is \\(\\frac{P(X = 1)}{P(X = 2)} = \\frac{1}{1-p} = \\frac{1}{1 - %.2f} = \\frac{1}{%.2f} = %d\\).",
				it.p, 1-it.p, it.r),
		}
	default:
		type memItem struct {
			p float64
			s int
			e int
		}
		items := []memItem{
			{0.20, 5, 5},
			{0.25, 4, 4},
			{0.10, 8, 10},
			{0.50, 3, 2},
			{0.05, 12, 20},
			{0.04, 15, 25},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A geometric random variable \\(X\\) has success probability \\(p = %.2f\\). Given that the first %d trials resulted in failure, what is the expected number of additional trials needed until the first success?", it.p, it.s),
			Answer:   fmt.Sprintf("%d", it.e),
			Explanation: fmt.Sprintf("By the memoryless property of the geometric distribution, prior failures do not affect future trials: \\(E[X - %d \\mid X > %d] = E[X] = \\frac{1}{p} = \\frac{1}{%.2f} = %d\\).",
				it.s, it.s, it.p, it.e),
		}
	}
}

type hypergeometricGen struct{}

func (g *hypergeometricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type expItem struct {
			nTotal int
			k      int
			n      int
		}
		items := []expItem{
			{20, 8, 5},
			{30, 12, 10},
			{40, 16, 15},
			{50, 20, 25},
			{24, 6, 12},
			{36, 9, 16},
			{50, 15, 20},
			{60, 24, 15},
		}
		it := items[rand.Intn(len(items))]
		mean := it.n * it.k / it.nTotal
		return generator.Problem{
			Question: fmt.Sprintf("An urn contains %d total items including %d marked items. If a random sample of %d items is drawn without replacement, find the expected number of marked items \\(E[X]\\).", it.nTotal, it.k, it.n),
			Answer:   fmt.Sprintf("%d", mean),
			Explanation: fmt.Sprintf("The expected number of marked items in hypergeometric sampling without replacement is \\(E[X] = n \\frac{K}{N} = %d \\times \\frac{%d}{%d} = %d\\).",
				it.n, it.k, it.nTotal, mean),
		}
	case 1:
		type failItem struct {
			nTotal int
			k      int
			n      int
		}
		items := []failItem{
			{20, 8, 5},
			{30, 12, 10},
			{40, 16, 15},
			{50, 20, 25},
			{60, 24, 15},
		}
		it := items[rand.Intn(len(items))]
		nonDef := it.nTotal - it.k
		expNonDef := it.n * nonDef / it.nTotal
		return generator.Problem{
			Question: fmt.Sprintf("A batch of %d items contains %d defective items and %d non-defective items. If %d items are sampled without replacement, what is the expected number of non-defective items in the sample?", it.nTotal, it.k, nonDef, it.n),
			Answer:   fmt.Sprintf("%d", expNonDef),
			Explanation: fmt.Sprintf("The expected number of non-defective items is \\(E[X_{\\text{non-def}}] = n \\frac{N - K}{N} = %d \\times \\frac{%d}{%d} = %d\\).",
				it.n, nonDef, it.nTotal, expNonDef),
		}
	case 2:
		type maxItem struct {
			nTotal int
			k      int
			n      int
			max    int
		}
		items := []maxItem{
			{25, 4, 8, 4},
			{30, 5, 9, 5},
			{20, 9, 6, 6},
			{40, 6, 12, 6},
			{35, 14, 8, 8},
			{50, 7, 15, 7},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("From a population of \\(N = %d\\) containing \\(K = %d\\) target items, a random sample of \\(n = %d\\) is selected without replacement. What is the theoretical maximum number of target items that can appear in the sample?", it.nTotal, it.k, it.n),
			Answer:   fmt.Sprintf("%d", it.max),
			Explanation: fmt.Sprintf("The number of observed target items cannot exceed the total target items in the population \\(K = %d\\) or the sample size \\(n = %d\\): \\(\\max(X) = \\min(n, K) = \\min(%d, %d) = %d\\).",
				it.k, it.n, it.n, it.k, it.max),
		}
	case 3:
		type minItem struct {
			nTotal int
			k      int
			n      int
			min    int
		}
		items := []minItem{
			{20, 15, 8, 3},
			{25, 20, 9, 4},
			{30, 24, 10, 4},
			{35, 28, 12, 5},
			{40, 32, 14, 6},
		}
		it := items[rand.Intn(len(items))]
		nonTarget := it.nTotal - it.k
		return generator.Problem{
			Question: fmt.Sprintf("A committee of %d members contains %d seniors and %d juniors. If a subcommittee of %d members is chosen at random without replacement, what is the minimum number of seniors that must be on the subcommittee?", it.nTotal, it.k, nonTarget, it.n),
			Answer:   fmt.Sprintf("%d", it.min),
			Explanation: fmt.Sprintf("Even if all %d juniors were selected, the remaining seats must be filled by seniors: \\(\\min(X) = \\max(0, n - (N - K)) = %d - (%d - %d) = %d - %d = %d\\).",
				nonTarget, it.n, it.nTotal, it.k, it.n, nonTarget, it.min),
		}
	default:
		type combItem struct {
			nTotal int
			n      int
			comb   int
		}
		items := []combItem{
			{6, 2, 15},
			{7, 2, 21},
			{7, 3, 35},
			{8, 2, 28},
			{8, 3, 56},
			{9, 2, 36},
			{10, 2, 45},
			{10, 3, 120},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("In a hypergeometric experiment with total population size \\(N = %d\\), how many distinct unordered samples of size \\(n = %d\\) can be drawn without replacement?", it.nTotal, it.n),
			Answer:   fmt.Sprintf("%d", it.comb),
			Explanation: fmt.Sprintf("The total number of possible samples is given by the combination formula \\(\\binom{N}{n} = \\binom{%d}{%d} = \\frac{%d!}{%d!(%d - %d)!} = %d\\).",
				it.nTotal, it.n, it.nTotal, it.n, it.nTotal, it.n, it.comb),
		}
	}
}

type normalDistGen struct{}

func (g *normalDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"How many parameters define the Normal distribution?", "2", "N(mean, variance) is parameterized by its mean and its variance."},
		{"About what proportion of Normal data falls within 2 SD of the mean? Answer as a decimal.", "0.95", "About 0.68 within 1 SD, 0.95 within 2 SD, 0.997 within 3 SD."},
		{"What is the mean of a standard Normal distribution?", "0", "Standard Normal: N(0,1) with μ=0, σ=1."},
		{"What is the variance of a standard Normal distribution?", "1", "Standard Normal: N(0,1) with μ=0, σ²=1."},
		{"If X = 70 comes from a Normal population with mean 50 and SD 10, what is the z-score?", "2", "Standardization: z = (x - mean)/SD = 30/10."},
		{"The Normal distribution is symmetric about its ____.", "mean", "The Normal distribution is symmetric about μ."},
		{"What is the shape of the Normal distribution?", "bell curve", "The Normal distribution has the characteristic bell curve shape."},
		{"Do sample means approach a Normal distribution as the sample size grows? (yes/no)", "yes", "CLT: sample means of independent identical variables approach Normal as n grows."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type poissonGen struct{}

func (g *poissonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type sdItem struct {
			lambda int
			sd     int
		}
		items := []sdItem{
			{4, 2},
			{9, 3},
			{16, 4},
			{25, 5},
			{36, 6},
			{49, 7},
			{64, 8},
			{81, 9},
			{100, 10},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) follows a Poisson distribution \\(\\operatorname{Poisson}(%d)\\). Find the standard deviation \\(\\sigma\\).", it.lambda),
			Answer:   fmt.Sprintf("%d", it.sd),
			Explanation: fmt.Sprintf("For a Poisson distribution, the variance equals the parameter: \\(\\operatorname{Var}(X) = \\lambda = %d\\). Therefore, the standard deviation is \\(\\sigma = \\sqrt{\\lambda} = \\sqrt{%d} = %d\\).",
				it.lambda, it.lambda, it.sd),
		}
	case 1:
		rates := []int{2, 3, 4, 5, 6, 7, 8}
		times := []int{2, 3, 4, 5}
		lambda := rates[rand.Intn(len(rates))]
		t := times[rand.Intn(len(times))]
		total := lambda * t
		return generator.Problem{
			Question: fmt.Sprintf("Events occur according to a Poisson process at an average rate of %d events per hour. How many events are expected to occur over a %d-hour interval?", lambda, t),
			Answer:   fmt.Sprintf("%d", total),
			Explanation: fmt.Sprintf("For a Poisson process with rate \\(\\lambda = %d\\), the expected number of events in \\(t = %d\\) hours is \\(E[N(t)] = \\lambda t = %d \\times %d = %d\\).",
				lambda, t, lambda, t, total),
		}
	case 2:
		l1List := []int{3, 5, 7, 10, 12, 15}
		l2List := []int{4, 6, 8, 11, 14, 20}
		l1 := l1List[rand.Intn(len(l1List))]
		l2 := l2List[rand.Intn(len(l2List))]
		sum := l1 + l2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X_1 \\sim \\operatorname{Poisson}(%d)\\) and \\(X_2 \\sim \\operatorname{Poisson}(%d)\\) be independent random variables. Find the variance of their sum \\(\\operatorname{Var}(X_1 + X_2)\\).", l1, l2),
			Answer:   fmt.Sprintf("%d", sum),
			Explanation: fmt.Sprintf("Since \\(X_1\\) and \\(X_2\\) are independent, \\(\\operatorname{Var}(X_1 + X_2) = \\operatorname{Var}(X_1) + \\operatorname{Var}(X_2) = \\lambda_1 + \\lambda_2 = %d + %d = %d\\).",
				l1, l2, sum),
		}
	case 3:
		type ratioItem struct {
			k      int
			lambda int
			ratio  int
		}
		items := []ratioItem{
			{2, 6, 3},
			{2, 8, 4},
			{2, 10, 5},
			{3, 12, 4},
			{3, 15, 5},
			{4, 20, 5},
			{4, 24, 6},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Poisson}(%d)\\). Find the probability ratio \\(\\frac{P(X = %d)}{P(X = %d)}\\).", it.lambda, it.k, it.k-1),
			Answer:   fmt.Sprintf("%d", it.ratio),
			Explanation: fmt.Sprintf("From the Poisson PMF, consecutive terms satisfy \\(\\frac{P(X = k)}{P(X = k-1)} = \\frac{\\lambda}{k}\\). Here \\(\\frac{P(X = %d)}{P(X = %d)} = \\frac{%d}{%d} = %d\\).",
				it.k, it.k-1, it.lambda, it.k, it.ratio),
		}
	default:
		lambdas := []int{2, 3, 4, 5, 6, 7, 8, 10}
		lam := lambdas[rand.Intn(len(lambdas))]
		m2 := lam + lam*lam
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Poisson}(%d)\\). Compute the second moment \\(E[X^2]\\).", lam),
			Answer:   fmt.Sprintf("%d", m2),
			Explanation: fmt.Sprintf("Using \\(\\operatorname{Var}(X) = E[X^2] - (E[X])^2\\), we find \\(E[X^2] = \\lambda + \\lambda^2 = %d + %d = %d\\).",
				lam, lam*lam, m2),
		}
	}
}

type studentTGen struct{}

// The answers here are "yes", "no", "standard Normal" and plain numbers, and the
// concept is declared symbolic, whose hint says "Answer in simplest form" — so a
// learner was asked a yes/no question and told to simplify. Text is the honest
// hint, and the grader has to accept words because SymPy cannot.
func (g *studentTGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Answer with the word or number the question asks for", InputMode: "text"}
}

func (g *studentTGen) Grade(expected, userAnswer string) grader.Result {
	e := strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(expected))), " ")
	a := strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(userAnswer))), " ")
	if a == "" {
		return grader.Result{Correct: false, Score: 0, Feedback: "Type the answer"}
	}
	if e == a {
		return grader.Result{Correct: true, Score: 1}
	}
	if ef, errE := strconv.ParseFloat(e, 64); errE == nil {
		if af, errA := strconv.ParseFloat(a, 64); errA == nil && af == ef {
			return grader.Result{Correct: true, Score: 1}
		}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect"}
}

func (g *studentTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Is Student's t defined by its degrees of freedom? (yes/no)", "yes", "t(ν) is parameterized by degrees of freedom ν."},
		{"As ν → ∞, the t distribution approaches ____.", "standard Normal", "As degrees of freedom increase, t approaches N(0,1)."},
		{"Compared to Normal, does Student's t have heavier tails? (yes/no)", "yes", "The t distribution has heavier tails than Normal, especially for small ν."},
		{"What is the mean of t(ν) for ν > 1?", "0", "The t distribution is symmetric about 0 with mean 0 (for ν > 1)."},
		{"Is the t distribution used instead of Normal when sigma is unknown? (yes/no)", "yes", "t is used when the population standard deviation is unknown."},
		{"The t(nu) variance is finite when nu exceeds which integer? (enter a number)", "2", "Var[t(ν)] = ν/(ν-2) for ν > 2."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type uniformDistGen struct{}

func (g *uniformDistGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		type meanItem struct {
			a int
			b int
		}
		items := []meanItem{
			{4, 16},
			{10, 34},
			{6, 26},
			{12, 48},
			{5, 35},
			{8, 40},
			{14, 50},
		}
		it := items[rand.Intn(len(items))]
		mean := (it.a + it.b) / 2
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Uniform}(%d, %d)\\) be a continuous uniform random variable. Find the expected value \\(E[X]\\).", it.a, it.b),
			Answer:   fmt.Sprintf("%d", mean),
			Explanation: fmt.Sprintf("The expected value of a uniform random variable on interval \\([a, b]\\) is \\(E[X] = \\frac{a + b}{2} = \\frac{%d + %d}{2} = %d\\).",
				it.a, it.b, mean),
		}
	case 1:
		type varItem struct {
			a int
			b int
		}
		items := []varItem{
			{0, 12},
			{2, 14},
			{1, 19},
			{0, 24},
			{4, 16},
			{5, 11},
			{10, 40},
			{7, 25},
		}
		it := items[rand.Intn(len(items))]
		diff := it.b - it.a
		v := diff * diff / 12
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Uniform}(%d, %d)\\). Calculate the variance \\(\\operatorname{Var}(X)\\).", it.a, it.b),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf("The variance of a uniform distribution on \\([a, b]\\) is \\(\\operatorname{Var}(X) = \\frac{(b - a)^2}{12} = \\frac{(%d - %d)^2}{12} = \\frac{%d}{12} = %d\\).",
				it.b, it.a, diff*diff, v),
		}
	case 2:
		type cdfItem struct {
			a   int
			b   int
			x   int
			pct int
		}
		items := []cdfItem{
			{0, 20, 8, 40},
			{10, 50, 20, 25},
			{5, 25, 15, 50},
			{0, 50, 35, 70},
			{20, 60, 50, 75},
			{10, 35, 30, 80},
			{0, 40, 24, 60},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("A random variable \\(X\\) is uniformly distributed on the interval \\([%d, %d]\\). What is the probability \\(P(X \\le %d)\\) expressed as an integer percentage from 0 to 100?", it.a, it.b, it.x),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("For a uniform distribution on \\([a, b]\\), the cumulative probability is \\(P(X \\le %d) = \\frac{%d - %d}{%d - %d} = \\frac{%d}{%d} = %.2f\\), which equals %d percent.",
				it.x, it.x, it.a, it.b, it.a, it.x-it.a, it.b-it.a, float64(it.pct)/100.0, it.pct),
		}
	case 3:
		type subItem struct {
			a   int
			b   int
			c   int
			d   int
			pct int
		}
		items := []subItem{
			{0, 100, 25, 65, 40},
			{10, 60, 20, 45, 50},
			{0, 50, 15, 30, 30},
			{20, 70, 30, 60, 60},
			{0, 80, 16, 56, 50},
			{0, 100, 15, 85, 70},
		}
		it := items[rand.Intn(len(items))]
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Uniform}(%d, %d)\\). Find the probability \\(P(%d \\le X \\le %d)\\) expressed as an integer percentage from 0 to 100.", it.a, it.b, it.c, it.d),
			Answer:   fmt.Sprintf("%d", it.pct),
			Explanation: fmt.Sprintf("The probability that \\(X\\) falls in sub-interval \\([%d, %d]\\) is \\(P(%d \\le X \\le %d) = \\frac{%d - %d}{%d - %d} = \\frac{%d}{%d} = %.2f\\), which equals %d percent.",
				it.c, it.d, it.c, it.d, it.d, it.c, it.b, it.a, it.d-it.c, it.b-it.a, float64(it.pct)/100.0, it.pct),
		}
	default:
		type m2Item struct {
			a  int
			b  int
			m2 int
		}
		items := []m2Item{
			{0, 6, 12},
			{0, 9, 27},
			{0, 12, 48},
			{0, 15, 75},
			{3, 6, 21},
			{2, 8, 28},
			{3, 9, 39},
			{4, 10, 52},
		}
		it := items[rand.Intn(len(items))]
		sumSq := it.a*it.a + it.a*it.b + it.b*it.b
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(X \\sim \\operatorname{Uniform}(%d, %d)\\). Compute the second moment \\(E[X^2]\\).", it.a, it.b),
			Answer:   fmt.Sprintf("%d", it.m2),
			Explanation: fmt.Sprintf("The second moment of a continuous uniform distribution on \\([a, b]\\) is \\(E[X^2] = \\frac{a^2 + ab + b^2}{3} = \\frac{%d}{3} = %d\\).",
				sumSq, it.m2),
		}
	}
}

type zTableGen struct{}

func (g *zTableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What z-score corresponds to the 50th percentile of N(0,1)?", "0", "The median of standard Normal is 0."},
		{"What z-score approximately corresponds to the 97.5th percentile?", "1.96", "P(Z < 1.96) ≈ 0.975, used for 95% confidence intervals."},
		{"What is the z-score for a 95% confidence interval?", "1.96", "95% CI uses z = 1.96 (approximately 2)."},
		{"Does a z-score measure how many SD a value is from the mean? (yes/no)", "yes", "z = (x - mean)/SD measures standard deviations from the mean."},
		{"About what proportion of data falls within z = +-1? Answer as a decimal.", "0.68", "About 0.68 of data falls within 1 standard deviation of the mean."},
		{"About what proportion of data falls within z = +-2? Answer as a decimal.", "0.95", "About 0.95 of data falls within 2 standard deviations of the mean."},
		{"About what proportion of data falls within z = +-3? Answer as a decimal.", "0.997", "About 0.997 of data falls within 3 standard deviations of the mean."},
		{"A z-score of -2 means the value is ____.", "2 SD below the mean", "Negative z: below mean; positive z: above mean."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// ----- stat.infer.* (inference) -----

type confidenceGen struct{}

func (g *confidenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Does a 95% confidence interval mean that 95% of such intervals contain the true parameter? (yes/no)", "yes", "In repeated sampling, 95% of CIs contain the population parameter."},
		{"A sample of n=100 has SD 20. With z=2, what is the margin of error?", "4", "CI = sample mean plus/minus margin; margin = z times SD/sqrt(n) = 2 times 2."},
		{"Does a larger sample give a narrower confidence interval? (yes/no)", "yes", "CI width goes like 1/sqrt(n), so larger samples give more precise estimates."},
		{"What happens to CI width when confidence level increases from 95% to 99%?", "it gets wider", "Higher confidence requires a wider interval."},
		{"A 95% CI is 10 plus/minus 2. What is the margin of error?", "2", "Margin of error = critical value times standard error."},
		{"For 50 of 100 voters with z=2, what is the margin of error as a decimal?", "0.1", "Margin = z times sqrt(phat(1-phat)/n) = 2 times 0.05."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type samplingGen struct{}

func (g *samplingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"The population SD is 12 and n=36. What is the standard error of the mean?", "2", "Standard error = population SD / sqrt(sample size) = 12/6."},
		{"Does the sampling distribution of the mean approach Normal as n grows? (yes/no)", "yes", "CLT: the sampling distribution of the mean approaches Normal for large n."},
		{"The population mean is 50. What is the mean of the sampling distribution of xbar?", "50", "xbar is an unbiased estimator of the population mean."},
		{"Is a sampling distribution the distribution of a statistic over all possible samples? (yes/no)", "yes", "The sampling distribution shows how a statistic varies from sample to sample."},
		{"Does a larger sample size reduce bias? (yes/no)", "no", "Bias depends on the sampling method, not sample size; bias is independent of n."},
		{"Is a sample a subset of the population? (yes/no)", "yes", "The population is the entire group; the sample is a subset drawn from it."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

// ----- stat.prob.* (additional probability) -----

type bayesGen struct{}

func (g *bayesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Is P(A|B) = P(B|A)P(A)/P(B) Bayes theorem? (yes/no)", "yes", "Bayes theorem relates conditional probabilities."},
		{"In Bayes' theorem, P(A) is called the ____ probability.", "prior", "P(A) is the prior probability before observing B."},
		{"In Bayes' theorem, P(A|B) is called the ____ probability.", "posterior", "P(A|B) is the updated probability after observing B."},
		{"What is the denominator in Bayes' theorem called?", "marginal likelihood", "P(B) normalizes the posterior; it is called the marginal likelihood or evidence."},
		{"In medical testing, P(test positive | has disease) is called ____.", "sensitivity", "Sensitivity = P(positive|disease), the true positive rate."},
		{"In medical testing, P(test negative | no disease) is called ____.", "specificity", "Specificity = P(negative|no disease), the true negative rate."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type continuousRVGen struct{}

func (g *continuousRVGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the probability of a continuous random variable taking any single value?", "0", "P(X = a) = 0 for continuous random variables."},
		{"What function gives probabilities for a continuous random variable via areas?", "probability density function (PDF)", "P(a ≤ X ≤ b) = ∫ₐᵇ f(x)dx."},
		{"Is the CDF equal to P(X <= x), the probability of being less than or equal to x? (yes/no)", "yes", "The CDF gives the probability of being less than or equal to x."},
		{"Is the CDF the integral of the PDF? (yes/no)", "yes", "CDF(x) is the integral of the PDF up to x, and PDF is the derivative of the CDF."},
		{"For a valid PDF, what must ∫_{-∞}^{∞} f(x)dx equal?", "1", "Total probability under the PDF must be 1."},
		{"Is the PDF of a continuous variable always ≤ 1?", "no", "The PDF can be > 1 as long as the total area is 1."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type discreteRVGen struct{}

func (g *discreteRVGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Does the PMF give P(X = x) for each possible x? (yes/no)", "yes", "PMF gives the probability of each specific value for a discrete variable."},
		{"For a valid PMF, what must Σ P(X=x) equal?", "1", "All probabilities must sum to 1."},
		{"Is the PMF for discrete variables and the PDF for continuous ones? (yes/no)", "yes", "PMF gives probabilities at specific points; PDF gives density."},
		{"Is the discrete CDF the sum of PMF values up to x? (yes/no)", "yes", "Discrete CDF sums PMF values up to x."},
		{"Can a discrete random variable take infinitely many values?", "yes", "Discrete variables can be countably infinite (e.g., Poisson, Geometric)."},
		{"X takes values 1 and 3 with probability 0.5 each. What is E[X]?", "2", "The mean is the probability-weighted sum: 1 times 0.5 plus 3 times 0.5."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type expectedValueGen struct{}

func (g *expectedValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the expected value of a constant c?", "c", "E[c] = c for any constant c."},
		{"If E[X] = 5, what is E[2X + 1]?", "11", "Linearity of expectation: E[aX+b] = aE[X] + b = 2 times 5 plus 1."},
		{"If E[X] = 3 and E[Y] = 4, what is E[X + Y]?", "7", "Linearity: expectation of sum = sum of expectations."},
		{"Is E[XY] always equal to E[X]E[Y]?", "no", "E[XY] = E[X]E[Y] only when X and Y are independent."},
		{"If E[X] = 2 and E[X^2] = 8, what is Var[X]?", "4", "Variance = second moment minus squared mean = 8 - 4."},
		{"What is E[X²] called?", "second moment", "E[X²] is the second moment about the origin."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type gaussianFuncGen struct{}

func (g *gaussianFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	tableEasy := []entry{
		{"Standard normal is Gaussian with mean? (enter a number)", "0", "Centered at zero by definition."},
		{"Gaussian e^(-x^2): maximum at x = ? (enter a number)", "0", "Peak at the mean, x = 0."},
		{"Is a Gaussian symmetric about its mean? (yes/no)", "yes", "Bell curve mirrors left-right."},
		{"How many parameters define a Gaussian (mean and variance)? (enter a number)", "2", "Mean plus variance."},
	}
	tableHard := []entry{
		{"Gaussian variance 4: standard deviation? (enter a number)", "2", "SD is the square root of variance."},
		{"Wider bell: larger or smaller variance? (type larger or smaller)", "larger", "Variance spreads the curve."},
		{"Is the sum of independent Gaussians Gaussian? (yes/no)", "yes", "Closed under convolution."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}
