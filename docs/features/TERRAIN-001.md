# TERRAIN-001: Basic Terrain Generation

## Objective
Implement realistic elevation generation using fractal noise with scientifically-based parameters that match Earth's terrain distribution. Generate base heightmaps using Diamond-Square algorithm with multi-octave noise and hypsometric curve validation.

## Technical Approach

### Fractal Terrain Generation
- **Diamond-Square Algorithm**: Primary terrain generation method with configurable roughness
- **Multi-Octave Noise**: Combine multiple noise frequencies (Hurst exponent H=0.8-0.9)
- **Spectral Synthesis**: Use real-world terrain power spectra with β ≈ 2
- **Hypsometric Curve**: Match Earth's elevation distribution (29% land, 71% ocean)

### Scientific Parameters
- **Elevation Range**: -11,000m to +9,000m (Mariana Trench to Everest)
- **Sea Level**: 0m with configurable adjustment
- **Land Ratio**: ~30% land coverage matching Earth
- **Fractal Dimension**: D ≈ 2.1 for realistic terrain roughness
- **Continental Shelf**: Gentle slopes (≈ 0.1°) at land-water transition

### Data Structure Design
```go
// Core terrain tile with geographic attributes
type HexTile struct {
    Coordinates     hex.AxialCoord
    Elevation       float64  // meters above sea level
    IsLand         bool     // land vs water classification
    DistanceToWater float64  // km to nearest water (for future use)
}

// Terrain generation configuration
type TerrainConfig struct {
    Seed        int64     // Random seed for reproducible generation
    SeaLevel    float64   // Elevation threshold for land/water (default: 0)
    LandRatio   float64   // Target percentage of land tiles (default: 0.3)
    NoiseParams NoiseParameters
}

// Multi-octave noise parameters
type NoiseParameters struct {
    Octaves     int     // Number of noise octaves (default: 6)
    Persistence float64 // Amplitude reduction per octave (default: 0.5)
    Lacunarity  float64 // Frequency increase per octave (default: 2.0)
    Scale       float64 // Initial noise scale (default: 0.01)
    HurstExp    float64 // Hurst exponent for fractal terrain (0.8-0.9)
}

// Statistical validation results
type TerrainStats struct {
    ElevationRange    [2]float64  // [min, max] elevation
    ElevationMean     float64     // Mean elevation
    ElevationStdDev   float64     // Standard deviation
    LandPercentage    float64     // Actual land coverage
    WaterPercentage   float64     // Actual water coverage
    HypsometricMatch  float64     // How well we match Earth's curve (0-1)
}
```

## API Design

### Core Functions
```go
package terrain

// Primary terrain generation function
func GenerateTerrain(grid *hex.Grid, config TerrainConfig) ([]*HexTile, error)

// Generate base heightmap using Diamond-Square algorithm
func GenerateHeightmap(width, height int, params NoiseParameters, seed int64) [][]float64

// Apply hypsometric curve adjustment to match Earth's distribution
func ApplyHypsometricCurve(heightmap [][]float64, targetLandRatio float64) [][]float64

// Convert heightmap to hex tiles with land/water classification
func HeightmapToHexTiles(heightmap [][]float64, grid *hex.Grid, seaLevel float64) []*HexTile

// Statistical validation and analysis
func ValidateTerrain(tiles []*HexTile) TerrainStats
func IsRealisticTerrain(stats TerrainStats) (bool, []string)
```

### Noise Generation
```go
// Diamond-Square fractal algorithm
func DiamondSquare(size int, roughness float64, seed int64) [][]float64

// Multi-octave noise combination
func MultiOctaveNoise(width, height int, params NoiseParameters, seed int64) [][]float64

// Spectral synthesis for realistic terrain frequencies
func SpectralSynthesis(width, height int, beta float64, seed int64) [][]float64
```

### Validation Functions
```go
// Check if elevation distribution matches Earth's hypsometric curve
func ValidateHypsometricCurve(elevations []float64) float64

// Ensure elevation ranges are within realistic bounds
func ValidateElevationRange(stats TerrainStats) bool

// Check for statistical anomalies (unrealistic spikes, etc.)
func DetectElevationAnomalies(tiles []*HexTile) []string
```

## File Structure
```
pkg/
  terrain/
    generator.go           // Main terrain generation engine
    generator_test.go      // Terrain generation tests
    heightmap.go          // HexTile and heightmap data structures
    heightmap_test.go     // Heightmap conversion tests
    noise.go              // Fractal noise algorithms
    noise_test.go         // Noise generation tests
    validation.go         // Statistical validation functions
    validation_test.go    // Validation logic tests
internal/
  noise/
    diamond_square.go     // Diamond-Square implementation
    diamond_square_test.go // Algorithm-specific tests
    spectral.go           // Spectral synthesis algorithms
cmd/
  hex-world/
    terrain.go           // Terrain-related CLI commands
```

## Testing Strategy

### Unit Tests
- **Noise Generation**: Verify Diamond-Square produces expected patterns
- **Statistical Properties**: Test noise parameters produce correct distributions
- **Hypsometric Validation**: Ensure curve matching works correctly
- **Land/Water Classification**: Verify sea level thresholding
- **Edge Cases**: Handle extreme parameters gracefully

