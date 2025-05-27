package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func PrettyURL(slug string, publishedAt time.Time) string {
	return fmt.Sprintf("/%s/%s/%s",
		publishedAt.Format("2006"),
		publishedAt.Format("01"),
		slug,
	)
}

// Mapping for Serbian Latin characters to ASCII
var transliterationMap = map[string]string{
	// Serbian Latin
	"đ": "dj", "č": "c", "ć": "c", "š": "s", "ž": "z",
	"Đ": "dj", "Č": "c", "Ć": "c", "Š": "s", "Ž": "z",

	// Serbian Cyrillic
	"А": "a", "Б": "b", "В": "v", "Г": "g", "Д": "d", "Ђ": "dj", "Е": "e",
	"Ж": "z", "З": "z", "И": "i", "Ј": "j", "К": "k", "Л": "l", "Љ": "lj",
	"М": "m", "Н": "n", "Њ": "nj", "О": "o", "П": "p", "Р": "r", "С": "s",
	"Т": "t", "Ћ": "c", "У": "u", "Ф": "f", "Х": "h", "Ц": "c", "Ч": "c",
	"Џ": "dz", "Ш": "s",

	"а": "a", "б": "b", "в": "v", "г": "g", "д": "d", "ђ": "dj", "е": "e",
	"ж": "z", "з": "z", "и": "i", "ј": "j", "к": "k", "л": "l", "љ": "lj",
	"м": "m", "н": "n", "њ": "nj", "о": "o", "п": "p", "р": "r", "с": "s",
	"т": "t", "ћ": "c", "у": "u", "ф": "f", "х": "h", "ц": "c", "ч": "c",
	"џ": "dz", "ш": "s",
}

func transliterateSerbian(s string) string {
	for src, target := range transliterationMap {
		s = strings.ReplaceAll(s, src, target)
	}
	return s
}

func Slugify(s string) string {
	// Normalize curly quotes and similar junk
	s = strings.ReplaceAll(s, "„", "")
	s = strings.ReplaceAll(s, "“", "")
	s = strings.ReplaceAll(s, "”", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "’", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "‘", "")

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
