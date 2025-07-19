package hex

// Topology defines how grid edges behave
type Topology int

const (
	TopologyRegion Topology = iota // Bounded edges, fewer neighbors at boundaries
	TopologyWorld                  // Toroidal wrapping, all hexes have 6 neighbors
)

// Grid represents a hexagonal grid with configurable topology
type Grid struct {
	config   GridConfig
	tiles    [][]interface{}
	coordMap map[AxialCoord]bool
}

// GridConfig defines the configuration for a hex grid
type GridConfig struct {
	Width, Height int
	Topology      Topology
}

// NewGrid creates a new hexagonal grid with the specified configuration
func NewGrid(config GridConfig) *Grid {
	tiles := make([][]interface{}, config.Height)
	for i := range tiles {
		tiles[i] = make([]interface{}, config.Width)
	}

	coordMap := make(map[AxialCoord]bool)
	
	// Pre-populate coordinate map for faster lookups
	for row := 0; row < config.Height; row++ {
		for col := 0; col < config.Width; col++ {
			coord := OffsetToAxial(col, row)
			coordMap[coord] = true
		}
	}

	return &Grid{
		config:   config,
		tiles:    tiles,
		coordMap: coordMap,
	}
}

// Topology returns the topology type of this grid
func (g *Grid) Topology() Topology {
	return g.config.Topology
}

// IsValid checks if a coordinate is valid within this grid
func (g *Grid) IsValid(coord AxialCoord) bool {
	if g.config.Topology == TopologyWorld {
		// In world topology, check if the wrapped coordinate is in the grid
		wrapped := g.WrapCoord(coord)
		return g.coordMap[wrapped]
	}
	
	// For region topology, check if coordinate is in our map
	return g.coordMap[coord]
}

// WrapCoord wraps a coordinate for world topology
func (g *Grid) WrapCoord(coord AxialCoord) AxialCoord {
	if g.config.Topology != TopologyWorld {
		return coord
	}

	// Convert to offset for easier wrapping calculation
	col, row := coord.ToOffset()
	
	// Wrap coordinates
	col = ((col % g.config.Width) + g.config.Width) % g.config.Width
	row = ((row % g.config.Height) + g.config.Height) % g.config.Height
	
	// Convert back to axial
	return OffsetToAxial(col, row)
}

// Get retrieves a value from the grid at the specified coordinate
func (g *Grid) Get(coord AxialCoord) interface{} {
	if g.config.Topology == TopologyWorld {
		coord = g.WrapCoord(coord)
	}
	
	if !g.IsValid(coord) {
		return nil
	}
	
	col, row := coord.ToOffset()
	return g.tiles[row][col]
}

// Set stores a value in the grid at the specified coordinate
func (g *Grid) Set(coord AxialCoord, value interface{}) {
	if g.config.Topology == TopologyWorld {
		coord = g.WrapCoord(coord)
	}
	
	if !g.IsValid(coord) {
		return
	}
	
	col, row := coord.ToOffset()
	g.tiles[row][col] = value
}

// AllCoords returns all valid coordinates in the grid
func (g *Grid) AllCoords() []AxialCoord {
	coords := make([]AxialCoord, 0, g.config.Width*g.config.Height)
	
	for row := 0; row < g.config.Height; row++ {
		for col := 0; col < g.config.Width; col++ {
			coord := OffsetToAxial(col, row)
			coords = append(coords, coord)
		}
	}
	
	return coords
}

// hexDirections are the 6 directions from any hex to its neighbors
var hexDirections = [6]AxialCoord{
	{1, 0}, {1, -1}, {0, -1}, {-1, 0}, {-1, 1}, {0, 1},
}

