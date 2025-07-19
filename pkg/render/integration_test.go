package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sean/hex-map/pkg/hex"
	"github.com/sean/hex-map/pkg/terrain"
)

func TestEndToEndTerrainVisualization(t *testing.T) {
	// Create test directory
	testDir := "test_integration"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Generate test terrain
	gridConfig := hex.GridConfig{Width: 10, Height: 10, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	terrainConfig := terrain.TerrainConfig{
		Seed:      12345,
		SeaLevel:  0.0,
		LandRatio: 0.3,
		NoiseParams: terrain.NoiseParameters{
			Octaves:     4,
			Persistence: 0.5,
			Lacunarity:  2.0,
			Scale:       0.05,
			HurstExp:    0.85,
		},
	}

	tiles, err := terrain.GenerateTerrain(grid, terrainConfig)
	if err != nil {
		t.Fatalf("Failed to generate terrain: %v", err)
	}

	// Test all color schemes
	schemes := []ColorScheme{
		SchemeElevation,
		SchemeRealistic,
		SchemeDebug,
		SchemeGrayscale,
	}

	schemeNames := []string{"elevation", "realistic", "debug", "grayscale"}

	for i, scheme := range schemes {
		schemeName := schemeNames[i]

		// Create renderer
		renderConfig := RenderConfig{
			Width:       200,
			Height:      200,
			HexSize:     8.0,
			Layers:      []RenderLayer{LayerElevation},
			ColorScheme: scheme,
			Quality:     90,
		}

		renderer := NewHexRenderer(grid, renderConfig)

		// Render terrain
		img, err := renderer.RenderTerrain(tiles)
		if err != nil {
			t.Errorf("Failed to render terrain with scheme %s: %v", schemeName, err)
			continue
		}

		if img == nil {
			t.Errorf("Rendered image is nil for scheme %s", schemeName)
			continue
		}

		// Export as JPEG
		jpegFile := filepath.Join(testDir, "test_scheme_"+schemeName+"_integration.jpg")
		err = renderer.ExportJPEG(jpegFile, 85)
		if err != nil {
			t.Errorf("Failed to export JPEG for scheme %s: %v", schemeName, err)
			continue
		}

		// Verify file was created
		if _, err := os.Stat(jpegFile); os.IsNotExist(err) {
			t.Errorf("JPEG file was not created for scheme %s", schemeName)
		}

		// Export as PNG
		pngFile := filepath.Join(testDir, "test_scheme_"+schemeName+"_integration.png")
		err = renderer.ExportPNG(pngFile)
		if err != nil {
			t.Errorf("Failed to export PNG for scheme %s: %v", schemeName, err)
			continue
		}

		// Verify file was created
		if _, err := os.Stat(pngFile); os.IsNotExist(err) {
			t.Errorf("PNG file was not created for scheme %s", schemeName)
		}
	}
}

func TestLargeGridPerformance(t *testing.T) {
	// Skip in short test mode
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Create a larger grid for performance testing
	gridConfig := hex.GridConfig{Width: 100, Height: 100, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	// Generate minimal terrain data for performance testing
	tiles := make([]*terrain.HexTile, 0, 10000)
	for _, coord := range grid.AllCoords() {
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   float64(coord.Q + coord.R*10), // Simple gradient
			IsLand:      (coord.Q+coord.R)%3 != 0,
		}
		tiles = append(tiles, tile)
	}

	// Test rendering performance
	renderConfig := RenderConfig{
		Width:       800,
		Height:      800,
		HexSize:     3.0,
		Layers:      []RenderLayer{LayerElevation},
		ColorScheme: SchemeElevation,
		Quality:     85,
	}

	renderer := NewHexRenderer(grid, renderConfig)

	// Measure rendering time
	_, err := renderer.RenderTerrain(tiles)
	if err != nil {
		t.Errorf("Failed to render large grid: %v", err)
	}

	// Create test directory for performance test output
	testDir := "test_performance"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Test export performance
	err = renderer.ExportJPEG(filepath.Join(testDir, "large_grid_test.jpg"), 85)
	if err != nil {
		t.Errorf("Failed to export large grid: %v", err)
	}
}

func TestMultiLayerRendering(t *testing.T) {
	testDir := "test_multilayer"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create test grid and terrain
	gridConfig := hex.GridConfig{Width: 8, Height: 8, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	// Create mixed land/water terrain
	tiles := make([]*terrain.HexTile, 0, 64)
	for i, coord := range grid.AllCoords() {
		elevation := float64(-100 + i*50) // Range from -100 to high positive
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   elevation,
			IsLand:      elevation > 0,
		}
		tiles = append(tiles, tile)
	}

	// Test multi-layer combinations
	layerCombinations := [][]RenderLayer{
		{LayerElevation},
		{LayerElevation, LayerWater},
		{LayerElevation, LayerDebugCoords},
		{LayerElevation, LayerWater, LayerDebugCoords},
	}

	for i, layers := range layerCombinations {
		renderConfig := RenderConfig{
			Width:       150,
			Height:      150,
			HexSize:     10.0,
			Layers:      layers,
			ColorScheme: SchemeElevation,
			Quality:     90,
		}

		renderer := NewHexRenderer(grid, renderConfig)

		_, err := renderer.RenderTerrain(tiles)
		if err != nil {
			t.Errorf("Failed to render layer combination %d: %v", i, err)
			continue
		}

		filename := filepath.Join(testDir, "multilayer_test_"+string(rune('A'+i))+".png")
		err = renderer.ExportPNG(filename)
		if err != nil {
			t.Errorf("Failed to export layer combination %d: %v", i, err)
		}
	}
}

func TestWorldTopologyRendering(t *testing.T) {
	testDir := "test_world_topology"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create world topology grid (wrapping edges)
	gridConfig := hex.GridConfig{Width: 12, Height: 8, Topology: hex.TopologyWorld}
	grid := hex.NewGrid(gridConfig)

	// Generate simple terrain
	terrainConfig := terrain.DefaultTerrainConfig()
	terrainConfig.Seed = 789

	tiles, err := terrain.GenerateTerrain(grid, terrainConfig)
	if err != nil {
		t.Fatalf("Failed to generate terrain for world topology: %v", err)
	}

	// Render with world topology
	renderConfig := RenderConfig{
		Width:       300,
		Height:      200,
		HexSize:     8.0,
		Layers:      []RenderLayer{LayerElevation, LayerDebugCoords},
		ColorScheme: SchemeDebug,
		Quality:     90,
	}

	renderer := NewHexRenderer(grid, renderConfig)

	_, err = renderer.RenderTerrain(tiles)
	if err != nil {
		t.Errorf("Failed to render world topology: %v", err)
	}

	err = renderer.ExportPNG(filepath.Join(testDir, "world_topology_test.png"))
	if err != nil {
		t.Errorf("Failed to export world topology: %v", err)
	}
}

func TestExtremeElevationValues(t *testing.T) {
	testDir := "test_extreme_values"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	gridConfig := hex.GridConfig{Width: 5, Height: 5, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	// Create tiles with extreme elevation values
	extremeElevations := []float64{
		-11000.0, // Mariana Trench depth
		-5000.0,  // Deep ocean
		0.0,      // Sea level
		8849.0,   // Everest height
		10000.0,  // Above Everest
	}

	tiles := make([]*terrain.HexTile, 0, 25)
	coords := grid.AllCoords()

	for i, coord := range coords {
		elevation := extremeElevations[i%len(extremeElevations)]
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   elevation,
			IsLand:      elevation > 0,
		}
		tiles = append(tiles, tile)
	}

	// Test all color schemes with extreme values
	schemes := []ColorScheme{SchemeElevation, SchemeRealistic, SchemeDebug}
	schemeNames := []string{"elevation", "realistic", "debug"}

	for i, scheme := range schemes {
		schemeName := schemeNames[i]

		renderConfig := RenderConfig{
			Width:       200,
			Height:      200,
			HexSize:     15.0,
			Layers:      []RenderLayer{LayerElevation},
			ColorScheme: scheme,
			Quality:     95,
		}

		renderer := NewHexRenderer(grid, renderConfig)

		_, err := renderer.RenderTerrain(tiles)
		if err != nil {
			t.Errorf("Failed to render extreme values with scheme %s: %v", schemeName, err)
			continue
		}

		filename := filepath.Join(testDir, "extreme_values_scheme_"+schemeName+"_test.jpg")
		err = renderer.ExportJPEG(filename, 95)
		if err != nil {
			t.Errorf("Failed to export extreme values with scheme %s: %v", schemeName, err)
		}
	}
}

func TestMemoryUsage(t *testing.T) {
	// Test that rendering doesn't cause memory leaks
	gridConfig := hex.GridConfig{Width: 50, Height: 50, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	// Generate simple terrain
	tiles := make([]*terrain.HexTile, 0, 2500)
	for _, coord := range grid.AllCoords() {
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   float64(coord.Q + coord.R),
			IsLand:      true,
		}
		tiles = append(tiles, tile)
	}

	renderConfig := RenderConfig{
		Width:       400,
		Height:      400,
		HexSize:     4.0,
		Layers:      []RenderLayer{LayerElevation},
		ColorScheme: SchemeElevation,
		Quality:     85,
	}

	// Render multiple times to check for memory leaks
	for i := 0; i < 10; i++ {
		renderer := NewHexRenderer(grid, renderConfig)
		_, err := renderer.RenderTerrain(tiles)
		if err != nil {
			t.Errorf("Failed to render iteration %d: %v", i, err)
		}
		// Renderer should be garbage collected after this loop iteration
	}
}
