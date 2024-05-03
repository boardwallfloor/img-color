package main

import "image/color"

type XYZ struct {
	X float64
	Y float64
	Z float64
}

var xyzCoefficient = [3][3]float64{
	{0.4124564, 0.3575761, 0.1804375},
	{0.2126568, 0.7151522, 0.0721890},
	{0.0193339, 0.2653060, 0.7143281},
}

func toXYZ(rgb color.RGBA) XYZ {
	normR, normG, normB := rgb.R/255.0, rgb.G/255.0, rgb.G/255.0
	var newXYZ XYZ
	for _, v := range xyzCoefficient {
		newXYZ.X += v[0] * float64(normR)
		newXYZ.Y += v[1] * float64(normG)
		newXYZ.Z += v[2] * float64(normB)
	}
	// linearize xyz with d65 as base
	return XYZ{}
}
