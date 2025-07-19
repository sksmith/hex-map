# Hex Map World Generation Design

## Overview
This document captures design decisions for a procedural world generation system using hexagonal tiles with realistic geological and hydrological simulation.

### Specifications
- **Hex Size**: 10km center-to-center (flat-top orientation)
- **Hex Radius**: ~5.77km (center to vertex)
- **Hex Apothem**: 5km (center to edge)
- **Area per Hex**: ~86.6 km²

## Core Components

### 1. Hex Grid System
- **Grid Structure**: Axial coordinate system for efficient neighbor lookups
- **Topology Support**: 
  - **World Maps**: Toroidal topology where top/bottom edges connect and left/right edges wrap
  - **Region Maps**: Bounded topology where edge hexes have fewer neighbors
- **Tile Properties**: 
  - Elevation (meters)
  - Temperature (°C)
  - Moisture/precipitation
  - Soil composition
  - Vegetation type
  - Water flow direction and volume

### 2. Terrain Generation (Real-World Models)
- **Base Heightmap**: 
  - Diamond-Square algorithm for fractal terrain
  - Spectral synthesis using real-world terrain power spectra (β ≈ 2)
  - Multi-fractal noise with Hurst exponent H = 0.8-0.9
- **Tectonic Simulation**: 
  - Simplified plate tectonics with convergent/divergent/transform boundaries
  - Isostatic adjustment for crustal thickness
  - Fault generation using Voronoi diagrams
- **Initial Landmass**: 
  - Hypsometric curve matching Earth's elevation distribution
  - Continental shelf modeling (slope ≈ 0.1°)
  - Fractal coastlines with dimension D ≈ 1.25

### 3. Erosion Simulation (Stream Power Model)
- **Hydraulic Erosion**: 
  - Stream power law: E = K × A^m × S^n (m≈0.5, n≈1)
  - Hack's Law for drainage networks: L = C × A^h (h≈0.6)
  - Manning's equation for flow velocity
- **Thermal Erosion**: 
  - Talus angle: 30-35° for loose material
  - Freeze-thaw cycles based on temperature crossing 0°C
- **Chemical Erosion**: 
  - Carbonate dissolution rates: 10-100 mm/kyr
  - Silicate weathering: 1-10 mm/kyr
- **Sediment Transport**: 
  - Hjulström-Sundborg diagram for particle transport
  - Settling velocity by Stokes' law

### 4. Hydrological System (Enhanced)

#### 4.1 Ocean Current Modeling
- **Surface Currents**: 
  - Driven by prevailing winds with 45° Coriolis deflection
  - Wind stress: τ = ρ_air × C_drag × |wind|² × wind_direction
  - Ekman transport: 90° deflection from wind direction
- **Thermohaline Circulation**:
  - Deep water formation in polar regions (temperature + salinity)
  - Upwelling along continental margins
  - Heat transport from equator to poles
- **Continental Boundaries**:
  - Western boundary currents (Gulf Stream, Kuroshio): narrow, fast
  - Eastern boundary currents (California, Canary): broad, slow
  - Coastal upwelling and downwelling
- **Gyre Formation**:
  - Subtropical gyres: clockwise NH, counterclockwise SH
  - Subpolar gyres: opposite rotation
  - Vector smoothing for realistic current patterns

#### 4.2 Terrestrial Hydrology
- **Rainfall Distribution**: 
  - Orographic enhancement: P = P₀(1 + kH), k≈0.0002/m
  - Köppen climate zones for precipitation patterns
  - Annual precipitation: 200-2000mm based on location
  - Ocean proximity modifier: coastal areas +20% moisture
- **River Formation**: 
  - D8 flow accumulation algorithm
  - Drainage density: 0.5-2.5 km/km²
  - Strahler stream ordering
  - Stream power calculation for channel formation
- **Lakes and Wetlands**: 
  - Depression filling algorithm
  - Perched water tables in clay layers
  - Wetland formation at <2% slope
  - Lake effect on local climate (+humidity, -temperature range)
- **Groundwater**: 
  - Darcy's law for aquifer flow: Q = -KA(dh/dl)
  - Hydraulic conductivity by soil type (10⁻⁸ to 10⁻² m/s)
  - Water table depth affects vegetation access
  - Springs at groundwater-surface intersections

### 5. Climate and Weather (Enhanced Energy Balance Model)

#### 5.1 Long-term Climate
- **Temperature Zones**: 
  - Lapse rate: -6.5°C/km elevation
  - Latitude temperature: T = T_eq × cos(lat)^0.25
  - Ocean moderation: ΔT = 10-20°C seasonal range
  - Aspect modifier: south-facing slopes +2°C (NH)
- **Prevailing Winds**: 
  - Hadley, Ferrel, and Polar cells
  - Trade winds: 15-30° latitude
  - Westerlies: 30-60° latitude
  - Terrain deflection around mountains
- **Rain Shadows**: 
  - Windward: 100% precipitation
  - Lee side: 20-50% precipitation
  - Föhn warming: +2°C/100m descent
