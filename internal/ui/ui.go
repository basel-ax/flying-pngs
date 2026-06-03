package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/basel-ax/flying-pngs/internal/animation"
	"github.com/basel-ax/flying-pngs/internal/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenSettings = iota
	screenAnimation
)

// settingRow describes one configurable parameter
type settingRow struct {
	label   string
	kind    string // "int", "float", "bool", "enum"
	min     float64
	max     float64
	step    float64
	options []string // for "enum"
}

var settingsDef = []settingRow{
	{label: "Format", kind: "enum", options: []string{"png", "svg"}},
	{label: "Collection", kind: "enum", options: nil}, // populated at runtime based on format
	{label: "Window Count", kind: "int", min: 50, max: 1000, step: 50},
	{label: "Speed", kind: "float", min: 1.0, max: 20.0, step: 0.5},
	{label: "Width", kind: "int", min: 320, max: 1920, step: 80},
	{label: "Height", kind: "int", min: 240, max: 1080, step: 60},
	{label: "Background", kind: "enum", options: []string{"black", "white", "transparent"}},
	{label: "Randomize Mode", kind: "bool"},
	{label: "Randomize Max Sec", kind: "int", min: 1, max: 300, step: 1},
	{label: "SVG Base Height", kind: "float", min: 16.0, max: 512.0, step: 8.0},
}

// Game implements ebiten.Game
type Game struct {
	cfg  *config.Config
	anim *animation.AnimationCanvas

	screen    int
	debugMode bool
	cursor    int // which setting row is highlighted
	keyRepeat int // cooldown for key repeat (frames)
	statusMsg string
}

// NewGame creates a new Game instance.
// The animation is NOT started — the user sees the settings screen first.
func NewGame(cfg *config.Config) *Game {
	// Populate collection options based on current format
	updateCollectionOptions(cfg)

	g := &Game{
		cfg:    cfg,
		anim:   animation.NewAnimationCanvas(cfg),
		screen: screenSettings,
	}
	return g
}

// updateCollectionOptions populates the collection dropdown based on format
func updateCollectionOptions(cfg *config.Config) {
	if cfg.Format == "svg" {
		settingsDef[1].options = config.DiscoverSvgCollections()
	} else {
		settingsDef[1].options = config.DiscoverCollections()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.Width, g.cfg.Height
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (g *Game) Update() error {
	switch g.screen {
	case screenSettings:
		return g.updateSettings()
	case screenAnimation:
		return g.updateAnimation()
	}
	return nil
}

func (g *Game) updateSettings() error {
	// Navigate with Up / Down
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyK) {
		g.cursor--
		if g.cursor < 0 {
			g.cursor = len(settingsDef) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		g.cursor = (g.cursor + 1) % len(settingsDef)
	}

	// Change value with Left / Right
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		dir := 1.0
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			dir = -1.0
		}
		g.adjustSetting(dir)
	}

	// Hold-to-repeat for Left/Right
	g.keyRepeat--
	if g.keyRepeat < 0 {
		g.keyRepeat = 0
	}
	if g.keyRepeat == 0 && (ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyArrowRight)) {
		dir := 1.0
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			dir = -1.0
		}
		g.adjustSetting(dir)
		g.keyRepeat = 6 // ~10 fps repeat
	}

	// Enter → start animation
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.startAnimation()
	}

	// D → toggle debug (works on settings screen too)
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debugMode = !g.debugMode
	}

	return nil
}

func (g *Game) updateAnimation() error {
	g.anim.Update()

	// P → pause/resume
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.anim.TogglePause()
	}

	// D → toggle debug overlay
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debugMode = !g.debugMode
	}

	// Escape → back to settings
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.anim.Stop()
		g.screen = screenSettings
		g.statusMsg = ""
	}

	return nil
}

