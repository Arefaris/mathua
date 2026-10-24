package statistics

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The batch conversions are checked for scaffolding, difficulty scaling,
// variety and give-aways, but nothing recomputed the answer from the numbers
// the question prints. That is the check that catches a formula error, and it
// is the check the statistics batch was merged without: a fixed table can be
// read by eye, a parameterised one cannot.
//
// Each case parses the printed stem, computes the answer independently and
// compares. Quartiles are compared as exact halves rather than rounded, so a
// data set whose half-sum is odd fails here instead of being silently
// truncated by integer division.

var (
	statProbeInts     = regexp.MustCompile(`-?\d+`)
	statProbeMeanSd   = regexp.MustCompile(`mean (-?\d+) and standard deviation (-?\d+)`)
	statProbeValue    = regexp.MustCompile(`value (-?\d+)`)
	statProbeSizeSd   = regexp.MustCompile(`size (-?\d+) has a standard deviation of (-?\d+)`)
	statProbeCritSe   = regexp.MustCompile(`critical value of (-?\d+) and a standard error of (-?\d+)`)
	statProbeFraction = regexp.MustCompile(`frac\{(-?\d+)\}\{(-?\d+)\}`)
)

// recomputeSetOf reads the data set a question prints between its braces.
func recomputeSetOf(t *testing.T, question string) []int {
	t.Helper()
	open := strings.Index(question, "{")
	shut := strings.LastIndex(question, "}")
	if open < 0 || shut <= open {
		t.Fatalf("question prints no data set: %q", question)
	}
	body := strings.ReplaceAll(question[open+1:shut], `\`, "")
	var out []int
	for _, token := range statProbeInts.FindAllString(body, -1) {
		n, err := strconv.Atoi(token)
		if err != nil {
			t.Fatalf("data set member %q is not an integer: %q", token, question)
		}
		out = append(out, n)
	}
	if len(out) == 0 {
		t.Fatalf("data set is empty: %q", question)
	}
	return out
}

// recomputeMedianOf returns the median of an even-length slice as an exact
// half, so a fractional median is visible rather than lost to truncation.
func recomputeMedianOf(sorted []int) float64 {
	if len(sorted)%2 == 1 {
		return float64(sorted[len(sorted)/2])
	}
	return float64(sorted[len(sorted)/2-1]+sorted[len(sorted)/2]) / 2
}

func recomputeGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func TestStatBatchAnswersRecomputeFromTheStem(t *testing.T) {
	reg := testStatRegistry(t)
	checked := 0
	for _, id := range convertedStatNumericIDs {
		for i := 0; i < 100; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: rand.Float64(),
				Seed:       int64(70000 + i),
			})
			if err != nil {
				t.Fatalf("%s: generate failed: %v", id, err)
			}
			want := ""
			switch id {
			case "stat.spread.mad":
				data := recomputeSetOf(t, p.Question)
				total := 0
				for _, x := range data {
					total += x
				}
				if total%len(data) != 0 {
					t.Fatalf("%s: the printed mean is not whole: %v", id, data)
				}
				mean := total / len(data)
				deviation := 0
				for _, x := range data {
					if x > mean {
						deviation += x - mean
					} else {
						deviation += mean - x
					}
				}
				if deviation%len(data) != 0 {
					t.Fatalf("%s: the deviation does not divide by n: %v", id, data)
				}
				want = strconv.Itoa(deviation / len(data))

			case "stat.spread.iqr", "stat.quartile.value":
				data := recomputeSetOf(t, p.Question)
				half := len(data) / 2
				lower := recomputeMedianOf(data[:half])
				upper := recomputeMedianOf(data[half:])
				if lower != float64(int(lower)) || upper != float64(int(upper)) {
					t.Fatalf("%s: a quartile is not a whole number: %v gives Q1=%v Q3=%v",
						id, data, lower, upper)
				}
				switch {
				case id == "stat.spread.iqr":
					if upper-lower <= 0 {
						t.Fatalf("%s: the printed set has no spread: %v", id, data)
					}
					want = strconv.Itoa(int(upper - lower))
				case strings.Contains(p.Question, "lower"):
					want = strconv.Itoa(int(lower))
				default:
					want = strconv.Itoa(int(upper))
				}

			case "stat.zscore.value":
				dist := statProbeMeanSd.FindStringSubmatch(p.Question)
				value := statProbeValue.FindStringSubmatch(p.Question)
				if dist == nil || value == nil {
					t.Fatalf("%s: question does not print mean, deviation and value: %q", id, p.Question)
				}
				mean, _ := strconv.Atoi(dist[1])
				sd, _ := strconv.Atoi(dist[2])
				x, _ := strconv.Atoi(value[1])
				if (x-mean)%sd != 0 {
					t.Fatalf("%s: the z-score is not whole: mean %d, sd %d, x %d", id, mean, sd, x)
				}
				want = strconv.Itoa((x - mean) / sd)

			case "stat.sampling.se":
				fields := statProbeSizeSd.FindStringSubmatch(p.Question)
				if fields == nil {
					t.Fatalf("%s: question does not print size and deviation: %q", id, p.Question)
				}
				n, _ := strconv.Atoi(fields[1])
				s, _ := strconv.Atoi(fields[2])
				root := 0
				for root*root < n {
					root++
				}
				if root*root != n {
					t.Fatalf("%s: the sample size is not a perfect square: %d", id, n)
				}
				if s%root != 0 {
					t.Fatalf("%s: the standard error is not whole: %d / sqrt(%d)", id, s, n)
				}
				want = strconv.Itoa(s / root)

			case "stat.infer.margin":
				fields := statProbeCritSe.FindStringSubmatch(p.Question)
				if fields == nil {
					t.Fatalf("%s: question does not print the critical value and standard error: %q", id, p.Question)
				}
				z, _ := strconv.Atoi(fields[1])
				se, _ := strconv.Atoi(fields[2])
				want = strconv.Itoa(z * se)

			case "stat.prob.conditional":
				pairs := statProbeFraction.FindAllStringSubmatch(p.Question, -1)
				if len(pairs) < 2 {
					t.Fatalf("%s: question does not print both fractions: %q", id, p.Question)
				}
				joint, jointDen := atoiPair(pairs[0])
				marginal, marginalDen := atoiPair(pairs[1])
				if jointDen != marginalDen || marginal == 0 || joint >= marginal {
					t.Fatalf("%s: the printed probabilities are not comparable: %q", id, p.Question)
				}
				divisor := recomputeGCD(joint, marginal)
				want = strconv.Itoa(joint/divisor) + "/" + strconv.Itoa(marginal/divisor)

			default:
				t.Fatalf("no recomputation is written for %s", id)
			}
			if p.Answer != want {
				t.Fatalf("%s: answer %q does not recompute to %q from the stem: %q",
					id, p.Answer, want, p.Question)
			}
			checked++
		}
	}
	if checked < 700 {
		t.Fatalf("only %d problems were recomputed, so the sweep is not covering the batch", checked)
	}
}

func atoiPair(pair []string) (int, int) {
	a, _ := strconv.Atoi(pair[1])
	b, _ := strconv.Atoi(pair[2])
	return a, b
}
