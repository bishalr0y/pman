package tui

import (
	"errors"
	"testing"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/bishalr0y/pman/internal/process"
)

func TestUpdate_listErrMsg_setsStatus(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{})
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
	m := NewModel(testTable(), nil, "test", Colors{})
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
	m := NewModel(testTable(), nil, "test", Colors{})
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
	m := NewModel(testTable(), nil, "test", Colors{})
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

func TestUpdate_emptySelection_noPanic(t *testing.T) {
	m := NewModel(testTable(), nil, "test", Colors{})
	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	got := updated.(model)
	if got.pendingKillPID != 0 || cmd != nil {
		t.Fatalf("expected no confirm on empty table, pending=%d cmd=%v", got.pendingKillPID, cmd)
	}
}