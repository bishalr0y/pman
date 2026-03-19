package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/bishalr0y/pman/internal/process"
)

type processKilledMsg struct {
	pid int32
	err error
}

func fetchProcesses() tea.Cmd {
	return func() tea.Msg {
		processes, err := process.ListProcesses()
		if err != nil {
			return err
		}
		return processes
	}
}

func killAndRefresh(pid int32) tea.Cmd {
	return func() tea.Msg {
		if err := process.KillProcessWithPID(pid); err != nil {
			return processKilledMsg{pid: pid, err: err}
		}

		processes, err := process.ListProcesses()
		if err != nil {
			return processKilledMsg{pid: pid, err: err}
		}
		return processes
	}
}
