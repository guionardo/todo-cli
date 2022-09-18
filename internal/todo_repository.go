package internal

type ToDoRepository interface {
	Add(item ToDoItem) error
	Complete(title string) error
}
