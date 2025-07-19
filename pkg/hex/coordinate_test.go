package hex

import (
	"math"
	"testing"
)

// TestAxialCoordCreation tests basic coordinate creation
func TestAxialCoordCreation(t *testing.T) {
	coord := NewAxialCoord(3, -1)
	if coord.Q != 3 || coord.R != -1 {
		t.Errorf("Expected (3, -1), got (%d, %d)", coord.Q, coord.R)
	}
}

// TestAxialToOffset tests axial to offset coordinate conversion
func TestAxialToOffset(t *testing.T) {
	tests := []struct {
		axial    AxialCoord
		expected [2]int // [col, row]
	}{
		{NewAxialCoord(0, 0), [2]int{0, 0}},
		{NewAxialCoord(1, 0), [2]int{1, 1}},   // even-q: row = r + (q+q&1)/2 = 0 + (1+1)/2 = 1
		{NewAxialCoord(0, 1), [2]int{0, 1}},   // even-q: row = r + (q+q&1)/2 = 1 + (0+0)/2 = 1
		{NewAxialCoord(1, 1), [2]int{1, 2}},   // even-q: row = r + (q+q&1)/2 = 1 + (1+1)/2 = 2
		{NewAxialCoord(-1, 1), [2]int{-1, 1}}, // even-q: row = r + (q+q&1)/2 = 1 + (-1+(-1&1))/2 = 1 + (-1+1)/2 = 1
		{NewAxialCoord(2, -1), [2]int{2, 0}},  // even-q: row = r + (q+q&1)/2 = -1 + (2+0)/2 = 0
	}

	for _, test := range tests {
		col, row := test.axial.ToOffset()
		if col != test.expected[0] || row != test.expected[1] {
			t.Errorf("ToOffset(%v) = (%d, %d), expected (%d, %d)",
				test.axial, col, row, test.expected[0], test.expected[1])
		}
	}
}

// TestOffsetToAxial tests offset to axial coordinate conversion
func TestOffsetToAxial(t *testing.T) {
	tests := []struct {
		col, row int
		expected AxialCoord
	}{
		{0, 0, NewAxialCoord(0, 0)},
		{1, 1, NewAxialCoord(1, 0)},  // even-q: r = row - (col+(col&1))/2 = 1 - (1+1)/2 = 0
		{0, 1, NewAxialCoord(0, 1)},  // even-q: r = row - (col+(col&1))/2 = 1 - (0+0)/2 = 1
		{1, 2, NewAxialCoord(1, 1)},  // even-q: r = row - (col+(col&1))/2 = 2 - (1+1)/2 = 1
		{-1, 1, NewAxialCoord(-1, 1)}, // even-q: r = row - (col+(col&1))/2 = 1 - (-1+1)/2 = 1
		{2, 0, NewAxialCoord(2, -1)}, // even-q: r = row - (col+(col&1))/2 = 0 - (2+0)/2 = -1
	}

	for _, test := range tests {
		result := OffsetToAxial(test.col, test.row)
		if result.Q != test.expected.Q || result.R != test.expected.R {
			t.Errorf("OffsetToAxial(%d, %d) = %v, expected %v",
				test.col, test.row, result, test.expected)
		}
	}
}

// TestAxialOffsetRoundTrip tests that axial ↔ offset conversion is symmetric
func TestAxialOffsetRoundTrip(t *testing.T) {
	coords := []AxialCoord{
		{0, 0}, {1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, -1},
		{2, -1}, {1, 1}, {-1, 1}, {-1, 2}, {-2, 1}, {-1, -1},
	}

	for _, original := range coords {
		col, row := original.ToOffset()
		roundTrip := OffsetToAxial(col, row)
		if roundTrip.Q != original.Q || roundTrip.R != original.R {
			t.Errorf("Round trip failed: %v → (%d,%d) → %v",
				original, col, row, roundTrip)
		}
	}
}

// TestPixelConversion tests axial to pixel coordinate conversion
func TestPixelConversion(t *testing.T) {
	hexSize := 10.0
	tests := []struct {
		axial AxialCoord
		x, y  float64
	}{
		{NewAxialCoord(0, 0), 0, 0},
		{NewAxialCoord(1, 0), 15, 8.66}, // x = 1.5*hexSize, y = sqrt(3)/2*hexSize
		{NewAxialCoord(0, 1), 0, 17.32}, // x = 0, y = sqrt(3)*hexSize
	}

	for _, test := range tests {
		x, y := test.axial.ToPixel(hexSize)
		if math.Abs(x-test.x) > 0.1 || math.Abs(y-test.y) > 0.1 {
			t.Errorf("ToPixel(%v, %f) = (%f, %f), expected (%f, %f)",
				test.axial, hexSize, x, y, test.x, test.y)
		}
	}
}

// TestPixelToAxial tests pixel to axial coordinate conversion
func TestPixelToAxial(t *testing.T) {
	hexSize := 10.0
	tests := []struct {
		x, y     float64
		expected AxialCoord
	}{
		{0, 0, NewAxialCoord(0, 0)},
		{15, 8.66, NewAxialCoord(1, 0)},
		{0, 17.32, NewAxialCoord(0, 1)},
	}

	for _, test := range tests {
		result := PixelToAxial(test.x, test.y, hexSize)
		if result.Q != test.expected.Q || result.R != test.expected.R {
			t.Errorf("PixelToAxial(%f, %f, %f) = %v, expected %v",
				test.x, test.y, hexSize, result, test.expected)
		}
	}
}

// TestPixelRoundTrip tests that axial ↔ pixel conversion is symmetric
func TestPixelRoundTrip(t *testing.T) {
	hexSize := 10.0
	coords := []AxialCoord{
		{0, 0}, {1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, -1},
		{2, -1}, {1, 1}, {-1, 1},
	}

	for _, original := range coords {
		x, y := original.ToPixel(hexSize)
		roundTrip := PixelToAxial(x, y, hexSize)
		if roundTrip.Q != original.Q || roundTrip.R != original.R {
			t.Errorf("Round trip failed: %v → (%f,%f) → %v",
				original, x, y, roundTrip)
		}
	}
}