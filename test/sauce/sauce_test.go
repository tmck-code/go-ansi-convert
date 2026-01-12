package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
)

func TestParseValidSAUCE(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected *convert.SAUCE
	}{
		{
			name: "Evoke 2025 ANSI art",
			input: append([]byte("Some ANSI art content here\x1a"), []byte{
				'S', 'A', 'U', 'C', 'E', '0', '0', // ID + Version
				'E', 'v', 'o', 'k', 'e', ' ', '2', '0', '2', '5', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Title (35 bytes)
				'A', 'r', 'l', 'e', 'q', 'u', 'i', 'n', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Author (20 bytes)
				'I', 'm', 'p', 'u', 'r', 'e', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Group (20 bytes)
				'2', '0', '2', '5', '0', '9', '2', '2', // Date (8 bytes)
				0xe5, 0x0f, 0x00, 0x00, // FileSize (4 bytes, little-endian: 0x0fe5 = 4069)
				0x01,       // DataType (Character)
				0x01,       // FileType (ANSi)
				0x50, 0x00, // TInfo1 (2 bytes, little-endian: 0x0050 = 80 - character width)
				0x19, 0x00, // TInfo2 (2 bytes, little-endian: 0x0019 = 25 - number of lines)
				0x00, 0x00, // TInfo3 (2 bytes)
				0x00, 0x00, // TInfo4 (2 bytes)
				0x00,                                                                                                                        // Comments (1 byte)
				0x04,                                                                                                                        // TFlags (1 byte - ANSiFlags: 0x04 = 9 pixel font)
				'I', 'B', 'M', ' ', 'V', 'G', 'A', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // TInfoS (22 bytes)
			}...),
			expected: &convert.SAUCE{
				ID:       "SAUCE",
				Version:  "00",
				Title:    "Evoke 2025",
				Author:   "Arlequin",
				Group:    "Impure",
				Date:     "20250922",
				FileSize: 4069,
				DataType: convert.DataTypeCharacter,
				FileType: convert.FileTypeCharacterANSI,
				TInfo1: convert.TInfoField{
					Name:  "Character width",
					Value: 80,
				},
				TInfo2: convert.TInfoField{
					Name:  "Number of lines",
					Value: 25,
				},
				TInfo3: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				TInfo4: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				Comments: 0,
				TFlags:   0x04,
				TInfoS:   "IBM VGA",
			},
		},
		{
			name: "ASCII text with iCE Color",
			input: append([]byte("ASCII content\x1a"), []byte{
				'S', 'A', 'U', 'C', 'E', '0', '0', // ID + Version
				'T', 'e', 's', 't', ' ', 'F', 'i', 'l', 'e', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Title (35 bytes)
				'T', 'e', 's', 't', 'e', 'r', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Author (20 bytes)
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Group (20 bytes, empty)
				'2', '0', '2', '6', '0', '1', '1', '1', // Date (8 bytes)
				0x00, 0x04, 0x00, 0x00, // FileSize (4 bytes, little-endian: 0x0400 = 1024)
				0x01,       // DataType (Character)
				0x00,       // FileType (ASCII)
				0x50, 0x00, // TInfo1 (80 character width)
				0x32, 0x00, // TInfo2 (50 lines)
				0x00, 0x00, // TInfo3
				0x00, 0x00, // TInfo4
				0x00,                                                                                                                      // Comments
				0x01,                                                                                                                      // TFlags (iCE Color enabled)
				'I', 'B', 'M', ' ', 'V', 'G', 'A', '5', '0', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // TInfoS (22 bytes)
			}...),
			expected: &convert.SAUCE{
				ID:       "SAUCE",
				Version:  "00",
				Title:    "Test File",
				Author:   "Tester",
				Group:    "",
				Date:     "20260111",
				FileSize: 1024,
				DataType: convert.DataTypeCharacter,
				FileType: convert.FileTypeCharacterASCII,
				TInfo1: convert.TInfoField{
					Name:  "Character width",
					Value: 80,
				},
				TInfo2: convert.TInfoField{
					Name:  "Number of lines",
					Value: 50,
				},
				TInfo3: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				TInfo4: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				Comments: 0,
				TFlags:   0x01,
				TInfoS:   "IBM VGA50",
			},
		},
		{
			name: "BinaryText file",
			input: append([]byte("BIN screen data here\x1a"), []byte{
				'S', 'A', 'U', 'C', 'E', '0', '0', // ID + Version
				'B', 'I', 'N', ' ', 'S', 'c', 'r', 'e', 'e', 'n', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Title (35 bytes)
				'A', 'r', 't', 'i', 's', 't', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Author (20 bytes)
				'C', 'r', 'e', 'w', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Group (20 bytes)
				'1', '9', '9', '5', '0', '3', '1', '5', // Date (8 bytes)
				0xf0, 0x0f, 0x00, 0x00, // FileSize (4 bytes)
				0x05,       // DataType (BinaryText)
				0x28,       // FileType (width = 40*2 = 80)
				0x00, 0x00, // TInfo1
				0x00, 0x00, // TInfo2
				0x00, 0x00, // TInfo3
				0x00, 0x00, // TInfo4
				0x00,                                                                                                                        // Comments
				0x00,                                                                                                                        // TFlags
				'I', 'B', 'M', ' ', 'E', 'G', 'A', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // TInfoS (22 bytes)
			}...),
			expected: &convert.SAUCE{
				ID:       "SAUCE",
				Version:  "00",
				Title:    "BIN Screen",
				Author:   "Artist",
				Group:    "Crew",
				Date:     "19950315",
				FileSize: 4080,
				DataType: convert.DataTypeBinaryText,
				FileType: 0x28,
				TInfo1: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				TInfo2: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				TInfo3: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				TInfo4: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				Comments: 0,
				TFlags:   0x00,
				TInfoS:   "IBM EGA",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _, err := convert.ParseSAUCE(tc.input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing SAUCE record: %v\n", err)
				os.Exit(1)
			}

			test.PrintSAUCETestResults(string(tc.input), tc.expected, result, t)
			test.Assert(tc.expected, result, t)
		})
	}
}

