package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func stdProcess(imgRef *vips.ImageRef) {
	imgByte, _, err := imgRef.ExportJpeg(vips.NewJpegExportParams())
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}
	img, err := jpeg.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}
	// Get image bounds (width and height)
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Iterate over each pixel in the image
	imgRGB := make([][]RGB, 0)
	for y := 0; y < height; y++ {
		tempArr := make([]RGB, 0)
		for x := 0; x < width; x++ {
			// Get the color of the pixel at coordinates (x, y)
			r, g, b, _ := img.At(x, y).RGBA()

			// Convert from 16-bit color to 8-bit color
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// Print the RGB values of the pixel
			// fmt.Printf("Pixel at (%d, %d): R=%d, G=%d, B=%d\n", x, y, r8, g8, b8)
			tempArr = append(tempArr, RGB{r8, g8, b8})
		}
		imgRGB = append(imgRGB, tempArr)
	}
	palette := []string{
		"#2e3440", "#3b4252", "#434c5e", "#4c566a",
		"#d8dee9", "#e5e9f0", "#eceff4", "#8fbcbb",
		"#88c0d0", "#81a1c1", "#5e81ac", "#bf616a",
		"#d08770", "#ebcb8b", "#a3be8c", "#b48ead",
	}
	paletteRGB := make([]color.RGBA, 0)
	mapPalette := make(map[color.RGBA]string, 0)
	for _, v := range palette {
		rgb, err := hexToRGB(v)
		if err != nil {
			log.Fatalf("skill issue at , :%s\n", err)
		}
		paletteRGB = append(paletteRGB, rgb)
		mapPalette[rgb] = v
	}
	// fmt.Println(paletteRGB)

	convertedByte := make([][]color.RGBA, 0)
	paletteCount := make(map[string]int, 0)
	for _, col := range imgRGB {
		tempArr := make([]color.RGBA, 0)
		for _, cell := range col {

			newRgb, _, _, err := getClosest(cell, paletteRGB)
			// fmt.Println(cell, newRgb)
			if err != nil {
				log.Fatalf("skill issue at , :%s\n", err)
			}
			// fmt.Println(newRgb)
			hex := mapPalette[newRgb]
			paletteCount[hex] += 1
			tempArr = append(tempArr, newRgb)
		}
		convertedByte = append(convertedByte, tempArr)
	}
	fmt.Println(paletteCount)

	newImage := image.NewRGBA(image.Rect(0, 0, imgRef.Width(), imgRef.Height()))
	for y := 0; y < imgRef.Height(); y++ {
		for x := 0; x < imgRef.Width(); x++ {
			newImage.Set(x, y, convertedByte[y][x])
		}
	}
	filename := "4k_conv.png"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Encode the image as PNG and write it to the file
	err = jpeg.Encode(file, newImage, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Image created and saved as %s\n", filename)
}

func createGradient() {
	// Define image dimensions
	width := 800
	height := 600

	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Define colors for the gradient
	startColor := color.RGBA{255, 0, 0, 255} // Red
	endColor := color.RGBA{0, 0, 255, 255}   // Blue
	colorRange := width                      // Number of pixels for the gradient

	// Create the gradient
	for x := 0; x < width; x++ {
		r := startColor.R + uint8(x*(int(endColor.R)-int(startColor.R))/colorRange)
		g := startColor.G + uint8(x*(int(endColor.G)-int(startColor.G))/colorRange)
		b := startColor.B + uint8(x*(int(endColor.B)-int(startColor.B))/colorRange)
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Create a PNG file
	file, err := os.Create("gradient.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Encode the image as PNG and write it to the file
	err = jpeg.Encode(file, img, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Gradient image created and saved as gradient.png")
}
