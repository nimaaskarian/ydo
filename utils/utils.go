package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
  log.Printf("File %q loaded\n", path)
  taskmap := make(core.TaskMap)
  content, _ := os.ReadFile(path)
  core.ParseYaml(taskmap, content)
  return taskmap
}

func Check(e error) {
  if e != nil {
    panic(e)
  }
}

func MustHaveTask(taskmap core.TaskMap, key string) {
  if !taskmap.HasTask(key) {
    log.Fatalf("No such task %q\n",key)
  }
}

func MustNotHaveTask(taskmap core.TaskMap, key string) {
  if taskmap.HasTask(key) {
    log.Fatalf("Task %q already exists\n",key)
  }
}

func WriteTaskmap(taskmap core.TaskMap, path string) {
  content, err := yaml.Marshal(taskmap)
  Check(err)
  err = os.WriteFile(path, content, 0644)
  Check(err)
  log.Println("Wrote to file")
}

func ReadYesNo(format string, a ...any) bool {
  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf(format, a...)
    line, err := reader.ReadString('\n')
    if err != nil {
      log.Fatalln("Error reading input:", err)
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
