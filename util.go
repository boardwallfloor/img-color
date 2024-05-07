package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"
)

func hexToRGB(hex string) (color.RGBA, error) {
	// Remove any leading '#' from the hex string
	hex = strings.TrimPrefix(hex, "#")

	// Parse the hex string into three parts: red, green, and blue
	red, err := strconv.ParseInt(hex[0:2], 16, 0)
	if err != nil {
		return color.RGBA{}, err
	}
	green, err := strconv.ParseInt(hex[2:4], 16, 0)
	if err != nil {
		return color.RGBA{}, err
	}
	blue, err := strconv.ParseInt(hex[4:6], 16, 0)
	if err != nil {
		return color.RGBA{}, err
	}

	// Create a color.RGBA struct with the parsed values
	rgba := color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: 255, // Alpha value (fully opaque)
	}

	return rgba, nil
}

func getClosest(rgb RGB, palette []color.RGBA) (color.RGBA, float64, int, error) {
	if len(palette) == 0 {
		return color.RGBA{}, 0.0, 0, errors.New("empty palette")
	}

	valR, valG, valB := float64(rgb.R)/255.0, float64(rgb.G)/255.0, float64(rgb.G)/255.0

	minDist := 0.0
	minCol := 0
	ass := false
	var lowestCount int

	for c, v := range palette {
		normR, normG, normB := float64(v.R)/255.0, float64(v.G)/255.0, float64(v.G)/255.0
		rCalc := math.Pow(float64(valR-normR), 2)
		gCalc := math.Pow(float64(valG-normG), 2)
		bCalc := math.Pow(float64(valB-normB), 2)
		dist := math.Sqrt(rCalc + gCalc + bCalc)
		// log.Fatalln(normR, normG, normB, valR, rgb.R, valG, rgb.G, valB, rgb.G, dist)
		if !ass {
			minDist = dist
			ass = true
			minCol = c
			lowestCount = 1
		} else if dist == minDist {
			lowestCount++
		} else if dist < minDist {
			minDist = dist
			minCol = c
			lowestCount = 1
		}
	}

	return palette[minCol], minDist, lowestCount, nil
}

func getRgb(xStart, xEnd, yStart, yEnd int, wg *sync.WaitGroup, img *vips.ImageRef) {
	defer wg.Done()
	rgb := make([][]float64, 0)
	for y := yStart; y < yEnd; y++ {
		for x := xStart; x < xEnd; x++ {
			val, err := img.GetPoint(x, y)
			if err != nil {
				log.Fatalf("skill issue at , :%s\n", err)
			}
			rgb = append(rgb, val)

		}
	}
	fmt.Println(len(rgb))
	fmt.Println("done")
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

func getRGBArray(height, width int, img image.Image) [][]RGB {
	imgRGB := make([][]RGB, height)
	for y := 0; y < height; y++ {
		tempArr := make([]RGB, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			tempArr[x] = RGB{r8, g8, b8}
		}
		imgRGB[y] = tempArr
	}
	return imgRGB
}
