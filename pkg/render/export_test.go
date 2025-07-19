package render

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sean/hex-map/pkg/hex"
	"github.com/sean/hex-map/pkg/terrain"
)

func TestEmbedMetadata(t *testing.T) {
	// Create test metadata
	metadata := RenderMetadata{
		Generator:    "hex-world v1.0",
		Timestamp:    time.Now().Format(time.RFC3339),
		WorldSeed:    42,
		Stage:        "post_terrain_generation",
		QualityScore: 0.85,
		KnownIssues:  []string{"elevation spike at (5,3)"},
	}

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with test pattern
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}

	// This should fail until we implement metadata embedding
	err := EmbedMetadata(img, metadata)
	if err != nil {
		t.Errorf("EmbedMetadata failed: %v", err)
	}

	// Test metadata extraction (should fail since not implemented yet)
	_, err = ExtractMetadata(img)
	if err == nil {
		t.Error("ExtractMetadata should fail since not implemented yet")
	}
}

func TestExportJPEGWithMetadata(t *testing.T) {
	// Create test directory
	testDir := "test_exports"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create test renderer
	config := hex.GridConfig{Width: 5, Height: 5, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{
		Width:   200,
		Height:  200,
		HexSize: 15.0,
		Quality: 90,
	}
	renderer := NewHexRenderer(grid, renderConfig)

	// Create test terrain
	tiles := []*terrain.HexTile{
		{Coordinates: hex.NewAxialCoord(0, 0), Elevation: 100.0, IsLand: true},
		{Coordinates: hex.NewAxialCoord(1, 0), Elevation: -50.0, IsLand: false},
	}

	// Render terrain
	_, err := renderer.RenderTerrain(tiles)
	if err != nil {
		t.Fatalf("Failed to render terrain: %v", err)
	}

	// Create metadata
	metadata := RenderMetadata{
		Generator: "test-renderer",
		WorldSeed: 123,
		Stage:     "test",
	}

	// Export with metadata
	filename := filepath.Join(testDir, "test_with_metadata.jpg")
	err = renderer.ExportJPEGWithMetadata(filename, metadata)
	if err != nil {
		t.Errorf("ExportJPEGWithMetadata failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("JPEG file was not created")
	}

	// Test reading metadata back (should fail since not implemented yet)
	_, err = ExtractMetadataFromFile(filename)
	if err == nil {
		t.Error("ExtractMetadataFromFile should fail since not implemented yet")
	}
}

func TestExportPNGWithMetadata(t *testing.T) {
	// Create test directory
	testDir := "test_exports"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create test renderer
	config := hex.GridConfig{Width: 3, Height: 3, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{Width: 150, Height: 150, HexSize: 20.0}
	renderer := NewHexRenderer(grid, renderConfig)

	// Create test terrain
	tiles := []*terrain.HexTile{
		{Coordinates: hex.NewAxialCoord(0, 0), Elevation: 200.0, IsLand: true},
	}

	// Render terrain
	_, err := renderer.RenderTerrain(tiles)
	if err != nil {
		t.Fatalf("Failed to render terrain: %v", err)
	}

	// Create metadata
	metadata := RenderMetadata{
		Generator: "test-renderer-png",
		WorldSeed: 456,
		Stage:     "test-png",
	}

	// Export with metadata
	filename := filepath.Join(testDir, "test_with_metadata.png")
	err = renderer.ExportPNGWithMetadata(filename, metadata)
	if err != nil {
		t.Errorf("ExportPNGWithMetadata failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("PNG file was not created")
	}

	// Test reading metadata back (should fail since not implemented yet)
	_, err = ExtractMetadataFromFile(filename)
	if err == nil {
		t.Error("ExtractMetadataFromFile should fail since not implemented yet")
	}
}

func TestInvalidExportPaths(t *testing.T) {
	// Create test renderer
	config := hex.GridConfig{Width: 2, Height: 2, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{Width: 100, Height: 100, HexSize: 10.0}
	renderer := NewHexRenderer(grid, renderConfig)

	// Try to export to invalid path
	err := renderer.ExportJPEG("/invalid/path/test.jpg", 85)
	if err == nil {
		t.Error("ExportJPEG should fail with invalid path")
	}

	err = renderer.ExportPNG("/invalid/path/test.png")
	if err == nil {
		t.Error("ExportPNG should fail with invalid path")
	}
}

func TestJPEGQualityValidation(t *testing.T) {
	config := hex.GridConfig{Width: 2, Height: 2, Topology: hex.TopologyRegion}
	grid := hex.NewGrid(config)

	renderConfig := RenderConfig{Width: 100, Height: 100, HexSize: 10.0}
	renderer := NewHexRenderer(grid, renderConfig)

	// Create test directory
	testDir := "test_quality"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Test invalid quality values
	invalidQualities := []int{-1, 0, 101, 200}
	for _, quality := range invalidQualities {
		filename := filepath.Join(testDir, "test.jpg")
		err := renderer.ExportJPEG(filename, quality)
		if err == nil {
			t.Errorf("ExportJPEG should fail with invalid quality %d", quality)
		}
	}

	// Test valid quality values
	validQualities := []int{1, 50, 85, 100}
	for _, quality := range validQualities {
		filename := filepath.Join(testDir, "test_valid.jpg")
		// Create a simple image first
		tiles := []*terrain.HexTile{
			{Coordinates: hex.NewAxialCoord(0, 0), Elevation: 100.0, IsLand: true},
		}
		renderer.RenderTerrain(tiles)

		err := renderer.ExportJPEG(filename, quality)
		if err != nil {
			t.Errorf("ExportJPEG should succeed with valid quality %d: %v", quality, err)
		}

		// Clean up
		os.Remove(filename)
	}
}

func TestMetadataJSONMarshalling(t *testing.T) {
	// Test that metadata can be marshalled/unmarshalled properly
	original := RenderMetadata{
		Generator:    "test-generator",
		Timestamp:    "2025-01-19T15:30:00Z",
		WorldSeed:    789,
		Stage:        "test-stage",
		QualityScore: 0.92,
		KnownIssues:  []string{"issue1", "issue2"},
	}

	// This should fail until we implement proper JSON handling
	jsonData, err := original.ToJSON()
	if err != nil {
		t.Errorf("Failed to marshal metadata to JSON: %v", err)
	}

	var reconstructed RenderMetadata
	err = reconstructed.FromJSON(jsonData)
	if err != nil {
		t.Errorf("Failed to unmarshal metadata from JSON: %v", err)
	}

	if reconstructed.Generator != original.Generator {
		t.Errorf("Generator not preserved: got %s, expected %s",
			reconstructed.Generator, original.Generator)
	}

	if reconstructed.WorldSeed != original.WorldSeed {
		t.Errorf("WorldSeed not preserved: got %d, expected %d",
			reconstructed.WorldSeed, original.WorldSeed)
	}

	if len(reconstructed.KnownIssues) != len(original.KnownIssues) {
		t.Errorf("KnownIssues length not preserved: got %d, expected %d",
			len(reconstructed.KnownIssues), len(original.KnownIssues))
	}
}

func TestExtractMetadataFromNonExistentFile(t *testing.T) {
	// Test error handling for non-existent files
	_, err := ExtractMetadataFromFile("non_existent_file.jpg")
	if err == nil {
		t.Error("ExtractMetadataFromFile should fail for non-existent files")
	}
}

func TestExtractMetadataFromFileWithoutMetadata(t *testing.T) {
	// Create a simple image file without metadata
	testDir := "test_no_metadata"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create a basic image without our custom metadata
	filename := filepath.Join(testDir, "no_metadata.png")

	// Save without metadata
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()

	// This would use standard PNG encoding without metadata
	// The extraction should handle this gracefully
	_, err = ExtractMetadataFromFile(filename)
	if err == nil {
		t.Error("ExtractMetadataFromFile should fail for files without metadata")
	}
}
