package config

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.yaml.in/yaml/v3"
)

//go:embed default.yaml
var defaultFS embed.FS

var (
	once         = new(sync.Once)
	cachedConfig Config
	cachedErr    error
)

type Colors struct {
	Banner        string `yaml:"banner"`
	Version       string `yaml:"version"`
	HelpKey       string `yaml:"help_key"`
	HelpDesc      string `yaml:"help_desc"`
	HelpSeparator string `yaml:"help_separator"`
	HeaderFg      string `yaml:"header_fg"`
	SelectedBg    string `yaml:"selected_bg"`
	SelectedFg    string `yaml:"selected_fg"`
	BorderFg      string `yaml:"border_fg"`
}

type Config struct {
	Colors Colors `yaml:"colors"`
}

// defaultConfig is parsed once at package init.
var defaultConfig Config

func init() {
	data, err := defaultFS.ReadFile("default.yaml")
	if err != nil {
		panic(fmt.Errorf("embedded default.yaml is missing: %w", err))
	}
	if err := yaml.Unmarshal(data, &defaultConfig); err != nil {
		panic(fmt.Errorf("failed to parse embedded default.yaml: %w", err))
	}
}

// ReadConfig reads ~/.config/pman/config.yaml, creating it from defaults if missing.
// Results are cached so subsequent calls are free.
func ReadConfig() (Config, error) {
	once.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			cachedErr = fmt.Errorf("home dir: %w", err)
			return
		}

		path := filepath.Join(homeDir, ".config/pman/config.yaml")

		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			// Write the embedded default yaml directly instead of re-marshaling
			raw, _ := defaultFS.ReadFile("default.yaml")
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				cachedErr = fmt.Errorf("create config dir %s: %w", dir, err)
				return
			}
			if err := os.WriteFile(path, raw, 0o644); err != nil {
				cachedErr = fmt.Errorf("write default config: %w", err)
				return
			}
			cachedConfig = defaultConfig
			return
		}

		if err != nil {
			cachedErr = fmt.Errorf("read %s: %w", path, err)
			return
		}

		var cfg Config
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			cachedErr = fmt.Errorf("parse %s: %w", path, err)
			return
		}
		cachedConfig = cfg
	})

	return cachedConfig, cachedErr
}