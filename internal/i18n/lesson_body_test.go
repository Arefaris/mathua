package i18n_test

// These tests hold every translated lesson body, data/lessons/<corpus>/<locale>/<path>.md, to
// its English original. A lesson is formulas, figures and links with prose around them, and a
// translation may change only the prose: one altered digit in a formula, one dropped figure or
// one lost licence line makes it a different lesson.
//
// The comparison itself (compareLessonBodies) works on two strings and has its own fixture test,
// which feeds it deliberately broken translations. A check that was never seen to fail is not
// evidence of anything.

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode"

	"github.com/chuma-beep/mathua/internal/i18n"
	"github.com/chuma-beep/mathua/internal/lessons"
)

const (
	lessonsDir = dataDir + "/lessons"
	// codeMark and mathMark stand in for a code span and a math span once they are taken out
	// of the text, so that what is left is prose and structure.
	codeMark = "\x01"
	mathMark = "\x02"
)

// lessonCorpora are the directories under data/lessons whose lessons can be translated.
var lessonCorpora = []string{"teaching", "algebrica", "authored"}

var (
	// An attribution or licence line is a licence term and is copied byte for byte: the
	// "Content sourced from" and "Content adapted from" lines of the vendored books, the
	// hand-authored notice, and any line that names a Creative Commons licence.
	attributionRe = regexp.MustCompile(`^>\s*(?:Content (?:sourced|adapted) from|Hand-authored lesson)|CC BY|Creative Commons`)
	headingRe     = regexp.MustCompile(`^ {0,3}(#{1,6})(?:[ \t]|$)`)
	listItemRe    = regexp.MustCompile(`^\s*(?:[-*+]|\d+[.)])\s+`)
	tableRowRe    = regexp.MustCompile(`^\s*\|`)
	// The target of a markdown link or image: `](target)`, with one level of parentheses
	// allowed inside it, as in a Wikipedia address.
	linkTargetRe = regexp.MustCompile(`\]\(\s*(<[^>]*>|(?:[^()\s]|\([^()]*\))*)`)
	htmlRefRe    = regexp.MustCompile(`(?i)\b(?:src|href)\s*=\s*["']([^"']*)["']`)
	urlRe        = regexp.MustCompile(`https?://[^\s)>\]"'<]+`)
	htmlTagRe    = regexp.MustCompile(`</?[A-Za-z][^>]*>`)
	blankLineRe  = regexp.MustCompile(`\n[ \t]*\n`)
	latinWordRe  = regexp.MustCompile(`[A-Za-z]{3,}`)
	digitRunRe   = regexp.MustCompile(`\d+`)
	letterRunRe  = regexp.MustCompile(`\p{L}+`)
)

// lessonFacts is everything about a lesson body that a translation must not change, plus the
// prose that it is meant to change.
type lessonFacts struct {
	math        []string // $…$, $$…$$, \(…\), \[…\], delimiters included
	code        []string // fenced blocks and inline code
	targets     []string // link and image targets, HTML src/href, and URLs
	images      int
	headings    [7]int // by level, 1 to 6
	paragraphs  int
	listItems   int
	tableRows   int
	attribution []string // attribution and licence lines
	strayDollar int      // a "$" that opens no math span
	prose       string   // what is left of the text
}

func parseLesson(md string) lessonFacts {
	md = strings.ReplaceAll(md, "\r\n", "\n")
	var f lessonFacts

	var kept []string
	for _, ln := range strings.Split(md, "\n") {
		if attributionRe.MatchString(ln) {
			f.attribution = append(f.attribution, ln)
			continue
		}
		kept = append(kept, ln)
	}
	rest := strings.Join(kept, "\n")

	rest = takeFencedCode(rest, &f.code)
	rest = takeInlineCode(rest, &f.code)
	rest, f.strayDollar = takeMath(rest, &f.math)

	for _, m := range linkTargetRe.FindAllStringSubmatch(rest, -1) {
		f.targets = append(f.targets, m[1])
	}
	for _, m := range htmlRefRe.FindAllStringSubmatch(rest, -1) {
		f.targets = append(f.targets, m[1])
	}
	f.targets = append(f.targets, urlRe.FindAllString(rest, -1)...)
	f.images = strings.Count(rest, "![")

	for _, ln := range strings.Split(rest, "\n") {
		switch {
		case headingRe.MatchString(ln):
			f.headings[len(headingRe.FindStringSubmatch(ln)[1])]++
		case listItemRe.MatchString(ln):
			f.listItems++
		}
		if tableRowRe.MatchString(ln) {
			f.tableRows++
		}
	}
	for _, p := range blankLineRe.Split(rest, -1) {
		if strings.TrimSpace(p) != "" {
			f.paragraphs++
		}
	}

	prose := linkTargetRe.ReplaceAllString(rest, "](")
	prose = urlRe.ReplaceAllString(prose, " ")
	prose = htmlTagRe.ReplaceAllString(prose, " ")
	prose = strings.NewReplacer(codeMark, " ", mathMark, " ").Replace(prose)
	f.prose = prose
	return f
}

