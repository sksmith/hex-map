package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sean/hex-map/pkg/hex"
	"github.com/sean/hex-map/pkg/render"
	"github.com/sean/hex-map/pkg/terrain"
)

func handleRender(args []string) {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	input := fs.String("input", "", "Input terrain JSON file")
	output := fs.String("output", "terrain_render.jpg", "Output image filename")
	mode := fs.String("mode", "elevation", "Render mode: elevation, debug, hillshade")
	width := fs.Int("width", 800, "Image width in pixels")
	height := fs.Int("height", 600, "Image height in pixels")
	hexSize := fs.Float64("hex-size", 5.0, "Hex size in pixels")
	quality := fs.Int("quality", 85, "JPEG quality (1-100)")
	scheme := fs.String("scheme", "elevation", "Color scheme: elevation, realistic, debug, grayscale")
	showCoords := fs.Bool("show-coords", false, "Show coordinate debug overlay")

	fs.Parse(args)

	if *input == "" {
		fmt.Println("Error: --input is required")
		fmt.Println("Usage: hex-world render --input=terrain.json --output=image.jpg")
		return
	}

	// Load terrain data
	terrainData, err := loadTerrainData(*input)
	if err != nil {
		fmt.Printf("Error loading terrain data: %v\n", err)
		return
	}

	// Parse color scheme
	var colorScheme render.ColorScheme
	switch *scheme {
	case "elevation":
		colorScheme = render.SchemeElevation
	case "realistic":
		colorScheme = render.SchemeRealistic
	case "debug":
		colorScheme = render.SchemeDebug
	case "grayscale":
		colorScheme = render.SchemeGrayscale
	default:
		fmt.Printf("Error: unknown color scheme '%s'\n", *scheme)
		return
	}

	// Determine layers based on mode
	var layers []render.RenderLayer
	switch *mode {
	case "elevation":
		layers = []render.RenderLayer{render.LayerElevation}
	case "debug":
		layers = []render.RenderLayer{render.LayerElevation}
		if *showCoords {
			layers = append(layers, render.LayerDebugCoords)
		}
	case "hillshade":
		layers = []render.RenderLayer{render.LayerElevation} // TODO: Add hillshading
	default:
		fmt.Printf("Error: unknown render mode '%s'\n", *mode)
		return
	}

	// Create grid configuration
	gridConfig := hex.GridConfig{
		Width:    50, // Default size, could be derived from terrain data
		Height:   50,
		Topology: hex.TopologyRegion,
	}
	grid := hex.NewGrid(gridConfig)

	// Create render configuration
	renderConfig := render.RenderConfig{
		Width:       *width,
		Height:      *height,
		HexSize:     *hexSize,
		Layers:      layers,
		ColorScheme: colorScheme,
		ShowDebug:   *showCoords,
		Quality:     *quality,
	}

	// Create renderer
	renderer := render.NewHexRenderer(grid, renderConfig)

	fmt.Printf("Rendering %s terrain (%d tiles)...\n", *mode, len(terrainData.Tiles))

	// Render terrain
	_, err = renderer.RenderTerrain(terrainData.Tiles)
	if err != nil {
		fmt.Printf("Error rendering terrain: %v\n", err)
		return
	}

	// Export image
	if strings.HasSuffix(*output, ".png") {
		err = renderer.ExportPNG(*output)
	} else {
		err = renderer.ExportJPEG(*output, *quality)
	}

	if err != nil {
		fmt.Printf("Error exporting image: %v\n", err)
		return
	}

	fmt.Printf("Image saved to %s\n", *output)
}

