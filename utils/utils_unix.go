// +build !windows,!darwin,unix

package utils

import (
  "path/filepath"
  "os"
  "os/exec"
	"syscall"
)

func ConfigDir() string {
  base := os.Getenv("XDG_CONFIG_HOME")
  if base == "" {
    base = filepath.Join(os.Getenv("HOME"), ".config")
  }
  return addYdoToDir(base)
}

func OpenURL(url string) error {
  cmd := exec.Command("xdg-open", url)
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Setpgid: true,
  }
  return cmd.Run()
}
