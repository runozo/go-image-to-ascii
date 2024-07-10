# Go-image-to-ascii

Convert a bitmap image into ASCII.

## Features

- works with webcam in realtime
- reads png, jpeg and webp images

![source image](examples/image.webp)

![resulting image](examples/image_ascii.png)

## Get the binary
You can obtain the binary executable for your operating system from the [releases page](https://github.com/runozo/go-image-to-ascii/releases)

## Build

```go build -ldflags="-s -w -v" ./...```

## Usage

The source image or webcam frames, will be stretched to fit the actual size in characters, of the active terminal.

### Webcam

```./go-image-to-ascii -webcam```

Press ```CTRL+C``` or ```ESC``` to exit.

### Static image

```./go-image-to-ascii -file image.webp```

## TODO

[ ] Add key bindings to change character density

