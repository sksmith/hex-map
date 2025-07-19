package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sean/hex-map/pkg/hex"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "demo-coords":
		handleDemoCoords(os.Args[2:])
	case "demo-distance":
		handleDemoDistance(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("hex-world - Hex Map World Generation Tool")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  demo-coords   --size=WxH --topology=TYPE    Show coordinate system demo")
	fmt.Println("  demo-distance --from=Q,R --to=Q,R --topology=TYPE   Show distance calculation")
	fmt.Println("")
	fmt.Println("Topology types: region (bounded) or world (toroidal)")
}

func handleDemoCoords(args []string) {
	fs := flag.NewFlagSet("demo-coords", flag.ExitOnError)
	size := fs.String("size", "10x8", "Grid size as WIDTHxHEIGHT")
	topology := fs.String("topology", "region", "Topology type: region or world")
	
	fs.Parse(args)
	
	// Parse size
	parts := strings.Split(*size, "x")
	if len(parts) != 2 {
		fmt.Println("Error: size must be in format WIDTHxHEIGHT (e.g., 10x8)")
		return
	}
	
	width, err1 := strconv.Atoi(parts[0])
	height, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		fmt.Println("Error: invalid size format")
		return
	}
	
	// Parse topology
	var topo hex.Topology
	switch *topology {
	case "region":
		topo = hex.TopologyRegion
	case "world":
		topo = hex.TopologyWorld
	default:
		fmt.Printf("Error: unknown topology '%s'. Use 'region' or 'world'\n", *topology)
		return
	}
	
	// Create grid and demonstrate
	config := hex.GridConfig{Width: width, Height: height, Topology: topo}
	grid := hex.NewGrid(config)
	
	fmt.Printf("Hex Grid Demo - %dx%d %s topology\n", width, height, *topology)
	fmt.Println(strings.Repeat("=", 50))
	
	// Show sample coordinates
	coords := grid.AllCoords()
	fmt.Printf("Total coordinates: %d\n", len(coords))
	
	// Show some sample coordinates with their properties
	sampleCoords := []hex.AxialCoord{
		coords[0],                    // first coordinate
		coords[len(coords)/2],       // middle coordinate  
		coords[len(coords)-1],       // last coordinate
	}
	
	fmt.Println("\nSample coordinates:")
	fmt.Println("Axial      | Offset  | Neighbors | Edge")
	fmt.Println("-----------|---------|-----------|-----")
	
	for _, coord := range sampleCoords {
		col, row := coord.ToOffset()
		neighbors := coord.Neighbors(grid)
		isEdge := coord.IsEdgeHex(grid)
		
		fmt.Printf("(%2d,%2d)    | (%d,%d)   | %d         | %v\n",
			coord.Q, coord.R, col, row, len(neighbors), isEdge)
	}
	
	// For world topology, show wrapping example
	if topo == hex.TopologyWorld {
		fmt.Println("\nWrapping examples:")
		wrapExamples := []hex.AxialCoord{
			hex.NewAxialCoord(-1, 0),
			hex.NewAxialCoord(width, 0),
			hex.NewAxialCoord(0, height),
		}
		
		for _, coord := range wrapExamples {
			wrapped := grid.WrapCoord(coord)
			col, row := coord.ToOffset()
			wCol, wRow := wrapped.ToOffset()
			fmt.Printf("(%2d,%2d) offset(%d,%d) → (%2d,%2d) offset(%d,%d)\n",
				coord.Q, coord.R, col, row, wrapped.Q, wrapped.R, wCol, wRow)
		}
	}
}

func handleDemoDistance(args []string) {
	fs := flag.NewFlagSet("demo-distance", flag.ExitOnError)
	fromStr := fs.String("from", "0,0", "Starting coordinate as Q,R")
	toStr := fs.String("to", "3,2", "Target coordinate as Q,R")
	topology := fs.String("topology", "region", "Topology type: region or world")
	
	fs.Parse(args)
	
	// Parse coordinates
	from, err := parseCoord(*fromStr)
	if err != nil {
		fmt.Printf("Error parsing 'from' coordinate: %v\n", err)
		return
	}
	
	to, err := parseCoord(*toStr)
	if err != nil {
		fmt.Printf("Error parsing 'to' coordinate: %v\n", err)
		return
	}
	
	// Parse topology
	var topo hex.Topology
	switch *topology {
	case "region":
		topo = hex.TopologyRegion
	case "world":
		topo = hex.TopologyWorld
	default:
		fmt.Printf("Error: unknown topology '%s'. Use 'region' or 'world'\n", *topology)
		return
	}
	
	// Create a reasonable grid size
	config := hex.GridConfig{Width: 10, Height: 8, Topology: topo}
	grid := hex.NewGrid(config)
	
	fmt.Printf("Distance Demo - %s topology\n", *topology)
	fmt.Println(strings.Repeat("=", 30))
	fmt.Printf("From: (%d,%d)\n", from.Q, from.R)
	fmt.Printf("To:   (%d,%d)\n", to.Q, to.R)
	
	// Calculate distance
	distance := from.DistanceTo(to, grid)
	fmt.Printf("Distance: %d hexes\n", distance)
	
	// Show path
	path := grid.ShortestPath(from, to)
	fmt.Printf("Path length: %d steps\n", len(path)-1)
	fmt.Println("Path:")
	for i, coord := range path {
		if i == 0 {
			fmt.Printf("  Start: (%d,%d)\n", coord.Q, coord.R)
		} else if i == len(path)-1 {
			fmt.Printf("  End:   (%d,%d)\n", coord.Q, coord.R)
		} else {
			fmt.Printf("  Step %d: (%d,%d)\n", i, coord.Q, coord.R)
		}
	}
	
	// For world topology, show if wrapping was used
	if topo == hex.TopologyWorld {
		directDistance := hexDistance(from, to)
		if distance < directDistance {
			fmt.Printf("\nWrapping used! Direct distance would be %d\n", directDistance)
		}
	}
}

func parseCoord(coordStr string) (hex.AxialCoord, error) {
	parts := strings.Split(coordStr, ",")
	if len(parts) != 2 {
		return hex.AxialCoord{}, fmt.Errorf("coordinate must be in format Q,R")
	}
	
	q, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	r, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return hex.AxialCoord{}, fmt.Errorf("invalid coordinate format")
	}
	
	return hex.NewAxialCoord(q, r), nil
}

// hexDistance calculates standard hex distance (duplicated here for demo)
func hexDistance(a, b hex.AxialCoord) int {
	return (abs(a.Q-b.Q) + abs(a.Q+a.R-b.Q-b.R) + abs(a.R-b.R)) / 2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}