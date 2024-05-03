package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math"
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
