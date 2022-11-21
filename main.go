package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	// Fetch Image Response
	res, err := http.Get("https://media.discordapp.net/attachments/1001564196015702026/1039702135467671582/unknown.png?width=400&height=400")
	if err != nil || res.StatusCode != 200 {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// Decode Image
	decodedImage, _, err := image.Decode(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Make Image RGBA Image (to allow for transparency)
	img := image.NewRGBA(decodedImage.Bounds())
	size := img.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			img.Set(x, y, decodedImage.At(x, y))
		}
	}

	// Remove the background
	RemoveBackground(img)

	duration := time.Since(start)
	fmt.Println(duration)

	outFile, _ := os.Create("output/changed.png")
	defer outFile.Close()
	png.Encode(outFile, img)
}

func RemoveBackground(image *image.RGBA) {
	maxX := image.Bounds().Max.X
	maxY := image.Bounds().Max.Y
	targetColor := image.At(0, 0)

	FillWithTargetColor(targetColor, 0, 0, color.RGBA{0, 0, 0, 0}, image)

	// if the other corner is already the target color, we dont run on the other corners
	if image.At(maxX-1, maxY-1) != targetColor {
		FillWithTargetColor(targetColor, maxX-1, maxY-1, color.RGBA{0, 0, 0, 0}, image)
		FillWithTargetColor(targetColor, 0, maxY-1, color.RGBA{0, 0, 0, 0}, image)
		FillWithTargetColor(targetColor, maxX-1, 0, color.RGBA{0, 0, 0, 0}, image)
	}
}

func FillWithTargetColor(targetColor color.Color, x int, y int, newColor color.Color, image *image.RGBA) {

	rx, gx, bx, ax := image.At(x, y).RGBA()

	tr, tg, tb, ta := targetColor.RGBA()

	a := math.Abs(float64(rx)-float64(tr)) + math.Abs(float64(gx)-float64(tg)) + math.Abs(float64(bx)-float64(tb)) + math.Abs(float64(ax)-float64(ta))

	if a < 30000 {
		image.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
		if x+1 < image.Bounds().Max.X {
			FillWithTargetColor(targetColor, x+1, y, newColor, image)
		}
		if x-1 > 0 {
			FillWithTargetColor(targetColor, x-1, y, newColor, image)
		}
		if y+1 < image.Bounds().Max.Y {
			FillWithTargetColor(targetColor, x, y+1, newColor, image)
		}
		if y-1 > 0 {
			FillWithTargetColor(targetColor, x, y-1, newColor, image)
		}
	}

}
