package process

import (
	"net"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func TestListAndKillProcessWithPID(t *testing.T) {
	// start a dummy process
	cmd := exec.Command("nc", "-l", "8765")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}
	defer cmd.Process.Kill()

	time.Sleep(100 * time.Millisecond)

	procs, err := ListProcesses()
	if err != nil {
		t.Fatalf("ListProcesses failed: %v", err)
	}

	// search for the dummy process
	var foundPID int32
	found := false
	for _, p := range procs {
		if p.ProcessID == int32(cmd.Process.Pid) {
			foundPID = p.ProcessID
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("process with PID %d not found", cmd.Process.Pid)
	}

	// terminate the dummy process (SIGTERM)
	if err := KillProcessWithPID(foundPID, false); err != nil {
		t.Fatalf("KillProcess failed: %v", err)
	}

	if err := cmd.Wait(); err == nil {
		t.Error("process should have been terminated")
	}
}

func TestListProcesses_multipleListenSocketsSamePID(t *testing.T) {
	const port1 = 18765
	const port2 = 18766

	l1, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port1)))
	if err != nil {
		t.Fatalf("listen on %d: %v", port1, err)
	}
	defer l1.Close()

	l2, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port2)))
	if err != nil {
		t.Fatalf("listen on %d: %v", port2, err)
	}
	defer l2.Close()

	time.Sleep(100 * time.Millisecond)

	procs, err := ListProcesses()
	if err != nil {
		t.Fatalf("ListProcesses failed: %v", err)
	}

	pid := int32(os.Getpid())
	ports := map[uint32]bool{}
	for _, p := range procs {
		if p.ProcessID != pid || p.Address != "127.0.0.1" {
			continue
		}
		if p.Port == port1 || p.Port == port2 {
			ports[p.Port] = true
		}
	}

	if !ports[port1] || !ports[port2] {
		t.Fatalf("expected PID %d on 127.0.0.1 ports %d and %d, found ports %v", pid, port1, port2, ports)
	}
}

func TestListAndKillProcessWithPort(t *testing.T) {
	const testPort = 8767
	// start a dummy process
	cmd := exec.Command("nc", "-l", strconv.Itoa(testPort))
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start process: %v", err)
	}
	defer cmd.Process.Kill()

	time.Sleep(100 * time.Millisecond)

	procs, err := ListProcesses()
	if err != nil {
		t.Fatalf("ListProcesses failed: %v", err)
	}

	// search for the process by port
	found := false
	for _, p := range procs {
		if int32(p.Port) == testPort {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("process with port %d not found", testPort)
	}

	// terminate the dummy process by port (SIGTERM)
	if err := KillProcessWithPort(testPort, false); err != nil {
		t.Fatalf("KillProcessWithPort failed: %v", err)
	}

	if err := cmd.Wait(); err == nil {
		t.Error("process should have been terminated")
	}
}
