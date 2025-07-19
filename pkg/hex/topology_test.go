package hex

import (
	"testing"
)

// TestRegionTopologyNeighbors tests neighbor lookup for bounded region maps
func TestRegionTopologyNeighbors(t *testing.T) {
	config := GridConfig{Width: 5, Height: 3, Topology: TopologyRegion}
	grid := NewGrid(config)

	tests := []struct {
		coord           AxialCoord
		expectedCount   int
		description     string
		shouldBeEdge    bool
	}{
		// Corner hexes have 2-3 neighbors (based on actual grid layout)
		{NewAxialCoord(0, 0), 3, "top-left corner", true},       // offset (0,0)
		{NewAxialCoord(4, -2), 3, "top-right corner", true},     // offset (4,0)
		{NewAxialCoord(0, 2), 2, "bottom-left corner", true},    // offset (0,2)
		{NewAxialCoord(4, 0), 2, "bottom-right corner", true},   // offset (4,2)
		
		// Edge hexes have 3-4 neighbors
		{NewAxialCoord(1, -1), 3, "top edge", true},             // offset (1,0)
		{NewAxialCoord(0, 1), 4, "left edge", true},             // offset (0,1)
		{NewAxialCoord(4, -1), 4, "right edge", true},           // offset (4,1)
		{NewAxialCoord(2, 1), 3, "bottom edge", true},           // offset (2,2)
		
		// Interior hexes have 6 neighbors
		{NewAxialCoord(2, 0), 6, "interior hex", false},         // offset (2,1)
	}

	for _, test := range tests {
		neighbors := test.coord.Neighbors(grid)
		if len(neighbors) != test.expectedCount {
			t.Errorf("%s at %v: expected %d neighbors, got %d",
				test.description, test.coord, test.expectedCount, len(neighbors))
		}

		isEdge := test.coord.IsEdgeHex(grid)
		if isEdge != test.shouldBeEdge {
			t.Errorf("%s at %v: expected edge=%v, got edge=%v",
				test.description, test.coord, test.shouldBeEdge, isEdge)
		}

		// Verify all neighbors are valid coordinates
		for _, neighbor := range neighbors {
			if !grid.IsValid(neighbor) {
				t.Errorf("%s at %v: invalid neighbor %v",
					test.description, test.coord, neighbor)
			}
		}
	}
}

// TestWorldTopologyNeighbors tests neighbor lookup for toroidal world maps
func TestWorldTopologyNeighbors(t *testing.T) {
	config := GridConfig{Width: 5, Height: 3, Topology: TopologyWorld}
	grid := NewGrid(config)

	// All hexes in world topology should have exactly 6 neighbors
	coords := []AxialCoord{
		{0, 0}, {4, 0}, {0, 2}, {4, 2}, // corners
		{1, 0}, {0, 1}, {4, 1}, {2, 2}, // edges
		{2, 1}, // interior
	}

	for _, coord := range coords {
		neighbors := coord.Neighbors(grid)
		if len(neighbors) != 6 {
			t.Errorf("World topology at %v: expected 6 neighbors, got %d",
				coord, len(neighbors))
		}

		// In world topology, no hex is considered an "edge"
		if coord.IsEdgeHex(grid) {
			t.Errorf("World topology at %v: no hex should be edge in world topology",
				coord)
		}

		// Verify all neighbors are valid after wrapping
		for _, neighbor := range neighbors {
			wrapped := grid.WrapCoord(neighbor)
			if !grid.IsValid(wrapped) {
				t.Errorf("World topology at %v: invalid wrapped neighbor %v → %v",
					coord, neighbor, wrapped)
			}
		}
	}
}

// TestCoordinateWrapping tests coordinate wrapping for world maps
func TestCoordinateWrapping(t *testing.T) {
	config := GridConfig{Width: 5, Height: 3, Topology: TopologyWorld}
	grid := NewGrid(config)

	tests := []struct {
		input    AxialCoord
		expected AxialCoord
	}{
		// No wrapping needed
		{NewAxialCoord(2, 1), NewAxialCoord(2, 1)},
		
		// Horizontal wrapping 
		{NewAxialCoord(-1, 1), NewAxialCoord(4, -1)},  // offset (-1,1) → (4,1)
		{NewAxialCoord(5, 1), NewAxialCoord(0, 1)},    // offset (5,4) → (0,1)
		{NewAxialCoord(6, 1), NewAxialCoord(1, 0)},    // offset (6,4) → (1,1)
		
		// Vertical wrapping
		{NewAxialCoord(2, -1), NewAxialCoord(2, -1)},  // offset (2,0) → (2,0) - already valid
		{NewAxialCoord(2, 3), NewAxialCoord(2, 0)},    // offset (2,4) → (2,1)
		{NewAxialCoord(2, 4), NewAxialCoord(2, 1)},    // offset (2,5) → (2,2)
		
		// Both coordinates need wrapping
		{NewAxialCoord(-1, -1), NewAxialCoord(4, 0)},  // offset (-1,-1) → (4,2)
		{NewAxialCoord(5, 3), NewAxialCoord(0, 0)},    // offset (5,6) → (0,0)
	}

	for _, test := range tests {
		result := grid.WrapCoord(test.input)
		if result.Q != test.expected.Q || result.R != test.expected.R {
			t.Errorf("WrapCoord(%v) = %v, expected %v",
				test.input, result, test.expected)
		}
	}
}