// takeFencedCode replaces each ``` or ~~~ block with one codeMark line and records the block.
func takeFencedCode(s string, code *[]string) string {
	var out, block []string
	var fence string
	for _, ln := range strings.Split(s, "\n") {
		t := strings.TrimSpace(ln)
		switch {
		case fence == "" && (strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~")):
			fence = strings.Repeat(t[:1], len(t)-len(strings.TrimLeft(t, t[:1])))
			block = []string{ln}
		case fence != "":
			block = append(block, ln)
			if strings.HasPrefix(t, fence) && strings.Trim(t, t[:1]) == "" {
				*code = append(*code, strings.Join(block, "\n"))
				out = append(out, codeMark)
				fence = ""
			}
		default:
			out = append(out, ln)
		}
	}
	if fence != "" { // never closed: the rest of the text is the block
		*code = append(*code, strings.Join(block, "\n"))
		out = append(out, codeMark)
	}
	return strings.Join(out, "\n")
}

// takeInlineCode replaces each `…` span (a run of n backticks to the next run of n) with codeMark.
func takeInlineCode(s string, code *[]string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] != '`' {
			b.WriteByte(s[i])
			i++
			continue
		}
		n := 1
		for i+n < len(s) && s[i+n] == '`' {
			n++
		}
		run := strings.Repeat("`", n)
		end := -1
		for j := i + n; j < len(s); {
			k := strings.Index(s[j:], "`")
			if k < 0 {
				break
			}
			m := 1
			for j+k+m < len(s) && s[j+k+m] == '`' {
				m++
			}
			if m == n {
				end = j + k
				break
			}
			j += k + m
		}
		if end < 0 {
			b.WriteString(run)
			i += n
			continue
		}
		*code = append(*code, s[i:end+n])
		b.WriteString(codeMark)
		i = end + n
	}
	return b.String()
}

// takeMath replaces each math span with mathMark and records it with its delimiters. A "\$" is
// a dollar sign, not a delimiter. An inline $…$ ends before the next blank line, so a lone
// dollar sign cannot swallow a paragraph; it is counted as stray instead.
func takeMath(s string, math *[]string) (string, int) {
	var b strings.Builder
	stray := 0
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case c == '\\' && i+1 < len(s):
			switch s[i+1] {
			case '(', '[':
				closer := `\)`
				if s[i+1] == '[' {
					closer = `\]`
				}
				if k := strings.Index(s[i+2:], closer); k >= 0 {
					end := i + 2 + k + 2
					*math = append(*math, s[i:end])
					b.WriteString(mathMark)
					i = end
					continue
				}
			}
			b.WriteString(s[i : i+2])
			i += 2
		case c == '$' && i+1 < len(s) && s[i+1] == '$':
			if k := strings.Index(s[i+2:], "$$"); k >= 0 {
				end := i + 2 + k + 2
				*math = append(*math, s[i:end])
				b.WriteString(mathMark)
				i = end
				continue
			}
			stray++
			b.WriteString("$$")
			i += 2
		case c == '$':
			end := closingDollar(s, i+1)
			if end < 0 {
				stray++
				b.WriteByte('$')
				i++
				continue
			}
			*math = append(*math, s[i:end+1])
			b.WriteString(mathMark)
			i = end + 1
		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String(), stray
}

// closingDollar is the index of the "$" that closes an inline span opened just before from, or -1
// if a blank line or the end of the text comes first.
func closingDollar(s string, from int) int {
	for j := from; j < len(s); j++ {
		switch s[j] {
		case '\\':
			j++
		case '$':
			return j
		case '\n':
			k := j + 1
			for k < len(s) && (s[k] == ' ' || s[k] == '\t') {
				k++
			}
			if k >= len(s) || s[k] == '\n' {
				return -1
			}
		}
	}
	return -1
}

