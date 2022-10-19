package model

import (
	"cli_kanban/datastore"
	"cli_kanban/form"
	"cli_kanban/styling"
	"cli_kanban/task"
	"cli_kanban/typedef"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /* MODEL MANAGEMENT */
// var Models []tea.Model

type ModelHandler typedef.Model

func New() *ModelHandler {
	return &ModelHandler{}
}

func (m *ModelHandler) MoveToNext() tea.Msg {
	selectedItem := m.Lists[m.Focused].SelectedItem()

	if selectedItem != nil {
		selectedTask := selectedItem.(*task.Task)
		m.Lists[selectedTask.Status()].RemoveItem(m.Lists[m.Focused].Index())
		selectedTask.Next() // increment the selectedTask.status field
		m.Lists[selectedTask.Status()].InsertItem(len(m.Lists[selectedTask.Status()].Items())-1, list.Item(selectedTask))
	}

	return nil
}

// TODO: Go to next list
func (m *ModelHandler) Next() {
	if m.Focused == typedef.Done {
		m.Focused = typedef.Todo
	} else {
		m.Focused++
	}
}

// TODO: go to previous list
func (m *ModelHandler) Prev() {
	if m.Focused == typedef.Todo {
		m.Focused = typedef.Done
	} else {
		m.Focused--
	}
}

// TODO: call this on tea.WindowSizeMsg
// on startup, grabs the size of the terminal window and adjust the list accordingly
func (m *ModelHandler) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/typedef.Divisor, height/2)
	defaultList.SetShowHelp(false)
	m.Lists = []list.Model{defaultList, defaultList, defaultList}

	columns := datastore.ReadFromStorage()

	// Init To Dos
	m.Lists[typedef.Todo].Title = "To Do"
	m.Lists[typedef.InProgress].Title = "In Progress"
	m.Lists[typedef.Done].Title = "Done"

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

	m.Lists[typedef.Todo].SetItems(todoItems)
	m.Lists[typedef.InProgress].SetItems(inProgressItems)
	m.Lists[typedef.Done].SetItems(doneItems)
}

func (m ModelHandler) Init() tea.Cmd {
	return nil // no timer to startup when the program starts
}

// update the list (passing the interactions and keypresses, no logic involved)
func (m ModelHandler) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // terminal dimensions on program startup
		if !m.Loaded { // if list is not loaded, initialize it
			styling.ColumnStyle.Width(msg.Width / typedef.Divisor)
			styling.FocusedStyle.Width(msg.Width / typedef.Divisor)
			styling.ColumnStyle.Height(msg.Height - typedef.Divisor)
			styling.FocusedStyle.Height(msg.Height - typedef.Divisor)
			m.initLists(msg.Width, msg.Height)
			m.Loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			datastore.WriteToStorage((typedef.Model)(m))
			return m, tea.Quit
		case "left", "a":
			m.Prev()
		case "right", "d":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		case "n": // saves the state of the current models to an array of models
			typedef.Models[typedef.ModelEnum] = m
			typedef.Models[typedef.FormEnum] = form.NewForm(m.Focused)
			return typedef.Models[typedef.FormEnum].Update(nil)
		case "x": // deletes an entry
			index := m.Lists[m.Focused].Index()
			m.Lists[m.Focused].RemoveItem(index)
		}
	case task.Task:
		task := msg
		return m, m.Lists[task.Status()].InsertItem(len(m.Lists[task.Status()].Items()), task)
	}

	var cmd tea.Cmd
	m.Lists[m.Focused], cmd = m.Lists[m.Focused].Update(msg) // update the list
	return m, cmd
}

func (m ModelHandler) View() string {
	if m.Quitting {
		return ""
	}

	if m.Loaded {
		todoView := m.Lists[typedef.Todo].View()
		inProgView := m.Lists[typedef.InProgress].View()
		doneView := m.Lists[typedef.Done].View()

		switch m.Focused {

		case typedef.InProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				styling.ColumnStyle.Render(todoView),
				styling.FocusedStyle.Render(inProgView),
				styling.ColumnStyle.Render(doneView),
			)
		case typedef.Done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				styling.ColumnStyle.Render(todoView),
				styling.ColumnStyle.Render(inProgView),
				styling.FocusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				styling.FocusedStyle.Render(todoView),
				styling.ColumnStyle.Render(inProgView),
				styling.ColumnStyle.Render(doneView),
			)
		}
	} else {
		return "loading..."
	}
}
