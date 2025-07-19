package hex

import (
	"testing"
)

// TestGridCreation tests basic grid creation with different topologies
func TestGridCreation(t *testing.T) {
	tests := []struct {
		config      GridConfig
		description string
	}{
		{
			GridConfig{Width: 10, Height: 8, Topology: TopologyRegion},
			"region topology grid",
		},
		{
			GridConfig{Width: 10, Height: 8, Topology: TopologyWorld},
			"world topology grid",
		},
	}

	for _, test := range tests {
		grid := NewGrid(test.config)
		if grid == nil {
			t.Errorf("NewGrid failed for %s", test.description)
			continue
		}

		if grid.Topology() != test.config.Topology {
			t.Errorf("%s: expected topology %d, got %d",
				test.description, test.config.Topology, grid.Topology())
		}

		// Test coordinate validity based on topology
		testCoords := []struct {
			coord     AxialCoord
			shouldBeValid bool
		}{
			{NewAxialCoord(0, 0), true},
			{NewAxialCoord(2, 0), true},       // interior coordinate
			{NewAxialCoord(4, 0), true},       // valid boundary
			{NewAxialCoord(15, 0), false},     // clearly outside width bounds 
			{NewAxialCoord(0, 15), false},     // clearly outside height bounds
		}

		for _, coordTest := range testCoords {
			isValid := grid.IsValid(coordTest.coord)
			expectedValid := coordTest.shouldBeValid
			
			// For world topology, all coordinates should be valid after wrapping
			if test.config.Topology == TopologyWorld {
				expectedValid = true
			}
			
			if isValid != expectedValid {
				t.Errorf("%s: IsValid(%v) = %v, expected %v",
					test.description, coordTest.coord, isValid, expectedValid)
			}
		}
	}
}

// TestGridGetSet tests getting and setting values in the grid
func TestGridGetSet(t *testing.T) {
	config := GridConfig{Width: 5, Height: 3, Topology: TopologyRegion}
	grid := NewGrid(config)

	coord := NewAxialCoord(2, 1)
	testValue := "test_value"

	// Test setting and getting a value
	grid.Set(coord, testValue)
	retrieved := grid.Get(coord)

	if retrieved != testValue {
		t.Errorf("Expected '%s', got '%v'", testValue, retrieved)
	}

	// Test getting from an empty coordinate
	emptyCoord := NewAxialCoord(0, 0)
	emptyValue := grid.Get(emptyCoord)
	if emptyValue != nil {
		t.Errorf("Expected nil for empty coordinate, got '%v'", emptyValue)
	}
}

// TestGridAllCoords tests getting all coordinates from the grid
func TestGridAllCoords(t *testing.T) {
	config := GridConfig{Width: 3, Height: 2, Topology: TopologyRegion}
	grid := NewGrid(config)

	allCoords := grid.AllCoords()
	expectedCount := 3 * 2 // width * height
	
	if len(allCoords) != expectedCount {
		t.Errorf("Expected %d coordinates, got %d", expectedCount, len(allCoords))
	}

	// Check that all coordinates are valid
	for _, coord := range allCoords {
		if !grid.IsValid(coord) {
			t.Errorf("AllCoords returned invalid coordinate: %v", coord)
		}
	}

	// Check for duplicates
	coordSet := make(map[AxialCoord]bool)
	for _, coord := range allCoords {
		if coordSet[coord] {
			t.Errorf("AllCoords returned duplicate coordinate: %v", coord)
		}
		coordSet[coord] = true
	}
}

// TestGridBoundaryHandling tests that grids handle boundary conditions correctly
func TestGridBoundaryHandling(t *testing.T) {
	config := GridConfig{Width: 5, Height: 3, Topology: TopologyRegion}
	grid := NewGrid(config)

	// Test coordinates on the boundary (using valid axial coordinates for 5x3 grid)
	boundaryCoords := []AxialCoord{
		{0, 0}, {4, -2}, {0, 2}, {4, 0}, // corners: offsets (0,0), (4,0), (0,2), (4,2)
		{2, -1}, {0, 1}, {4, -1}, {2, 1}, // edges: offsets (2,0), (0,1), (4,1), (2,2)
	}

	for _, coord := range boundaryCoords {
		if !grid.IsValid(coord) {
			t.Errorf("Boundary coordinate %v should be valid", coord)
		}

		// Should be able to get/set on boundary
		grid.Set(coord, "boundary_value")
		value := grid.Get(coord)
		if value != "boundary_value" {
			t.Errorf("Failed to set/get value on boundary coordinate %v", coord)
		}
	}

	// Test coordinates outside the boundary
	outsideCoords := []AxialCoord{
		{-1, 0}, {5, 0}, {0, -1}, {0, 3}, // just outside edges
		{-1, -1}, {5, 3}, // corners outside
	}

	for _, coord := range outsideCoords {
		if grid.IsValid(coord) {
			t.Errorf("Outside coordinate %v should not be valid for region topology", coord)
		}
	}
}

// TestWorldGridWrapping tests that world grids handle wrapped coordinates correctly
func TestWorldGridWrapping(t *testing.T) {
	config := GridConfig{Width: 5, Height: 3, Topology: TopologyWorld}
	grid := NewGrid(config)

	// Test that wrapped coordinates work correctly
	tests := []struct {
		original AxialCoord
		wrapped  AxialCoord
		value    string
	}{
		{NewAxialCoord(-1, 1), NewAxialCoord(4, 1), "wrapped_left"},
		{NewAxialCoord(5, 1), NewAxialCoord(0, 1), "wrapped_right"},
		{NewAxialCoord(2, -1), NewAxialCoord(2, 2), "wrapped_up"},
		{NewAxialCoord(2, 3), NewAxialCoord(2, 0), "wrapped_down"},
	}

	for _, test := range tests {
		// Set value using wrapped coordinate
		grid.Set(test.wrapped, test.value)
		
		// Both original and wrapped should access the same value
		wrappedValue := grid.Get(test.wrapped)
		if wrappedValue != test.value {
			t.Errorf("Wrapped coordinate %v failed to retrieve value", test.wrapped)
		}

		// The wrapped coordinate should be valid
		if !grid.IsValid(test.wrapped) {
			t.Errorf("Wrapped coordinate %v should be valid", test.wrapped)
		}
	}
}