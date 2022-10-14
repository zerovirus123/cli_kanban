package main

import "cli_kanban/typedef"

type Task struct { // implementing the list.item inteface
	status      typedef.Status
	title       string
	description string
}

func NewTask(status typedef.Status, title, description string) Task {
	return Task{title: title, description: description, status: status}
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

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}