- **Seasonal Variation**: 
  - Axial tilt effects: ±23.5° solar angle
  - Monsoon patterns in tropical regions

#### 5.2 Dynamic Weather Systems
- **Air Pressure Systems**:
  - High pressure: clear skies, stable conditions, subsiding air
  - Low pressure: clouds, precipitation, rising air
  - Pressure gradient force: F = -1/ρ × ∇P
  - Geostrophic wind: wind parallel to isobars
- **Weather Fronts**:
  - Cold fronts: steep, fast-moving, thunderstorms
  - Warm fronts: gradual, widespread precipitation
  - Occluded fronts: complex, mixed precipitation
  - Front propagation: 30-50 km/h typical speed
- **Daily Weather Updates**:
  - Temperature: base_climate ± daily_variation ± weather_system_effect
  - Humidity: evaporation + advection + precipitation
  - Cloud formation: relative humidity >80% threshold
  - Precipitation: cloud water content > saturation threshold
- **Weather Events**:
  - Thunderstorms: high CAPE (convective available potential energy)
  - Blizzards: low temperature + high wind + precipitation
  - Droughts: persistent high pressure blocking precipitation
  - Hurricanes: warm ocean + low wind shear + Coriolis effect

### 6. Soil Composition Modeling

#### 6.1 Soil Formation Factors
- **Parent Material**: Base rock type influences initial soil chemistry
  - Granite → acidic, well-drained soils
  - Limestone → alkaline, clay-rich soils
  - Volcanic → fertile, high mineral content
  - Sedimentary → varied based on composition
- **Climate**: Temperature and precipitation drive weathering rates
  - High temp + high precip = deep, leached soils
  - Low temp + low precip = shallow, mineral-rich soils
- **Topography**: Slope affects erosion and water retention
  - Steep slopes: shallow, rocky, fast drainage
  - Gentle slopes: deeper soils, better retention
  - Depressions: organic accumulation, poor drainage
- **Time**: Soil development over geological timescales
- **Biota**: Vegetation and decomposer activity

#### 6.2 Soil Profile Generation
```
soil_texture = determine_texture(rock_type, climate, slope)
soil_depth = calculate_depth(slope, erosion_rate, age)
fertility = base_fertility(rock_type) * climate_modifier * organic_modifier
drainage = drainage_class(texture, slope, clay_content)
erosion_risk = erosion_potential(slope, precipitation, vegetation_cover)
```

#### 6.3 Soil Classification Rules
- **Desert/Arid**: Sandy, fast drainage, low fertility, high erosion risk
- **Rainforest**: Clay-rich, poor drainage, low fertility (leached), low erosion risk
- **Temperate Forest**: Loamy, moderate drainage, high fertility, low erosion risk
- **Grassland**: Deep, fertile, moderate drainage, moderate erosion risk
- **Mountain/Alpine**: Rocky, shallow, fast drainage, low fertility, high erosion risk
- **Wetland/Floodplain**: Organic/alluvial, deep, poor drainage, high fertility

### 7. Biome Generation (Enhanced Whittaker Model)
- **Primary Classification**: 
  - Whittaker diagram: Temperature vs Precipitation
  - Holdridge life zones for elevation refinement
  - Treeline: ~10°C July isotherm
- **Secondary Modifiers**:
  - Soil quality affects vegetation density
  - Slope limits forest establishment (>30° = alpine/grassland)
  - Drainage affects wetland vs terrestrial biomes
- **Vegetation Growth**: 
  - NPP = 3000 × (1 - e^(-0.000664×P)) × soil_fertility_modifier
  - LAI (Leaf Area Index): 0.5-8.0 by biome type
  - Growing degree days (GDD) for phenology
- **Ecosystem Interactions**: 
  - Liebig's law of the minimum (nutrients, water, temperature)
  - Carrying capacity by NPP and soil fertility
  - Succession stages: pioneer → intermediate → climax

## Technical Architecture

### Data Structures

