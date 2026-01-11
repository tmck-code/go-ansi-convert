package test

import (
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
)

func TestConvertCursorCodes(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]convert.ANSILineToken
	}{
		{
			name:  "Collects cursor tokens",
			input: "\x1b[32mXX\x1b[10CYY",
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[32m", BG: "", T: "XX          YY"},
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
