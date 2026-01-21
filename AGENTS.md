# AGENTS.md

This repository contains a tiny p5.js project that renders a Windows-logo “flying” screensaver style animation. This AGENTS.md is written to help AI coding agents understand how to run, inspect, and improve the project with safe, targeted changes.

## Project overview

- Location: `flying-pngs/`
- Entry point: `flying-pngs/index.html` (loads p5.js and `flying-pngs/sketch.js`)
- Visual: 15 PNG assets loaded from the `png/` folder (anonymous-mask.png, da-vinci.png, grand-master-key.png, hacker.png, hanukkah-glass.png, jain-symbol.png, kaaba.png, lilith-symbol.png, lucifer-sigil.png, magic-map.png, mezuzah.png, praying-beads.png, seed-of-life.png, spirituality.png, star-of-david.png)
- Audio: disabled (no error.mp3 or sound usage)
- Tech stack: HTML, JavaScript (no build step), uses p5.js

## How to run

- Simple static serve (no build step):
  - npx serve -s flying-pngs
  - Or python3 -m http.server 8000 from the repository root or inside `flying-pngs`
- Open the served URL in a web browser to view the animation
- Quick check: verify the PNG assets appear, the Agent Controls panel loads, and no sound plays

## Key files and roles

- `flying-pngs/index.html` — HTML shell; loads p5.js, `sketch.js`, and `agent-controls.js`
- `flying-pngs/sketch.js` — animation logic (loads 15 assets from `png/`, pause/resume, persistence)
  - `preload()` loads PNG images from `png/`
  - `setup()` initializes the canvas and exposes runtime setters
  - `draw()` renders frames by updating and drawing all window objects
  - `keyPressed()` toggles pause/resume with P
  - `Window` class handles 3D-like projection, position, and rendering of each logo
- `flying-pngs/agent-controls.js` — runtime UI panel, URL parameter handling, and localStorage persistence
- `flying-pngs/png/` — image assets (15 PNGs)

## Agent-facing workflow

Agents should aim to improve the project with safe, incremental changes that preserve the visual intent while improving robustness, performance, or add configurability.

- Non-invasive improvements
  - Add runtime configurability via URL parameters (e.g., count, speed, size) without editing code
  - Implement a responsive canvas that adapts to window size while preserving aspect ratio
  - Respect devicePixelRatio for crisp rendering on high-DPI displays
  - Provide a simple fallback if assets fail to load (e.g., draw a basic shape instead of logos)

- Performance and robustness
  - Allow the number of windows to be capped, and adjust speed based on a parameter or performance hints
  - Add error handling around asset loading; log helpful warnings and fall back gracefully
  - Minimize draw calls or optimize math in projection to reduce CPU usage on low-end devices

- UX and UX-only improvements
  - A tiny UI overlay is available to tweak speed, count, and size at runtime
  - Keyboard toggle: press P to pause/resume the animation
  - LocalStorage-backed auto-load toggle to restore the last session

## Example agent prompts

- Add URL-configured runtime parameters
  - Prompt: "Expose runtime parameters via URL: `?count=<n>&speed=<s>` with sane defaults and no breaking changes."
- Make canvas responsive
  - Prompt: "Add a resize handler that scales the canvas to fit the window while preserving aspect ratio and current content style."
- Use devicePixelRatio for crisp rendering
  - Prompt: "Scale the canvas using `devicePixelRatio` so the logos remain crisp on high-DPI displays."
- Asset failover
  - Prompt: "If any image fails to load, fall back to a simple vector logo or a placeholder shape and log a descriptive warning."
- Runtime tuning
  - Prompt: "Introduce a minimal, non-intrusive UI or URL parameters to adjust `windowsNum`, `speed`, and canvas size at runtime."
- Persistence
  - Prompt: "Persist speed, count, width, and height in localStorage and add a toggle to auto-load the last session."
- Pause/resume
  - Prompt: "Add a keyboard toggle (P) to pause/resume animation and show a PAUSED indicator."
- Testing and verification prompts
  - Prompt: "Provide smoke tests that validate: a) page loads, b) no runtime errors on load, c) parameter changes apply, d) animation remains smooth under common browser conditions."

## How to structure changes (agent guidance)

- Make small, isolated edits:
  - Prefer editing `flying-pngs/sketch.js` to adjust behavior
  - Prefer adding new files (e.g., a tiny helper module) only if it significantly improves readability or tests
- Validate with manual checks:
  - Run the app locally in a couple of browsers and screen sizes
  - Try edge cases: very high/low counts, large window sizes, missing assets
- Documentation and traceability:
  - Update this AGENTS.md with any new agent-driven changes
  - Add brief notes on any notable trade-offs (e.g., performance vs. visual density)

## Nested AGENTS.md

If you add subpackages under `flying-pngs/` or elsewhere, you can place an AGENTS.md there as well. Agents will generally read the nearest AGENTS.md in the directory tree, and the closest one to a modified file should take precedence.

## PR and quality guidelines

- Before committing, verify the page loads and the animation runs without errors
- If you modify assets or import paths, ensure the asset paths remain valid
- Keep changes focused and small; prefer incremental improvements
- Document any non-obvious decisions in this AGENTS.md

## Known considerations

- This project does not have automated tests; visual/manual checks are the primary verification method
- The codebase is intentionally small and approachable, designed for quick agent-driven experimentation

## References and related tooling

- Environment is browser-based; no build pipeline is required
- If you adopt agent tooling in your workflow, configure it to read this AGENTS.md as part of the project context