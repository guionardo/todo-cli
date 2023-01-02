package ctx

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"time"

	"github.com/guionardo/todo-cli/pkg/logger"
)

func ChainedContext(functions ...func(c *cli.Context) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		for _, f := range functions {
			err := f(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// AssertLocalConfig Assert configuration exists and is valid
func AssertLocalConfig(c *cli.Context) (err error) {
	c2 := ContextFromCli(c)
	c.Context = context.WithValue(c.Context, Key, c2)
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
		return errors.New("invalid ID")
	}
	c2.CurrentToDo = c2.Collection.Get(c2.Id)
	if c2.CurrentToDo == nil {
		return fmt.Errorf("to-do #%d not found", c2.Id)
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
		return fmt.Errorf("to-do #%d not found", c2.Id)
	}
	return nil
}

func AssertSychronization(c *cli.Context) error {
	c2 := ContextFromCtx(c)
	if c2.CancelSync || c2.LocalConfig.Gist.Authorization == "" {
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
