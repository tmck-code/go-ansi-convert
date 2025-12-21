package test

import (
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
)

func TestTokeniseRedundant(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]convert.ANSILineToken
	}{
		{
			name:  "Tokenise consecutive identical color codes",
			input: strings.Repeat("\x1b[0;37;46m ", 4),
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenised := convert.TokeniseANSIString(tc.input)

			Assert(tc.expected, tokenised, t)
		})
	}
}

func TestOptimiseANSITokens(t *testing.T) {
	testCases := []struct {
		name     string
		input    [][]convert.ANSILineToken
		expected [][]convert.ANSILineToken
	}{
		{
			name: "Optimise redundant color codes",
			input: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
				},
			},
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[46m", T: "    "},
				},
			},
		},
		{
			name: "Optimise background color changes",
			input: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[42m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[42m", T: " "},
				},
			},
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[46m", T: "  "},
					{FG: "", BG: "\x1b[42m", T: "  "},
				},
			},
		},
		{
			name: "Optimise foreground color changes",
			input: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[42m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[42m", T: " "},
					{FG: "\x1b[32m", BG: "\x1b[42m", T: " "},
					{FG: "\x1b[32m", BG: "\x1b[42m", T: " "},
				},
			},
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[42m", T: "  "},
					{FG: "\x1b[32m", BG: "", T: "  "},
				},
			},
		},
		{
			name: "Optimise mixed color changes",
			input: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[46m", T: " "},
					{FG: "\x1b[32m", BG: "\x1b[42m", T: " "},
					{FG: "\x1b[32m", BG: "\x1b[42m", T: " "},
				},
			},
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[37m", BG: "\x1b[46m", T: "  "},
					{FG: "\x1b[32m", BG: "\x1b[42m", T: "  "},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			optimised := convert.OptimiseANSITokens(tc.input)
			Assert(tc.expected, optimised, t)
		})
	}
}
