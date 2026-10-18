package ai

import (
	"fmt"
	"strings"
)

// tutorSystemPrompt is instruction-only. It never carries learner data, so it is
// the same bytes for every learner and for every provider.
const tutorSystemPrompt = `You are Mathua's study tutor for mathematics.
Answer the learner's question in short, clear steps.
Use only the concept context you are given.
Do not invent grades, mastery numbers or progress the context does not state.
If the context is insufficient, say what is missing instead of guessing.
Never state a final numeric answer for a graded question; guide the method.`

// maxExcerptRunes bounds authored course text included in a prompt, so one long
// lesson cannot push a request past a provider's context window.
const maxExcerptRunes = 2000

// MasteryLine is one concept's mastery state as shown to the learner. It carries
// no identifier: a concept id and a state word are the same for every learner.
type MasteryLine struct {
	ConceptID string
	State     string
}

// TutorContext is the material a tutor request may use.
//
// Every identifier field is present on this struct on purpose: a caller that
// holds a learner's name, email or student id must hand it to the redactor and
// watch it dropped, rather than be trusted to leave it out. The zero value is a
// valid, empty context.
type TutorContext struct {
	// LearnerMessage is the question the learner typed.
	LearnerMessage string
	// LearnerName, LearnerEmail and StudentID are dropped by Messages. They are
	// declared so the redaction is testable and so a reviewer can see the
	// contract in the type.
	LearnerName  string
	LearnerEmail string
	StudentID    string
	// Domain and ConceptIDs identify what the learner is working on, not who
	// they are.
	Domain     string
	ConceptIDs []string
	// Mastery is the learner's current state per concept.
	Mastery []MasteryLine
	// LessonExcerpt is authored course text. It is public content, not learner
	// data.
	LessonExcerpt string
	// ExpectedAnswer is the graded answer. A model must never receive it, so
	// Messages drops it unconditionally. The field exists so that a caller
	// holding the answer cannot pass it through by accident, and so
	// TestTutorMessagesNeverLeakExpectedAnswer can prove the drop.
	ExpectedAnswer string
}

// Class is the data class of the messages Messages produces. A tutor request
// always includes the learner's own words and their progress, so it is
// pseudonymous: no identifier, but not public either.
func (c TutorContext) Class() DataClass { return DataPseudonymous }

// Messages renders the context as provider messages.
//
// It uses only the non-identifying fields. LearnerName, LearnerEmail,
// StudentID and ExpectedAnswer are ignored, which is what makes the function
// safe to call with a fully populated context from the server.
func (c TutorContext) Messages() []Message {
	var b strings.Builder
	if c.Domain != "" {
		fmt.Fprintf(&b, "Domain: %s\n", c.Domain)
	}
	if len(c.ConceptIDs) > 0 {
		fmt.Fprintf(&b, "Concepts: %s\n", strings.Join(c.ConceptIDs, ", "))
	}
	if len(c.Mastery) > 0 {
		b.WriteString("Mastery:\n")
		for _, m := range c.Mastery {
			fmt.Fprintf(&b, "- %s: %s\n", m.ConceptID, m.State)
		}
	}
	if c.LessonExcerpt != "" {
		fmt.Fprintf(&b, "Lesson excerpt:\n%s\n", truncateRunes(c.LessonExcerpt, maxExcerptRunes))
	}

	msgs := make([]Message, 0, 3)
	msgs = append(msgs, Message{Role: RoleSystem, Content: tutorSystemPrompt})
	if ctx := strings.TrimSpace(b.String()); ctx != "" {
		msgs = append(msgs, Message{Role: RoleSystem, Content: ctx})
	}
	msgs = append(msgs, Message{Role: RoleUser, Content: c.LearnerMessage})
	return msgs
}

// truncateRunes cuts s to at most n runes, appending an ellipsis when it was
// shortened. It counts runes rather than bytes so a multi-byte character is
// never split in half.
func truncateRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
