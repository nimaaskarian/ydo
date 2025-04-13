package gtasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
	"gopkg.in/yaml.v3"
)

type Gtasks struct {
	CredentialsFile string
	TokenFile       string
	CacheFile       string
	CacheExpire     string
	srv             *tasks.Service
	Cache           GtasksCache
}

type GtasksCache struct {
	Tasklists   map[string]*tasks.TaskList             `yaml:",omitempty"`
	Tasks       map[string][]*tasks.Task               `yaml:",omitempty"`
	Time        time.Time                              `yaml:",omitempty"`
	EditedLists EditedIdCache[string, *tasks.TaskList] `yaml:",omitempty"`
	EditedTasks EditedIdCache[[2]string, CreatedTask]  `yaml:",omitempty"`
}
type CreatedTask struct {
	Task *tasks.Task
	Id   string
}

type EditedIdCache[T any, E any] struct {
	Deleted []T
	Updated []T
	Created []E
}

func (gtasks *Gtasks) Login() error {
	c, err := gtasks.ReadCredentials()
	if err != nil {
		return err
	}
	client := gtasks.getClient(c)
	srv, err := tasks.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return errors.New("Unable to retrieve tasks Client\n" + err.Error())
	}
	gtasks.srv = srv
	return nil
}

func (gtasks *Gtasks) SaveCache(now time.Time) error {
	slog.Debug("Saving yaml gtasks cache file", "path", gtasks.CacheFile)
	gtasks.Cache.Time = now
	content, err := yaml.Marshal(gtasks.Cache)
	if err != nil {
		return err
	}
	if err := os.WriteFile(gtasks.CacheFile, content, 0644); err != nil {
		slog.Error("Failed writing the to gtasks cache file.", "path", gtasks.CacheFile)
		return err
	}
	return nil
}

func (gtasks *Gtasks) ReadCache(now time.Time) error {
	slog.Debug("Loading cache file", "path", gtasks.CacheFile)
	content, _ := os.ReadFile(gtasks.CacheFile)
	if err := yaml.Unmarshal(content, &gtasks.Cache); err != nil {
		return err
	}
	date, err := utils.ParseDuration(gtasks.CacheExpire, gtasks.Cache.Time)
	if err != nil {
		panic("Cache expire format is incorrect")
	}
	if date.Before(now) {
		gtasks.Cache = GtasksCache{}
		return errors.New("Cache is older than one day")
	}
	return nil
}

func (gtasks *Gtasks) DeleteTaskCache(task_id string, list_id string) {
	gtasks.Cache.EditedTasks.Deleted = append(gtasks.Cache.EditedTasks.Deleted, [...]string{list_id, task_id})
	gtasks.Cache.Tasks[list_id] = slices.DeleteFunc(gtasks.Cache.Tasks[list_id], func(task *tasks.Task) bool {
		return task.Id == task_id
	})
}

func (gtasks *Gtasks) AddTaskCache(task *tasks.Task, list_id string) {
	gtasks.Cache.EditedTasks.Created = append(gtasks.Cache.EditedTasks.Created, CreatedTask{Task: task, Id: list_id})
	if gtasks.Cache.Tasks == nil {
		gtasks.Cache.Tasks = make(map[string][]*tasks.Task, 1)
	}
	gtasks.Cache.Tasks[list_id] = append(gtasks.Cache.Tasks[list_id], task)
}

func (gtasks *Gtasks) DoTaskCache(list_id string , task_id string) error {
	gtasks.Cache.EditedTasks.Updated = append(gtasks.Cache.EditedTasks.Updated, [...]string{list_id, task_id})
	if gtasks.Cache.Tasks == nil {
		gtasks.Cache.Tasks = make(map[string][]*tasks.Task, 1)
	}
  tasks_tmp := gtasks.Cache.Tasks[list_id]
  index := slices.IndexFunc(tasks_tmp, func(task *tasks.Task) bool {
    return task.Id == task_id
  })
  if index == -1 {
    return errors.New("Task not found")
  }
  tasks_tmp[index].Status = "completed"
  return nil
}

