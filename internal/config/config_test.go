package config

import (
	"os"
	"testing"
)

// TestDefaultConfig tests that default configuration values are correct
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.WindowCount != 500 {
		t.Errorf("Expected WindowCount=500, got %d", cfg.WindowCount)
	}
	if cfg.Speed != 4.0 {
		t.Errorf("Expected Speed=4.0, got %v", cfg.Speed)
	}
	if cfg.Width != 1280 {
		t.Errorf("Expected Width=1280, got %d", cfg.Width)
	}
	if cfg.Height != 720 {
		t.Errorf("Expected Height=720, got %d", cfg.Height)
	}
	if cfg.BackgroundType != "black" {
		t.Errorf("Expected BackgroundType=black, got %s", cfg.BackgroundType)
	}
	if cfg.RandomizeMode != false {
		t.Errorf("Expected RandomizeMode=false, got %v", cfg.RandomizeMode)
	}
	if cfg.RandomizeMaxN != 10 {
		t.Errorf("Expected RandomizeMaxN=10, got %d", cfg.RandomizeMaxN)
	}
	if cfg.AutoLoadLast != true {
		t.Errorf("Expected AutoLoadLast=true, got %v", cfg.AutoLoadLast)
	}
	// Collection default depends on what's available in ./collection/
	// Just ensure it doesn't panic
	_ = cfg.Collection
}

// TestLoadDefaultConfig tests loading configuration when no file exists
func TestLoadDefaultConfig(t *testing.T) {
	// Create a temporary directory for config file
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", oldHome) }()
	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Unexpected error loading config: %v", err)
	}

	// Should return default config
	if cfg.WindowCount != 500 {
		t.Errorf("Expected WindowCount=500, got %d", cfg.WindowCount)
	}
}

// TestSaveAndLoadConfig tests saving and loading configuration
func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for config file
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", oldHome) }()
	os.Setenv("HOME", tmpDir)

	// Create a custom config
	cfg := Config{
		WindowCount:    100,
		Speed:          5.5,
		Width:          800,
		Height:         600,
		BackgroundType: "white",
		RandomizeMode:  true,
		RandomizeMaxN:  30,
		AutoLoadLast:   false,
		Collection:     "symbol",
		Format:         "svg",
		ImageSize:      128.0,
	}

	// Save the config
	if err := Save(cfg); err != nil {
		t.Fatalf("Unexpected error saving config: %v", err)
	}

	// Load the config back
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Unexpected error loading config: %v", err)
	}

	// Check that values match
	if loadedCfg.WindowCount != cfg.WindowCount {
		t.Errorf("Expected WindowCount=%d, got %d", cfg.WindowCount, loadedCfg.WindowCount)
	}
	if loadedCfg.Speed != cfg.Speed {
		t.Errorf("Expected Speed=%v, got %v", cfg.Speed, loadedCfg.Speed)
	}
	if loadedCfg.Width != cfg.Width {
		t.Errorf("Expected Width=%d, got %d", cfg.Width, loadedCfg.Width)
	}
	if loadedCfg.Height != cfg.Height {
		t.Errorf("Expected Height=%d, got %d", cfg.Height, loadedCfg.Height)
	}
	if loadedCfg.BackgroundType != cfg.BackgroundType {
		t.Errorf("Expected BackgroundType=%s, got %s", cfg.BackgroundType, loadedCfg.BackgroundType)
	}
	if loadedCfg.RandomizeMode != cfg.RandomizeMode {
		t.Errorf("Expected RandomizeMode=%v, got %v", cfg.RandomizeMode, loadedCfg.RandomizeMode)
	}
	if loadedCfg.RandomizeMaxN != cfg.RandomizeMaxN {
		t.Errorf("Expected RandomizeMaxN=%d, got %d", cfg.RandomizeMaxN, loadedCfg.RandomizeMaxN)
	}
	if loadedCfg.AutoLoadLast != cfg.AutoLoadLast {
		t.Errorf("Expected AutoLoadLast=%v, got %v", cfg.AutoLoadLast, loadedCfg.AutoLoadLast)
	}
	if loadedCfg.Collection != cfg.Collection {
		t.Errorf("Expected Collection=%s, got %s", cfg.Collection, loadedCfg.Collection)
	}
	if loadedCfg.Format != cfg.Format {
		t.Errorf("Expected Format=%s, got %s", cfg.Format, loadedCfg.Format)
	}
	if loadedCfg.ImageSize != cfg.ImageSize {
		t.Errorf("Expected ImageSize=%v, got %v", cfg.ImageSize, loadedCfg.ImageSize)
	}
}

