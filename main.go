package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/nsf/termbox-go"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

// const density = "Ñ@#W$9876543210?!abc;:+=-,._          "
const density = "       _.,-=+:;cba!?0123456789$W#@Ñ"

type Pixel struct {
	R int
	G int
	B int
	A int
}

// var filename = flag.String("image", "image.webp", "filename of image to convert to ascii (png, jpeg, webp)")

func main() {
	// You can register other formats here
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jfif", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("webp", "rif", webp.Decode, webp.DecodeConfig)

	flag.Parse()

	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "filename")
		os.Exit(1)
	}

	filename := os.Args[1]

	if err := termbox.Init(); err != nil {
		panic(err)
	}

	termWidth, termHeight := termbox.Size()
	termbox.Close()
	fmt.Println("width:", termWidth, "height:", termHeight)

	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	defer file.Close()

	pixels, err := getPixels(file, termWidth, termHeight)

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	fmt.Println(pixelsToAscii(termWidth, termHeight, pixels, density))
}

func getPixels(file io.Reader, termWidth, termHeight int) ([]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	// Set the image size to the terminal size
	resizedImage := image.NewRGBA(image.Rect(0, 0, termWidth, termHeight))

	// Resize
	draw.NearestNeighbor.Scale(resizedImage, resizedImage.Rect, img, img.Bounds(), draw.Over, nil)

	bounds := resizedImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make([]Pixel, 0, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels = append(pixels, rgbaToPixel(resizedImage.At(x, y).RGBA()))
		}
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func pixelsToAscii(width int, height int, pixels []Pixel, density string) string {
	var buffer bytes.Buffer
	slope := float32((len(density) - 1)) / 255.0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := y*width + x
			avg := (pixels[index].R + pixels[index].G + pixels[index].B) / 3 // 0-255
			// fmt.Println(slope)
			buffer.WriteString(string(density[int32(float32(avg)*slope)]))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}
