package tui

import (
	"strconv"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bishalr0y/pman/internal/process"
)

var bannerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7287fd")).
	Bold(true)

var banner = "   ___  __ _  ___ ____ \n" +
	"  / _ \\/  ' \\/ _ `/ _ \\\n" +
	" / .__/_/_/_/\\_,_/_//_/\n" +
	"/_/                      "

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
		{k.Up, k.Down, k.Kill}, // first column
		{k.Refresh, k.Quit},    // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Kill: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "kill"),
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
	h := help.New()

	h.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Bold(true)

	h.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#babbf1")).
		Italic(true)

	h.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#209fb5"))

	return &model{
		table:     table,
		processes: processes,
		keys:      keys,
		help:      h,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchProcesses(), autorefresh())
}

func autorefresh() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return t
	})
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

	case time.Time:
		return m, tea.Batch(fetchProcesses(), autorefresh())

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
	header := bannerStyle.Render(banner)

	tableView := baseStyle.Render(m.table.View())

	helpView := "  " + m.help.View(m.keys) + "\n"

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableView,
		helpView,
	)

	return tea.NewView(content)
}
