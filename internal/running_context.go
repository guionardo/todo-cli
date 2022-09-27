package internal

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/guionardo/todo-cli/pkg/github"
	"github.com/urfave/cli/v2"
)

type (
	RunningContext struct {
		CollectionFileName string
		DebugMode          bool
		Collection         *ToDoCollection
		GistAPI            *github.GitHubGistAPI
	}

	Config struct {
		ToDoListName   string `yaml:"todo_list_name"`
		GistId         string `yaml:"gist_id"`
		Authorization  string `yaml:"authorization"`
		ConfigFileName string `yaml:"-"`
	}
)

const (
	CollectionFileName = "todo-cli.yaml"
)

var (
	DefaultCollectionFilePath string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home dir: %s", err)
	}
	configPath := path.Join(home, ".config", "todo-cli")
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(configPath, 0755)
	}
	if err != nil {
		log.Fatalf("Error creating config dir %s: %s", configPath, err)
	}
	DefaultCollectionFilePath = path.Join(configPath, CollectionFileName)
}

func (c Config) String() string {
	return fmt.Sprintf("ToDoListName: %s\nGistId: %s\nAuthorization: %s\nFile: %s\n", c.ToDoListName, c.GistId, c.Authorization, c.ConfigFileName)
}

func NewRunningContext(c *cli.Context) *RunningContext {
	debugMode := c.Bool("debug")
	collectionFileName := c.String("config")
	if collectionFileName == "" {
		collectionFileName = DefaultCollectionFilePath
		if debugMode {
			log.Printf("Using default collection file %s", collectionFileName)
		}
	}

	return &RunningContext{
		CollectionFileName: collectionFileName,
		DebugMode:          debugMode,
	}
}

func GetRunningContext(c *cli.Context) *RunningContext {
	ctx := c.Context.Value("running_context")
	if ctx == nil {
		log.Fatalf("Running context not found")
	}
	return ctx.(*RunningContext)
}

func CollectionFile() (string, error) {
	if todo_collection := os.Getenv("TODO_COLLECTION"); todo_collection != "" {
		if _, err := os.Stat(todo_collection); err == nil {
			_, err := ParseCollectionFile(todo_collection)
			if err == nil {
				return todo_collection, nil
			}
			log.Printf("Error parsing collection file %s: %s", todo_collection, err)
		}
	}
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := path.Join(home, ".config", "todo-cli")
		if _, err = os.Stat(configPath); os.IsNotExist(err) {
			err = os.MkdirAll(configPath, 0755)
		}
		if err == nil {
			return path.Join(configPath, "todo.yaml"), nil
		}
	}
	return "", err
}

func (c *RunningContext) AssertExist() *RunningContext {
	if stat, err := os.Stat(c.CollectionFileName); err != nil || stat.IsDir() {
		log.Fatalf("Collection file %s not found - Run setup", c.CollectionFileName)
	}
	collection, err := ParseCollectionFile(c.CollectionFileName)
	if err != nil {
		log.Fatalf("Error parsing collection file: %s", err)
	}
	c.Collection = collection
	if c.DebugMode {
		log.Printf("Collection file %s loaded", c.CollectionFileName)
	}
	return c
}

func (c *RunningContext) GetCollection() *ToDoCollection {
	if c.Collection == nil {
		collection, err := ParseCollectionFile(c.CollectionFileName)
		if err != nil {
			log.Fatalf("Error parsing collection file: %s", err)
		}
		c.Collection = collection
	}
	return c.Collection
}

func (c *RunningContext) GetGistAPI() *github.GitHubGistAPI {
	if c.GistAPI == nil {
		auth := c.GetCollection().Config.Authorization
		c.GistAPI = github.NewGitHubGistAPI(auth, c.DebugMode)
	}
	return c.GistAPI
}

func (c *RunningContext) DebugLog(format string, args ...interface{}) {
	if c.DebugMode {
		log.Printf(format, args...)
	}
}
