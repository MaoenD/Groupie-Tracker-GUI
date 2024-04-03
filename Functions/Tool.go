package Functions

import (
	"fmt"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// getAverageColor calculates the average color of an image.
func getAverageColor(img *canvas.Image) (r, g, b, a uint32) {
	// Get the width and height of the image
	width := int(img.MinSize().Width)
	height := int(img.MinSize().Height)
	var totalR, totalG, totalB, totalA uint64

	// Iterate over all pixels of the image to calculate the sum of color components
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the color of the current pixel
			colorRGBA := img.Image.At(x, y)
			r, g, b, a := colorRGBA.RGBA()
			// Add the color components to the total sum
			totalR += uint64(r >> 8)
			totalG += uint64(g >> 8)
			totalB += uint64(b >> 8)
			totalA += uint64(a >> 8)
		}
	}

	// Calculate the rounded averages
	pixelCount := uint64(width * height)
	r = uint32(totalR / pixelCount)
	g = uint32(totalG / pixelCount)
	b = uint32(totalB / pixelCount)
	a = uint32(totalA / pixelCount)

	// Check if the average color components are close to 255
	maxComponent := uint32(255)
	tolerance := uint32(50) // Tolerance value to adjust the color
	if r >= maxComponent-tolerance && g >= maxComponent-tolerance && b >= maxComponent-tolerance {
		// Reduce the red, green, and blue components to avoid pure white
		r = uint32(math.Max(0, float64(r)-float64(tolerance)))
		g = uint32(math.Max(0, float64(g)-float64(tolerance)))
		b = uint32(math.Max(0, float64(b)-float64(tolerance)))
	}

	return r, g, b, a
}

// LoadImageResource loads an image resource from the specified path.
func LoadImageResource(path string) fyne.Resource {
	// Load a resource (image) from the specified path
	image, err := fyne.LoadResourceFromPath(path)
	// Check if there was an error loading the image
	if err != nil {
		// Print an error message if loading the image failed
		fmt.Println("Error loading icon:", err)
		// Return nil if there was an error loading the image
		return nil
	}
	// Return the successfully loaded resource (image)
	return image
}

// checkMemberName checks if a member's name contains the search text.
func checkMemberName(members []string, searchText string) bool {
	// Iterate over each member in the list
	for _, member := range members {
		// Check if the member's name contains the search text (case-insensitive)
		if strings.Contains(strings.ToLower(member), searchText) {
			return true // Return true if a match is found
		}
	}
	return false // Return false if no match is found
}

// checkConcertLocation checks if a concert's location contains the search text.
func checkConcertLocation(concerts []Concert, searchText string) bool {
	// Iterate over each concert in the list
	for _, concert := range concerts {
		// Check if the concert's location contains the search text (case-insensitive)
		for _, location := range concert.Locations {
			if strings.Contains(strings.ToLower(string(location)), searchText) {
				return true // Return true if a match is found
			}
		}
	}
	return false // Return false if no match is found
}

// contains checks if a string is present in a slice of strings.
func contains(slice []string, str string) bool {
	// Check if a string is present in a slice of strings
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
