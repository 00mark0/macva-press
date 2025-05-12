package utils

import (
	"context"
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/microcosm-cc/bluemonday"
)

// HTMLComponent wraps sanitized HTML content as a templ component
type HTMLComponent struct {
	content string
}

// Render implements the templ.Component interface
func (h HTMLComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := io.WriteString(w, h.content)
	return err
}

// ParseHTML sanitizes HTML and returns it as a templ.Component
func ParseHTML(content string) templ.Component {
	// Create a policy that allows common HTML elements and attributes
	p := bluemonday.UGCPolicy()

	// Allow all the elements and attributes that TinyMCE generates
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6", "p", "ul", "ol", "li",
		"a", "img", "strong", "em", "blockquote", "hr", "div", "br")
	p.AllowAttrs("class", "data-start", "data-end").Globally()
	p.AllowAttrs("href", "target", "rel").OnElements("a")
	p.AllowAttrs("src", "alt", "width", "height").OnElements("img")

	// Sanitize the content
	sanitized := p.Sanitize(content)

	// Return as a custom component that implements templ.Component
	return HTMLComponent{content: sanitized}
}

// ParseHTMLToText sanitizes HTML and returns its plain text representation.
func ParseHTMLToText(content string) string {
	// Create a policy that allows common HTML elements and attributes
	p := bluemonday.UGCPolicy()

	// Allow the elements that TinyMCE typically generates
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6", "p", "ul", "ol", "li",
		"a", "img", "strong", "em", "blockquote", "hr", "div", "br")
	p.AllowAttrs("class", "data-start", "data-end").Globally()
	p.AllowAttrs("href", "target", "rel").OnElements("a")
	p.AllowAttrs("src", "alt", "width", "height").OnElements("img")

	// Sanitize the content first
	sanitized := p.Sanitize(content)

	// Use goquery to parse the sanitized HTML and extract text.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sanitized))
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		// Fallback: return sanitized content if parsing fails
		return sanitized
	}

	// Return the concatenated text from the document.
	return doc.Text()
}
