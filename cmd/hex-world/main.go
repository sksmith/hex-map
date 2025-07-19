package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sean/hex-map/pkg/hex"
	"github.com/sean/hex-map/pkg/terrain"
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
	case "generate-terrain":
		handleGenerateTerrain(os.Args[2:])
	case "terrain-stats":
		handleTerrainStats(os.Args[2:])
	case "validate-terrain":
		handleValidateTerrain(os.Args[2:])
	case "demo-terrain":
		handleDemoTerrain(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("hex-world - Hex Map World Generation Tool")
	fmt.Println("")
	fmt.Println("Hex Grid Commands:")
	fmt.Println("  demo-coords     --size=WxH --topology=TYPE              Show coordinate system demo")
	fmt.Println("  demo-distance   --from=Q,R --to=Q,R --topology=TYPE     Show distance calculation")
	fmt.Println("")
	fmt.Println("Terrain Generation Commands:")
	fmt.Println("  generate-terrain --size=WxH --seed=N --output=FILE      Generate terrain and save to JSON")
	fmt.Println("  terrain-stats   FILE.json                               Show terrain statistics")
	fmt.Println("  validate-terrain FILE.json [--strict]                   Validate terrain realism")
	fmt.Println("  demo-terrain    --size=WxH [--seed=N]                    Quick terrain demo with stats")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --topology=TYPE     region (bounded) or world (toroidal)")
	fmt.Println("  --size=WxH          Grid dimensions (e.g., 100x100)")
	fmt.Println("  --seed=N            Random seed for reproducible generation")
	fmt.Println("  --output=FILE       Output filename for JSON data")
	fmt.Println("  --land-ratio=N      Target land percentage (0.0-1.0, default: 0.29)")
	fmt.Println("  --sea-level=N       Sea level in meters (default: 0)")
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

// Terrain generation commands

func handleGenerateTerrain(args []string) {
	fs := flag.NewFlagSet("generate-terrain", flag.ExitOnError)
	size := fs.String("size", "100x100", "Grid size as WIDTHxHEIGHT")
	seed := fs.Int64("seed", 42, "Random seed for terrain generation")
	output := fs.String("output", "terrain.json", "Output filename for JSON data")
	topology := fs.String("topology", "region", "Topology type: region or world")
	landRatio := fs.Float64("land-ratio", 0.29, "Target land percentage (0.0-1.0)")
	seaLevel := fs.Float64("sea-level", 0.0, "Sea level in meters")
	
	fs.Parse(args)
	
	// Parse grid size
	width, height, err := parseSize(*size)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Parse topology
	topo, err := parseTopology(*topology)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Create grid
	gridConfig := hex.GridConfig{Width: width, Height: height, Topology: topo}
	grid := hex.NewGrid(gridConfig)
	
	// Configure terrain generation
	terrainConfig := terrain.TerrainConfig{
		Seed:        *seed,
		SeaLevel:    *seaLevel,
		LandRatio:   *landRatio,
		NoiseParams: terrain.DefaultNoiseParameters(),
	}
	
	fmt.Printf("Generating %dx%d terrain (seed: %d)...\n", width, height, *seed)
	
	// Generate terrain
	tiles, err := terrain.GenerateTerrain(grid, terrainConfig)
	if err != nil {
		fmt.Printf("Error generating terrain: %v\n", err)
		return
	}
	
	// Calculate statistics
	stats := terrain.ValidateTerrain(tiles)
	
	// Save to JSON
	terrainData := struct {
		Config terrain.TerrainConfig `json:"config"`
		Stats  terrain.TerrainStats  `json:"stats"`
		Tiles  []*terrain.HexTile    `json:"tiles"`
	}{
		Config: terrainConfig,
		Stats:  stats,
		Tiles:  tiles,
	}
	
	file, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(terrainData); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}
	
	fmt.Printf("Terrain saved to %s\n", *output)
	fmt.Printf("Land coverage: %.1f%% (%d/%d tiles)\n", 
		stats.LandPercentage, stats.LandTiles, stats.TotalTiles)
	fmt.Printf("Elevation range: %.1fm to %.1fm\n", 
		stats.ElevationRange[0], stats.ElevationRange[1])
}

