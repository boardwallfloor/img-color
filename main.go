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
	img, err := vips.NewImageFromFile("./4k.jpg")
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}
	stdProcess(img)
	durTime := time.Since(startTime)
	fmt.Printf("Process took %s s\n", durTime.String())
}
