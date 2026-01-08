package test

import (
	"os"
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"golang.org/x/text/encoding/charmap"
)

func TestConvertAns(t *testing.T) {
	// Read the original .ans file in CP437 encoding
	inputBytes, err := os.ReadFile("data/arl-evoke.ans")
	if err != nil {
		t.Fatalf("Failed to read input file: %v", err)
	}
	
	// Decode from CP437 to UTF-8
	decoder := charmap.CodePage437.NewDecoder()
	inputUTF8, err := decoder.Bytes(inputBytes)
	if err != nil {
		t.Fatalf("Failed to decode CP437: %v", err)
	}
	input := string(inputUTF8)
	
	// Strip SAUCE metadata (everything after \x1a)
	if idx := strings.IndexByte(input, 0x1a); idx >= 0 {
		input = input[:idx]
	}

	// Read the expected converted .ansi file
	expectedBytes, err := os.ReadFile("data/arl-evoke.converted.ansi")
	if err != nil {
		t.Fatalf("Failed to read expected output file: %v", err)
	}
	expected := string(expectedBytes)

	// Convert the input
	result := convert.ConvertAns(input)

	// Assert they match
	Assert(expected, result, t)
}
