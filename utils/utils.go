package utils

import (
	"os"
	"path/filepath"

	"github.com/nimaaskarian/ydo/core"
	"runtime"
)

const APP_NAME = "ydo";

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
  return filepath.Join(base, APP_NAME)
}

func LoadTodos(path string) core.TodoMap {
  todomap := make(core.TodoMap)
  content, _ := os.ReadFile(path)
  core.ParseYaml(todomap, content)
  return todomap
}

func Check(e error) {
  if e != nil {
    panic(e)
  }
}
