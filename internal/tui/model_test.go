package tui

import (
	"errors"
	"testing"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/bishalr0y/pman/internal/process"
)

func TestUpdate_listErrMsg_setsStatus(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{}, nil)
	updated, _ := m.Update(listErrMsg{err: errors.New("list failed")})
	got := updated.(model)
	if got.status != "Error listing processes: list failed" {
		t.Fatalf("status = %q, want list error message", got.status)
	}
}

func testTable() table.Model {
	columns := []table.Column{
		{Title: "PORT", Width: 8},
		{Title: "BIND", Width: 15},
		{Title: "PID", Width: 8},
		{Title: "PROCESS NAME", Width: 15},
		{Title: "PROTOCOL", Width: 10},
		{Title: "MEMORY", Width: 10},
		{Title: "USERNAME", Width: 10},
		{Title: "STARTED AT", Width: 22},
	}
	return table.New(table.WithColumns(columns))
}

func TestUpdate_processList_clearsStatus(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{}, nil)
	m.status = "Error listing processes: list failed"

	updated, _ := m.Update([]process.Process{{
		Port: 8080, ProcessID: 1, ProcessName: "test",
		Address: "127.0.0.1", Protocol: "tcp",
		Memory: "1 MB", Username: "u", StartedAt: "now",
	}})
	got := updated.(model)
	if got.status != "" {
		t.Fatalf("status = %q, want empty after successful list", got.status)
	}
}

func seedTableRow(m *model) {
	m.table.SetRows([]table.Row{
		{"8080", "127.0.0.1", "42", "test", "tcp", "1 MB", "u", "now"},
	})
}

func TestUpdate_enter_startsKillConfirm(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{}, nil)
	seedTableRow(m)

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)
	if got.pendingKillPID != 42 {
		t.Fatalf("pendingKillPID = %d, want 42", got.pendingKillPID)
	}
	if got.status != "Kill PID 42? (y/n)" {
		t.Fatalf("status = %q, want confirm prompt", got.status)
	}
	if cmd != nil {
		t.Fatal("expected no cmd while awaiting confirmation")
	}
}

func TestUpdate_confirmCancel_doesNotKill(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{}, nil)
	seedTableRow(m)

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)

	updated, cmd := got.Update(tea.KeyPressMsg{Text: "n"})
	got = updated.(model)
	if got.pendingKillPID != 0 {
		t.Fatalf("pendingKillPID = %d, want 0 after cancel", got.pendingKillPID)
	}
	if cmd != nil {
		t.Fatal("expected no kill cmd on cancel")
	}
}

func TestUpdate_processList_hidesUDPByDefault(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{}, nil)

	updated, _ := m.Update([]process.Process{
		{Port: 1, ProcessID: 1, ProcessName: "tcp-svc", Protocol: "tcp", Address: "127.0.0.1", Memory: "1", Username: "u", StartedAt: "now"},
		{Port: 2, ProcessID: 2, ProcessName: "udp-svc", Protocol: "udp", Address: "127.0.0.1", Memory: "1", Username: "u", StartedAt: "now"},
	})
	got := updated.(model)

	if len(got.table.Rows()) != 1 {
		t.Fatalf("expected 1 visible row (tcp only), got %d", len(got.table.Rows()))
	}
	if got.table.Rows()[0][3] != "tcp-svc" {
		t.Fatalf("expected tcp row, got %v", got.table.Rows()[0])
	}
}

func TestUpdate_toggleUDP_showsUDPRows(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{}, nil)

	updated, _ := m.Update([]process.Process{
		{Port: 1, ProcessID: 1, ProcessName: "tcp-svc", Protocol: "tcp", Address: "127.0.0.1", Memory: "1", Username: "u", StartedAt: "now"},
		{Port: 2, ProcessID: 2, ProcessName: "udp-svc", Protocol: "udp", Address: "127.0.0.1", Memory: "1", Username: "u", StartedAt: "now"},
	})
	got := updated.(model)

	updated, _ = got.Update(tea.KeyPressMsg{Text: "u"})
	got = updated.(model)
	if len(got.table.Rows()) != 2 {
		t.Fatalf("expected 2 rows with UDP shown, got %d", len(got.table.Rows()))
	}
}

func TestUpdate_hideProcessNamesFromConfig(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{}, []string{"cloudflared"})

	updated, _ := m.Update([]process.Process{
		{Port: 1, ProcessID: 1, ProcessName: "cloudflared", Protocol: "tcp", Address: "0.0.0.0", Memory: "1", Username: "u", StartedAt: "now"},
		{Port: 2, ProcessID: 2, ProcessName: "node", Protocol: "tcp", Address: "127.0.0.1", Memory: "1", Username: "u", StartedAt: "now"},
	})
	got := updated.(model)
	if len(got.table.Rows()) != 1 || got.table.Rows()[0][3] != "node" {
		t.Fatalf("expected only node visible, got %v", got.table.Rows())
	}
}

func TestUpdate_emptySelection_noPanic(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{}, nil)
	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)
	if got.pendingKillPID != 0 || cmd != nil {
		t.Fatalf("expected no confirm on empty table, pending=%d cmd=%v", got.pendingKillPID, cmd)
	}
}