// diffMultiset reports what a has that b lacks, and what b has that a lacks.
func diffMultiset(a, b []string) (lost, added []string) {
	count := map[string]int{}
	for _, x := range a {
		count[x]++
	}
	for _, x := range b {
		count[x]--
	}
	for x, n := range count {
		for ; n > 0; n-- {
			lost = append(lost, x)
		}
		for ; n < 0; n++ {
			added = append(added, x)
		}
	}
	sort.Strings(lost)
	sort.Strings(added)
	return lost, added
}

func sample(xs []string) string {
	const max = 4
	var out []string
	for i, x := range xs {
		if i == max {
			out = append(out, fmt.Sprintf("… and %d more", len(xs)-max))
			break
		}
		if len(x) > 80 {
			x = x[:80] + "…"
		}
		out = append(out, fmt.Sprintf("%q", x))
	}
	return strings.Join(out, ", ")
}

// compareLessonBodies lists everything wrong with tr as a translation of en into loc. An empty
// result means the formulas, code, figures, links, structure and licence lines are the original's
// and the prose is in the right alphabet.
func compareLessonBodies(loc i18n.Locale, en, tr string) []string {
	e, t := parseLesson(en), parseLesson(tr)
	var problems []string
	add := func(format string, args ...interface{}) { problems = append(problems, fmt.Sprintf(format, args...)) }

	for _, kind := range []struct {
		name   string
		en, tr []string
	}{
		{"math span", e.math, t.math},
		{"code", e.code, t.code},
		{"link or figure target", e.targets, t.targets},
	} {
		lost, added := diffMultiset(kind.en, kind.tr)
		if len(lost) > 0 {
			add("%d %s(s) of the original are missing or changed: %s", len(lost), kind.name, sample(lost))
		}
		if len(added) > 0 {
			add("%d %s(s) are not in the original: %s", len(added), kind.name, sample(added))
		}
	}
	if e.images != t.images {
		add("%d image reference(s) for %d in the original", t.images, e.images)
	}
	for level := 1; level <= 6; level++ {
		if e.headings[level] != t.headings[level] {
			add("%d level-%d heading(s) for %d in the original", t.headings[level], level, e.headings[level])
		}
	}
	for _, c := range []struct {
		name   string
		en, tr int
	}{
		{"paragraph", e.paragraphs, t.paragraphs},
		{"list item", e.listItems, t.listItems},
		{"table row", e.tableRows, t.tableRows},
		{"stray dollar sign", e.strayDollar, t.strayDollar},
	} {
		if c.en != c.tr {
			add("%d %s(s) for %d in the original", c.tr, c.name, c.en)
		}
	}
	// A licence term is kept as the original wrote it, and as a whole line.
	lostLines, _ := diffMultiset(e.attribution, t.attribution)
	if len(lostLines) > 0 {
		add("attribution or licence line(s) missing or altered, they must be copied byte for byte: %s", sample(lostLines))
	}

	// Numbers written in the prose stay as digits. A translation may add some ("2nd" becomes
	// "2-й" and keeps its 2), so only the original's digits that it lacks are reported.
	_, missingDigits := diffMultiset(t.digitsOfProse(), e.digitsOfProse())
	if len(missingDigits) > 0 {
		add("number(s) of the prose are missing: %s", sample(missingDigits))
	}

	for _, p := range letterProblems(loc, t.prose) {
		add("prose %s", p)
	}
	if bad := untranslatedParagraphs(t.prose); len(bad) > 0 {
		add("%d paragraph(s) look untranslated: %s", len(bad), sample(bad))
	}
	if mixed := mixedScriptWords(t.prose); len(mixed) > 0 {
		add("word(s) that mix Latin and Cyrillic letters: %s", sample(mixed))
	}
	return problems
}

// digitsOfProse returns the digit runs of the prose, with the ones that the translation may
// legitimately add left to the caller.
func (f lessonFacts) digitsOfProse() []string { return digitRunRe.FindAllString(f.prose, -1) }

// untranslatedParagraphs returns the paragraphs that have several Latin words and no Cyrillic at
// all, which is what a paragraph left in English looks like once its formulas are taken out.
func untranslatedParagraphs(prose string) []string {
	var out []string
	for _, p := range blankLineRe.Split(prose, -1) {
		if len(latinWordRe.FindAllString(p, -1)) >= 4 && !hasCyrillic(p) {
			out = append(out, strings.Join(strings.Fields(p), " "))
		}
	}
	return out
}

