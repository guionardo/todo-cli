package todo

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	Id         string    `yaml:"id"`
	Index      int       `yaml:"index"`
	Title      string    `yaml:"title"`
	Completed  bool      `yaml:"completed"`
	DueTo      time.Time `yaml:"due_to"`
	LastAction time.Time `yaml:"last_action"`
	Tags       []string  `yaml:"tags"`
	UpdatedAt  time.Time `yaml:"updated_at"`
	ParentId   string    `yaml:"parent_id"`
	Deleted    bool      `yaml:"-"`
	Level      int       `yaml:"level"`
}

func NewItemId() string {
	return uuid.New().String()
}

func (i *Item) Equal(other *Item) bool {
	return (i != nil &&
		other != nil &&
		i.Id == other.Id &&
		i.Title == other.Title &&
		i.Completed == other.Completed &&
		i.DueTo.Equal(other.DueTo) &&
		i.LastAction.Equal(other.LastAction) &&
		i.UpdatedAt.Equal(other.UpdatedAt) &&
		i.ParentId == other.ParentId &&
		i.Deleted == other.Deleted &&
		i.Level == other.Level)
}

func (i *Item) SetTags(tags []string) {
	tags = ParseTags(tags)
	if EqualTags(i.Tags, tags) {
		return
	}
	i.Tags = tags
	i.UpdatedAt = time.Now()
}
