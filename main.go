package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
)

type RGB struct {
	R uint8
	G uint8
	B uint8
}

// type Filterer interface {
// 	Filter(*vips.ImageRef)
// }
//
// type ImageFilter struct {
// 	saveDir     string
// 	filePath    string
// 	newFileName string
// 	config
// }

type mode string

const (
	modeProfiling mode = "profiling"
	modeC2k       mode = "c2k"
	modeImglib    mode = "imglib"
	modeVips      mode = "vips"
)

var modePtr = flag.String("mode", "", "Select the program mode (profiling, c2k, imglib, vips)")

func main() {
	flag.Parse()

	var selectedMode mode
	switch *modePtr {
	case string(modeProfiling):
		selectedMode = modeProfiling
	case string(modeC2k):
		selectedMode = modeC2k
	case string(modeImglib):
		selectedMode = modeImglib
	case string(modeVips):
		selectedMode = modeVips
	default:
		fmt.Println("Invalid mode specified. Please use profiling, c2k, imglib, or vips.")
		return
	}

	fmt.Println("Selected mode:", selectedMode)
	switch selectedMode {
	case modeC2k:
		getc2k()
	case modeImglib:
		runImgLib()
	case modeVips:
		// Perform actions specific to vips mode
		fmt.Println("vips mode activated.")
		// ... your vips logic here
	}

	// Enable CPU profiling
	if selectedMode == modeProfiling {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
}

func runImgLib() {
	vips.Startup(nil)
	defer vips.Shutdown()

	palette := []string{
		"#2e3440", "#3b4252", "#434c5e", "#4c566a",
		"#d8dee9", "#e5e9f0", "#eceff4", "#8fbcbb",
		"#88c0d0", "#81a1c1", "#5e81ac", "#bf616a",
		"#d08770", "#ebcb8b", "#a3be8c", "#b48ead",
	}

	startTime := time.Now()
	img, err := vips.NewImageFromFile("./4k.jpg")
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}
	imgProcess := ProcessStdLib{newFileName: "4k_conv", filePath: "./sample/4k.jpg", saveDir: "./output/"}
	imgProcess.Process(img, palette)
	durTime := time.Since(startTime)
	fmt.Printf("Process took %s s\n", durTime.String())
}
