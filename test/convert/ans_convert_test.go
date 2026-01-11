package test

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/text/encoding/charmap"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
)

func TestConvertAnsFiles(t *testing.T) {
	testCases := []struct {
		name          string
		inputFpath    string
		expectedFpath string
	}{
		{
			// https://16colo.rs/pack/impure91/arl-evoke.ans
			"arl-evoke",
			"../data/arl-evoke.ans",
			"../data/arl-evoke.converted.ansi",
		},
		{
			// https://16colo.rs/pack/impure89/xz-gibson.ans
			"xz-gibson",
			"../data/xz-gibson.ans",
			"../data/xz-gibson.converted.ansi",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read the original .ans file in CP437 encoding

			inputBytes, err := os.ReadFile(tc.inputFpath)
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
			expectedBytes, err := os.ReadFile(tc.expectedFpath)
			if err != nil {
				t.Fatalf("Failed to read expected output file: %v", err)
			}
			expected := string(expectedBytes)

			// Convert the input
			sauce, _, err := convert.ParseSAUCEFromFile(tc.inputFpath)
			if err != nil {
				t.Fatalf("Failed to parse SAUCE record: %v", err)
			}
			result := convert.ConvertAns(input, *sauce)
			// Assert they match
			test.Assert(expected, result, t)
		})
	}
}
