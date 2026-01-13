package test

import (
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/convert"
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

					convert.ANSILineToken{FG: "\x1b[38;5;16m", BG: "\x1b[49m", T: "▄▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;142m", BG: "\x1b[48;5;16m", T: "▄▄▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;16m", BG: "\x1b[49m", T: "▄▄"},
				},
				{
					convert.ANSILineToken{FG: "\x1b[38;5;16m", BG: "\x1b[49m", T: "     ▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;58m", BG: "\x1b[48;5;16m", T: "▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;70m", BG: "\x1b[48;5;58m", T: "▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;70m", BG: "\x1b[48;5;70m", T: " "},
					convert.ANSILineToken{FG: "\x1b[38;5;70m", BG: "\x1b[48;5;227m", T: "    "},
					convert.ANSILineToken{FG: "\x1b[38;5;227m", BG: "\x1b[48;5;237m", T: "▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;237m", BG: "\x1b[48;5;16m", T: "▄"},
					convert.ANSILineToken{FG: "\x1b[38;5;16m", BG: "\x1b[49m", T: "▄"},
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
		{
			name:  "Tokenises each new character with a preceding color",
			input: strings.Repeat("\x1b[32m"+"X", 3),
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[32m", BG: "", T: "X"},
					convert.ANSILineToken{FG: "\x1b[32m", BG: "", T: "X"},
					convert.ANSILineToken{FG: "\x1b[32m", BG: "", T: "X"},
				},
			},
		},
		{
			name:  "Skips tokenising colours if they are not followed by any text characters",
			input: "\x1b[0m" + strings.Repeat("\x1b[32m", 3) + "\x1b[33m text here!",
			expected: [][]convert.ANSILineToken{
				{
					convert.ANSILineToken{FG: "\x1b[33m", BG: "", T: " text here!"},
				},
			},
		},
		{
			name:  "Handles control characters and truecolor codes",
			input: "\x1b[37m\x1b[1;168;168;168t  \x1b[1m\x1b[1;224;224;224txxxx\x1b[0m\x1b[1;168;168;168t.²²'\x1b[34m\x1b[1;24;56;88t.xX²xx²²²x²XXXx. \x1b[37m\x1b[1;168;168;168t`XXXXXXl   \x1b[35m\x1b[1;144;16;64t    \x1b[45;30m\x1b[0;168;48;76tâ\x1b[40;35m\x1b[1;144;16;64tâ\"¯¯\"\x1b[45;30m\x1b[0;168;48;76tx\x1b[0m\r\n",
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[38;2;168;168;168m", BG: "", T: "  "},
					{FG: "\x1b[1m\x1b[38;2;224;224;224m", BG: "", T: "xxxx"},
					{FG: "\x1b[0m", BG: "", T: ""},
					{FG: "\x1b[38;2;168;168;168m", BG: "", T: ".²²'"},
					{FG: "\x1b[38;2;24;56;88m", BG: "", T: ".xX²xx²²²x²XXXx. "},
					{FG: "\x1b[38;2;168;168;168m", BG: "", T: "`XXXXXXl   "},
					{FG: "\x1b[38;2;144;16;64m", BG: "\x1b[49m", T: "    "},
					{FG: "\x1b[30m", BG: "\x1b[45m", T: "â"},
					{FG: "\x1b[38;2;144;16;64m", BG: "\x1b[40m", T: "â\"¯¯\""},
					{FG: "\x1b[30m", BG: "\x1b[45m", T: "x"},
					{FG: "\x1b[0m", BG: "", T: ""},
				},
			},
		},
		{
			name:  "tokenise long line",
			input: "\x1b[0;31;40m \x1b[0;31;40m \x1b[0;31;40m▐\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;30;41m░\x1b[0;31;40m▓\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m▌\x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m ",
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[31m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▐"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[30m", BG: "\x1b[41m", T: "░"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▓"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▌"},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
				},
			},
		},
		{
			name:  "Tokenise long line 2",
			input: "\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m▄\x1b[0;31;40m▄\x1b[0;31;40m▄\x1b[0;31;40m▄\x1b[0;31;40m \x1b[0;37;40m \x1b[0;30;43m▐\x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;30;44m█\x1b[0;30;44m█\x1b[0;30;44m█\x1b[0;30;44m█\x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;37;40m \x1b[0;30;43m▓\x1b[0;37;40m \x1b[0;31;40m█\x1b[0;37;40m \x1b[0;37;40m \x1b[0;31;40m█\x1b[0;31;40m▌\x1b[0;37;40m \x1b[0;30;43m▓\x1b[0;30;43m█\x1b[0;31;40m█\x1b[0;37;40m \x1b[0;37;40m \x1b[0;31;40m▐\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█\x1b[0;31;40m█",
			expected: [][]convert.ANSILineToken{
				{
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▄"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▄"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▄"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▄"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[30m", BG: "\x1b[43m", T: "▐"},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[30m", BG: "\x1b[44m", T: "█"},
					{FG: "\x1b[30m", BG: "\x1b[44m", T: "█"},
					{FG: "\x1b[30m", BG: "\x1b[44m", T: "█"},
					{FG: "\x1b[30m", BG: "\x1b[44m", T: "█"},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[30m", BG: "\x1b[43m", T: "▓"},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▌"},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[30m", BG: "\x1b[43m", T: "▓"},
					{FG: "\x1b[30m", BG: "\x1b[43m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[37m", BG: "\x1b[40m", T: " "},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "▐"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
					{FG: "\x1b[31m", BG: "\x1b[40m", T: "█"},
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

func TestTokeniseANSIFile(t *testing.T) {
	testCases := []struct {
		name     string
		filepath string
		expected [][]convert.ANSILineToken
	}{
		{
			name:     "File with empty lines",
			filepath: "../data/bhe-peaceofmind.txt",
			expected: [][]convert.ANSILineToken{
				{{FG: "", BG: "", T: ""}},
				{{FG: "", BG: "", T: ""}},
				{{FG: "", BG: "", T: "                              .                ."}},
				{{FG: "", BG: "", T: "                              .       /\\       ."}},
				{{FG: "", BG: "", T: "                              :     _/  \\_     :"}},
				{{FG: "", BG: "", T: "                         ___  |____/  __  \\____|  ___"}},
				{{FG: "", BG: "", T: "       _____ ___ ____ __/  \\\\ // .___/  \\___. \\\\ //  \\__ ____ ___ _____"}},
				{{FG: "", BG: "", T: "      / _ __\\\\__\\\\___\\\\_.   \\\\/ <|_ _ /\\ _ _|> \\//   ._//___//__//__ _ \\"}},
				{{FG: "", BG: "", T: "   ,_/   _______________| __ '    /Y \\\\// Y\\    ' __ |______________    \\_,"}},
				{{FG: "", BG: "", T: " _ __ ___\\==============|/  \\____/ .  \\/  . \\____/  \\|=============/___ __ _"}},
				{{FG: "", BG: "", T: " . _____________________________.  _________   _______    _______________. ."}},
				{{FG: "", BG: "", T: " . \\_______      \\_________     |__\\___     \\_/   ___/____\\_________     | ."}},
				{{FG: "", BG: "", T: " :  | .  |/       / .   | /     |.  ___      | .  \\|       | .   | /     | :"}},
				{{FG: "", BG: "", T: " |  | .  ________/| .   |/      |.  \\ |      | .   |       | .   |/      | |"}},
				{{FG: "", BG: "", T: " |  | :: \\      | | ::  /_______|::  \\|      | ::  '       | ::  /_______| |"}},
				{{FG: "", BG: "", T: " |  | :::.\\     | | :: _      | :::.  \\      | :::.        | :: _      |   |"}},
				{{FG: "", BG: "", T: " |  |______\\    | |_____\\     |________\\     |___          |_____\\     |   |"}},
				{{FG: "", BG: "", T: " :     .  . \\___| _ _____\\____|         \\____|   \\_________|   .  \\____|   :"}},
				{{FG: "", BG: "", T: " _ ___  \\  \\/  __  _ ___/                           \\___ _  __  \\/  .  ___ _"}},
				{{FG: "", BG: "", T: " .._  \\  \\    .\\/. \\\\\\_      _____________________    _/// .\\/.    /  /  _.."}},
				{{FG: "", BG: "", T: " || \\__\\__)  _|  |_   /    ./         \\_         /    \\   _|  |_  (__/__/ ||"}},
				{{FG: "", BG: "", T: " :| (___\\   _\\    /_  \\.   |    ___    |     ___/___ ./  _\\    /_   /___) |:"}},
				{{FG: "", BG: "", T: " . \\ \\   \\  \\  /\\  /  //   | .  | /    | .         / \\\\  \\  /\\  /  /   / / ."}},
				{{FG: "", BG: "", T: " . \\\\ \\__ \\  \\ \\/ /  /.    | .  |/     | .    ____/   .\\  \\ \\/ /  / __/ // ."}},
				{{FG: "", BG: "", T: " |\\ '\\__/ /  /    \\  \\     | :: '      | ::      |     /  /    \\  \\ \\__/' /|"}},
				{{FG: "", BG: "", T: " ' _ ____/  /  /\\  \\  '    | :::.  ____| ::: ____|    '  /  /\\  \\  \\____ _ '"}},
				{{FG: "", BG: "", T: ". _________/  /  \\  \\_ __  |______/    |____/       __ _/  /  \\  \\__________ ."}},
				{{FG: "", BG: "", T: " \\\\  ________/    '     /                           \\     '    \\_________  //"}},
				{{FG: "", BG: "", T: " // / ___   ______   ________ ________ ______  ._________________    ___ \\ \\\\"}},
				{{FG: "", BG: "", T: " \\\\ \\ \\/  ./      \\_/      \\ \\\\       \\_     \\_|     \\________   \\_.  \\/ / //"}},
				{{FG: "", BG: "", T: "  '\\ \\    | .    _   _      \\__________/.   .  |     ./.   | /     |    / /'"}},
				{{FG: "", BG: "", T: "   / /___ | .     \\_/        | .      | .   |  |     | .   |/      | ___\\ \\"}},
				{{FG: "", BG: "", T: "  /___  / | ::     |         | ::     | ::  |  |     | ::  |       | \\  ___\\"}},
				{{FG: "", BG: "", T: "  __ / /  | :::.   |         | :::.   | ::: |  '     | ::__|       |  \\ \\ __"}},
				{{FG: "", BG: "", T: " .\\_ \\ \\  |________|      ___|____    |_____|        |___\\ '      _|  / / _/."}},
				{{FG: "", BG: "", T: "  _ __\\ \\___ __ _  |_____/        \\___|     |________|    \\______/  _/ /__ _"}},
				{{FG: "", BG: "", T: " (_\\\\__   _                                                        _   __//_)"}},
				{{FG: "", BG: "", T: "    '(_\\  \\  - -- ---+-[ iT's hARd tO fiND yOUr pEACe ]-+--- -- -  /  /_)'"}},
				{{FG: "", BG: "", T: "       '\\__\\________________   ._.    __    ._.   ________________/__/'"}},
				{{FG: "", BG: "", T: "         bHe'===============\\_.  |___ \\/ ___|  ._/==============='sE!"}},
				{{FG: "", BG: "", T: "                   '----------|____  \\__/  ____|----------'"}},
				{{FG: "", BG: "", T: "                              :    \\_    _/    :"}},
				{{FG: "", BG: "", T: "                              :      \\  /      :"}},
				{{FG: "", BG: "", T: "                              .       \\/       ."}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, data, err := convert.ParseSAUCEFromFile(tc.filepath)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tc.filepath, err)
			}
			result := convert.TokeniseANSIString(string(data))
			test.PrintANSITestResults(data, tc.expected, result, t)
			test.Assert(tc.expected, result, t)
		})
	}
}
