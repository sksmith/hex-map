package terrain

import (
	"github.com/sean/hex-map/pkg/hex"
)

// HexTile represents a single hex tile with terrain attributes
type HexTile struct {
	Coordinates     hex.AxialCoord `json:"coordinates"`
	Elevation       float64        `json:"elevation"`        // meters above sea level
	IsLand         bool           `json:"is_land"`          // land vs water classification
	DistanceToWater float64        `json:"distance_to_water"` // km to nearest water (future use)
}

// TerrainConfig contains all parameters for terrain generation
type TerrainConfig struct {
	Seed        int64           `json:"seed"`         // Random seed for reproducible generation
	SeaLevel    float64         `json:"sea_level"`    // Elevation threshold for land/water
	LandRatio   float64         `json:"land_ratio"`   // Target percentage of land tiles
	NoiseParams NoiseParameters `json:"noise_params"` // Multi-octave noise configuration
}

// NoiseParameters controls the fractal noise generation
type NoiseParameters struct {
	Octaves     int     `json:"octaves"`     // Number of noise octaves
	Persistence float64 `json:"persistence"` // Amplitude reduction per octave
	Lacunarity  float64 `json:"lacunarity"`  // Frequency increase per octave
	Scale       float64 `json:"scale"`       // Initial noise scale
	HurstExp    float64 `json:"hurst_exp"`   // Hurst exponent for fractal terrain
}

// TerrainStats provides statistical analysis of generated terrain
type TerrainStats struct {
	ElevationRange   [2]float64 `json:"elevation_range"`    // [min, max] elevation
	ElevationMean    float64    `json:"elevation_mean"`     // Mean elevation
	ElevationStdDev  float64    `json:"elevation_std_dev"`  // Standard deviation
	LandPercentage   float64    `json:"land_percentage"`    // Actual land coverage
	WaterPercentage  float64    `json:"water_percentage"`   // Actual water coverage
	HypsometricMatch float64    `json:"hypsometric_match"`  // Earth curve match (0-1)
	TotalTiles       int        `json:"total_tiles"`        // Total number of tiles
	LandTiles        int        `json:"land_tiles"`         // Number of land tiles
	WaterTiles       int        `json:"water_tiles"`        // Number of water tiles
}

// DefaultTerrainConfig returns scientifically-based default parameters
func DefaultTerrainConfig() TerrainConfig {
	return TerrainConfig{
		Seed:      42,
		SeaLevel:  0.0,
		LandRatio: 0.29, // Earth's actual land coverage
		NoiseParams: NoiseParameters{
			Octaves:     6,
			Persistence: 0.5,
			Lacunarity:  2.0,
			Scale:       0.01,
			HurstExp:    0.85, // Typical terrain roughness
		},
	}
}

// DefaultNoiseParameters returns default fractal noise settings
func DefaultNoiseParameters() NoiseParameters {
	return NoiseParameters{
		Octaves:     6,
		Persistence: 0.5,
		Lacunarity:  2.0,
		Scale:       0.01,
		HurstExp:    0.85,
	}
}

// Validate checks if terrain configuration parameters are reasonable
func (tc TerrainConfig) Validate() error {
	if tc.LandRatio < 0.0 || tc.LandRatio > 1.0 {
		return &TerrainError{"land_ratio must be between 0.0 and 1.0"}
	}
	
	if tc.NoiseParams.Octaves < 1 || tc.NoiseParams.Octaves > 10 {
		return &TerrainError{"octaves must be between 1 and 10"}
	}
	
	if tc.NoiseParams.Persistence <= 0.0 || tc.NoiseParams.Persistence > 1.0 {
		return &TerrainError{"persistence must be between 0.0 and 1.0"}
	}
	
	if tc.NoiseParams.Lacunarity <= 1.0 {
		return &TerrainError{"lacunarity must be greater than 1.0"}
	}
	
	if tc.NoiseParams.HurstExp < 0.0 || tc.NoiseParams.HurstExp > 1.0 {
		return &TerrainError{"hurst_exp must be between 0.0 and 1.0"}
	}
	
	return nil
}

// TerrainError represents an error in terrain generation
type TerrainError struct {
	Message string
}

func (e *TerrainError) Error() string {
	return "terrain error: " + e.Message
}

// Constants for realistic terrain bounds (based on Earth)
const (
	ElevationMin     = -11000.0 // Mariana Trench depth
	ElevationMax     = 8849.0   // Everest height
	SeaLevelDefault  = 0.0      // Standard sea level
	LandRatioEarth   = 0.29     // Earth's land coverage
	HurstExponent    = 0.85     // Typical terrain roughness
	FractalDimension = 2.15     // Realistic terrain complexity
)

// IsRealistic checks if a HexTile has realistic terrain values
func (ht *HexTile) IsRealistic() bool {
	return ht.Elevation >= ElevationMin && ht.Elevation <= ElevationMax
}

// ClassifyLandWater determines if a tile is land or water based on sea level
func (ht *HexTile) ClassifyLandWater(seaLevel float64) {
	ht.IsLand = ht.Elevation > seaLevel
}

// GetDepth returns water depth (positive value) if underwater, 0 if above sea level
func (ht *HexTile) GetDepth(seaLevel float64) float64 {
	if ht.Elevation < seaLevel {
		return seaLevel - ht.Elevation
	}
	return 0.0
}

// GetHeight returns height above sea level if on land, 0 if underwater
func (ht *HexTile) GetHeight(seaLevel float64) float64 {
	if ht.Elevation > seaLevel {
		return ht.Elevation - seaLevel
	}
	return 0.0
}