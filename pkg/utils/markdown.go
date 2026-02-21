package utils

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	md     goldmark.Markdown
	policy *bluemonday.Policy
)

func init() {
	// Configure Goldmark with common extensions
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,         // GitHub Flavored Markdown (tables, strikethrough, autolinks, task lists)
			extension.Typographer, // Smart quotes, dashes, etc.
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Auto-generate heading IDs
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // Convert newlines to <br>
			html.WithXHTML(),     // XHTML-compliant output
		),
	)

	// Configure Bluemonday with a strict UGC (User Generated Content) policy
	// This allows common formatting but strips dangerous elements
	policy = bluemonday.UGCPolicy()

	// Allow additional safe attributes for better formatting
	policy.AllowAttrs("class").OnElements("code", "pre", "span", "div")
	policy.AllowAttrs("id").OnElements("h1", "h2", "h3", "h4", "h5", "h6")
}

// MarkdownToHTML converts Markdown text to sanitized HTML.
// It uses Goldmark for parsing and Bluemonday for sanitization.
// Safe for rendering user-generated content.
func MarkdownToHTML(markdown string) string {
	if markdown == "" {
		return ""
	}

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		// On error, return empty string (or could return escaped plain text)
		return ""
	}

	// Sanitize the generated HTML to prevent XSS
	return policy.Sanitize(buf.String())
}
