package main // mandatory import

import (
	"cli_kanban/form"
	"cli_kanban/model"
	"cli_kanban/typedef"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	typedef.Models = []tea.Model{model.New(), form.NewForm(typedef.Todo)}
	m := typedef.Models[typedef.ModelEnum] // opens the main model as the default
	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
