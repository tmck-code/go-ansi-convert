package test

import (
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
)

func TestUnicodeStringLength(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "ASCII text",
			input:    "Hello World",
			expected: 11,
		},
		{
			name:     "Block/box drawing characters",
			input:    "┌─────┐│     │└─────┘",
			expected: 21,
		},
		{
			name:     "Full block characters",
			input:    "█▀▄▌▐░▒▓",
			expected: 8,
		},
		{
			name:     "Box drawing with double lines",
			input:    "╔═══╗║   ║╚═══╝",
			expected: 15,
		},
		{
			name:     "Japanese hiragana",
			input:    "こんにちは", // "hello" in hiragana
			expected: 10,      // 5 chars × 2 width each
		},
		{
			name:     "Japanese katakana",
			input:    "カタカナ", // katakana
			expected: 8,      // 4 chars × 2 width each
		},
		{
			name:     "Japanese kanji",
			input:    "日本語", // "Japanese language"
			expected: 6,     // 3 chars × 2 width each
		},
		{
			name:     "Mixed ASCII and Japanese",
			input:    "Hello世界", // "Hello world"
			expected: 9,         // 5 ASCII + 2 kanji (4 width)
		},
		{
			name:     "ANSI colored ASCII",
			input:    "\x1b[38;5;129mHello\x1b[0m",
			expected: 5, // ANSI codes don't count
		},
		{
			name:     "ANSI colored Japanese",
			input:    "\x1b[38;5;160m日本語\x1b[0m",
			expected: 6, // 3 kanji × 2 width, ANSI codes don't count
		},
		{
			name:     "ANSI colored box characters",
			input:    "\x1b[38;5;196m█\x1b[48;5;16m▀▄\x1b[0m",
			expected: 3, // 3 block chars, ANSI codes don't count
		},
		{
			name:     "Complex ANSI with mixed characters",
			input:    "\x1b[38;5;129mHello\x1b[0m \x1b[38;5;160m世界\x1b[0m \x1b[38;5;46m█▀\x1b[0m",
			expected: 13,
			// Hello 世界 █▀
			// "Hello " (6) + "世界 " (5) + "█▀" (2) = 13
		},
		{
			name:     "Multiple ANSI codes in sequence",
			input:    "\x1b[38;5;160m\x1b[48;5;16mTest\x1b[0m",
			expected: 4,
		},
		{
			name:     "Japanese with FG and BG colors",
			input:    "\x1b[38;5;129m\x1b[48;5;160mこんにちは\x1b[0m",
			expected: 10, // 5 hiragana chars × 2 width
		},
		{
			name:     "Box art",
			input:    " ▄  █ ▄███▄   █    █    ████▄       ▄ ▄   ████▄ █▄▄▄▄ █     ██▄",
			expected: 63,
		},
		{
			name:     "Mixed width Unicode with ANSI",
			input:    "\x1b[38;5;46m┌──┐\x1b[0m \x1b[38;5;160m日本\x1b[0m \x1b[38;5;129mABC\x1b[0m",
			expected: 13, // "┌──┐" (4) + " " (1) + "日本" (4) + " " (1) + "ABC" (3) = 13
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "Only ANSI codes",
			input:    "\x1b[38;5;129m\x1b[48;5;160m\x1b[0m",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.UnicodeStringLength(tc.input)
			if result != tc.expected {
				t.Errorf("UnicodeStringLength(%q) = %d; want %d", tc.input, result, tc.expected)
			}
		})
	}
}

// Test ANSI sanitisation ------------------------------------------------------