func (gtasks *Gtasks) Sync(now time.Time) error {
	slog.Debug("Syncing tasks")
  slog.Debug("edited lists", "updated", gtasks.Cache.EditedLists.Updated, "created", gtasks.Cache.EditedLists.Created, "deleted", gtasks.Cache.EditedLists.Deleted)
  slog.Debug("edited tasks", "updated", gtasks.Cache.EditedTasks.Updated, "created", gtasks.Cache.EditedTasks.Created, "deleted", gtasks.Cache.EditedTasks.Deleted)
	for _, updated := range gtasks.Cache.EditedLists.Updated {
		_, err := gtasks.srv.Tasklists.Patch(updated, gtasks.Cache.Tasklists[updated]).Do()
		if err != nil {
			return err
		}
	}
	for _, created := range gtasks.Cache.EditedLists.Created {
		_, err := gtasks.srv.Tasklists.Insert(created).Do()
		if err != nil {
			return err
		}
	}
	for _, tasklist_id := range gtasks.Cache.EditedLists.Deleted {
		err := gtasks.srv.Tasklists.Delete(tasklist_id).Do()
		if err != nil {
			return err
		}
	}
	for _, updated := range gtasks.Cache.EditedTasks.Updated {
		tasks_tmp := gtasks.Cache.Tasks[updated[0]]
		index := slices.IndexFunc(tasks_tmp, func(task *tasks.Task) bool {
			return task.Id == updated[1]
		})
		if index == -1 {
			return errors.New("Updated task doesn't exist")
		}
		_, err := gtasks.srv.Tasks.Patch(updated[0], updated[1], tasks_tmp[index]).Do()
		if err != nil {
			return err
		}
	}
	for _, created := range gtasks.Cache.EditedTasks.Created {
		_, err := gtasks.srv.Tasks.Insert(created.Id, created.Task).Do()
		if err != nil {
			return err
		}
	}
	for _, deleted := range gtasks.Cache.EditedTasks.Deleted {
		err := gtasks.srv.Tasks.Delete(deleted[0], deleted[1]).Do()
		if err != nil {
			return err
		}
	}
	gtasks.Cache.EditedLists = EditedIdCache[string, *tasks.TaskList]{}
	gtasks.Cache.EditedTasks = EditedIdCache[[2]string, CreatedTask]{}
	gtasks.UpdateTasklists()
	gtasks.UpdateAllTasks(now)
	return nil
}

var default_credentials = `{
  "installed": {
    "client_id": "1076675530781-6bh6k30efrf796adk3s9t6ks5eee53j4.apps.googleusercontent.com",
    "project_id": "ydo-open-source",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "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",
    "redirect_uris": [
      "http://localhost:8183"
    ]
  }
}`

func (gtasks *Gtasks) ReadCredentials() (*oauth2.Config, error) {
	content, err := os.ReadFile(gtasks.CredentialsFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, CredentialsError{}
	}
	if len(content) == 0 {
		slog.Info("Reading the default credentials")
		content = []byte(default_credentials)
	}
	config, err := google.ConfigFromJSON(content, tasks.TasksScope)
	if err != nil {
		return nil, CredentialsError{}
	}
	return config, nil
}

type CredentialsError struct{}

func (err CredentialsError) Error() string {
	return "Unable to parse client secret file"
}

func (gtasks *Gtasks) getClient(c *oauth2.Config) *http.Client {
	tok, err := gtasks.tokenFromFile()
	if err != nil {
		tok, err = tokenFromWeb(c)
		gtasks.saveToken(tok)
	}
	return c.Client(context.Background(), tok)
}

