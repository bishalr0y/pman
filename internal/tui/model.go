package tui

import (
	"strconv"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bishalr0y/pman/internal/process"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table     table.Model
	processes []process.Process
}

func NewModel(table table.Model, processes []process.Process) *model {
	return &model{
		table:     table,
		processes: processes,
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
		case "q", "ctrl+c":
			// quit
			return m, tea.Quit
		case "r":
			// refresh (refetch the processes)
			return m, fetchProcesses()
		case "enter":
			// kill the process
			processID, err := strconv.ParseInt(m.table.SelectedRow()[1], 10, 32)
			if err != nil {
				return m, tea.Printf("failed to parse the processID: %v", err)
			}
			return m, killAndRefresh(int32(processID))
		}

	case []process.Process:
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
