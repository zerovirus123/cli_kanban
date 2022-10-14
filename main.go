package main // mandatory import

import (
	"cli_kanban/typedef"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	models = []tea.Model{New(), NewForm(typedef.Todo)}
	m := models[typedef.Model] // opens the main model as the default
	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