func (gtasks *Gtasks) tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(gtasks.TokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (gtasks *Gtasks) saveToken(token *oauth2.Token) error {
	f, err := os.OpenFile(gtasks.TokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.New("Unable to cache oauth token")
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

var authCode string

func RedirectCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authCode = r.URL.Query().Get("code")
	if authCode != "" {
		w.Write([]byte("Got the code!"))
	} else {
		w.Write([]byte("Unable reading the code. You will get an error"))
	}
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func tokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	address := "localhost:8183"
	config.RedirectURL = "http://" + address
	router := httprouter.New()
	router.GET("/", RedirectCode)
	server := &http.Server{Handler: router, Addr: address}
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Go to the following link in your browser (if its not already opened):\n" + authURL)
	utils.OpenURL(authURL)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		ctx, release := context.WithTimeout(context.Background(), time.Second)
		defer release()
		server.Shutdown(ctx)
	}()
	server.ListenAndServe()
	fmt.Println("Caught the code, here it is: ", authCode)
	if authCode == "" {
		return nil, errors.New("Unable to retrieve auth code")
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// update tasklists and sort them by their due date
func (gtasks *Gtasks) UpdateTasklists() error {
	slog.Debug("retrieving task lists")
	r, err := gtasks.srv.Tasklists.List().Do()
	if err != nil {
		return errors.New("Unable to retrieve task lists\n" + err.Error())
	}
	gtasks.makeTasklistsMap(r.Items)
	return nil
}

func (gtasks *Gtasks) makeTasklistsMap(items []*tasks.TaskList) {
	if gtasks.Cache.Tasklists == nil {
		gtasks.Cache.Tasklists = make(map[string]*tasks.TaskList, len(items))
	}
	for _, item := range items {
		gtasks.Cache.Tasklists[item.Id] = item
	}
}

// update and sort tasks
func (gtasks *Gtasks) UpdateTasks(id string) error {
	if gtasks.Cache.Tasks == nil {
		gtasks.Cache.Tasks = make(map[string][]*tasks.Task, 1)
	}
	r, err := gtasks.srv.Tasks.List(id).ShowHidden(true).Do()
	if err != nil {
		return nil
	}
	gtasks.Cache.Tasks[id] = r.Items
	slices.SortStableFunc(gtasks.Cache.Tasks[id], func(a *tasks.Task, b *tasks.Task) int {
		a_due, _ := time.Parse(time.RFC3339, a.Due)
		b_due, _ := time.Parse(time.RFC3339, b.Due)
		a_updated, _ := time.Parse(time.RFC3339, a.Updated)
		b_updated, _ := time.Parse(time.RFC3339, b.Updated)
		return a_due.Compare(b_due)*2 + a_updated.Compare(b_updated)
	})
	return nil
}

func (gtasks *Gtasks) Load(now time.Time) error {
	err := gtasks.ReadCache(now)
	if err != nil {
		slog.Debug("Reading cache file failed.", "err", err)
		// if updating gtasks failed (for some connection problem) use the cache with
		// no errors returned
		gtasks.UpdateTasklists()

		if err := gtasks.UpdateAllTasks(now); err == nil {
			return nil
		}
	}
	return nil
}

func (gtasks *Gtasks) UpdateAllTasks(now time.Time) error {
	slog.Debug("Updating all tasklists data")
	for _, tasklist := range gtasks.Cache.Tasklists {
		err := gtasks.UpdateTasks(tasklist.Id)
		if err != nil {
			return err
		}
	}
	gtasks.SaveCache(now)
	return nil
}

type GtasksFilter func(*tasks.Task) bool
func (gtasks *Gtasks) PrintMarkdown(mc *core.MarkdownConfig, filter GtasksFilter) error {
	for _, item := range gtasks.Cache.Tasklists {
		mc.PrintListPrefix()
		fmt.Println(item.Title)
		tasks := gtasks.Cache.Tasks[item.Id]
		parent_stacks := map[string][]int{}
		for i, task := range tasks {
			if task.Parent != "" {
				parent_stacks[task.Parent] = append(parent_stacks[task.Parent], i)
			}
		}
		for _, task := range tasks {
			if task.Parent == "" {
				printTask(task, mc, 1, parent_stacks, tasks)
			}
		}
	}
	return nil
}

func IsCompleted(task *tasks.Task) bool {
  return task.Status == "completed"
}

func IsPending(task *tasks.Task) bool {
  return task.Status == "needsAction"
}

func printTask(task *tasks.Task, mc *core.MarkdownConfig, depth uint, parent_stacks map[string][]int, tasks []*tasks.Task) {
	mc.PrintIndent(depth)
	if task.Status == "completed" {
		mc.PrintDonePrefix()
	} else {
		mc.PrintUndonePrefix()
	}

	fmt.Println(task.Title)
	for line := range strings.Lines(task.Notes) {
		mc.PrintIndent(depth + 1)
		fmt.Println(line)
	}
	stack := parent_stacks[task.Id]
	if len(stack) > 0 {
		for _, index := range stack {
			printTask(tasks[index], mc, depth+1, parent_stacks, tasks)
		}
	}
}
