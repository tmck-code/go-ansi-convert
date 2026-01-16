package test

import (
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/convert"
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
		{
			name: "collects more cursor tokens",
			input: strings.Join(
				[]string{
					"\x1b[40m\x1b[32C\x1b[0;1;34mSpa\x1b[35mce: \x1b[32m212 \x1b[33mMegs\x1b[2H\x1b[32C\x1b[35mSpe\x1b[34med\x1b[s\r",
					"\x1b[u: \x1b[33m9600+ \x1b[32mBaud\x1b[3H   ▄▄█████▄▄▄▄▄▄▄\x1b[15C\x1b[31mSys\x1b[30mOp: \x1b[36mLe \x1b[s\r",
					"\x1b[u\x1b[37mChat\x1b[5C▄▄▄\x1b[7C\x1b[31m█    ▄▄▄\x1b[0;31m▄▄\x1b[1;30m▄\x1b[4H \x1b[32m▄█████████████\x1b[s\r",
				},
				"\n",
			),
			expected: [][]convert.ANSILineToken{
				{
					{FG: "", BG: "\x1b[40m", Control: "", T: "                                "},
					{FG: "\x1b[0;1;34m", BG: "\x1b[40m", Control: "", T: "Spa"},
					{FG: "\x1b[35m", BG: "\x1b[40m", Control: "", T: "ce: "},
					{FG: "\x1b[32m", BG: "\x1b[40m", Control: "", T: "212 "},
					{FG: "\x1b[33m", BG: "\x1b[40m", Control: "", T: "Megs"},
					{FG: "", BG: "", Control: "\x1b[2H", T: ""},
					{FG: "\x1b[33m", BG: "\x1b[40m", Control: "", T: "                                "},
					{FG: "\x1b[35m", BG: "\x1b[40m", Control: "", T: "Spe"},
					{FG: "\x1b[34m", BG: "\x1b[40m", Control: "", T: "ed"},
					{FG: "", BG: "", Control: "\x1b[s", T: ""},
				},
				{
					{FG: "", BG: "", Control: "\x1b[u", T: ""},
					{FG: "\x1b[34m", BG: "\x1b[40m", Control: "", T: ": "},
					{FG: "\x1b[33m", BG: "\x1b[40m", Control: "", T: "9600+ "},
					{FG: "\x1b[32m", BG: "\x1b[40m", Control: "", T: "Baud"},
					{FG: "", BG: "", Control: "\x1b[3H", T: ""},
					{FG: "\x1b[32m", BG: "\x1b[40m", Control: "", T: "   ▄▄█████▄▄▄▄▄▄▄               "},
					{FG: "\x1b[31m", BG: "\x1b[40m", Control: "", T: "Sys"},
					{FG: "\x1b[30m", BG: "\x1b[40m", Control: "", T: "Op: "},
					{FG: "\x1b[36m", BG: "\x1b[40m", Control: "", T: "Le "},
					{FG: "", BG: "", Control: "\x1b[s", T: ""},
				},
				{
					{FG: "", BG: "", Control: "\x1b[u", T: ""},
					{FG: "\x1b[37m", BG: "\x1b[40m", Control: "", T: "Chat     ▄▄▄       "},
					{FG: "\x1b[31m", BG: "\x1b[40m", Control: "", T: "█    ▄▄▄"},
					{FG: "\x1b[0;31m", BG: "\x1b[40m", Control: "", T: "▄▄"},
					{FG: "\x1b[1;30m", BG: "\x1b[40m", Control: "", T: "▄"},
					{FG: "", BG: "", Control: "\x1b[4H", T: ""},
					{FG: "\x1b[1;30m", BG: "\x1b[40m", Control: "", T: " "},
					{FG: "\x1b[32m", BG: "\x1b[40m", Control: "", T: "▄█████████████"},
					{FG: "", BG: "", Control: "\x1b[s", T: ""},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.TokeniseANSIString(tc.input)

			test.PrintANSITestResults(tc.input, tc.expected, result, t, false)
			test.Assert(tc.expected, result, t)
		})
	}
}
