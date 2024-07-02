package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nsf/termbox-go"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/prop"

	// This is required to register camera adapter
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	// This is required to register camera adapter
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	// Note: If you don't have a camera or your adapters are not supported,
	//       you can always swap your adapters with our dummy adapters below.
	// _ "github.com/pion/mediadevices/pkg/driver/videotest"
)

// const density = "Ñ@#W$9876543210?!abc;:+=-,._          "
const density = "         _.,-=+:;cba!?0123456789$W#@Ñ"
const coldef = termbox.ColorDefault

type Pixel struct {
	R int
	G int
	B int
	A int
}

var webcam = flag.Bool("webcam", false, "capture image from webcam")
var filename = flag.String("file", "", "read image from file")

func main() {
	// You can register other formats here
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jfif", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("webp", "rif", webp.Decode, webp.DecodeConfig)

	flag.Parse()

	if *filename == "" && !*webcam {
		fmt.Println("Usage:", os.Args[0], "-file <filename>|-webcam")
		os.Exit(1)
	}

	if err := termbox.Init(); err != nil {
		panic(err)
	}

	termWidth, termHeight := termbox.Size()
	defer termbox.Close()
	fmt.Println("width:", termWidth, "height:", termHeight)

	if *filename != "" {
		file, err := os.Open(*filename)

		if err != nil {
			fmt.Println("Error: File could not be opened")
			os.Exit(1)
		}

		defer file.Close()

		srcImage, _, err := image.Decode(file)

		if err != nil {
			fmt.Println("Error: Image could not be decoded")
			os.Exit(1)
		}

		pixels := getPixels(srcImage, termWidth, termHeight)

		fmt.Println(pixelsToAscii(termWidth, termHeight, pixels, density))
	}

	if *webcam {
		stream, _ := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
			Video: func(constraint *mediadevices.MediaTrackConstraints) {
				// Query for ideal resolutions
				constraint.Width = prop.Int(600)
				constraint.Height = prop.Int(400)
			},
		})

		// Since track can represent audio as well, we need to cast it to
		// *mediadevices.VideoTrack to get video specific functionalities
		track := stream.GetVideoTracks()[0]
		videoTrack := track.(*mediadevices.VideoTrack)
		defer videoTrack.Close()

		// Create a new video reader to get the decoded frames. Release is used
		// to return the buffer to hold frame back to the source so that the buffer
		// can be reused for the next frames.
		videoReader := videoTrack.NewReader(false)
		termbox.Clear(coldef, coldef)
		for {
			frame, release, _ := videoReader.Read()
			pixels := getPixels(frame, termWidth, termHeight)
			asciiPixels := pixelsToAscii(termWidth, termHeight, pixels, density)
			fmt.Println(asciiPixels)
			termbox.Clear(coldef, coldef)
			/*
				for y := 0; y < termHeight; y++ {
					for x := 0; x < termWidth; x++ {
						index := y*termWidth + x
						termbox.SetChar(x, y, rune(asciiPixels[index]))
					}
				}
				termbox.Flush()
			*/
			release()
		}

		// Since frame is the standard image.Image, it's compatible with Go standard
		// library. For example, capturing the first frame and store it as a jpeg image.
		// output, _ := os.Create("frame.jpg")
		// jpeg.Encode(output, frame, nil)
	}
}

func getPixels(srcImage image.Image, termWidth, termHeight int) []Pixel {

	// Set the image size to the terminal size
	resizedImage := image.NewRGBA(image.Rect(0, 0, termWidth, termHeight))

	// Resize
	draw.NearestNeighbor.Scale(resizedImage, resizedImage.Rect, srcImage, srcImage.Bounds(), draw.Over, nil)

	bounds := resizedImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make([]Pixel, 0, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels = append(pixels, rgbaToPixel(resizedImage.At(x, y).RGBA()))
		}
	}

	return pixels
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
