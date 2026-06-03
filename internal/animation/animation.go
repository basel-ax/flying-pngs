package animation

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/basel-ax/flying-pngs/internal/assets"
	"github.com/basel-ax/flying-pngs/internal/config"
)

// AnimationCanvas holds the state for the flying PNGs animation
type AnimationCanvas struct {
	cfg               *config.Config
	windows           []*Window
	imgs              []*ebiten.Image
	speed             float64
	baseSpeed         float64
	windowsNum        int
	paused            bool
	started           bool
	whiteMode         bool
	randomizeMode     bool
	randomizeN        float64
	isSlowed          bool
	nextRandomizeTime int64
	lastUpdate        int64
	tintColors        [][]int
	darkTintColors    [][]int
}

// DebugInfo returns a formatted string of the current animation state for debugging
func (ac *AnimationCanvas) DebugInfo() string {
	assetCount := len(ac.imgs)
	return fmt.Sprintf(
		"[DEBUG]\nstarted=%v paused=%v\nwindows=%d imgs=%d\nspeed=%.1f baseSpeed=%.1f\nwhiteMode=%v randomize=%v\nisSlowed=%v\nnextRandzIn=%.1fs\nFPS=%.1f",
		ac.started, ac.paused,
		len(ac.windows), assetCount,
		ac.speed, ac.baseSpeed,
		ac.whiteMode, ac.randomizeMode,
		ac.isSlowed,
		float64(ac.nextRandomizeTime-time.Now().UnixNano())/float64(time.Second),
		ebiten.ActualFPS(),
	)
}

// IsStarted returns whether the animation has started
func (ac *AnimationCanvas) IsStarted() bool { return ac.started }

// IsPaused returns whether the animation is paused
func (ac *AnimationCanvas) IsPaused() bool { return ac.paused }

// GetConfig returns the underlying config
func (ac *AnimationCanvas) GetConfig() *config.Config { return ac.cfg }

// AssetCount returns number of loaded images
func (ac *AnimationCanvas) AssetCount() int { return len(ac.imgs) }

// Window represents a single flying window/logo
type Window struct {
	x, y, z, pz    float64
	img            *ebiten.Image
	colorIndex     int
	width, height  int
	tintColors     [][]int
	darkTintColors [][]int
}

// NewAnimationCanvas creates a new animation canvas with the given configuration
func NewAnimationCanvas(cfg *config.Config) *AnimationCanvas {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Tint colors for dark background (bright colors)
	tintColors := [][]int{
		{255, 255, 255}, // white
		{0, 174, 255},   // blue
		{255, 69, 58},   // red
		{50, 215, 75},   // green
		{255, 214, 10},  // yellow
		{191, 90, 242},  // purple
		{100, 210, 255}, // light blue
		{255, 159, 10},  // orange
		{172, 142, 104}, // tan
		{0, 199, 190},   // teal
	}

	// Dark tint colors for white background (dark colors)
	darkTintColors := [][]int{
		{0, 0, 0},      // black
		{0, 80, 180},   // dark blue
		{180, 30, 20},  // dark red
		{20, 130, 40},  // dark green
		{160, 120, 0},  // dark yellow
		{100, 40, 160}, // dark purple
		{0, 100, 160},  // dark cyan
		{180, 90, 0},   // dark orange
		{100, 70, 50},  // dark brown
		{0, 120, 110},  // dark teal
	}

	ac := &AnimationCanvas{
		cfg:            cfg,
		speed:          cfg.Speed,
		baseSpeed:      cfg.Speed,
		windowsNum:     cfg.WindowCount,
		whiteMode:      cfg.BackgroundType == "white",
		randomizeMode:  cfg.RandomizeMode,
		randomizeN:     float64(cfg.RandomizeMaxN),
		tintColors:     tintColors,
		darkTintColors: darkTintColors,
	}

	// Load assets
	if err := ac.loadAssets(); err != nil {
		// In a real app, we might want to handle this better
		// For now, we'll just panic - in production we'd show an error screen
		panic("Failed to load assets: " + err.Error())
	}

	// Initialize windows (but don't start animation yet)
	ac.windows = make([]*Window, 0, ac.windowsNum)

	return ac
}

