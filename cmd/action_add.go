package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/guionardo/todo-cli/internal"
	"github.com/urfave/cli/v2"
)

var (
	AddCommand = &cli.Command{
		Name:     "add",
		Usage:    "Add a new todo item",
		Aliases:  []string{"a"},
		Action:   ActionAdd,
		Category: "Tasks",
		// SkipFlagParsing: true,
		Flags: []cli.Flag{
			&cli.TimestampFlag{
				Name:     "due-date",
				Aliases:  []string{"d"},
				Usage:    "Due date for the todo item",
				Layout:   "2006-01-02",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "tags",
				Aliases:  []string{"t"},
				Usage:    "Tags for the todo item",
				Required: false,
			},
		},
	}
)

func ActionAdd(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("Missing title")
	}
	title := c.Args().Get(0)

	context := internal.GetRunningContext(c).AssertExist()

	var dueDate = c.Timestamp("due-date")

	tags := c.StringSlice("tags")
	item := &internal.ToDoItem{
		Id:        internal.NewItemId(),
		Title:     title,
		DueTo:     *dueDate,
		Tags:      tags,
		UpdatedAt: time.Now(),
	}

	context.Collection.Add(item)
	err := context.Collection.Save(context.CollectionFileName)
	if err == nil {
		log.Printf("Add todo: %s", item.ToMarkDown())
		context.Collection.GISTSync(context.DebugMode)
	} else {
		log.Printf("Error saving collection: %s", err)
	}
	return nil
}
