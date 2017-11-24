package main

import (
	"bytes"
	"errors"
	"image/png"
	"math/rand"
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

func generateRandomString(mode string, length int) string {
	var letters = []rune(mode)
	var runes = make([]rune, length)

	for i := range runes {
		runes[i] = letters[rand.Intn(len(letters))]
	}

	return string(runes)
}

func TestMinVersionReturnsMinVersion(t *testing.T) {

	numericMode := "0123456789"
	alphanumericMode := "ABVDEFGHIJKLMNOPQRSTUVWXYZ0123456789-./']+$"
	bytesMode := "abcdefghijklmnopqrstuvwxyzABVDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	table := []struct {
		input    string
		encoding Mode
		excepted Version
	}{
		{"555-2834", modeByte, Version(1)},
		{"555-2834", dataAnalysis("555-2834"), Version(1)},
		{"5552834", dataAnalysis("5552834"), Version(1)},
		{generateRandomString(alphanumericMode, 24), modeAlphanumeric, Version(1)},
		{generateRandomString(alphanumericMode, 30), modeAlphanumeric, Version(2)},
		{generateRandomString(alphanumericMode, 154), modeAlphanumeric, Version(5)},
		{generateRandomString(alphanumericMode, 2200), modeAlphanumeric, Version(28)},
		{generateRandomString(numericMode, 2200), modeNumeric, Version(21)},
		{generateRandomString(bytesMode, 2200), modeByte, Version(35)},
	}

	for _, test := range table {
		result, _ := minVersion(test.input, test.encoding)
		if result != test.excepted {
			t.Errorf("Expected Version(%d) but got Version(%d) for input of len(%d)", test.excepted, result, len(test.input))
		}
	}

	_, err := minVersion(generateRandomString(alphanumericMode, 10000), modeAlphanumeric)

	if err == nil || err.Error() != "input too large for QR Code" {
		t.Errorf("Expected size exception, got %v", err)
	}

}

func equalSlice(a []interface{}, b []interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestDataAnalysisReturnsTheCorrectMode(t *testing.T) {
	table := []struct {
		input    string
		expected Mode
	}{
		{"12345", modeNumeric},
		{"%*ABC1 ./3123$+-:", modeAlphanumeric},
		{"abcABC, 1234", modeByte},
		{"", modeByte},
		{"Привет мир", modeByte},
	}

	for _, test := range table {
		result := dataAnalysis(test.input)

		if result != test.expected {
			t.Errorf("DataAnalysis(%v) expected %v but got %v", test.input, test.expected, result)
		}
	}
}
