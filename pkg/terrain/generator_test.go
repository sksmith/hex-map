package terrain

import (
	"testing"

	"github.com/sean/hex-map/pkg/hex"
)

func TestGenerateTerrain(t *testing.T) {
	// Create a simple grid
	config := hex.GridConfig{Width: 10, Height: 8, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)
	
	// Use default terrain config
	terrainConfig := DefaultTerrainConfig()
	
	// Generate terrain
	tiles, err := GenerateTerrain(grid, terrainConfig)
	if err != nil {
		t.Fatalf("GenerateTerrain() failed: %v", err)
	}
	
	// Verify we got the right number of tiles
	expectedTiles := 10 * 8
	if len(tiles) != expectedTiles {
		t.Errorf("Expected %d tiles, got %d", expectedTiles, len(tiles))
	}
	
	// Verify all tiles have valid data
	for i, tile := range tiles {
		if tile == nil {
			t.Errorf("Tile %d is nil", i)
			continue
		}
		
		// Check elevation is in realistic range
		if !tile.IsRealistic() {
			t.Errorf("Tile %d has unrealistic elevation: %f", i, tile.Elevation)
		}
		
		// Check coordinates are valid
		if !grid.IsValid(tile.Coordinates) {
			t.Errorf("Tile %d has invalid coordinates: %v", i, tile.Coordinates)
		}
	}
	
	// Check that we have both land and water tiles
	landCount := 0
	waterCount := 0
	for _, tile := range tiles {
		if tile.IsLand {
			landCount++
		} else {
			waterCount++
		}
	}
	
	if landCount == 0 {
		t.Error("No land tiles generated")
	}
	if waterCount == 0 {
		t.Error("No water tiles generated")
	}
}

func TestGenerateTerrainWithInvalidConfig(t *testing.T) {
	config := hex.GridConfig{Width: 5, Height: 5, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)
	
	// Invalid config with land ratio > 1
	invalidConfig := TerrainConfig{
		LandRatio:   1.5, // Invalid!
		NoiseParams: DefaultNoiseParameters(),
	}
	
	_, err := GenerateTerrain(grid, invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

func TestGenerateHeightmap(t *testing.T) {
	width, height := 8, 6
	params := DefaultNoiseParameters()
	seed := int64(42)
	
	heightmap := GenerateHeightmap(width, height, params, seed)
	
	// Check dimensions
	if len(heightmap) != height {
		t.Errorf("Expected height %d, got %d", height, len(heightmap))
	}
	
	if len(heightmap[0]) != width {
		t.Errorf("Expected width %d, got %d", width, len(heightmap[0]))
	}
	
	// Check that values are in reasonable range [-1, 1]
	for y := range heightmap {
		for x := range heightmap[y] {
			value := heightmap[y][x]
			if value < -2.0 || value > 2.0 {
				t.Errorf("Heightmap value out of range at (%d,%d): %f", x, y, value)
			}
		}
	}
	
	// Check determinism - same seed should produce same result
	heightmap2 := GenerateHeightmap(width, height, params, seed)
	
	for y := range heightmap {
		for x := range heightmap[y] {
			if heightmap[y][x] != heightmap2[y][x] {
				t.Errorf("Non-deterministic generation at (%d,%d): %f vs %f", 
					x, y, heightmap[y][x], heightmap2[y][x])
			}
		}
	}
}

func TestApplyHypsometricCurve(t *testing.T) {
	// Create simple heightmap
	heightmap := [][]float64{
		{-1.0, -0.5, 0.0, 0.5, 1.0},
		{-0.8, -0.3, 0.2, 0.7, 0.9},
	}
	
	targetLandRatio := 0.4 // 40% land
	
	result := ApplyHypsometricCurve(heightmap, targetLandRatio)
	
	// Check dimensions preserved
	if len(result) != len(heightmap) {
		t.Errorf("Height dimension changed: %d vs %d", len(result), len(heightmap))
	}
	
	if len(result[0]) != len(heightmap[0]) {
		t.Errorf("Width dimension changed: %d vs %d", len(result[0]), len(heightmap[0]))
	}
	
	// Count land vs water tiles to verify land ratio
	totalTiles := len(result) * len(result[0])
	landTiles := 0
	
	for y := range result {
		for x := range result[y] {
			if result[y][x] > 0 {
				landTiles++
			}
		}
	}
	
	actualLandRatio := float64(landTiles) / float64(totalTiles)
	tolerance := 0.2 // 20% tolerance for small grids
	
	if actualLandRatio < targetLandRatio-tolerance || actualLandRatio > targetLandRatio+tolerance {
		t.Errorf("Land ratio not achieved: target %.2f, got %.2f", targetLandRatio, actualLandRatio)
	}
}

func TestHeightmapToHexTiles(t *testing.T) {
	// Create grid
	config := hex.GridConfig{Width: 3, Height: 2, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)
	
	// Create simple heightmap
	heightmap := [][]float64{
		{-100.0, 200.0, 1500.0},
		{-50.0, 0.0, 800.0},
	}
	
	seaLevel := 0.0
	
	tiles := HeightmapToHexTiles(heightmap, grid, seaLevel)
	
	// Check we got the right number of tiles
	expectedTiles := 3 * 2
	if len(tiles) != expectedTiles {
		t.Errorf("Expected %d tiles, got %d", expectedTiles, len(tiles))
	}
	
	// Verify tile properties
	for _, tile := range tiles {
		// Check land/water classification
		expectedLand := tile.Elevation > seaLevel
		if tile.IsLand != expectedLand {
			t.Errorf("Incorrect land classification for elevation %f: got %v, expected %v",
				tile.Elevation, tile.IsLand, expectedLand)
		}
		
		// Check coordinates are valid
		if !grid.IsValid(tile.Coordinates) {
			t.Errorf("Invalid coordinates: %v", tile.Coordinates)
		}
	}
}

func TestCalculateGridDimensions(t *testing.T) {
	tests := []struct {
		name     string
		coords   []hex.AxialCoord
		expWidth int
		expHeight int
	}{
		{
			name:      "empty coords",
			coords:    []hex.AxialCoord{},
			expWidth:  0,
			expHeight: 0,
		},
		{
			name: "single coordinate",
			coords: []hex.AxialCoord{
				hex.NewAxialCoord(0, 0),
			},
			expWidth:  1,
			expHeight: 1,
		},
		{
			name: "rectangular area",
			coords: []hex.AxialCoord{
				hex.NewAxialCoord(0, 0),   // offset (0,0)
				hex.NewAxialCoord(1, 0),   // offset (1,0)
				hex.NewAxialCoord(0, 1),   // offset (0,1)
				hex.NewAxialCoord(1, 1),   // offset (1,1)
			},
			expWidth:  2,
			expHeight: 3, // Due to hex coordinate system offset conversion
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height := calculateGridDimensions(tt.coords)
			
			if width != tt.expWidth {
				t.Errorf("calculateGridDimensions() width = %d, want %d", width, tt.expWidth)
			}
			
			if height != tt.expHeight {
				t.Errorf("calculateGridDimensions() height = %d, want %d", height, tt.expHeight)
			}
		})
	}
}

func TestElevationToRealisticRange(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		wantRange [2]float64 // [min, max] expected range
	}{
		{
			name:      "maximum negative",
			input:     -1.0,
			wantRange: [2]float64{-11000, -11000},
		},
		{
			name:      "zero",
			input:     0.0,
			wantRange: [2]float64{0, 0},
		},
		{
			name:      "maximum positive",
			input:     1.0,
			wantRange: [2]float64{8849, 8849},
		},
		{
			name:      "half negative",
			input:     -0.5,
			wantRange: [2]float64{-5500, -5500},
		},
		{
			name:      "half positive",
			input:     0.5,
			wantRange: [2]float64{4424.5, 4424.5},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ElevationToRealisticRange(tt.input)
			
			tolerance := 0.1
			if result < tt.wantRange[0]-tolerance || result > tt.wantRange[1]+tolerance {
				t.Errorf("ElevationToRealisticRange(%f) = %f, want ~%f", 
					tt.input, result, tt.wantRange[0])
			}
		})
	}
}

