package ctx

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/guionardo/todo-cli/pkg/config"
	"github.com/guionardo/todo-cli/pkg/consts"
	"github.com/guionardo/todo-cli/pkg/todo"
	"github.com/guionardo/todo-cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

type (
	KeyType struct{}
	Context struct {
		DataFolder          string
		LocalConfigFile     string
		LocalCollectionFile string
		LocalConfig         *config.LocalConfig
		Collection          *todo.Collection
		Error               error
		ExitMessage         string
		ExitWarning         bool
		Id                  int // Id is the first argument that is a number
		CurrentToDo         *todo.Item
		Args                []string // Args is the list of arguments that are not numbers
		ErrorPrefix         string
	}
)

var Key = KeyType{}

var globalContext Context

func ContextFromCli(c *cli.Context) *Context {
	if len(globalContext.DataFolder) == 0 {
		dataFolder := c.String("data-folder")
		if dataFolder == "" {
			dataFolder = utils.GetDefaultDataFolder()
		}
		globalContext = contextFromDataFolder(dataFolder)
		globalContext.Args = make([]string, 0)
		for i := 0; i < c.NArg(); i++ {
			id, err := strconv.Atoi(c.Args().Get(i))
			if err == nil {
				globalContext.Id = id
			} else {
				globalContext.Args = append(globalContext.Args, c.Args().Get(i))
			}
		}
		if id, err := strconv.Atoi(c.Args().Get(0)); err == nil {
			globalContext.Id = id
		}
		if c.IsSet("id") {
			if id := c.Int("id"); id > 0 {
				globalContext.Id = id
			}
		}
	}
	return &globalContext
}

func ContextFromCtx(c *cli.Context) *Context {
	return c.Context.Value(Key).(*Context)
}

func contextFromDataFolder(dataFolder string) Context {
	stat, err := os.Stat(dataFolder)
	if err == nil && !stat.IsDir() {
		panic(fmt.Errorf("data folder %s is not a directory", dataFolder))
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(dataFolder, 0755)
	}
	if err != nil {
		panic(fmt.Errorf("error creating config dir %s: %s", dataFolder, err))
	}

	context := Context{
		DataFolder:          dataFolder,
		LocalConfigFile:     path.Join(dataFolder, consts.DefaultLocalConfigFile),
		LocalCollectionFile: path.Join(dataFolder, consts.DefaultLocalCollectionFile),
	}

	cfg, err := config.LoadLocalConfig(context.LocalConfigFile)
	if err != nil {
		context.LocalConfig = &config.LocalConfig{
			ToDoListName: "todo",
			Gist:         config.GetDefaultGistConfig(),
			Backup:       config.GetDefaultBackupConfig(dataFolder),
		}
		context.Error = err
	}
	context.LocalConfig = &cfg

	collection, err := todo.LoadCollection(context.LocalCollectionFile)
	if err != nil {
		context.Error = err
		return context
	}
	context.Collection = &collection
	return context
}

// SetExitError sets the error and the exit message.
func (ctx *Context) SetExitError(err error, format string, args ...interface{}) error {
	ctx.Error = err
	ctx.ExitMessage = fmt.Sprintf(format, args...)
	return err
}

func (ctx *Context) SetExitSuccess(format string, args ...interface{}) error {
	ctx.ExitMessage = fmt.Sprintf(format, args...)
	return nil
}

func (ctx *Context) SetExitWarning(format string, args ...interface{}) error {
	ctx.ExitMessage = fmt.Sprintf(format, args...)
	ctx.ExitWarning = true
	return nil
}
