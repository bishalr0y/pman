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
	Address     string `json:"address"`
	Port        uint32 `json:"port"`
	ProcessID   int32  `json:"pid"`
	ProcessName string `json:"name"`
	Protocol    string `json:"protocol"`
	Memory      string `json:"memory"`
	Username    string `json:"username"`
	StartedAt   string `json:"started_at"`
}

type pidInfo struct {
	name      string
	username  string
	memory    string
	startedAt string
}

func formatBindAddress(ip string) string {
	if ip == "" {
		return "0.0.0.0"
	}
	return ip
}

func socketKey(pid int32, protocol, address string, port uint32) string {
	return fmt.Sprintf("%d|%s|%s|%d", pid, protocol, address, port)
}

func lookupPIDInfo(pid int32, cache map[int32]pidInfo) (pidInfo, bool) {
	info, ok := cache[pid]
	return info, ok
}

func fetchPIDInfo(pid int32) (pidInfo, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return pidInfo{}, err
	}

	name, _ := proc.Name()
	username, err := proc.Username()
	if err != nil || username == "" {
		username = "unknown"
	}

	createTime, _ := proc.CreateTime()
	startedAt := formatStartTime(createTime)

	memInfo, err := proc.MemoryInfo()
	var memory string
	if err == nil {
		memory = formatMemory(memInfo.RSS)
	} else {
		memory = "unknown"
	}

	return pidInfo{
		name:      name,
		username:  username,
		memory:    memory,
		startedAt: startedAt,
	}, nil
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

func listConnections() ([]gnet.ConnectionStat, error) {
	var connections []gnet.ConnectionStat
	for _, kind := range []string{"inet", "inet6"} {
		conns, err := gnet.Connections(kind)
		if err != nil {
			return nil, err
		}
		connections = append(connections, conns...)
	}
	return connections, nil
}

func ListProcesses() ([]Process, error) {
	var processes []Process

	connections, err := listConnections()
	if err != nil {
		return processes, err
	}

	seenSocket := make(map[string]bool)
	pidCache := make(map[int32]pidInfo)

	for _, conn := range connections {

		if conn.Pid == 0 {
			continue
		}

		// include TCP LISTEN + UDP
		if conn.Status != "LISTEN" && conn.Type != syscall.SOCK_DGRAM {
			continue
		}

		protocol := getProtocol(conn.Type)
		address := formatBindAddress(conn.Laddr.IP)
		key := socketKey(conn.Pid, protocol, address, conn.Laddr.Port)
		if seenSocket[key] {
			continue
		}
		seenSocket[key] = true

		info, ok := lookupPIDInfo(conn.Pid, pidCache)
		if !ok {
			var err error
			info, err = fetchPIDInfo(conn.Pid)
			if err != nil {
				continue
			}
			pidCache[conn.Pid] = info
		}

		processes = append(processes, Process{
			Address:     address,
			Port:        conn.Laddr.Port,
			ProcessID:   conn.Pid,
			ProcessName: info.name,
			Protocol:    protocol,
			Memory:      info.memory,
			Username:    info.username,
			StartedAt:   info.startedAt,
		})
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