### Property-Based Tests
- **Elevation Bounds**: All generated elevations within realistic range
- **Land Ratio**: Generated terrain matches target land percentage within tolerance
- **Statistical Distribution**: Elevation follows expected fractal properties
- **Deterministic**: Same seed produces identical terrain
- **Grid Independence**: Works with both region and world topologies

### Integration Tests
- **Full Pipeline**: Generate terrain end-to-end with validation
- **Different Grid Sizes**: Test performance and quality at various scales
- **Parameter Sensitivity**: Verify parameter changes produce expected effects
- **Memory Usage**: Profile memory consumption for large grids

### Validation Tests
- **Earth Realism**: Generated terrain passes Earth-like checks
- **Hypsometric Accuracy**: Elevation distribution within 5% of Earth's curve
- **No Anomalies**: No unrealistic elevation spikes or artifacts
- **Continental Shelves**: Smooth land-water transitions

## Implementation Phases

### Phase 1: Core Data Structures (Day 1)
1. Implement HexTile struct with elevation data
2. Create TerrainConfig and NoiseParameters types
3. Set up basic terrain package structure
4. Write initial data structure tests

### Phase 2: Noise Generation (Day 1-2)
1. Implement Diamond-Square algorithm
2. Add multi-octave noise combination
3. Create spectral synthesis for realistic frequencies
4. Comprehensive noise generation tests

### Phase 3: Terrain Generation Pipeline (Day 2)
1. Main GenerateTerrain function
2. Heightmap to hex tile conversion
3. Land/water classification logic
4. Integration with existing hex grid system

### Phase 4: Hypsometric Validation (Day 2-3)
1. Earth's hypsometric curve reference data
2. Curve matching and adjustment algorithms
3. Statistical validation functions
4. Realism checks and anomaly detection

### Phase 5: CLI Integration (Day 3)
1. Extend hex-world CLI with terrain commands
2. JSON export/import for terrain data
3. Statistical reporting and validation output
4. Demo commands and usage examples

## Scientific Accuracy

### Real-World Validation
- **Elevation Distribution**: Match Earth's hypsometric curve within 5%
- **Fractal Properties**: Power spectrum follows β ≈ 2 (realistic terrain)
- **Land Coverage**: 29-31% land ratio matching Earth's continents
- **Continental Margins**: Smooth transitions at coastlines

### Parameter Ranges
```go
// Scientifically validated parameter ranges
ELEVATION_MIN     = -11000.0  // Mariana Trench depth
ELEVATION_MAX     = 8849.0    // Everest height  
SEA_LEVEL_DEFAULT = 0.0       // Standard sea level
LAND_RATIO_EARTH  = 0.29      // Earth's land coverage
HURST_EXPONENT    = 0.85      // Typical terrain roughness
FRACTAL_DIMENSION = 2.15      // Realistic terrain complexity
```

### Quality Metrics
- **Hypsometric Match**: >0.95 correlation with Earth's curve
- **Elevation Variance**: Within 20% of real-world terrain
- **Land Distribution**: 29% ± 2% land coverage
- **No Artifacts**: Zero unrealistic elevation spikes or discontinuities

## CLI Commands

### Terrain Generation
```bash
# Generate basic terrain
./hex-world generate-terrain --size=100x100 --seed=42 --output=terrain.json

# Generate with custom parameters
./hex-world generate-terrain \
  --size=200x150 \
  --seed=12345 \
  --land-ratio=0.35 \
  --sea-level=-50 \
  --roughness=0.8 \
  --output=custom_terrain.json
```

### Statistical Analysis
```bash
# Show terrain statistics
./hex-world terrain-stats terrain.json

# Validate terrain realism
./hex-world validate-terrain terrain.json --strict

# Export elevation data for analysis
./hex-world export-elevation terrain.json --format=csv
```

### Demo Commands
```bash
# Quick terrain demo
./hex-world demo-terrain --size=50x50 --show-stats

# Compare different seeds
./hex-world compare-terrain --seeds=42,123,456 --size=100x100
```

## Success Criteria

### Functional Requirements
1. ✅ Generate realistic elevation maps using Diamond-Square algorithm
2. ✅ Produce terrain matching Earth's hypsometric curve within 5%
3. ✅ Support configurable land/water ratios and sea levels
4. ✅ Integration with existing hex grid and topology systems
5. ✅ Export terrain data in JSON format with metadata

### Quality Requirements
1. ✅ >90% test coverage for all terrain generation code
2. ✅ Performance: Generate 100x100 terrain in <1 second
3. ✅ Memory efficiency: Handle 1000x1000 grids without issues
4. ✅ Deterministic: Same seed always produces identical terrain
5. ✅ No artifacts: Generated terrain passes all realism checks

### Validation Requirements
1. ✅ Elevation ranges within realistic Earth bounds
2. ✅ Statistical properties match fractal terrain expectations
3. ✅ Land/water distribution matches target ratios
4. ✅ Smooth transitions at continental margins
5. ✅ No unrealistic elevation spikes or discontinuities

## Dependencies
- **HEX-001**: Hex grid system for coordinate management
- **Go Standard Library**: math/rand for seeded random generation
- **Future Integration**: Prepares for VIZ-001 (visualization) and HYDRO-001 (water flow)

## Performance Targets
- **Small Grids** (50x50): <100ms generation time
- **Medium Grids** (200x200): <1s generation time  
- **Large Grids** (1000x1000): <30s generation time
- **Memory Usage**: <50MB for 1000x1000 grid
- **Validation**: <10ms for statistical analysis of any grid size