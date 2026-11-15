package tui

import (
	"errors"
	"testing"

	"charm.land/bubbles/v2/table"
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