func handleTerrainStats(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Please provide a terrain JSON file")
		fmt.Println("Usage: hex-world terrain-stats FILE.json")
		return
	}
	
	filename := args[0]
	
	// Load terrain data
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	
	var terrainData struct {
		Config terrain.TerrainConfig `json:"config"`
		Stats  terrain.TerrainStats  `json:"stats"`
		Tiles  []*terrain.HexTile    `json:"tiles"`
	}
	
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&terrainData); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	
	// Display comprehensive statistics
	stats := terrainData.Stats
	config := terrainData.Config
	
	fmt.Printf("Terrain Statistics for %s\n", filename)
	fmt.Println(strings.Repeat("=", 50))
	
	fmt.Println("Generation Parameters:")
	fmt.Printf("  Seed: %d\n", config.Seed)
	fmt.Printf("  Sea Level: %.1fm\n", config.SeaLevel)
	fmt.Printf("  Target Land Ratio: %.1f%%\n", config.LandRatio*100)
	fmt.Printf("  Noise Octaves: %d\n", config.NoiseParams.Octaves)
	fmt.Printf("  Persistence: %.2f\n", config.NoiseParams.Persistence)
	
	fmt.Println("\nElevation Statistics:")
	fmt.Printf("  Range: %.1fm to %.1fm (span: %.1fm)\n", 
		stats.ElevationRange[0], stats.ElevationRange[1], 
		stats.ElevationRange[1]-stats.ElevationRange[0])
	fmt.Printf("  Mean: %.1fm\n", stats.ElevationMean)
	fmt.Printf("  Standard Deviation: %.1fm\n", stats.ElevationStdDev)
	
	fmt.Println("\nLand/Water Distribution:")
	fmt.Printf("  Total Tiles: %d\n", stats.TotalTiles)
	fmt.Printf("  Land: %d tiles (%.1f%%)\n", stats.LandTiles, stats.LandPercentage)
	fmt.Printf("  Water: %d tiles (%.1f%%)\n", stats.WaterTiles, stats.WaterPercentage)
	
	fmt.Println("\nQuality Metrics:")
	fmt.Printf("  Hypsometric Match: %.1f%% (Earth-like curve)\n", stats.HypsometricMatch*100)
	
	// Check realism
	isRealistic, issues := terrain.IsRealisticTerrain(stats)
	if isRealistic {
		fmt.Println("  Realism Check: ✅ PASS")
	} else {
		fmt.Println("  Realism Check: ❌ FAIL")
		for _, issue := range issues {
			fmt.Printf("    - %s\n", issue)
		}
	}
}

func handleValidateTerrain(args []string) {
	fs := flag.NewFlagSet("validate-terrain", flag.ExitOnError)
	strict := fs.Bool("strict", false, "Use strict validation criteria")
	
	fs.Parse(args)
	
	if len(fs.Args()) == 0 {
		fmt.Println("Error: Please provide a terrain JSON file")
		fmt.Println("Usage: hex-world validate-terrain FILE.json [--strict]")
		return
	}
	
	filename := fs.Args()[0]
	
	// Load terrain data
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	
	var terrainData struct {
		Tiles []*terrain.HexTile `json:"tiles"`
	}
	
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&terrainData); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	
	fmt.Printf("Validating terrain from %s\n", filename)
	fmt.Println(strings.Repeat("=", 40))
	
	// Run validation
	stats := terrain.ValidateTerrain(terrainData.Tiles)
	isRealistic, issues := terrain.IsRealisticTerrain(stats)
	
	// Detect anomalies
	anomalies := terrain.DetectElevationAnomalies(terrainData.Tiles)
	
	// Report results
	fmt.Printf("Total tiles validated: %d\n", len(terrainData.Tiles))
	
	if isRealistic && len(anomalies) == 0 {
		fmt.Println("Status: ✅ VALID - Terrain passes all realism checks")
	} else {
		fmt.Println("Status: ❌ INVALID - Issues detected")
		
		if !isRealistic {
			fmt.Println("\nRealism Issues:")
			for _, issue := range issues {
				fmt.Printf("  - %s\n", issue)
			}
		}
		
		if len(anomalies) > 0 {
			fmt.Println("\nElevation Anomalies:")
			for _, anomaly := range anomalies {
				fmt.Printf("  - %s\n", anomaly)
			}
		}
	}
	
	// In strict mode, additional checks
	if *strict {
		fmt.Println("\nStrict Mode Validation:")
		
		// Check hypsometric curve match
		if stats.HypsometricMatch < 0.95 {
			fmt.Printf("  ❌ Hypsometric curve match too low: %.1f%% (strict requires >95%%)\n", 
				stats.HypsometricMatch*100)
		} else {
			fmt.Println("  ✅ Hypsometric curve match acceptable")
		}
		
		// Check land ratio precision
		targetLandRatio := 29.0 // Earth's land percentage
		landRatioDiff := abs(int(stats.LandPercentage - targetLandRatio))
		if landRatioDiff > 1 {
			fmt.Printf("  ❌ Land ratio deviation too high: %.1f%% (target: %.1f%%)\n", 
				stats.LandPercentage, targetLandRatio)
		} else {
			fmt.Println("  ✅ Land ratio within acceptable range")
		}
	}
}

