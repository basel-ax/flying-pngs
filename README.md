# Flying PNGs Screensaver (Go Desktop Version)

A Go desktop application recreation of the classic flying Windows screensaver with configurable animation behavior.

## Quick Start

```bash
# Clone and run
git clone https://github.com/basel-ax/flying-pngs
cd flying-pngs/flying-pngs-go
go run ./cmd/flying-pngs
```

**Or build first:**
```bash
go build -o flying-pngs ./cmd/flying-pngs
./flying-pngs
```

**Controls:** Press `P` to pause/resume

---

## Table of Contents

1. [Features](#features)
2. [Repository Layout](#repository-layout)
3. [Local Development](#local-development)
   - [Prerequisites](#prerequisites)
   - [Build and Run](#build-and-run)
4. [Configuration](#configuration)
5. [Keyboard Shortcuts](#keyboard-shortcuts)
6. [Randomize Mode](#randomize-mode)
7. [Building for Distribution](#building-for-distribution)
8. [Contributing](#contributing)

---

## Features

- Desktop application built with Go using Ebiten for rendering
- 15 PNG sprites streamed from `png/` and randomly assigned to each particle
- Adjustable particle speed, count, and canvas dimensions via config file or defaults
- Configuration persistence with automatic loading of last session
- Keyboard pause/resume shortcut (`P`)
- **Randomize Mode** — automatically toggles animation speed to ×1/5 at random intervals (0–N seconds)
- Cross-platform support (Linux, Windows, macOS)

## Repository Layout

```
flying-pngs/
├── flying-pngs-go/                 # Go desktop application
│   ├── cmd/flying-pngs/main.go     # Application entry point
│   ├── internal/                   # Internal packages
│   │   ├── config/                 # Configuration management
│   │   ├── ui/                     # Game interface (Ebiten)
│   │   ├── animation/              # Animation logic (Ebiten)
│   │   └── assets/                 # PNG asset loading
│   └── go.mod                      # Go module definition
├── png/                            # PNG sprites (shared with web version)
├── tmp/                            # Original web version (archived)
├── README.md                       # Project documentation
└── AGENTS.md                       # Agent-facing instructions
```

## Local Development

### Prerequisites

- Go 1.24+
- X11 development libraries (Linux): `libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev`

### Build and Run

```bash
cd flying-pngs-go
go run ./cmd/flying-pngs
```

Or build first:
```bash
go build -o flying-pngs ./cmd/flying-pngs
./flying-pngs
```

## Configuration

Settings are saved to `~/.flying-pngs/config.json` and restored on launch.

| Setting | Range | Default |
|---------|-------|---------|
| Number of Windows | 50-1000 | 500 |
| Speed | 1.0-20.0 | 4.0 |
| Width | 320-1920 | 1280 |
| Height | 240-1080 | 720 |
| Background Type | black, white, transparent | black |
| Randomize Mode | true/false | false |
| Max Interval (s) | 1-300 | 10 |
| Auto-load Last Session | true/false | true |

## Keyboard Shortcuts

- `P` — Pause/resume animation

## Randomize Mode

When enabled, a timer fires at random intervals (0–N seconds), toggling animation speed between normal and ×1/5. A "SLOW ×1/5" indicator appears during slow phases.

## Building for Distribution

```bash
cd flying-pngs-go
GOOS=linux GOARCH=amd64 go build -o flying-pngs-linux ./cmd/flying-pngs
GOOS=windows GOARCH=amd64 go build -o flying-pngs-windows.exe ./cmd/flying-pngs
GOOS=darwin GOARCH=amd64 go build -o flying-pngs-macos ./cmd/flying-pngs
```

## Contributing

- Keep changes focused and maintainable
- Update `AGENTS.md` with new workflows
- Place new assets in `flying-pngs/png/`

Enjoy the retro vibes!