package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/guionardo/todo-cli/internal"
	"github.com/urfave/cli/v2"
)

var (
	AddTodoTitle   string
	AddTodoDueDate cli.Timestamp
	AddTodoTags    = cli.NewStringSlice()
	AddCommand     = &cli.Command{
		Name:    "add",
		Usage:   "Add a new todo item",
		Aliases: []string{"a"},
		Action:  ActionAdd,

		// SkipFlagParsing: true,
		Flags: []cli.Flag{
			&cli.TimestampFlag{
				Name:        "due-date",
				Aliases:     []string{"d"},
				Usage:       "Due date for the todo item",
				Layout:      "2006-01-02",
				Required:    false,
				Destination: &AddTodoDueDate,
			},
			&cli.StringSliceFlag{
				Name:        "tags",
				Aliases:     []string{"t"},
				Usage:       "Tags for the todo item",
				Required:    false,
				Destination: AddTodoTags,
			},
		},
	}
)

func ActionAdd(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("Missing title")
	}
	AddTodoTitle = c.Args().Get(0)

	collection := GetCollection(c)

	var dueDate time.Time
	if AddTodoDueDate.Value() != nil {
		d := AddTodoDueDate.Value()
		dueDate = *d
	}

	item := &internal.ToDoItem{
		Title: AddTodoTitle,
		DueTo: dueDate,
		Tags:  AddTodoTags.Value(),
	}
	// TODO: DueTo and Tags are not being saved
	collection.Add(item)
	err := collection.Save(collection.FileName)
	if err == nil {
		log.Printf("Add todo: %s", item.ToMarkDown())
	} else {
		log.Printf("Error saving collection: %s", err)
	}
	return nil
}
