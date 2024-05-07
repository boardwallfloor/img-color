package main

import (
	"fmt"
	"image/color"
)

type XYZ struct {
	X      float64
	Y      float64
	Z      float64
	scaled bool
}

var xyzCoefficient = [3][3]float64{
	{0.4124564, 0.2126568, 0.0193339},
	{0.3575761, 0.7151522, 0.1191920},
	{0.1804375, 0.0721890, 0.9503041},
}

var xyzD5Linear = [3]float64{0.95047, 1.00000, 1.08883}

func toXYZ(rgb color.RGBA) XYZ {
	valR, valG, valB := float64(rgb.R)/255.0, float64(rgb.G)/255.0, float64(rgb.G)/255.0
	var newXYZ XYZ
	for _, v := range xyzCoefficient {
		newXYZ.X += v[0] * float64(valR)
		newXYZ.Y += v[1] * float64(valG)
		newXYZ.Z += v[2] * float64(valB)
		fmt.Println(v[2] * float64(valB))
	}
	// linearize xyz with d65 as base
	newXYZ.X /= xyzD5Linear[0]
	newXYZ.Y /= xyzD5Linear[1]
	newXYZ.Z /= xyzD5Linear[2]

	return newXYZ
}

func getc2k() {
	xyzVal := toXYZ(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	fmt.Println(xyzVal.X, xyzVal.Y, xyzVal.Z)
}
