package exceptions

import (
	"errors"
	"fmt"
)

type TodoException struct {
	Err          error
	CancelSaving bool
	CancelSync   bool
}

func NewException(err error) *TodoException {
	return &TodoException{Err: err,
		CancelSaving: true,
		CancelSync:   true,
	}
}

func (e *TodoException) Error() string {
	return e.Err.Error()
}

func TodoNotFoundError(id int) error {
	return fmt.Errorf("to-do #%d not found", id)
}

func ParentTodoNotFoundError(id int) error {
	return fmt.Errorf("parent to-do #%d not found", id)
}

func InvalidTodoIdError(id int) error {
	return fmt.Errorf("invalid ID %d", id)
}

func NoSetupError(err error) error {
	return fmt.Errorf("no setup: %s - run setup again", err)
}

func NoChangedTodoError(index int) error {
	return NewException(fmt.Errorf("to-do item #%d no change", index))
}

func TodoItemHasChildrenError(index int, children int) error {
	return fmt.Errorf("update failed. To-do item #%d has %d children", index, children)
}

var (
	MissingTodoIdError = errors.New("missing ID")
	NoNeedBackupError  = errors.New("current collection is the same as last backup. Skipping backup")
)
