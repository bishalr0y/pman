# pman

<p align="center">
  <img src="assets/gopher_transparent.png" alt="pman mascot" width="500"/>
</p>

A terminal-based process manager with a beautiful TUI for viewing and managing listening network processes.

![platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-blue)
![Go](https://img.shields.io/badge/Go-1.26+-00ADD8)

## Preview

TCP listeners are shown by default (press `u` in the TUI to include UDP). Each row is one listen socket with **PORT**, **BIND**, and **PID**.

![pman screenshot](assets/screenshot.png?v=7)

## Features

### TUI Features
- One row per listen socket (same PID can appear on multiple ports or binds)
- Columns: PORT, BIND, PID, process name, protocol, memory, user, started at
- TCP-only view by default; `u` toggles UDP listeners
- Terminal-sized table (resizes with your window)
- Human-readable memory (KB, MB, GB)
- Auto-refresh every 2 seconds; manual refresh with `r`
- Vim-style navigation (`j`/`k` or arrow keys)
- Confirm before kill (`Enter` then `y`); SIGTERM by default
- Theme colors via `~/.config/pman/config.yaml`
- Hide process names in the TUI via `display.hide_process_names` (JSON is always complete)

### CLI Commands
- `pman` - Launch interactive TUI (default)
- `pman json` - Output process list as JSON
- `pman kill [--force] <pid>` - Terminate process by PID (SIGTERM; `--force` for SIGKILL)
- `pman killport [--force] <port>` - Terminate process on port (SIGTERM; `--force` for SIGKILL)
- `pman version` - Show version information

## Installation

### Using Homebrew (Recommended on macOS)

```bash
brew tap bishalr0y/homebrew-bishalr0y
brew install pman
```

### Using go install

```bash
go install github.com/bishalr0y/pman@latest
```

### From Source

```bash
git clone https://github.com/bishalr0y/pman.git
cd pman
go install ./cmd/...
```

## Usage

### TUI Mode

```bash
pman
```

### Keybindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Confirm kill in a floating dialog (process name, PID, and all listen ports, including UDP hidden by the default filter) |
| `r` | Refresh process list |
| `u` | Toggle UDP listeners (TCP-only by default) |
| `q` / `Esc` / `Ctrl+C` | Quit |

### CLI Examples

```bash
# Full process list (all protocols); each object includes "address" and "port"
pman json

# Script-friendly: failures exit 1
pman kill not-a-pid   # exit 1

# Terminate by PID (SIGTERM)
pman kill 1234

# Force kill (SIGKILL)
pman kill --force 1234

# Terminate listener on a port (first match)
pman killport 8080
```

## Configuration

pman reads config from `~/.config/pman/config.yaml`. On first run, it auto-creates the file with defaults:

```yaml
display:
  hide_process_names: []

colors:
  banner: "#8bd5ca"
  version: "#7287fd"
  help_key: "#a6e3a1"
  help_desc: "#babbf1"
  help_separator: "#8bd5ca"
  header_fg: "#7287fd"
  selected_bg: "#7287fd"
  selected_fg: "#f9e2af"
  border_fg: "#6c7086"
```

`display.hide_process_names` hides processes from the **TUI** by exact name (case-insensitive). `pman json` always lists everything. Example:

```yaml
display:
  hide_process_names:
    - cloudflared
    - rapportd
```

Color keys control specific UI elements — change any hex value to customize the theme.

## Requirements

- Go 1.26 or later (to build from source)
- macOS or Linux

## Releasing

Goreleaser publishes GitHub assets and Homebrew when a GitHub release is created.

```bash
# 1. bump const Version in cmd/version.go and add a CHANGELOG.md section
# 2. commit and push to main
# 3. tag and publish:
gh release create vX.Y.Z --title "vX.Y.Z" --notes "..." --target main
```

## Regenerating the README screenshot

```bash
go build -o bin/pman .
vhs scripts/screenshot.tape   # writes assets/screenshot.png
```

## License

MIT
