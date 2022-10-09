package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* MODEL MANAGEMENT */
var models []tea.Model

/* MAIN MODEL */

type Model struct {
	quitting bool
	loaded   bool   // wait before anything is displayed, to make sure that the lists have been initialized
	focused  status // holds the list that is currently in focus
	lists    []list.Model
}

func New() *Model {
	return &Model{}
}

func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()

	if selectedItem != nil {
		selectedTask := selectedItem.(Task)
		m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
		selectedTask.Next() // increment the selectedTask.status field
		m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	}

	return nil
}

// TODO: Go to next list
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

// TODO: go to previous list
func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

// TODO: call this on tea.WindowSizeMsg
// on startup, grabs the size of the terminal window and adjust the list accordingly
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	columns := ReadFromStorage()

	// Init To Dos
	m.lists[todo].Title = "To Do"
	m.lists[inProgress].Title = "In Progress"
	m.lists[done].Title = "Done"

	todoItems := []list.Item{}
	inProgressItems := []list.Item{}
	doneItems := []list.Item{}

	for _, value := range columns.Todo {
		task := Task{status: todo, title: value.Title, description: value.Description}
		todoItems = append(todoItems, task)
	}

	for _, value := range columns.InProgress {
		task := Task{status: inProgress, title: value.Title, description: value.Description}
		inProgressItems = append(inProgressItems, task)
	}

	for _, value := range columns.Done {
		task := Task{status: done, title: value.Title, description: value.Description}
		doneItems = append(doneItems, task)
	}

	m.lists[todo].SetItems(todoItems)
	m.lists[inProgress].SetItems(inProgressItems)
	m.lists[done].SetItems(doneItems)

}

func (m Model) Init() tea.Cmd {
	return nil // no timer to startup when the program starts
}

// update the list (passing the interactions and keypresses, no logic involved)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // terminal dimensions on program startup
		if !m.loaded { // if list is not loaded, initialize it
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			WriteToStorage(m)
			return m, tea.Quit
		case "left", "a":
			m.Prev()
		case "right", "d":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		case "n": // saves the state of the current models to an array of models
			models[model] = m
			models[form] = NewForm(m.focused)
			return models[form].Update(nil)
		case "x": //deletes an entry
			index := m.lists[m.focused].Index()
			m.lists[m.focused].RemoveItem(index)
		}
	case Task:
		task := msg
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg) //update the list
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.loaded {
		todoView := m.lists[todo].View()
		inProgView := m.lists[inProgress].View()
		doneView := m.lists[done].View()

		switch m.focused {

		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		}
	} else {
		return "loading..."
	}
}