// mixedScriptWords returns the words written with both Latin and Cyrillic letters, which in a
// Russian text is a Latin letter typed inside a Russian word.
func mixedScriptWords(prose string) []string {
	var out []string
	for _, w := range letterRunRe.FindAllString(prose, -1) {
		latin, cyr := false, false
		for _, r := range w {
			latin = latin || unicode.Is(unicode.Latin, r)
			cyr = cyr || unicode.Is(unicode.Cyrillic, r)
		}
		if latin && cyr {
			out = append(out, w)
		}
	}
	return out
}

// validateLessonTranslations checks every translated lesson under root (a data/lessons
// directory) in the given locales: the body against its English original, then the lesson
// against the steps of the concepts that read it. It returns the problems, each starting with
// the file, and the number of translated lessons it looked at.
func validateLessonTranslations(root string, locales []i18n.Locale) (problems []string, checked int) {
	report := func(file, format string, args ...interface{}) {
		problems = append(problems, file+": "+fmt.Sprintf(format, args...))
	}
	readRaw := func(p string) (string, bool) {
		b, err := os.ReadFile(p)
		return string(b), err == nil
	}

	var concepts map[string][]string // lesson source -> concept ids
	if raw, ok := readRaw(filepath.Join(root, "lessons.json")); ok {
		concepts = map[string][]string{}
		var mapping []struct {
			ConceptID string `json:"concept_id"`
			Source    string `json:"source"`
		}
		if err := jsonUnmarshal(raw, &mapping); err != nil {
			report("lessons.json", "%v", err)
		}
		for _, m := range mapping {
			concepts[m.Source] = append(concepts[m.Source], m.ConceptID)
		}
	}
	loader, loadErr := lessons.Load(root)

	for _, loc := range locales {
		for _, corpus := range lessonCorpora {
			dir := filepath.Join(root, corpus, string(loc))
			if _, err := os.Stat(dir); err != nil {
				continue
			}
			_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				rel, _ := filepath.Rel(dir, path)
				rel = filepath.ToSlash(rel)
				file := fmt.Sprintf("%s/%s/%s", corpus, loc, rel)
				if !strings.HasSuffix(rel, ".md") {
					report(file, "only .md lessons belong in a locale directory")
					return nil
				}
				source := corpus + "/" + rel
				original, ok := readRaw(filepath.Join(root, filepath.FromSlash(source)))
				if !ok {
					report(file, "there is no English lesson %s to translate", source)
					return nil
				}
				translation, _ := readRaw(path)
				checked++
				for _, p := range compareLessonBodies(loc, original, translation) {
					report(file, "%s", p)
				}
				for _, p := range stepSectionProblems(root, loc, source, translation, concepts[source], loader, loadErr) {
					report(file, "%s", p)
				}
				return nil
			})
		}
	}
	return problems, checked
}

// stepSectionProblems checks that the steps of every concept that reads the lesson find their
// Russian (or other) section as a heading of the translation, and that the server does serve
// the translation for it.
func stepSectionProblems(root string, loc i18n.Locale, source, translation string, ids []string, loader *lessons.Loader, loadErr error) []string {
	var problems []string
	add := func(format string, args ...interface{}) { problems = append(problems, fmt.Sprintf(format, args...)) }
	if len(ids) == 0 {
		return []string{"no concept reads this lesson (lessons.json), so a learner could never see the translation"}
	}
	if loadErr != nil {
		return []string{fmt.Sprintf("the lessons did not load: %v", loadErr)}
	}
	for _, id := range ids {
		var steps []kpStep
		raw, err := os.ReadFile(filepath.Join(root, "kp", string(loc), id+".json"))
		if err != nil {
			add("concept %s has no %s steps in kp/%s/%s.json", id, loc, loc, id)
			continue
		}
		if err := jsonUnmarshal(string(raw), &steps); err != nil {
			add("kp/%s/%s.json: %v", loc, id, err)
			continue
		}
		for i, s := range steps {
			if s.Section != "" && !lessons.SectionResolves(translation, s.Section) {
				add("concept %s step %d: its %s section %q is not a heading of this lesson", id, i, loc, s.Section)
			}
		}
		// The lesson the server picks for the concept may be another file; then this one is
		// never shown for it and there is nothing more to hold it to.
		if en := loader.Lesson(id); en != nil && en.File() == source && loader.LessonIn(id, loc).Locale() != loc {
			add("concept %s: the server would not serve this translation and falls back to English", id)
		}
	}
	return problems
}

