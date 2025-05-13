package utils

import "strings"

func GenerateTitleTag(title string) string {
	maxLen := 60
	if len(title) <= maxLen {
		return title
	}

	// Try to cut on a word boundary
	trimmed := title[:maxLen]
	lastSpace := strings.LastIndex(trimmed, " ")
	if lastSpace > 0 {
		trimmed = trimmed[:lastSpace]
	}

	return trimmed + "..."
}

func GenerateMetaDescription(content string) string {
	// Fallback if content is empty or only whitespace
	if strings.TrimSpace(content) == "" {
		return "- МАЧВА ПРЕС БОГАТИЋ - NOVOSTI IZ MAČVE - BOGATIĆ - ŠABAC"
	}

	// Try to parse HTML content to plain text
	plainText := ParseHTMLToText(content)

	// If parsing fails or returns an empty string, fallback to a default
	if strings.TrimSpace(plainText) == "" {
		return "- МАЧВА ПРЕС БОГАТИЋ - NOVOSTI IZ MAČVE - BOGATIĆ - ŠABAC"
	}

	// Split the plain text into sentences
	sentences := strings.SplitAfter(plainText, ". ")

	// Initialize description variable
	var description string

	// Add sentences to the description while keeping it under the 160 char limit
	for _, sentence := range sentences {
		if len(description)+len(sentence) > 160 {
			break
		}
		description += sentence
	}

	// If we couldn't build a description, fall back to the default
	if len(description) == 0 {
		return "- МАЧВА ПРЕС БОГАТИЋ - NOVOSTI IZ MAČVE - BOGATIĆ - ŠABAC"
	}

	// Return the description trimmed to remove any extra spaces at the ends
	return strings.TrimSpace(description)
}
