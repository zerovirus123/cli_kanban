package form

import (
	"cli_kanban/task"
	"cli_kanban/typedef"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormHandler typedef.Form

func (m FormHandler) CreateTask() tea.Msg {
	task := task.NewTask(m.Focused, m.Title.Value(), m.Description.Value())
	return task
}

func (m FormHandler) Init() tea.Cmd {
	return nil
}

func (m FormHandler) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.Title.Focused() {
				m.Title.Blur()
				m.Description.Focus()
				return m, textarea.Blink
			} else {
				typedef.Models[typedef.FormEnum] = m
				return typedef.Models[typedef.ModelEnum], m.CreateTask
			}
		}
	}

	if m.Title.Focused() {
		m.Title, cmd = m.Title.Update(msg)
		return m, cmd
	} else {
		m.Description, cmd = m.Description.Update(msg)
		return m, cmd
	}
}

func NewForm(focused typedef.Status) *FormHandler {
	form := &FormHandler{Focused: focused}
	form.Title = textinput.New()
	form.Title.Focus()
	form.Description = textarea.New()
	return form
}

func (m FormHandler) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.Title.View(), m.Description.View())
}
