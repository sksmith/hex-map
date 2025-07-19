package terrain

import (
	"testing"

	"github.com/sean/hex-map/pkg/hex"
)

func TestValidateTerrain(t *testing.T) {
	// Create sample terrain tiles
	tiles := []*HexTile{
		{Coordinates: hex.NewAxialCoord(0, 0), Elevation: -2000, IsLand: false},
		{Coordinates: hex.NewAxialCoord(1, 0), Elevation: 100, IsLand: true},
		{Coordinates: hex.NewAxialCoord(0, 1), Elevation: 1500, IsLand: true},
		{Coordinates: hex.NewAxialCoord(1, 1), Elevation: -500, IsLand: false},
	}
	
	stats := ValidateTerrain(tiles)
	
	// Check basic statistics
	if stats.TotalTiles != 4 {
		t.Errorf("Expected 4 total tiles, got %d", stats.TotalTiles)
	}
	
	if stats.LandTiles != 2 {
		t.Errorf("Expected 2 land tiles, got %d", stats.LandTiles)
	}
	
	if stats.WaterTiles != 2 {
		t.Errorf("Expected 2 water tiles, got %d", stats.WaterTiles)
	}
	
	expectedLandPercentage := 50.0
	if stats.LandPercentage != expectedLandPercentage {
		t.Errorf("Expected land percentage %.1f, got %.1f", expectedLandPercentage, stats.LandPercentage)
	}
	
	// Check elevation range
	expectedMin := -2000.0
	expectedMax := 1500.0
	
	if stats.ElevationRange[0] != expectedMin {
		t.Errorf("Expected min elevation %.1f, got %.1f", expectedMin, stats.ElevationRange[0])
	}
	
	if stats.ElevationRange[1] != expectedMax {
		t.Errorf("Expected max elevation %.1f, got %.1f", expectedMax, stats.ElevationRange[1])
	}
	
	// Check that hypsometric match is calculated (should be between 0 and 1)
	if stats.HypsometricMatch < 0 || stats.HypsometricMatch > 1 {
		t.Errorf("Hypsometric match should be between 0 and 1, got %f", stats.HypsometricMatch)
	}
}

func TestValidateTerrainEmpty(t *testing.T) {
	stats := ValidateTerrain([]*HexTile{})
	
	// Should handle empty input gracefully
	if stats.TotalTiles != 0 {
		t.Errorf("Expected 0 total tiles for empty input, got %d", stats.TotalTiles)
	}
	
	if stats.HypsometricMatch != 0 {
		t.Errorf("Expected 0 hypsometric match for empty input, got %f", stats.HypsometricMatch)
	}
}

func TestIsRealisticTerrain(t *testing.T) {
	tests := []struct {
		name       string
		stats      TerrainStats
		wantValid  bool
		wantIssues int
	}{
		{
			name: "realistic terrain",
			stats: TerrainStats{
				ElevationRange:   [2]float64{-3000, 3000},
				LandPercentage:   30.0,
				HypsometricMatch: 0.9,
				ElevationStdDev:  2000.0,
			},
			wantValid:  true,
			wantIssues: 0,
		},
		{
			name: "too deep ocean",
			stats: TerrainStats{
				ElevationRange:   [2]float64{-15000, 3000}, // Deeper than Mariana Trench * 1.2
				LandPercentage:   30.0,
				HypsometricMatch: 0.9,
				ElevationStdDev:  2000.0,
			},
			wantValid:  false,
			wantIssues: 1,
		},
		{
			name: "too high mountains",
			stats: TerrainStats{
				ElevationRange:   [2]float64{-3000, 15000}, // Higher than Everest * 1.2
				LandPercentage:   30.0,
				HypsometricMatch: 0.9,
				ElevationStdDev:  2000.0,
			},
			wantValid:  false,
			wantIssues: 1,
		},
		{
			name: "unrealistic land ratio",
			stats: TerrainStats{
				ElevationRange:   [2]float64{-3000, 3000},
				LandPercentage:   80.0, // Way too much land
				HypsometricMatch: 0.9,
				ElevationStdDev:  2000.0,
			},
			wantValid:  false,
			wantIssues: 1,
		},
		{
			name: "poor hypsometric match",
			stats: TerrainStats{
				ElevationRange:   [2]float64{-3000, 3000},
				LandPercentage:   30.0,
				HypsometricMatch: 0.5, // Poor Earth-like distribution
				ElevationStdDev:  2000.0,
			},
			wantValid:  false,
			wantIssues: 1,
		},
		{
			name: "too little elevation variance",
			stats: TerrainStats{
				ElevationRange:   [2]float64{-3000, 3000},
				LandPercentage:   30.0,
				HypsometricMatch: 0.9,
				ElevationStdDev:  500.0, // Too little variance
			},
			wantValid:  false,
			wantIssues: 1,
		},
		{
			name: "multiple issues",
			stats: TerrainStats{
				ElevationRange:   [2]float64{-15000, 15000}, // Both too deep and too high
				LandPercentage:   80.0,                       // Too much land
				HypsometricMatch: 0.5,                        // Poor match
				ElevationStdDev:  500.0,                      // Too little variance
			},
			wantValid:  false,
			wantIssues: 5, // All 5 issues should be detected
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, issues := IsRealisticTerrain(tt.stats)
			
			if isValid != tt.wantValid {
				t.Errorf("IsRealisticTerrain() validity = %v, want %v", isValid, tt.wantValid)
			}
			
			if len(issues) != tt.wantIssues {
				t.Errorf("IsRealisticTerrain() issues count = %d, want %d", len(issues), tt.wantIssues)
				t.Logf("Issues: %v", issues)
			}
		})
	}
}

