package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDateEmpty(t *testing.T) {
  date, err := ParseDate("", time.Time{})
  assert.Nil(t, err)
  assert.True(t, date.IsZero())
}

func TestParseDateInvalid(t *testing.T) {
  date, err := ParseDate("tomorrow/invalid", time.Time{})
  assert.ErrorContains(t, err, "Invalid time")
  assert.True(t, date.IsZero())

  date, err = ParseDate("2025/12/12/invalid", time.Time{})
  assert.ErrorContains(t, err, "Invalid date")
  assert.True(t, date.IsZero())

  date, err = ParseDate("2025/12/12/invalid", time.Time{})
  assert.ErrorContains(t, err, "Invalid date")
  assert.True(t, date.IsZero())

  date, err = ParseDate("invalid/invalid", time.Time{})
  assert.ErrorContains(t, err, "Invalid time")
  assert.True(t, date.IsZero())
}

func TestParseDateAbsolute(t *testing.T) {
  actual, err := ParseDate("2025-12-19/8", time.Time{})
  assert.Nil(t, err)
  expected, _ := time.Parse("2006-01-02 15:04:05","2025-12-19 8:00:00")
  assert.Equal(t, expected, actual)
}

func TestParseDateRelative(t *testing.T) {
  now, _ := time.Parse("2006-01-02 15:04:05","2025-03-20 17:00:00")
  tests := [...][3]string{
    {"tomorrow/8", "2025-03-21 8:00:00"},
    {"today/8", "2025-03-20 8:00:00"},
    {"yesterday/", "2025-03-19 00:00:00"},
    {"fri/8", "2025-03-21 8:00:00"},
    {"thu/8:00:20", "2025-03-27 8:00:20"},
    {"sat/18:43:21", "2025-03-22 18:43:21"},
    {"sat/now", "2025-03-22 17:00:00"},
    {"sat", "2025-03-22 00:00:00"},
    {"sun", "2025-03-23 00:00:00"},
    {"mon", "2025-03-24 00:00:00"},
    {"tue", "2025-03-25 00:00:00"},
    {"wed", "2025-03-26 00:00:00"},
    {"2025-12-19/now", "2025-12-19 17:00:00"},
    {"nyear/8", "2026-03-20 8:00:00"},
    {"nmonth/8", "2025-04-20 8:00:00"},
    {"nweek/8", "2025-03-27 8:00:00"},
    {"later", "3025-03-20 17:00:00"},
  }
  for _, test := range tests {
    actual, err := ParseDate(test[0], now)
    assert.Nil(t, err)
    expected, _ := time.Parse("2006-01-02 15:04:05", test[1])
    if !assert.Equal(t, expected, actual) {
      break
    }
  }
}

func TestFormatDuration(t *testing.T) {
  expected := FormatDuration(time.Hour + time.Second * 59)
  assert.Equal(t, "1h", expected)
  expected = FormatDuration(time.Hour + time.Second*32 + time.Minute*12)
  assert.Equal(t, "1h 12min", expected)

  expected = FormatDuration(time.Hour*49 + time.Second*32 + time.Minute*12)
  assert.Equal(t, "2d 1h", expected)

  expected = FormatDuration(time.Hour*72 + time.Second*32 + time.Minute*12)
  assert.Equal(t, "3d", expected)

  expected = FormatDuration(time.Hour*245 + time.Second*32 + time.Minute*12)
  assert.Equal(t, "1w 3d", expected)
  expected = FormatDuration(time.Hour*3600)
  assert.Equal(t, "5m", expected)
  expected = FormatDuration(time.Hour*3600*2)
  assert.Equal(t, "10m", expected)
  expected = FormatDuration(time.Hour*7920)
  assert.Equal(t, "11m", expected)
  expected = FormatDuration(time.Hour*7944)
  assert.Equal(t, "11m 1d", expected)
  expected = FormatDuration(time.Hour*12360 + time.Second*32 + time.Minute*12 + 24*time.Hour)
  assert.Equal(t, "1y 5m", expected)
  expected = FormatDuration(time.Second*10)
  assert.Equal(t, "10s", expected)
  expected = FormatDuration(time.Minute*10 + time.Second*10)
  assert.Equal(t, "10min 10s", expected)
  expected = FormatDuration(0)
  assert.Equal(t, "0s", expected)
}

