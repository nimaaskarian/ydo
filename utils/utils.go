package utils

import (
"bufio"
"errors"
"fmt"
"log/slog"
"math"
"os"
"os/exec"
"path/filepath"
"strconv"
"strings"
"syscall"
"time"

"runtime"
)

const (
  Windows = "windows"
  Darwin = "darwin"
)

const (
  IsWindows = runtime.GOOS == Windows
  IsDarwin  = runtime.GOOS == Darwin
)

func ConfigDir() string {
  var base string;
  if IsWindows {
    base = os.Getenv("APPDATA")
  } else if IsDarwin {
    base = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
  } else {
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

func ParseDate(date string) time.Time {
  now := time.Now()
  today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
  weekday := now.Weekday()
  var target_weekday time.Weekday

  switch strings.ToLower(date) {
  case "today":
    return today
  case "tomorrow":
    return today.Add(24*time.Hour)
  case "yesterday":
    return today.Add(-24*time.Hour)
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
    date, err := time.Parse("2006-01-02", date)
    if err != nil {
      slog.Error("Due date specified is invalid. Use YYYY-MM-DD format, today, tomorrow or a week day")
      os.Exit(1)
    }
    return date
  }
  day := 24*time.Hour
  count_days := (7 + target_weekday-weekday) % 7
  if count_days == 0 {
    count_days = 7
  }
  return today.Add(time.Duration(count_days)*day)
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

func DeepCopyMap[K comparable, V any](m map[K]V) (out map[K]V) {
  out = make(map[K]V, len(m))
  var key K
  for key = range m {
    out[key] = m[key]
  }
  return out
}

func OpenURL(url string) error {
  if IsWindows {
    return exec.Command("cmd.exe", "/C", "start "+url).Run()
  }
  if IsDarwin {
    return exec.Command("open", url).Run()
  }
  cmd := exec.Command("xdg-open", url)
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Setpgid: true,
  }
  return cmd.Run()
}

func Filter[T any] (arr[]T, test func(T) bool) (out []T) {
  for _, item := range arr {
    if test(item) {
      out = append(out, item)
    }
  }
  return out
}
