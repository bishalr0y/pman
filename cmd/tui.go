package cmd

import (
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bishalr0y/pman/internal/config"
	"github.com/bishalr0y/pman/internal/tui"
)

func tuiInit() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	colors := tui.Colors{
		Banner:        cfg.Colors.Banner,
		Version:       cfg.Colors.Version,
		HelpKey:       cfg.Colors.HelpKey,
		HelpDesc:      cfg.Colors.HelpDesc,
		HelpSeparator: cfg.Colors.HelpSeparator,
		HeaderFg:      cfg.Colors.HeaderFg,
		SelectedBg:    cfg.Colors.SelectedBg,
		SelectedFg:    cfg.Colors.SelectedFg,
		BorderFg:      cfg.Colors.BorderFg,
	}

	columns := []table.Column{
		{Title: "PORT", Width: 10},
		{Title: "PID", Width: 10},
		{Title: "PROCESS NAME", Width: 15},
		{Title: "PROTOCOL", Width: 10},
		{Title: "MEMORY", Width: 10},
		{Title: "USERNAME", Width: 10},
		{Title: "STARTED AT", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(100),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color(colors.BorderFg)).
		Foreground(lipgloss.Color(colors.HeaderFg)).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(colors.SelectedFg)).
		Background(lipgloss.Color(colors.SelectedBg)).
		Bold(false)
	t.SetStyles(s)

	m := tui.NewModel(t, nil, Version, colors)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}