func jsonUnmarshal(text string, v interface{}) error { return json.Unmarshal([]byte(text), v) }

func TestLessonTranslationsMirrorTheEnglishOnes(t *testing.T) {
	problems, checked := validateLessonTranslations(lessonsDir, i18n.Overlays())
	t.Logf("%d translated lessons checked", checked)
	for _, p := range problems {
		t.Error(p)
	}
}

// goodLesson is a small lesson with every feature the validator looks at, and goodTranslation is
// a correct Russian version of it.
const goodLesson = `> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Adding fractions

To add fractions, find a common denominator. See the [rules](https://example.org/rules).

## The idea

Rewrite $\frac{1}{2}$ as $\frac{2}{4}$, then add:

$$
\frac{2}{4}+\frac{1}{4}=\frac{3}{4}
$$

![A pie cut in four](svg/pie.svg)

Use ` + "`add(1, 2)`" + ` in 3 steps.

| Step | Result |
| --- | --- |
| one | $1$ |

- First, at most \(5\) pieces.
- Second, price \$5.

## Summary

Check \[x=2\].
`

const goodTranslation = `> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0

# Сложение дробей

Чтобы сложить дроби, найдите общий знаменатель. Смотрите [правила](https://example.org/rules).

## Идея

Перепишите $\frac{1}{2}$ как $\frac{2}{4}$, а затем сложите:

$$
\frac{2}{4}+\frac{1}{4}=\frac{3}{4}
$$

![Пирог, разрезанный на четыре части](svg/pie.svg)

Используйте ` + "`add(1, 2)`" + ` за 3 шага.

| Шаг | Результат |
| --- | --- |
| один | $1$ |

- Первое: не более \(5\) кусков.
- Второе: цена \$5.

## Итог

Проверьте \[x=2\].
`

func TestLessonValidatorAcceptsAFaithfulTranslation(t *testing.T) {
	if got := compareLessonBodies(i18n.RU, goodLesson, goodTranslation); len(got) != 0 {
		t.Errorf("a faithful translation was rejected: %v", got)
	}
	// Line endings are not content.
	crlf := strings.ReplaceAll(goodTranslation, "\n", "\r\n")
	if got := compareLessonBodies(i18n.RU, goodLesson, crlf); len(got) != 0 {
		t.Errorf("CRLF translation was rejected: %v", got)
	}
}

// Each broken translation below is the good one with one thing changed, and the validator has to
// name that thing.
func TestLessonValidatorRejectsBrokenTranslations(t *testing.T) {
	mutate := func(old, new string) string {
		if strings.Count(goodTranslation, old) != 1 {
			t.Fatalf("fixture bug: %q occurs %d times", old, strings.Count(goodTranslation, old))
		}
		return strings.Replace(goodTranslation, old, new, 1)
	}
	for _, tc := range []struct {
		name string
		tr   string
		want string
	}{
		{"altered inline formula", mutate(`$\frac{2}{4}$, а затем`, `$\frac{2}{5}$, а затем`), "math span"},
		{"translated word inside a formula", mutate(`\frac{3}{4}`, `\дробь{3}{4}`), "math span"},
		{"dropped display formula", mutate("$$\n\\frac{2}{4}+\\frac{1}{4}=\\frac{3}{4}\n$$\n\n", ""), "math span"},
		{"changed \\( \\) formula", mutate(`\(5\)`, `\(6\)`), "math span"},
		{"changed \\[ \\] formula", mutate(`\[x=2\]`, `\[x=3\]`), "math span"},
		{"dropped dollar escape", mutate(`цена \$5`, `цена $5`), "stray dollar"},
		{"altered inline code", mutate("`add(1, 2)`", "`сложить(1, 2)`"), "code"},
		{"dropped inline code", mutate("Используйте `add(1, 2)` за", "Используйте add за"), "code"},
		{"changed image path", mutate(`svg/pie.svg`, `svg/pirog.svg`), "target"},
		{"dropped image", mutate("![Пирог, разрезанный на четыре части](svg/pie.svg)\n\n", ""), "image"},
		{"changed link target", mutate(`https://example.org/rules`, `https://example.org/pravila`), "target"},
		{"dropped heading", mutate("## Итог\n\nПроверьте \\[x=2\\].\n", "Проверьте \\[x=2\\].\n"), "heading"},
		{"heading level changed", mutate("## Идея", "### Идея"), "heading"},
		{"dropped attribution line", mutate("> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0\n\n", ""), "attribution"},
		{"translated attribution line", mutate("Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0", "Содержание взято из [ORCCA](https://pcc.edu/orcca) — CC BY 4.0"), "attribution"},
		{"attribution with another dash", mutate("College Algebra) — CC BY", "College Algebra) - CC BY"), "attribution"},
		{"dropped table row", mutate("| один | $1$ |\n", ""), "table row"},
		{"dropped list item", mutate("- Второе: цена \\$5.\n", ""), "list item"},
		{"dropped paragraph", mutate("Чтобы сложить дроби, найдите общий знаменатель. Смотрите [правила](https://example.org/rules).\n\n", ""), "paragraph"},
		{"lost a number", mutate("за 3 шага", "за три шага"), "number"},
		{"not translated at all", goodLesson, "no ru text"},
		{"one paragraph left in English", mutate("Чтобы сложить дроби, найдите общий знаменатель. Смотрите [правила]", "To add fractions, find a common denominator. See the [rules]"), "untranslated"},
		{"Ukrainian letter in Russian prose", mutate("найдите общий знаменатель", "знайдіть общий знаменатель"), "not ru"},
		{"Latin letter inside a Russian word", mutate("Идея", "Иdея"), "mix"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := compareLessonBodies(i18n.RU, goodLesson, tc.tr)
			if len(got) == 0 {
				t.Fatalf("the broken translation was accepted")
			}
			if !strings.Contains(strings.Join(got, "\n"), tc.want) {
				t.Errorf("the problems do not mention %q:\n%s", tc.want, strings.Join(got, "\n"))
			}
		})
	}
}

