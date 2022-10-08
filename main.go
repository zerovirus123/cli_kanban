package main // mandatory import

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
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
	loaded  bool   // wait before anything is displayed, to make sure that the lists have been initialized
	focused status // holds the list that is currently in focus
	lists   []list.Model
	err     error // display error to user at some point
}

func New() *Model {
	return &Model{}
}

// TODO: call this on tea.WindowSizeMsg
// on startup, grabs the size of the terminal window and adjust the list accordingly
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)
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
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}

	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg) //update the list
	return m, cmd
}

func (m Model) View() string {

	if m.loaded {
		return lipgloss.JoinHorizontal( // style definitions for nice terminal layouts
			lipgloss.Left,
			m.lists[todo].View(),
			m.lists[inProgress].View(),
			m.lists[done].View(),
		)
	} else {
		return "loading..."
	}

}

func main() {
	m := New()
	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
