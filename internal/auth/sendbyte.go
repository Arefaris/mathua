package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	stdhtml "html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// SendByte is an HTTP email API. SMTP remains supported for deployments that
// already have credentials; SendByte is chosen when its key is present, since
// that is the provider this deployment uses.
//
// The API contract is a single POST:
//
//	POST https://api.sendbyte.africa/v1/emails
//	Authorization: Bearer <key>
//	{"from":"...","to":"...","subject":"...","html":"..."}
//
// A sandbox key reports Queued/Sent/Delivered through webhooks but never
// reaches an inbox. If reset mail still does not arrive with a key configured,
// check which kind of key it is before changing this code.

func sendByteConfigured() bool {
	return strings.TrimSpace(os.Getenv("SENDBYTE_API_KEY")) != ""
}

func sendByteBaseURL() string {
	if v := strings.TrimSpace(os.Getenv("SENDBYTE_API_URL")); v != "" {
		return strings.TrimRight(v, "/")
	}
	return "https://api.sendbyte.africa"
}

func sendByteFrom() string {
	if v := strings.TrimSpace(os.Getenv("SENDBYTE_FROM")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("SMTP_FROM"))
}

func sendSendByte(to, subject, body string) error {
	key := strings.TrimSpace(os.Getenv("SENDBYTE_API_KEY"))
	if key == "" {
		return errors.New("sendbyte: SENDBYTE_API_KEY is not set")
	}
	from := sendByteFrom()
	if from == "" {
		return errors.New("sendbyte: SENDBYTE_FROM is not set")
	}
	payload, err := json.Marshal(map[string]string{
		"from":    from,
		"to":      to,
		"subject": subject,
		"html":    plainTextToHTML(body),
	})
	if err != nil {
		return fmt.Errorf("sendbyte: encode request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, sendByteBaseURL()+"/v1/emails", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("sendbyte: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sendbyte: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		// The provider explains refusals (unverified sender, sandbox key, bad
		// address) in the body; surface it so the next attempt is diagnosable.
		snippet, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("sendbyte: %s: %s", res.Status, strings.TrimSpace(string(snippet)))
	}
	return nil
}

// plainTextToHTML keeps the mail templates readable as paragraphs. The bodies
// are authored as plain text with \r\n line breaks, and turning them into one
// <p> per non-empty line preserves the shape without an HTML template.
func plainTextToHTML(body string) string {
	var out strings.Builder
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out.WriteString("<p>")
		out.WriteString(stdhtml.EscapeString(line))
		out.WriteString("</p>")
	}
	return out.String()
}
