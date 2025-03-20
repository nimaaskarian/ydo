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
    slog.Error("Couldn't create config directory. Using the current directory.", "config_dir", ".")
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

func ParseDate(s string, now time.Time) (time.Time, error) {
  if s == "" {
    return time.Time{}, nil
  }
  date_time := strings.Split(s, "/")
  var time_duration time.Duration
  if len(date_time) == 2 {
    var t time.Time
    switch date_time[1] {
    case "":
      t = time.Time{}
    case "now":
      t = time.Date(0,0,0, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
    default:
      var err error
      for _, format := range [...]string{"15:04:05", "15:04", "15"} {
        t, err = time.Parse(format, date_time[1])
        if err == nil {
          break
        }
      }
      if err != nil {
        return time.Time{}, fmt.Errorf("Invalid time %q. Time is a string with format HH:MM:SS, HH:MM or HH", date_time[1])
      }
    }
    time_duration = time.Hour*time.Duration(t.Hour()) + time.Minute*time.Duration(t.Minute()) + time.Second*time.Duration(t.Second()) + time.Nanosecond*time.Duration(t.Nanosecond())
  }
  today_with_time := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time_duration)
  weekday := now.Weekday()
  var target_weekday time.Weekday

  switch strings.ToLower(date_time[0]) {
  case "today":
    return today_with_time, nil
  case "tomorrow":
    return today_with_time.AddDate(0, 0, 1), nil
  case "yesterday":
    return today_with_time.AddDate(0, 0, -1), nil
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
  case "later":
    // max YYYY-MM-DD date possible
    return now.AddDate(1000, 0, 0), nil
  case "nyear":
    return today_with_time.AddDate(1, 0, 0), nil
  case "nmonth":
    return today_with_time.AddDate(0, 1, 0), nil
  case "nweek":
    return today_with_time.AddDate(0, 0, 7), nil
  default:
    date, err := time.Parse("2006-01-02", date_time[0])
    if err != nil {
      return date, fmt.Errorf("Invalid date %q. Date is a YY-MM-DD, weekday, yesterday, today, tomorrow or later", date_time[0])
    }
    return date.Add(time.Duration(time_duration)), nil
  }
  day := 24*time.Hour
  count_days := (7 + target_weekday-weekday) % 7
  if count_days == 0 {
    count_days = 7
  }
  return today_with_time.Add(time.Duration(count_days)*day), nil
}

func makeDatePartString(value uint64, indicator string) string {
  if value > 0 {
    return fmt.Sprintf("%d%s", value, indicator)
  }
  return ""
}

func joinSpace(a,b string) string {
  if b == "" {
    return a
  } else {
    return a+" "+b
  }
}

func FormatDuration(diff time.Duration) string {
  rounded_seconds := uint64(math.Round(diff.Seconds()))
  rounded_minutes := rounded_seconds / 60
  rounded_hours := rounded_minutes / 60
  rounded_days := rounded_hours / 24
  rounded_weeks := rounded_days / 7
  rounded_months := rounded_days / 30
  rounded_years := rounded_days / 365
  seconds := makeDatePartString(rounded_seconds%60, "s")
  minutes := makeDatePartString(rounded_minutes%60, "min")
  hours := makeDatePartString(rounded_hours%24, "h")
  days := makeDatePartString(rounded_days%7, "d")
  weeks := makeDatePartString(rounded_weeks%4, "w")
  months := makeDatePartString(rounded_months%12, "m")
  years := makeDatePartString(rounded_years, "y")
  if rounded_years > 0 {
    return joinSpace(years, months)
  }
  if rounded_months > 0 {
    return joinSpace(months, makeDatePartString(rounded_days%30, "d"))
  }
  if rounded_weeks > 0 {
    return joinSpace(weeks, days)
  }
  if rounded_days > 0 {
    return joinSpace(days, hours)
  }
  if rounded_hours > 0 {
    return joinSpace(hours, minutes)
  }
  if rounded_minutes > 0 {
    return joinSpace(minutes, seconds)
  }
  return fmt.Sprintf("%ds", rounded_seconds)
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
