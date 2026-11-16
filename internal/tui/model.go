package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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

var banner = ` ██▓███   ███▄ ▄███▓ ▄▄▄       ███▄    █ 
▓██░  ██▒▓██▒▀█▀ ██▒▒████▄     ██ ▀█   █ 
▓██░ ██▓▒▓██    ▓██░▒██  ▀█▄  ▓██  ▀█ ██▒
▒██▄█▓▒ ▒▒██    ▒██ ░██▄▄▄▄██ ▓██▒  ▐▌██▒
▒██▒ ░  ░▒██▒   ░██▒ ▓█   ▓██▒▒██░   ▓██░
▒▓▒░ ░  ░░ ▒░   ░  ░ ▒▒   ▓▒█░░ ▒░   ▒ ▒ 
░▒ ░     ░  ░      ░  ▒   ▒▒ ░░ ░░   ░ ▒░
░░       ░      ░     ░   ▒      ░   ░ ░ 
                ░         ░  ░         ░`

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Kill      key.Binding
	Refresh   key.Binding
	ToggleUDP key.Binding
	Quit      key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.ToggleUDP, k.Kill, k.Refresh, k.Quit}
}

// FullHelp returns keybindings for the expaned help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Kill},
		{k.Refresh, k.ToggleUDP, k.Quit},
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
	ToggleUDP: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "toggle udp"),
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
	table            table.Model
	processes        []process.Process
	status           string
	pendingKillPID   int32
	showUDP          bool
	hideProcessNames []string
	keys             keyMap
	help             help.Model
	version          string
	colors           Colors
}

func statusStyle(c Colors) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.SelectedFg)).
		Bold(true)
}

func NewModel(table table.Model, processes []process.Process, version string, colors Colors, hideProcessNames []string) *model {
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
		table:            table,
		processes:        processes,
		hideProcessNames: hideProcessNames,
		keys:             keys,
		help:             h,
		version:          version,
		colors:           colors,
	}
}

func (m *model) displayFilter() process.Filter {
	return process.Filter{
		TCPOnly:          !m.showUDP,
		HideProcessNames: m.hideProcessNames,
	}
}

type killSocket struct {
	proto  string
	port   uint32
	hidden bool
}

type killTarget struct {
	name    string
	pid     int32
	sockets []killSocket
}

func (m model) killTargetFor(pid int32, selected table.Row) killTarget {
	t := killTarget{pid: pid}
	seen := make(map[string]struct{})
	filter := m.displayFilter()
	for _, p := range m.processes {
		if p.ProcessID != pid {
			continue
		}
		if t.name == "" {
			t.name = p.ProcessName
		}
		key := fmt.Sprintf("%s:%d", p.Protocol, p.Port)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		hidden := len(filter.Apply([]process.Process{p})) == 0
		t.sockets = append(t.sockets, killSocket{proto: p.Protocol, port: p.Port, hidden: hidden})
	}
	if t.name == "" && len(selected) > 3 {
		t.name = selected[3]
	}
	if len(t.sockets) == 0 && len(selected) > 0 {
		proto := "tcp"
		if len(selected) > 4 && selected[4] != "" {
			proto = selected[4]
		}
		port, _ := strconv.ParseUint(selected[0], 10, 32)
		t.sockets = append(t.sockets, killSocket{proto: proto, port: uint32(port)})
	}
	sort.Slice(t.sockets, func(i, j int) bool {
		if t.sockets[i].proto != t.sockets[j].proto {
			return t.sockets[i].proto < t.sockets[j].proto
		}
		return t.sockets[i].port < t.sockets[j].port
	})
	return t
}

func (t killTarget) dialog(c Colors) string {
	const inner = 36

	title := lipgloss.NewStyle().
		Width(inner).
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color(orColor(c.SelectedFg, "#f9e2af"))).
		Background(lipgloss.Color(orColor(c.SelectedBg, "#7287fd"))).
		Render("Confirm kill")

	muted := lipgloss.NewStyle().Foreground(lipgloss.Color(orColor(c.HelpDesc, "#babbf1")))
	name := lipgloss.NewStyle().
		Width(inner).
		Bold(true).
		Foreground(lipgloss.Color(orColor(c.Banner, "#8bd5ca")))
	section := lipgloss.NewStyle().
		Width(inner).
		Foreground(lipgloss.Color(orColor(c.HeaderFg, "#7287fd"))).
		Bold(true)
	protoCol := lipgloss.NewStyle().Width(6).Foreground(lipgloss.Color(orColor(c.HeaderFg, "#7287fd")))
	portCol := lipgloss.NewStyle().Width(7).Align(lipgloss.Right)
	noteCol := muted.Italic(true)
	key := lipgloss.NewStyle().Width(8).Bold(true).Foreground(lipgloss.Color(orColor(c.HelpKey, "#a6e3a1")))
	gap := muted.Render(strings.Repeat("─", inner))

	rows := []string{title, ""}
	if t.name != "" {
		rows = append(rows, name.Render(t.name))
	}
	rows = append(rows, muted.Width(inner).Render(fmt.Sprintf("PID  %d", t.pid)), gap)
	if len(t.sockets) == 0 {
		rows = append(rows, muted.Width(inner).Render("No listen sockets recorded."))
	} else {
		rows = append(rows, section.Render("Listening"))
		for _, s := range t.sockets {
			note := ""
			if s.hidden {
				note = "  " + noteCol.Render("hidden")
			}
			line := lipgloss.JoinHorizontal(lipgloss.Top,
				protoCol.Render(s.proto),
				portCol.Render(strconv.FormatUint(uint64(s.port), 10)),
				note,
			)
			rows = append(rows, lipgloss.NewStyle().Width(inner).Render(line))
		}
	}
	rows = append(rows, gap,
		lipgloss.NewStyle().Width(inner).Render(
			lipgloss.JoinHorizontal(lipgloss.Top, key.Render("y/enter"), muted.Render("kill")),
		),
		lipgloss.NewStyle().Width(inner).Render(
			lipgloss.JoinHorizontal(lipgloss.Top, key.Render("n/esc"), muted.Render("cancel")),
		),
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(orColor(c.SelectedBg, "#7287fd"))).
		Padding(1, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func orColor(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func (m *model) applyFilterToTable() {
	visible := m.displayFilter().Apply(m.processes)
	rows := make([]table.Row, len(visible))
	for i, p := range visible {
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
		case "u":
			m.showUDP = !m.showUDP
			m.applyFilterToTable()
			if m.showUDP {
				m.status = "Showing TCP and UDP"
			} else {
				m.status = "Showing TCP only"
			}
			return m, nil
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
			return m, nil
		}

	case []process.Process:
		m.processes = msg
		if m.pendingKillPID == 0 && strings.HasPrefix(m.status, "Error listing processes:") {
			m.status = ""
		}
		m.applyFilterToTable()

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
	if m.pendingKillPID == 0 && m.status != "" {
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

	if m.pendingKillPID != 0 {
		dialog := m.killTargetFor(m.pendingKillPID, m.table.SelectedRow()).dialog(m.colors)
		base := lipgloss.NewLayer(content)
		dw, dh := lipgloss.Width(dialog), lipgloss.Height(dialog)
		bw, bh := lipgloss.Width(content), lipgloss.Height(content)
		x := max((bw-dw)/2, 0)
		y := max((bh-dh)/2, 1)
		overlay := lipgloss.NewLayer(dialog).X(x).Y(y).Z(1)
		content = lipgloss.NewCompositor(base, overlay).Render()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}