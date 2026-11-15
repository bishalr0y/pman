package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/bishalr0y/pman/internal/process"
)

type processKilledMsg struct {
	pid int32
	err error
}

type listErrMsg struct {
	err error
}

func fetchProcesses() tea.Cmd {
	return func() tea.Msg {
		processes, err := process.ListProcesses()
		if err != nil {
			return listErrMsg{err: err}
		}
		return processes
	}
}

func killAndRefresh(pid int32, force bool) tea.Cmd {
	return func() tea.Msg {
		if err := process.KillProcessWithPID(pid, force); err != nil {
			return processKilledMsg{pid: pid, err: err}
		}

		processes, err := process.ListProcesses()
		if err != nil {
			return listErrMsg{err: err}
		}
		return processes
	}
}
