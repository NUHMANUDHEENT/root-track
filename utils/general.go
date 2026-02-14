package utils

import "time"

func SameDate(a, b time.Time) bool {
	yearA, monthA, dayA := a.Date()
	yearB, monthB, dayB := b.Date()

	return yearA == yearB && monthA == monthB && dayA == dayB
}