func TestValidateElevationRange(t *testing.T) {
	tests := []struct {
		name  string
		stats TerrainStats
		want  bool
	}{
		{
			name: "valid range",
			stats: TerrainStats{
				ElevationRange: [2]float64{-5000, 5000},
			},
			want: true,
		},
		{
			name: "too low minimum",
			stats: TerrainStats{
				ElevationRange: [2]float64{-15000, 5000},
			},
			want: false,
		},
		{
			name: "too high maximum",
			stats: TerrainStats{
				ElevationRange: [2]float64{-5000, 15000},
			},
			want: false,
		},
		{
			name: "extreme range",
			stats: TerrainStats{
				ElevationRange: [2]float64{ElevationMin, ElevationMax},
			},
			want: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateElevationRange(tt.stats)
			if result != tt.want {
				t.Errorf("ValidateElevationRange() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDetectElevationAnomalies(t *testing.T) {
	// Create tiles with many extreme outliers to trigger the 1% threshold
	tiles := []*HexTile{}
	
	// Add 98 normal tiles (very similar values to ensure low std dev)
	for i := 0; i < 98; i++ {
		tiles = append(tiles, &HexTile{Elevation: 100.0})
	}
	
	// Add 2 extremely different outliers (>1% of 100 tiles)
	tiles = append(tiles, &HexTile{Elevation: 50000.0})  // Extreme outlier
	tiles = append(tiles, &HexTile{Elevation: -50000.0}) // Extreme outlier
	
	anomalies := DetectElevationAnomalies(tiles)
	
	// Should detect the outliers
	if len(anomalies) == 0 {
		t.Error("Expected to detect elevation anomalies, got none")
	}
	
	// Test with flat terrain
	flatTiles := []*HexTile{
		{Elevation: 100},
		{Elevation: 100},
		{Elevation: 100},
		{Elevation: 100},
	}
	
	flatAnomalies := DetectElevationAnomalies(flatTiles)
	
	// Should detect flat terrain
	hasFlat := false
	for _, anomaly := range flatAnomalies {
		if anomaly == "terrain too flat (insufficient elevation variation)" {
			hasFlat = true
			break
		}
	}
	
	if !hasFlat {
		t.Error("Expected to detect flat terrain anomaly")
	}
	
	// Test with extreme range
	extremeTiles := []*HexTile{
		{Elevation: -12000}, // Deeper than any ocean
		{Elevation: 10000},  // Higher than Everest
	}
	
	extremeAnomalies := DetectElevationAnomalies(extremeTiles)
	
	// Should detect extreme range
	hasExtreme := false
	for _, anomaly := range extremeAnomalies {
		if anomaly == "elevation range exceeds Earth's total range" {
			hasExtreme = true
			break
		}
	}
	
	if !hasExtreme {
		t.Error("Expected to detect extreme elevation range")
	}
}

func TestCalculateHypsometricMatch(t *testing.T) {
	// Test with Earth-like distribution
	earthLikeElevations := []float64{
		-6000, -4000, -2000, -500, -100, // Ocean depths
		50, 200, 500, 1000, 2000,        // Land heights
	}
	
	match := calculateHypsometricMatch(earthLikeElevations)
	
	// Should have reasonable match with Earth
	if match < 0.5 {
		t.Errorf("Earth-like elevations should have good hypsometric match, got %f", match)
	}
	
	// Test with flat distribution
	flatElevations := []float64{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	
	flatMatch := calculateHypsometricMatch(flatElevations)
	
	// Should have poor match
	if flatMatch > 0.7 {
		t.Errorf("Flat elevations should have poor hypsometric match, got %f", flatMatch)
	}
	
	// Test with empty input
	emptyMatch := calculateHypsometricMatch([]float64{})
	if emptyMatch != 0 {
		t.Errorf("Empty input should return 0, got %f", emptyMatch)
	}
}

func TestGetElevationPercentiles(t *testing.T) {
	// Create sample tiles
	tiles := []*HexTile{
		{Elevation: 0},
		{Elevation: 100},
		{Elevation: 200},
		{Elevation: 300},
		{Elevation: 400},
	}
	
	percentiles := []float64{0.0, 0.25, 0.5, 0.75, 1.0}
	result := GetElevationPercentiles(tiles, percentiles)
	
	if len(result) != len(percentiles) {
		t.Errorf("Expected %d percentiles, got %d", len(percentiles), len(result))
	}
	
	// Check that percentiles are in ascending order
	for i := 1; i < len(result); i++ {
		if result[i] < result[i-1] {
			t.Errorf("Percentiles not in ascending order: %v", result)
		}
	}
	
	// Test with empty input
	emptyResult := GetElevationPercentiles([]*HexTile{}, percentiles)
	if emptyResult != nil {
		t.Errorf("Expected nil for empty input, got %v", emptyResult)
	}
}

func TestStatisticalHelperFunctions(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	
	// Test mean calculation
	mean := calculateMean(values)
	expectedMean := 3.0
	if mean != expectedMean {
		t.Errorf("calculateMean() = %f, want %f", mean, expectedMean)
	}
	
	// Test standard deviation
	stdDev := calculateStdDev(values, mean)
	if stdDev <= 0 {
		t.Errorf("calculateStdDev() should be positive, got %f", stdDev)
	}
	
	// Test correlation
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 6, 8, 10} // Perfect positive correlation
	
	correlation := calculateCorrelation(x, y)
	if correlation < 0.99 { // Should be very close to 1
		t.Errorf("calculateCorrelation() = %f, want ~1.0", correlation)
	}
	
	// Test with negative correlation
	yNeg := []float64{10, 8, 6, 4, 2}
	negCorrelation := calculateCorrelation(x, yNeg)
	if negCorrelation > -0.99 { // Should be very close to -1
		t.Errorf("calculateCorrelation() = %f, want ~-1.0", negCorrelation)
	}
	
	// Test min/max finding
	min, max := findMinMaxFloat64(values)
	if min != 1.0 || max != 5.0 {
		t.Errorf("findMinMaxFloat64() = (%f, %f), want (1.0, 5.0)", min, max)
	}
}

func TestValidateHypsometricCurve(t *testing.T) {
	// Test with realistic elevation distribution
	elevations := []float64{
		-6000, -4000, -3000, -2000, -1000, // Deep ocean to shallow
		-500, -200, -100, -50, 0,           // Continental shelf
		50, 100, 200, 500, 1000,            // Low land
		2000, 3000, 4000, 5000, 6000,       // Mountains
	}
	
	match := ValidateHypsometricCurve(elevations)
	
	// Should be between 0 and 1
	if match < 0 || match > 1 {
		t.Errorf("ValidateHypsometricCurve() = %f, should be between 0 and 1", match)
	}
	
	// Test with empty input
	emptyMatch := ValidateHypsometricCurve([]float64{})
	if emptyMatch != 0 {
		t.Errorf("ValidateHypsometricCurve() with empty input = %f, want 0", emptyMatch)
	}
}