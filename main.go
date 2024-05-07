package main

import (
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

func main() {
	// Enable CPU profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	vips.Startup(nil)
	defer vips.Shutdown()

	startTime := time.Now()
	img, err := vips.NewImageFromFile("./sample/4k.jpg")
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}
	palette := []string{
		"#2e3440", "#3b4252", "#434c5e", "#4c566a",
		"#d8dee9", "#e5e9f0", "#eceff4", "#8fbcbb",
		"#88c0d0", "#81a1c1", "#5e81ac", "#bf616a",
		"#d08770", "#ebcb8b", "#a3be8c", "#b48ead",
	}
	stdImg := ProcessStdLib{saveDir: "./output/", filePath: "./sample/4k.jpg", newFileName: "4k_conv"}
	stdImg.Process(img, palette)

	durTime := time.Since(startTime)
	fmt.Printf("Process took %s s\n", durTime.String())
}
