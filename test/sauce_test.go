package test

import (
	"fmt"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
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
				ID:         "SAUCE",
				Version:    "00",
				Title:      "Evoke 2025",
				Author:     "Arlequin",
				Group:      "Impure",
				Date:       "20250922",
				FileSize:   4069,
				DataType:   convert.DataTypeCharacter,
				FileType:   convert.FileTypeCharacterANSI,
				TInfo1Name: "Character width",
				TInfo2Name: "Number of lines",
				TInfo3Name: "0",
				TInfo4Name: "0",
				TInfo1:     80,
				TInfo2:     25,
				TInfo3:     0,
				TInfo4:     0,
				Comments:   0,
				TFlags:     0x04,
				TInfoS:     "IBM VGA",
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
				ID:         "SAUCE",
				Version:    "00",
				Title:      "Test File",
				Author:     "Tester",
				Group:      "",
				Date:       "20260111",
				FileSize:   1024,
				DataType:   convert.DataTypeCharacter,
				FileType:   convert.FileTypeCharacterASCII,
				TInfo1Name: "Character width",
				TInfo2Name: "Number of lines",
				TInfo3Name: "0",
				TInfo4Name: "0",
				TInfo1:     80,
				TInfo2:     50,
				TInfo3:     0,
				TInfo4:     0,
				Comments:   0,
				TFlags:     0x01,
				TInfoS:     "IBM VGA50",
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
				ID:         "SAUCE",
				Version:    "00",
				Title:      "BIN Screen",
				Author:     "Artist",
				Group:      "Crew",
				Date:       "19950315",
				FileSize:   4080,
				DataType:   convert.DataTypeBinaryText,
				FileType:   0x28,
				TInfo1Name: "0",
				TInfo2Name: "0",
				TInfo3Name: "0",
				TInfo4Name: "0",
				TInfo1:     0,
				TInfo2:     0,
				TInfo3:     0,
				TInfo4:     0,
				Comments:   0,
				TFlags:     0x00,
				TInfoS:     "IBM EGA",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.ParseSAUCE(tc.input)

			PrintSimpleTestResults(fmt.Sprintf("%q", tc.input), fmt.Sprintf("%+v", tc.expected), fmt.Sprintf("%+v", result))
			Assert(tc.expected, result, t)

			// 		if tc.expected == nil {
			// 			if result != nil {
			// 				t.Errorf("ParseSAUCE() = %+v; want nil", result)
			// 			}
			// 			return
			// 		}

			// 		if result == nil {
			// 			t.Fatalf("ParseSAUCE() = nil; want %+v", tc.expected)
			// 		}

			// 		if result.ID != tc.expected.ID {
			// 			t.Errorf("ID = %q; want %q", result.ID, tc.expected.ID)
			// 		}
			// 		if result.Version != tc.expected.Version {
			// 			t.Errorf("Version = %q; want %q", result.Version, tc.expected.Version)
			// 		}
			// 		if result.Title != tc.expected.Title {
			// 			t.Errorf("Title = %q; want %q", result.Title, tc.expected.Title)
			// 		}
			// 		if result.Author != tc.expected.Author {
			// 			t.Errorf("Author = %q; want %q", result.Author, tc.expected.Author)
			// 		}
			// 		if result.Group != tc.expected.Group {
			// 			t.Errorf("Group = %q; want %q", result.Group, tc.expected.Group)
			// 		}
			// 		if result.Date != tc.expected.Date {
			// 			t.Errorf("Date = %q; want %q", result.Date, tc.expected.Date)
			// 		}
			// 		if result.FileSize != tc.expected.FileSize {
			// 			t.Errorf("FileSize = %d; want %d", result.FileSize, tc.expected.FileSize)
			// 		}
			// 		if result.DataType != tc.expected.DataType {
			// 			t.Errorf("DataType = %d; want %d", result.DataType, tc.expected.DataType)
			// 		}
			// 		if result.FileType != tc.expected.FileType {
			// 			t.Errorf("FileType = %d; want %d", result.FileType, tc.expected.FileType)
			// 		}
			// 		if result.TInfo1 != tc.expected.TInfo1 {
			// 			t.Errorf("TInfo1 = %d; want %d", result.TInfo1, tc.expected.TInfo1)
			// 		}
			// 		if result.TInfo2 != tc.expected.TInfo2 {
			// 			t.Errorf("TInfo2 = %d; want %d", result.TInfo2, tc.expected.TInfo2)
			// 		}
			// 		if result.TInfo3 != tc.expected.TInfo3 {
			// 			t.Errorf("TInfo3 = %d; want %d", result.TInfo3, tc.expected.TInfo3)
			// 		}
			// 		if result.TInfo4 != tc.expected.TInfo4 {
			// 			t.Errorf("TInfo4 = %d; want %d", result.TInfo4, tc.expected.TInfo4)
			// 		}
			// 		if result.Comments != tc.expected.Comments {
			// 			t.Errorf("Comments = %d; want %d", result.Comments, tc.expected.Comments)
			// 		}
			// 		if result.TFlags != tc.expected.TFlags {
			// 			t.Errorf("TFlags = 0x%02x; want 0x%02x", result.TFlags, tc.expected.TFlags)
			// 		}
			// 		if result.TInfoS != tc.expected.TInfoS {
			// 			t.Errorf("TInfoS = %q; want %q", result.TInfoS, tc.expected.TInfoS)
			// 		}
		})
	}
}

func TestParseInvalidSAUCE(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected *convert.SAUCE
	}{
		{
			name:     "No SAUCE present",
			input:    []byte("This is just regular data without SAUCE metadata at the end"),
			expected: nil,
		},
		{
			name:     "Data too short",
			input:    []byte("Short"),
			expected: nil,
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
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.ParseSAUCE(tc.input)

			Assert(tc.expected, result, t)
		})
	}
}
