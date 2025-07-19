package render

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"

	"github.com/sean/hex-map/pkg/hex"
	"github.com/sean/hex-map/pkg/terrain"
)

// RenderConfig controls visualization output
type RenderConfig struct {
	Width       int           // Image width in pixels
	Height      int           // Image height in pixels
	HexSize     float64       // Hex radius in pixels
	Layers      []RenderLayer // Active rendering layers
	ColorScheme ColorScheme   // Color mapping scheme
	ShowDebug   bool          // Enable debug overlays
	Quality     int           // JPEG quality (1-100)
}

// RenderLayer defines what to visualize
type RenderLayer int

const (
	LayerElevation RenderLayer = iota
	LayerWater
	LayerHillshade
	LayerDebugCoords
	LayerDebugNeighbors
	LayerValidation
)

// ColorScheme defines color mapping approaches
type ColorScheme int

const (
	SchemeElevation ColorScheme = iota // Terrain-style elevation colors
	SchemeRealistic                    // Earth-like realistic colors
	SchemeDebug                        // High-contrast debug colors
	SchemeGrayscale                    // Grayscale for scientific analysis
)

// HexRenderer is the main rendering engine
type HexRenderer struct {
	config RenderConfig
	grid   *hex.Grid
	canvas *image.RGBA
	bounds image.Rectangle
}

// RenderMetadata contains information embedded in exported images
type RenderMetadata struct {
	Generator    string               `json:"generator"`
	Timestamp    string               `json:"timestamp"`
	WorldSeed    int64                `json:"world_seed"`
	Stage        string               `json:"generation_stage"`
	ViewConfig   RenderConfig         `json:"view_config"`
	TerrainStats terrain.TerrainStats `json:"terrain_stats"`
	QualityScore float64              `json:"quality_score"`
	KnownIssues  []string             `json:"known_issues"`
}

// NewHexRenderer creates a new renderer for hex grid
func NewHexRenderer(grid *hex.Grid, config RenderConfig) *HexRenderer {
	bounds := image.Rect(0, 0, config.Width, config.Height)
	canvas := image.NewRGBA(bounds)

	// Initialize with a neutral background color
	draw.Draw(canvas, bounds, &image.Uniform{color.RGBA{240, 248, 255, 255}}, image.Point{}, draw.Src)

	return &HexRenderer{
		config: config,
		grid:   grid,
		canvas: canvas,
		bounds: bounds,
	}
}

// RenderTerrain renders terrain data to image
func (r *HexRenderer) RenderTerrain(tiles []*terrain.HexTile) (*image.RGBA, error) {
	// Clear canvas
	draw.Draw(r.canvas, r.bounds, &image.Uniform{color.RGBA{240, 248, 255, 255}}, image.Point{}, draw.Src)

	// Render each active layer
	for _, layer := range r.config.Layers {
		err := r.RenderLayer(layer, tiles)
		if err != nil {
			return nil, fmt.Errorf("failed to render layer %v: %w", layer, err)
		}
	}

	return r.canvas, nil
}

// RenderLayer renders a single layer
func (r *HexRenderer) RenderLayer(layer RenderLayer, tiles []*terrain.HexTile) error {
	switch layer {
	case LayerElevation:
		return r.renderElevationLayer(tiles)
	case LayerWater:
		return r.renderWaterLayer(tiles)
	case LayerDebugCoords:
		return r.renderDebugCoords()
	default:
		return fmt.Errorf("unsupported layer: %v", layer)
	}
}

// renderElevationLayer renders elevation data as colored hexes
func (r *HexRenderer) renderElevationLayer(tiles []*terrain.HexTile) error {
	for _, tile := range tiles {
		if tile == nil {
			continue
		}

		// Get tile color based on elevation
		tileColor := r.MapElevationToColor(tile.Elevation, r.config.ColorScheme)

		// Render the hex
		r.renderHex(tile.Coordinates, tileColor)
	}
	return nil
}

// renderWaterLayer renders water-specific features
func (r *HexRenderer) renderWaterLayer(tiles []*terrain.HexTile) error {
	for _, tile := range tiles {
		if tile == nil || tile.IsLand {
			continue
		}

		// Use a blue gradient based on depth
		depth := tile.GetDepth(0.0)
		intensity := math.Min(depth/1000.0, 1.0) // Normalize to 1km depth
		blue := uint8(50 + intensity*150)        // Range: 50-200

		waterColor := color.RGBA{0, 100, blue, 200} // Semi-transparent water
		r.renderHex(tile.Coordinates, waterColor)
	}
	return nil
}

// renderDebugCoords renders coordinate labels on hexes
func (r *HexRenderer) renderDebugCoords() error {
	coords := r.grid.AllCoords()

	for _, coord := range coords {
		x, y := r.hexToPixel(coord)

		// Draw a small marker for each hex center
		r.setPixelSafe(int(x), int(y), color.RGBA{255, 0, 0, 255})
		r.setPixelSafe(int(x+1), int(y), color.RGBA{255, 0, 0, 255})
		r.setPixelSafe(int(x), int(y+1), color.RGBA{255, 0, 0, 255})
		r.setPixelSafe(int(x-1), int(y), color.RGBA{255, 0, 0, 255})
		r.setPixelSafe(int(x), int(y-1), color.RGBA{255, 0, 0, 255})
	}

	return nil
}

