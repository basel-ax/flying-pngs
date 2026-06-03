package assets

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	ebitenimg "github.com/hajimehoshi/ebiten/v2"
)

// LoadPNGAssets loads all PNG assets from the png/ directory relative to basePath
// and converts them to Ebiten images. If loading fails for a file, that file is skipped.
func LoadPNGAssets(basePath string) ([]*ebitenimg.Image, error) {
	var results []*ebitenimg.Image
	pattern := filepath.Join(basePath, "png", "*.png")

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob error for %q: %w", pattern, err)
	}

	log.Printf("[INFO] Asset search: pattern=%q found=%d files", pattern, len(files))

	if len(files) == 0 {
		absPath, _ := filepath.Abs(pattern)
		return nil, fmt.Errorf("no PNG files found matching %q (abs: %s)", pattern, absPath)
	}

	for _, p := range files {
		f, err := os.Open(p)
		if err != nil {
			log.Printf("[WARN] Cannot open %s: %v", p, err)
			continue
		}
		img, format, err := image.Decode(f)
		f.Close()
		if err != nil {
			log.Printf("[WARN] Cannot decode %s: %v", p, err)
			continue
		}
		bounds := img.Bounds()
		log.Printf("[INFO] Loaded %s: format=%s size=%dx%d", p, format, bounds.Dx(), bounds.Dy())
		ebImg := ebitenimg.NewImageFromImage(img)
		results = append(results, ebImg)
	}

	log.Printf("[INFO] Asset loading complete: %d/%d images loaded", len(results), len(files))
	return results, nil
}
