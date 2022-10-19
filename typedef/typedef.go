package typedef

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Status int

const Divisor = 4 // divide the total width of the viewport by 4

const ( //indices to determine which list is focused
	Todo Status = iota
	InProgress
	Done
)

const (
	ModelEnum Status = iota
	FormEnum
)

/* MODEL MANAGEMENT */
var Models []tea.Model

type Model struct {
	Quitting bool
	Loaded   bool   // wait before anything is displayed, to make sure that the lists have been initialized
	Focused  Status // holds the list that is currently in focus
	Lists    []list.Model
}

/* FORM MODEL */
type Form struct {
	Focused     Status
	Title       textinput.Model
	Description textarea.Model
}
