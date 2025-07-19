package terrain

import (
	"math"
	"sort"
)

// ValidateTerrain performs comprehensive statistical analysis of generated terrain
func ValidateTerrain(tiles []*HexTile) TerrainStats {
	if len(tiles) == 0 {
		return TerrainStats{}
	}
	
	// Extract elevation data
	elevations := make([]float64, len(tiles))
	landCount := 0
	waterCount := 0
	
	for i, tile := range tiles {
		elevations[i] = tile.Elevation
		if tile.IsLand {
			landCount++
		} else {
			waterCount++
		}
	}
	
	// Calculate basic statistics
	minElev, maxElev := findMinMaxFloat64(elevations)
	meanElev := calculateMean(elevations)
	stdDev := calculateStdDev(elevations, meanElev)
	
	// Calculate percentages
	totalTiles := len(tiles)
	landPercentage := float64(landCount) / float64(totalTiles) * 100.0
	waterPercentage := float64(waterCount) / float64(totalTiles) * 100.0
	
	// Calculate hypsometric curve match
	hypsometricMatch := calculateHypsometricMatch(elevations)
	
	return TerrainStats{
		ElevationRange:   [2]float64{minElev, maxElev},
		ElevationMean:    meanElev,
		ElevationStdDev:  stdDev,
		LandPercentage:   landPercentage,
		WaterPercentage:  waterPercentage,
		HypsometricMatch: hypsometricMatch,
		TotalTiles:       totalTiles,
		LandTiles:        landCount,
		WaterTiles:       waterCount,
	}
}

// IsRealisticTerrain checks if terrain passes Earth-realism validation
func IsRealisticTerrain(stats TerrainStats) (bool, []string) {
	var issues []string
	
	// Check elevation range
	if stats.ElevationRange[0] < ElevationMin*1.2 { // Allow 20% tolerance
		issues = append(issues, "minimum elevation too low (deeper than Mariana Trench)")
	}
	if stats.ElevationRange[1] > ElevationMax*1.2 {
		issues = append(issues, "maximum elevation too high (higher than Everest)")
	}
	
	// Check land/water ratio (Earth is ~29% land)
	if stats.LandPercentage < 20.0 || stats.LandPercentage > 40.0 {
		issues = append(issues, "land percentage outside realistic range (20-40%)")
	}
	
	// Check hypsometric curve match
	if stats.HypsometricMatch < 0.8 {
		issues = append(issues, "elevation distribution doesn't match Earth's hypsometric curve")
	}
	
	// Check for reasonable elevation variance
	expectedStdDev := 2000.0 // Approximately Earth's elevation std dev
	if stats.ElevationStdDev < expectedStdDev*0.5 || stats.ElevationStdDev > expectedStdDev*2.0 {
		issues = append(issues, "elevation variance outside realistic range")
	}
	
	return len(issues) == 0, issues
}

// ValidateHypsometricCurve checks how well elevation distribution matches Earth's
func ValidateHypsometricCurve(elevations []float64) float64 {
	return calculateHypsometricMatch(elevations)
}

// ValidateElevationRange ensures all elevations are within realistic bounds
func ValidateElevationRange(stats TerrainStats) bool {
	return stats.ElevationRange[0] >= ElevationMin && 
		   stats.ElevationRange[1] <= ElevationMax
}