```
HexTile {
  // Static Geographic Attributes
  coordinates: AxialCoord
  latitude: float64
  longitude: float64
  elevation: float64 (meters)
  slope: float64 (computed from neighbors)
  aspect: float64 (degrees, direction slope faces)
  base_rock_type: RockType (Granite, Limestone, Volcanic, Sedimentary)
  distance_to_water: float64 (km)
  tectonic_region: TectonicType (PlateBoundary, Shield, Hotspot, Rift)
  is_edge_hex: bool (true if on map boundary for region maps)
  
  // Derived Static Attributes
  prevailing_wind: Vector2D (from latitude + terrain deflection)
  ocean_current: Vector2D (if ocean tile)
  rain_shadow_effect: float64 (0.0-1.0, from mountain positioning)
  solar_exposure: float64 (from aspect + latitude)
  
  // Climate Attributes (Long-term averages)
  avg_annual_temp: float64 (Celsius)
  avg_annual_precip: float64 (mm)
  seasonal_temp_range: float64 (Celsius)
  humidity_index: float64 (0.0-1.0)
  aridity_index: float64 (precip/potential_evapotranspiration)
  
  // Dynamic Weather Attributes
  current_temp: float64 (Celsius)
  humidity: float64 (0.0-1.0)
  precipitation: float64 (mm/day)
  cloud_cover: float64 (0.0-1.0)
  wind_speed: float64 (km/h)
  wind_direction: float64 (degrees)
  air_pressure: float64 (hPa)
  weather_event: Optional<WeatherEvent>
  
  // Hydrological Attributes
  water_volume: float64
  flow_direction: Direction
  flow_accumulation: float64
  stream_order: int (Strahler ordering)
  
  // Ecological Attributes
  biome: BiomeType
  vegetation_density: float64
  soil_profile: SoilProfile
}

SoilProfile {
  texture: SoilTexture (Sandy, Loamy, Clay, Peaty, Rocky)
  depth: SoilDepth (Shallow, Moderate, Deep)
  fertility: float64 (0.0-1.0)
  drainage: DrainageType (Fast, Moderate, Poor)
  erosion_risk: ErosionRisk (Low, Moderate, High)
  organic_content: float64 (0.0-1.0)
  ph_level: float64 (3.0-9.0)
  nutrient_levels: {
    nitrogen: float64
    phosphorus: float64
    potassium: float64
  }
}

WeatherEvent {
  type: EventType (Thunderstorm, Blizzard, Drought, Hurricane, etc.)
  intensity: float64 (0.0-1.0)
  duration: int (days)
  affected_radius: int (hex count)
}
```

### Processing Pipeline

#### Phase 1: World Generation
1. Generate base terrain heightmap using fractal noise
2. Apply tectonic forces and fault generation
3. Calculate static geographic attributes (slope, aspect, tectonic regions)
4. Compute derived static attributes (prevailing winds, rain shadows)

#### Phase 2: Oceanic Systems
5. Generate ocean currents from wind patterns + Coriolis effect
6. Model oceanic heat transport and circulation patterns
7. Calculate distance to water for all land hexes

#### Phase 3: Climate Calculation
8. Establish long-term climate averages (temperature, precipitation)
9. Apply orographic effects and rain shadow calculations
10. Simulate water flow and drainage networks
11. Apply erosion cycles with sediment transport

#### Phase 4: Soil and Biome Generation
12. Generate soil profiles based on climate, geology, and slope
13. Classify biomes using enhanced Whittaker model
14. Establish vegetation density and ecosystem interactions

#### Phase 5: Dynamic Weather Engine
15. Initialize weather systems with pressure gradients
16. Implement daily weather simulation with fronts
17. Add seasonal variation and weather events

### Snapshot/Checkpoint System
The pipeline supports capturing intermediate states for visualization and debugging:

```
PipelineSnapshot {
  stage_name: string
  timestamp: datetime
  hex_data: HexTile[]
  statistics: {
    elevation_range: [min, max]
    water_coverage: percentage
    biome_distribution: map<BiomeType, count>
    erosion_volume: float
  }
  visualization_hints: {
    primary_layer: enum (elevation|temperature|moisture|flow)
    color_scheme: string
    highlight_changes: bool
  }
}
```

**Snapshot Points**:
- After initial terrain generation
- Post-tectonic simulation and slope calculation
- After ocean current establishment
- Post-climate calculation (long-term averages)
- After each erosion iteration (configurable)
- Post-soil profile generation
- Post-hydrology establishment
- Final biome assignment
- Daily weather snapshots (during dynamic simulation)

**Visualization Options**:
- **Terrain**: Elevation heatmap with hillshading, slope analysis
- **Hydrology**: Water flow accumulation, ocean currents, stream networks
- **Climate**: Temperature/moisture gradients, pressure systems, wind patterns
- **Erosion**: Change detection (before/after), sediment transport
- **Soil**: Fertility maps, drainage classification, erosion risk
- **Weather**: Current conditions, weather fronts, storm tracking
- **Biomes**: Distribution map, vegetation density, ecosystem health
- **Combined**: Multi-layer overlays, 3D terrain mesh export
- **Time Series**: Animation of weather systems, seasonal changes

### Dual Visualization System

#### Human Visual Output (JPEG/PNG)
```
VisualRenderer {
  // Layer rendering with proper color mapping
  render_elevation_map(colormap="terrain", hillshade=true)
  render_flow_map(streamlines=true, width_by_volume=true)
  render_climate_map(temp_precip_overlay=true)
  render_biome_map(realistic_colors=true)
  render_soil_map(fertility_texture_overlay=true)
  
  // Composite views
  render_overview(layers=[elevation, water, biomes])
  render_debug_overlay(show_hex_coords=true, show_values=true)
  
  // Export formats
  save_jpeg(filename, quality=95)
  save_png(filename, transparent_background=false)
}
```

