package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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