func TestSanitiseUnicodeString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single line without ANSI codes",
			input:    "Hello World",
			expected: "Hello World\x1b[0m",
		},
		{
			name:     "Single line with ANSI but no reset",
			input:    "\x1b[38;5;129mColored text",
			expected: "\x1b[38;5;129mColored text\x1b[0m",
		},
		{
			name:     "Single line with ANSI and existing reset",
			input:    "\x1b[38;5;129mColored text\x1b[0m",
			expected: "\x1b[38;5;129mColored text\x1b[0m",
		},
		{
			name: "Multi-line without ANSI codes",
			input: strings.Join(
				[]string{
					"Line 1",
					"Line 2",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"Line 1\x1b[0m",
					"Line 2\x1b[0m",
				},
				"\n",
			),
		},
		{
			name: "Multi-line with ANSI codes",
			input: strings.Join(
				[]string{
					"\x1b[38;5;129mLine 1",
					"Line 2",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"\x1b[38;5;129mLine 1\x1b[0m",
					"\x1b[38;5;129mLine 2\x1b[0m",
				},
				"\n",
			),
		},
		{
			name: "Multi-line with FG and BG colors",
			input: strings.Join(
				[]string{
					"\x1b[38;5;129m\x1b[48;5;160mLine 1",
					"Line 2",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"\x1b[38;5;129m\x1b[48;5;160mLine 1\x1b[0m",
					"\x1b[38;5;129m\x1b[48;5;160mLine 2\x1b[0m",
				},
				"\n",
			),
		},
		{
			name:     "Line with reset followed by color",
			input:    "\x1b[38;5;129mText\x1b[0m\x1b[38;5;160mMore",
			expected: "\x1b[38;5;129mText\x1b[38;5;160mMore\x1b[0m",
		},
		{
			name: "Multi-line with reset in middle",
			input: strings.Join(
				[]string{
					"\x1b[38;5;129mLine 1\x1b[0m",
					"Line 2",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"\x1b[38;5;129mLine 1\x1b[0m",
					"Line 2\x1b[0m",
				},
				"\n",
			),
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name: "Multi-line with color continuation and reset",
			input: strings.Join(
				[]string{
					"\x1b[38;5;160m▄\x1b[38;5;46m▄",
					"▄\x1b[38;5;190m▄",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"\x1b[38;5;160m▄\x1b[38;5;46m▄\x1b[0m",
					"\x1b[38;5;46m▄\x1b[38;5;190m▄\x1b[0m",
				},
				"\n",
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.SanitiseUnicodeString(tc.input, false)
			PrintSimpleTestResults(tc.input, tc.expected, result)
			Assert(tc.expected, result, t)
		})
	}
}

func TestSanitiseUnicodeStringWithJustify(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single line without ANSI codes - no justification needed",
			input:    "Hello World",
			expected: "Hello World\x1b[0m",
		},
		{
			name: "Multi-line without ANSI codes - different lengths",
			input: strings.Join(
				[]string{
					"Short",
					"Longer line",
					"Med",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"Short\x1b[0m      ",
					"Longer line\x1b[0m",
					"Med\x1b[0m        ",
				},
				"\n",
			),
		},
		{
			name: "Multi-line with ANSI codes - different lengths",
			input: strings.Join(
				[]string{
					"\x1b[38;5;129mShort",
					"Longer line",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"\x1b[38;5;129mShort\x1b[0m      ",
					"\x1b[38;5;129mLonger line\x1b[0m",
				},
				"\n",
			),
		},
		{
			name: "Multi-line with Unicode characters - different widths",
			input: strings.Join(
				[]string{
					"Hello",
					"世界",
					"Test",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"Hello\x1b[0m",
					"世界\x1b[0m ",
					"Test\x1b[0m ",
				},
				"\n",
			),
		},
		{
			name: "Multi-line with ANSI and Unicode - complex case",
			input: strings.Join(
				[]string{
					"\x1b[38;5;160mABC\x1b[0m   ",
					"\x1b[38;5;46m世界\x1b[0m  ",
					"\x1b[38;5;46mLonger\x1b[0m",
					"\x1b[38m\x1b[48;5;160mこんにちは\x1b[0m",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"\x1b[38;5;160mABC\x1b[0m       ",
					"\x1b[38;5;46m世界\x1b[0m      ",
					"\x1b[38;5;46mLonger\x1b[0m    ",
					"\x1b[38m\x1b[48;5;160mこんにちは\x1b[0m",
				},
				"\n",
			),
		},
		{
			name: "Multi-line with box characters",
			input: strings.Join(
				[]string{
					"█▀",
					"████",
					"▄",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"█▀\x1b[0m  ",
					"████\x1b[0m",
					"▄\x1b[0m   ",
				},
				"\n",
			),
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name: "Multi-line - all same length",
			input: strings.Join(
				[]string{
					"AAAA",
					"BBBB",
					"CCCC",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"AAAA\x1b[0m",
					"BBBB\x1b[0m",
					"CCCC\x1b[0m",
				},
				"\n",
			),
		},
		{
			name: "Single char vs longer line",
			input: strings.Join(
				[]string{
					"A",
					"Long line here",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"A\x1b[0m             ",
					"Long line here\x1b[0m",
				},
				"\n",
			),
		},
		{
			name: "Multi-line with FG and BG colors",
			input: strings.Join(
				[]string{
					"\x1b[38;5;129m\x1b[48;5;160mAB",
					"Longer",
				},
				"\n",
			),
			expected: strings.Join(
				[]string{
					"\x1b[38;5;129m\x1b[48;5;160mAB\x1b[0m    ",
					"\x1b[38;5;129m\x1b[48;5;160mLonger\x1b[0m",
				},
				"\n",
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.SanitiseUnicodeString(tc.input, true)
			PrintSimpleTestResults(tc.input, tc.expected, result)
			Assert(tc.expected, result, t)
		})
	}
}

