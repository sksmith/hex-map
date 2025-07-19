package noise

import (
	"math"
	"testing"
)

func TestDiamondSquare(t *testing.T) {
	size := 129 // 2^7 + 1
	roughness := 0.5
	seed := int64(42)
	
	heightmap := DiamondSquare(size, roughness, seed)
	
	// Check dimensions
	if len(heightmap) != size {
		t.Errorf("Expected size %d, got %d", size, len(heightmap))
	}
	
	for i, row := range heightmap {
		if len(row) != size {
			t.Errorf("Row %d has wrong size: expected %d, got %d", i, size, len(row))
		}
	}
	
	// Check that values are in reasonable range
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			value := heightmap[y][x]
			if math.IsNaN(value) || math.IsInf(value, 0) {
				t.Errorf("Invalid value at (%d,%d): %f", x, y, value)
			}
			
			// Values should be roughly in [-2, 2] range for most realistic terrain
			if value < -5.0 || value > 5.0 {
				t.Errorf("Value out of expected range at (%d,%d): %f", x, y, value)
			}
		}
	}
	
	// Check determinism - same seed should produce same result
	heightmap2 := DiamondSquare(size, roughness, seed)
	
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if heightmap[y][x] != heightmap2[y][x] {
				t.Errorf("Non-deterministic generation at (%d,%d): %f vs %f", 
					x, y, heightmap[y][x], heightmap2[y][x])
			}
		}
	}
	
	// Different seeds should produce different results
	heightmap3 := DiamondSquare(size, roughness, seed+1)
	
	different := false
	for y := 0; y < size && !different; y++ {
		for x := 0; x < size; x++ {
			if heightmap[y][x] != heightmap3[y][x] {
				different = true
				break
			}
		}
	}
	
	if !different {
		t.Error("Different seeds should produce different terrain")
	}
}

func TestDiamondSquareInvalidSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid size, but didn't panic")
		}
	}()
	
	// Should panic for size that's not (2^n + 1)
	DiamondSquare(100, 0.5, 42)
}

func TestIsPowerOfTwoPlusOne(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{1, false},   // 2^0 = 1, but we need 2^0 + 1 = 2
		{2, false},   // Not 2^n + 1
		{3, true},    // 2^1 + 1 = 3
		{5, true},    // 2^2 + 1 = 5
		{9, true},    // 2^3 + 1 = 9
		{17, true},   // 2^4 + 1 = 17
		{33, true},   // 2^5 + 1 = 33
		{65, true},   // 2^6 + 1 = 65
		{129, true},  // 2^7 + 1 = 129
		{257, true},  // 2^8 + 1 = 257
		{10, false},  // Not 2^n + 1
		{16, false},  // 2^4 = 16, but we need 2^4 + 1 = 17
		{128, false}, // 2^7 = 128, but we need 2^7 + 1 = 129
	}
	
	for _, tt := range tests {
		t.Run(string(rune(tt.n)), func(t *testing.T) {
			result := isPowerOfTwoPlusOne(tt.n)
			if result != tt.want {
				t.Errorf("isPowerOfTwoPlusOne(%d) = %v, want %v", tt.n, result, tt.want)
			}
		})
	}
}

func TestMultiOctaveNoise(t *testing.T) {
	width, height := 50, 40
	octaves := 4
	persistence := 0.5
	lacunarity := 2.0
	scale := 0.01
	seed := int64(42)
	
	result := MultiOctaveNoise(width, height, octaves, persistence, lacunarity, scale, seed)
	
	// Check dimensions
	if len(result) != height {
		t.Errorf("Expected height %d, got %d", height, len(result))
	}
	
	if len(result[0]) != width {
		t.Errorf("Expected width %d, got %d", width, len(result[0]))
	}
	
	// Check that values are normalized (should be roughly in [-1, 1])
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := result[y][x]
			if math.IsNaN(value) || math.IsInf(value, 0) {
				t.Errorf("Invalid value at (%d,%d): %f", x, y, value)
			}
			
			// Multi-octave noise should be normalized to roughly [-1, 1]
			if value < -2.0 || value > 2.0 {
				t.Errorf("Value possibly out of normalized range at (%d,%d): %f", x, y, value)
			}
		}
	}
	
	// Test determinism
	result2 := MultiOctaveNoise(width, height, octaves, persistence, lacunarity, scale, seed)
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if result[y][x] != result2[y][x] {
				t.Errorf("Non-deterministic generation at (%d,%d): %f vs %f", 
					x, y, result[y][x], result2[y][x])
			}
		}
	}
}

