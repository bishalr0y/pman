# pman

A terminal-based process manager with a beautiful TUI for viewing and managing listening network processes.

![pman](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-blue)

## Features

- View all processes listening on network ports
- Display PORT, PID, Process Name, and Username
- Kill processes directly from the TUI
- Auto-refresh every 2 seconds
- Manual refresh with `r`
- Vim-style navigation (`j`/`k` or arrow keys)

## Installation

### Using go install (Recommended)

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

```bash
pman
```

### Keybindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Kill selected process |
| `r` | Refresh process list |
| `q` / `Ctrl+C` | Quit |

## Requirements

- Go 1.25 or later
- macOS or Linux

## License

MIT
