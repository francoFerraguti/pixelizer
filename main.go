package main

import (
	"flag"
	"github.com/liteByte/frango"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

func main() {
	fileNameFlag := flag.String("filename", "sample.png", "")
	pixelWidthFlag := flag.Int("pixelWidth", 2, "")
	pixelHeightFlag := flag.Int("pixelHeight", 2, "")
	flag.Parse()

	frango.SeedRandom()

	createSampleImage("sample.png", 100, 100)

	sampleImage := loadImage(*fileNameFlag)

	pixelizedImage := pixelizeImage(sampleImage, *pixelWidthFlag, *pixelHeightFlag)

	outputName := "output/" + frango.GetRandomString(8) + "_" + strconv.Itoa(*pixelWidthFlag) + "_" + strconv.Itoa(*pixelHeightFlag) + ".png"

	outputFile, err := os.Create(outputName)
	defer outputFile.Close()

	if err != nil {
		log.Println(err.Error())
	}

	png.Encode(outputFile, pixelizedImage)
}

func pixelizeImage(img image.Image, pixelWidth, pixelHeight int) image.Image {

	width, height, widthInFakePixels, heightInFakePixels := getImageData(img, pixelWidth, pixelHeight)

	log.Println("Original image dimensions: " + strconv.Itoa(width) + "x" + strconv.Itoa(height))
	log.Println("Width in fake pixels: " + strconv.Itoa(widthInFakePixels))
	log.Println("Height in fake pixels: " + strconv.Itoa(heightInFakePixels))

	rArray, gArray, bArray, aArray := createRGBAArrays(widthInFakePixels, heightInFakePixels) //These are 4 arrays with the size of [widthInFaxePixels, heightInFakePixels]

	for y := 0; y < height; y++ { //This double loop adds the color of every pixel to the portion of the fakePixelArray that we created before
		for x := 0; x < width; x++ {

			fakePixelX := 0
			fakePixelY := 0
			c := 0

			r, g, b, a := img.At(x, y).RGBA() //Gets the color of the current pixel

			for i := 0; i <= width; i += pixelWidth {
				c++

				if x < i {
					fakePixelX = c - 2 //Assigns the correct fakePixel position to the pixel
					break
				}
			}

			c = 0

			for i := 0; i <= height; i += pixelHeight {
				c++

				if y < i {
					fakePixelY = c - 2 //Assigns the correct fakePixel position to the pixel
					break
				}
			}

			rArray[fakePixelX][fakePixelY] += int(r / 0x101) //Adds the color
			gArray[fakePixelX][fakePixelY] += int(g / 0x101) //Adds the color
			bArray[fakePixelX][fakePixelY] += int(b / 0x101) //Adds the color
			aArray[fakePixelX][fakePixelY] += int(a / 0x101) //Adds the color
		}
	}

	for y := 0; y < heightInFakePixels; y++ { //This double loop gets the average of the colors we just added
		for x := 0; x < widthInFakePixels; x++ {
			rArray[x][y] = rArray[x][y] / pixelWidth / pixelHeight
			gArray[x][y] = gArray[x][y] / pixelWidth / pixelHeight
			bArray[x][y] = bArray[x][y] / pixelWidth / pixelHeight
			aArray[x][y] = aArray[x][y] / pixelWidth / pixelHeight
		}
	}

	log.Println("Finished creating the fakePixel arrays")

	finalImage := image.NewRGBA(image.Rect(0, 0, width, height)) //Creates an empty image to paste the new one above

	for y := 0; y < height; y++ { //It's the same double loop as the one above, but this time it pastes the fakePixel color onto a new image
		for x := 0; x < width; x++ {

			fakePixelX := 0
			fakePixelY := 0
			c := 0

			for i := 0; i <= width; i += pixelWidth {
				c++

				if x < i {
					fakePixelX = c - 2
					break
				}
			}

			c = 0

			for i := 0; i <= height; i += pixelHeight {
				c++

				if y < i {
					fakePixelY = c - 2
					break
				}
			}

			pixelColor := prepareColor(rArray[fakePixelX][fakePixelY], gArray[fakePixelX][fakePixelY], bArray[fakePixelX][fakePixelY], aArray[fakePixelX][fakePixelY])
			finalImage.SetRGBA(x, y, pixelColor)
		}
	}

	log.Println("Finished pixelizing image")

	return finalImage
}

func getImageData(img image.Image, pixelWidth, pixelHeight int) (int, int, int, int) {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	for width%pixelWidth != 0 {
		width--
	}

	for height%pixelHeight != 0 {
		height--
	}

	widthInFakePixels := width / pixelWidth
	heightInFakePixels := height / pixelHeight

	return width, height, widthInFakePixels, heightInFakePixels
}

func convertImageToRGBA(img image.Image) image.RGBA {
	bounds := img.Bounds()
	imageRGBA := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(imageRGBA, imageRGBA.Bounds(), img, bounds.Min, draw.Src)
	return *imageRGBA
}

func createRGBAArrays(width, height int) ([][]int, [][]int, [][]int, [][]int) {
	rArray := make([][]int, width)
	for i := range rArray {
		rArray[i] = make([]int, height)
	}

	gArray := make([][]int, width)
	for i := range gArray {
		gArray[i] = make([]int, height)
	}

	bArray := make([][]int, width)
	for i := range bArray {
		bArray[i] = make([]int, height)
	}

	aArray := make([][]int, width)
	for i := range aArray {
		aArray[i] = make([]int, height)
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