// TestValidateAndClamp tests that configuration values are properly clamped
func TestValidateAndClamp(t *testing.T) {
	tests := []struct {
		name     string
		input    Config
		expected Config
	}{
		{
			name: "values within bounds",
			input: Config{
				WindowCount:    500,
				Speed:          10.0,
				Width:          1024,
				Height:         768,
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  10,
				AutoLoadLast:   true,
			},
			expected: Config{
				WindowCount:    500,
				Speed:          10.0,
				Width:          1024,
				Height:         768,
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  10,
				AutoLoadLast:   true,
			},
		},
		{
			name: "values below minimum",
			input: Config{
				WindowCount:    10,  // below 50
				Speed:          0.5, // below 1.0
				Width:          100, // below 320
				Height:         100, // below 240
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  10,
				AutoLoadLast:   true,
			},
			expected: Config{
				WindowCount:    50,
				Speed:          1.0,
				Width:          320,
				Height:         240,
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  10,
				AutoLoadLast:   true,
			},
		},
		{
			name: "values above maximum",
			input: Config{
				WindowCount:    2000, // above 1000
				Speed:          50.0, // above 20.0
				Width:          5000, // above 1920
				Height:         3000, // above 1080
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  10,
				AutoLoadLast:   true,
			},
			expected: Config{
				WindowCount:    1000,
				Speed:          20.0,
				Width:          1920,
				Height:         1080,
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  10,
				AutoLoadLast:   true,
			},
		},
		{
			name: "invalid background type",
			input: Config{
				WindowCount:    500,
				Speed:          4.0,
				Width:          1280,
				Height:         720,
				BackgroundType: "invalid",
				RandomizeMode:  false,
				RandomizeMaxN:  10,
				AutoLoadLast:   true,
			},
			expected: Config{
				WindowCount:    500,
				Speed:          4.0,
				Width:          1280,
				Height:         720,
				BackgroundType: "black", // should default to black
				RandomizeMode:  false,
				RandomizeMaxN:  10,
				AutoLoadLast:   true,
			},
		},
		{
			name: "invalid randomize max n",
			input: Config{
				WindowCount:    500,
				Speed:          4.0,
				Width:          1280,
				Height:         720,
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  0, // below 1
				AutoLoadLast:   true,
			},
			expected: Config{
				WindowCount:    500,
				Speed:          4.0,
				Width:          1280,
				Height:         720,
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  1,
				AutoLoadLast:   true,
			},
		},
		{
			name: "invalid randomize max n high",
			input: Config{
				WindowCount:    500,
				Speed:          4.0,
				Width:          1280,
				Height:         720,
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  500, // above 300
				AutoLoadLast:   true,
			},
			expected: Config{
				WindowCount:    500,
				Speed:          4.0,
				Width:          1280,
				Height:         720,
				BackgroundType: "black",
				RandomizeMode:  false,
				RandomizeMaxN:  300,
				AutoLoadLast:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy to avoid modifying the original
			input := tt.input
			// Validate and clamp
			if err := validateAndClamp(&input); err != nil {
				t.Fatalf("Unexpected error in validateAndClamp: %v", err)
			}

			// Check window count
			if input.WindowCount != tt.expected.WindowCount {
				t.Errorf("WindowCount: expected %d, got %d", tt.expected.WindowCount, input.WindowCount)
			}

			// Check speed
			if input.Speed != tt.expected.Speed {
				t.Errorf("Speed: expected %v, got %v", tt.expected.Speed, input.Speed)
			}

			// Check width
			if input.Width != tt.expected.Width {
				t.Errorf("Width: expected %d, got %d", tt.expected.Width, input.Width)
			}

			// Check height
			if input.Height != tt.expected.Height {
				t.Errorf("Height: expected %d, got %d", tt.expected.Height, input.Height)
			}

			// Check background type
			if input.BackgroundType != tt.expected.BackgroundType {
				t.Errorf("BackgroundType: expected %s, got %s", tt.expected.BackgroundType, input.BackgroundType)
			}

			// Check randomize max n
			if input.RandomizeMaxN != tt.expected.RandomizeMaxN {
				t.Errorf("RandomizeMaxN: expected %d, got %d", tt.expected.RandomizeMaxN, input.RandomizeMaxN)
			}

			// Check auto load last (should remain unchanged)
			if input.AutoLoadLast != tt.expected.AutoLoadLast {
				t.Errorf("AutoLoadLast: expected %v, got %v", tt.expected.AutoLoadLast, input.AutoLoadLast)
			}

			// Check randomize mode (should remain unchanged)
			if input.RandomizeMode != tt.expected.RandomizeMode {
				t.Errorf("RandomizeMode: expected %v, got %v", tt.expected.RandomizeMode, input.RandomizeMode)
			}
		})
	}
}
