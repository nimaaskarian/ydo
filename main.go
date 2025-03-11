package main
import (
  "github.com/nimaaskarian/ydo/core"
  "fmt"
)

func main() {
  todomap := make(core.TodoMap)
  todomap.YamlParseMap(
`
1:
  message: hello
  priority: 1
  deps: [2]
  done: true
2:
  message: how you doing
  priority: 12
3:
  message: doing fine?
`);
  val:=todomap.NextId();
  fmt.Printf("%v %d\n", todomap, val);
}
