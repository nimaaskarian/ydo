package utils

import (
	"log"
	"os"
	"path/filepath"

	"runtime"

	"github.com/nimaaskarian/ydo/core"
	"gopkg.in/yaml.v3"
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

func LoadTasks(path string) core.TaskMap {
  todomap := make(core.TaskMap)
  content, _ := os.ReadFile(path)
  core.ParseYaml(todomap, content)
  return todomap
}

func Check(e error) {
  if e != nil {
    panic(e)
  }
}

func MustHaveTask(taskmap core.TaskMap, key string) {
  if _, ok := taskmap[key]; !ok {
    log.Fatalf("No such todo %q\n",key)
  }
}

func MustNotHaveTask(taskmap core.TaskMap, key string) {
  if _, ok := taskmap[key]; ok {
    log.Fatalf("Task %q already exists\n",key)
  }
}

func WriteTaskmap(taskmap core.TaskMap, path string) {
    content, err := yaml.Marshal(taskmap)
    Check(err)
    err = os.WriteFile(path, content, 0644)
    Check(err)

}
