package assets

import (
	"encoding/xml"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ebitenimg "github.com/hajimehoshi/ebiten/v2"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"

	"github.com/basel-ax/flying-pngs/internal/logger"
)

// defaultColor is used to replace CSS "currentColor" references in SVG files,
// since oksvg does not support CSS color keywords.
const defaultColor = "white"

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

	logger.Info("Asset search: pattern=%q found=%d files", pattern, len(files))

	if len(files) == 0 {
		absPath, _ := filepath.Abs(pattern)
		return nil, fmt.Errorf("no PNG files found matching %q (abs: %s)", pattern, absPath)
	}

	for _, p := range files {
		f, err := os.Open(p)
		if err != nil {
			logger.Warn("Cannot open %s: %v", p, err)
			continue
		}
		img, format, err := image.Decode(f)
		f.Close()
		if err != nil {
			logger.Warn("Cannot decode %s: %v", p, err)
			continue
		}
		bounds := img.Bounds()
		logger.Info("Loaded %s: format=%s size=%dx%d", p, format, bounds.Dx(), bounds.Dy())
		ebImg := ebitenimg.NewImageFromImage(img)
		results = append(results, ebImg)
	}

	logger.Info("Asset loading complete: %d/%d images loaded", len(results), len(files))
	return results, nil
}

// LoadSVGAssets loads all SVG assets from the svg_collection/<collectionName>/ directory
// relative to basePath and converts them to Ebiten images.
// SVGs are rasterized at a fixed high resolution (512px) for smooth vector lines;
// the actual display size is controlled separately via the Window's baseSize field.
func LoadSVGAssets(basePath, collectionName string) ([]*ebitenimg.Image, error) {
	var results []*ebitenimg.Image
	pattern := filepath.Join(basePath, "svg_collection", collectionName, "*.svg")

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob error for %q: %w", pattern, err)
	}

	logger.Info("SVG asset search: pattern=%q found=%d files", pattern, len(files))

	if len(files) == 0 {
		absPath, _ := filepath.Abs(pattern)
		return nil, fmt.Errorf("no SVG files found matching %q (abs: %s)", pattern, absPath)
	}

	// Fixed high resolution for quality - display size is controlled separately
	const masterSize = 512.0

	for _, p := range files {
		img, err := loadSVGAsImage(p, masterSize)
		if err != nil {
			logger.Warn("Cannot load SVG %s: %v", p, err)
			continue
		}
		results = append(results, img)
	}

	logger.Info("SVG asset loading complete: %d/%d images loaded (master size: %.0fpx)", len(results), len(files), masterSize)
	return results, nil
}

// loadSVGAsImage loads a single SVG file and rasterizes it to an Ebiten image.
// The SVG content is preprocessed to work around oksvg limitations:
//   - CSS "currentColor" values are replaced with a concrete color
//   - <symbol>+<use> patterns are inlined
//   - em-based dimensions are converted to pixel values
func loadSVGAsImage(path string, baseHeight float64) (*ebitenimg.Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Preprocess SVG XML to fix oksvg limitations
	content, err := preprocessSVG(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess SVG: %w", err)
	}

	icon, err := oksvg.ReadIconStream(strings.NewReader(content))
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

	logger.Info("Loaded SVG %s: size=%dx%d", path, targetW, targetH)
	return ebitenimg.NewImageFromImage(img), nil
}

// preprocessSVG applies XML-level transformations to make SVG content
// compatible with the oksvg parser. It handles:
//   - Replacing CSS "currentColor" with a concrete color
//   - Inlining <symbol> content referenced by <use> elements
//   - Converting em-based width/height to pixel values from viewBox
func preprocessSVG(content string) (string, error) {
	// Replace currentColor (case-insensitive) with a concrete color
	reColor := regexp.MustCompile(`(?i)currentColor`)
	content = reColor.ReplaceAllString(content, defaultColor)

	// Parse XML to handle <symbol>+<use> and em units
	var root xmlNode
	if err := xml.Unmarshal([]byte(content), &root); err != nil {
		// If XML parsing fails, return the color-fixed content as-is
		// and let oksvg report the error
		return content, nil
	}

	// Resolve <symbol>+<use> patterns
	resolveSymbols(&root)

	// Fix em-based dimensions using viewBox
	fixEmUnits(&root)

	out, err := xml.Marshal(root)
	if err != nil {
		return content, nil
	}
	return string(out), nil
}

