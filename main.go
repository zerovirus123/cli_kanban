package main // mandatory import

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const ( //indices to determine which list is focused
	todo status = iota
	inProgress
	done
)

/* CUSTOM ITEM */

type Task struct { // implementing the list.item inteface
	status      status
	title       string
	description string
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
	lists list.Model
	err   error // display error to user at some point
}

func New() *Model {
	return &Model{}
}

// TODO: call this on tea.WindowSizeMsg
// on startup, grabs the size of the terminal window and adjust the list accordingly
func (m *Model) initLists(width, height int) {
	m.lists = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.lists.Title = "To Do"
	m.lists.SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t shirts"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil // no timer to startup when the program starts
}

// update the list (passing the interactions and keypresses, no logic involved)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // terminal dimensions on program startup
		m.initLists(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.lists, cmd = m.lists.Update(msg) //update the list
	return m, cmd
}

func (m Model) View() string {
	return m.lists.View()
}

func main() {
	m := New()
	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
