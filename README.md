# Hex Map World Generation

A procedural world generation system using hexagonal tiles with realistic geological and hydrological simulation.

## Overview

This project implements a comprehensive world generation pipeline that creates realistic terrains, climate systems, hydrology, and biomes using scientifically-based models. The system supports both world-scale generation (with toroidal wrapping) and regional maps (bounded areas).

## Current Status

### ✅ Completed Deliverables

#### HEX-001: Hex Grid Foundation
- **Axial coordinate system** with efficient neighbor lookups
- **Dual topology support**: World maps (toroidal) and Region maps (bounded)
- **Coordinate conversions**: Axial ↔ Offset ↔ Pixel
- **Distance calculations** with shortest path for world maps
- **Comprehensive test coverage**: 18 tests, 100% passing

## Quick Start

### Build and Test
```bash
# Clone and build
git clone <repository>
cd hex-map
go build ./cmd/hex-world

# Run tests
go test ./...

# Demo the coordinate system
./hex-world demo-coords --size=20x20 --topology=world
./hex-world demo-distance --from=0,0 --to=9,0 --topology=world
```

### Features

#### Coordinate System
- **Flat-top hexagons** with 10km center-to-center spacing
- **Axial coordinates** for efficient computation
- **Topology awareness**: World vs Region behavior

#### Dual Topology Support
- **World Maps**: Toroidal wrapping where edges connect
  - All hexes have exactly 6 neighbors
  - Shortest path calculations consider wrapping
  - Perfect for planet-scale generation
- **Region Maps**: Bounded areas with edge effects
  - Edge hexes have 2-5 neighbors  
  - Natural boundaries for continental/island generation

## Architecture

### Project Structure
```
pkg/
  hex/                 # Core hex grid system
    coordinate.go      # Axial coordinate implementation
    topology.go        # Grid with topology support
    *_test.go         # Comprehensive tests
cmd/
  hex-world/           # CLI application
    main.go           # Demo commands
docs/
  features/           # Feature documentation
    HEX-001.md       # Hex grid foundation spec
```

### Design Philosophy

The system follows a **12-step deliverable workflow** ensuring:
- ✅ Comprehensive planning and documentation
- ✅ Test-driven development with >90% coverage
- ✅ Scientific accuracy with real-world parameters
- ✅ Modular architecture with clear dependencies

## Roadmap

### 🔄 Next: TERRAIN-001
- Diamond-Square algorithm for fractal terrain
- Multi-octave noise generation
- Hypsometric curve matching Earth's distribution
- Land/water designation

### 🔮 Future Deliverables
- **VIZ-001**: Visualization and rendering system
- **HYDRO-001**: Water flow and drainage networks
- **CLIMATE-001**: Temperature and precipitation modeling
- **SOIL-001**: Soil profile generation
- **BIOME-001**: Biome classification system
- And 7 more deliverables...

## Documentation

- **[DESIGN.md](DESIGN.md)**: Complete technical specification
- **[HEX-001 Feature Plan](docs/features/HEX-001.md)**: Hex grid implementation details

## Contributing

This project follows the 12-step deliverable workflow outlined in DESIGN.md. Each feature is:
1. Planned with comprehensive documentation
2. Test-driven with failing tests first
3. Implemented with scientific accuracy
4. Validated with extensive testing

## License

[License details to be added]
