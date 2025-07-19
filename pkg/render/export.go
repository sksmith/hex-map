package render

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

// EmbedMetadata embeds metadata in image (simplified implementation)
func EmbedMetadata(img *image.RGBA, metadata RenderMetadata) error {
	// For now, this is a placeholder
	// Real implementation would embed metadata in image headers/EXIF data
	return nil
}

// ExtractMetadata extracts metadata from image (simplified implementation)
func ExtractMetadata(img *image.RGBA) (*RenderMetadata, error) {
	// For now, return an error as we don't have embedded metadata yet
	return nil, fmt.Errorf("metadata extraction not yet implemented")
}

// ExtractMetadataFromFile extracts metadata from file
func ExtractMetadataFromFile(filename string) (*RenderMetadata, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Try to decode as JPEG
	_, err = jpeg.Decode(file)
	if err == nil {
		// It's a JPEG, but we don't have metadata extraction yet
		return nil, fmt.Errorf("metadata extraction from JPEG not yet implemented")
	}

	// Reset file position
	file.Seek(0, 0)

	// Try to decode as PNG
	_, err = png.Decode(file)
	if err == nil {
		// It's a PNG, but we don't have metadata extraction yet
		return nil, fmt.Errorf("metadata extraction from PNG not yet implemented")
	}

	return nil, fmt.Errorf("unsupported file format or corrupted file")
}

// ToJSON converts metadata to JSON
func (rm *RenderMetadata) ToJSON() ([]byte, error) {
	return json.Marshal(rm)
}

// FromJSON converts metadata from JSON
func (rm *RenderMetadata) FromJSON(data []byte) error {
	return json.Unmarshal(data, rm)
}