// Neighbors returns all valid neighbors of a coordinate based on grid topology
func (c AxialCoord) Neighbors(grid *Grid) []AxialCoord {
	neighbors := make([]AxialCoord, 0, 6)
	
	for _, direction := range hexDirections {
		neighbor := AxialCoord{
			Q: c.Q + direction.Q,
			R: c.R + direction.R,
		}
		
		if grid.config.Topology == TopologyWorld {
			// In world topology, all neighbors are valid (after wrapping)
			wrapped := grid.WrapCoord(neighbor)
			neighbors = append(neighbors, wrapped)
		} else {
			// In region topology, only add if the neighbor is valid
			if grid.IsValid(neighbor) {
				neighbors = append(neighbors, neighbor)
			}
		}
	}
	
	return neighbors
}

// IsEdgeHex returns true if the coordinate is on the edge of a region map
// For world maps, no hex is considered an "edge"
func (c AxialCoord) IsEdgeHex(grid *Grid) bool {
	if grid.config.Topology == TopologyWorld {
		return false
	}
	
	// A hex is an edge hex if it has fewer than 6 neighbors
	neighbors := c.Neighbors(grid)
	return len(neighbors) < 6
}

// DistanceTo calculates the distance between two coordinates
// For world topology, considers wrapping for shortest path
func (c AxialCoord) DistanceTo(other AxialCoord, grid *Grid) int {
	if grid.config.Topology == TopologyRegion {
		// Standard hex distance for region topology
		return hexDistance(c, other)
	}
	
	// For world topology, consider wrapped distances
	minDist := hexDistance(c, other)
	
	// Try all possible wrapped versions of 'other'
	for dq := -1; dq <= 1; dq++ {
		for dr := -1; dr <= 1; dr++ {
			wrappedOther := AxialCoord{
				Q: other.Q + dq*grid.config.Width,
				R: other.R + dr*grid.config.Height,
			}
			dist := hexDistance(c, wrappedOther)
			if dist < minDist {
				minDist = dist
			}
		}
	}
	
	return minDist
}

// hexDistance calculates the standard hex distance between two coordinates
func hexDistance(a, b AxialCoord) int {
	return (abs(a.Q-b.Q) + abs(a.Q+a.R-b.Q-b.R) + abs(a.R-b.R)) / 2
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ShortestPath returns the shortest path between two coordinates
// For world maps, considers wrapping
func (g *Grid) ShortestPath(from, to AxialCoord) []AxialCoord {
	if g.config.Topology == TopologyRegion {
		return hexPathRegion(from, to)
	}
	
	// For world topology, find the wrapped version of 'to' that gives shortest distance
	bestTo := to
	minDist := hexDistance(from, to)
	
	// Try all possible wrapped versions of 'to' - need to check more offsets
	for dCol := -1; dCol <= 1; dCol++ {
		for dRow := -1; dRow <= 1; dRow++ {
			// Create wrapped target in offset space then convert to axial
			toCol, toRow := to.ToOffset()
			wrappedCol := toCol + dCol*g.config.Width
			wrappedRow := toRow + dRow*g.config.Height
			wrappedTo := OffsetToAxial(wrappedCol, wrappedRow)
			
			dist := hexDistance(from, wrappedTo)
			if dist < minDist {
				minDist = dist
				bestTo = wrappedTo
			}
		}
	}
	
	// Generate path to best target, then wrap coordinates back to valid range
	path := hexPathRegion(from, bestTo)
	for i := range path {
		path[i] = g.WrapCoord(path[i])
	}
	
	return path
}

// hexPathRegion generates a simple path between two coordinates (without wrapping)
func hexPathRegion(from, to AxialCoord) []AxialCoord {
	distance := hexDistance(from, to)
	if distance == 0 {
		return []AxialCoord{from}
	}
	
	path := make([]AxialCoord, distance+1)
	path[0] = from
	path[distance] = to
	
	// Simple linear interpolation path
	for i := 1; i < distance; i++ {
		t := float64(i) / float64(distance)
		q := float64(from.Q)*(1-t) + float64(to.Q)*t
		r := float64(from.R)*(1-t) + float64(to.R)*t
		path[i] = axialRound(q, r)
	}
	
	return path
}