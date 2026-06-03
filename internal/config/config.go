package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	// Collection: must be a valid discovered collection
	collections := DiscoverCollections()
	valid := false
	for _, c := range collections {
		if cfg.Collection == c {
			valid = true
			break
		}
	}
	if !valid && len(collections) > 0 {
		cfg.Collection = collections[0]
	}

	return nil
}

// DiscoverCollections scans the collection/ directory for subdirectories
// and returns their names sorted alphabetically.
func DiscoverCollections() []string {
	collectionDir := filepath.Join(".", "collection")
	entries, err := os.ReadDir(collectionDir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		return nil
	}
	return names
}

// CollectionPath returns the filesystem path for a given collection name.
func CollectionPath(name string) string {
	return filepath.Join(".", "collection", name)
}

// EnsureCollectionDir ensures the collection directory exists.
func EnsureCollectionDir() error {
	return os.MkdirAll(filepath.Join(".", "collection"), 0755)
}

// ValidateCollection checks if a collection name is valid (exists as a subdirectory with PNGs).
func ValidateCollection(name string) error {
	p := CollectionPath(name)
	info, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("collection %q not found: %w", name, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("collection %q is not a directory", name)
	}
	matches, _ := filepath.Glob(filepath.Join(p, "*.png"))
	if len(matches) == 0 {
		return fmt.Errorf("collection %q contains no PNG files", name)
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
