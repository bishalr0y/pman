package main

import (
	"fmt"
	"strconv"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table     table.Model
	processes []Process
}

type processKilledMsg struct {
	pid int32
	err error
}

func fetchProcesses() tea.Cmd {
	return func() tea.Msg {
		processes, err := ListProcesses()
		if err != nil {
			return err
		}
		return processes
	}
}

func killAndRefresh(pid int32) tea.Cmd {
	return func() tea.Msg {
		if err := KillProcess(pid); err != nil {
			return processKilledMsg{pid: pid, err: err}
		}

		processes, err := ListProcesses()
		if err != nil {
			return processKilledMsg{pid: pid, err: err}
		}
		return processes
	}
}

func (m model) Init() tea.Cmd {
	return fetchProcesses()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
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
			return m, killAndRefresh(int32(processID))
		}

	case []Process:
		processes := msg
		m.processes = processes
		rows := make([]table.Row, len(processes))
		for i, p := range processes {
			rows[i] = table.Row{
				strconv.FormatUint(uint64(p.Port), 10),
				strconv.FormatInt(int64(p.ProcessID), 10),
				p.ProcessName,
				p.Username,
			}
		}
		m.table.SetRows(rows)

	case processKilledMsg:
		if msg.err != nil {
			return m, tea.Printf("failed to kill process: %v", msg.err)
		}
		return m, tea.Printf("killed process: %d", msg.pid)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")
}

func main() {
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
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
