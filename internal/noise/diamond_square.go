package noise

import (
	"math"
	"math/rand"
)

// DiamondSquare generates fractal terrain using the Diamond-Square algorithm
// Size must be (2^n + 1) for proper algorithm operation
func DiamondSquare(size int, roughness float64, seed int64) [][]float64 {
	// Validate size is (2^n + 1)
	if !isPowerOfTwoPlusOne(size) {
		panic("DiamondSquare: size must be (2^n + 1), e.g., 129, 257, 513")
	}
	
	rng := rand.New(rand.NewSource(seed))
	heightmap := make([][]float64, size)
	for i := range heightmap {
		heightmap[i] = make([]float64, size)
	}
	
	// Initialize corners with random values
	heightmap[0][0] = rng.Float64()*2 - 1         // Top-left
	heightmap[0][size-1] = rng.Float64()*2 - 1    // Top-right
	heightmap[size-1][0] = rng.Float64()*2 - 1    // Bottom-left
	heightmap[size-1][size-1] = rng.Float64()*2 - 1 // Bottom-right
	
	// Current step size starts at full grid and halves each iteration
	stepSize := size - 1
	scale := roughness
	
	for stepSize > 1 {
		halfStep := stepSize / 2
		
		// Diamond step: set center points of squares
		for y := halfStep; y < size; y += stepSize {
			for x := halfStep; x < size; x += stepSize {
				// Average the four corner values
				avg := (heightmap[y-halfStep][x-halfStep] + // Top-left
					heightmap[y-halfStep][x+halfStep] + // Top-right
					heightmap[y+halfStep][x-halfStep] + // Bottom-left
					heightmap[y+halfStep][x+halfStep]) / 4.0 // Bottom-right
				
				// Add random offset scaled by current roughness
				heightmap[y][x] = avg + (rng.Float64()*2-1)*scale
			}
		}
		
		// Square step: set center points of diamonds
		for y := 0; y < size; y += halfStep {
			for x := (y+halfStep)%stepSize; x < size; x += stepSize {
				// Calculate diamond center by averaging neighbors
				// Handle edge wrapping for seamless terrain
				avg := diamondAverage(heightmap, x, y, halfStep, size)
				
				// Add random offset
				heightmap[y][x] = avg + (rng.Float64()*2-1)*scale
			}
		}
		
		// Reduce step size and roughness for next iteration
		stepSize /= 2
		scale *= roughness // Scale factor controls how rough the terrain is
	}
	
	return heightmap
}

// diamondAverage calculates the average of diamond neighbors with edge wrapping
func diamondAverage(heightmap [][]float64, x, y, halfStep, size int) float64 {
	count := 0
	sum := 0.0
	
	// Check four diamond neighbors (up, down, left, right)
	neighbors := [][2]int{
		{x, y - halfStep}, // Up
		{x, y + halfStep}, // Down
		{x - halfStep, y}, // Left
		{x + halfStep, y}, // Right
	}
	
	for _, neighbor := range neighbors {
		nx, ny := neighbor[0], neighbor[1]
		
		// Handle edge wrapping for seamless terrain
		if nx < 0 {
			nx = size - 1
		} else if nx >= size {
			nx = 0
		}
		
		if ny < 0 {
			ny = size - 1
		} else if ny >= size {
			ny = 0
		}
		
		// Only include if the neighbor has been set (non-zero or explicitly set)
		if nx >= 0 && nx < size && ny >= 0 && ny < size {
			sum += heightmap[ny][nx]
			count++
		}
	}
	
	if count > 0 {
		return sum / float64(count)
	}
	return 0.0
}

// isPowerOfTwoPlusOne checks if n is of the form (2^k + 1)
func isPowerOfTwoPlusOne(n int) bool {
	if n < 3 {
		return false // Minimum is 2^1 + 1 = 3
	}
	n-- // Convert to 2^k
	return n > 0 && (n&(n-1)) == 0
}

// MultiOctaveNoise combines multiple octaves of Diamond-Square noise
func MultiOctaveNoise(width, height int, octaves int, persistence, lacunarity, scale float64, seed int64) [][]float64 {
	// Find the smallest power-of-two-plus-one size that fits our target
	noiseSize := nextPowerOfTwoPlusOne(max(width, height))
	
	result := make([][]float64, height)
	for i := range result {
		result[i] = make([]float64, width)
	}
	
	amplitude := 1.0
	frequency := scale
	maxValue := 0.0
	
	for octave := 0; octave < octaves; octave++ {
		// Generate noise for this octave
		octaveSeed := seed + int64(octave*1000)
		octaveNoise := DiamondSquare(noiseSize, 0.5, octaveSeed)
		
		// Add this octave to the result
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Sample from the noise using frequency scaling
				noiseX := int(float64(x) * frequency) % noiseSize
				noiseY := int(float64(y) * frequency) % noiseSize
				
				if noiseX < 0 {
					noiseX += noiseSize
				}
				if noiseY < 0 {
					noiseY += noiseSize
				}
				
				result[y][x] += octaveNoise[noiseY][noiseX] * amplitude
			}
		}
		
		maxValue += amplitude
		amplitude *= persistence
		frequency *= lacunarity
	}
	
	// Normalize to [-1, 1] range
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			result[y][x] /= maxValue
		}
	}
	
	return result
}

// nextPowerOfTwoPlusOne finds the smallest (2^n + 1) >= size
func nextPowerOfTwoPlusOne(size int) int {
	if size <= 1 {
		return 3 // Minimum is 2^1 + 1 = 3
	}
	
	// If already a power of two plus one, return it
	if isPowerOfTwoPlusOne(size) {
		return size
	}
	
	// Find next power of two plus one
	n := 2
	for n+1 < size {
		n *= 2
	}
	return n + 1
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SpectralSynthesis generates terrain using spectral synthesis with power law
// Beta controls the power spectrum: β ≈ 2 gives realistic terrain
func SpectralSynthesis(width, height int, beta float64, seed int64) [][]float64 {
	rng := rand.New(rand.NewSource(seed))
	
	// Create frequency domain representation
	freqWidth := width / 2
	freqHeight := height / 2
	
	result := make([][]float64, height)
	for i := range result {
		result[i] = make([]float64, width)
	}
	
	// Generate in frequency domain
	for fy := 0; fy < freqHeight; fy++ {
		for fx := 0; fx < freqWidth; fx++ {
			// Calculate frequency magnitude
			freq := math.Sqrt(float64(fx*fx + fy*fy))
			if freq == 0 {
				freq = 1 // Avoid division by zero
			}
			
			// Power law amplitude: A(f) = 1/f^(β/2)
			amplitude := 1.0 / math.Pow(freq, beta/2.0)
			
			// Random phase
			phase := rng.Float64() * 2 * math.Pi
			
			// Generate spatial domain value (simplified inverse FFT)
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					spatial := 2*math.Pi*(float64(fx*x)/float64(width) + float64(fy*y)/float64(height))
					result[y][x] += amplitude * math.Cos(spatial + phase)
				}
			}
		}
	}
	
	// Normalize to [-1, 1]
	minVal, maxVal := findMinMax(result)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			result[y][x] = 2*(result[y][x]-minVal)/(maxVal-minVal) - 1
		}
	}
	
	return result
}

// findMinMax finds the minimum and maximum values in a 2D array
func findMinMax(data [][]float64) (float64, float64) {
	if len(data) == 0 || len(data[0]) == 0 {
		return 0, 0
	}
	
	minVal := data[0][0]
	maxVal := data[0][0]
	
	for _, row := range data {
		for _, val := range row {
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
		}
	}
	
	return minVal, maxVal
}