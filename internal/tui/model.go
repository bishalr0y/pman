package tui

import (
	"fmt"
	"strconv"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bishalr0y/pman/internal/process"
)

func bannerStyle(c Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.Banner)).
		Bold(true)
}

func baseStyle(c Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(c.BorderFg))
}

func versionStyle(c Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.Version)).
		Italic(true)
}

var banner = `
 ██▓███   ███▄ ▄███▓ ▄▄▄       ███▄    █ 
▓██░  ██▒▓██▒▀█▀ ██▒▒████▄     ██ ▀█   █ 
▓██░ ██▓▒▓██    ▓██░▒██  ▀█▄  ▓██  ▀█ ██▒
▒██▄█▓▒ ▒▒██    ▒██ ░██▄▄▄▄██ ▓██▒  ▐▌██▒
▒██▒ ░  ░▒██▒   ░██▒ ▓█   ▓██▒▒██░   ▓██░
▒▓▒░ ░  ░░ ▒░   ░  ░ ▒▒   ▓▒█░░ ▒░   ▒ ▒ 
░▒ ░     ░  ░      ░  ▒   ▒▒ ░░ ░░   ░ ▒░
░░       ░      ░     ░   ▒      ░   ░ ░ 
                ░         ░  ░         ░ 
`

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
		key.WithHelp("enter", "confirm kill"),
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

// space for ASCII banner, help line, table header, and borders (not data rows)
const layoutChromeHeight = 15

const tableMinHeight = 8
const tableMinWidth = 98

func tableSizeForWindow(width, height int) (w, h int) {
	h = height - layoutChromeHeight
	if h < tableMinHeight {
		h = tableMinHeight
	}
	w = width - 2
	if w < tableMinWidth {
		w = tableMinWidth
	}
	return w, h
}

type model struct {
	table          table.Model
	processes      []process.Process
	status         string
	pendingKillPID int32
	keys           keyMap
	help           help.Model
	version        string
	colors         Colors
}

func statusStyle(c Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.SelectedFg)).
		Bold(true)
}

func NewModel(table table.Model, processes []process.Process, version string, colors Colors) *model {
	h := help.New()

	h.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.HelpKey)).
		Bold(true)

	h.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.HelpDesc)).
		Italic(true)

	h.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.HelpSeparator))

	return &model{
		table:     table,
		processes: processes,
		keys:      keys,
		help:      h,
		version:   version,
		colors:    colors,
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
	case tea.WindowSizeMsg:
		w, h := tableSizeForWindow(msg.Width, msg.Height)
		m.table.SetWidth(w)
		m.table.SetHeight(h)

	case listErrMsg:
		m.status = "Error listing processes: " + msg.err.Error()

	case tea.KeyMsg:
		if m.pendingKillPID != 0 {
			switch msg.String() {
			case "y", "enter":
				pid := m.pendingKillPID
				m.pendingKillPID = 0
				m.status = ""
				return m, killAndRefresh(pid, false)
			case "n", "esc":
				m.pendingKillPID = 0
				m.status = ""
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			default:
				return m, nil
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, tea.Quit
		case "r":
			return m, fetchProcesses()
		case "enter":
			selected := m.table.SelectedRow()
			if len(selected) < 3 {
				return m, nil
			}
			processID, err := strconv.ParseInt(selected[2], 10, 32)
			if err != nil {
				m.status = fmt.Sprintf("Invalid PID: %v", err)
				return m, nil
			}
			m.pendingKillPID = int32(processID)
			m.status = fmt.Sprintf("Kill PID %d? (y/n)", processID)
			return m, nil
		}

	case []process.Process:
		m.status = ""
		processes := msg
		m.processes = processes
		rows := make([]table.Row, len(processes))
		for i, p := range processes {
			rows[i] = table.Row{
				strconv.FormatUint(uint64(p.Port), 10),
				p.Address,
				strconv.FormatInt(int64(p.ProcessID), 10),
				p.ProcessName,
				p.Protocol,
				p.Memory,
				p.Username,
				p.StartedAt,
			}
		}
		m.table.SetRows(rows)

	case time.Time:
		return m, tea.Batch(fetchProcesses(), autorefresh())

	case processKilledMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Failed to kill PID %d: %v", msg.pid, msg.err)
			return m, nil
		}
		m.status = fmt.Sprintf("Killed PID %d", msg.pid)
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		bannerStyle(m.colors).Render(banner),
		versionStyle(m.colors).Render("v"+m.version),
	)

	tableView := baseStyle(m.colors).Render(m.table.View())

	var statusView string
	if m.status != "" {
		statusView = statusStyle(m.colors).Render("  "+m.status) + "\n"
	}

	helpView := "  " + m.help.View(m.keys) + "\n"

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableView,
		statusView,
		helpView,
	)

	return tea.NewView(content)
}