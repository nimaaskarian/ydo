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
	"time"
)

func addYdoToDir(base string) string {
	dir := filepath.Join(base, "ydo")
	if err := os.Mkdir(dir, 0755); err != nil {
		if errors.Is(err, os.ErrExist) {
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

func ParseDue(input string, now time.Time) (time.Time, error) {
	time, err := ParseDuration(input, now)
	if err != nil {
		time, err = parseDate(input, now)
	}
	return time, err
}

func ParseDuration(input string, now time.Time) (time.Time, error) {
	if input == "" {
		return time.Time{}, nil
	}
	coefficient := 1
	if strings.HasPrefix(input, "-") {
		coefficient = -1
		input = input[1:]
	}
	index := strings.IndexFunc(input, IsNotDigit)
	num, err := strconv.Atoi(input[:index])
	num *= coefficient
	if err != nil {
		return time.Time{}, err
	}
	var base time.Duration
	switch input[index:] {
	case "months", "mnths", "ms", "m":
		return AddDate(now, 0, num, 0), nil
	case "year", "y", "ys", "yrs":
		return AddDate(now, num, 0, 0), nil
	case "weeks", "wks", "w", "ws":
		base = time.Hour * 24 * 7
	case "day", "d", "days", "ds":
		base = time.Hour * 24
	case "hrs", "hours", "h", "hs":
		base = time.Hour
	case "minutes", "min", "mins":
		base = time.Minute
	case "s", "seconds", "scnds":
		base = time.Second
	default:
		return time.Time{}, errors.New("Invalid duration")
	}
	return now.Add(base * time.Duration(num)), nil
}

func parseDate(s string, now time.Time) (time.Time, error) {
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
			t = time.Date(0, 0, 0, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
		default:
			var err error
			for _, format := range [...]string{"15:04:05", "15:04", "15"} {
				t, err = time.Parse(format, date_time[1])
				if err == nil {
					break
				}
			}
			if err != nil {
				return time.Time{}, fmt.Errorf("Invalid time %q. Time is a string with format H:M:S, H:M or H", date_time[1])
			}
		}
		time_duration = time.Hour*time.Duration(t.Hour()) + time.Minute*time.Duration(t.Minute()) + time.Second*time.Duration(t.Second()) + time.Nanosecond*time.Duration(t.Nanosecond())
	}
	today_with_time := NaiveDate(now).Add(time_duration)
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
		// yeah. like you gonna do that in a thousand years
		return AddDate(now, 1000, 0, 0), nil
	default:
		date, err := ParseYmd(date_time[0], time.Local)
		if err != nil {
			return date, fmt.Errorf("Invalid date %q. Date is a Y-M-D, weekday, yesterday, today, tomorrow or later", date_time[0])
		}
		return date.Add(time.Duration(time_duration)), nil
	}
	day := 24 * time.Hour
	count_days := (7 + target_weekday - weekday) % 7
	if count_days == 0 {
		count_days = 7
	}
	return today_with_time.Add(time.Duration(count_days) * day), nil
}

func NaiveDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func NaiveDateEqual(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Year()
}

func makeDatePartString(value uint64, indicator string) string {
	if value > 0 {
		return fmt.Sprintf("%d%s", value, indicator)
	}
	return ""
}

func joinSpace(a, b string) string {
	if b == "" {
		return a
	} else {
		return a + " " + b
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

func DeepCopyMap[K comparable, V any](m map[K]*V) (out map[K]V) {
	out = make(map[K]V, len(m))
	var key K
	for key = range m {
		out[key] = *m[key]
	}
	return out
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// kinda copied from github.com/charmbracelet/x/editor's Cmd function, and
// debloated for my usecase
func EditorCmd(path string) (*exec.Cmd, error) {
	if os.Getenv("SNAP_REVISION") != "" {
		return nil, fmt.Errorf("Did you install with Snap? ydo is sandboxed and unable to open an editor. Please install ydo with Go or another package manager to enable editing.")
	}
	editor, args := getEditor()
	args = append(args, path)
	return exec.Command(editor, args...), nil
}

func getEditor() (string, []string) {
	editor := strings.Fields(os.Getenv("EDITOR"))
	if len(editor) > 1 {
		return editor[0], editor[1:]
	}
	if len(editor) == 1 {
		return editor[0], []string{}
	}
	return DEFAULT_EDITOR, []string{}
}

// set the command's stdout, stderr and stdin to os's
func CmdStdOs(c *exec.Cmd) {
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
}

func IsNotDigit(r rune) bool {
	return r > '9' || r < '0'
}

func ParseWeekday(s string) (time.Weekday, error) {
  for day := time.Sunday; day <= time.Saturday; day++ {
    day_s := day.String()
    if strings.HasPrefix(day_s, s) || strings.HasPrefix(strings.ToLower(day_s), s) {
      return day, nil
    }
  }
  return time.Sunday, errors.New("Not a weekday")
}

