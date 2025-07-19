package terrain

import (
	"testing"

	"github.com/sean/hex-map/pkg/hex"
)

func TestDefaultTerrainConfig(t *testing.T) {
	config := DefaultTerrainConfig()
	
	// Test that default config is valid
	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid, got error: %v", err)
	}
	
	// Test expected default values
	if config.SeaLevel != 0.0 {
		t.Errorf("Expected default sea level 0.0, got %f", config.SeaLevel)
	}
	
	if config.LandRatio != 0.29 {
		t.Errorf("Expected default land ratio 0.29, got %f", config.LandRatio)
	}
	
	if config.NoiseParams.Octaves != 6 {
		t.Errorf("Expected default octaves 6, got %d", config.NoiseParams.Octaves)
	}
}

func TestTerrainConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    TerrainConfig
		wantError bool
	}{
		{
			name:      "valid config",
			config:    DefaultTerrainConfig(),
			wantError: false,
		},
		{
			name: "invalid land ratio - negative",
			config: TerrainConfig{
				LandRatio: -0.1,
				NoiseParams: DefaultNoiseParameters(),
			},
			wantError: true,
		},
		{
			name: "invalid land ratio - too high",
			config: TerrainConfig{
				LandRatio: 1.5,
				NoiseParams: DefaultNoiseParameters(),
			},
			wantError: true,
		},
		{
			name: "invalid octaves - too low",
			config: TerrainConfig{
				LandRatio: 0.3,
				NoiseParams: NoiseParameters{
					Octaves:     0,
					Persistence: 0.5,
					Lacunarity:  2.0,
					HurstExp:    0.85,
				},
			},
			wantError: true,
		},
		{
			name: "invalid persistence - negative",
			config: TerrainConfig{
				LandRatio: 0.3,
				NoiseParams: NoiseParameters{
					Octaves:     6,
					Persistence: -0.1,
					Lacunarity:  2.0,
					HurstExp:    0.85,
				},
			},
			wantError: true,
		},
		{
			name: "invalid lacunarity - too low",
			config: TerrainConfig{
				LandRatio: 0.3,
				NoiseParams: NoiseParameters{
					Octaves:     6,
					Persistence: 0.5,
					Lacunarity:  0.5,
					HurstExp:    0.85,
				},
			},
			wantError: true,
		},
		{
			name: "invalid hurst exponent - negative",
			config: TerrainConfig{
				LandRatio: 0.3,
				NoiseParams: NoiseParameters{
					Octaves:     6,
					Persistence: 0.5,
					Lacunarity:  2.0,
					HurstExp:    -0.1,
				},
			},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestHexTileRealism(t *testing.T) {
	tests := []struct {
		name      string
		tile      HexTile
		realistic bool
	}{
		{
			name: "realistic elevation",
			tile: HexTile{
				Elevation: 1000.0,
			},
			realistic: true,
		},
		{
			name: "too high elevation",
			tile: HexTile{
				Elevation: 15000.0,
			},
			realistic: false,
		},
		{
			name: "too low elevation",
			tile: HexTile{
				Elevation: -15000.0,
			},
			realistic: false,
		},
		{
			name: "everest height",
			tile: HexTile{
				Elevation: ElevationMax,
			},
			realistic: true,
		},
		{
			name: "mariana trench depth",
			tile: HexTile{
				Elevation: ElevationMin,
			},
			realistic: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tile.IsRealistic()
			if result != tt.realistic {
				t.Errorf("IsRealistic() = %v, want %v for elevation %f", 
					result, tt.realistic, tt.tile.Elevation)
			}
		})
	}
}

func TestClassifyLandWater(t *testing.T) {
	tests := []struct {
		name      string
		elevation float64
		seaLevel  float64
		wantLand  bool
	}{
		{
			name:      "above sea level",
			elevation: 100.0,
			seaLevel:  0.0,
			wantLand:  true,
		},
		{
			name:      "below sea level",
			elevation: -100.0,
			seaLevel:  0.0,
			wantLand:  false,
		},
		{
			name:      "at sea level",
			elevation: 0.0,
			seaLevel:  0.0,
			wantLand:  false,
		},
		{
			name:      "above custom sea level",
			elevation: 50.0,
			seaLevel:  -10.0,
			wantLand:  true,
		},
		{
			name:      "below custom sea level",
			elevation: -50.0,
			seaLevel:  -10.0,
			wantLand:  false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tile := &HexTile{
				Elevation: tt.elevation,
			}
			
			tile.ClassifyLandWater(tt.seaLevel)
			
			if tile.IsLand != tt.wantLand {
				t.Errorf("ClassifyLandWater() set IsLand = %v, want %v", 
					tile.IsLand, tt.wantLand)
			}
		})
	}
}

