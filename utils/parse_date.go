//go:build !jalali
// +build !jalali

package utils

import (
	"time"
)
func parseYmd(s string, loc *time.Location) (time.Time, error){
  return time.ParseInLocation("2006-01-02", s, loc)
}

