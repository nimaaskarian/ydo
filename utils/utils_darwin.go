//go:build !windows && darwin && !unix
// +build !windows,darwin,!unix

package utils

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	base := filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	return addYdoToDir(base)
}

func OpenURL(url string) error {
	return exec.Command("open", url).Run()
}
