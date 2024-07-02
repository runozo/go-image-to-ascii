# Go-image-to-ascii

Convert a bitmap image into ASCII.

![source image](examples/image.webp)

![resulting image](examples/image_ascii.png)

## Build

```go build -o go-image-to-ascii main.go```

## Build for other platforms

```make windows```

```make linux```

```make darwin```

## Usage

The source image will be stretched to fit the actual size in characters of the terminal. 

```./go-image-to-ascii image.webp```

## Help

```./go-image-to-ascii -h```

## Todo

[] fetch image from webcam
[] live video rendering on terminal



