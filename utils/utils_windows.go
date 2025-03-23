// +build windows,!darwin,!unix

package utils

import (
	"os/exec"
  "os"
)

func ConfigDir() string {
  return addYdoToDir(os.Getenv("APPDATA"))
}

func OpenURL(url string) error {
  return exec.Command("cmd.exe", "/C", "start "+url).Run()
}
