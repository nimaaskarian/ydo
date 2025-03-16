package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"runtime"
)

func ConfigDir() string {
  var base string;
  switch runtime.GOOS {
  case "windows":
    base = os.Getenv("APPDATA")
  case "darwin":
    base = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
  default:
    if base=os.Getenv("XDG_CONFIG_HOME"); base == "" {
      base = filepath.Join(os.Getenv("HOME"), ".config")
    }
  }
  dir := filepath.Join(base, "ydo")
  if err := os.Mkdir(dir, 0755); err != nil {
    if errors.Is(err, os.ErrExist)  {
      stat, _ := os.Stat(dir)
      if stat.IsDir() {
        return dir
      }
    }
    slog.Error("Couldn't create config directory. Using the current directory.", "config_dir", dir)
    return "."
  }
  return dir
}

func ReadYesNo(format string, a ...any) bool {
  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf(format, a...)
    line, err := reader.ReadString('\n')
    if err != nil {
      slog.Error("Error reading input:", "err", err)
      panic(err)
    }
    lower_line := strings.ToLower(strings.TrimSpace(line))
    if strings.HasPrefix("yes", lower_line) {
      return true
    }
    if strings.HasPrefix("no", lower_line) {
      return false
    }
  }
}

func ParseDate(date string) (time.Time, error) {
  now := time.Now()
  today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
  weekday := now.Weekday()
  var target_weekday time.Weekday

  switch strings.ToLower(date) {
  case "today":
    return today, nil
  case "tomorrow":
    return today.Add(24*time.Hour), nil
  case "sunday", "sun":
    target_weekday = time.Sunday
  case "monday", "mon":
    target_weekday = time.Monday
  case "tuesday", "tue":
    target_weekday = time.Tuesday
  case "wednesday", "wed":
    target_weekday = time.Wednesday
  case "thursday", "thu":
    target_weekday = time.Thursday
  case "friday", "fri":
    target_weekday = time.Friday
  case "saturday", "sat":
    target_weekday = time.Saturday
  default:
    return time.Parse("2006-01-02", date)
  }
  day := 24*time.Hour
  count_days := (7 + target_weekday-weekday) % 7
  if count_days == 0 {
    count_days = 7
  }
  return today.Add(time.Duration(count_days)*day), nil
}

const (
	SecondsInMinutes = 60
	SecondsInHour   = 3600
	SecondsInDay    = 86400
	SecondsInWeek   = 604800
)

func FormatDuration(diff time.Duration) string {
	rounded_seconds := int(math.Round(diff.Seconds()))
	rounded_minutes := rounded_seconds / SecondsInMinutes
	rounded_hours := rounded_minutes / SecondsInMinutes
	rounded_days := rounded_hours / 24
	rounded_weeks := rounded_days / 7
	if rounded_weeks > 0 {
		return strconv.Itoa(rounded_weeks) + "w"
	} else if rounded_days > 0 {
		return strconv.Itoa(rounded_days) + "d"
	} else if rounded_hours > 0 {
    minutes := rounded_minutes%60
    if minutes != 0 {
      return strconv.Itoa(rounded_hours) + "h" + strconv.Itoa(minutes) + "m"
    }
    return strconv.Itoa(rounded_hours) + "h"
	} else if rounded_minutes > 0 {
    seconds := rounded_seconds%SecondsInMinutes
    if seconds != 0 {
      return strconv.Itoa(rounded_minutes) + "m" + strconv.Itoa(seconds) + "s"
    }
    return strconv.Itoa(rounded_minutes) + "m"
	} else {
		return strconv.Itoa(rounded_seconds) + "s"
	}
}

func DeepCopyMap[K comparable, V any](m map[K]V) map[K]V {
  out := make(map[K]V, len(m))
  var key K
  for key = range m {
    out[key] = m[key]
  }
  return out
}
