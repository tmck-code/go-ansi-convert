package test

import (
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
)

func TestANSITokenise(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]convert.ANSILineToken
	}{
		{
			name:  "Single line with no colour",
			input: "         ▄▄          ▄▄",
			expected: [][]convert.ANSILineToken{
				{
					{FG: "", BG: "", T: "         ▄▄          ▄▄"},
				},
			},
		},
		{
			// purple fg, red bg
			name:  "Single line with fg and bg",
			input: "\x1b[38;5;129mAAA\x1b[48;5;160mXX",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{T: "AAA", FG: "\x1b[38;5;129m", BG: "\x1b[49m"},
					convert.ANSILineToken{T: "XX", FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m"},
				},
			},
		},
		{
			// purple fg, red bg
			name:  "Longer single line with fg and bg",
			input: "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{T: "AAA    ", FG: "\x1b[38;5;129m", BG: "\x1b[49m"},
					convert.ANSILineToken{T: " XX ", FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m"},
				},
			},
		},
		{
			name: "Multi-line",
			// line 1 : purple fg,                  line 2: red bg
			input: "\x1b[38;5;160m▄\x1b[38;5;46m▄\n▄\x1b[38;5;190m▄",
			expected: [][]convert.ANSILineToken{
				{ // Line 1
					convert.ANSILineToken{FG: "\x1b[38;5;160m", BG: "", T: "▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;46m", BG: "", T: "▄"},
				},
				{ // Line 2
					convert.ANSILineToken{FG: "\x1b[38;5;46m", BG: "", T: "▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;190m", BG: "", T: "▄"},
				},
			},
		},
		// purple fg, red bg
		{
			name:  "Single line with spaces",
			input: "\x1b[38;5;129mAAA  \x1b[48;5;160m  XX  \x1b[0m",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{T: "AAA  ", FG: "\x1b[38;5;129m", BG: "\x1b[49m"},
					convert.ANSILineToken{T: "  XX  ", FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m"},
				},
			},
		},
		{
			name:  "Single line with existing ANSI reset",
			input: "\x1b[38;5;129mAAA\x1b[48;5;160mXX\x1b[0m",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{T: "AAA", FG: "\x1b[38;5;129m", BG: "\x1b[49m"},
					convert.ANSILineToken{T: "XX", FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m"},
				},
			},
		},
		{
			name: "Top of Egg",
			input: "    \x1b[49m   \x1b[38;5;16m▄▄\x1b[48;5;16m\x1b[38;5;142m▄▄▄\x1b[49m\x1b[38;5;16m▄▄\n" +
				"     ▄\x1b[48;5;16m\x1b[38;5;58m▄\x1b[48;5;58m\x1b[38;5;70m▄\x1b[48;5;70m \x1b[48;5;227m    \x1b[48;5;237m\x1b[38;5;227m▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[49m\x1b[38;5;16m▄",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "", BG: "", T: "    "},
					convert.ANSILineToken{FG: "", BG: "\x1b[49m", T: "   "},

					convert.ANSILineToken{FG: "\u001b[38;5;16m", BG: "\x1b[49m", T: "▄▄"},
					convert.ANSILineToken{FG: "\u001b[38;5;142m", BG: "\u001b[48;5;16m", T: "▄▄▄"},
					convert.ANSILineToken{FG: "\u001b[38;5;16m", BG: "\x1b[49m", T: "▄▄"},
				},
				{
					convert.ANSILineToken{FG: "\x1b[38;5;16m", BG: "\x1b[49m", T: "     ▄"},
					convert.ANSILineToken{FG: "\u001b[38;5;58m", BG: "\u001b[48;5;16m", T: "▄"},
					convert.ANSILineToken{FG: "\u001b[38;5;70m", BG: "\u001b[48;5;58m", T: "▄"},
					convert.ANSILineToken{FG: "\u001b[38;5;70m", BG: "\u001b[48;5;70m", T: " "},
					convert.ANSILineToken{FG: "\u001b[38;5;70m", BG: "\u001b[48;5;227m", T: "    "},
					convert.ANSILineToken{FG: "\u001b[38;5;227m", BG: "\u001b[48;5;237m", T: "▄"},
					convert.ANSILineToken{FG: "\u001b[38;5;237m", BG: "\u001b[48;5;16m", T: "▄"},
					convert.ANSILineToken{FG: "\u001b[38;5;16m", BG: "\u001b[49m", T: "▄"},
				},
			},
		},
		{
			name: "Multi-line with trailing spaces",
			// The AAA has a purple fg
			// The XX has a red bg
			input: "  \x1b[38;5;129mAAA \x1b[48;5;160m XY \x1b[0m     ",
			// The AAA should still have a purple fg
			// The XX should still have a red bg
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "", BG: "", T: "  "},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAA "},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XY "},
					convert.ANSILineToken{FG: "\x1b[0m", BG: "", T: "     "},
				},
			},
		},
		{
			name: "Lines with FG continuation (spaces)",
			// purple fg, red bg
			// the 4 spaces after AAA should have a purple fg, and no bg
			input: "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			// expected := "\x1b[0m\x1b[48;5;160m\x1b[38;5;129m XX \x1b[38;5;129m\x1b[49m    AAA\x1b[0m"
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAA    "},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX "},
				},
			},
		},
		{
			name: "Lines with BG continuation (spaces)",
			// purple fg, red bg
			// the 4 spaces after AAA should have a purple fg, and no bg
			input: "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m",
			// expected := "\x1b[0m\x1b[48;5;160m\x1b[38;5;129m XX \x1b[38;5;129m\x1b[49m    AAA\x1b[0m"
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAA    "},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX "},
				},
			},
		},
		{
			name:  "Reset followed by new color on same line",
			input: "\x1b[38;5;129mText\x1b[0m\x1b[38;5;160mMore",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "", T: "Text"},
					convert.ANSILineToken{FG: "\x1b[38;5;160m", BG: "", T: "More"},
				},
			},
		},
		{
			name:  "Multi-line where line 1 ends with reset",
			input: "\x1b[38;5;129mLine 1\x1b[0m\nLine 2",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "", T: "Line 1"},
				},
				{
					convert.ANSILineToken{FG: "", BG: "", T: "Line 2"},
				},
			},
		},
		{
			name:  "Multi-line where line 1 ends with reset and line 2 has color",
			input: "\x1b[38;5;129mLine 1\x1b[0m\n\x1b[38;5;160mLine 2",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "", T: "Line 1"},
				},
				{
					convert.ANSILineToken{FG: "\x1b[38;5;160m", BG: "", T: "Line 2"},
				},
			},
		},
		{
			name:  "Reset clears both FG and BG colors",
			input: "\x1b[38;5;129m\x1b[48;5;160mColored\x1b[0mPlain",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: "Colored"},
					convert.ANSILineToken{FG: "\x1b[0m", BG: "", T: "Plain"},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.TokeniseANSIString(tc.input)
			test.PrintANSITestResults(tc.input, tc.expected, result, t)
			test.Assert(tc.expected, result, t)
		})
	}
}
