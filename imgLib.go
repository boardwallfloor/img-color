package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

type ProcessStdLib struct {
	saveDir     string
	filePath    string
	newFileName string
}

// type ImgProcess struct {
// 	prc *ProcessStdLib
// }

// func NewImgProcess(opts ...func(*ProcessStdLib)) *ProcessStdLib {
// 	prc := &ProcessStdLib{}
// 	for _,v := range opts {
// 		v(prc)
// 	}
// 	return &ProcessStdLib{}
// }

// func SetFilePath()

func (prc *ProcessStdLib) Process(imgRef *vips.ImageRef, palette []string) {
	imgByte, _, err := imgRef.ExportPng(vips.NewPngExportParams())
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}

	img, err := png.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

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

	paletteRGB := make([]color.RGBA, len(palette))
	for i, v := range palette {
		rgb, err := hexToRGB(v)
		if err != nil {
			log.Fatalf("skill issue at , :%s\n", err)
		}
		paletteRGB[i] = rgb
	}

	convertedByte := make([][]color.RGBA, height)
	for y, col := range imgRGB {
		tempArr := make([]color.RGBA, width)
		for x, cell := range col {

			newRgb, _, _, err := getClosest(cell, paletteRGB)
			if err != nil {
				log.Fatalf("skill issue at , :%s\n", err)
			}
			tempArr[x] = newRgb
		}
		convertedByte[y] = tempArr
	}

	newImage := image.NewRGBA(image.Rect(0, 0, imgRef.Width(), imgRef.Height()))
	for y := 0; y < imgRef.Height(); y++ {
		for x := 0; x < imgRef.Width(); x++ {
			newImage.Set(x, y, convertedByte[y][x])
		}
	}
	file, err := os.Create(fmt.Sprintf("%s.png", prc.newFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Encode the image as PNG and write it to the file
	err = png.Encode(file, newImage)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Image created and saved as %s\n", prc.newFileName)
}
