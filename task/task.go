package task

import (
	"cli_kanban/typedef"
)

type Task struct { // implementing the list.item inteface
	status      typedef.Status
	title       string
	description string
}

func NewTask(status typedef.Status, title string, description string) *Task {
	return &Task{status: status, title: title, description: description}
}

func (t *Task) Next() {
	if t.status == typedef.Done {
		t.status = typedef.Todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) GetTitle() string {
	return t.title
}

func (t Task) GetDescription() string {
	return t.description
}

func (t Task) GetStatus() typedef.Status {
	return t.status
}
