package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"golang.org/x/text/encoding/charmap"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
)

func TestConvertAnsStrings(t *testing.T) {
	testCases := []struct {
		name           string
		inputString    string
		inputSAUCE     convert.SAUCE
		expectedString string
	}{
		{
			name:        "ANSI cursor codes",
			inputString: "\x1b[42m\x1b[10Cxxx\n",
			inputSAUCE: convert.SAUCE{
				ID:       "SAUCE",
				DataType: 1,
				FileType: 1,
				TInfo1: convert.TInfoField{
					Name: "Character Width", Value: 13,
				},
				TInfo2: convert.TInfoField{
					Name: "Number of lines", Value: 1,
				},
			},
			expectedString: "\x1b[42m          xxx\x1b[0m\n",
		},
		{
			name:        "Small single-line",
			inputString: "\x1b[31m\x1b[40m123\x1b[36m\x1b[43mabc\x1b[0m",
			inputSAUCE: convert.SAUCE{
				ID:       "SAUCE",
				DataType: 1,
				FileType: 1,
				TInfo1: convert.TInfoField{
					Name: "Character Width", Value: 3,
				},
				TInfo2: convert.TInfoField{
					Name: "Number of lines", Value: 2,
				},
			},
			expectedString: strings.Join(
				[]string{
					"\x1b[31m\x1b[40m123\x1b[0m",
					"\x1b[36m\x1b[43mabc\x1b[0m",
					"",
				},
				"\n",
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.ConvertAns(tc.inputString, tc.inputSAUCE)
			test.PrintSimpleTestResults(
				fmt.Sprintf("%+v\x1b[0m", tc.inputString),
				fmt.Sprintf("%+v\x1b[0m", tc.expectedString),
				fmt.Sprintf("%+v\x1b[0m", result),
			)
			test.Assert(tc.expectedString, result, t)
		})
	}
}

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
			"smallTwoLines",
			"../data/smallTwoLines.ans",
			"../data/smallTwoLines.converted.ansi",
		},
		// {
		// 	// https://16colo.rs/pack/impure89/xz-gibson.ans
		// 	"xz-gibson",
		// 	"../data/xz-gibson.ans",
		// 	"../data/xz-gibson.converted.ansi",
		// },
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