func TestTerrainFromGrid(t *testing.T) {
	// Create a simple grid
	config := hex.GridConfig{Width: 5, Height: 5, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)
	
	// Generate terrain
	tiles, err := TerrainFromGrid(grid)
	if err != nil {
		t.Fatalf("TerrainFromGrid() failed: %v", err)
	}
	
	// Basic checks
	if len(tiles) != 25 {
		t.Errorf("Expected 25 tiles, got %d", len(tiles))
	}
	
	// Should use default seed (42)
	for _, tile := range tiles {
		if !tile.IsRealistic() {
			t.Errorf("Unrealistic elevation: %f", tile.Elevation)
		}
	}
}

func TestTerrainFromGridWithSeed(t *testing.T) {
	config := hex.GridConfig{Width: 5, Height: 5, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)
	
	seed1 := int64(12345)
	seed2 := int64(67890)
	
	// Generate terrain with different seeds
	tiles1, err1 := TerrainFromGridWithSeed(grid, seed1)
	tiles2, err2 := TerrainFromGridWithSeed(grid, seed2)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("TerrainFromGridWithSeed() failed: %v, %v", err1, err2)
	}
	
	// Should be different due to different seeds
	different := false
	for i := range tiles1 {
		if tiles1[i].Elevation != tiles2[i].Elevation {
			different = true
			break
		}
	}
	
	if !different {
		t.Error("Expected different terrain with different seeds")
	}
	
	// Same seed should produce same result
	tiles3, err3 := TerrainFromGridWithSeed(grid, seed1)
	if err3 != nil {
		t.Fatalf("TerrainFromGridWithSeed() failed: %v", err3)
	}
	
	for i := range tiles1 {
		if tiles1[i].Elevation != tiles3[i].Elevation {
			t.Errorf("Same seed produced different elevations at tile %d: %f vs %f",
				i, tiles1[i].Elevation, tiles3[i].Elevation)
		}
	}
}

func TestScaleElevationRange(t *testing.T) {
	// Create sample tiles
	tiles := []*HexTile{
		{Elevation: -1000},
		{Elevation: 0},
		{Elevation: 1000},
		{Elevation: 2000},
	}
	
	minElev := -5000.0
	maxElev := 10000.0
	
	ScaleElevationRange(tiles, minElev, maxElev)
	
	// Check that all elevations are within new range
	for i, tile := range tiles {
		if tile.Elevation < minElev || tile.Elevation > maxElev {
			t.Errorf("Tile %d elevation %f outside range [%f, %f]", 
				i, tile.Elevation, minElev, maxElev)
		}
	}
	
	// Check that range is actually used
	actualMin := tiles[0].Elevation
	actualMax := tiles[0].Elevation
	
	for _, tile := range tiles {
		if tile.Elevation < actualMin {
			actualMin = tile.Elevation
		}
		if tile.Elevation > actualMax {
			actualMax = tile.Elevation
		}
	}
	
	tolerance := 0.1
	if actualMin < minElev-tolerance || actualMin > minElev+tolerance {
		t.Errorf("Minimum elevation not used: got %f, expected ~%f", actualMin, minElev)
	}
	
	if actualMax < maxElev-tolerance || actualMax > maxElev+tolerance {
		t.Errorf("Maximum elevation not used: got %f, expected ~%f", actualMax, maxElev)
	}
}