// TestDistanceCalculation tests distance calculations for both topologies
func TestDistanceCalculation(t *testing.T) {
	regionConfig := GridConfig{Width: 10, Height: 8, Topology: TopologyRegion}
	worldConfig := GridConfig{Width: 10, Height: 8, Topology: TopologyWorld}
	regionGrid := NewGrid(regionConfig)
	worldGrid := NewGrid(worldConfig)

	tests := []struct {
		from, to      AxialCoord
		regionDist    int
		worldDist     int
		description   string
	}{
		{
			NewAxialCoord(0, 0), NewAxialCoord(2, 1),
			3, 3, "same distance for both topologies (no wrapping benefit)",
		},
		{
			NewAxialCoord(0, 0), NewAxialCoord(9, 0),
			9, 1, "horizontal wrapping beneficial in world topology",
		},
		{
			NewAxialCoord(1, 0), NewAxialCoord(1, 7),
			7, 1, "vertical wrapping beneficial in world topology", 
		},
		{
			NewAxialCoord(0, 0), NewAxialCoord(9, 7),
			16, 2, "both wrappings beneficial in world topology",
		},
	}

	for _, test := range tests {
		regionDist := test.from.DistanceTo(test.to, regionGrid)
		worldDist := test.from.DistanceTo(test.to, worldGrid)

		if regionDist != test.regionDist {
			t.Errorf("%s: region distance from %v to %v = %d, expected %d",
				test.description, test.from, test.to, regionDist, test.regionDist)
		}

		if worldDist != test.worldDist {
			t.Errorf("%s: world distance from %v to %v = %d, expected %d",
				test.description, test.from, test.to, worldDist, test.worldDist)
		}
	}
}

// TestDistanceSymmetry tests that distance calculations are symmetric
func TestDistanceSymmetry(t *testing.T) {
	configs := []GridConfig{
		{Width: 10, Height: 8, Topology: TopologyRegion},
		{Width: 10, Height: 8, Topology: TopologyWorld},
	}

	testCoords := []AxialCoord{
		{0, 0}, {9, 0}, {0, 7}, {9, 7}, {5, 3},
	}

	for _, config := range configs {
		grid := NewGrid(config)
		topologyName := "region"
		if config.Topology == TopologyWorld {
			topologyName = "world"
		}

		for _, from := range testCoords {
			for _, to := range testCoords {
				distAB := from.DistanceTo(to, grid)
				distBA := to.DistanceTo(from, grid)
				if distAB != distBA {
					t.Errorf("%s topology: distance not symmetric %v→%v=%d, %v→%v=%d",
						topologyName, from, to, distAB, to, from, distBA)
				}
			}
		}
	}
}

// TestShortestPath tests pathfinding for world maps with wrapping
func TestShortestPath(t *testing.T) {
	config := GridConfig{Width: 5, Height: 3, Topology: TopologyWorld}
	grid := NewGrid(config)

	tests := []struct {
		from, to     AxialCoord
		maxPathLen   int
		description  string
	}{
		{
			NewAxialCoord(0, 0), NewAxialCoord(4, -2),  // offset (0,0) to (4,0) - should wrap
			2, "should wrap horizontally (distance 1, path length ≤ 2)",
		},
		{
			NewAxialCoord(2, 0), NewAxialCoord(2, -1),  // offset (2,1) to (2,0) - should wrap
			2, "should wrap vertically (distance 1, path length ≤ 2)",
		},
	}

	for _, test := range tests {
		path := grid.ShortestPath(test.from, test.to)
		if len(path) > test.maxPathLen {
			t.Errorf("%s: path length %d > max expected %d, path: %v",
				test.description, len(path), test.maxPathLen, path)
		}

		// Verify path connectivity
		if len(path) > 1 {
			for i := 0; i < len(path)-1; i++ {
				dist := path[i].DistanceTo(path[i+1], grid)
				if dist != 1 {
					t.Errorf("%s: path not connected at step %d: %v to %v (distance %d)",
						test.description, i, path[i], path[i+1], dist)
				}
			}
		}
	}
}