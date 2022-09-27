package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/google/uuid"
)

type ToDoItem struct {
	Id         string    `yaml:"id"`
	Index      int       `yaml:"index"`
	Title      string    `yaml:"title"`
	Completed  bool      `yaml:"completed"`
	DueTo      time.Time `yaml:"due_to"`
	LastAction time.Time `yaml:"last_action"`
	Tags       []string  `yaml:"tags"`
	UpdatedAt  time.Time `yaml:"updated_at"`
}

const (
	TimeFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
)

func (item *ToDoItem) ItemColor() func(string) string {
	if item.Completed {
		return color.InGreen
	}
	if item.DueTo.IsZero() || item.DueTo.After(time.Now()) {
		return func(s string) string {
			return color.InBold(color.InYellow(s))
		}
	}
	if !item.DueTo.IsZero() && item.DueTo.Before(time.Now()) {
		return func(s string) string {
			return color.InBold(color.InRed(s))
		}
	}
	return color.InWhite
}

func (item *ToDoItem) String() string {
	var clr = item.ItemColor()

	return clr(item.StringNoColor())
}

func (item *ToDoItem) StringNoColor() string {
	completed := " "
	if item.Completed {
		completed = "✓"
	}
	tags := ""
	if len(item.Tags) > 0 {
		tags = fmt.Sprintf("(%s)", strings.Join(item.Tags, " "))
	}
	lastAction := ""
	if !item.LastAction.IsZero() {
		lastAction = fmt.Sprintf("Last action: %s ", item.LastAction.Format(DateTimeFormat))
	}
	dueTo := ""
	if !item.DueTo.IsZero() {
		dueTo = fmt.Sprintf("(%s) ", item.DueTo.Format(TimeFormat))
	}
	return fmt.Sprintf("#%03d %s %s %s %s %s", item.Index, completed, tags, dueTo, lastAction, item.Title)
}

func NewItemId() string {
	return uuid.New().String()
}

func (item *ToDoItem) NotifyText() string {
	if item.Completed {
		return fmt.Sprintf("Completed @ %s", item.UpdatedAt.Format(DateTimeFormat))
	}
	if !item.DueTo.IsZero() {
		if item.DueTo.Before(time.Now()) {
			return fmt.Sprintf("Overdue %d days", int(item.DueTo.Sub(time.Now()).Hours()/24))
		}
		daysToComplete := time.Now().Sub(item.DueTo).Hours() / 24
		return fmt.Sprintf("Due in %d days @ %s", int(daysToComplete), item.DueTo.Format(DateTimeFormat))
	}
	lastAction := item.LastAction
	if lastAction.IsZero() {
		lastAction = item.UpdatedAt
	}
	if lastAction.IsZero() {
		lastAction = time.Now()
	}
	days := time.Now().Sub(lastAction).Hours() / 24

	return fmt.Sprintf("New (%d days)", int(days))
}
