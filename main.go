package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
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

func inputAnalysis(data string) (Version, error) {
	table := []struct {
		version   Version
		maxLength int
	}{
		{Version(1), 25},
		{Version(2), 47},
		{Version(3), 77},
		{Version(4), 114},
		{Version(5), 154},
		{Version(6), 195},
		{Version(7), 224},
		{Version(8), 279},
		{Version(9), 335},
		{Version(10), 359},
		{Version(11), 468},
		{Version(12), 535},
		{Version(13), 619},
		{Version(14), 667},
		{Version(15), 758},
		{Version(16), 854},
		{Version(17), 938},
		{Version(18), 1046},
		{Version(19), 1153},
		{Version(20), 1249},
		{Version(21), 1352},
		{Version(22), 1460},
		{Version(23), 1588},
		{Version(24), 1704},
		{Version(25), 1853},
		{Version(26), 1990},
		{Version(27), 2132},
		{Version(28), 2223},
		{Version(29), 2369},
		{Version(30), 2520},
		{Version(31), 2677},
		{Version(32), 2840},
		{Version(33), 3009},
		{Version(34), 3183},
		{Version(35), 3351},
		{Version(36), 3537},
		{Version(37), 3729},
		{Version(38), 3927},
		{Version(39), 4087},
		{Version(40), 4296},
	}

	for _, v := range table {
		if v.maxLength >= len(data) {
			return v.version, nil
		}
	}

	return Version(0), errors.New("input too large for QR Code")
}

type Mode int

const (
	modeNumeric Mode = iota
	modeAlphanumeric
	modeByte
)

func dataAnalysis(input string) Mode {
	// Check if input is just a number
	if _, err := strconv.Atoi(input); err == nil {
		return modeNumeric
	}

	// Check if input is in limited set of alphanumeric latin characters
	if regexp.MustCompile(`^[ 0-9A-Z$%*+-./:]+$`).MatchString(input) {
		return modeAlphanumeric
	}

	// Check if input is valid ISO 8859-1 string
	// No need to check if the default is byte
	// if regexp.MustCompile(`^[a-zA-Z0-9\xc0-\xd6 ,]*$`).MatchString(input) {
	// 	return modeByte
	// }

	return modeByte
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