func handleDemoTerrain(args []string) {
	fs := flag.NewFlagSet("demo-terrain", flag.ExitOnError)
	size := fs.String("size", "50x50", "Grid size as WIDTHxHEIGHT")
	seed := fs.Int64("seed", 42, "Random seed for terrain generation")
	topology := fs.String("topology", "region", "Topology type: region or world")
	
	fs.Parse(args)
	
	// Parse grid size
	width, height, err := parseSize(*size)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Parse topology
	topo, err := parseTopology(*topology)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Terrain Generation Demo - %dx%d grid (seed: %d)\n", width, height, *seed)
	fmt.Println(strings.Repeat("=", 50))
	
	// Create grid
	gridConfig := hex.GridConfig{Width: width, Height: height, Topology: topo}
	grid := hex.NewGrid(gridConfig)
	
	// Generate terrain with default config
	terrainConfig := terrain.DefaultTerrainConfig()
	terrainConfig.Seed = *seed
	
	fmt.Println("Generating terrain...")
	tiles, err := terrain.GenerateTerrain(grid, terrainConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Analyze results
	stats := terrain.ValidateTerrain(tiles)
	isRealistic, issues := terrain.IsRealisticTerrain(stats)
	
	fmt.Println("\nGeneration Results:")
	fmt.Printf("  Total tiles: %d\n", stats.TotalTiles)
	fmt.Printf("  Land coverage: %.1f%% (%d tiles)\n", stats.LandPercentage, stats.LandTiles)
	fmt.Printf("  Water coverage: %.1f%% (%d tiles)\n", stats.WaterPercentage, stats.WaterTiles)
	
	fmt.Println("\nElevation Analysis:")
	fmt.Printf("  Range: %.0fm to %.0fm\n", stats.ElevationRange[0], stats.ElevationRange[1])
	fmt.Printf("  Mean: %.0fm\n", stats.ElevationMean)
	fmt.Printf("  Std Dev: %.0fm\n", stats.ElevationStdDev)
	
	fmt.Println("\nQuality Assessment:")
	fmt.Printf("  Hypsometric Match: %.1f%%\n", stats.HypsometricMatch*100)
	if isRealistic {
		fmt.Println("  Realism Check: ✅ PASS")
	} else {
		fmt.Println("  Realism Check: ❌ FAIL")
		for _, issue := range issues {
			fmt.Printf("    - %s\n", issue)
		}
	}
	
	// Show a few sample tiles
	fmt.Println("\nSample Terrain Tiles:")
	fmt.Println("Coordinate  | Elevation | Type | Depth/Height")
	fmt.Println("------------|-----------|------|-------------")
	
	sampleIndices := []int{0, len(tiles)/4, len(tiles)/2, 3*len(tiles)/4, len(tiles)-1}
	for _, i := range sampleIndices {
		if i < len(tiles) {
			tile := tiles[i]
			tileType := "Water"
			depthHeight := fmt.Sprintf("%.0fm deep", tile.GetDepth(0))
			
			if tile.IsLand {
				tileType = "Land"
				depthHeight = fmt.Sprintf("%.0fm high", tile.GetHeight(0))
			}
			
			fmt.Printf("(%2d,%2d)      | %8.0f  | %-5s | %s\n",
				tile.Coordinates.Q, tile.Coordinates.R, tile.Elevation, tileType, depthHeight)
		}
	}
}

// Helper functions

func parseSize(sizeStr string) (int, int, error) {
	parts := strings.Split(sizeStr, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("size must be in format WIDTHxHEIGHT (e.g., 100x100)")
	}
	
	width, err1 := strconv.Atoi(parts[0])
	height, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("invalid size format")
	}
	
	if width <= 0 || height <= 0 {
		return 0, 0, fmt.Errorf("size dimensions must be positive")
	}
	
	return width, height, nil
}

func parseTopology(topologyStr string) (hex.Topology, error) {
	switch topologyStr {
	case "region":
		return hex.TopologyRegion, nil
	case "world":
		return hex.TopologyWorld, nil
	default:
		return hex.TopologyRegion, fmt.Errorf("unknown topology '%s'. Use 'region' or 'world'", topologyStr)
	}
}