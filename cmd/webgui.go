package cmd

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"log/slog"
	"net"
	"net/http"
	"reflect"
	"slices"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

//go:embed webgui/templates/*
//go:embed webgui/static
var embed_fs embed.FS

var (
  address string
  port int
  changed bool
)

var tmpls *template.Template
func init() {
  rootCmd.AddCommand(webguiCmd)
  webguiCmd.Flags().StringVarP(&address, "address", "a", "127.0.0.1", "address to listen and serve to")
  webguiCmd.Flags().IntVarP(&port, "port", "p", 8485, "port to listen and serve to")
}

func updateChanged() {
  if !changed {
    changed = !reflect.DeepEqual(old_taskmap, taskmap)
  }
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  updateChanged()
  err := tmpls.ExecuteTemplate(w, "index.html", map[string]any { "Taskmap": taskmap,
    "SeenKeys": make(map[string]bool, len(taskmap) ),
    "Keys": sorted_keys,
    "Url": r.URL.String(),
    "Changed": changed,
  })
  if err != nil {
    slog.Error("Error in executing the template", "err", err)
  }
}

type TaskData struct {
  Taskmap core.TaskMap
  Key string
  SeenKeys map[string]bool
  Filter map[string]bool
}

func dict(args ...any) map[string]any {
  m := make(map[string]any)
  for i := 0; i < len(args); i += 2 {
      m[args[i].(string)] = args[i+1]
  }
  return m
}

func add2map(m map[string]any, args ...any) map[string]any {
  for i := 0; i < len(args); i += 2 {
      m[args[i].(string)] = args[i+1]
  }
  return m
}


func see(seen map[string]bool, key string) string {
  seen[key] = true
  return ""
}

func Task(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  w.Header().Add("Content-Type", "text/html")
  err := tmpls.ExecuteTemplate(w, "task.html", TaskData { Key: ps.ByName("key"), Taskmap: taskmap, SeenKeys: map[string]bool{} })
  if err != nil {
    slog.Error("Error in executing the template", "err", err)
  }
}

func makeFilter(f func (t core.Task, tm core.TaskMap) bool) map[string]bool {
  filter := make(map[string]bool, len(taskmap))
  for key, task := range taskmap {
    filter[key] = f(task, taskmap)
  }
  return filter
}

func DoTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  key := ps.ByName("key")
  task := taskmap[key]
  if !task.AutoComplete {
    taskmap.Do(key)
  }

  url := r.URL.Query().Get("redirect")
  if url == "" {
    url = "/"
  }
  w.Header().Add("HX-Location", url)
}

func Todo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  updateChanged()
    err := tmpls.ExecuteTemplate(w, "index.html", map[string]any { "Taskmap": taskmap,
    "SeenKeys": make(map[string]bool, len(taskmap) ),
    "Keys": sorted_keys,
    "Url": r.URL.String(),
    "Changed": changed,
    "Filter": makeFilter(core.Task.IsNotDone),
  })
  if err != nil {
    slog.Error("Error in executing the template", "err", err)
  }
}

func Write(cmd *cobra.Command) func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    rootCmd.PersistentPostRun(cmd, []string{})
    old_taskmap = utils.DeepCopyMap(taskmap)
    changed = false
  }
}

func UndoTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  key := ps.ByName("key")
  taskmap.Undo(key)
  url := r.URL.Query().Get("redirect")
  if url == "" {
    url = "/"
  }
  w.Header().Add("HX-Location", url)
}

var sorted_keys []string
var webguiCmd = &cobra.Command{
  Use: "webgui",
  Short: "run the webgui on the current file",
  Run: func(cmd *cobra.Command, args []string) {
    func_map := template.FuncMap{
      "dict": dict,
      "see": see,
      "add2map": add2map,
    }
    tmpls = template.Must(template.New("").Funcs(func_map).ParseFS(embed_fs, "webgui/templates/*"))
    sorted_keys = make([]string, 0, len(taskmap))
    for key := range taskmap {
      sorted_keys = append(sorted_keys, key)
    }
    slices.SortFunc(sorted_keys, func(k1, k2 string) int {
      return taskmap[k1].CreatedAt.Compare(taskmap[k2].CreatedAt)
    })
    server_root, err := fs.Sub(embed_fs, "webgui/static")
    if err != nil {
      slog.Error("Error happend when subbing embed_fs", "err", err)
    }
    static_server := http.FileServer(http.FS(server_root))

    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/todo", Todo)
    router.GET("/task/:key", Task)
    router.PUT("/write", Write(cmd))
    router.PUT("/do/:key", DoTask)
    router.PUT("/undo/:key", UndoTask)
    router.Handler("GET", "/static/*filepath", http.StripPrefix("/static/", static_server))

    address = net.JoinHostPort(address, strconv.Itoa(port))
    go utils.OpenURL("http://"+address)
    log.Fatal(http.ListenAndServe(address, router))
  },
}