func TestParseInvalidSAUCE(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "No SAUCE present",
			input:    []byte("This is just regular data without SAUCE metadata at the end"),
			expected: "data too short to contain SAUCE record",
		},
		{
			name:     "Data too short",
			input:    []byte("Short"),
			expected: "data too short to contain SAUCE record",
		},
		{
			name: "Invalid ID",
			input: append([]byte("dummy\x1a"), []byte{
				'N', 'O', 'T', 'I', 'T', '0', '0', // Wrong ID
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Title (35 bytes)
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Author (20 bytes)
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Group (20 bytes)
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', // Date (8 bytes)
				0x00, 0x00, 0x00, 0x00, // FileSize
				0x00, 0x00, // DataType, FileType
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // TInfo1-4
				0x00, 0x00, // Comments, TFlags
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // TInfoS (22 bytes)
			}...),
			expected: "no valid SAUCE record found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := convert.ParseSAUCE(tc.input)
			var result string
			if err != nil {
				result = err.Error()
			}

			test.PrintSimpleTestResults(string(tc.input), tc.expected, result, t)
			test.Assert(tc.expected, result, t)
		})
	}
}

func TestParseSAUCEFiles(t *testing.T) {
	testCases := []struct {
		path     string
		expected *convert.SAUCE
	}{
		{
			"../data/arl-evoke.ans",
			&convert.SAUCE{
				ID:       "SAUCE",
				Version:  "00",
				Title:    "Evoke 2025",
				Author:   "Arlequin",
				Group:    "Impure",
				Date:     "20250922",
				FileSize: 4069,
				DataType: convert.DataTypeCharacter,
				FileType: convert.FileTypeCharacterANSI,
				TInfo1: convert.TInfoField{
					Name:  "Character width",
					Value: 80,
				},
				TInfo2: convert.TInfoField{
					Name:  "Number of lines",
					Value: 25,
				},
				TInfo3: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				TInfo4: convert.TInfoField{
					Name:  "0",
					Value: 0,
				},
				Comments: 0,
				TFlags:   0x04,
				TInfoS:   "IBM VGA",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			data, err := os.ReadFile(tc.path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", tc.path, err)
				os.Exit(1)
			}

			result, _, err := convert.ParseSAUCE(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing SAUCE from file %s: %v\n", tc.path, err)
				os.Exit(1)
			}

			test.PrintSAUCETestResults(tc.path, tc.expected, result, t)
			test.Assert(tc.expected, result, t)
		})
	}
}
