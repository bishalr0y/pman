package main

import (
	"fmt"

	gnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type Process struct {
	Port        uint32
	ProcessID   int32
	ProcessName string
	Username    string
}

func ListProcesses() ([]Process, error) {
	var processes []Process
	connections, err := gnet.Connections("inet")
	if err != nil {
		panic(err)
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

func KillProcess(processID int32) error {
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
