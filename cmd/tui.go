package cmd

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bishalr0y/pman/internal/tui"
)

func tui_init() {
	columns := []table.Column{
		{Title: "PORT", Width: 10},
		{Title: "PID", Width: 10},
		{Title: "PROCESS NAME", Width: 15},
		{Title: "USERNAME", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(50),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("240")).
		Foreground(lipgloss.Color(tui.ColorLavender)).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color(tui.ColorLavender)).
		Bold(false)
	t.SetStyles(s)

	m := tui.NewModel(t, nil, Version)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
