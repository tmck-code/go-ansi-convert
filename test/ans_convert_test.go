package test

import (
	"os"
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
	"golang.org/x/text/encoding/charmap"
)

func TestConvertAnsFiles(t *testing.T) {
	testCases := []struct {
		inputFpath    string
		expectedFpath string
	}{
		{
			"data/arl-evoke.ans",
			"data/arl-evoke.converted.ansi",
		},
		{
			"data/xz-gibson.ans",
			"data/empty_file.ansi",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.inputFpath, func(t *testing.T) {
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
