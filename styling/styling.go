package styling

import "github.com/charmbracelet/lipgloss"

var (
	ColumnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())

	FocusedStyle = lipgloss.NewStyle(). // styling for the focused column
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	HelpStyle = lipgloss.NewStyle(). // styling for the help text on the bottom
			Foreground(lipgloss.Color("241"))
)
