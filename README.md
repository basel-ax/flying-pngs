# Flying PNGs Screensaver (Go Desktop Version)

A Go desktop application recreation of the classic flying Windows screensaver with configurable animation behavior.

## First Run

```bash
git clone https://github.com/basel-ax/flying-pngs
cd flying-pngs

# Initialize go module (only needed first time)
go mod tidy

# Run the application
go run ./cmd
```

**Or build:**
```bash
go build -o flying-pngs ./cmd
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
6. [Debug Mode](#debug-mode)
7. [Randomize Mode](#randomize-mode)
8. [Building for Distribution](#building-for-distribution)
9. [Troubleshooting](#troubleshooting)
10. [Contributing](#contributing)

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
├── cmd/flying-pngs/main.go     # Application entry point
├── internal/                   # Internal packages
│   ├── config/                 # Configuration management
│   ├── ui/                     # Game interface (Ebiten)
│   ├── animation/              # Animation logic (Ebiten)
│   └── assets/                 # PNG asset loading
├── go.mod                      # Go module definition
├── png/                        # PNG sprites (shared with web version)
├── tmp/                        # Original web version (archived)
├── README.md                   # Project documentation
└── AGENTS.md                   # Agent-facing instructions
```

## Local Development

### Prerequisites

- Go 1.24+
- X11 development libraries (Linux only)

**Install X11 dependencies on Debian/Ubuntu:**
```bash
sudo apt-get install -y libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libxext-dev libxfixes-dev libxxf86vm-dev libgl1-mesa-dev
```

**Install on Fedora/RHEL:**
```bash
sudo dnf install libX11-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXext-devel libXfixes-devel libXxf86vm-devel mesa-libGL-devel
```

> **Note:** Without these libraries, compilation will fail with errors like:
> ```
> fatal error: X11/extensions/Xrandr.h: No such file or directory
> ```
> See [Troubleshooting](#troubleshooting) below for more details.

### Build and Run

```bash
cd flying-pngs
go run ./cmd
```

Or build first:
```bash
go build -o flying-pngs ./cmd
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

| Key | Context | Action |
|-----|---------|--------|
| `Up`/`Down` | Settings | Navigate between settings |
| `Left`/`Right` | Settings | Change selected value |
| `Enter`/`Space` | Settings | Start animation |
| `D` | Anywhere | Toggle debug overlay |
| `P` | Animation | Pause/resume |
| `Esc` | Animation | Return to settings |

## Randomize Mode

When enabled, a timer fires at random intervals (0–N seconds), toggling animation speed between normal and ×1/5. A "SLOW ×1/5" indicator appears during slow phases.

## Debug Mode

By default, the application only prints critical messages (`[WARN]` and `[ERROR]`). To enable verbose `[INFO]` logging (asset loading details, animation state, etc.), pass the `--debug` flag:

```bash
# Run with debug logging
go run ./cmd --debug

# Or with a built binary
./flying-pngs --debug
```

You can also toggle the on-screen debug overlay at any time by pressing `D`.

| Log Level | Shown by Default | Shown with `--debug` |
|-----------|------------------|----------------------|
| `[ERROR]` | Yes | Yes |
| `[WARN]`  | Yes | Yes |
| `[INFO]`  | No  | Yes |

## Building for Distribution

```bash
GOOS=linux GOARCH=amd64 go build -o flying-pngs-linux ./cmd
GOOS=windows GOARCH=amd64 go build -o flying-pngs-windows.exe ./cmd
GOOS=darwin GOARCH=amd64 go build -o flying-pngs-macos ./cmd
```

## Troubleshooting

### Missing X11 Development Libraries

**Symptom:** Build fails with one of these errors:
```
fatal error: X11/extensions/Xrandr.h: No such file or directory
fatal error: X11/Xlib.h: No such file or directory
fatal error: X11/extensions/XInput2.h: No such file or directory
/usr/bin/ld: cannot find -lX11
```

**Cause:** Ebiten uses GLFW for windowing, which requires X11 development headers on Linux. These are not installed by default on most distributions.

**Fix:**
```bash
# Debian / Ubuntu
sudo apt-get install -y libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev \
    libxi-dev libxext-dev libxfixes-dev libxxf86vm-dev libgl1-mesa-dev

# Fedora / RHEL
sudo dnf install libX11-devel libXrandr-devel libXcursor-devel libXinerama-devel \
    libXi-devel libXext-devel libXfixes-devel libXxf86vm-devel mesa-libGL-devel

# Arch Linux
sudo pacman -S libx11 libxrandr libxcursor libxinerama libxi libxext libxfixes libxxf86vm mesa
```

### No Images Appear After Starting

**Symptom:** Settings screen works, but after pressing Enter/Space, only a black (or gray) background is shown.

**Cause:** The app cannot find the PNG assets. It loads from `./png/*.png` relative to the current working directory.

**Fix:** Always run from the project root:
```bash
cd /path/to/flying-pngs
go run ./cmd   # or ./flying-pngs
```

### Transparent Background Still Shows Black

True window transparency requires a compositor (e.g. `picom` on X11). Without a compositor, the "transparent" background falls back to dark gray. Install and enable a compositor for full transparency support.

## Contributing

- Keep changes focused and maintainable
- Update `AGENTS.md` with new workflows
- Place new assets in `flying-pngs/png/`

Enjoy the retro vibes!