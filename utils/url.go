package utils

import (
	"regexp"
	"strings"
)

// Mapping for Serbian Latin characters to ASCII
var transliterationMap = map[string]string{
	"đ": "dj", "č": "c", "ć": "c", "š": "s", "ž": "z",
	"Đ": "dj", "Č": "c", "Ć": "c", "Š": "s", "Ž": "z",
}

func transliterateSerbian(s string) string {
	for src, target := range transliterationMap {
		s = strings.ReplaceAll(s, src, target)
	}
	return s
}

func Slugify(s string) string {
	s = transliterateSerbian(s)
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	// Optional: limit length
	if len(s) > 80 {
		s = s[:80]
		s = strings.TrimRight(s, "-")
	}

	return s
}