#### Machine-Readable Debug Output (JSON/CSV)
```
DebugAnalyzer {
  // Statistical summaries for validation
  elevation_stats: {min, max, mean, std_dev, percentiles}
  flow_network_stats: {total_streams, max_order, drainage_density}
  climate_validation: {temp_range_by_latitude, precip_correlation}
  biome_distribution: {area_percentages, transition_zones}
  soil_properties: {fertility_by_biome, drainage_issues}
  
  // Anomaly detection
  identify_elevation_outliers(threshold=3_std_dev)
  detect_flow_discontinuities()
  find_climate_inconsistencies()
  validate_biome_transitions()
  check_soil_logic_errors()
  
  // Neighbor analysis for debugging
  analyze_hex_neighborhood(hex_id) -> {
    center: HexTile,
    neighbors: [HexTile],
    gradients: {elevation, temperature, moisture},
    flow_connectivity: bool,
    transition_smoothness: float
  }
  
  // Export formats
  export_summary_stats(format="json")
  export_hex_grid_csv(selected_attributes=[])
  export_validation_report(issues_only=true)
}

### Debug Output Formats for AI Analysis

#### 1. Validation Report (JSON)
```json
{
  "generation_stage": "post_erosion",
  "timestamp": "2025-01-19T14:30:00Z",
  "world_stats": {
    "hex_count": 10000,
    "land_percentage": 68.5,
    "elevation": {"min": -2840, "max": 4567, "mean": 245, "std": 892},
    "temperature": {"min": -15.2, "max": 32.1, "mean": 12.4, "std": 8.7},
    "precipitation": {"min": 89, "max": 3240, "mean": 987, "std": 445}
  },
  "validation_checks": {
    "elevation_continuity": {"status": "PASS", "outliers": 3},
    "flow_network": {"status": "WARNING", "orphaned_streams": 7},
    "climate_gradients": {"status": "PASS", "abrupt_transitions": 12},
    "biome_logic": {"status": "FAIL", "desert_in_high_precip": 5}
  },
  "anomalies": [
    {
      "type": "elevation_spike",
      "location": [45, -23],
      "value": 4567,
      "neighbors_avg": 1200,
      "severity": "HIGH"
    }
  ]
}
```

#### 2. Hex Neighborhood Analysis (for specific problem areas)
```json
{
  "center_hex": {
    "id": "45_-23",
    "elevation": 4567,
    "temperature": -8.2,
    "precipitation": 1200,
    "biome": "ALPINE",
    "soil": {"texture": "Rocky", "fertility": 0.1}
  },
  "neighbors": [
    {"id": "46_-23", "elevation": 1200, "temp_diff": 12.4},
    {"id": "45_-22", "elevation": 1150, "temp_diff": 11.8}
  ],
  "gradients": {
    "elevation_max_diff": 3367,
    "temperature_lapse_expected": -22.4,
    "temperature_lapse_actual": -8.2,
    "flow_direction": "SW",
    "flow_accumulation": 245
  },
  "issues": [
    "Elevation spike too extreme",
    "Temperature lapse rate inconsistent",
    "Should be snow/ice at this elevation"
  ]
}
```

#### 3. Flow Network Analysis (CSV for spreadsheet analysis)
```csv
hex_id,elevation,flow_to,accumulation,stream_order,velocity,issues
"0_0",1250,"1_0",15.2,1,0.8,""
"1_0",1245,"2_0",30.4,2,1.2,""
"15_-8",890,"15_-7",245.7,3,2.1,"orphaned_from_upstream"
```

#### 4. Cross-Section Analysis (for terrain validation)
```json
{
  "transect": "latitude_45N",
  "start": [-180, 45],
  "end": [180, 45],
  "sample_points": [
    {"lon": -180, "elev": 234, "temp": 8.2, "precip": 567},
    {"lon": -150, "elev": 1890, "temp": -4.1, "precip": 1200},
    {"lon": -120, "elev": 45, "temp": 12.8, "precip": 234}
  ],
  "analysis": {
    "rain_shadow_detected": true,
    "mountain_range_width": 180,
    "temperature_correlation": 0.85
  }
}
```

### Validation Rules and Statistical Checks

#### Expected Ranges and Relationships
```
VALIDATION_RULES = {
  "elevation": {
    "global_range": [-11000, 9000],  // Mariana Trench to Everest
    "neighbor_max_diff": 2000,       // Max elevation change between neighbors
    "land_sea_transition": [-200, 200] // Reasonable coastal transitions
  },
  
  "temperature": {
    "by_latitude": {
      "equator": [15, 35],
      "temperate": [-5, 25], 
      "polar": [-40, 5]
    },
    "lapse_rate": [-5, -8],           // °C per 1000m elevation
    "seasonal_range": [5, 40]         // Max annual temperature range
  },
  
  "precipitation": {
    "global_range": [0, 5000],        // mm/year
    "desert_threshold": 250,
    "rainforest_threshold": 2000,
    "orographic_enhancement": [1.0, 3.0] // Windward side multiplier
  },
  
  "flow_networks": {
    "max_accumulation_jump": 1000,    // Sudden flow increases indicate errors
    "stream_order_progression": true,  // Must increase downstream
    "outlet_requirement": true        // All streams must reach ocean/lake
  },
  
  "biome_logic": {
    "temperature_precipitation_bounds": {
      "DESERT": {"temp": [-10, 50], "precip": [0, 600]},
      "RAINFOREST": {"temp": [18, 35], "precip": [1500, 5000]},
      "TUNDRA": {"temp": [-15, 5], "precip": [100, 800]}
    },
    "elevation_limits": {
      "treeline": [1500, 4000],        // Varies by latitude
      "alpine_minimum": 2000
    }
  }
}
```

#### Automated Quality Checks
```
quality_check_pipeline() {
  // Phase 1: Basic data integrity
  check_null_values()
  check_coordinate_validity()
  check_neighbor_connectivity()
  
  // Phase 2: Physical realism
  validate_elevation_gradients()
  validate_temperature_lapse_rates()
  validate_precipitation_patterns()
  validate_flow_connectivity()
  
  // Phase 3: Logical consistency
  check_biome_climate_match()
  check_soil_parent_material_match()
  check_vegetation_climate_compatibility()
  
  // Phase 4: Statistical analysis
  analyze_distribution_shapes()
  detect_spatial_autocorrelation_breaks()
  identify_unrealistic_clusters()
  
  return QualityReport {
    overall_score: float,
    critical_issues: [],
    warnings: [],
    recommendations: []
  }
}

