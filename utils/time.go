package utils

import (
	"fmt"
	"time"
)

var Loc, _ = time.LoadLocation("Europe/Belgrade")

func TimeAgo(t time.Time) string {
	now := time.Now().In(Loc)
	diff := now.Sub(t)

	switch {
	case diff < time.Second:
		return "сада"
	case diff < time.Minute:
		sec := int(diff.Seconds())
		if sec < 5 {
			return fmt.Sprintf("пре %d с", sec)
		}
		return fmt.Sprintf("пре %d с", sec)
	case diff < time.Hour:
		min := int(diff.Minutes())
		if min == 1 {
			return "пре 1 мин"
		}
		return fmt.Sprintf("пре %d мин", min)
	case diff < 24*time.Hour:
		h := int(diff.Hours())
		if h == 1 {
			return "пре 1 сат"
		} else if h < 5 {
			return fmt.Sprintf("пре %d сата", h)
		}
		return fmt.Sprintf("пре %d сат", h)
	case diff < 30*24*time.Hour:
		d := int(diff.Hours() / 24)
		if d == 1 {
			return "пре 1 д"
		}
		return fmt.Sprintf("пре %d д", d)
	case diff < 365*24*time.Hour:
		m := int(diff.Hours() / 24 / 30)
		if m == 1 {
			return "пре 1 мес"
		}
		return fmt.Sprintf("пре %d мес", m)
	default:
		y := int(diff.Hours() / 24 / 365)
		if y == 1 {
			return "пре 1 г"
		}
		return fmt.Sprintf("пре %d г", y)
	}
}
