package tui

import (
	"strconv"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bishalr0y/pman/internal/process"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Kill    key.Binding
	Refresh key.Binding
	Quit    key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit, k.Kill, k.Refresh}
}

// FullHelp returns keybindings for the expaned help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},      // first column
		{k.Kill, k.Refresh}, // second column
		{k.Quit},            // third column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Kill: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "kill process"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	table     table.Model
	processes []process.Process
	keys      keyMap
	help      help.Model
}

func NewModel(table table.Model, processes []process.Process) *model {
	return &model{
		table:     table,
		processes: processes,
		keys:      keys,
		help:      help.New(),
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
	helpView := m.help.View(m.keys)
	// return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")

	return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + helpView + "\n")
}