### Visualization Metadata Format

#### Image Metadata (embedded in JPEG/PNG)
```json
{
  "generator": "hex-world-gen v1.0",
  "timestamp": "2025-01-19T14:30:00Z",
  "world_seed": 42,
  "generation_stage": "post_biome_assignment",
  "view_config": {
    "layer_type": "composite_overview",
    "center_lat": 0,
    "center_lon": 0,
    "zoom_level": 1.0,
    "visible_layers": ["elevation", "water", "biomes"],
    "color_scheme": "realistic_earth"
  },
  "legend": {
    "elevation": {"min": -2840, "max": 4567, "unit": "meters"},
    "temperature": {"min": -15.2, "max": 32.1, "unit": "celsius"},
    "biomes": ["Ocean", "Desert", "Grassland", "Forest", "Alpine", "Tundra"]
  },
  "quality_score": 0.85,
  "known_issues": ["3 elevation outliers", "7 orphaned streams"]
}
```

#### Debugging Session Context
```json
{
  "session_id": "debug_2025_0119_001",
  "problem_description": "Rivers not flowing properly in mountain regions",
  "focus_area": {"lat_range": [40, 50], "lon_range": [-120, -110]},
  "generated_outputs": {
    "visual": "debug_mountain_rivers.jpg",
    "validation": "validation_report.json", 
    "neighborhood": "hex_analysis_45_-23.json",
    "flow_csv": "flow_network_mountain_region.csv"
  },
  "ai_analysis_prompt": "Focus on elevation gradients, flow accumulation, and stream order consistency in the mountain region. Look for: 1) Unrealistic elevation spikes, 2) Broken flow connectivity, 3) Streams flowing uphill"
}
```
```

## Implementation Deliverables

### Deliverable 1: Hex Grid Foundation (HEX-001)
**Goal**: Establish core hex grid system with coordinate management and topology support
**Demo Value**: Show working hex coordinate system, neighbor lookups, distance calculations, world vs region topology
**Dependencies**: None
**Core Features**:
- Axial coordinate system implementation
- Hex-to-pixel and pixel-to-hex conversion
- Configurable topology: World (toroidal wrapping) vs Region (bounded edges)
- Smart neighbor lookup algorithms (handles edge cases for both topologies)
- Distance calculations between hexes (shortest path on world maps)
- Basic hex grid data structure with topology awareness
**Test Coverage**: Coordinate math, neighbor detection, boundary conditions, wrapping behavior
**Success Criteria**: Can create NxM hex grid with either topology, query neighbors correctly, calculate distances

### Deliverable 2: Basic Terrain Generation (TERRAIN-001)  
**Goal**: Generate realistic elevation using fractal noise
**Demo Value**: Show elevation heatmap with realistic mountain ranges and valleys
**Dependencies**: HEX-001
**Core Features**:
- Diamond-Square algorithm implementation
- Multi-octave noise generation (Hurst exponent H=0.8-0.9)
- Hypsometric curve matching Earth's distribution
- Basic land/water designation
**Test Coverage**: Elevation ranges, noise consistency, statistical validation
**Success Criteria**: Realistic-looking terrain with proper elevation distribution

### Deliverable 3: Core Visualization System (VIZ-001)
**Goal**: Render hex grids with multiple layer support
**Demo Value**: Show elevation maps, debug overlays, coordinate display
**Dependencies**: HEX-001, TERRAIN-001
**Core Features**:
- JPEG/PNG export with metadata
- Color mapping for elevation data
- Hillshading for terrain visualization
- Debug overlay with hex coordinates
- Basic validation output (JSON stats)
**Test Coverage**: Image generation, color mapping accuracy, metadata embedding
**Success Criteria**: Can export publication-quality terrain maps

### Deliverable 4: Enhanced Terrain (TERRAIN-002)
**Goal**: Add slope, aspect, and tectonic features
**Demo Value**: Show slope analysis, aspect-based solar exposure, fault systems
**Dependencies**: TERRAIN-001, VIZ-001
**Core Features**:
- Slope calculation from elevation gradients
- Aspect determination (compass direction of slope)
- Basic tectonic simulation with Voronoi faults
- Continental shelf modeling
**Test Coverage**: Gradient calculations, aspect accuracy, fault generation
**Success Criteria**: Realistic slope/aspect maps, identifiable tectonic features

### Deliverable 5: Water Flow System (HYDRO-001)
**Goal**: Implement D8 flow accumulation and stream networks
**Demo Value**: Show realistic river networks flowing to oceans
**Dependencies**: TERRAIN-002, VIZ-001
**Core Features**:
- D8 flow direction algorithm
- Flow accumulation calculation
- Strahler stream ordering
- Basic stream network extraction
- Flow visualization with variable width
**Test Coverage**: Flow connectivity, stream ordering, mass conservation
**Success Criteria**: Realistic drainage networks, all streams reach outlets

### Deliverable 6: Climate Foundation (CLIMATE-001)
**Goal**: Calculate temperature and precipitation patterns
**Demo Value**: Show realistic climate zones and rain shadows
**Dependencies**: TERRAIN-002, HYDRO-001
**Core Features**:
- Latitude-based temperature calculation
- Elevation lapse rate application
- Orographic precipitation enhancement
- Basic rain shadow modeling
- Köppen climate zone assignment
**Test Coverage**: Temperature gradients, precipitation patterns, climate realism
**Success Criteria**: Earth-like climate distribution, proper rain shadows

### Deliverable 7: Soil Profile System (SOIL-001)
**Goal**: Generate realistic soil profiles from geology and climate
**Demo Value**: Show soil fertility maps, drainage classification
**Dependencies**: TERRAIN-002, CLIMATE-001
**Core Features**:
- Soil texture determination (parent material + climate)
- Fertility calculation with multiple factors
- Drainage classification (fast/moderate/poor)
- Erosion risk assessment
- pH and nutrient level modeling
**Test Coverage**: Soil-climate relationships, fertility accuracy, drainage logic
**Success Criteria**: Realistic soil distribution matching climate/geology

### Deliverable 8: Biome Classification (BIOME-001)
**Goal**: Classify biomes using enhanced Whittaker model
**Demo Value**: Show realistic biome distribution map
**Dependencies**: CLIMATE-001, SOIL-001
**Core Features**:
- Whittaker diagram implementation
- Soil and slope modifiers
- Treeline calculation by latitude
- Biome transition smoothing
- Vegetation density calculation
**Test Coverage**: Biome boundaries, climate-biome correlation, transition realism
**Success Criteria**: Earth-like biome distribution, smooth transitions

### Deliverable 9: Ocean Current System (HYDRO-002)
**Goal**: Model realistic ocean currents with Coriolis effects
**Demo Value**: Show ocean circulation patterns, coastal upwelling
**Dependencies**: CLIMATE-001, HYDRO-001
**Core Features**:
- Wind-driven surface currents
- Coriolis deflection (45°)
- Gyre formation (subtropical/subpolar)
- Thermohaline circulation basics
- Current visualization arrows
**Test Coverage**: Current directions, Coriolis accuracy, gyre formation
**Success Criteria**: Realistic ocean circulation patterns

### Deliverable 10: Enhanced Visualization (VIZ-002)
**Goal**: Advanced debugging and multi-layer visualization
**Demo Value**: Show comprehensive debugging tools, validation reports
**Dependencies**: All previous deliverables
**Core Features**:
- Multi-layer composite rendering
- Validation report generation
- Neighborhood analysis tools
- Cross-section analysis
- Statistical anomaly detection
**Test Coverage**: Validation accuracy, debug tool functionality, report generation
**Success Criteria**: Comprehensive debugging toolset for world validation

### Deliverable 11: Weather Systems (WEATHER-001)
**Goal**: Dynamic weather simulation with pressure systems
**Demo Value**: Show weather fronts, storms, daily variation
**Dependencies**: CLIMATE-001, HYDRO-002
**Core Features**:
- Air pressure system modeling
- Weather front detection/propagation
- Daily weather updates
- Storm generation (thunderstorms, hurricanes)
- Seasonal variation
**Test Coverage**: Pressure calculations, front movement, weather realism
**Success Criteria**: Realistic weather patterns and storm systems

### Deliverable 12: Erosion Simulation (EROSION-001)
**Goal**: Implement stream power erosion with sediment transport
**Demo Value**: Show terrain evolution, valley carving, sediment deposition
**Dependencies**: HYDRO-001, SOIL-001
**Core Features**:
- Stream power law implementation
- Sediment capacity calculation
- Erosion/deposition cycling
- Terrain modification over time
- Before/after visualization
**Test Coverage**: Mass conservation, erosion rates, realistic valley formation
**Success Criteria**: Realistic terrain evolution, proper sediment transport

### Feature Dependencies Flow:
```
HEX-001 → TERRAIN-001 → VIZ-001
    ↓         ↓          ↓
