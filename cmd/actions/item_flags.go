package actions

import "github.com/urfave/cli/v2"

var ItemFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "title",
		Aliases:  []string{"t"},
		Usage:    "Title of the task",
		Required: true,
	},
	&cli.TimestampFlag{
		Name:     "due-date",
		Aliases:  []string{"d"},
		Usage:    "Due date for the todo item",
		Layout:   "2006-01-02",
		Required: false,
	},
	&cli.StringSliceFlag{
		Name:     "tags",
		Usage:    "Tags for the todo item",
		Required: false,
		Aliases:  []string{"t"},
	},
	&cli.IntFlag{
		Name:     "parent-id",
		Usage:    "Parent ID for the todo item",
		Required: false,
		Aliases:  []string{"p"},
	},
}
