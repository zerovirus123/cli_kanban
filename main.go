package main // mandatory import

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	models = []tea.Model{New(), NewForm(todo)}
	m := models[model] // opens the main model as the default
	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