TERRAIN-002 → HYDRO-001 → CLIMATE-001
    ↓         ↓          ↓
SOIL-001 → BIOME-001   HYDRO-002
    ↓         ↓          ↓
    VIZ-002 ← WEATHER-001
    ↓
EROSION-001
```

### 12-Step Deliverable Workflow

Each deliverable follows this exact cycle:

#### Steps 1-2: Planning Phase
1. **Review Design & README**: Study this document and current README for context
2. **Review Implementation**: Examine existing codebase, understand current architecture

#### Steps 3-5: Setup Phase  
3. **Create Feature Branch**: `git checkout -b feature/[DELIVERABLE-ID]`
4. **Create Feature Plan**: Document in `docs/features/[DELIVERABLE-ID].md`:
   ```markdown
   # [DELIVERABLE-ID]: [Title]
   
   ## Objective
   [What this deliverable accomplishes]
   
   ## Technical Approach
   [How it will be implemented]
   
   ## API Design
   [Go interfaces, structs, function signatures]
   
   ## File Structure
   [New files, modified files]
   
   ## Testing Strategy
   [Test cases, validation criteria]
   ```

5. **Create Failing Tests**: Write comprehensive test cases that define expected behavior

#### Steps 6-8: Implementation Phase
6. **Implement Feature**: 
   - Write idiomatic Go code with proper documentation
   - Commit frequently with descriptive messages
   - Ensure all tests pass
   - Follow existing code patterns and conventions

7. **Add Additional Tests**: Edge cases, integration tests, performance tests

8. **Clean Up**: Remove deprecated functions, unused imports, dead code

#### Steps 9-12: Integration Phase
9. **Run All Tests**: `go test ./...` - ensure no regressions
10. **Commit & Push**: Final commit with clean history
11. **Create PR**: 
    - Title: `[DELIVERABLE-ID]: [Brief description]`
    - Description: Demo instructions, testing notes, validation data
    - Request review and demo

12. **Post-Approval**: Update design doc, README, merge to main

### Demo Requirements Per Deliverable

#### HEX-001 Demo:
```bash
# Show coordinate system working with region topology
./hex-world demo-coords --size=20x20 --topology=region
# Output: Grid with coordinate labels, neighbor highlighting, edge hexes marked

