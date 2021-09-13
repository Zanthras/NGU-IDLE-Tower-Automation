package main

import "C"
import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/go-vgo/robotgo"
)

func makePNG() {

	w := 20
	h := 30

	upLeft := image.Point{0, 0}
	lowRight := image.Point{w, h}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			pxColor := robotgo.GetPixelColor(x+LEFT+EnemyHealth.X, y+TOP+EnemyHealth.Y)
			r, _ := strconv.ParseInt(pxColor[:2], 16, 32)
			g, _ := strconv.ParseInt(pxColor[2:4], 16, 32)
			b, _ := strconv.ParseInt(pxColor[4:], 16, 32)
			pixel := color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
			img.Set(x, y, pixel)
			log.Println("Pixilized", x+EnemyHealth.X, y+EnemyHealth.Y, x, y)
		}
	}
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

func snagRect(rect RECT, filename string) {
	rect.Top = rect.Top + int32(TOP)
	rect.Bottom = rect.Bottom + int32(TOP)
	rect.Left = rect.Left + int32(LEFT)
	rect.Right = rect.Right + int32(LEFT)
	bitmap := robotgo.CaptureScreen(int(rect.Left), int(rect.Top), int(rect.Right-rect.Left), int(rect.Bottom-rect.Top))
	defer robotgo.FreeBitmap(bitmap)
	robotgo.SaveBitmap(bitmap, filename)
}

func getRectColors(rect RECT) map[string]int {
	rect.Top = rect.Top + int32(TOP)
	rect.Bottom = rect.Bottom + int32(TOP)
	rect.Left = rect.Left + int32(LEFT)
	rect.Right = rect.Right + int32(LEFT)
	bitmap := robotgo.CaptureScreen(int(rect.Left), int(rect.Top), int(rect.Right-rect.Left), int(rect.Bottom-rect.Top))
	defer robotgo.FreeBitmap(bitmap)
	h := int(rect.Bottom) - int(rect.Top)
	w := int(rect.Right) - int(rect.Left)
	mapping := make(map[string]int)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			pxColor := robotgo.GetColors(bitmap, x, y)
			//log.Println(pxColor)
			if _, found := mapping[pxColor]; !found {
				mapping[pxColor] = 0
			}
			mapping[pxColor]++
		}
	}
	return mapping
}
