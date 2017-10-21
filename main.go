package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

func main() {
	createSampleImage("sample.png", 100, 100)

	sampleImage := loadImage("gioconda_2.jpg")

	pixelizedImage := pixelizeImage(sampleImage, 4, 6)

	outputFile, err := os.Create("output.png")
	defer outputFile.Close()

	if err != nil {
		log.Println(err.Error())
	}

	png.Encode(outputFile, pixelizedImage)
}

func pixelizeImage(img image.Image, pixelWidth, pixelHeight int) image.Image {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	widthInFakePixels := width / pixelWidth
	heightInFakePixels := height / pixelHeight

	log.Println("Width in fake pixels: " + strconv.Itoa(widthInFakePixels))
	log.Println("Height in fake pixels: " + strconv.Itoa(heightInFakePixels))

	rArray, gArray, bArray, aArray := createRGBAArrays(widthInFakePixels, heightInFakePixels)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			fakePixelX := 0
			fakePixelY := 0

			r, g, b, a := img.At(x, y).RGBA()

			for i := 0; i < width; i += pixelWidth {
				if x < i {
					fakePixelX = i / pixelWidth
					break
				}
			}

			for i := 0; i < height; i += pixelHeight {
				if y < i {
					fakePixelY = i / pixelHeight
					break
				}
			}

			rArray[fakePixelX][fakePixelY] += int(r / 0x101)
			gArray[fakePixelX][fakePixelY] += int(g / 0x101)
			bArray[fakePixelX][fakePixelY] += int(b / 0x101)
			aArray[fakePixelX][fakePixelY] += int(a / 0x101)
		}
	}

	for y := 0; y < heightInFakePixels; y++ {
		for x := 0; x < widthInFakePixels; x++ {
			rArray[x][y] = rArray[x][y] / pixelWidth / pixelHeight
			gArray[x][y] = gArray[x][y] / pixelWidth / pixelHeight
			bArray[x][y] = bArray[x][y] / pixelWidth / pixelHeight
			aArray[x][y] = aArray[x][y] / pixelWidth / pixelHeight
		}
	}

	finalImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			fakePixelX := 0
			fakePixelY := 0

			for i := 0; i < width; i += pixelWidth {
				if x < i {
					fakePixelX = i / pixelWidth
					break
				}
			}

			for i := 0; i < height; i += pixelHeight {
				if y < i {
					fakePixelY = i / pixelHeight
					break
				}
			}

			pixelColor := prepareColor(rArray[fakePixelX][fakePixelY], gArray[fakePixelX][fakePixelY], bArray[fakePixelX][fakePixelY], aArray[fakePixelX][fakePixelY])
			finalImage.SetRGBA(x, y, pixelColor)
		}
	}

	return finalImage
}

func createRGBAArrays(width, height int) ([][]int, [][]int, [][]int, [][]int) {
	rArray := make([][]int, height)
	for i := range rArray {
		rArray[i] = make([]int, width)
	}

	gArray := make([][]int, height)
	for i := range gArray {
		gArray[i] = make([]int, width)
	}

	bArray := make([][]int, height)
	for i := range bArray {
		bArray[i] = make([]int, width)
	}

	aArray := make([][]int, height)
	for i := range aArray {
		aArray[i] = make([]int, width)
	}

	return rArray, gArray, bArray, aArray
}

func createSampleImage(filename string, width, height int) {
	myImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixelColor := prepareColor(x+y, x+y, x+y, 255)
			myImage.SetRGBA(x, y, pixelColor)
		}
	}

	outputFile, err := os.Create(filename)
	defer outputFile.Close()

	if err != nil {
		log.Println(err.Error())
	}

	png.Encode(outputFile, myImage)
}

func loadImage(filename string) image.Image {
	fImg1, _ := os.Open(filename)
	defer fImg1.Close()
	img, _, _ := image.Decode(fImg1)
	return img
}

func prepareColor(r, g, b, a int) color.RGBA {
	if r < 0 {
		r = 0
	}
	if g < 0 {
		g = 0
	}
	if b < 0 {
		b = 0
	}
	if a < 0 {
		a = 0
	}
	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}
	if a > 255 {
		a = 255
	}

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}
