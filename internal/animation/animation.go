package animation

import (
	"fmt"
	"image/color"
	"log"
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
	drawFrames        int // counts frames drawn for startup debug
}

// DebugInfo returns a formatted string of the current animation state for debugging
func (ac *AnimationCanvas) DebugInfo() string {
	return fmt.Sprintf(
		"[DEBUG]\nstarted=%v paused=%v\nwindows=%d imgs=%d\nspeed=%.1f baseSpeed=%.1f\nwhiteMode=%v bg=%s\nrandomize=%v slowed=%v\nnextRandzIn=%.1fs\nFPS=%.1f frames=%d",
		ac.started, ac.paused,
		len(ac.windows), len(ac.imgs),
		ac.speed, ac.baseSpeed,
		ac.whiteMode, ac.cfg.BackgroundType,
		ac.randomizeMode, ac.isSlowed,
		float64(ac.nextRandomizeTime-time.Now().UnixNano())/float64(time.Second),
		ebiten.ActualFPS(),
		ac.drawFrames,
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
		log.Printf("[ERROR] Failed to load assets: %v", err)
		panic("Failed to load assets: " + err.Error())
	}
	if len(ac.imgs) == 0 {
		log.Printf("[ERROR] No PNG images found in ./collection/%s/ directory", cfg.Collection)
		panic(fmt.Sprintf("No PNG images found in ./collection/%s/ directory — check working directory", cfg.Collection))
	}
	log.Printf("[INFO] NewAnimationCanvas: %d images loaded, %d windows, %dx%d, speed=%.1f",
		len(ac.imgs), ac.windowsNum, cfg.Width, cfg.Height, cfg.Speed)

	ac.windows = make([]*Window, 0, ac.windowsNum)
	return ac
}

// loadAssets loads all PNG assets from the assets package
func (ac *AnimationCanvas) loadAssets() error {
	ebImages, err := assets.LoadPNGAssets(".", ac.cfg.Collection)
	if err != nil {
		return fmt.Errorf("asset glob error: %w", err)
	}
	if len(ebImages) == 0 {
		return fmt.Errorf("no PNG files matched ./collection/%s/*.png", ac.cfg.Collection)
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
	log.Printf("[INFO] Animation started: %d windows created", ac.windowsNum)
}

// Stop stops the animation
func (ac *AnimationCanvas) Stop() {
	ac.started = false
}

// SetSpeed sets the animation speed
func (ac *AnimationCanvas) SetSpeed(speed float64) {
	ac.baseSpeed = speed
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
		for i := len(ac.windows); i < ac.windowsNum; i++ {
			ac.windows = append(ac.windows, NewWindow(ac.imgs, ac.tintColors, ac.darkTintColors, ac.cfg.Width, ac.cfg.Height))
		}
	} else if len(ac.windows) > ac.windowsNum {
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
		ac.isSlowed = false
		ac.speed = ac.baseSpeed
		ac.nextRandomizeTime = time.Now().UnixNano() + time.Duration(rand.Float64()*ac.randomizeN*float64(time.Second)).Nanoseconds()
	} else {
		ac.isSlowed = false
		ac.speed = ac.baseSpeed
	}
}

// SetRandomizeN sets the maximum interval for randomize mode
func (ac *AnimationCanvas) SetRandomizeN(n float64) {
	ac.randomizeN = n
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
	if ac.randomizeMode {
		if now >= ac.nextRandomizeTime {
			ac.isSlowed = !ac.isSlowed
			if ac.isSlowed {
				ac.speed = ac.baseSpeed / 5
			} else {
				ac.speed = ac.baseSpeed
			}
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
		// Dark gray — true transparency requires compositor support
		screen.Fill(&color.RGBA{30, 30, 30, 255})
	default: // "black"
		screen.Fill(&color.RGBA{0, 0, 0, 255})
	}

	if !ac.started {
		ebitenutil.DebugPrint(screen, "Set parameters and press Start")
		return
	}

	ac.drawFrames++

	// Draw all windows
	for _, w := range ac.windows {
		w.Draw(screen, ac.whiteMode, ac.cfg.Width, ac.cfg.Height)
	}

	// Show render info for first 120 frames (~2 seconds)
	if ac.drawFrames <= 120 {
		info := fmt.Sprintf("Rendering: %d windows, %d images, frame %d",
			len(ac.windows), len(ac.imgs), ac.drawFrames)
		ebitenutil.DebugPrintAt(screen, info, 10, ac.cfg.Height-30)
	}

	// Draw pause indicator if paused
	if ac.paused {
		ebitenutil.DebugPrintAt(screen, "PAUSED (P)", ac.cfg.Width-100, ac.cfg.Height-20)
	}

	// Draw randomize mode indicator if active and slowed
	if ac.randomizeMode && ac.isSlowed && !ac.paused {
		ebitenutil.DebugPrintAt(screen, "SLOW x1/5", 10, ac.cfg.Height-20)
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
	w.Reset(imgs, width, height)
	return w
}

// Reset resets the window to a new random state
func (w *Window) Reset(imgs []*ebiten.Image, width, height int) {
	// z: depth from 1 (very close) to width/2 (far away)
	w.z = rand.Float64()*float64(width)/2 + 1
	w.pz = w.z
	// Scale x/y with z so objects at every depth fill the full screen
	fl := float64(width) / 2
	scale := w.z / fl
	w.x = (rand.Float64()*2 - 1) * fl * scale
	w.y = (rand.Float64()*2 - 1) * float64(height) / 2 * scale
	w.img = imgs[rand.Intn(len(imgs))]
	w.colorIndex = rand.Intn(len(w.tintColors))
}

// Update updates the window's position based on speed
func (w *Window) Update(speed float64) {
	w.z -= speed

	if w.z < 1 {
		w.z = float64(w.width)/2 + rand.Float64()*50
		// Scale x/y with z so objects spawn spread across the full screen
		fl := float64(w.width) / 2
		scale := w.z / fl
		w.x = (rand.Float64()*2 - 1) * fl * scale
		w.y = (rand.Float64()*2 - 1) * float64(w.height) / 2 * scale
		w.pz = w.z
	}
}

// Draw draws the window to the screen using perspective projection
func (w *Window) Draw(screen *ebiten.Image, whiteMode bool, screenWidth, screenHeight int) {
	if w.z < 0.5 {
		return // safety: avoid division by near-zero
	}

	// Perspective projection: objects farther away appear closer to center and smaller
	fl := float64(screenWidth) / 2 // focal length
	scale := fl / w.z

	// Screen position (centered projection)
	sx := w.x*scale + float64(screenWidth)/2
	sy := w.y*scale + float64(screenHeight)/2

	// Image size based on depth
	targetSize := 64.0 * scale // base size 64px at scale=1
	if targetSize < 2 {
		return // too small to see
	}

	imgW := w.img.Bounds().Dx()
	imgH := w.img.Bounds().Dy()
	if imgW == 0 || imgH == 0 {
		return
	}

	// Select color palette
	var palette [][]int
	if whiteMode {
		palette = w.darkTintColors
	} else {
		palette = w.tintColors
	}
	clr := palette[w.colorIndex]

	// Draw the image scaled and tinted
	opts := &ebiten.DrawImageOptions{}
	imgScale := targetSize / float64(imgW)
	opts.GeoM.Scale(imgScale, imgScale)
	opts.GeoM.Translate(sx-targetSize/2, sy-targetSize/2)

	opts.ColorM.Scale(
		float64(clr[0])/255,
		float64(clr[1])/255,
		float64(clr[2])/255,
		1,
	)

	screen.DrawImage(w.img, opts)
}