func handleRenderWithMetadata(args []string) {
	fs := flag.NewFlagSet("render-with-metadata", flag.ExitOnError)
	input := fs.String("input", "", "Input terrain JSON file")
	output := fs.String("output", "terrain_render.jpg", "Output image filename")
	width := fs.Int("width", 800, "Image width in pixels")
	height := fs.Int("height", 600, "Image height in pixels")
	hexSize := fs.Float64("hex-size", 5.0, "Hex size in pixels")
	quality := fs.Int("quality", 85, "JPEG quality (1-100)")

	fs.Parse(args)

	if *input == "" {
		fmt.Println("Error: --input is required")
		return
	}

	// Load terrain data
	terrainData, err := loadTerrainData(*input)
	if err != nil {
		fmt.Printf("Error loading terrain data: %v\n", err)
		return
	}

	// Create grid and renderer
	gridConfig := hex.GridConfig{Width: 50, Height: 50, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	renderConfig := render.RenderConfig{
		Width:       *width,
		Height:      *height,
		HexSize:     *hexSize,
		Layers:      []render.RenderLayer{render.LayerElevation},
		ColorScheme: render.SchemeElevation,
		Quality:     *quality,
	}

	renderer := render.NewHexRenderer(grid, renderConfig)

	// Render terrain
	_, err = renderer.RenderTerrain(terrainData.Tiles)
	if err != nil {
		fmt.Printf("Error rendering terrain: %v\n", err)
		return
	}

	// Create metadata
	metadata := render.RenderMetadata{
		Generator:    "hex-world v1.0",
		Timestamp:    time.Now().Format(time.RFC3339),
		WorldSeed:    terrainData.Config.Seed,
		Stage:        "terrain_visualization",
		ViewConfig:   renderConfig,
		TerrainStats: terrainData.Stats,
		QualityScore: 0.9, // TODO: Calculate actual quality score
		KnownIssues:  []string{},
	}

	// Export with metadata
	if strings.HasSuffix(*output, ".png") {
		err = renderer.ExportPNGWithMetadata(*output, metadata)
	} else {
		err = renderer.ExportJPEGWithMetadata(*output, metadata)
	}

	if err != nil {
		fmt.Printf("Error exporting image with metadata: %v\n", err)
		return
	}

	fmt.Printf("Image with metadata saved to %s\n", *output)
}

func handleDemoRender(args []string) {
	fs := flag.NewFlagSet("demo-render", flag.ExitOnError)
	size := fs.String("size", "20x20", "Grid size as WIDTHxHEIGHT")
	seed := fs.Int64("seed", 42, "Random seed for terrain generation")
	outputDir := fs.String("output-dir", "demo_renders", "Output directory for demo images")

	fs.Parse(args)

	// Parse grid size
	width, height, err := parseSize(*size)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Create output directory
	err = os.MkdirAll(*outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	fmt.Printf("Rendering Demo - %dx%d grid (seed: %d)\n", width, height, *seed)
	fmt.Println(strings.Repeat("=", 50))

	// Generate terrain
	gridConfig := hex.GridConfig{Width: width, Height: height, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	terrainConfig := terrain.DefaultTerrainConfig()
	terrainConfig.Seed = *seed

	fmt.Println("Generating terrain...")
	tiles, err := terrain.GenerateTerrain(grid, terrainConfig)
	if err != nil {
		fmt.Printf("Error generating terrain: %v\n", err)
		return
	}

	// Common render config
	baseRenderConfig := render.RenderConfig{
		Width:   600,
		Height:  600,
		HexSize: 8.0,
		Quality: 90,
	}

	// Demo different color schemes
	schemes := []struct {
		name   string
		scheme render.ColorScheme
	}{
		{"elevation", render.SchemeElevation},
		{"realistic", render.SchemeRealistic},
		{"debug", render.SchemeDebug},
		{"grayscale", render.SchemeGrayscale},
	}

	for _, s := range schemes {
		fmt.Printf("Rendering %s color scheme...\n", s.name)

		renderConfig := baseRenderConfig
		renderConfig.ColorScheme = s.scheme
		renderConfig.Layers = []render.RenderLayer{render.LayerElevation}

		renderer := render.NewHexRenderer(grid, renderConfig)

		_, err = renderer.RenderTerrain(tiles)
		if err != nil {
			fmt.Printf("Error rendering %s: %v\n", s.name, err)
			continue
		}

		filename := fmt.Sprintf("%s/terrain_%s.jpg", *outputDir, s.name)
		err = renderer.ExportJPEG(filename, 90)
		if err != nil {
			fmt.Printf("Error exporting %s: %v\n", s.name, err)
			continue
		}

		fmt.Printf("  Saved: %s\n", filename)
	}

	// Demo with debug coordinates
	fmt.Println("Rendering debug version with coordinates...")
	debugConfig := baseRenderConfig
	debugConfig.ColorScheme = render.SchemeElevation
	debugConfig.Layers = []render.RenderLayer{render.LayerElevation, render.LayerDebugCoords}
	debugConfig.ShowDebug = true

	debugRenderer := render.NewHexRenderer(grid, debugConfig)
	_, err = debugRenderer.RenderTerrain(tiles)
	if err != nil {
		fmt.Printf("Error rendering debug version: %v\n", err)
	} else {
		debugFilename := fmt.Sprintf("%s/terrain_debug.png", *outputDir)
		err = debugRenderer.ExportPNG(debugFilename)
		if err != nil {
			fmt.Printf("Error exporting debug version: %v\n", err)
		} else {
			fmt.Printf("  Saved: %s\n", debugFilename)
		}
	}

	fmt.Println("\nDemo complete! Check the output directory for rendered images.")
}

// Helper function to load terrain data from JSON
func loadTerrainData(filename string) (*struct {
	Config terrain.TerrainConfig `json:"config"`
	Stats  terrain.TerrainStats  `json:"stats"`
	Tiles  []*terrain.HexTile    `json:"tiles"`
}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var terrainData struct {
		Config terrain.TerrainConfig `json:"config"`
		Stats  terrain.TerrainStats  `json:"stats"`
		Tiles  []*terrain.HexTile    `json:"tiles"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&terrainData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &terrainData, nil
}
