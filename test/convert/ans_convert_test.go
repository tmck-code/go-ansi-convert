package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/convert"
	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/parse"
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
		{
			name:        "Line with truecolor codes",
			inputString: "\x1b[69C\x1b[1m\x1b[1;224;224;224t_\x1b[0m\x1b[1;168;168;168t.\x1b[1m\x1b[1;224;224;224txxXXx...\x1b[0m\x1b[1;168;168;168t \x1b[1m\x1b[1;224;224;224txXXX/\x1b[0m\x1b[1;168;168;168tXXXxxxXXxxx.",
			inputSAUCE: convert.SAUCE{
				ID:       "SAUCE",
				DataType: 1,
				FileType: 1,
				TInfo1:   convert.TInfoField{Name: "Character Width", Value: 170},
				TInfo2:   convert.TInfoField{Name: "Number of lines", Value: 1},
			},
			expectedString: "                                                                     \x1b[1m\x1b[38;2;224;224;224m_\x1b[0m\x1b[38;2;168;168;168m.\x1b[1m\x1b[38;2;224;224;224mxxXXx...\x1b[0m\x1b[38;2;168;168;168m \x1b[1m\x1b[38;2;224;224;224mxXXX/\x1b[0m\x1b[38;2;168;168;168mXXXxxxXXxxx.                                                                         \x1b[0m\n",
		},
		{
			name:        "Line with truecolor codes 2",
			inputString: "\x1b[37m\x1b[1;168;168;168t  \x1b[1m\x1b[1;224;224;224txxxx\x1b[0m\x1b[1;168;168;168t.²²'\x1b[34m\x1b[1;24;56;88t.xX²xx²²²x²XXXx. \x1b[37m\x1b[1;168;168;168t`XXXXXXl   \x1b[35m\x1b[1;144;16;64t    \x1b[45;30m\x1b[0;168;48;76tâ\x1b[40;35m\x1b[1;144;16;64tâ\"¯¯\"\x1b[45;30m\x1b[0;168;48;76tx\x1b[0m\r\n",
			inputSAUCE: convert.SAUCE{
				ID:       "SAUCE",
				DataType: 1,
				FileType: 1,
				TInfo1:   convert.TInfoField{Name: "Character Width", Value: 60},
				TInfo2:   convert.TInfoField{Name: "Number of lines", Value: 1},
			},
			expectedString: "\x1b[38;2;168;168;168m  \x1b[1m\x1b[38;2;224;224;224mxxxx\x1b[0m\x1b[38;2;168;168;168m.²²'\x1b[38;2;24;56;88m.xX²xx²²²x²XXXx. \x1b[38;2;168;168;168m`XXXXXXl   \x1b[38;2;144;16;64m\x1b[49m    \x1b[30m\x1b[48;2;168;48;76mâ\x1b[38;2;144;16;64m\x1b[40mâ\"¯¯\"\x1b[30m\x1b[48;2;168;48;76mx\x1b[0m           \x1b[0m\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.ConvertAns(tc.inputString, tc.inputSAUCE)
			test.PrintSimpleTestResults(
				fmt.Sprintf("%+v\x1b[0m", tc.inputString),
				fmt.Sprintf("%+v\x1b[0m", tc.expectedString),
				fmt.Sprintf("%+v\x1b[0m", result),
				t,
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
		{
			// https://16colo.rs/pack/impure89/xz-gibson.ans
			"xz-gibson",
			"../data/xz-gibson.ans",
			"../data/xz-gibson.converted.ansi",
		},
		{
			// https://16colo.rs/pack/impure91/bhe-peaceofmind.txt
			"bhe-peaceofmind",
			"../data/bhe-peaceofmind.txt",
			"../data/bhe-peaceofmind.converted.ansi",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read the original .ans file in CP437 encoding
			data, err := os.ReadFile(tc.inputFpath)
			if err != nil {
				t.Fatalf("Failed to read input file: %v", err)
			}
			encoding := parse.DetectEncoding(data)
			sauce, input, err := convert.ParseSAUCE(data, encoding)
			if err != nil {
				decodedData, decodeErr := parse.DecodeFileContents(data, encoding)
				if decodeErr != nil {
					t.Fatalf("Failed to parse SAUCE or decode file data: %v", err)
				}
				input = decodedData
			}

			// Read the expected converted .ansi file
			expectedBytes, err := os.ReadFile(tc.expectedFpath)
			if err != nil {
				t.Fatalf("Failed to read expected output file: %v", err)
			}
			expected := string(expectedBytes)

			fmt.Println("Converting", tc.inputFpath, "with detected SAUCE:", sauce)

			result := convert.ConvertAns(input, *sauce)
			test.Assert(expected, result, t)
		})
	}
}
