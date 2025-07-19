package terrain

import (
	"math"
	"sort"

	"github.com/sean/hex-map/internal/noise"
	"github.com/sean/hex-map/pkg/hex"
)

// GenerateTerrain creates a complete terrain with elevation and land/water classification
func GenerateTerrain(grid *hex.Grid, config TerrainConfig) ([]*HexTile, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	
	// Get grid dimensions for heightmap generation
	coords := grid.AllCoords()
	if len(coords) == 0 {
		return nil, &TerrainError{"empty grid provided"}
	}
	
	// Determine bounding box for heightmap
	width, height := calculateGridDimensions(coords)
	
	// Generate base heightmap using multi-octave noise
	heightmap := GenerateHeightmap(width, height, config.NoiseParams, config.Seed)
	
	// Apply hypsometric curve to match Earth's elevation distribution
	heightmap = ApplyHypsometricCurve(heightmap, config.LandRatio)
	
	// Convert heightmap to hex tiles with proper coordinate mapping
	tiles := HeightmapToHexTiles(heightmap, grid, config.SeaLevel)
	
	return tiles, nil
}

// GenerateHeightmap creates a fractal heightmap using Diamond-Square algorithm
func GenerateHeightmap(width, height int, params NoiseParameters, seed int64) [][]float64 {
	return noise.MultiOctaveNoise(width, height, params.Octaves, 
		params.Persistence, params.Lacunarity, params.Scale, seed)
}

// ApplyHypsometricCurve adjusts elevation distribution to match Earth's curve
func ApplyHypsometricCurve(heightmap [][]float64, targetLandRatio float64) [][]float64 {
	if targetLandRatio <= 0 || targetLandRatio >= 1 {
		return heightmap // No adjustment needed for extreme ratios
	}
	
	// Flatten heightmap for sorting
	var elevations []float64
	for _, row := range heightmap {
		elevations = append(elevations, row...)
	}
	
	// Sort elevations to find percentile thresholds
	sort.Float64s(elevations)
	
	// Find the elevation that gives us the target land ratio
	seaLevelIndex := int(float64(len(elevations)) * (1.0 - targetLandRatio))
	if seaLevelIndex >= len(elevations) {
		seaLevelIndex = len(elevations) - 1
	}
	seaLevelThreshold := elevations[seaLevelIndex]
	
	// Apply Earth's hypsometric curve transformation
	result := make([][]float64, len(heightmap))
	for i := range result {
		result[i] = make([]float64, len(heightmap[i]))
		copy(result[i], heightmap[i])
	}
	
	// Transform elevations to match Earth's distribution
	for y := range result {
		for x := range result[y] {
			originalElev := result[y][x]
			
			if originalElev <= seaLevelThreshold {
				// Ocean depths: apply cubic curve for deep ocean basins
				ratio := originalElev / seaLevelThreshold
				if ratio < 0 {
					ratio = 0
				}
				depth := math.Pow(ratio, 3) * 6000 // Max depth ~6000m
				result[y][x] = -depth
			} else {
				// Land elevations: apply power curve for mountain peaks
				ratio := (originalElev - seaLevelThreshold) / (1.0 - seaLevelThreshold)
				if ratio > 1 {
					ratio = 1
				}
				// Power curve creates realistic mountain distribution
				height := math.Pow(ratio, 2.5) * 8800 // Max height ~8800m (Everest)
				result[y][x] = height
			}
		}
	}
	
	return result
}

// HeightmapToHexTiles converts a heightmap to hex tiles with land/water classification
func HeightmapToHexTiles(heightmap [][]float64, grid *hex.Grid, seaLevel float64) []*HexTile {
	coords := grid.AllCoords()
	tiles := make([]*HexTile, len(coords))
	
	height := len(heightmap)
	width := 0
	if height > 0 {
		width = len(heightmap[0])
	}
	
	for i, coord := range coords {
		// Map hex coordinate to heightmap indices
		col, row := coord.ToOffset()
		
		// Ensure we're within heightmap bounds
		x := col % width
		y := row % height
		if x < 0 {
			x += width
		}
		if y < 0 {
			y += height
		}
		
		elevation := heightmap[y][x]
		
		tile := &HexTile{
			Coordinates:     coord,
			Elevation:       elevation,
			DistanceToWater: 0, // Will be calculated later
		}
		
		// Classify as land or water based on sea level
		tile.ClassifyLandWater(seaLevel)
		
		tiles[i] = tile
	}
	
	return tiles
}

// calculateGridDimensions determines the bounding box for a set of coordinates
func calculateGridDimensions(coords []hex.AxialCoord) (width, height int) {
	if len(coords) == 0 {
		return 0, 0
	}
	
	minCol, maxCol := math.MaxInt32, math.MinInt32
	minRow, maxRow := math.MaxInt32, math.MinInt32
	
	for _, coord := range coords {
		col, row := coord.ToOffset()
		
		if col < minCol {
			minCol = col
		}
		if col > maxCol {
			maxCol = col
		}
		if row < minRow {
			minRow = row
		}
		if row > maxRow {
			maxRow = row
		}
	}
	
	width = maxCol - minCol + 1
	height = maxRow - minRow + 1
	
	return width, height
}

// ElevationToRealisticRange scales normalized elevation [-1,1] to Earth-like range
func ElevationToRealisticRange(normalizedElev float64) float64 {
	if normalizedElev < 0 {
		// Ocean depths: -11000m to 0m
		return normalizedElev * 11000
	} else {
		// Land heights: 0m to 8849m
		return normalizedElev * 8849
	}
}

// TerrainFromGrid generates terrain for an existing hex grid with default parameters
func TerrainFromGrid(grid *hex.Grid) ([]*HexTile, error) {
	config := DefaultTerrainConfig()
	return GenerateTerrain(grid, config)
}

// TerrainFromGridWithSeed generates terrain with a specific seed
func TerrainFromGridWithSeed(grid *hex.Grid, seed int64) ([]*HexTile, error) {
	config := DefaultTerrainConfig()
	config.Seed = seed
	return GenerateTerrain(grid, config)
}

// ScaleElevationRange scales all elevations in tiles to a specific range
func ScaleElevationRange(tiles []*HexTile, minElev, maxElev float64) {
	if len(tiles) == 0 {
		return
	}
	
	// Find current range
	currentMin := tiles[0].Elevation
	currentMax := tiles[0].Elevation
	
	for _, tile := range tiles {
		if tile.Elevation < currentMin {
			currentMin = tile.Elevation
		}
		if tile.Elevation > currentMax {
			currentMax = tile.Elevation
		}
	}
	
	currentRange := currentMax - currentMin
	if currentRange == 0 {
		return // All elevations are the same
	}
	
	targetRange := maxElev - minElev
	
	// Scale all elevations
	for _, tile := range tiles {
		normalized := (tile.Elevation - currentMin) / currentRange
		tile.Elevation = minElev + normalized*targetRange
		
		// Reclassify land/water after scaling
		tile.ClassifyLandWater(0.0) // Assume sea level is 0
	}
}