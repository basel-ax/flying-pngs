# AGENTS.md

This repository contains a Go desktop application that renders a Windows-logo "flying" screensaver style animation. This AGENTS.md is written to help AI coding agents understand how to run, inspect, and improve the project with safe, targeted changes.

## Project overview

- Location: `flying-pngs-go/`
- Entry point: `flying-pngs-go/cmd/flying-pngs/main.go` (Go application using Ebiten for rendering)
- Visual: 15 PNG assets loaded from the `flying-pngs/png/` folder (same assets as the web version)
- Audio: disabled (no sound usage)
- Tech stack: Go, Ebiten (game engine/rendering)

## How to run

- Build and run (requires X11 display):
  ```
  cd flying-pngs-go
  go run ./cmd/flying-pngs
  ```
- Or build first:
  ```
  cd flying-pngs-go
  go build -o flying-pngs ./cmd/flying-pngs
  ./flying-pngs
  ```
- Configuration is loaded from config file or uses defaults:
  - Number of windows (50-1000)
  - Speed (1.0-20.0)
  - Width and height (with valid ranges)
  - Background type (black, white, transparent)
  - Randomize mode (with max interval setting)
  - Auto-load last session
- Click "Start" to begin the animation
- Press 'P' to pause/resume the animation

## Key files and roles

- `flying-pngs-go/cmd/flying-pngs/main.go` — application entry point; initializes Fyne app, loads config, shows settings window
- `flying-pngs-go/internal/config/config.go` — configuration management (loading, saving, defaults, validation)
- `flying-pngs-go/internal/ui/ui.go` — settings window and animation window implementations using Fyne
- `flying-pngs-go/internal/animation/animation.go` — animation logic adapted from the original sketch.js (using Ebiten)
- `flying-pngs-go/internal/assets/assets.go` — PNG asset loading utility
- `flying-pngs-go/internal/` — internal packages containing the core logic
- `flying-pngs-go/cmd/` — main application package
- `flying-pngs/png/` — image assets (15 PNGs, shared with web version)

## Agent-facing workflow

Agents should aim to improve the project with safe, incremental changes that preserve the visual intent while improving robustness, performance, or add configurability.

- Non-invasive improvements
  - Add runtime configurability via command line flags or environment variables
  - Implement a responsive canvas that adapts to window size while preserving aspect ratio
  - Respect devicePixelRatio for crisp rendering on high-DPI displays
  - Provide a simple fallback if assets fail to load (e.g., draw a basic shape instead of logos)
  - Add keyboard shortcuts for common actions
  - Add system tray integration

- Performance and robustness
  - Allow the number of windows to be capped, and adjust speed based on a parameter or performance hints
  - Add error handling around asset loading; log helpful warnings and fall back gracefully
  - Minimize draw calls or optimize math in projection to reduce CPU usage on low-end devices
  - Add profiling capabilities to identify performance bottlenecks

- UX and UX-only improvements
  - A settings panel is available to tweak speed, count, size, and background at runtime
  - Keyboard toggle: press P to pause/resume the animation
  - LocalStorage-backed auto-load toggle to restore the last session (via config file)
  - Add tooltips to explain each setting
  - Add a reset to defaults button in settings

## Example agent prompts

- Add command-line configuration
  - Prompt: "Add command-line flags to override configuration values (e.g., --count, --speed, --width, --height) with sane defaults."
- Make canvas responsive
  - Prompt: "Add a resize handler that scales the canvas to fit the window while preserving aspect ratio and current content style."
- Use devicePixelRatio for crisp rendering
  - Prompt: "Ensure the Ebiten rendering respects devicePixelRatio so the logos remain crisp on high-DPI displays."
- Asset failover
  - Prompt: "If any image fails to load, fall back to a simple vector logo or a placeholder shape and log a descriptive warning."
- Runtime tuning via settings
  - Prompt: "Ensure all runtime parameters can be adjusted via the settings window and are persisted correctly."
- Persistence
  - Prompt: "Verify that speed, count, width, height, background type, randomize mode, and auto-load settings are persisted between sessions."
- Pause/resume
  - Prompt: "Verify that the P key correctly pauses and resumes the animation with appropriate visual feedback."
- System tray integration
  - Prompt: "Add system tray icon with menu options to show/hide window, pause/resume, and quit."

## How to structure changes (agent guidance)

- Make small, isolated edits:
  - Prefer editing specific internal packages to adjust behavior
  - Prefer adding new files only if it significantly improves readability or tests
- Validate with manual checks:
  - Run the application locally and verify it builds and runs without errors
  - Try edge cases: very high/low counts, large window sizes, missing assets
  - Test on different platforms if possible (Linux, Windows, macOS)
- Documentation and traceability:
  - Update this AGENTS.md with any new agent-driven changes
  - Add brief notes on any notable trade-offs (e.g., performance vs. visual density)

## Nested AGENTS.md

If you add subpackages under `flying-pngs-go/` or elsewhere, you can place an AGENTS.md there as well. Agents will generally read the nearest AGENTS.md in the directory tree, and the closest one to a modified file should take precedence.

## PR and quality guidelines

- Before committing, verify the application builds and runs without errors
- If you modify assets or import paths, ensure the asset paths remain valid
- Keep changes focused and small; prefer incremental improvements
- Document any non-obvious decisions in this AGENTS.md

## Known considerations

- This project does not have automated tests at the GUI level; visual/manual checks are the primary verification method
- Unit tests exist for configuration and utility functions
- The codebase is designed to be approachable, with clear separation of concerns

## References and related tooling

- Environment is desktop-based; requires Go 1.24+ and development libraries for X11 (on Linux)
- If you adopt agent tooling in your workflow, configure it to read this AGENTS.md as part of the project context
- Required system libraries on Debian/Ubuntu: libx11-dev, libxrandr-dev, libxcursor-dev, libxinerama-dev, libxi-dev