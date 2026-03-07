package process

import (
	"os/exec"
	"testing"
	"time"
)

func TestListAndKillProcess(t *testing.T) {
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

	// kill the dummy process
	if err := KillProcess(foundPID); err != nil {
		t.Fatalf("KillProcess failed: %v", err)
	}

	if err := cmd.Wait(); err == nil {
		t.Error("process should have been terminated")
	}
}
