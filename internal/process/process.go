package process

import (
	"fmt"
	"syscall"
	"time"

	gnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// Process struct (cleaned up)
type Process struct {
	Port        uint32 `json:"port"`
	ProcessID   int32  `json:"pid"`
	ProcessName string `json:"name"`
	Username    string `json:"username"`
	Protocol    string `json:"protocol"`
	StartedAt   string `json:"started_at"`
	Memory      string `json:"memory"`
}

// Helper: protocol detection
func getProtocol(connType uint32) string {
	switch connType {
	case syscall.SOCK_STREAM:
		return "tcp"
	case syscall.SOCK_DGRAM:
		return "udp"
	default:
		return "unknown"
	}
}

// Helper: format time
func formatStartTime(ms int64) string {
	t := time.Unix(0, ms*int64(time.Millisecond))
	return t.Local().Format("2006-01-02 15:04:05")
}

// Helper: human-readable memory format
func formatMemory(bytes uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func ListProcesses() ([]Process, error) {
	var processes []Process

	connections, err := gnet.Connections("inet")
	if err != nil {
		return processes, err
	}

	// avoid duplicate PID processing
	seen := make(map[int32]bool)

	for _, conn := range connections {

		if conn.Pid == 0 {
			continue
		}

		// include TCP LISTEN + UDP
		if conn.Status != "LISTEN" && conn.Type != syscall.SOCK_DGRAM {
			continue
		}

		// skip duplicate PIDs
		if seen[conn.Pid] {
			continue
		}
		seen[conn.Pid] = true

		proc, err := process.NewProcess(conn.Pid)
		if err != nil {
			continue
		}

		// basic info
		name, _ := proc.Name()
		username, _ := proc.Username()

		// start time
		createTime, _ := proc.CreateTime()
		startedAt := formatStartTime(createTime)

		// memory usage
		memInfo, err := proc.MemoryInfo()
		var memory string
		if err == nil {
			memory = formatMemory(memInfo.RSS)
		} else {
			memory = "unknown"
		}

		p := Process{
			Port:        conn.Laddr.Port,
			ProcessID:   conn.Pid,
			ProcessName: name,
			Username:    username,
			Protocol:    getProtocol(conn.Type),
			StartedAt:   startedAt,
			Memory:      memory,
		}

		processes = append(processes, p)
	}

	return processes, nil
}

// Kill by PID
func KillProcessWithPID(processID int32) error {
	proc, err := process.NewProcess(processID)
	if err != nil {
		return fmt.Errorf("process not found")
	}

	if err := proc.Kill(); err != nil {
		return err
	}
	return nil
}

// Kill by Port
func KillProcessWithPort(port int32) error {
	processes, err := ListProcesses()
	if err != nil {
		return err
	}

	for _, p := range processes {
		if int32(p.Port) == port {
			return KillProcessWithPID(p.ProcessID)
		}
	}

	return fmt.Errorf("process not found")
}
