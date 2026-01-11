package test

import (
	"fmt"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
)

type AdjustANSILineWidthsParams struct {
	lines       [][]convert.ANSILineToken
	targetWidth int
	targetLines int
}

func TestAdjustANSILineWidths(t *testing.T) {
	testCases := []struct {
		name     string
		input    AdjustANSILineWidthsParams
		expected [][]convert.ANSILineToken
	}{
		{
			name: "Split lines to match target width and lines",
			input: AdjustANSILineWidthsParams{
				lines: [][]convert.ANSILineToken{
					{
						convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAA"},
						convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX "},
						convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[49m", T: "BBBBB"},
						convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[48;5;28m", T: "YY"},
					},
				},
				targetWidth: 7,
				targetLines: 2,
			},
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAA"},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX "},
				},
				{
					convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[49m", T: "BBBBB"},
					convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[48;5;28m", T: "YY"},
				},
			},
		},
		{
			name: "Split token text if needed when adjusting line widths",
			input: AdjustANSILineWidthsParams{
				lines: [][]convert.ANSILineToken{
					{
						convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAAAAAA"},
						convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX"},
					},
				},
				targetWidth: 5,
				targetLines: 2,
			},
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAAAA"},
				},
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AA"},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX"},
				},
			},
		},
		{
			name: "Pad lines to match target width and lines",
			input: AdjustANSILineWidthsParams{
				lines: [][]convert.ANSILineToken{
					{
						// length 7
						convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAA"},
						convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX "},
					},
					{
						// length 8
						convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[49m", T: "BBBBB"},
						convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[48;5;28m", T: "YYZ"},
					},
				},
				targetWidth: 10,
				targetLines: 2,
			},
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[49m", T: "AAA"},
					convert.ANSILineToken{FG: "\x1b[38;5;129m", BG: "\x1b[48;5;160m", T: " XX "},
					convert.ANSILineToken{FG: "\x1b[0m", BG: "\x1b[0m", T: "   "},
				},
				{
					convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[49m", T: "BBBBB"},
					convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[48;5;28m", T: "YYZ"},
					convert.ANSILineToken{FG: "\x1b[0m", BG: "\x1b[0m", T: "  "},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convert.AdjustANSILineWidths(tc.input.lines, tc.input.targetWidth, tc.input.targetLines)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			test.PrintSimpleTestResults(
				fmt.Sprintf("targetWidth: %+v, targetLines: %+v\n%+v\x1b[0m", tc.input.targetWidth, tc.input.targetLines, tc.input.lines),
				fmt.Sprintf("%+v\x1b[0m", tc.expected),
				fmt.Sprintf("%+v\x1b[0m", result),
			)

			test.Assert(tc.expected, result, t)
		})
	}
}
