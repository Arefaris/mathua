package generator

import "testing"

func TestChoiceOptionsParsesWhatMultipleChoiceEmits(t *testing.T) {
	correct := `\exists x \forall y (x \ge y)`
	stem := MultipleChoice("Which statement is true?", correct, "1 - P(A)", "down on the left", "$x^2<0$")

	options := ChoiceOptions(stem.Question)
	if len(options) != 4 {
		t.Fatalf("parsed %d options from %q, want 4", len(options), stem.Question)
	}
	for i, letter := range []string{"A", "B", "C", "D"} {
		if options[letter] == "" {
			t.Errorf("option %s is missing from %v", letter, options)
		}
		if i > 0 && options[letter] == options[string(rune('A'+i-1))] {
			t.Errorf("options %s and %s are identical: %q", letter, string(rune('A'+i-1)), options[letter])
		}
	}
	found := false
	for _, text := range options {
		if text == correct {
			found = true
		}
	}
	if !found {
		t.Errorf("the correct option %q is not among %v", correct, options)
	}
}

// A stem that is prose, not a list, must stay prose: a false positive here
// would let a learner answer a numeric question by typing a letter.
func TestChoiceOptionsIgnoresProse(t *testing.T) {
	cases := []string{
		"",
		"What is 2 + 2?",
		"Solve for x: A) is not an option here",
		"Choose one.\n\nA) only one option",
		"Choose one.\n\nB) starts at B\n\nC) then C",
	}
	for _, question := range cases {
		if got := ChoiceOptions(question); got != nil {
			t.Errorf("ChoiceOptions(%q) = %v, want nil", question, got)
		}
	}
}
