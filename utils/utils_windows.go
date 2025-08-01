package utils

import (
	"os"
	"os/exec"
)

func ConfigDir() string {
	return addYdoToDir(os.Getenv("APPDATA"))
}

func OpenURL(url string) error {
	return exec.Command("cmd.exe", "/C", "start "+url).Run()
}
