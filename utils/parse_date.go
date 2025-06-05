//go:build !jalali
// +build !jalali

package utils

import (
	"time"
)

func ParseYmd(s string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(time.DateOnly, s, loc)
}

func AddDate(t time.Time, years, months, days int) time.Time {
	return t.AddDate(years, months, days)
}
