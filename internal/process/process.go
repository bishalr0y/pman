package process

import (
	"fmt"

	gnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// TODO: add more fields (status, memory, start time, protocol)
type Process struct {
	Port        uint32 `json:"port"`
	ProcessID   int32  `json:"pid"`
	ProcessName string `json:"name"`
	Username    string `json:"username"`
}

func ListProcesses() ([]Process, error) {
	var processes []Process
	connections, err := gnet.Connections("inet")
	if err != nil {
		return processes, err
	}

	seen := make(map[int32]bool) // avoid duplicate PID lookups

	for _, conn := range connections {
		if conn.Status == "LISTEN" && conn.Pid != 0 {

			// prevent repeating same PID multiple times
			if seen[conn.Pid] {
				continue
			}
			seen[conn.Pid] = true

			proc, err := process.NewProcess(conn.Pid)
			if err != nil {
				continue
			}

			processName, err := proc.Name()
			if err != nil {
				return []Process{}, err
			}

			username, err := proc.Username()
			if err != nil {
				return []Process{}, err
			}

			process := Process{
				Port:        conn.Laddr.Port,
				ProcessID:   conn.Pid,
				ProcessName: processName,
				Username:    username,
			}

			processes = append(processes, process)
		}
	}
	return processes, nil
}

func KillProcessWithPID(processID int32) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		id := p.Pid
		if id == processID {
			if err := p.Kill(); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("process not found\n")
}

func KillProcessWithPort(port int32) error {
	processes, err := ListProcesses()
	if err != nil {
		return err
	}
	for _, process := range processes {
		p := process.Port
		if port == int32(p) {
			if err := KillProcessWithPID(process.ProcessID); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("process not found\n")
}
