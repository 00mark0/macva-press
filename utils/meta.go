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
	plainText := ParseHTMLToText(content)
	sentences := strings.SplitAfter(plainText, ". ")

	var description string
	for _, sentence := range sentences {
		if len(description)+len(sentence) > 160 {
			break
		}
		description += sentence
	}

	return strings.TrimSpace(description)
}