# Show world topology with wrapping
./hex-world demo-coords --size=20x20 --topology=world
# Output: Grid showing wrapped neighbors at edges

# Show distance calculations  
./hex-world demo-distance --from=0,0 --to=5,3 --topology=region
./hex-world demo-distance --from=0,0 --to=19,0 --topology=world
# Output: Distance calculation with path visualization, wrapping on world maps
```

#### TERRAIN-001 Demo:
```bash
# Generate and export terrain
./hex-world generate-terrain --size=100x100 --seed=42
# Output: terrain_seed42.jpg with elevation heatmap

# Show statistics
./hex-world terrain-stats terrain_seed42.json
# Output: Elevation range, distribution, validation report
```

#### VIZ-001 Demo:
```bash
# Show multiple render modes
./hex-world render --input=terrain.json --mode=elevation
./hex-world render --input=terrain.json --mode=debug --show-coords=true
# Output: Different visualization styles

# Show metadata
./hex-world metadata terrain_elevation.jpg
# Output: Embedded metadata and generation parameters
```

### Expected File Structure by Deliverable

#### HEX-001:
```
pkg/
  hex/
    coordinate.go          // Axial coordinate system
    coordinate_test.go     // Coordinate math tests
    grid.go               // Grid data structure  
    grid_test.go          // Grid operations tests
    neighbor.go           // Neighbor algorithms
    neighbor_test.go      // Neighbor logic tests
cmd/
  hex-world/
    main.go              // CLI entry point
    demo.go              // Demo commands
```

#### TERRAIN-001:
```
pkg/
  terrain/
    generator.go         // Noise generation
    generator_test.go    // Terrain generation tests
    heightmap.go         // Elevation data structure
    noise.go            // Fractal noise algorithms
    validation.go        // Statistical validation
internal/
  noise/
    diamond_square.go    // Diamond-square implementation
```

#### VIZ-001:
```
pkg/
  render/
    renderer.go          // Core rendering engine
    colormap.go          // Color mapping functions
    export.go           // JPEG/PNG export with metadata
    hillshade.go        // Hillshading algorithms
  validation/
    stats.go            // Statistical analysis
    report.go           // Validation report generation
