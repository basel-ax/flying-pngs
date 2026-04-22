# AGENTS.md

This repository contains a Go desktop application that renders a Windows-logo "flying" screensaver style animation. This AGENTS.md is written to help AI coding agents understand how to run, inspect, and improve the project with safe, targeted changes.

## Project overview

- Location: project root
- Entry point: `cmd/main.go` (Go application using Ebiten for rendering)
- Visual: 15 PNG assets loaded from the `png/` folder (same assets as the web version)
- Audio: disabled (no sound usage)
- Tech stack: Go, Ebiten (game engine/rendering)

## How to run

- Build and run (requires X11 display):
  ```
  cd flying-pngs
  go run ./cmd
  ```
- Or build first:
  ```
  go build -o flying-pngs ./cmd
  ./flying-pngs
  ```
- Configuration is loaded from config file or uses defaults:
  - Number of windows (50-1000)
  - Speed (1.0-20.0)
  - Width and height (with valid ranges)
  - Background type (black, white, transparent)
  - Randomize mode (with max interval setting)
  - Auto-load last session
- Press 'P' to pause/resume the animation

## Key files and roles

- `cmd/main.go` — application entry point; creates game instance
- `internal/config/config.go` — configuration management (loading, saving, defaults, validation)
- `internal/ui/ui.go` — game interface implementation (Ebiten)
- `internal/animation/animation.go` — animation logic using Ebiten
- `internal/assets/assets.go` — PNG asset loading utility
- `internal/` — internal packages containing the core logic
- `png/` — image assets (15 PNGs, shared with web version)

## Agent-facing workflow

Agents should aim to improve the project with safe, incremental changes that preserve the visual intent while improving robustness, performance, or add configurability.

## Build and test

```bash
go build -o flying-pngs ./cmd
go test ./internal/config/...
```

## Known considerations

- This project does not have automated tests at the GUI level; visual/manual checks are the primary verification method
- Unit tests exist for configuration and utility functions
- The codebase is designed to be approachable, with clear separation of concerns
- X11 display required to run the application

## Agent-facing changes

If you add subpackages under `internal/` or elsewhere, you can place an AGENTS.md there as well. Agents will generally read the nearest AGENTS.md in the directory tree, and the closest one to a modified file should take precedence.