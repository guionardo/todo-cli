package ctx

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/guionardo/todo-cli/pkg/backup"
	"github.com/guionardo/todo-cli/pkg/consts"
	"github.com/guionardo/todo-cli/pkg/github"
	"github.com/guionardo/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)



type Context struct {
	DataFolder          string
	LocalConfigFile     string
	LocalCollectionFile string
	LocalConfig         *LocalConfig
	Collection          *todo.ToDoCollection
	Error               error
	ExitMessage         string
	Id                  int // Id is the first argument that is a number
	CurrentToDo         *todo.ToDoItem
	Args                []string // Args is the list of arguments that are not numbers
	CancelSaving        bool
	CancelSync          bool
}

func GetDefaultDataFolder() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("Error getting user home dir: %s", err))
	}
	return path.Join(home, ".config", "todo-cli")
}

func ContextFromCli(c *cli.Context) *Context {
	dataFolder := c.String("data-folder")
	if dataFolder == "" {
		dataFolder = GetDefaultDataFolder()
	}
	c2 := ContextFromDataFolder(dataFolder)
	c2.Args = make([]string, 0)
	for i := 0; i < c.NArg(); i++ {
		id, err := strconv.Atoi(c.Args().Get(i))
		if err == nil {
			c2.Id = id
		} else {
			c2.Args = append(c2.Args, c.Args().Get(i))
		}
	}
	if id, err := strconv.Atoi(c.Args().Get(0)); err == nil {
		c2.Id = id
	}
	if c.IsSet("id") {
		if id := c.Int("id"); id > 0 {
			c2.Id = id
		}
	}
	return c2
}

func ContextFromCtx(c *cli.Context) *Context {
	return c.Context.Value("ctx").(*Context)
}

func ContextFromDataFolder(dataFolder string) *Context {
	stat, err := os.Stat(dataFolder)
	if err == nil && !stat.IsDir() {
		panic(fmt.Errorf("Data folder %s is not a directory", dataFolder))
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(dataFolder, 0755)
	}
	if err != nil {
		panic(fmt.Errorf("Error creating config dir %s: %s", dataFolder, err))
	}

	context := &Context{
		DataFolder:          dataFolder,
		LocalConfigFile:     path.Join(dataFolder, consts.DefaultLocalConfigFile),
		LocalCollectionFile: path.Join(dataFolder, consts.DefaultLocalCollectionFile),
	}

	config, err := LoadLocalConfig(context.LocalConfigFile)
	if err != nil {
		context.LocalConfig = &LocalConfig{
			ToDoListName: "todo",
			Gist:         github.GetDefaultGistConfig(),
			Backup:       backup.GetDefaultBackupConfig(dataFolder),
		}
		context.Error = err
	}
	context.LocalConfig = &config

	collection, err := todo.LoadCollection(context.LocalCollectionFile)
	if err != nil {
		context.Error = err
		return context
	}
	context.Collection = &collection
	return context
}

func (c *Context) SetExit(err error, format string, args ...interface{}) error {
	if err != nil {
		c.ExitMessage = err.Error()
	} else {
		c.ExitMessage = fmt.Sprintf(format, args...)
	}
	return err
}