```

### Testing Standards

#### Unit Tests (Required):
- All public functions must have tests
- Test coverage >80% for core algorithms
- Table-driven tests for mathematical functions
- Property-based testing for coordinate operations

#### Integration Tests (Required):
- End-to-end pipeline testing
- File I/O validation
- Image generation verification
- Cross-platform compatibility

#### Validation Tests (Required):  
- Statistical distribution checks
- Physical realism validation
- Performance benchmarks
- Memory usage profiling

### Commit Message Standards
```
[DELIVERABLE-ID]: Brief description

Detailed explanation of changes:
- Added coordinate system implementation
- Implemented neighbor lookup algorithms  
- Added comprehensive test coverage

Tests: All passing (47 tests, 0 failures)
Performance: Distance calc ~0.1ms per operation
Validation: Coordinate math matches reference implementation
```

## Key Algorithms

### Ocean Current Generation
```
// Surface current from wind stress
wind_stress = air_density * drag_coeff * wind_speed^2 * wind_direction
surface_current = wind_stress / (water_density * coriolis_parameter)
// Apply 45° Coriolis deflection
deflected_current = rotate_vector(surface_current, 45°)
// Smooth with neighbors for realistic flow
smoothed_current = gaussian_blur(deflected_current, radius=3)
```

### Water Flow (Enhanced D8 Algorithm)
```
for each hex in sorted_by_elevation_desc:
  flow_dir = lowest_neighbor(hex)
  downstream[flow_dir].accumulation += hex.accumulation + hex.rainfall
  downstream[flow_dir].velocity = manning_equation(slope, accumulation)
  // Calculate stream order (Strahler)
  if accumulation > threshold: assign_stream_order(hex)
```

### Stream Power Erosion
```
stream_power = water_density * gravity * discharge * slope
erosion_rate = K * (drainage_area^0.5) * (slope^1.0)
carrying_capacity = C * velocity^2
if carrying_capacity < sediment_load:
  deposition_rate = (sediment_load - carrying_capacity) * settling_velocity
sediment_flux = erosion - deposition
```

### Soil Profile Generation
```
base_fertility = rock_fertility_map[rock_type]
climate_modifier = min(temperature/25, precipitation/1000) 
organic_modifier = vegetation_density * decomposition_rate
fertility = base_fertility * climate_modifier * organic_modifier
drainage = texture_drainage[soil_texture] * slope_factor
erosion_risk = slope * precipitation * (1 - vegetation_cover)
```

### Weather System Updates
```
// Daily pressure update
pressure_change = -divergence(wind_field) * time_step
new_pressure = old_pressure + pressure_change
// Geostrophic wind from pressure gradient
pressure_gradient = gradient(pressure_field)
geostrophic_wind = -pressure_gradient / (density * coriolis)
// Weather front detection
temperature_gradient = gradient(temperature_field)
if |temperature_gradient| > threshold: mark_as_front(hex)
```

### Enhanced Biome Classification
```
base_biome = whittaker_classification(temp, precip)
// Soil and slope modifiers
if slope > 30: base_biome = modify_for_steep_terrain(base_biome)
if soil.drainage == "Poor": base_biome = modify_for_wetland(base_biome)
if soil.fertility < 0.3: reduce_forest_density(base_biome)
// Elevation modifiers
if elevation > treeline_elevation(latitude): base_biome = ALPINE
return base_biome
```

## Real-World Parameters and Constants

### Physical Constants
- **Gravity**: 9.81 m/s²
- **Water Density**: 1000 kg/m³
- **Sediment Density**: 2650 kg/m³
- **Air Pressure at Sea Level**: 101.325 kPa
- **Stefan-Boltzmann**: 5.67×10⁻⁸ W/m²K⁴

### Geological Parameters
- **Erosion Rate Constant (K)**: 1×10⁻⁶ to 1×10⁻⁴ /year
- **Soil Porosity**: 0.3-0.5 (dimensionless)
- **Rock Density**: 2200-2800 kg/m³
- **Continental Crust Thickness**: 30-50 km
- **Isostatic Adjustment Rate**: 1-10 mm/year

### Hydrological Parameters
- **Manning's n**: 0.025-0.1 (natural channels)
- **Infiltration Rate**: 0.1-50 mm/hour (by soil type)
- **Evapotranspiration**: 500-1500 mm/year
- **Runoff Coefficient**: 0.1-0.95 (by surface type)

### Climate Parameters
- **Lapse Rate**: 6.5°C/km
- **Relative Humidity**: 40-80% typical
- **Wind Speed**: 2-20 m/s (surface)
- **Solar Constant**: 1361 W/m²
- **Albedo**: 0.1-0.9 (water to snow)

## Rendering Considerations
- LOD system for distant hexes
- Efficient batching for hex mesh rendering
- Texture atlasing for biome visuals
- Water shader with flow visualization
- Snapshot comparison view (side-by-side or overlay)
