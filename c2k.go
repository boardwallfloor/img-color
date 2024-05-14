package main

import (
	"fmt"
	"image/color"
	"math"
)

type XYZ struct {
	X      float64
	Y      float64
	Z      float64
	scaled bool
}

type LAB struct {
	L float64
	a float64
	b float64
}

var xyzCoefficient = [3][3]float64{
	{0.4124564, 0.2126568, 0.0193339},
	{0.3575761, 0.7151522, 0.1191920},
	{0.1804375, 0.0721890, 0.9503041},
}

var xyzD5Linear = [3]float64{0.95047, 1.00000, 1.08883}

const (
	kappa   = 903.3
	epsilon = 0.008856
)

func toXYZ(rgb color.RGBA) XYZ {
	valR, valG, valB := float64(rgb.R)/255.0, float64(rgb.G)/255.0, float64(rgb.G)/255.0
	var newXYZ XYZ
	for _, v := range xyzCoefficient {
		newXYZ.X += v[0] * float64(valR)
		newXYZ.Y += v[1] * float64(valG)
		newXYZ.Z += v[2] * float64(valB)
		fmt.Println(v[2] * float64(valB))
	}

	return newXYZ
}

func toLab(xyz XYZ) LAB {
	var labVal LAB
	// linearize xyz with d65 as base
	xyz.X /= xyzD5Linear[0]
	xyz.Y /= xyzD5Linear[1]
	xyz.Z /= xyzD5Linear[2]

	// clamping
	LabTreshold := math.Pow(xyz.Y, 1/3)
	if LabTreshold > epsilon {
		labVal.L = 116*xyz.Y - 16
	} else {
		labVal.L = kappa * xyz.Y
	}
	if xyz.X > epsilon {
		xyz.X = math.Pow(xyz.X, 1/3)
	} else {
		xyz.X = (xyz.X*kappa + 16) / 116
	}
	if xyz.Y > epsilon {
		xyz.Y = math.Pow(xyz.Y, 1/3)
	} else {
		xyz.Y = (xyz.Y*kappa + 16) / 116
	}
	if xyz.Z > epsilon {
		xyz.Z = math.Pow(xyz.Z, 1/3)
	} else {
		xyz.Z = (xyz.Z*kappa + 16) / 116
	}

	labVal.a = 500 * (xyz.X - xyz.Y)
	labVal.b = 200 * (xyz.Y - xyz.Z)
	return labVal
}

func getc2k() {
	xyzVal := toXYZ(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	fmt.Println(xyzVal.X, xyzVal.Y, xyzVal.Z)
	lab := toLab(xyzVal)
	fmt.Println(lab.L, lab.a, lab.b)
}

// Steps on how to get c2k value
// function CIEDE2000(L1, a1, b1, L2, a2, b2):
//     // Step 1: Calculate CIEDE2000 components
//     C1 = sqrt(a1^2 + b1^2)
//     C2 = sqrt(a2^2 + b2^2)
//     C_bar = (C1 + C2) / 2
//
//     G = 0.5 * (1 - sqrt(C_bar^7 / (C_bar^7 + 25^7)))
//
//     a1_prime = (1 + G) * a1
//     a2_prime = (1 + G) * a2
//
//     C1_prime = sqrt(a1_prime^2 + b1^2)
//     C2_prime = sqrt(a2_prime^2 + b2^2)
//
//     h1_prime = atan2(b1, a1_prime)
//     if h1_prime < 0:
//         h1_prime += 2 * pi
//
//     h2_prime = atan2(b2, a2_prime)
//     if h2_prime < 0:
//         h2_prime += 2 * pi
//
//     // Step 2: Calculate ΔL_prime, ΔC_prime, and ΔH_prime
//     ΔL_prime = L2 - L1
//     ΔC_prime = C2_prime - C1_prime
//
//     Δh_prime = h2_prime - h1_prime
//     if C1_prime * C2_prime == 0:
//         Δh_prime = 0
//     else if abs(h1_prime - h2_prime) <= pi:
//         Δh_prime = h2_prime - h1_prime
//     else if h2_prime - h1_prime > pi:
//         Δh_prime = h2_prime - h1_prime - 2 * pi
//     else:
//         Δh_prime = h2_prime - h1_prime + 2 * pi
//
//     ΔH_prime = 2 * sqrt(C1_prime * C2_prime) * sin(Δh_prime / 2)
//
//     // Step 3: Calculate CIEDE2000 color difference
//     L_bar_prime = (L1 + L2) / 2
//     C_bar_prime = (C1_prime + C2_prime) / 2
//
//     h_bar_prime = (h1_prime + h2_prime) / 2
//     if C1_prime * C2_prime == 0:
//         h_bar_prime = h1_prime + h2_prime
//     else if abs(h1_prime - h2_prime) > pi:
//         h_bar_prime = (h1_prime + h2_prime + 2 * pi) / 2
//
//     T = 1 - 0.17 * cos(h_bar_prime - pi / 6) + 0.24 * cos(2 * h_bar_prime) + 0.32 * cos(3 * h_bar_prime + pi / 30) - 0.2 * cos(4 * h_bar_prime - 63 * pi / 180)
//
//     Δθ = pi / 6 * exp(-((h_bar_prime - 275 * pi / 180) / (25 * pi / 180))^2)
//
//     R_C = 2 * sqrt(C_bar_prime^7 / (C_bar_prime^7 + 25^7))
//
//     S_L = 1 + (0.015 * (L_bar_prime - 50)^2) / sqrt(20 + (L_bar_prime - 50)^2)
//     S_C = 1 + 0.045 * C_bar_prime
//     S_H = 1 + 0.015 * C_bar_prime * T
//
//     R_T = -sin(2 * Δθ) * R_C
//
//     ΔE00 = sqrt((ΔL_prime / S_L)^2 + (ΔC_prime / S_C)^2 + (ΔH_prime / S_H)^2 + R_T * (ΔC_prime / S_C) * (ΔH_prime / S_H))
//
//     return ΔE00
