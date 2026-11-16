package process

import "testing"

func TestFilter_TCPOnly(t *testing.T) {
	all := []Process{
		{ProcessName: "a", Protocol: "tcp", Port: 1},
		{ProcessName: "b", Protocol: "udp", Port: 2},
	}
	got := (Filter{TCPOnly: true}).Apply(all)
	if len(got) != 1 || got[0].Protocol != "tcp" {
		t.Fatalf("TCPOnly: got %#v, want single tcp row", got)
	}
}

func TestFilter_HideProcessNames(t *testing.T) {
	all := []Process{
		{ProcessName: "cloudflared", Protocol: "tcp", Port: 1},
		{ProcessName: "node", Protocol: "tcp", Port: 2},
	}
	got := (Filter{HideProcessNames: []string{"Cloudflared"}}).Apply(all)
	if len(got) != 1 || got[0].ProcessName != "node" {
		t.Fatalf("hide names: got %#v", got)
	}
}

func TestFilter_Combined(t *testing.T) {
	all := []Process{
		{ProcessName: "cloudflared", Protocol: "udp", Port: 1},
		{ProcessName: "cloudflared", Protocol: "tcp", Port: 2},
		{ProcessName: "node", Protocol: "tcp", Port: 3},
		{ProcessName: "node", Protocol: "udp", Port: 4},
	}
	got := (Filter{TCPOnly: true, HideProcessNames: []string{"cloudflared"}}).Apply(all)
	if len(got) != 1 || got[0].Port != 3 {
		t.Fatalf("combined filter: got %#v", got)
	}
}