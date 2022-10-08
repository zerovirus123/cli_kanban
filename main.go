package main // mandatory import

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4 // divide the total width of the viewport by 4

const ( //indices to determine which list is focused
	todo status = iota
	inProgress
	done
)

/* MODEL MANAGEMENT */
var models []tea.Model

const (
	model status = iota
	form
)

/* STYLING  */
var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())

	focusedStyle = lipgloss.NewStyle(). // styling for the focused column
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	helpStyle = lipgloss.NewStyle(). // styling for the help text on the bottom
			Foreground(lipgloss.Color("241"))
)

/* CUSTOM ITEM */

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

/* MAIN MODEL */

type Model struct {
	quitting bool
	loaded   bool   // wait before anything is displayed, to make sure that the lists have been initialized
	focused  status // holds the list that is currently in focus
	lists    []list.Model
	err      error // display error to user at some point
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

	// Init To Dos
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t shirts"},
	})

	// Init in progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "write code", description: "don't worry, it's Go"},
	})

	// Init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "stay cool", description: "as a cucumber"},
	})
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

/* FORM MODEL */
type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
}

func (m Form) CreateTask() tea.Msg {
	task := NewTask(m.focused, m.title.Value(), m.description.Value())
	return task
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				models[form] = m
				return models[model], m.CreateTask
			}
		}
	}

	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	}
}

func NewForm(focused status) *Form {
	form := &Form{focused: focused}
	form.title = textinput.New()
	form.title.Focus()
	form.description = textarea.New()
	return form
}

func (m Form) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.description.View())
}

func main() {
	models = []tea.Model{New(), NewForm(todo)}
	m := models[model] // opens the main model as the default
	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
