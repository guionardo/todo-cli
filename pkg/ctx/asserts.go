package ctx

import (
	"context"
	"time"

	"github.com/guionardo/todo-cli/pkg/exceptions"
	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

// RequiredTodoId Asserts that the ID is valid and exists in the collection
func RequiredTodoId(c *cli.Context) error {
	c2 := ContextFromCtx(c)
	if c2.Id < 1 {
		return exceptions.InvalidTodoIdError(c2.Id)
	}
	c2.CurrentToDo = c2.Collection.Get(c2.Id)
	if c2.CurrentToDo == nil {
		return exceptions.TodoNotFoundError(c2.Id)
	}
	return nil
}

// OptionalId Asserts that the ID is valid and exists in the collection
func OptionalId(c *cli.Context) error {
	if !c.IsSet("id") {
		return nil
	}
	c2 := ContextFromCtx(c)
	if c2.Id < 1 {
		return exceptions.InvalidTodoIdError(c2.Id)
	}
	c2.CurrentToDo = c2.Collection.Get(c2.Id)
	if c2.CurrentToDo == nil {
		return exceptions.TodoNotFoundError(c2.Id)
	}
	return nil
}

func AssertSychronization(c *cli.Context) error {
	c2 := ContextFromCtx(c)

	if err, ok := c.Err().(*exceptions.TodoException); ok && err.CancelSync {
		logger.Debugf("Canceling sync due to error: %s", err.Error())
		return nil
	}
	if c2.LocalConfig.Gist.Authorization == "" {
		logger.Debugf("Canceling sync due to no authorization")
		return nil
	}

	logs, err := c2.GistSync()
	if err != nil {
		logger.Warnf("Error syncing GIST: %s", err)
		return nil
	} else {
		for _, log := range logs {
			logger.Infof(" %s", log)
		}
	}

	return nil
}

// LocalConfigRequired Assert configuration exists and is valid
func LocalConfigRequired(c *cli.Context) (err error) {
	c2 := ContextFromCli(c)
	c.Context = context.WithValue(c.Context, Key, c2)
	err = c2.Error
	if err != nil {
		err = exceptions.NoSetupError(err)
	}

	return
}

func AssertAutoSychronization(c *cli.Context) error {
	c2 := ContextFromCtx(c)
	passedTimeForSync := time.Now().Sub(c2.LocalConfig.Gist.LastSync) - c2.LocalConfig.Gist.AutoSyncInterval
	if passedTimeForSync <= 0 {
		return nil
	}
	if !c2.LocalConfig.Gist.AutoSync {
		logger.Warnf("AutoSync is disabled. Passed %v for sync", passedTimeForSync)
		return nil
	}
	return AssertSychronization(c)
}

func AssertSave(c *cli.Context) (err error) {
	c2 := ContextFromCtx(c)
	if err, ok := c.Err().(*exceptions.TodoException); ok && err.CancelSaving {
		logger.Debugf("Canceling saving due to error: %s", err.Error())
		return nil
	}	
	if err = c2.Collection.Save(c2.LocalCollectionFile); err == nil {
		err = c2.LocalConfig.Save(c2.LocalConfigFile)
	}
	return
}
