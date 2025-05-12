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
		return "sada"
	case diff < time.Minute:
		sec := int(diff.Seconds())
		if sec < 5 {
			return fmt.Sprintf("pre %d s", sec)
		}
		return fmt.Sprintf("pre %d s", sec)
	case diff < time.Hour:
		min := int(diff.Minutes())
		if min == 1 {
			return "pre 1 min"
		}
		return fmt.Sprintf("pre %d min", min)
	case diff < 24*time.Hour:
		h := int(diff.Hours())
		if h == 1 {
			return "pre 1 h"
		} else if h < 5 {
			return fmt.Sprintf("pre %d h", h)
		}
		return fmt.Sprintf("pre %d h", h)
	case diff < 30*24*time.Hour:
		d := int(diff.Hours() / 24)
		if d == 1 {
			return "pre 1 d"
		}
		return fmt.Sprintf("pre %d d", d)
	case diff < 365*24*time.Hour:
		m := int(diff.Hours() / 24 / 30)
		if m == 1 {
			return "pre 1 m"
		}
		return fmt.Sprintf("pre %d m", m)
	default:
		y := int(diff.Hours() / 24 / 365)
		if y == 1 {
			return "pre 1 g"
		}
		return fmt.Sprintf("pre %d g", y)
	}
}
