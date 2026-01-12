package test

import (
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
)

// Test ANSI sanitisation ------------------------------------------------------

func TestSanitiseUnicodeString(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Single line without ANSI codes",
			input:    []string{"Hello World"},
			expected: []string{"Hello World\x1b[0m"},
		},
		{
			name:     "Single line with ANSI but no reset",
			input:    []string{"\x1b[38;5;129mColored text"},
			expected: []string{"\x1b[38;5;129mColored text\x1b[0m"},
		},
		{
			name:     "Single line with ANSI and existing reset",
			input:    []string{"\x1b[38;5;129mColored text\x1b[0m"},
			expected: []string{"\x1b[38;5;129mColored text\x1b[0m"},
		},
		{
			name:     "Multi-line without ANSI codes",
			input:    []string{"Line 1", "Line 2"},
			expected: []string{"Line 1\x1b[0m", "Line 2\x1b[0m"},
		},
		{
			name:     "Multi-line with ANSI codes",
			input:    []string{"\x1b[38;5;129mLine 1", "Line 2"},
			expected: []string{"\x1b[38;5;129mLine 1\x1b[0m", "\x1b[38;5;129mLine 2\x1b[0m"},
		},
		{
			name: "Multi-line with FG and BG colors",
			input: []string{
				"\x1b[38;5;129m\x1b[48;5;160mLine 1",
				"Line 2",
			},
			expected: []string{
				"\x1b[38;5;129m\x1b[48;5;160mLine 1\x1b[0m",
				"\x1b[38;5;129m\x1b[48;5;160mLine 2\x1b[0m",
			},
		},
		{
			name:     "Line with reset followed by color",
			input:    []string{"\x1b[38;5;129mText\x1b[0m\x1b[38;5;160mMore"},
			expected: []string{"\x1b[38;5;129mText\x1b[38;5;160mMore\x1b[0m"},
		},
		{
			name: "Multi-line with reset in middle",
			input: []string{
				"\x1b[38;5;129mLine 1\x1b[0m",
				"Line 2",
			},
			expected: []string{
				"\x1b[38;5;129mLine 1\x1b[0m",
				"Line 2\x1b[0m",
			},
		},
		{
			name:     "Empty string",
			input:    []string{""},
			expected: []string{""},
		},
		{
			name: "Multi-line with color continuation and reset",
			input: []string{
				"\x1b[38;5;160m▄\x1b[38;5;46m▄",
				"▄\x1b[38;5;190m▄",
			},
			expected: []string{
				"\x1b[38;5;160m▄\x1b[38;5;46m▄\x1b[0m",
				"\x1b[38;5;46m▄\x1b[38;5;190m▄\x1b[0m",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputStr := strings.Join(tc.input, "\n")
			expectedStr := strings.Join(tc.expected, "\n")

			result := convert.SanitiseUnicodeString(inputStr, false)

			test.PrintSimpleTestResults(inputStr, expectedStr, result, t)
			test.Assert(expectedStr, result, t)
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
			test.PrintSimpleTestResults(tc.input, tc.expected, result, t)
			test.Assert(tc.expected, result, t)
		})
	}
}
