# Flying PNGs Screensaver

A lightweight p5.js recreation of the classic flying Windows screensaver, now driven by a collection of PNG assets (religious and mystical iconography) with a minimal runtime control panel for tweaking animation behaviour.

## Features

- Pure client-side implementation using p5.js (no build pipeline required)
- 15 PNG sprites streamed from `png/` and randomly assigned to each particle
- Adjustable particle speed, count, and canvas dimensions at runtime
- LocalStorage-backed persistence with optional "auto-load last session" toggle
- Keyboard pause/resume shortcut (`P`)
- **Randomize Mode** — automatically toggles animation speed to ×1/5 at random intervals (0–N seconds), with an on-screen indicator and full persistence support
- Works out-of-the-box on any static hosting provider (e.g., Cloudflare Pages)

## Repository layout

```
flying-pngs/
├── agent-controls.js   # Runtime control panel and persistence wiring
├── index.html          # App entry point (loads p5.js, sketch.js, and control panel)
├── package.json        # npm scripts for convenient local serving
├── serve.json          # serve v14 config (public root, SPA rewrite, cache headers)
├── sketch.js           # Animation logic and runtime setup
├── png/                # PNG sprites used by the animation
├── README.md           # Project documentation (this file)
└── AGENTS.md           # Agent-facing instructions for coding assistants
```

## Local development

You must serve the project over HTTP(S); loading the HTML file directly with `file://` will trigger CORS errors when p5.js tries to load the PNG assets.

### Recommended: run from inside the `flying-pngs/` directory

```sh
cd flying-pngs

# Using the npm start script (always correct, no path confusion)
npm start

# Or directly with serve (serve.json is auto-detected)
npx serve .
```

### Alternative: run from the repository root (parent directory)

```sh
npx serve flying-pngs
# or use any static file server (Python example):
python3 -m http.server 8000 --directory flying-pngs
```

> **Why `npx serve flying-pngs` from inside the directory causes a 404**  
> `serve` resolves its path argument relative to your current working directory. If your CWD is already `flying-pngs/`, the argument `flying-pngs` resolves to `flying-pngs/flying-pngs/` — a directory that does not exist. Always use `npx serve .` (or `npm start`) when you are already inside the project directory.

Then open the printed URL (for example `http://localhost:3000/`).  
The runtime control panel appears in the bottom-right corner.

### Runtime controls

| Control | Description |
|---|---|
| **Speed slider** (1–20) | Controls how fast sprites fly toward the viewer |
| **Count slider** (50–1000) | Number of simultaneous sprites |
| **Width / Height inputs** (320–1920 / 240–1080) | Canvas dimensions |
| **Auto-load toggle** | Reads/writes all values to LocalStorage on load |
| **Reset to defaults button** | Restores all settings to their default values |
| **Randomize Mode checkbox** | Enables/disables the speed randomize mode |
| **Max interval N (s) input** (1–300) | Maximum seconds between each randomize event |
| **Keyboard `P`** | Pause/resume the animation |

#### Randomize Mode

When enabled, a timer fires at a random interval between 0 and N seconds. On each tick the animation speed toggles between normal and ×1/5 (slow), then a new random interval is scheduled. A yellow **"SLOW ×1/5"** indicator is shown in the bottom-left corner of the canvas during slow phases. The mode state and N value are persisted to LocalStorage and restored with the rest of the session when auto-load is on.

### URL parameters

All major settings can be seeded via URL query string. Values are clamped to valid ranges and take precedence over LocalStorage.

```
?count=300&speed=6&width=1024&height=768&randomizeMode=true&randomizeN=5
```

| Parameter | Type | Range | Default |
|---|---|---|---|
| `count` | integer | 50–1000 | 500 |
| `speed` | integer | 1–20 | 4 |
| `width` | integer | 320–1920 | 1280 |
| `height` | integer | 240–1080 | 720 |
| `randomizeMode` | boolean (`true`/`false`) | — | false |
| `randomizeN` | number | 1–300 | 10 |

## Deploying to Cloudflare Pages

1. **Create a new Pages project** and connect the Git repository.
2. **Build settings**
   - Build command: _leave empty_ (no build step required)
   - Build output directory: `flying-pngs`
   - Root directory: _root of the repository_ (Cloudflare Pages will publish the specified output directory)
3. **Environment variables**
   - None required (all logic is client-side)
4. **Preview & production branches**
   - The default settings work out-of-the-box; each commit to your main branch triggers a redeploy.

After deployment, the app is served from your Cloudflare Pages domain. If you prefer to host from the repository root instead of `flying-pngs/`, add a redirecting `index.html` at the root:

```html
<!doctype html>
<meta http-equiv="refresh" content="0; url=./flying-pngs/">
<script>window.location.replace("./flying-pngs/");</script>
```

## Manual smoke tests

Run these checks before shipping changes or after deploying:

1. Page loads without console errors.
2. All PNG sprites render (no broken images).
3. Control panel sliders and inputs update the animation in real time.
4. URL parameters override defaults on load.
5. Toggling "auto-load" off prevents persisted values from auto-applying after a refresh.
6. Pressing `P` pauses/resumes and shows the `PAUSED (P)` indicator.
7. Reload the page and confirm values persist only when auto-load is enabled.
8. Enable Randomize Mode — confirm the animation periodically slows to a crawl and the yellow `SLOW ×1/5` label appears, then speed restores.
9. Change N to a small value (e.g. 2 s) — confirm transitions happen more frequently.
10. Disable Randomize Mode mid-animation — confirm speed immediately returns to normal and the indicator disappears.
11. Reload with `?randomizeMode=true&randomizeN=3` in the URL — confirm the mode starts active with the correct N value.

## Contributing

- Keep changes focused and maintainable.
- Update `AGENTS.md` with any new workflows or instructions useful for AI coding assistants.
- If you add assets, place them in `flying-pngs/png/` and adjust `ASSET_FILES` in `sketch.js`.
- Prefer small, incremental pull requests with accompanying notes in this README when deployment, runtime, or persistence behaviour changes.

Enjoy the retro vibes!