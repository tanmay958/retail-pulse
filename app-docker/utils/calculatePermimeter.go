package utils

import (
	"fmt"
	"image"
	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"net/http"
)

func CalculatePerimeter(imageURL string) (float64, error) {
	// Fetch the image from the URL
	resp, err := http.Get(imageURL)
	if err != nil {
		return 0, fmt.Errorf("error fetching image: %w", err)
	}
	defer resp.Body.Close()

	// Decode the image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error decoding image: %w", err)
	}

	// Get image dimensions
	bounds := img.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	// Calculate perimeter
	perimeter := 2 * (width + height)
	return perimeter, nil
}
