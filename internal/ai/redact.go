package ai

import (
	"fmt"
	"strings"
)

// assistantSystemPrompt is instruction-only. It never carries learner data, so
// it is the same bytes for every visitor and for every provider.
//
// It describes the product rather than a subject. The first version of this
// prompt addressed a study tutor and was fed concept ids and mastery lines;
// the feature it actually backs answers "what is Mathua" and "where do I
// start", so it is an assistant, and a prompt that says "tutor" invites the
// model to answer a maths question instead of the question that was asked.
//
// The last line is not decoration: the endpoint returns plain text and the app
// renders plain text, so Markdown and LaTeX reach the reader as asterisks and
// backslashes. A reply of four short sentences also keeps one request inside a
// budget that a free provider will actually serve.
const assistantSystemPrompt = `You are Mathua's assistant. Mathua teaches mathematics by concept: each concept has lessons, a short mastery check, and scheduled review.
Answer questions about the product and help a visitor work out where to start.
When the answer depends on what the person already knows, ask one clarifying question before you suggest anything.
Suggest only destinations from the list you are given. Never invent a course, feature, price or date.
When you name a destination or an area of the site, write its id in double square brackets, like [[calculus]], and the app turns that into a link. Use only ids from the list you are given.
Do not claim to know anyone's progress, results or grades.
Reply in plain prose, at most four short sentences. No Markdown, no LaTeX, no lists.`

// maxDescriptionRunes bounds one catalogue description. The list is authored
// product data, but a long description would still push a prompt past a small
// provider's context window, and the assistant needs the name, not the essay.
const maxDescriptionRunes = 200

// Destination is one entry of the public destination catalogue. It carries no
// learner data: the same list is shown to every visitor on the onboard screen.
type Destination struct {
	ID          string
	Name        string
	Description string
}

// Turn is one earlier message in the same conversation, replayed so a
// clarifying question from the assistant is answerable. Role is narrowed to the
// two sides of a conversation: a replayed system message would be a second
// instruction set the client controls, which is exactly what the fixed prompt
// above exists to prevent.
type Turn struct {
	Role    string
	Content string
}

// AssistantContext is the material an assistant request may use.
//
// Every identifier field is present on this struct on purpose: a caller that
// holds a visitor's name, email or student id must hand it to the redactor and
// watch it dropped, rather than be trusted to leave it out. The zero value is a
// valid, empty context.
type AssistantContext struct {
	// Message is the question the visitor typed.
	Message string
	// LearnerName, LearnerEmail and StudentID are dropped by Messages. They are
	// declared so the redaction is testable and so a reviewer can see the
	// contract in the type.
	LearnerName  string
	LearnerEmail string
	StudentID    string
	// Destinations is the public catalogue the assistant may suggest from, so
	// it can name a real place to start instead of inventing one.
	Destinations []Destination
	// History is the earlier turns of this conversation, oldest first. It is the
	// visitor's own words and the assistant's own replies, so it carries no data
	// the current message does not already carry.
	History []Turn
	// ExpectedAnswer is a graded answer. A model must never receive it, so
	// Messages drops it unconditionally. The field exists so that a caller
	// holding the answer cannot pass it through by accident, and so
	// TestAssistantMessagesNeverLeakExpectedAnswer can prove the drop.
	ExpectedAnswer string
}

// Class is the data class of the messages Messages produces. Even with no
// progress data at all, the request carries the visitor's own words, so it is
// pseudonymous: no identifier, but not public either.
func (c AssistantContext) Class() DataClass { return DataPseudonymous }

// Messages renders the context as provider messages.
//
// It uses only the non-identifying fields. LearnerName, LearnerEmail,
// StudentID and ExpectedAnswer are ignored, which is what makes the function
// safe to call with a fully populated context from the server.
func (c AssistantContext) Messages() []Message {
	msgs := make([]Message, 0, len(c.History)+3)
	msgs = append(msgs, Message{Role: RoleSystem, Content: assistantSystemPrompt})
	if catalogue := c.Catalogue(); catalogue != "" {
		msgs = append(msgs, Message{Role: RoleSystem, Content: catalogue})
	}
	for _, t := range c.History {
		switch t.Role {
		case string(RoleUser):
			msgs = append(msgs, Message{Role: RoleUser, Content: t.Content})
		case string(RoleAssistant):
			msgs = append(msgs, Message{Role: RoleAssistant, Content: t.Content})
		}
	}
	msgs = append(msgs, Message{Role: RoleUser, Content: c.Message})
	return msgs
}

// Catalogue renders the destination list as one system message. It is empty
// when the server has no catalogue, which is a valid state: the assistant then
// answers product questions and asks where the visitor wants to go, rather than
// naming a destination it cannot vouch for.
func (c AssistantContext) Catalogue() string {
	if len(c.Destinations) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Destinations you may suggest:\n")
	for _, d := range c.Destinations {
		fmt.Fprintf(&b, "- %s: %s", d.ID, d.Name)
		if d.Description != "" {
			fmt.Fprintf(&b, " — %s", truncateRunes(d.Description, maxDescriptionRunes))
		}
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
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
