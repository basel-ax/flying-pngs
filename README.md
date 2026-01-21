# flying-windows
![flying_windows](https://github.com/miresk/flying-windows/blob/main/flying-windows-new.gif)

Javascript version (p5.js) of an old Windows screensaver. I updated it so it's using new Windows logo.
The project is inspired by Starfield animation -> https://thecodingtrain.com/CodingChallenges/001-starfield.html

[DEMO](https://editor.p5js.org/miresk/sketches/DzxMND-Rw)

[Written tutorial](https://miro.substack.com/p/flying-windows-screensaver-in-javascript)

## Running

- Serve the `flying-pngs/` directory:
  - `npx serve -s flying-pngs`
  - Or `python3 -m http.server 8000` from the repo root or inside `flying-pngs`
- Open the served URL in a browser.

## Runtime controls

- A small “Agent Controls” panel appears in the bottom-right:
  - Speed (1–20)
  - Count (50–1000)
  - Width (320–1920) and Height (240–1080)
  - Reset to defaults
  - Auto-load last session (localStorage toggle)
- URL params are supported:
  - `?count=300&speed=6&width=1024&height=768`
- Default canvas size: 1280x720
- Pause/resume: press **P** to toggle

## Persistence

- Speed, count, width, and height are saved in `localStorage`.
- “Auto-load last session” controls whether stored values are applied on reload.

## Testing (manual smoke tests)

- Page loads without errors
- Agent Controls panel appears
- Adjusting speed/count/size updates the animation
- URL params apply correctly
- Pressing **P** pauses/resumes the animation
- Settings persist across reloads when auto-load is enabled


