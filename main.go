package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/davidbyttow/govips/v2/vips"
)

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

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()

	img, err := vips.NewImageFromFile("./Designer.jpeg")
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}
	imgByte, err := img.ToBytes()
	if err != nil {
		log.Fatalf("skill issue at , :%s\n", err)
	}

	heightImg := img.Height()
	widthImg := img.Width()
	fmt.Printf("Image size width:%d height:%d\n", widthImg, heightImg)
	fmt.Println(len(imgByte))
	// chunking
	// so i know the part of image that are being processed, so i know there image parts are and it's order
	// so can i just use it to order it
	// if i know the area are 512 - 1024 that's mean it's 2/2 part of the whole array, where as if it
	chunkAmount := 3
	chunkHeight := heightImg / chunkAmount
	chunkWidth := widthImg / chunkAmount
	chunkHeightChange := heightImg % chunkAmount
	chunkWidthChange := widthImg % chunkAmount
	chunkHeightArr := make([][]int, 0)
	chunkWidthArr := make([][]int, 0)
	for i := 0; i < heightImg; i += chunkHeight {
		if i+chunkHeight+chunkHeightChange == heightImg {
			chunkHeightArr = append(chunkHeightArr, []int{i, i + chunkHeight + chunkHeightChange})
			i += chunkHeightChange
		} else {
			chunkHeightArr = append(chunkHeightArr, []int{i, i + chunkHeight})
		}
	}
	for i := 0; i < widthImg; i += chunkWidth {
		if i+chunkWidth+chunkWidthChange == widthImg {
			chunkWidthArr = append(chunkWidthArr, []int{i, i + chunkWidth + chunkWidthChange})
			i += chunkWidthChange
		} else {
			chunkWidthArr = append(chunkWidthArr, []int{i, i + chunkWidth})
		}
	}
	fmt.Println(chunkHeightArr)
	fmt.Println(chunkWidthArr)

	coordArr := make([][]int, 0)
	for h := 0; h < len(chunkHeightArr); h++ {
		for w := 0; w < len(chunkWidthArr); w++ {
			fmt.Println(chunkHeightArr[h], chunkWidthArr[w])
			coord := make([]int, 0)
			coord = append(coord, chunkHeightArr[h]...)
			coord = append(coord, chunkWidthArr[w]...)
			coordArr = append(coordArr, coord)
		}
	}

	wg := sync.WaitGroup{}
	startTime := time.Now()
	for _, v := range coordArr {
		wg.Add(1)
		go getRgb(v[0], v[1], v[2], v[3], &wg, img)
	}
	fmt.Println(coordArr)
	wg.Wait()
	fmt.Printf("Took %s s\n", time.Now().Sub(startTime).String())
}
