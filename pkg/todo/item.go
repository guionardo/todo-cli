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
