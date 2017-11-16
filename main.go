package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
)

// GenerateQRCode generates a qr code
func GenerateQRCode(w io.Writer, code string, version Version) error {
	size := version.PatternSize()
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	return png.Encode(w, img)
}

// Version is an unsigned short int with a value 1-4
type Version uint8

// PatternSize returns the width of the pattern for this version in modules
func (v Version) PatternSize() int {
	return 4*int(v) + 17
}

func main() {
	fmt.Println("Hello QR")

	file, err := os.Open("qrcode.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = GenerateQRCode(file, "555-2368", Version(1))
	if err != nil {
		log.Fatal(err)
	}

}
