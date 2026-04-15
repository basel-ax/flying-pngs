# Flying PNGs Screensaver (Go Desktop Version)

A Go desktop application recreation of the classic flying Windows screensaver, now driven by a collection of PNG assets (religious and mystical iconography) with configurable animation behavior.

## Features

- Desktop application built with Go using Ebiten for rendering
- 15 PNG sprites streamed from `png/` and randomly assigned to each particle
- Adjustable particle speed, count, and canvas dimensions via config file or defaults
- Configuration persistence with automatic loading of last session
- Keyboard pause/resume shortcut (`P`)
- **Randomize Mode** — automatically toggles animation speed to ×1/5 at random intervals (0–N seconds), with an on-screen indicator and full persistence support
- Cross-platform support (Linux, Windows, macOS)

## Repository layout

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
├── png/                            # PNG sprites used by the animation (shared with web version)
├── tmp/                            # Original web version files (HTML/JS)
│   ├── index.html
│   ├── sketch.js
│   ├── agent-controls.js
│   ├── package.json
│   └── serve.json
├── README.md                       # Project documentation (this file)
└── AGENTS.md                       # Agent-facing instructions for coding assistants
```

## Local development

### Prerequisites

- Go 1.24+
- X11 development libraries (on Linux): libx11-dev, libxrandr-dev, libxcursor-dev, libxinerama-dev, libxi-dev

### Build and run

```bash
# Navigate to the Go application directory
cd flying-pngs-go

# Build and run directly
go run ./cmd/flying-pngs

# Or build first, then run
go build -o flying-pngs ./cmd/flying-pngs
./flying-pngs
```

The application will first show a settings window where you can configure:
- Number of windows (50-1000)
- Speed (1.0-20.0)
- Width and height (with valid ranges)
- Background type (black, white, transparent)
- Randomize mode (with max interval setting)
- Auto-load last session

Click "Start" to begin the animation.
Press 'P' to pause/resume the animation.

## Configuration

All settings are automatically saved to `~/.flying-pngs/config.json` and restored on subsequent launches when auto-load is enabled.

### Settings ranges

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

- `P`: Pause/resume the animation

## Randomize Mode

When enabled via the settings window, a timer fires at a random interval between 0 and N seconds. On each tick the animation speed toggles between normal and ×1/5 (slow), then a new random interval is scheduled. A yellow "SLOW ×1/5" indicator is shown in the bottom-left corner of the canvas during slow phases. The mode state and N value are persisted with the rest of the session when auto-load is on.

## Building for Distribution

To create a binary for your platform:

```bash
cd flying-pngs-go
GOOS=linux GOARCH=amd64 go build -o flying-pngs-linux ./cmd/flying-pngs
GOOS=windows GOARCH=amd64 go build -o flying-pngs-windows.exe ./cmd/flying-pngs
GOOS=darwin GOARCH=amd64 go build -o flying-pngs-macos ./cmd/flying-pngs
```

## Contributing

- Keep changes focused and maintainable.
- Update `AGENTS.md` with any new workflows or instructions useful for AI coding assistants.
- If you add assets, place them in `flying-pngs/png/` (shared with web version).
- Prefer small, incremental changes with clear intent.

## Enjoy the retro vibes!