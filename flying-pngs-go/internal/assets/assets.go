package assets

import (
	"image"
	_ "image/png"
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
		return nil, err
	}
	for _, p := range files {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			continue
		}
		ebImg := ebitenimg.NewImageFromImage(img)
		results = append(results, ebImg)
	}
	return results, nil
}