func TestNextPowerOfTwoPlusOne(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{1, 3},    // smallest is 2^1 + 1 = 3
		{2, 3},    // 2^1 + 1 = 3
		{3, 3},    // already 2^1 + 1 = 3
		{4, 5},    // 2^2 + 1 = 5
		{5, 5},    // already 2^2 + 1 = 5
		{8, 9},    // 2^3 + 1 = 9
		{9, 9},    // already 2^3 + 1 = 9
		{16, 17},  // 2^4 + 1 = 17
		{17, 17},  // already 2^4 + 1 = 17
		{32, 33},  // 2^5 + 1 = 33
		{64, 65},  // 2^6 + 1 = 65
		{128, 129}, // 2^7 + 1 = 129
	}
	
	for _, tt := range tests {
		t.Run(string(rune(tt.input)), func(t *testing.T) {
			result := nextPowerOfTwoPlusOne(tt.input)
			if result != tt.want {
				t.Errorf("nextPowerOfTwoPlusOne(%d) = %d, want %d", tt.input, result, tt.want)
			}
		})
	}
}

func TestSpectralSynthesis(t *testing.T) {
	width, height := 32, 24
	beta := 2.0 // Typical value for realistic terrain
	seed := int64(42)
	
	result := SpectralSynthesis(width, height, beta, seed)
	
	// Check dimensions
	if len(result) != height {
		t.Errorf("Expected height %d, got %d", height, len(result))
	}
	
	if len(result[0]) != width {
		t.Errorf("Expected width %d, got %d", width, len(result[0]))
	}
	
	// Check that values are normalized to [-1, 1]
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := result[y][x]
			if math.IsNaN(value) || math.IsInf(value, 0) {
				t.Errorf("Invalid value at (%d,%d): %f", x, y, value)
			}
			
			// Should be normalized to [-1, 1]
			if value < -1.1 || value > 1.1 {
				t.Errorf("Value out of normalized range at (%d,%d): %f", x, y, value)
			}
		}
	}
	
	// Test determinism
	result2 := SpectralSynthesis(width, height, beta, seed)
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if math.Abs(result[y][x]-result2[y][x]) > 1e-10 {
				t.Errorf("Non-deterministic generation at (%d,%d): %f vs %f", 
					x, y, result[y][x], result2[y][x])
			}
		}
	}
	
	// Different seeds should produce different results
	result3 := SpectralSynthesis(width, height, beta, seed+1)
	
	different := false
	for y := 0; y < height && !different; y++ {
		for x := 0; x < width; x++ {
			if math.Abs(result[y][x]-result3[y][x]) > 1e-10 {
				different = true
				break
			}
		}
	}
	
	if !different {
		t.Error("Different seeds should produce different terrain")
	}
}

func TestFindMinMax(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]float64
		wantMin  float64
		wantMax  float64
	}{
		{
			name: "simple case",
			data: [][]float64{
				{1.0, 2.0, 3.0},
				{4.0, 5.0, 6.0},
			},
			wantMin: 1.0,
			wantMax: 6.0,
		},
		{
			name: "negative values",
			data: [][]float64{
				{-5.0, -2.0, 1.0},
				{-10.0, 3.0, 7.0},
			},
			wantMin: -10.0,
			wantMax: 7.0,
		},
		{
			name: "single value",
			data: [][]float64{
				{42.0},
			},
			wantMin: 42.0,
			wantMax: 42.0,
		},
		{
			name:    "empty data",
			data:    [][]float64{},
			wantMin: 0.0,
			wantMax: 0.0,
		},
		{
			name:    "empty row",
			data:    [][]float64{{}},
			wantMin: 0.0,
			wantMax: 0.0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := findMinMax(tt.data)
			
			if min != tt.wantMin {
				t.Errorf("findMinMax() min = %f, want %f", min, tt.wantMin)
			}
			
			if max != tt.wantMax {
				t.Errorf("findMinMax() max = %f, want %f", max, tt.wantMax)
			}
		})
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{1, 2, 2},
		{5, 3, 5},
		{-1, -5, -1},
		{0, 0, 0},
		{-10, 10, 10},
	}
	
	for _, tt := range tests {
		result := max(tt.a, tt.b)
		if result != tt.want {
			t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.want)
		}
	}
}

// Benchmark tests to ensure reasonable performance

func BenchmarkDiamondSquare(b *testing.B) {
	size := 129 // 2^7 + 1
	roughness := 0.5
	seed := int64(42)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DiamondSquare(size, roughness, seed)
	}
}

func BenchmarkMultiOctaveNoise(b *testing.B) {
	width, height := 100, 100
	octaves := 6
	persistence := 0.5
	lacunarity := 2.0
	scale := 0.01
	seed := int64(42)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiOctaveNoise(width, height, octaves, persistence, lacunarity, scale, seed)
	}
}

func BenchmarkSpectralSynthesis(b *testing.B) {
	width, height := 64, 64
	beta := 2.0
	seed := int64(42)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SpectralSynthesis(width, height, beta, seed)
	}
}