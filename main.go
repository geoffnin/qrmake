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

func minVersion(data string, mode Mode) (Version, error) {
	table := []struct {
		version               Version
		maxAlphanumericLength int
		maxNumericLength      int
		maxBytesLength        int
	}{
		{Version(1), 25, 41, 17},
		{Version(2), 47, 77, 32},
		{Version(3), 77, 127, 53},
		{Version(4), 114, 187, 73},
		{Version(5), 154, 255, 106},
		{Version(6), 195, 322, 134},
		{Version(7), 224, 370, 154},
		{Version(8), 279, 461, 192},
		{Version(9), 335, 552, 230},
		{Version(10), 359, 652, 271},
		{Version(11), 468, 772, 321},
		{Version(12), 535, 883, 367},
		{Version(13), 619, 1022, 425},
		{Version(14), 667, 1101, 458},
		{Version(15), 758, 1250, 520},
		{Version(16), 854, 1408, 586},
		{Version(17), 938, 1548, 644},
		{Version(18), 1046, 1725, 718},
		{Version(19), 1153, 1903, 792},
		{Version(20), 1249, 2061, 858},
		{Version(21), 1352, 2232, 929},
		{Version(22), 1460, 2409, 1003},
		{Version(23), 1588, 2620, 1091},
		{Version(24), 1704, 2812, 1171},
		{Version(25), 1853, 3057, 1273},
		{Version(26), 1990, 3283, 1367},
		{Version(27), 2132, 3517, 1465},
		{Version(28), 2223, 3669, 1528},
		{Version(29), 2369, 3909, 1628},
		{Version(30), 2520, 4158, 1732},
		{Version(31), 2677, 4417, 1840},
		{Version(32), 2840, 4686, 1952},
		{Version(33), 3009, 4965, 2068},
		{Version(34), 3183, 5253, 2188},
		{Version(35), 3351, 5529, 2303},
		{Version(36), 3537, 5836, 2431},
		{Version(37), 3729, 6153, 2563},
		{Version(38), 3927, 6479, 2699},
		{Version(39), 4087, 6743, 2809},
		{Version(40), 4296, 7089, 2953},
	}

	for _, v := range table {
		if mode == modeAlphanumeric {
			if v.maxAlphanumericLength >= len(data) {
				return v.version, nil
			}
		} else if mode == modeNumeric {
			if v.maxNumericLength >= len(data) {
				return v.version, nil
			}
		} else if mode == modeByte {
			if v.maxBytesLength >= len(data) {
				return v.version, nil
			}
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