// DetectElevationAnomalies finds unrealistic elevation patterns
func DetectElevationAnomalies(tiles []*HexTile) []string {
	var anomalies []string
	
	if len(tiles) == 0 {
		return anomalies
	}
	
	// Extract elevations for statistical analysis
	elevations := make([]float64, len(tiles))
	for i, tile := range tiles {
		elevations[i] = tile.Elevation
	}
	
	mean := calculateMean(elevations)
	stdDev := calculateStdDev(elevations, mean)
	
	// Detect extreme outliers (more than 3 standard deviations)
	outlierThreshold := 3.0
	outlierCount := 0
	
	for _, elev := range elevations {
		if math.Abs(elev-mean) > outlierThreshold*stdDev {
			outlierCount++
		}
	}
	
	if outlierCount > len(elevations)/100 { // More than 1% outliers
		anomalies = append(anomalies, "too many elevation outliers detected")
	}
	
	// Check for unrealistic elevation spikes
	minElev, maxElev := findMinMaxFloat64(elevations)
	if maxElev-minElev > 15000 { // Larger than Earth's range
		anomalies = append(anomalies, "elevation range exceeds Earth's total range")
	}
	
	// Check for flat terrain (no variation)
	if stdDev < 10.0 { // Less than 10m variation
		anomalies = append(anomalies, "terrain too flat (insufficient elevation variation)")
	}
	
	return anomalies
}

// calculateHypsometricMatch computes how well elevation distribution matches Earth's curve
func calculateHypsometricMatch(elevations []float64) float64 {
	if len(elevations) == 0 {
		return 0.0
	}
	
	// Sort elevations for percentile calculation
	sorted := make([]float64, len(elevations))
	copy(sorted, elevations)
	sort.Float64s(sorted)
	
	// Earth's hypsometric curve percentiles (approximate)
	earthPercentiles := []float64{
		-6000, // 10th percentile (deep ocean)
		-4000, // 20th percentile
		-2000, // 30th percentile  
		-500,  // 40th percentile
		-100,  // 50th percentile
		50,    // 60th percentile
		200,   // 70th percentile (land)
		500,   // 80th percentile
		1000,  // 90th percentile
		2000,  // 95th percentile
	}
	
	// Calculate our terrain's percentiles
	percentileIndices := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.95}
	ourPercentiles := make([]float64, len(percentileIndices))
	
	for i, p := range percentileIndices {
		index := int(p * float64(len(sorted)))
		if index >= len(sorted) {
			index = len(sorted) - 1
		}
		ourPercentiles[i] = sorted[index]
	}
	
	// Calculate correlation between our curve and Earth's curve
	correlation := calculateCorrelation(ourPercentiles, earthPercentiles)
	
	// Convert correlation to 0-1 range (correlation can be -1 to 1)
	return (correlation + 1.0) / 2.0
}

// Helper functions for statistical calculations

func findMinMaxFloat64(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	
	min := values[0]
	max := values[0]
	
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	
	return min, max
}

func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	
	return sum / float64(len(values))
}

func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}
	
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	
	variance := sumSquares / float64(len(values)-1)
	return math.Sqrt(variance)
}

func calculateCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}
	
	// Calculate means
	meanX := calculateMean(x)
	meanY := calculateMean(y)
	
	// Calculate correlation coefficient
	numerator := 0.0
	sumXSquares := 0.0
	sumYSquares := 0.0
	
	for i := 0; i < len(x); i++ {
		xDiff := x[i] - meanX
		yDiff := y[i] - meanY
		
		numerator += xDiff * yDiff
		sumXSquares += xDiff * xDiff
		sumYSquares += yDiff * yDiff
	}
	
	denominator := math.Sqrt(sumXSquares * sumYSquares)
	if denominator == 0 {
		return 0
	}
	
	return numerator / denominator
}

// GetElevationPercentiles calculates elevation percentiles for analysis
func GetElevationPercentiles(tiles []*HexTile, percentiles []float64) []float64 {
	if len(tiles) == 0 {
		return nil
	}
	
	elevations := make([]float64, len(tiles))
	for i, tile := range tiles {
		elevations[i] = tile.Elevation
	}
	
	sort.Float64s(elevations)
	
	result := make([]float64, len(percentiles))
	for i, p := range percentiles {
		index := int(p * float64(len(elevations)))
		if index >= len(elevations) {
			index = len(elevations) - 1
		}
		if index < 0 {
			index = 0
		}
		result[i] = elevations[index]
	}
	
	return result
}