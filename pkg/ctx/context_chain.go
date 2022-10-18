package ctx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/guionardo/todo-cli/pkg/logger"
	"github.com/urfave/cli/v2"
)

func ChainedContext(funcs ...func(c *cli.Context) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		for _, f := range funcs {
			err := f(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Assert configuration exists and is valid
func AssertLocalConfig(c *cli.Context) (err error) {
	c2 := ContextFromCli(c)
	c.Context = context.WithValue(c.Context, "context", c2)
	err = c2.Error
	if err != nil {
		err = errors.New("Error loading configuration - Run setup again: " + err.Error())
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
	if c2.CancelSaving {
		return
	}
	if err = c2.Collection.Save(c2.LocalCollectionFile); err == nil {
		err = c2.LocalConfig.Save(c2.LocalConfigFile)
	}
	return
}

// AssertValidId checks if
func AssertValidId(c *cli.Context) error {
	c2 := ContextFromCtx(c)
	if c2.Id < 1 {
		return errors.New("Invalid ID")
	}
	c2.CurrentToDo = c2.Collection.Get(c2.Id)
	if c2.CurrentToDo == nil {
		return fmt.Errorf("To-do #%d not found", c2.Id)
	}
	return nil
}

func OptionalId(c *cli.Context) error {
	c2 := ContextFromCtx(c)
	if c2.Id < 1 {
		return nil
	}
	c2.CurrentToDo = c2.Collection.Get(c2.Id)
	if c2.CurrentToDo == nil {
		return fmt.Errorf("To-do #%d not found", c2.Id)
	}
	return nil
}

func AssertSychronization(c *cli.Context) error {
	c2 := ContextFromCtx(c)
	if c2.CancelSync || c2.LocalConfig.Gist.Authorization == "" {
		return nil
	}

	diffCount, log, err := c2.Collection.GistSync(&c2.LocalConfig.Gist)
	if err != nil {
		logger.Warnf("Error syncing GIST: %s", err)
		return nil
	}
	if diffCount == 0 {
		return nil
	}

	if err = c2.Collection.Save(c2.LocalCollectionFile); err != nil {
		logger.Warnf("Error saving collection - syncing failed: %s", err)
		return nil
	}
	logger.Infof("Synced %d changes from GIST", diffCount)
	for _, l := range log {
		logger.Infof(" %s", l)
	}
	c2.LocalConfig.Save(c2.LocalConfigFile)

	return nil
}