// adjustSetting changes the value of the currently highlighted setting
func (g *Game) adjustSetting(dir float64) {
	row := settingsDef[g.cursor]
	switch g.cursor {
	case 0: // Format
		idx := formatIndex(g.cfg.Format)
		idx += int(dir)
		if idx < 0 {
			idx = len(row.options) - 1
		}
		if idx >= len(row.options) {
			idx = 0
		}
		g.cfg.Format = row.options[idx]
		// Update collection options and reset collection
		updateCollectionOptions(g.cfg)
		if len(settingsDef[1].options) > 0 {
			g.cfg.Collection = settingsDef[1].options[0]
		} else {
			g.cfg.Collection = ""
		}
	case 1: // Collection
		idx := collectionIndex(g.cfg.Collection)
		idx += int(dir)
		if idx < 0 {
			idx = len(row.options) - 1
		}
		if idx >= len(row.options) {
			idx = 0
		}
		if len(row.options) > 0 {
			g.cfg.Collection = row.options[idx]
		}
	case 2: // Window Count
		v := float64(g.cfg.WindowCount) + dir*row.step
		if v < row.min {
			v = row.min
		}
		if v > row.max {
			v = row.max
		}
		g.cfg.WindowCount = int(v)
	case 3: // Speed
		g.cfg.Speed += dir * row.step
		if g.cfg.Speed < row.min {
			g.cfg.Speed = row.min
		}
		if g.cfg.Speed > row.max {
			g.cfg.Speed = row.max
		}
	case 4: // Width
		v := float64(g.cfg.Width) + dir*row.step
		if v < row.min {
			v = row.min
		}
		if v > row.max {
			v = row.max
		}
		g.cfg.Width = int(v)
	case 5: // Height
		v := float64(g.cfg.Height) + dir*row.step
		if v < row.min {
			v = row.min
		}
		if v > row.max {
			v = row.max
		}
		g.cfg.Height = int(v)
	case 6: // Background
		idx := bgIndex(g.cfg.BackgroundType)
		idx += int(dir)
		if idx < 0 {
			idx = len(row.options) - 1
		}
		if idx >= len(row.options) {
			idx = 0
		}
		g.cfg.BackgroundType = row.options[idx]
	case 7: // Randomize Mode
		g.cfg.RandomizeMode = !g.cfg.RandomizeMode
	case 8: // Randomize Max N
		v := float64(g.cfg.RandomizeMaxN) + dir*row.step
		if v < row.min {
			v = row.min
		}
		if v > row.max {
			v = row.max
		}
		g.cfg.RandomizeMaxN = int(v)
	case 9: // SVG Base Height
		g.cfg.SvgBaseHeight += dir * row.step
		if g.cfg.SvgBaseHeight < row.min {
			g.cfg.SvgBaseHeight = row.min
		}
		if g.cfg.SvgBaseHeight > row.max {
			g.cfg.SvgBaseHeight = row.max
		}
	}
}

func bgIndex(bg string) int {
	for i, o := range settingsDef[6].options {
		if o == bg {
			return i
		}
	}
	return 0
}

func collectionIndex(name string) int {
	for i, o := range settingsDef[1].options {
		if o == name {
			return i
		}
	}
	return 0
}

func formatIndex(format string) int {
	for i, o := range settingsDef[0].options {
		if o == format {
			return i
		}
	}
	return 0
}

// startAnimation applies config and starts the animation
func (g *Game) startAnimation() {
	// Re-create animation canvas with possibly-updated config
	g.anim = animation.NewAnimationCanvas(g.cfg)
	g.anim.Start()
	g.screen = screenAnimation
	g.statusMsg = "Running — ESC to return to settings"

	// Persist settings
	_ = config.Save(*g.cfg)
}

// ─── Draw ────────────────────────────────────────────────────────────────────

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.screen {
	case screenSettings:
		g.drawSettings(screen)
	case screenAnimation:
		g.anim.Draw(screen)
		if g.debugMode {
			g.drawDebugOverlay(screen)
		}
	}
}

func (g *Game) drawSettings(screen *ebiten.Image) {
	// Background
	screen.Fill(&color.RGBA{20, 20, 30, 255})

	var lines []string
	lines = append(lines, "=== FLYING PNGs — Settings ===")
	lines = append(lines, "")

	for i, row := range settingsDef {
		val := g.settingValueString(i)
		marker := "  "
		if i == g.cursor {
			marker = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%-22s  %s", marker, row.label+":", val))
	}

	lines = append(lines, "")
	lines = append(lines, "Controls:")
	lines = append(lines, "  Up/Down  - select setting")
	lines = append(lines, "  Left/Right - change value")
	lines = append(lines, "  Enter/Space - start animation")
	lines = append(lines, "  D - toggle debug overlay")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Assets loaded: %d  |  Config: %dx%d",
		g.anim.AssetCount(), g.cfg.Width, g.cfg.Height))

	if g.statusMsg != "" {
		lines = append(lines, "")
		lines = append(lines, g.statusMsg)
	}

	// Draw debug info on settings screen too
	if g.debugMode {
		lines = append(lines, "")
		lines = append(lines, g.anim.DebugInfo())
	}

	ebitenutil.DebugPrint(screen, strings.Join(lines, "\n"))
}

func (g *Game) settingValueString(idx int) string {
	switch idx {
	case 0:
		return g.cfg.Format
	case 1:
		if g.cfg.Collection == "" {
			return "(none)"
		}
		return g.cfg.Collection
	case 2:
		return fmt.Sprintf("%d", g.cfg.WindowCount)
	case 3:
		return fmt.Sprintf("%.1f", g.cfg.Speed)
	case 4:
		return fmt.Sprintf("%d", g.cfg.Width)
	case 5:
		return fmt.Sprintf("%d", g.cfg.Height)
	case 6:
		return g.cfg.BackgroundType
	case 7:
		if g.cfg.RandomizeMode {
			return "ON"
		}
		return "OFF"
	case 8:
		return fmt.Sprintf("%d", g.cfg.RandomizeMaxN)
	case 9:
		return fmt.Sprintf("%.0f", g.cfg.SvgBaseHeight)
	}
	return "?"
}

func (g *Game) drawDebugOverlay(screen *ebiten.Image) {
	// Semi-transparent background for readability
	w := g.cfg.Width
	overlay := ebiten.NewImage(220, 180)
	overlay.Fill(&color.RGBA{0, 0, 0, 180})

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(w-230), 10)
	screen.DrawImage(overlay, opts)

	ebitenutil.DebugPrintAt(screen, g.anim.DebugInfo(), w-225, 15)
}
