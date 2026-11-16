package process

import "strings"

// Filter trims a process list for TUI display.
type Filter struct {
	TCPOnly          bool
	HideProcessNames []string
}

// Apply returns processes matching filter options.
func (f Filter) Apply(processes []Process) []Process {
	if !f.TCPOnly && len(f.HideProcessNames) == 0 {
		return processes
	}

	hide := make(map[string]struct{}, len(f.HideProcessNames))
	for _, name := range f.HideProcessNames {
		hide[strings.ToLower(strings.TrimSpace(name))] = struct{}{}
	}

	out := make([]Process, 0, len(processes))
	for _, p := range processes {
		if f.TCPOnly && p.Protocol != "tcp" {
			continue
		}
		if _, hidden := hide[strings.ToLower(p.ProcessName)]; hidden {
			continue
		}
		out = append(out, p)
	}
	return out
}