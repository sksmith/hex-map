package hex

import (
	"math"
)

// AxialCoord represents a hexagon coordinate using axial coordinates (q, r)
// This is the most efficient coordinate system for hex grids
type AxialCoord struct {
	Q, R int
}

// NewAxialCoord creates a new axial coordinate
func NewAxialCoord(q, r int) AxialCoord {
	return AxialCoord{Q: q, R: r}
}

// ToOffset converts axial coordinates to offset coordinates (col, row)
// Uses flat-top hexagon orientation with even-q offset layout
func (c AxialCoord) ToOffset() (col, row int) {
	col = c.Q
	row = c.R + (c.Q+(c.Q&1))/2
	return col, row
}

// OffsetToAxial converts offset coordinates to axial coordinates
// Uses flat-top hexagon orientation with even-q offset layout
func OffsetToAxial(col, row int) AxialCoord {
	q := col
	r := row - (col+(col&1))/2
	return AxialCoord{Q: q, R: r}
}

// ToPixel converts axial coordinates to pixel coordinates
// Uses flat-top hexagon orientation
func (c AxialCoord) ToPixel(hexSize float64) (x, y float64) {
	x = hexSize * (3.0/2.0 * float64(c.Q))
	y = hexSize * (math.Sqrt(3.0)/2.0*float64(c.Q) + math.Sqrt(3.0)*float64(c.R))
	return x, y
}

// PixelToAxial converts pixel coordinates to axial coordinates
// Uses flat-top hexagon orientation
func PixelToAxial(x, y, hexSize float64) AxialCoord {
	q := (2.0/3.0) * x / hexSize
	r := (-1.0/3.0*x + math.Sqrt(3.0)/3.0*y) / hexSize
	return axialRound(q, r)
}

// axialRound rounds fractional axial coordinates to the nearest hex
func axialRound(q, r float64) AxialCoord {
	s := -q - r
	
	rq := math.Round(q)
	rr := math.Round(r)
	rs := math.Round(s)
	
	qDiff := math.Abs(rq - q)
	rDiff := math.Abs(rr - r)
	sDiff := math.Abs(rs - s)
	
	if qDiff > rDiff && qDiff > sDiff {
		rq = -rr - rs
	} else if rDiff > sDiff {
		rr = -rq - rs
	}
	
	return AxialCoord{Q: int(rq), R: int(rr)}
}