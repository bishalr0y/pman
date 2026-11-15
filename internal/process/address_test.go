package process

import "testing"

func TestFormatBindAddress(t *testing.T) {
	tests := map[string]string{
		"":          "0.0.0.0",
		"*":         "0.0.0.0",
		"0.0.0.0":   "0.0.0.0",
		"127.0.0.1": "127.0.0.1",
		"::":        "::",
		"::1":       "::1",
	}
	for in, want := range tests {
		if got := formatBindAddress(in); got != want {
			t.Errorf("formatBindAddress(%q) = %q, want %q", in, got, want)
		}
	}
}