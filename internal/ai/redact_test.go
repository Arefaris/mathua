package ai

import (
	"strings"
	"testing"
)

func joinMessages(msgs []Message) string {
	var b strings.Builder
	for _, m := range msgs {
		b.WriteString(m.Content)
		b.WriteString("\n")
	}
	return b.String()
}

func TestTutorMessagesNeverLeakExpectedAnswer(t *testing.T) {
	const sentinel = "EXPECTED-ANSWER-SENTINEL-7f3"
	c := TutorContext{
		LearnerMessage: "how do I solve this?",
		ExpectedAnswer: sentinel,
		ConceptIDs:     []string{"algebra.linear.one-variable"},
	}
	for _, m := range c.Messages() {
		if strings.Contains(m.Content, sentinel) {
			t.Fatalf("expected answer leaked into a %s message: %q", m.Role, m.Content)
		}
	}
}

func TestTutorMessagesDropIdentifiers(t *testing.T) {
	c := TutorContext{
		LearnerMessage: "help with fractions",
		LearnerName:    "NAME-SENTINEL-Alice",
		LearnerEmail:   "EMAIL-SENTINEL-alice@example.test",
		StudentID:      "ID-SENTINEL-12345",
		Domain:         "algebra",
	}
	got := joinMessages(c.Messages())
	for _, forbidden := range []string{"NAME-SENTINEL-Alice", "EMAIL-SENTINEL-alice@example.test", "ID-SENTINEL-12345"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("identifier %q leaked into tutor messages", forbidden)
		}
	}
}

func TestTutorMessagesKeepLearningContext(t *testing.T) {
	c := TutorContext{
		LearnerMessage: "explain completing the square",
		Domain:         "algebra",
		ConceptIDs:     []string{"algebra.quadratic.completing-square", "algebra.quadratic.roots"},
		Mastery:        []MasteryLine{{ConceptID: "algebra.quadratic.roots", State: "learning"}},
		LessonExcerpt:  "A quadratic has the form ax^2 + bx + c = 0.",
	}
	msgs := c.Messages()
	if len(msgs) < 3 {
		t.Fatalf("want system, context and user messages, got %d", len(msgs))
	}
	if msgs[0].Role != RoleSystem {
		t.Errorf("first message role = %s, want system", msgs[0].Role)
	}
	if last := msgs[len(msgs)-1]; last.Role != RoleUser || last.Content != "explain completing the square" {
		t.Errorf("last message = %+v", last)
	}
	got := joinMessages(msgs)
	for _, want := range []string{
		"algebra",
		"algebra.quadratic.completing-square",
		"algebra.quadratic.roots: learning",
		"ax^2 + bx + c = 0",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("tutor messages missing %q", want)
		}
	}
}

func TestTutorMessagesOmitEmptyContextBlock(t *testing.T) {
	c := TutorContext{LearnerMessage: "hello"}
	msgs := c.Messages()
	if len(msgs) != 2 {
		t.Fatalf("want exactly system+user when there is no context, got %d: %+v", len(msgs), msgs)
	}
}

func TestTutorContextClassIsPseudonymous(t *testing.T) {
	if got := (TutorContext{}).Class(); got != DataPseudonymous {
		t.Errorf("TutorContext class = %s, want pseudonymous", got)
	}
}

func TestTruncateRunes(t *testing.T) {
	if got := truncateRunes("abc", 5); got != "abc" {
		t.Errorf("short string changed: %q", got)
	}
	if got := truncateRunes("abcdef", 3); got != "abc..." {
		t.Errorf("truncation = %q, want abc...", got)
	}
	// A multi-byte rune must not be split.
	if got := truncateRunes("aeiou", 2); got != "ae..." {
		t.Errorf("runes = %q", got)
	}
	if got := truncateRunes("abc", 0); got != "" {
		t.Errorf("zero limit = %q", got)
	}
}
