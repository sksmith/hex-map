package render

import (
	"image/color"
	"testing"
)

func TestTerrainColorScheme(t *testing.T) {
	// This should fail until we implement TerrainColorScheme
	colorMap := TerrainColorScheme()

	if colorMap.SeaLevel != 0.0 {
		t.Errorf("Expected sea level 0.0, got %f", colorMap.SeaLevel)
	}

	if len(colorMap.Breakpoints) == 0 {
		t.Error("TerrainColorScheme should have breakpoints")
	}

	// Verify breakpoints are sorted by elevation
	for i := 1; i < len(colorMap.Breakpoints); i++ {
		if colorMap.Breakpoints[i].Elevation <= colorMap.Breakpoints[i-1].Elevation {
			t.Error("Color breakpoints should be sorted by elevation")
		}
	}
}

func TestRealisticEarthScheme(t *testing.T) {
	// This should fail until we implement RealisticEarthScheme
	colorMap := RealisticEarthScheme()

	if len(colorMap.Breakpoints) == 0 {
		t.Error("RealisticEarthScheme should have breakpoints")
	}

	// Should have water colors for negative elevations
	hasWaterColors := false
	for _, bp := range colorMap.Breakpoints {
		if bp.Elevation < 0 {
			hasWaterColors = true
			break
		}
	}
	if !hasWaterColors {
		t.Error("RealisticEarthScheme should include water colors for negative elevations")
	}
}

func TestDebugColorScheme(t *testing.T) {
	// This should fail until we implement DebugColorScheme
	colorMap := DebugColorScheme()

	if len(colorMap.Breakpoints) == 0 {
		t.Error("DebugColorScheme should have breakpoints")
	}

	// Debug colors should be high contrast
	for _, bp := range colorMap.Breakpoints {
		if bp.Color.R == 0 && bp.Color.G == 0 && bp.Color.B == 0 {
			t.Error("Debug colors should not be black (low contrast)")
		}
	}
}

func TestInterpolateColor(t *testing.T) {
	c1 := color.RGBA{255, 0, 0, 255} // Red
	c2 := color.RGBA{0, 255, 0, 255} // Green

	tests := []struct {
		ratio    float64
		expected color.RGBA
	}{
		{0.0, color.RGBA{255, 0, 0, 255}},   // Should be c1
		{1.0, color.RGBA{0, 255, 0, 255}},   // Should be c2
		{0.5, color.RGBA{127, 127, 0, 255}}, // Should be halfway
	}

	for _, test := range tests {
		// This should fail until we implement InterpolateColor
		result := InterpolateColor(c1, c2, test.ratio)
		if result != test.expected {
			t.Errorf("InterpolateColor(%v, %v, %f) = %v, expected %v",
				c1, c2, test.ratio, result, test.expected)
		}
	}
}

func TestElevationToColor(t *testing.T) {
	// Create a simple test color map
	colorMap := ElevationColorMap{
		SeaLevel: 0.0,
		Breakpoints: []ColorBreakpoint{
			{-1000.0, color.RGBA{0, 0, 255, 255}},  // Deep blue for deep water
			{0.0, color.RGBA{0, 128, 255, 255}},    // Light blue for sea level
			{500.0, color.RGBA{0, 255, 0, 255}},    // Green for low land
			{2000.0, color.RGBA{139, 69, 19, 255}}, // Brown for high land
		},
	}

	tests := []struct {
		elevation float64
		expected  color.RGBA
	}{
		{-1000.0, color.RGBA{0, 0, 255, 255}},  // Exact match deep water
		{0.0, color.RGBA{0, 128, 255, 255}},    // Exact match sea level
		{250.0, color.RGBA{0, 191, 127, 255}},  // Interpolated between sea level and low land
		{2000.0, color.RGBA{139, 69, 19, 255}}, // Exact match high land
		{3000.0, color.RGBA{139, 69, 19, 255}}, // Above highest breakpoint
	}

	for _, test := range tests {
		// This should fail until we implement ElevationToColor
		result := ElevationToColor(test.elevation, colorMap)
		if result != test.expected {
			t.Errorf("ElevationToColor(%f, colorMap) = %v, expected %v",
				test.elevation, result, test.expected)
		}
	}
}

func TestColorBreakpointValidation(t *testing.T) {
	// Test that color maps have reasonable ranges
	schemes := []ElevationColorMap{
		TerrainColorScheme(),
		RealisticEarthScheme(),
		DebugColorScheme(),
	}

	for i, scheme := range schemes {
		if len(scheme.Breakpoints) < 3 {
			t.Errorf("Color scheme %d should have at least 3 breakpoints", i)
		}

		// Should cover reasonable elevation range
		minElevation := scheme.Breakpoints[0].Elevation
		maxElevation := scheme.Breakpoints[len(scheme.Breakpoints)-1].Elevation

		if minElevation >= 0 {
			t.Errorf("Color scheme %d should include negative elevations for water", i)
		}

		if maxElevation <= 1000 {
			t.Errorf("Color scheme %d should include elevations > 1000m for mountains", i)
		}
	}
}

func TestColorInterpolationBoundaries(t *testing.T) {
	c1 := color.RGBA{100, 150, 200, 255}
	c2 := color.RGBA{200, 100, 50, 255}

	// Test edge cases
	tests := []struct {
		ratio    float64
		expected color.RGBA
	}{
		{-0.1, color.RGBA{100, 150, 200, 255}}, // Below 0 should clamp to c1
		{1.1, color.RGBA{200, 100, 50, 255}},   // Above 1 should clamp to c2
		{0.0, color.RGBA{100, 150, 200, 255}},  // Exactly 0
		{1.0, color.RGBA{200, 100, 50, 255}},   // Exactly 1
	}

	for _, test := range tests {
		result := InterpolateColor(c1, c2, test.ratio)
		if result != test.expected {
			t.Errorf("InterpolateColor with ratio %f = %v, expected %v",
				test.ratio, result, test.expected)
		}
	}
}

func TestElevationToColorEdgeCases(t *testing.T) {
	// Test with minimal color map
	colorMap := ElevationColorMap{
		SeaLevel: 0.0,
		Breakpoints: []ColorBreakpoint{
			{-100.0, color.RGBA{0, 0, 255, 255}},
			{100.0, color.RGBA{0, 255, 0, 255}},
		},
	}

	// Test elevations outside the range
	tests := []struct {
		elevation float64
		expected  color.RGBA
	}{
		{-500.0, color.RGBA{0, 0, 255, 255}}, // Below lowest should use lowest color
		{500.0, color.RGBA{0, 255, 0, 255}},  // Above highest should use highest color
	}

	for _, test := range tests {
		result := ElevationToColor(test.elevation, colorMap)
		if result != test.expected {
			t.Errorf("ElevationToColor(%f) = %v, expected %v",
				test.elevation, result, test.expected)
		}
	}
}
