package main
import (
  "github.com/nimaaskarian/ydo/core"
  "github.com/nimaaskarian/ydo/utils"
  "os"
	"path/filepath"
  "fmt"
)

func main() {
  default_todo := core.Todo{};
  dir := utils.ConfigDir()
  path := filepath.Join(dir, "default-todo.yaml")
  content, _ := os.ReadFile(path)
  core.ParseYaml(&default_todo, content)
  fmt.Printf("%v\n", default_todo)
  path = filepath.Join(dir, "todos.yaml")
  content, _ = os.ReadFile(path)

  todomap := make(core.TodoMap)
  core.ParseYaml(todomap, content)
  fmt.Printf("%v\n", todomap)
}