// loadAssets loads all PNG assets from the assets package
func (ac *AnimationCanvas) loadAssets() error {
	// Load PNG assets using our assets package
	ebImages, err := assets.LoadPNGAssets(".")
	if err != nil {
		return err
	}
	ac.imgs = ebImages
	return nil
}

// Start begins the animation
func (ac *AnimationCanvas) Start() {
	if ac.started {
		return
	}
	ac.started = true
	ac.windows = make([]*Window, ac.windowsNum)
	for i := 0; i < ac.windowsNum; i++ {
		ac.windows[i] = NewWindow(ac.imgs, ac.tintColors, ac.darkTintColors, ac.cfg.Width, ac.cfg.Height)
	}
	ac.nextRandomizeTime = time.Now().UnixNano() + time.Duration(rand.Float64()*ac.randomizeN*float64(time.Second)).Nanoseconds()
	ac.lastUpdate = time.Now().UnixNano()
}

// Stop stops the animation
func (ac *AnimationCanvas) Stop() {
	ac.started = false
}

// SetSpeed sets the animation speed
func (ac *AnimationCanvas) SetSpeed(speed float64) {
	ac.baseSpeed = speed
	// Only apply directly to speed when not currently slowed
	if ac.isSlowed && ac.randomizeMode {
		ac.speed = speed / 5
	} else {
		ac.speed = speed
	}
}

// SetCount sets the number of windows
func (ac *AnimationCanvas) SetCount(count int) {
	ac.windowsNum = count
	if len(ac.windows) < ac.windowsNum {
		// Need to add more windows
		for i := len(ac.windows); i < ac.windowsNum; i++ {
			ac.windows = append(ac.windows, NewWindow(ac.imgs, ac.tintColors, ac.darkTintColors, ac.cfg.Width, ac.cfg.Height))
		}
	} else if len(ac.windows) > ac.windowsNum {
		// Need to remove windows
		ac.windows = ac.windows[:ac.windowsNum]
	}
}

// SetCanvasSize updates the canvas size
func (ac *AnimationCanvas) SetCanvasSize(width, height int) {
	ac.cfg.Width = width
	ac.cfg.Height = height
}

// SetRandomizeMode enables or disables randomize mode
func (ac *AnimationCanvas) SetRandomizeMode(enabled bool) {
	ac.randomizeMode = enabled
	if ac.randomizeMode {
		// Start fresh: not slowed, schedule first event
		ac.isSlowed = false
		ac.speed = ac.baseSpeed
		ac.nextRandomizeTime = time.Now().UnixNano() + time.Duration(rand.Float64()*ac.randomizeN*float64(time.Second)).Nanoseconds()
	} else {
		// Restore full speed when mode is turned off
		ac.isSlowed = false
		ac.speed = ac.baseSpeed
	}
}

// SetRandomizeN sets the maximum interval for randomize mode
func (ac *AnimationCanvas) SetRandomizeN(n float64) {
	ac.randomizeN = n
	// Re-schedule so the new N takes effect immediately for the next tick
	ac.nextRandomizeTime = time.Now().UnixNano() + time.Duration(rand.Float64()*ac.randomizeN*float64(time.Second)).Nanoseconds()
}

// SetWhiteMode enables or disables white mode (inverts colors)
func (ac *AnimationCanvas) SetWhiteMode(enabled bool) {
	ac.whiteMode = enabled
}

// TogglePause pauses or resumes the animation
func (ac *AnimationCanvas) TogglePause() {
	ac.paused = !ac.paused
}

// Update updates the animation state
func (ac *AnimationCanvas) Update() {
	if !ac.started || ac.paused {
		return
	}

	now := time.Now().UnixNano()
	ac.lastUpdate = now

	// Update randomize mode
	if ac.randomizeMode && !ac.paused {
		if now >= ac.nextRandomizeTime {
			// Toggle slow/normal
			ac.isSlowed = !ac.isSlowed
			if ac.isSlowed {
				ac.speed = ac.baseSpeed / 5
			} else {
				ac.speed = ac.baseSpeed
			}
			// Schedule next toggle at a new random interval within [0, N] seconds
			ac.nextRandomizeTime = now + time.Duration(rand.Float64()*ac.randomizeN*float64(time.Second)).Nanoseconds()
		}
	}

	// Update all windows
	for _, w := range ac.windows {
		w.Update(ac.speed)
	}
}