// Test ANSI line reversal -----------------------------------------------------
//
// These are smaller "unit" tests for ANSI line reversal.
// - reverse individual lines
// - reverse multiple (newline separated) lines

func TestFlipHorizontal(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]convert.ANSILineToken
	}{
		{
			name: "Single line with ANSI colours",
			// The AAA has a purple fg, and the XX has a red bg
			input: "\x1b[38;5;129mAAA \x1b[48;5;160m XX \x1b[0m",
			// expected: "\x1b[0m\x1b[38;5;129m\x1b[48;5;160m XX \x1b[49m AAA\x1b[0m",
			expected: [][]convert.ANSILineToken{
				{
					{"", "", ""},
					{"\x1b[38;5;129m", "\x1b[48;5;160m", " XX "},
					{"\x1b[38;5;129m", "\x1b[49m", " AAA"},
				},
			},
		},
		{
			name: "Multi-line with ANSI colours",
			// purple fg, red bg
			// the 4 spaces after AAA should have a purple fg, and no bg
			input: "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			// expected: "\x1b[0m\x1b[38;5;129m\x1b[48;5;160m XX \x1b[0m\x1b[38;5;129m    AAA\x1b[0m",
			expected: [][]convert.ANSILineToken{
				{
					{FG: "", BG: "", T: ""},
					{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX "},
					{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "    AAA"},
				},
			},
		},
		{
			name: "Multi-line with trailing spaces",
			// The AAA has a purple fg, the XX has a red bg
			input: "  \x1b[38;5;129mAAA \x1b[48;5;160m XY \x1b[0m  ",
			// The AAA should still have a purple fg, and the XX should still have a red bg
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "", BG: "", T: ""},
					convert.ANSILineToken{FG: "\x1b[0m", BG: "", T: "  "},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " YX "},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: " AAA"},
					convert.ANSILineToken{FG: "", BG: "", T: "  "},
				},
			},
		},
		{
			name:  "Multi-line with colour continuation",
			input: "\x1b[38;5;160m▄ \x1b[38;5;46m▄\n▄ \x1b[38;5;190m▄",
			// expected: "\x1b[0m\x1b[38;5;46m▄\x1b[38;5;160m ▄\n\x1b[38;5;190m▄\x1b[38;5;46m ▄\x1b[0m",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "", BG: "", T: ""},
					convert.ANSILineToken{FG: "\x1b[38;5;46m", BG: "", T: "▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;160m", BG: "", T: " ▄"},
				},
				{
					convert.ANSILineToken{FG: "", BG: "", T: ""},
					convert.ANSILineToken{FG: "\x1b[38;5;190m", BG: "", T: "▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;46m", BG: "", T: " ▄"},
				},
			},
		},
		{
			// 	// Test ANSI pokemon reversal --------------------------------------------------
			// 	//
			// 	// These are larger "integration" tests for reversing ANSI strings.
			// 	// - reverse pokemon sprite (with & without ANSI colours)
			name: "flip pikachu without colour",
			input: strings.Join(
				[]string{
					"         ▄▄          ▄▄",
					"        ▄▄▄     ▄▄▄▄▄▄ ▄▄",
					"       ▄  ▄▀ ▄▄▄  ▄▄   ▄▀",
					"     ▄▄▄   ▄▄  ▄▄    ▄▀",
					"    ▄▄   ▄▄▄  ▄ ▀▄  ▄▄",
					"    ▀▄▄   ▄▄▄   ▄▄▄ ▄▀",
					"    ▀▄▄▄▄▄ ▄▄   ▄▄▄▄▄",
					"           ▄▄▄▄  ▄▄▀",
					"         ▀▄▄▄    ▄▀",
					"             ▀▀▄▀",
				},
				"\n",
			),
			expected: [][]convert.ANSILineToken{
				{{"", "", "  "}, {"", "", "▄▄          ▄▄         "}},
				{{"", "", ""}, {"", "", "▄▄ ▄▄▄▄▄▄     ▄▄▄        "}},
				{{"", "", ""}, {"", "", "▀▄   ▄▄  ▄▄▄ ▀▄  ▄       "}},
				{{"", "", "  "}, {"", "", "▀▄    ▄▄  ▄▄   ▄▄▄     "}},
				{{"", "", "   "}, {"", "", "▄▄  ▄▀ ▄  ▄▄▄   ▄▄    "}},
				{{"", "", "   "}, {"", "", "▀▄ ▄▄▄   ▄▄▄   ▄▄▀    "}},
				{{"", "", "    "}, {"", "", "▄▄▄▄▄   ▄▄ ▄▄▄▄▄▀    "}},
				{{"", "", "     "}, {"", "", "▀▄▄  ▄▄▄▄           "}},
				{{"", "", "      "}, {"", "", "▀▄    ▄▄▄▀         "}},
				{{"", "", "        "}, {"", "", "▀▄▀▀             "}},
			},
		},
		{
			name: "flip pikachu with colour",
			input: strings.Join(
				[]string{
					"    \x1b[49m     \x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;232m▄ \x1b[49m         \x1b[38;5;16m▄▄",
					"        ▄\x1b[48;5;16m\x1b[38;5;94m▄\x1b[48;5;232m▄\x1b[48;5;16m \x1b[49m    \x1b[38;5;16m▄▄▄▄\x1b[48;5;16m\x1b[38;5;214m▄\x1b[48;5;214m\x1b[38;5;94m▄\x1b[48;5;94m \x1b[48;5;16m▄\x1b[49m\x1b[38;5;16m▄",
					"       ▄\x1b[48;5;16m \x1b[48;5;94m \x1b[48;5;58m▄\x1b[49m▀ ▄\x1b[48;5;16m\x1b[38;5;214m▄▄\x1b[48;5;232m  \x1b[38;5;94m▄\x1b[48;5;214m▄\x1b[48;5;94m   \x1b[38;5;16m▄\x1b[49m▀",
					"     ▄\x1b[48;5;16m\x1b[38;5;214m▄\x1b[48;5;94m▄\x1b[48;5;214m   \x1b[48;5;16m▄\x1b[38;5;58m▄\x1b[48;5;214m  \x1b[38;5;232m▄\x1b[48;5;232m\x1b[38;5;94m▄\x1b[48;5;94m    \x1b[38;5;16m▄\x1b[49m▀",
					"    ▄\x1b[48;5;16m\x1b[38;5;214m▄\x1b[48;5;214m   \x1b[38;5;94m▄\x1b[48;5;94m\x1b[38;5;231m▄\x1b[48;5;214m\x1b[38;5;16m▄  \x1b[48;5;58m\x1b[38;5;214m▄\x1b[48;5;16m \x1b[49m\x1b[38;5;16m▀\x1b[48;5;94m▄  \x1b[48;5;16m\x1b[38;5;94m▄\x1b[49m\x1b[38;5;16m▄",
					"    ▀\x1b[48;5;214m▄\x1b[48;5;58m\x1b[38;5;214m▄\x1b[48;5;214m   \x1b[48;5;16m▄\x1b[48;5;232m\x1b[38;5;196m▄\x1b[48;5;214m▄   \x1b[48;5;16m\x1b[38;5;214m▄\x1b[49m\x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;94m▄\x1b[48;5;94m \x1b[38;5;16m▄\x1b[49m▀",
					"    ▀\x1b[48;5;94m▄\x1b[48;5;232m▄\x1b[48;5;94m▄\x1b[48;5;214m\x1b[38;5;94m▄▄ \x1b[48;5;196m\x1b[38;5;214m▄\x1b[38;5;232m▄\x1b[48;5;214m   \x1b[48;5;88m\x1b[38;5;214m▄\x1b[48;5;232m▄\x1b[48;5;52m\x1b[38;5;232m▄\x1b[48;5;16m\x1b[38;5;52m▄\x1b[49m\x1b[38;5;16m▄",
					"        \x1b[48;5;16m \x1b[48;5;94m  \x1b[48;5;232m\x1b[38;5;94m▄\x1b[48;5;214m\x1b[38;5;232m▄▄\x1b[48;5;232m\x1b[38;5;214m▄\x1b[48;5;214m  \x1b[48;5;88m▄\x1b[48;5;232m\x1b[38;5;16m▄\x1b[49m▀",
					"         ▀\x1b[48;5;94m▄▄\x1b[48;5;214m▄    \x1b[48;5;232m▄\x1b[49m▀",
					"             ▀▀\x1b[48;5;214m▄\x1b[49m▀\x1b[39m\x1b[39m",
				},
				"\n",
			),
			expected: [][]convert.ANSILineToken{
				{{"", "", "  "}, {"\x1b[38;5;16m", "\x1b[49m", "▄▄"}, {"\x1b[38;5;232m", "\x1b[49m", "         "}, {"\x1b[38;5;232m", "\x1b[48;5;16m", " ▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"", "\x1b[49m", "     "}, {"", "", "    "}},
				{{"", "", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;94m", " "}, {"\x1b[38;5;94m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄▄▄▄"}, {"\x1b[38;5;94m", "\x1b[49m", "    "}, {"\x1b[38;5;94m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;94m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄        "}},
				{{"", "", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;94m", "   "}, {"\x1b[38;5;94m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;232m", "  "}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄ ▀"}, {"\x1b[38;5;16m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", " "}, {"\x1b[38;5;16m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;16m", "\x1b[49m", "▄       "}},
				{{"", "", "  "}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;94m", "    "}, {"\x1b[38;5;94m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;232m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;58m", "\x1b[48;5;214m", "  "}, {"\x1b[38;5;58m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;214m", "   "}, {"\x1b[38;5;214m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄     "}},
				{{"", "", "   "}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "  ▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;214m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;214m", "  ▄"}, {"\x1b[38;5;231m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;214m", "   "}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄    "}},
				{{"", "", "   "}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;94m", " "}, {"\x1b[38;5;94m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;196m", "\x1b[48;5;214m", "   ▄"}, {"\x1b[38;5;196m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;214m", "   "}, {"\x1b[38;5;214m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀    "}},
				{{"", "", "    "}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;52m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;232m", "\x1b[48;5;52m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;88m", "▄"}, {"\x1b[38;5;232m", "\x1b[48;5;214m", "   "}, {"\x1b[38;5;232m", "\x1b[48;5;196m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;196m", "▄"}, {"\x1b[38;5;94m", "\x1b[48;5;214m", " ▄▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀    "}},
				{{"", "", "     "}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;88m", "▄"}, {"\x1b[38;5;214m", "\x1b[48;5;214m", "  "}, {"\x1b[38;5;214m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;232m", "\x1b[48;5;214m", "▄▄"}, {"\x1b[38;5;94m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "  "}, {"\x1b[38;5;16m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;16m", "\x1b[49m", "        "}},
				{{"", "", "      "}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;232m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;214m", "    ▄"}, {"\x1b[38;5;16m", "\x1b[48;5;94m", "▄▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀         "}},
				{{"", "", "        "}, {"\x1b[39m", "\x1b[49m", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▀"}, {"\x1b[38;5;16m", "\x1b[48;5;214m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀▀             "}},
			},
		},
		{
			name: "flip egg with colour",
			input: strings.Join(
				[]string{
					"     \x1b[49m   \x1b[38;5;16m▄▄\x1b[48;5;16m\x1b[38;5;142m▄▄▄\x1b[49m\x1b[38;5;16m▄▄",
					"      ▄\x1b[48;5;16m\x1b[38;5;58m▄\x1b[48;5;58m\x1b[38;5;70m▄\x1b[48;5;70m \x1b[48;5;227m    \x1b[48;5;237m\x1b[38;5;227m▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[49m\x1b[38;5;16m▄",
					"     ▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[48;5;70m\x1b[38;5;227m▄▄\x1b[48;5;227m    \x1b[38;5;70m▄▄\x1b[48;5;142m \x1b[48;5;16m\x1b[38;5;237m▄\x1b[49m\x1b[38;5;16m▄",
					"     \x1b[48;5;16m \x1b[48;5;227m       \x1b[48;5;70m\x1b[38;5;227m▄\x1b[38;5;58m▄\x1b[48;5;58m \x1b[48;5;142m \x1b[48;5;16m \x1b[49m",
					"     \x1b[48;5;16m \x1b[48;5;142m\x1b[38;5;237m▄\x1b[48;5;227m\x1b[38;5;142m▄\x1b[48;5;70m  \x1b[48;5;227m▄▄\x1b[38;5;58m▄\x1b[48;5;142m▄▄ \x1b[38;5;237m▄\x1b[48;5;16m \x1b[49m",
					"      \x1b[48;5;16m \x1b[48;5;142m▄   \x1b[48;5;58m    \x1b[38;5;234m▄\x1b[48;5;16m \x1b[49m",
					"       \x1b[38;5;16m▀▀\x1b[48;5;142m▄▄▄\x1b[48;5;58m▄▄\x1b[49m▀▀\x1b[39m\x1b[39m",
				},
				"\n",
			),
			expected: [][]convert.ANSILineToken{
				{{"", "", "   "}, {"\x1b[38;5;16m", "\x1b[49m", "▄▄"}, {"\x1b[38;5;142m", "\x1b[48;5;16m", "▄▄▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄▄"}, {"", "\x1b[49m", "   "}, {"", "", "     "}},
				{{"", "", " "}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;237m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;227m", "\x1b[48;5;237m", "▄"}, {"\x1b[38;5;70m", "\x1b[48;5;227m", "    "}, {"\x1b[38;5;70m", "\x1b[48;5;70m", " "}, {"\x1b[38;5;70m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;58m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄      "}},
				{{"", "", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▄"}, {"\x1b[38;5;237m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;70m", "\x1b[48;5;142m", " "}, {"\x1b[38;5;70m", "\x1b[48;5;227m", "▄▄"}, {"\x1b[38;5;227m", "\x1b[48;5;227m", "    "}, {"\x1b[38;5;227m", "\x1b[48;5;70m", "▄▄"}, {"\x1b[38;5;237m", "\x1b[48;5;16m", "▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▄     "}},
				{{"", "", ""}, {"\x1b[38;5;58m", "\x1b[49m", ""}, {"\x1b[38;5;58m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;58m", "\x1b[48;5;142m", " "}, {"\x1b[38;5;58m", "\x1b[48;5;58m", " "}, {"\x1b[38;5;58m", "\x1b[48;5;70m", "▄"}, {"\x1b[38;5;227m", "\x1b[48;5;70m", "▄"}, {"\x1b[38;5;16m", "\x1b[48;5;227m", "       "}, {"\x1b[38;5;16m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;16m", "\x1b[49m", "     "}},
				{{"", "", ""}, {"\x1b[38;5;237m", "\x1b[49m", ""}, {"\x1b[38;5;237m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;237m", "\x1b[48;5;142m", "▄"}, {"\x1b[38;5;58m", "\x1b[48;5;142m", " ▄▄"}, {"\x1b[38;5;58m", "\x1b[48;5;227m", "▄"}, {"\x1b[38;5;142m", "\x1b[48;5;227m", "▄▄"}, {"\x1b[38;5;142m", "\x1b[48;5;70m", "  "}, {"\x1b[38;5;142m", "\x1b[48;5;227m", "▄"}, {"\x1b[38;5;237m", "\x1b[48;5;142m", "▄"}, {"\x1b[38;5;58m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;58m", "\x1b[49m", "     "}},
				{{"", "", " "}, {"\x1b[38;5;234m", "\x1b[49m", ""}, {"\x1b[38;5;234m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;234m", "\x1b[48;5;58m", "▄"}, {"\x1b[38;5;237m", "\x1b[48;5;58m", "    "}, {"\x1b[38;5;237m", "\x1b[48;5;142m", "   ▄"}, {"\x1b[38;5;237m", "\x1b[48;5;16m", " "}, {"\x1b[38;5;237m", "\x1b[49m", "      "}},
				{{"", "", "  "}, {"\x1b[39m", "\x1b[49m", ""}, {"\x1b[38;5;16m", "\x1b[49m", "▀▀"}, {"\x1b[38;5;16m", "\x1b[48;5;58m", "▄▄"}, {"\x1b[38;5;16m", "\x1b[48;5;142m", "▄▄▄"}, {"\x1b[38;5;16m", "\x1b[49m", "▀▀"}, {"\x1b[38;5;234m", "\x1b[49m", "       "}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.FlipHorizontal(convert.TokeniseANSIString(tc.input))
			PrintANSITestResults(tc.input, tc.expected, result, t)
			Assert(tc.expected, result, t)
		})
	}
}
