package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/gdamore/tcell"
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

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	screen.SetStyle(tcell.StyleDefault)
	// defer screen.Fini()

	termWidth, termHeight := screen.Size()
	screen.Fini()
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

		fmt.Println(imageToAscii(srcImage, termWidth, termHeight, density))
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
		// mainloop:
		for {
			termWidth, termHeight = screen.Size()
			// read a frame from the webcam
			frame, release, err := videoReader.Read()
			if err != nil {
				fmt.Println("Error: Video could not be read")
				os.Exit(1)
			}

			flushImageToScreen(screen, frame, termWidth, termHeight, density)

			release()

			// poll for keyboard events in another goroutine
			go func() {
				for {
					switch event := screen.PollEvent().(type) {
					case *tcell.EventKey:
						if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
							screen.Fini()
							os.Exit(0)
						}
					}
				}
			}()
		}
	}
}

// imageToAscii converts the source image into ASCII art based on the given terminal width, terminal height, and density values.
//
// Parameters:
//   - srcImage: the source image to convert
//   - termWidth: the width of the terminal
//   - termHeight: the height of the terminal
//   - density: the density string used to map pixel values to ASCII characters
//
// Return type: string
func imageToAscii(srcImage image.Image, termWidth, termHeight int, density string) string {
	var buffer bytes.Buffer
	slope := float32((len(density) - 1)) / 255.0
	// Set the image size to the terminal size
	resizedImage := image.NewRGBA(image.Rect(0, 0, termWidth, termHeight))

	// Resize
	draw.NearestNeighbor.Scale(resizedImage, resizedImage.Bounds(), srcImage, srcImage.Bounds(), draw.Over, nil)

	for y := 0; y < termHeight; y++ {
		for x := 0; x < termWidth; x++ {
			// get rgba values
			pixel := rgbaToPixel(resizedImage.At(x, y).RGBA())
			avg := (pixel.R + pixel.G + pixel.B) / 3 // 0-255
			buffer.WriteString(string(density[int32(float32(avg)*slope)]))
		}
	}
	return buffer.String()
}

// flushImageToScreen flushes the given image to the screen.
//
// Parameters:
// - screen: the tcell screen to draw on.
// - frame: the image to draw.
// - termWidth: the width of the terminal.
// - termHeight: the height of the terminal.
// - density: the density string used to map pixel values to ASCII characters.
func flushImageToScreen(screen tcell.Screen, frame image.Image, termWidth, termHeight int, density string) {
	screen.Fill(' ', tcell.StyleDefault)
	asciiPixels := imageToAscii(frame, termWidth, termHeight, density)

	for y := 0; y < termHeight; y++ {
		for x := 0; x < termWidth; x++ {
			screen.SetCell(x, y, tcell.StyleDefault, rune(asciiPixels[y*termWidth+x]))
		}
	}
	screen.Show()
}

// rgbaToPixel converts RGBA values to a Pixel struct.
//
// Parameters:
// - r: the red component of the RGBA value.
// - g: the green component of the RGBA value.
// - b: the blue component of the RGBA value.
// - a: the alpha component of the RGBA value.
// Return type: Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}
