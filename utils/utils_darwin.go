// +build !windows,darwin,!unix

package utils

import (
  "path/filepath"
  "os"
)

func ConfigDir() string {
  base := filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
  return addYdoToDir(base)
}

func OpenURL(url string) error {
  return exec.Command("open", url).Run()
}
