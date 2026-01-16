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
			encoding := parse.DetectEncoding(tc.input)
			result, _, err := convert.ParseSAUCE(tc.input, encoding)
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
			encoding := parse.DetectEncoding(tc.input)
			_, _, err := convert.ParseSAUCE(tc.input, encoding)
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
				t.Fatalf("Error reading file %s: %v", tc.path, err)
			}
			encoding := parse.DetectEncoding(data)
			result, _, err := convert.ParseSAUCE(data, encoding)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing SAUCE from file %s: %v\n", tc.path, err)
				os.Exit(1)
			}

			test.PrintSAUCETestResults(tc.path, tc.expected, result, t)
			test.Assert(tc.expected, result, t)
		})
	}
}

func TestParseFileDataWithEncoding(t *testing.T) {
	testCases := []struct {
		name             string
		path             string
		expectedEncoding string
		expectedData     string
	}{
		{
			name: "ISO-8859-1 encoded file",
			path: "../data/h7-matt.ans",
			expectedData: strings.Join(
				[]string{
					"\x1b[0;40;37m\x1b[0;0;0;0t\x1b[1;171;171;171t",
					"",
					"",
					"",
					"\x1b[34C\x1b[1;30m\x1b[1;87;87;87t ___________",
					"\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[30C\x1b[1;30m\x1b[1;87;87;87t ___| \x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[9C\x1b[1;30m\x1b[1;87;87;87t/.___",
					"\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[30C\x1b[1;30m\x1b[1;87;87;87t(\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t   \x1b[1;30m\x1b[1;87;87;87t!\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[11C\x1b[1;30m\x1b[1;87;87;87t!\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t   \x1b[1;30m\x1b[1;87;87;87t)",
					"\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[30C\x1b[1;30m\x1b[1;87;87;87t(\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t   \x1b[1;30m\x1b[1;87;87;87t:\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t   \x1b[35m\x1b[1;171;0;171t______ \x1b[37m\x1b[1;171;171;171t \x1b[1;30m\x1b[1;87;87;87t:\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t  \x1b[1;30m\x1b[1;87;87;87t )",
					"\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[30C\x1b[1;30m\x1b[1;87;87;87t(\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t   \x1b[1;30m\x1b[1;87;87;87t.\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t  \x1b[35m\x1b[1;171;0;171t/\x1b[34m\x1b[1;0;0;171t¯¯¯¯¯¯\x1b[35m\x1b[1;171;0;171t\\ \x1b[1;30m\x1b[1;87;87;87t.\x1b[32m\x1b[1;87;255;87t  \x1b[30m\x1b[1;87;87;87t )",
					"\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[25C\x1b[35m\x1b[1;171;0;171t   _______ /\x1b[34m\x1b[1;0;0;171t·\x1b[35m\x1b[1;171;0;171t   \x1b[36m\x1b[1;0;171;171t \x1b[35m\x1b[1;171;0;171t   \\   \x1b[1;32m\x1b[1;87;255;87t \x1b[30m\x1b[1;87;87;87t·\x1b[0;35m\x1b[0;0;0;0t\x1b[1;171;0;171t___    ___",
					"\x1b[37m\x1b[1;171;171;171t\x1b[19C\x1b[35m\x1b[1;171;0;171t  ____  /\x1b[34m\x1b[1;0;0;171t¯¯¯¯¯¯\x1b[35m\x1b[1;171;0;171t//\x1b[34m\x1b[1;0;0;171t/\x1b[35m\x1b[1;171;0;171t\x1b[6C\x1b[36m\x1b[1;0;171;171t  \x1b[35m\x1b[1;171;0;171t \\  _/\x1b[34m\x1b[1;0;0;171t¯¯¯\x1b[35m\x1b[1;171;0;171t\\__/\x1b[34m\x1b[1;0;0;171t¯¯¯\x1b[35m\x1b[1;171;0;171t\\__",
					"\x1b[37m\x1b[1;171;171;171t\x1b[16C\x1b[1;32m\x1b[1;87;255;87t \x1b[0;34m\x1b[0;0;0;0t\x1b[1;0;0;171t/\x1b[36m\x1b[1;0;171;171t \x1b[34m\x1b[1;0;0;171t/\x1b[36m\x1b[1;0;171;171t/\x1b[34m\x1b[1;0;0;171t·¯¯¯\x1b[36m\x1b[1;0;171;171t\\/  _   //\x1b[34m\x1b[1;0;0;171t/\x1b[36m\x1b[1;0;171;171t    \x1b[1;34m\x1b[1;87;87;255t/\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t\\\x1b[5C\\/\x1b[34m\x1b[1;0;0;171t \x1b[36m\x1b[1;0;171;171t\x1b[5C_/\x1b[34m\x1b[1;0;0;171t·\x1b[36m\x1b[1;0;171;171t    _/",
					"\x1b[37m\x1b[1;171;171;171t\x1b[17C\x1b[1;32m\x1b[1;87;255;87t \x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t_/\x1b[34m\x1b[1;0;0;171t/\x1b[36m\x1b[1;0;171;171t  _   \x1b[34m\x1b[1;0;0;171t/\x1b[36m\x1b[1;0;171;171t/\x1b[34m\x1b[1;0;0;171t¯ \x1b[36m\x1b[1;0;171;171t  \\\x1b[5C¯¯¯¯ \x1b[34m\x1b[1;0;0;171t_\x1b[36m\x1b[1;0;171;171t   \\\x1b[34m\x1b[1;0;0;171t \x1b[36m\x1b[1;0;171;171t   \x1b[34m\x1b[1;0;0;171t/\x1b[36m\x1b[1;0;171;171t//\x1b[34m\x1b[1;0;0;171t/\x1b[36m\x1b[1;0;171;171t    /\x1b[34m\x1b[1;0;0;171t¯",
					"\x1b[37m\x1b[1;171;171;171t\x1b[15C\x1b[1;30m\x1b[1;87;87;87t   \x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t\\\x1b[1;32m\x1b[1;87;255;87t_____\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t\\\x1b[1;32m\x1b[1;87;255;87t__\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t/\x1b[1;32m\x1b[1;87;255;87t_____/____\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t/¯¯¯¯\\\x1b[34m\x1b[1;0;0;171t\\\x1b[36m\x1b[1;0;171;171t    \\\x1b[1;32m\x1b[1;87;255;87t___\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t/\\\x1b[1;32m\x1b[1;87;255;87t_____\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t/\x1b[34m\x1b[1;0;0;171t/\x1b[1;32m\x1b[1;87;255;87t \x1b[0;34m\x1b[0;0;0;0t\x1b[1;0;0;171t/",
					"\x1b[37m\x1b[1;171;171;171t\x1b[19C\x1b[34m\x1b[1;0;0;171t¯¯¯¯¯¯¯¯¯¯¯\x1b[1;30m\x1b[1;87;87;87t_\x1b[0;34m\x1b[0;0;0;0t\x1b[1;0;0;171t¯¯¯¯¯¯¯\x1b[37m\x1b[1;171;171;171t\x1b[6C\x1b[36m\x1b[1;0;171;171t\\\x1b[1;32m\x1b[1;87;255;87t_____\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t\\\x1b[1;30m\x1b[1;87;87;87tH7/dS!\x1b[0;34m\x1b[0;0;0;0t\x1b[1;0;0;171t¯¯¯",
					"\x1b[37m\x1b[1;171;171;171t\x1b[17C\x1b[41;30m\x1b[0;171;0;0t\x1b[1;0;0;0t \x1b[40;37m\x1b[0;0;0;0t\x1b[1;171;171;171t \x1b[41;30m\x1b[0;171;0;0t\x1b[1;0;0;0t happy birthday mAtt! \x1b[40;37m\x1b[0;0;0;0t\x1b[1;171;171;171t \x1b[41;30m\x1b[0;171;0;0t\x1b[1;0;0;0t \x1b[40;37m\x1b[0;0;0;0t\x1b[1;171;171;171t  \x1b[34m\x1b[1;0;0;171t¯\x1b[1;30m\x1b[1;87;87;87t.\x1b[0;34m\x1b[0;0;0;0t\x1b[1;0;0;171t¯¯¯\x1b[1;30m\x1b[1;87;87;87t_",
					"\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[30C\x1b[1;30m\x1b[1;87;87;87t( \x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t  \x1b[1;30m\x1b[1;87;87;87t:\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[11C\x1b[1;30m\x1b[1;87;87;87t:\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t   \x1b[1;30m\x1b[1;87;87;87t)",
					"\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[30C\x1b[1;30m\x1b[1;87;87;87t(___|_\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[9C\x1b[1;30m\x1b[1;87;87;87t |___)",
					"\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[33C\x1b[1;30m\x1b[1;87;87;87t '/__________|",
					"",
				}, "\r\n"),
		},

		// {
		// 	name:             "ISO-8859-1 encoded file with cursor and truecolor codes",
		// 	path:             "../data/xz-gibson.ans",
		// 	expectedEncoding: "iso-8859-1",
		// 	expectedData: strings.Join(
		// 		[]string{
		// 			"\x1b[0;40;37m\x1b[0;0;0;0t\x1b[1;171;171;171t\r",
		// 			"\r",
		// 		},
		// 		"\n",
		// 	),
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			data, err := os.ReadFile(tc.path)
			if err != nil {
				t.Fatalf("Error reading file %s: %v\n", tc.path, err)
			}
			encoding := parse.DetectEncoding(data)
			_, result, err := convert.ParseSAUCE(data, encoding)
			if err != nil {
				t.Fatalf("Error parsing SAUCE from %s: %v\n", tc.path, err)
			}
			test.PrintSimpleTestResults(tc.path, tc.expectedData, result, t)
			test.Assert(tc.expectedData, result, t)
		})
	}
}
