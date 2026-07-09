package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	if defaultConfig.Colors.Banner == "" {
		t.Fatal("expected banner color to be set in default config")
	}
	if defaultConfig.Colors.Version == "" {
		t.Fatal("expected version color to be set in default config")
	}
}

func TestReadConfig_CreatesDefault(t *testing.T) {
	saved := once
	once = new(sync.Once)
	t.Cleanup(func() { once = saved })

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	cfg, err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() returned error: %v", err)
	}

	if cfg.Colors.Banner == "" {
		t.Fatal("expected banner color after creating default config")
	}

	configPath := filepath.Join(tmpHome, ".config/pman/config.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("expected config file to be created at %s: %v", configPath, err)
	}

	cfg2, err := ReadConfig()
	if err != nil {
		t.Fatalf("second ReadConfig() returned error: %v", err)
	}
	if cfg2.Colors.Banner != cfg.Colors.Banner {
		t.Error("cached config differs from first read")
	}
}

func TestReadConfig_ExistingFile(t *testing.T) {
	saved := once
	once = new(sync.Once)
	t.Cleanup(func() { once = saved })

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	configDir := filepath.Join(tmpHome, ".config/pman")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configContent := `colors:
  banner: "#ff0000"
  version: "#00ff00"
  help_key: "#a6e3a1"
  help_desc: "#babbf1"
  help_separator: "#8bd5ca"
  header_fg: "#7287fd"
  selected_bg: "#7287fd"
  selected_fg: "#f9e2af"
  border_fg: "#6c7086"
`
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() returned error: %v", err)
	}

	if cfg.Colors.Banner != "#ff0000" {
		t.Errorf("expected banner #ff0000, got %q", cfg.Colors.Banner)
	}
	if cfg.Colors.Version != "#00ff00" {
		t.Errorf("expected version #00ff00, got %q", cfg.Colors.Version)
	}
}
