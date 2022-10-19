package main

import (
	"cli_kanban/task"
	"cli_kanban/typedef"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* MODEL MANAGEMENT */
var models []tea.Model

/* MAIN MODEL */

type Model struct {
	quitting bool
	loaded   bool           // wait before anything is displayed, to make sure that the lists have been initialized
	focused  typedef.Status // holds the list that is currently in focus
	lists    []list.Model
}

func New() *Model {
	return &Model{}
}

func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()

	if selectedItem != nil {
		selectedTask := selectedItem.(*task.Task)
		m.lists[selectedTask.Status()].RemoveItem(m.lists[m.focused].Index())
		selectedTask.Next() // increment the selectedTask.status field
		m.lists[selectedTask.Status()].InsertItem(len(m.lists[selectedTask.Status()].Items())-1, list.Item(selectedTask))
	}

	return nil
}

// TODO: Go to next list
func (m *Model) Next() {
	if m.focused == typedef.Done {
		m.focused = typedef.Todo
	} else {
		m.focused++
	}
}

// TODO: go to previous list
func (m *Model) Prev() {
	if m.focused == typedef.Todo {
		m.focused = typedef.Done
	} else {
		m.focused--
	}
}

// TODO: call this on tea.WindowSizeMsg
// on startup, grabs the size of the terminal window and adjust the list accordingly
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/typedef.Divisor, height/2)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	columns := ReadFromStorage()

	// Init To Dos
	m.lists[typedef.Todo].Title = "To Do"
	m.lists[typedef.InProgress].Title = "In Progress"
	m.lists[typedef.Done].Title = "Done"

	todoItems := []list.Item{}
	inProgressItems := []list.Item{}
	doneItems := []list.Item{}

	for _, value := range columns.Todo {
		task := task.NewTask(typedef.Todo, value.Title, value.Description)
		todoItems = append(todoItems, task)
	}

	for _, value := range columns.InProgress {
		task := task.NewTask(typedef.InProgress, value.Title, value.Description)
		inProgressItems = append(inProgressItems, task)
	}

	for _, value := range columns.Done {
		task := task.NewTask(typedef.Done, value.Title, value.Description)
		doneItems = append(doneItems, task)
	}

	m.lists[typedef.Todo].SetItems(todoItems)
	m.lists[typedef.InProgress].SetItems(inProgressItems)
	m.lists[typedef.Done].SetItems(doneItems)
}

func (m Model) Init() tea.Cmd {
	return nil // no timer to startup when the program starts
}

// update the list (passing the interactions and keypresses, no logic involved)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // terminal dimensions on program startup
		if !m.loaded { // if list is not loaded, initialize it
			columnStyle.Width(msg.Width / typedef.Divisor)
			focusedStyle.Width(msg.Width / typedef.Divisor)
			columnStyle.Height(msg.Height - typedef.Divisor)
			focusedStyle.Height(msg.Height - typedef.Divisor)
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
			models[typedef.Model] = m
			models[typedef.Form] = NewForm(m.focused)
			return models[typedef.Form].Update(nil)
		case "x": // deletes an entry
			index := m.lists[m.focused].Index()
			m.lists[m.focused].RemoveItem(index)
		}
	case task.Task:
		task := msg
		return m, m.lists[task.Status()].InsertItem(len(m.lists[task.Status()].Items()), task)
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg) // update the list
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.loaded {
		todoView := m.lists[typedef.Todo].View()
		inProgView := m.lists[typedef.InProgress].View()
		doneView := m.lists[typedef.Done].View()

		switch m.focused {

		case typedef.InProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		case typedef.Done:
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
