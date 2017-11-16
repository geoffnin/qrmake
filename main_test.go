package main

import (
	"bytes"
	"errors"
	"image/png"
	"testing"
)

func TestGenerateQRCodeReturnsPNG(t *testing.T) {
	result := new(bytes.Buffer)
	GenerateQRCode(result, "555-2368", Version(1))

	if result.Len() == 0 {
		t.Errorf("Generated QRCode has no data")
	}

	_, err := png.Decode(result)

	if err != nil {
		t.Errorf("Generated QR Code is not a valid png: %s", err)
	}
}

func TestVersionDeterminesSize(t *testing.T) {
	table := []struct {
		version uint8
		size    int
	}{
		{1, 21},
		{2, 25},
		{6, 41},
		{7, 45},
		{14, 73},
		{40, 177},
	}

	for _, test := range table {
		version := Version(test.version)
		if version.PatternSize() != test.size {
			t.Errorf("Version %d, expected %d but got %d", test.version, test.size, version.PatternSize())
		}
	}
}

type errorWriter struct{}

func (w *errorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("expected error")
}

func TestGenerateQRCodePropgatesError(t *testing.T) {
	w := new(errorWriter)
	err := GenerateQRCode(w, "555-2834", Version(1))

	if err == nil || err.Error() != "expected error" {
		t.Errorf("GeneratedQRCode does not propogate correctly, got %v", err)
	}
}
