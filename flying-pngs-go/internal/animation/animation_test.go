package animation

import (
	"testing"

	"flying-pngs-go/internal/config"
)

// TestNewAnimationCanvas tests that a new animation canvas is created correctly
func TestNewAnimationCanvas(t *testing.T) {
	cfg := config.DefaultConfig()
	ac := NewAnimationCanvas(&cfg)

	if ac == nil {
		t.Error("NewAnimationCanvas returned nil")
	}

	if ac.cfg.Width != cfg.Width {
		t.Errorf("Expected width %d, got %d", cfg.Width, ac.cfg.Width)
	}
	if ac.cfg.Height != cfg.Height {
		t.Errorf("Expected height %d, got %d", cfg.Height, ac.cfg.Height)
	}
	if ac.windowsNum != cfg.WindowCount {
		t.Errorf("Expected windowsNum %d, got %d", cfg.WindowCount, ac.windowsNum)
	}
	if ac.speed != cfg.Speed {
		t.Errorf("Expected speed %v, got %v", cfg.Speed, ac.speed)
	}
	if ac.baseSpeed != cfg.Speed {
		t.Errorf("Expected baseSpeed %v, got %v", cfg.Speed, ac.baseSpeed)
	}
	if ac.whiteMode != false {
		t.Errorf("Expected whiteMode false, got %v", ac.whiteMode)
	}
	if ac.randomizeMode != cfg.RandomizeMode {
		t.Errorf("Expected randomizeMode %v, got %v", cfg.RandomizeMode, ac.randomizeMode)
	}
	if ac.randomizeN != float64(cfg.RandomizeMaxN) {
		t.Errorf("Expected randomizeN %v, got %v", float64(cfg.RandomizeMaxN), ac.randomizeN)
	}
}

// TestSetSpeed tests setting the animation speed
func TestSetSpeed(t *testing.T) {
	cfg := config.DefaultConfig()
	ac := NewAnimationCanvas(&cfg)

	ac.SetSpeed(10.0)
	if ac.baseSpeed != 10.0 {
		t.Errorf("Expected baseSpeed 10.0, got %v", ac.baseSpeed)
	}
	if ac.speed != 10.0 {
		t.Errorf("Expected speed 10.0, got %v", ac.speed)
	}

	// Test with randomize mode enabled and slowed
	ac.randomizeMode = true
	ac.isSlowed = true
	ac.SetSpeed(5.0)
	if ac.baseSpeed != 5.0 {
		t.Errorf("Expected baseSpeed 5.0, got %v", ac.baseSpeed)
	}
	if ac.speed != 1.0 { // should be 5.0 / 5 = 1.0 when slowed
		t.Errorf("Expected speed 1.0, got %v", ac.speed)
	}
}

// TestSetCount tests setting the number of windows
func TestSetCount(t *testing.T) {
	cfg := config.DefaultConfig()
	ac := NewAnimationCanvas(&cfg)

	ac.SetCount(100)
	if ac.windowsNum != 100 {
		t.Errorf("Expected windowsNum 100, got %d", ac.windowsNum)
	}
	if len(ac.windows) != 100 {
		t.Errorf("Expected windows length 100, got %d", len(ac.windows))
	}

	// Test decreasing count
	ac.SetCount(50)
	if ac.windowsNum != 50 {
		t.Errorf("Expected windowsNum 50, got %d", ac.windowsNum)
	}
	if len(ac.windows) != 50 {
		t.Errorf("Expected windows length 50, got %d", len(ac.windows))
	}
}

// TestTogglePause tests pausing and resuming the animation
func TestTogglePause(t *testing.T) {
	cfg := config.DefaultConfig()
	ac := NewAnimationCanvas(&cfg)

	if ac.paused != false {
		t.Errorf("Expected paused false, got %v", ac.paused)
	}

	ac.TogglePause()
	if ac.paused != true {
		t.Errorf("Expected paused true, got %v", ac.paused)
	}

	ac.TogglePause()
	if ac.paused != false {
		t.Errorf("Expected paused false, got %v", ac.paused)
	}
}
