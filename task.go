package main

type Task struct { // implementing the list.item inteface
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) Task {
	return Task{title: title, description: description, status: status}
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
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
