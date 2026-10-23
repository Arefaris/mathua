package generator

import (
	"regexp"
	"strings"
)

// choiceOptionLine matches an option line that MultipleChoice appends to a stem.
//
// MultipleChoice writes each option as "\n\nC) text" with no newline inside the
// option, so a line-anchored pattern is enough and a stem that merely mentions
// "A)" mid-sentence cannot be mistaken for an option list.
var choiceOptionLine = regexp.MustCompile(`(?m)^([A-E])\)\s+(.+)$`)

// ChoiceOptions returns the options a multiple-choice stem carries, keyed by
// letter, or nil when the stem has no option list.
//
// This exists because the browser cannot be trusted to send an option's text
// back unchanged: tapping a choice fills the answer control, and a math field
// rewrites what it renders. Spaces disappear, LaTeX becomes glyphs and U+2212
// becomes a hyphen, so a correct tap was graded wrong for four different
// reasons (reports #154, #160, #163, #172, #179 on 23 October 2026). The client
// now sends the letter, and this is what turns it back into the authored text
// the answer is stored against.
//
// Letters must run A, B, C, … from the start, which is what MultipleChoice
// emits, so prose containing a stray "B)" is not read as a list.
func ChoiceOptions(question string) map[string]string {
	if question == "" {
		return nil
	}
	matches := choiceOptionLine.FindAllStringSubmatch(question, -1)
	if len(matches) < 2 {
		return nil
	}
	options := make(map[string]string, len(matches))
	for i, m := range matches {
		if m[1] != string(rune('A'+i)) {
			return nil
		}
		options[m[1]] = strings.TrimSpace(m[2])
	}
	return options
}
