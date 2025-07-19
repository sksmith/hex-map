package render

import (
	"image/color"
	"os"
	"testing"

	"github.com/sean/hex-map/pkg/hex"
	"github.com/sean/hex-map/pkg/terrain"
)

func BenchmarkElevationToColor(b *testing.B) {
	colorMap := TerrainColorScheme()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		elevation := float64(i%10000 - 5000) // Range from -5000 to 5000
		ElevationToColor(elevation, colorMap)
	}
}

func BenchmarkInterpolateColor(b *testing.B) {
	c1 := color.RGBA{255, 0, 0, 255}
	c2 := color.RGBA{0, 255, 0, 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ratio := float64(i%100) / 100.0
		InterpolateColor(c1, c2, ratio)
	}
}

func BenchmarkRenderSmallGrid(b *testing.B) {
	// Setup
	gridConfig := hex.GridConfig{Width: 10, Height: 10, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	tiles := make([]*terrain.HexTile, 0, 100)
	for _, coord := range grid.AllCoords() {
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   float64(coord.Q + coord.R*10),
			IsLand:      coord.Q%2 == 0,
		}
		tiles = append(tiles, tile)
	}

	renderConfig := RenderConfig{
		Width:       200,
		Height:      200,
		HexSize:     8.0,
		Layers:      []RenderLayer{LayerElevation},
		ColorScheme: SchemeElevation,
		Quality:     85,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer := NewHexRenderer(grid, renderConfig)
		_, err := renderer.RenderTerrain(tiles)
		if err != nil {
			b.Fatalf("Render failed: %v", err)
		}
	}
}

func BenchmarkRenderMediumGrid(b *testing.B) {
	// Setup
	gridConfig := hex.GridConfig{Width: 50, Height: 50, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	tiles := make([]*terrain.HexTile, 0, 2500)
	for _, coord := range grid.AllCoords() {
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   float64(coord.Q + coord.R*10),
			IsLand:      coord.Q%2 == 0,
		}
		tiles = append(tiles, tile)
	}

	renderConfig := RenderConfig{
		Width:       400,
		Height:      400,
		HexSize:     4.0,
		Layers:      []RenderLayer{LayerElevation},
		ColorScheme: SchemeElevation,
		Quality:     85,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer := NewHexRenderer(grid, renderConfig)
		_, err := renderer.RenderTerrain(tiles)
		if err != nil {
			b.Fatalf("Render failed: %v", err)
		}
	}
}

func BenchmarkHexToPixelConversion(b *testing.B) {
	gridConfig := hex.GridConfig{Width: 10, Height: 10, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	renderConfig := RenderConfig{
		Width:   400,
		Height:  400,
		HexSize: 10.0,
	}
	renderer := NewHexRenderer(grid, renderConfig)

	coords := []hex.AxialCoord{
		hex.NewAxialCoord(0, 0),
		hex.NewAxialCoord(5, 3),
		hex.NewAxialCoord(-2, 7),
		hex.NewAxialCoord(8, -1),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		coord := coords[i%len(coords)]
		renderer.hexToPixel(coord)
	}
}

func BenchmarkMultiLayerRender(b *testing.B) {
	// Setup
	gridConfig := hex.GridConfig{Width: 20, Height: 20, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	tiles := make([]*terrain.HexTile, 0, 400)
	for _, coord := range grid.AllCoords() {
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   float64(coord.Q + coord.R*5),
			IsLand:      coord.Q%3 != 0,
		}
		tiles = append(tiles, tile)
	}

	renderConfig := RenderConfig{
		Width:       300,
		Height:      300,
		HexSize:     6.0,
		Layers:      []RenderLayer{LayerElevation, LayerWater, LayerDebugCoords},
		ColorScheme: SchemeElevation,
		Quality:     85,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer := NewHexRenderer(grid, renderConfig)
		_, err := renderer.RenderTerrain(tiles)
		if err != nil {
			b.Fatalf("Multi-layer render failed: %v", err)
		}
	}
}

func BenchmarkJPEGExport(b *testing.B) {
	// Setup
	gridConfig := hex.GridConfig{Width: 15, Height: 15, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	tiles := make([]*terrain.HexTile, 0, 225)
	for _, coord := range grid.AllCoords() {
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   float64(coord.Q + coord.R*10),
			IsLand:      true,
		}
		tiles = append(tiles, tile)
	}

	renderConfig := RenderConfig{
		Width:       250,
		Height:      250,
		HexSize:     7.0,
		Layers:      []RenderLayer{LayerElevation},
		ColorScheme: SchemeElevation,
		Quality:     85,
	}

	// Pre-render the terrain
	renderer := NewHexRenderer(grid, renderConfig)
	_, err := renderer.RenderTerrain(tiles)
	if err != nil {
		b.Fatalf("Failed to render terrain for export benchmark: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := "benchmark_temp.jpg"
		err := renderer.ExportJPEG(filename, 85)
		if err != nil {
			b.Fatalf("JPEG export failed: %v", err)
		}
		// Clean up (in real benchmark, we'd do this in b.StopTimer())
		os.Remove(filename)
	}
}

func BenchmarkColorSchemeComparison(b *testing.B) {
	schemes := []ColorScheme{
		SchemeElevation,
		SchemeRealistic,
		SchemeDebug,
		SchemeGrayscale,
	}

	gridConfig := hex.GridConfig{Width: 10, Height: 10, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(gridConfig)

	tiles := make([]*terrain.HexTile, 0, 100)
	for _, coord := range grid.AllCoords() {
		tile := &terrain.HexTile{
			Coordinates: coord,
			Elevation:   float64(coord.Q*100 + coord.R*50),
			IsLand:      true,
		}
		tiles = append(tiles, tile)
	}

	schemeNames := []string{"elevation", "realistic", "debug", "grayscale"}

	for i, scheme := range schemes {
		b.Run(schemeNames[i], func(b *testing.B) {
			renderConfig := RenderConfig{
				Width:       200,
				Height:      200,
				HexSize:     8.0,
				Layers:      []RenderLayer{LayerElevation},
				ColorScheme: scheme,
				Quality:     85,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderer := NewHexRenderer(grid, renderConfig)
				_, err := renderer.RenderTerrain(tiles)
				if err != nil {
					b.Fatalf("Render failed for scheme %v: %v", scheme, err)
				}
			}
		})
	}
}