// The same checks, run on files: a translated lesson that breaks the rules fails the whole
// validator, and so does one whose steps cannot find their headings.
func TestLessonValidatorFailsOnABrokenLessonsDirectory(t *testing.T) {
	root := t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("lessons.json", `[{"concept_id":"frac.add","source":"teaching/frac.add.md"}]`)
	write("teaching/frac.add.md", goodLesson)
	write("kp/frac.add.json", `[{"label":"Idea","section":"The idea","subgoals":["a","b","c"]},{"label":"Sum","section":"Summary","subgoals":["a","b","c"]}]`)
	write("kp/ru/frac.add.json", `[{"label":"Идея","section":"Идея","subgoals":["а","б","в"]},{"label":"Итог","section":"Итог","subgoals":["а","б","в"]}]`)

	// The faithful translation passes, and is served.
	write("teaching/ru/frac.add.md", goodTranslation)
	if problems, n := validateLessonTranslations(root, []i18n.Locale{i18n.RU}); n != 1 || len(problems) != 0 {
		t.Fatalf("faithful translation: %d checked, problems %v", n, problems)
	}

	// A step section that is not a heading of the translation.
	write("kp/ru/frac.add.json", `[{"label":"Идея","section":"Идея урока","subgoals":["а","б","в"]},{"label":"Итог","section":"Итог","subgoals":["а","б","в"]}]`)
	problems, _ := validateLessonTranslations(root, []i18n.Locale{i18n.RU})
	if got := strings.Join(problems, "\n"); !strings.Contains(got, `section "Идея урока" is not a heading`) || !strings.Contains(got, "falls back to English") {
		t.Errorf("a step section that matches no heading was not reported:\n%s", got)
	}
	write("kp/ru/frac.add.json", `[{"label":"Идея","section":"Идея","subgoals":["а","б","в"]},{"label":"Итог","section":"Итог","subgoals":["а","б","в"]}]`)

	// A broken body.
	write("teaching/ru/frac.add.md", strings.Replace(goodTranslation, `$\frac{1}{2}$`, `$\frac{1}{3}$`, 1))
	problems, _ = validateLessonTranslations(root, []i18n.Locale{i18n.RU})
	if got := strings.Join(problems, "\n"); !strings.Contains(got, "teaching/ru/frac.add.md") || !strings.Contains(got, "math span") {
		t.Errorf("a broken body was not reported:\n%s", got)
	}

	// A translation with no original, a lesson no concept reads, and a stray file.
	write("teaching/ru/frac.add.md", goodTranslation)
	write("teaching/ru/frac.orphan.md", goodTranslation)
	write("teaching/ru/notes.txt", "x")
	problems, _ = validateLessonTranslations(root, []i18n.Locale{i18n.RU})
	got := strings.Join(problems, "\n")
	for _, want := range []string{"there is no English lesson teaching/frac.orphan.md", "only .md lessons"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}
