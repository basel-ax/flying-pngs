package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config holds all user-configurable settings for the flying PNGs application
type Config struct {
	WindowCount    int     `json:"windowCount"`
	Speed          float64 `json:"speed"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	BackgroundType string  `json:"backgroundType"` // "black", "white", or "transparent"
	RandomizeMode  bool    `json:"randomizeMode"`
	RandomizeMaxN  int     `json:"randomizeMaxN"` // maximum seconds between randomize events
	AutoLoadLast   bool    `json:"autoLoadLast"`  // whether to load last session on startup
	Collection     string  `json:"collection"`    // name of the collection subdirectory to use
	Format         string  `json:"format"`        // "png" or "svg"
	ImageSize      float64 `json:"imageSize"`     // base display size for flying images in pixels
}

// DefaultConfig returns the default configuration values
func DefaultConfig() Config {
	collections := DiscoverCollections()
	defaultCollection := ""
	if len(collections) > 0 {
		defaultCollection = collections[0]
	}
	return Config{
		WindowCount:    500,
		Speed:          4.0,
		Width:          1280,
		Height:         720,
		BackgroundType: "black",
		RandomizeMode:  false,
		RandomizeMaxN:  10,
		AutoLoadLast:   true,
		Collection:     defaultCollection,
		Format:         "png",
		ImageSize:      64.0,
	}
}

// Load attempts to load configuration from the user's config file
// If the file doesn't exist or is invalid, it returns the default configuration
func Load() (Config, error) {
	configPath, err := configFilePath()
	if err != nil {
		return DefaultConfig(), err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// File doesn't exist, return default config
		return DefaultConfig(), nil
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), err
	}

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	// Validate and clamp values
	if err := validateAndClamp(&cfg); err != nil {
		return DefaultConfig(), err
	}

	return cfg, nil
}

// Save saves the configuration to the user's config file
func Save(cfg Config) error {
	configPath, err := configFilePath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	// Write JSON
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// validateAndClamp ensures all config values are within valid ranges
func validateAndClamp(cfg *Config) error {
	// Window count: 50-1000
	if cfg.WindowCount < 50 {
		cfg.WindowCount = 50
	} else if cfg.WindowCount > 1000 {
		cfg.WindowCount = 1000
	}

	// Speed: 1.0-20.0
	if cfg.Speed < 1.0 {
		cfg.Speed = 1.0
	} else if cfg.Speed > 20.0 {
		cfg.Speed = 20.0
	}

	// Width: 320-1920
	if cfg.Width < 320 {
		cfg.Width = 320
	} else if cfg.Width > 1920 {
		cfg.Width = 1920
	}

	// Height: 240-1080
	if cfg.Height < 240 {
		cfg.Height = 240
	} else if cfg.Height > 1080 {
		cfg.Height = 1080
	}

	// BackgroundType: must be one of the valid types
	validBackgrounds := map[string]bool{
		"black":       true,
		"white":       true,
		"transparent": true,
	}
	if !validBackgrounds[cfg.BackgroundType] {
		cfg.BackgroundType = "black"
	}

	// RandomizeMaxN: 1-300
	if cfg.RandomizeMaxN < 1 {
		cfg.RandomizeMaxN = 1
	} else if cfg.RandomizeMaxN > 300 {
		cfg.RandomizeMaxN = 300
	}

	// Format: must be "png" or "svg"
	if cfg.Format != "png" && cfg.Format != "svg" {
		cfg.Format = "png"
	}

	// Collection: must be a valid discovered collection for the chosen format
	var collections []string
	if cfg.Format == "svg" {
		collections = DiscoverSvgCollections()
	} else {
		collections = DiscoverCollections()
	}
	valid := false
	for _, c := range collections {
		if cfg.Collection == c {
			valid = true
			break
		}
	}
	if !valid {
		if len(collections) > 0 {
			cfg.Collection = collections[0]
		} else {
			cfg.Collection = ""
		}
	}

	// ImageSize: 8-256 pixels
	if cfg.ImageSize < 8 {
		cfg.ImageSize = 8
	} else if cfg.ImageSize > 256 {
		cfg.ImageSize = 256
	}

	return nil
}

// DiscoverCollections scans the collection/ directory for subdirectories
// that contain at least one PNG file and returns their names sorted alphabetically.
func DiscoverCollections() []string {
	return discoverCollectionsIn("collection", "*.png")
}

// DiscoverSvgCollections scans the svg_collection/ directory for subdirectories
// that contain at least one SVG file and returns their names sorted alphabetically.
func DiscoverSvgCollections() []string {
	return discoverCollectionsIn("svg_collection", "*.svg")
}

// discoverCollectionsIn scans a directory for subdirectories that contain
// at least one file matching the given glob pattern (e.g. "*.png", "*.svg").
func discoverCollectionsIn(dir, filePattern string) []string {
	baseDir := filepath.Join(projectRoot(), dir)
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		matches, _ := filepath.Glob(filepath.Join(baseDir, e.Name(), filePattern))
		if len(matches) > 0 {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		return nil
	}
	return names
}

// projectRoot returns the project root directory
func projectRoot() string {
	// Try to find the project root by looking for go.mod
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		// Go up from internal/config/config.go to project root
		return filepath.Join(filepath.Dir(filename), "..", "..")
	}
	// Fallback to current directory
	return "."
}

// CollectionPath returns the filesystem path for a given collection name.
func CollectionPath(name string) string {
	return filepath.Join(projectRoot(), "collection", name)
}

// SvgCollectionPath returns the filesystem path for a given SVG collection name.
func SvgCollectionPath(name string) string {
	return filepath.Join(projectRoot(), "svg_collection", name)
}

// EnsureCollectionDir ensures the collection directory exists.
func EnsureCollectionDir() error {
	return os.MkdirAll(filepath.Join(projectRoot(), "collection"), 0755)
}

// ValidateCollection checks if a collection name is valid (exists as a subdirectory with PNGs).
func ValidateCollection(name string) error {
	return validateCollectionIn("collection", name, "*.png")
}

// ValidateSvgCollection checks if an SVG collection name is valid (exists as a subdirectory with SVGs).
func ValidateSvgCollection(name string) error {
	return validateCollectionIn("svg_collection", name, "*.svg")
}

// validateCollectionIn checks if a collection exists and contains files with the given extension
func validateCollectionIn(dir, name, pattern string) error {
	p := filepath.Join(projectRoot(), dir, name)
	info, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("collection %q not found: %w", name, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("collection %q is not a directory", name)
	}
	matches, _ := filepath.Glob(filepath.Join(p, pattern))
	if len(matches) == 0 {
		return fmt.Errorf("collection %q contains no matching files", name)
	}
	return nil
}

// configFilePath returns the path to the configuration file
func configFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".flying-pngs", "config.json"), nil
}
