package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Result represents the segmentation result
type Result struct {
	OriginalImage  string `json:"original_image"`
	SegmentedImage string `json:"segmented_image"`
	Message        string `json:"message"`
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// performImageSegmentation performs basic image segmentation
func performImageSegmentation(inputPath string, outputPath string) error {
	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening image: %v", err)
	}
	defer file.Close()

	// Decode the image
	var img image.Image
	var decodeErr error

	if strings.HasSuffix(strings.ToLower(inputPath), ".png") {
		img, decodeErr = png.Decode(file)
	} else {
		img, decodeErr = jpeg.Decode(file)
	}

	if decodeErr != nil {
		return fmt.Errorf("error decoding image: %v", decodeErr)
	}

	// Get image bounds
	bounds := img.Bounds()
	// width := bounds.Max.X - bounds.Min.X
	// height := bounds.Max.Y - bounds.Min.Y

	// Create a new RGBA image
	segmented := image.NewRGBA(bounds)

	// Simple thresholding for segmentation
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			r, g, b, _ := color.RGBAModel.Convert(pixel).RGBA()
			
			// Calculate grayscale value
			gray := (r + g + b) / 3
			
			// Simple threshold
			if gray > 32768 { // 32768 is middle value (65535/2)
				segmented.Set(x, y, color.RGBA{255, 255, 255, 255}) // White
			} else {
				segmented.Set(x, y, color.RGBA{0, 0, 0, 255}) // Black
			}
		}
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer out.Close()

	// Encode and save the segmented image
	if strings.HasSuffix(strings.ToLower(outputPath), ".png") {
		err = png.Encode(out, segmented)
	} else {
		err = jpeg.Encode(out, segmented, &jpeg.Options{Quality: 90})
	}

	if err != nil {
		return fmt.Errorf("error encoding output image: %v", err)
	}

	return nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 10MB max memory
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
		return
	}

	// Create unique filenames for original and segmented images
	originalPath := filepath.Join(uploadsDir, "original_"+handler.Filename)
	segmentedPath := filepath.Join(uploadsDir, "segmented_"+handler.Filename)

	// Save original file
	dst, err := os.Create(originalPath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Perform image segmentation
	err = performImageSegmentation(originalPath, segmentedPath)
	if err != nil {
		http.Error(w, "Error performing segmentation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	result := Result{
		OriginalImage:  "/uploads/original_" + handler.Filename,
		SegmentedImage: "/uploads/segmented_" + handler.Filename,
		Message:        "Image segmentation completed successfully",
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	// Serve static files from the uploads directory
	fs := http.FileServer(http.Dir("uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// Handle upload endpoint
	http.HandleFunc("/api/upload", enableCORS(uploadHandler))

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