// renderHex renders a single hex at the given coordinate with the specified color
func (r *HexRenderer) renderHex(coord hex.AxialCoord, hexColor color.RGBA) {
	centerX, centerY := r.hexToPixel(coord)
	size := r.config.HexSize

	// Generate hex vertices
	vertices := make([][2]float64, 6)
	for i := 0; i < 6; i++ {
		angle := math.Pi / 3.0 * float64(i) // 60 degrees per vertex
		x := centerX + size*math.Cos(angle)
		y := centerY + size*math.Sin(angle)
		vertices[i] = [2]float64{x, y}
	}

	// Simple hex fill using scanline approach
	minY := math.MaxFloat64
	maxY := -math.MaxFloat64
	for _, v := range vertices {
		if v[1] < minY {
			minY = v[1]
		}
		if v[1] > maxY {
			maxY = v[1]
		}
	}

	// Fill hex with solid color (simplified polygon fill)
	for y := int(minY); y <= int(maxY); y++ {
		if r.pointInHex(centerX, centerY, size, centerX, float64(y)) {
			for x := int(centerX - size); x <= int(centerX+size); x++ {
				if r.pointInHex(centerX, centerY, size, float64(x), float64(y)) {
					r.setPixelSafe(x, y, hexColor)
				}
			}
		}
	}
}

// pointInHex checks if a point is inside a hex
func (r *HexRenderer) pointInHex(hexX, hexY, hexSize, pointX, pointY float64) bool {
	dx := math.Abs(pointX - hexX)
	dy := math.Abs(pointY - hexY)

	// Simple approximation using a circle for now
	dist := math.Sqrt(dx*dx + dy*dy)
	return dist <= hexSize*0.9 // Slightly smaller to avoid overlap
}

// setPixelSafe safely sets a pixel color with bounds checking
func (r *HexRenderer) setPixelSafe(x, y int, c color.RGBA) {
	if x >= 0 && x < r.config.Width && y >= 0 && y < r.config.Height {
		r.canvas.Set(x, y, c)
	}
}

// MapElevationToColor applies color mapping to elevation data
func (r *HexRenderer) MapElevationToColor(elevation float64, scheme ColorScheme) color.RGBA {
	var colorMap ElevationColorMap

	switch scheme {
	case SchemeElevation:
		colorMap = TerrainColorScheme()
	case SchemeRealistic:
		colorMap = RealisticEarthScheme()
	case SchemeDebug:
		colorMap = DebugColorScheme()
	case SchemeGrayscale:
		// Convert elevation to grayscale
		normalized := (elevation + 5000.0) / 15000.0 // Normalize to [0,1]
		if normalized < 0 {
			normalized = 0
		}
		if normalized > 1 {
			normalized = 1
		}
		gray := uint8(normalized * 255)
		return color.RGBA{gray, gray, gray, 255}
	default:
		colorMap = TerrainColorScheme()
	}

	return ElevationToColor(elevation, colorMap)
}

// hexToPixel converts hex coordinate to pixel coordinate
func (r *HexRenderer) hexToPixel(coord hex.AxialCoord) (float64, float64) {
	// Use standard flat-top hex to pixel conversion
	size := r.config.HexSize

	// Flat-top hex conversion
	x := size * (3.0 / 2.0 * float64(coord.Q))
	y := size * (math.Sqrt(3.0)/2.0*float64(coord.Q) + math.Sqrt(3.0)*float64(coord.R))

	// Center in image and add offset
	centerX := float64(r.config.Width) / 2.0
	centerY := float64(r.config.Height) / 2.0

	return centerX + x, centerY + y
}

// ExportJPEG exports the rendered image as JPEG
func (r *HexRenderer) ExportJPEG(filename string, quality int) error {
	if quality < 1 || quality > 100 {
		return fmt.Errorf("JPEG quality must be between 1 and 100, got %d", quality)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	options := &jpeg.Options{Quality: quality}
	err = jpeg.Encode(file, r.canvas, options)
	if err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

// ExportPNG exports the rendered image as PNG
func (r *HexRenderer) ExportPNG(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	err = png.Encode(file, r.canvas)
	if err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// ExportJPEGWithMetadata exports JPEG with embedded metadata
func (r *HexRenderer) ExportJPEGWithMetadata(filename string, metadata RenderMetadata) error {
	// For now, just export the JPEG (metadata embedding to be implemented)
	return r.ExportJPEG(filename, r.config.Quality)
}

// ExportPNGWithMetadata exports PNG with embedded metadata
func (r *HexRenderer) ExportPNGWithMetadata(filename string, metadata RenderMetadata) error {
	// For now, just export the PNG (metadata embedding to be implemented)
	return r.ExportPNG(filename)
}
