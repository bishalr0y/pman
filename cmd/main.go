package main

import (
	"fmt"
	"os"
	"strconv"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			processID, err := strconv.ParseInt(m.table.SelectedRow()[1], 10, 32)
			if err != nil {
				return m, tea.Printf("failed to parse the processID: %v", err)
			}
			if err := KillProcess(int32(processID)); err != nil {
				return m, tea.Printf("failed to kill the process: %v", err)
			}
			return m, tea.Batch(
				tea.Printf("killed process: %d!", processID),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")
}

func main() {
	KillProcess(86420)
	columns := []table.Column{
		{Title: "PORT", Width: 10},
		{Title: "PID", Width: 10},
		{Title: "PROCESS NAME", Width: 15},
		{Title: "USERNAME", Width: 10},
	}
	rows := []table.Row{}

	processes, err := ListProcesses()
	if err != nil {
		// TODO: handle the error message
		fmt.Printf("Error: %v", err)
		os.Exit(0)
	}

	for _, p := range processes {
		row := table.Row{
			strconv.FormatUint(uint64(p.Port), 10),
			strconv.FormatInt(int64(p.ProcessID), 10),
			p.ProcessName,
			p.Username,
		}

		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(50),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