// Draw draws the animation to the screen
func (ac *AnimationCanvas) Draw(screen *ebiten.Image) {
	// Clear screen with appropriate background color
	switch ac.cfg.BackgroundType {
	case "white":
		screen.Fill(&color.RGBA{255, 255, 255, 255})
	case "transparent":
		screen.Fill(&color.RGBA{0, 0, 0, 0})
	default: // "black"
		screen.Fill(&color.RGBA{0, 0, 0, 255})
	}

	if !ac.started {
		// Draw "Set parameters and press Start" message
		ebitenutil.DebugPrint(screen, "Set parameters and press Start")
		return
	}

	// Draw all windows
	for _, w := range ac.windows {
		w.Draw(screen, ac.whiteMode, ac.cfg.Width, ac.cfg.Height)
	}

	// Draw pause indicator if paused
	if ac.paused {
		ebitenutil.DebugPrintAt(screen, "PAUSED (P)", ac.cfg.Width-80, ac.cfg.Height-20)
	}

	// Draw randomize mode indicator if active and slowed
	if ac.randomizeMode && ac.isSlowed && !ac.paused {
		ebitenutil.DebugPrintAt(screen, "SLOW ×1/5", 10, ac.cfg.Height-20)
	}
}

// NewWindow creates a new window with random initial values
func NewWindow(imgs []*ebiten.Image, tintColors [][]int, darkTintColors [][]int, width, height int) *Window {
	w := &Window{
		tintColors:     tintColors,
		darkTintColors: darkTintColors,
		width:          width,
		height:         height,
	}
	w.Reset(imgs, tintColors, darkTintColors, width, height)
	return w
}

// Reset resets the window to a new random state
func (w *Window) Reset(imgs []*ebiten.Image, tintColors [][]int, darkTintColors [][]int, width, height int) {
	w.x = rand.Float64()*float64(-width) - rand.Float64()*float64(width)
	w.y = rand.Float64()*float64(-height) - rand.Float64()*float64(height)
	w.z = rand.Float64() * float64(width)
	w.pz = w.z
	w.img = imgs[rand.Intn(len(imgs))]
	w.colorIndex = rand.Intn(len(tintColors))
}

// Update updates the window's position based on speed
func (w *Window) Update(speed float64) {
	w.z -= speed

	if w.z < 1 {
		w.z = float64(w.width) / 2
		w.x = rand.Float64()*float64(-w.width) - rand.Float64()*float64(w.width)
		w.y = rand.Float64()*float64(-w.height) - rand.Float64()*float64(w.height)
		w.pz = w.z
	}
}

// Draw draws the window to the screen
func (w *Window) Draw(screen *ebiten.Image, whiteMode bool, screenWidth, screenHeight int) {
	// Calculate 2D projection
	sx := (w.x/w.z)*float64(screenWidth/2) + float64(screenWidth/2)
	sy := (w.y/w.z)*float64(screenHeight/2) + float64(screenHeight/2)

	// Calculate size based on depth (closer = larger)
	r := (w.z/float64(w.width/2))*22 + 4

	// Select color palette
	var palette [][]int
	if whiteMode {
		palette = w.darkTintColors
	} else {
		palette = w.tintColors
	}
	color := palette[w.colorIndex]

	// Draw the image scaled and tinted
	opts := &ebiten.DrawImageOptions{}
	scaleX := r / float64(w.img.Bounds().Dx())
	scaleY := r / float64(w.img.Bounds().Dy())
	opts.GeoM.Scale(scaleX, scaleY)
	opts.GeoM.Translate(sx-r/2, sy-r/2)

	opts.ColorM.Scale(
		float64(color[0])/255,
		float64(color[1])/255,
		float64(color[2])/255,
		1,
	)

	screen.DrawImage(w.img, opts)
}