// xmlNode is a generic XML element for SVG preprocessing.
type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Content  []byte     `xml:",chardata"`
	Children []xmlNode  `xml:",any"`
}

// resolveSymbols finds <symbol> elements and inlines their content
// into <use> elements that reference them via href/xlink:href.
func resolveSymbols(root *xmlNode) {
	// Collect all symbol elements by ID
	symbols := make(map[string]*xmlNode)
	collectSymbols(root, symbols)

	// Replace <use> elements with inlined symbol content
	replaceUse(root, symbols)

	// Remove <symbol> wrapper elements (their content is now inlined)
	removeSymbols(root)
}

func collectSymbols(node *xmlNode, symbols map[string]*xmlNode) {
	if node.XMLName.Local == "symbol" {
		for _, attr := range node.Attrs {
			if attr.Name.Local == "id" {
				symbols[attr.Value] = node
				break
			}
		}
	}
	for i := range node.Children {
		collectSymbols(&node.Children[i], symbols)
	}
}

func replaceUse(node *xmlNode, symbols map[string]*xmlNode) {
	for i := range node.Children {
		child := &node.Children[i]
		if child.XMLName.Local == "use" {
			// Find the referenced symbol ID
			var refID string
			for _, attr := range child.Attrs {
				if attr.Name.Local == "href" || (attr.Name.Space == "http://www.w3.org/1999/xlink" && attr.Name.Local == "href") {
					refID = strings.TrimPrefix(attr.Value, "#")
					break
				}
			}
			if sym, ok := symbols[refID]; ok {
				// Replace <use> with a <g> containing the symbol's children
				child.XMLName = xml.Name{Local: "g"}
				child.Attrs = nil
				child.Children = make([]xmlNode, len(sym.Children))
				copy(child.Children, sym.Children)
			}
		}
		replaceUse(child, symbols)
	}
}

func removeSymbols(node *xmlNode) {
	filtered := node.Children[:0]
	for _, child := range node.Children {
		c := child
		removeSymbols(&c)
		if c.XMLName.Local != "symbol" {
			filtered = append(filtered, c)
		}
	}
	node.Children = filtered
}

// fixEmUnits converts em-based width/height attributes to pixel values
// derived from the viewBox attribute.
func fixEmUnits(node *xmlNode) {
	if node.XMLName.Local == "svg" {
		var (
			viewBoxW, viewBoxH float64
			hasViewBox         bool
			widthIdx           = -1
			heightIdx          = -1
		)

		for i, attr := range node.Attrs {
			switch attr.Name.Local {
			case "viewBox":
				fmt.Sscanf(attr.Value, "%f %f %f %f", new(float64), new(float64), &viewBoxW, &viewBoxH)
				if viewBoxW > 0 && viewBoxH > 0 {
					hasViewBox = true
				}
			case "width":
				widthIdx = i
			case "height":
				heightIdx = i
			}
		}

		if hasViewBox {
			if widthIdx >= 0 && strings.HasSuffix(node.Attrs[widthIdx].Value, "em") {
				node.Attrs[widthIdx].Value = fmt.Sprintf("%g", viewBoxW)
			}
			if heightIdx >= 0 && strings.HasSuffix(node.Attrs[heightIdx].Value, "em") {
				node.Attrs[heightIdx].Value = fmt.Sprintf("%g", viewBoxH)
			}
		}
	}
	for i := range node.Children {
		fixEmUnits(&node.Children[i])
	}
}
