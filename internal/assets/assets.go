package assets

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	ebitenimg "github.com/hajimehoshi/ebiten/v2"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// LoadPNGAssets loads all PNG assets from the collection/<collectionName>/ directory
// relative to basePath and converts them to Ebiten images.
// If loading fails for a file, that file is skipped.
func LoadPNGAssets(basePath, collectionName string) ([]*ebitenimg.Image, error) {
	var results []*ebitenimg.Image
	pattern := filepath.Join(basePath, "collection", collectionName, "*.png")

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

// LoadSVGAssets loads all SVG assets from the svg_collection/<collectionName>/ directory
// relative to basePath and converts them to Ebiten images.
// SVGs are rasterized at a high master resolution for quality; the baseHeight parameter
// controls the target display size, while the actual rasterization uses masterSize
// to ensure smooth lines at any scale.
func LoadSVGAssets(basePath, collectionName string, baseHeight float64) ([]*ebitenimg.Image, error) {
	var results []*ebitenimg.Image
	pattern := filepath.Join(basePath, "svg_collection", collectionName, "*.svg")

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob error for %q: %w", pattern, err)
	}

	log.Printf("[INFO] SVG asset search: pattern=%q found=%d files", pattern, len(files))

	if len(files) == 0 {
		absPath, _ := filepath.Abs(pattern)
		return nil, fmt.Errorf("no SVG files found matching %q (abs: %s)", pattern, absPath)
	}

	// Rasterize at a high resolution for quality. The baseHeight controls display size,
	// but we render much larger so Ebiten scales down (preserving smooth vector lines).
	masterSize := baseHeight * 8 // 8x oversampling for crisp rendering
	if masterSize < 512 {
		masterSize = 512
	}

	for _, p := range files {
		img, err := loadSVGAsImage(p, masterSize)
		if err != nil {
			log.Printf("[WARN] Cannot load SVG %s: %v", p, err)
			continue
		}
		results = append(results, img)
	}

	log.Printf("[INFO] SVG asset loading complete: %d/%d images loaded (master size: %.0fpx)", len(results), len(files), masterSize)
	return results, nil
}

// loadSVGAsImage loads a single SVG file and rasterizes it to an Ebiten image
func loadSVGAsImage(path string, baseHeight float64) (*ebitenimg.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	icon, err := oksvg.ReadIconStream(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SVG: %w", err)
	}

	// Get SVG dimensions
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	if w == 0 || h == 0 {
		w, h = int(baseHeight), int(baseHeight)
	}

	// Calculate scale to match baseHeight
	scale := baseHeight / float64(h)
	targetW := int(float64(w) * scale)
	targetH := int(baseHeight)

	if targetW == 0 || targetH == 0 {
		return nil, fmt.Errorf("invalid dimensions: %dx%d", targetW, targetH)
	}

	// Create RGBA image
	img := image.NewRGBA(image.Rect(0, 0, targetW, targetH))

	// Rasterize SVG
	scanner := rasterx.NewScannerGV(targetW, targetH, img, img.Bounds())
	raster := rasterx.NewDasher(targetW, targetH, scanner)

	icon.SetTarget(0, 0, float64(targetW), float64(targetH))
	icon.Draw(raster, 1.0)

	log.Printf("[INFO] Loaded SVG %s: size=%dx%d", path, targetW, targetH)
	return ebitenimg.NewImageFromImage(img), nil
}
