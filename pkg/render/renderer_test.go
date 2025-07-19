package render

import (
	"image/color"
	"testing"

	"github.com/sean/hex-map/pkg/hex"
	"github.com/sean/hex-map/pkg/terrain"
)

func TestNewHexRenderer(t *testing.T) {
	// Create test grid
	config := hex.GridConfig{Width: 10, Height: 8, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	// Create renderer config
	renderConfig := RenderConfig{
		Width:       800,
		Height:      600,
		HexSize:     20.0,
		Layers:      []RenderLayer{LayerElevation},
		ColorScheme: SchemeElevation,
		ShowDebug:   false,
		Quality:     85,
	}

	// This should fail until we implement NewHexRenderer
	renderer := NewHexRenderer(grid, renderConfig)
	if renderer == nil {
		t.Error("NewHexRenderer should return a valid renderer")
	}

	if renderer.config.Width != 800 {
		t.Errorf("Expected width 800, got %d", renderer.config.Width)
	}
}

func TestRenderTerrain(t *testing.T) {
	// Create test data
	config := hex.GridConfig{Width: 5, Height: 5, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{
		Width:       400,
		Height:      400,
		HexSize:     15.0,
		Layers:      []RenderLayer{LayerElevation},
		ColorScheme: SchemeElevation,
	}

	renderer := NewHexRenderer(grid, renderConfig)

	// Create test terrain tiles
	tiles := []*terrain.HexTile{
		{
			Coordinates: hex.NewAxialCoord(0, 0),
			Elevation:   100.0,
			IsLand:      true,
		},
		{
			Coordinates: hex.NewAxialCoord(1, 0),
			Elevation:   -50.0,
			IsLand:      false,
		},
	}

	// This should fail until we implement RenderTerrain
	img, err := renderer.RenderTerrain(tiles)
	if err != nil {
		t.Errorf("RenderTerrain failed: %v", err)
	}

	if img == nil {
		t.Error("RenderTerrain should return a valid image")
	}

	if img.Bounds().Dx() != 400 || img.Bounds().Dy() != 400 {
		t.Errorf("Expected image size 400x400, got %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestMapElevationToColor(t *testing.T) {
	config := hex.GridConfig{Width: 1, Height: 1, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{ColorScheme: SchemeElevation}
	renderer := NewHexRenderer(grid, renderConfig)

	tests := []struct {
		elevation float64
		scheme    ColorScheme
		expected  color.RGBA // Expected colors based on our actual color scheme
	}{
		{-1000.0, SchemeElevation, color.RGBA{0, 99, 227, 255}}, // Interpolated between -2000 and -200
		{0.0, SchemeElevation, color.RGBA{244, 164, 96, 255}},   // Exact match - sea level tan/beach
		{500.0, SchemeElevation, color.RGBA{34, 139, 34, 255}},  // Exact match - low hills green
		{2000.0, SchemeElevation, color.RGBA{139, 69, 19, 255}}, // Exact match - mountains brown
	}

	for _, test := range tests {
		// This should fail until we implement MapElevationToColor
		result := renderer.MapElevationToColor(test.elevation, test.scheme)
		if result != test.expected {
			t.Errorf("MapElevationToColor(%f, %v) = %v, expected %v",
				test.elevation, test.scheme, result, test.expected)
		}
	}
}

func TestRenderLayer(t *testing.T) {
	config := hex.GridConfig{Width: 3, Height: 3, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{
		Width:   300,
		Height:  300,
		HexSize: 20.0,
	}
	renderer := NewHexRenderer(grid, renderConfig)

	tiles := []*terrain.HexTile{
		{Coordinates: hex.NewAxialCoord(0, 0), Elevation: 100.0, IsLand: true},
		{Coordinates: hex.NewAxialCoord(1, 0), Elevation: -50.0, IsLand: false},
	}

	// Test different layers
	layers := []RenderLayer{LayerElevation, LayerWater, LayerDebugCoords}

	for _, layer := range layers {
		// This should fail until we implement RenderLayer
		err := renderer.RenderLayer(layer, tiles)
		if err != nil {
			t.Errorf("RenderLayer failed for layer %v: %v", layer, err)
		}
	}
}

func TestExportJPEG(t *testing.T) {
	config := hex.GridConfig{Width: 2, Height: 2, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{Width: 200, Height: 200, HexSize: 25.0}
	renderer := NewHexRenderer(grid, renderConfig)

	// Create a simple test image first
	tiles := []*terrain.HexTile{
		{Coordinates: hex.NewAxialCoord(0, 0), Elevation: 100.0, IsLand: true},
	}

	// Render terrain to create the image
	_, err := renderer.RenderTerrain(tiles)
	if err != nil {
		t.Errorf("Failed to render terrain for export test: %v", err)
	}

	// This should fail until we implement ExportJPEG
	err = renderer.ExportJPEG("test_output.jpg", 85)
	if err != nil {
		t.Errorf("ExportJPEG failed: %v", err)
	}

	// Test invalid quality values
	err = renderer.ExportJPEG("test_invalid.jpg", 150)
	if err == nil {
		t.Error("ExportJPEG should fail with quality > 100")
	}
}

func TestExportPNG(t *testing.T) {
	config := hex.GridConfig{Width: 2, Height: 2, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{Width: 200, Height: 200, HexSize: 25.0}
	renderer := NewHexRenderer(grid, renderConfig)

	// Create a simple test image first
	tiles := []*terrain.HexTile{
		{Coordinates: hex.NewAxialCoord(0, 0), Elevation: 100.0, IsLand: true},
	}

	// Render terrain to create the image
	_, err := renderer.RenderTerrain(tiles)
	if err != nil {
		t.Errorf("Failed to render terrain for export test: %v", err)
	}

	// This should fail until we implement ExportPNG
	err = renderer.ExportPNG("test_output.png")
	if err != nil {
		t.Errorf("ExportPNG failed: %v", err)
	}
}

func TestHexToPixel(t *testing.T) {
	config := hex.GridConfig{Width: 10, Height: 10, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{Width: 400, Height: 400, HexSize: 20.0}
	renderer := NewHexRenderer(grid, renderConfig)

	// Test relative to center position
	centerX := float64(renderConfig.Width) / 2.0
	centerY := float64(renderConfig.Height) / 2.0

	tests := []struct {
		coord     hex.AxialCoord
		expectedX float64
		expectedY float64
	}{
		{hex.NewAxialCoord(0, 0), centerX, centerY},               // Center
		{hex.NewAxialCoord(1, 0), centerX + 30.0, centerY + 17.3}, // Hex to the right
		{hex.NewAxialCoord(0, 1), centerX, centerY + 34.6},        // Hex below
	}

	for _, test := range tests {
		x, y := renderer.hexToPixel(test.coord)
		if abs(x-test.expectedX) > 1.0 || abs(y-test.expectedY) > 1.0 {
			t.Errorf("hexToPixel(%v) = (%.1f, %.1f), expected (%.1f, %.1f)",
				test.coord, x, y, test.expectedX, test.expectedY)
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Test that renderer handles empty tile list gracefully
func TestRenderTerrainEmpty(t *testing.T) {
	config := hex.GridConfig{Width: 1, Height: 1, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{Width: 100, Height: 100, HexSize: 10.0}
	renderer := NewHexRenderer(grid, renderConfig)

	// Empty tile list should not crash
	img, err := renderer.RenderTerrain([]*terrain.HexTile{})
	if err != nil {
		t.Errorf("RenderTerrain with empty tiles failed: %v", err)
	}

	if img == nil {
		t.Error("RenderTerrain should return valid image even with empty tiles")
	}
}

// Test that renderer handles invalid coordinates gracefully
func TestRenderTerrainInvalidCoords(t *testing.T) {
	config := hex.GridConfig{Width: 2, Height: 2, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{Width: 100, Height: 100, HexSize: 10.0}
	renderer := NewHexRenderer(grid, renderConfig)

	// Tile with coordinates outside grid bounds
	tiles := []*terrain.HexTile{
		{Coordinates: hex.NewAxialCoord(10, 10), Elevation: 100.0, IsLand: true},
	}

	// Should handle gracefully without crashing
	img, err := renderer.RenderTerrain(tiles)
	if err != nil {
		t.Errorf("RenderTerrain with invalid coordinates failed: %v", err)
	}

	if img == nil {
		t.Error("RenderTerrain should return valid image even with invalid coordinates")
	}
}
