package render

import (
	"image/color"
	"math"
)

// ElevationColorMap maps elevation ranges to colors
type ElevationColorMap struct {
	SeaLevel    float64
	Breakpoints []ColorBreakpoint
}

type ColorBreakpoint struct {
	Elevation float64
	Color     color.RGBA
}

// TerrainColorScheme returns standard topographic color mapping
func TerrainColorScheme() ElevationColorMap {
	return ElevationColorMap{
		SeaLevel: 0.0,
		Breakpoints: []ColorBreakpoint{
			{-11000.0, color.RGBA{0, 0, 128, 255}},   // Deep ocean - dark blue
			{-2000.0, color.RGBA{0, 64, 192, 255}},   // Deep water - medium blue
			{-200.0, color.RGBA{0, 128, 255, 255}},   // Shallow water - light blue
			{0.0, color.RGBA{244, 164, 96, 255}},     // Beach/coast - tan
			{100.0, color.RGBA{144, 238, 144, 255}},  // Low plains - light green
			{500.0, color.RGBA{34, 139, 34, 255}},    // Low hills - green
			{1000.0, color.RGBA{154, 205, 50, 255}},  // Hills - yellow-green
			{2000.0, color.RGBA{139, 69, 19, 255}},   // Mountains - brown
			{3000.0, color.RGBA{160, 82, 45, 255}},   // High mountains - saddle brown
			{5000.0, color.RGBA{169, 169, 169, 255}}, // Very high - gray
			{8000.0, color.RGBA{255, 255, 255, 255}}, // Peaks - white
		},
	}
}

// RealisticEarthScheme returns Earth-like realistic color mapping
func RealisticEarthScheme() ElevationColorMap {
	return ElevationColorMap{
		SeaLevel: 0.0,
		Breakpoints: []ColorBreakpoint{
			{-11000.0, color.RGBA{0, 0, 128, 255}},   // Deep ocean - navy
			{-3000.0, color.RGBA{0, 50, 150, 255}},   // Ocean depths
			{-1000.0, color.RGBA{0, 102, 204, 255}},  // Ocean
			{-100.0, color.RGBA{102, 178, 255, 255}}, // Coastal waters
			{0.0, color.RGBA{244, 196, 161, 255}},    // Sand/coast
			{50.0, color.RGBA{200, 230, 200, 255}},   // Coastal plains
			{200.0, color.RGBA{34, 139, 34, 255}},    // Lowlands - forest green
			{800.0, color.RGBA{0, 100, 0, 255}},      // Foothills - dark green
			{1500.0, color.RGBA{160, 82, 45, 255}},   // Low mountains - brown
			{3000.0, color.RGBA{105, 105, 105, 255}}, // High mountains - dim gray
			{6000.0, color.RGBA{248, 248, 255, 255}}, // Snow peaks - ghost white
		},
	}
}

// DebugColorScheme returns high-contrast colors for debugging
func DebugColorScheme() ElevationColorMap {
	return ElevationColorMap{
		SeaLevel: 0.0,
		Breakpoints: []ColorBreakpoint{
			{-5000.0, color.RGBA{0, 0, 255, 255}},  // Bright blue for water
			{0.0, color.RGBA{255, 255, 0, 255}},    // Yellow for sea level
			{1000.0, color.RGBA{0, 255, 0, 255}},   // Bright green for low land
			{3000.0, color.RGBA{255, 0, 0, 255}},   // Red for high land
			{8000.0, color.RGBA{255, 0, 255, 255}}, // Magenta for peaks
		},
	}
}

// InterpolateColor linearly interpolates between two colors
func InterpolateColor(c1, c2 color.RGBA, ratio float64) color.RGBA {
	// Clamp ratio to [0, 1]
	if ratio < 0.0 {
		ratio = 0.0
	}
	if ratio > 1.0 {
		ratio = 1.0
	}

	// Linear interpolation for each channel
	r := uint8(float64(c1.R)*(1.0-ratio) + float64(c2.R)*ratio)
	g := uint8(float64(c1.G)*(1.0-ratio) + float64(c2.G)*ratio)
	b := uint8(float64(c1.B)*(1.0-ratio) + float64(c2.B)*ratio)
	a := uint8(float64(c1.A)*(1.0-ratio) + float64(c2.A)*ratio)

	return color.RGBA{r, g, b, a}
}

// ElevationToColor maps elevation to color using the provided color map
func ElevationToColor(elevation float64, colorMap ElevationColorMap) color.RGBA {
	if len(colorMap.Breakpoints) == 0 {
		return color.RGBA{0, 0, 0, 255} // Black if no breakpoints
	}

	// Handle edge cases
	if elevation <= colorMap.Breakpoints[0].Elevation {
		return colorMap.Breakpoints[0].Color
	}

	lastIdx := len(colorMap.Breakpoints) - 1
	if elevation >= colorMap.Breakpoints[lastIdx].Elevation {
		return colorMap.Breakpoints[lastIdx].Color
	}

	// Find the two breakpoints to interpolate between
	for i := 0; i < len(colorMap.Breakpoints)-1; i++ {
		bp1 := colorMap.Breakpoints[i]
		bp2 := colorMap.Breakpoints[i+1]

		if elevation >= bp1.Elevation && elevation <= bp2.Elevation {
			// Calculate interpolation ratio
			range_ := bp2.Elevation - bp1.Elevation
			if math.Abs(range_) < 1e-6 {
				return bp1.Color // Avoid division by zero
			}

			ratio := (elevation - bp1.Elevation) / range_
			return InterpolateColor(bp1.Color, bp2.Color, ratio)
		}
	}

	// Fallback (should not reach here)
	return colorMap.Breakpoints[lastIdx].Color
}