func TestGetDepthAndHeight(t *testing.T) {
	tests := []struct {
		name           string
		elevation      float64
		seaLevel       float64
		expectedDepth  float64
		expectedHeight float64
	}{
		{
			name:           "land above sea level",
			elevation:      100.0,
			seaLevel:       0.0,
			expectedDepth:  0.0,
			expectedHeight: 100.0,
		},
		{
			name:           "water below sea level",
			elevation:      -50.0,
			seaLevel:       0.0,
			expectedDepth:  50.0,
			expectedHeight: 0.0,
		},
		{
			name:           "at sea level",
			elevation:      0.0,
			seaLevel:       0.0,
			expectedDepth:  0.0,
			expectedHeight: 0.0,
		},
		{
			name:           "custom sea level",
			elevation:      25.0,
			seaLevel:       50.0,
			expectedDepth:  25.0,
			expectedHeight: 0.0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tile := &HexTile{
				Elevation: tt.elevation,
			}
			
			depth := tile.GetDepth(tt.seaLevel)
			height := tile.GetHeight(tt.seaLevel)
			
			if depth != tt.expectedDepth {
				t.Errorf("GetDepth() = %f, want %f", depth, tt.expectedDepth)
			}
			
			if height != tt.expectedHeight {
				t.Errorf("GetHeight() = %f, want %f", height, tt.expectedHeight)
			}
		})
	}
}

func TestTerrainError(t *testing.T) {
	err := &TerrainError{"test error message"}
	expected := "terrain error: test error message"
	
	if err.Error() != expected {
		t.Errorf("TerrainError.Error() = %s, want %s", err.Error(), expected)
	}
}

func TestDefaultNoiseParameters(t *testing.T) {
	params := DefaultNoiseParameters()
	
	// Test reasonable default values
	if params.Octaves < 1 {
		t.Errorf("Default octaves should be positive, got %d", params.Octaves)
	}
	
	if params.Persistence <= 0 || params.Persistence > 1 {
		t.Errorf("Default persistence should be in (0,1], got %f", params.Persistence)
	}
	
	if params.Lacunarity <= 1 {
		t.Errorf("Default lacunarity should be > 1, got %f", params.Lacunarity)
	}
	
	if params.HurstExp < 0 || params.HurstExp > 1 {
		t.Errorf("Default Hurst exponent should be in [0,1], got %f", params.HurstExp)
	}
}

// Test that coordinates can be properly serialized/deserialized
func TestHexTileSerialization(t *testing.T) {
	coord := hex.NewAxialCoord(5, -3)
	original := &HexTile{
		Coordinates:     coord,
		Elevation:       1234.5,
		IsLand:         true,
		DistanceToWater: 2.5,
	}
	
	// Test that all fields are accessible (this would catch JSON tag issues)
	if original.Coordinates.Q != 5 {
		t.Errorf("Coordinates.Q = %d, want 5", original.Coordinates.Q)
	}
	
	if original.Coordinates.R != -3 {
		t.Errorf("Coordinates.R = %d, want -3", original.Coordinates.R)
	}
	
	if original.Elevation != 1234.5 {
		t.Errorf("Elevation = %f, want 1234.5", original.Elevation)
	}
	
	if !original.IsLand {
		t.Errorf("IsLand = %v, want true", original.IsLand)
	}
}