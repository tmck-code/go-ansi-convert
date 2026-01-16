package test

import (
	"os"
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/convert"
	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/parse"
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
					convert.ANSILineToken{FG: "", BG: "\x1b[49m", T: "    "},
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
					{FG: "\x1b[30m", BG: "\x1b[48;2;168;48;76m", T: "â"},
					{FG: "\x1b[38;2;144;16;64m", BG: "\x1b[40m", T: "â\"¯¯\""},
					{FG: "\x1b[30m", BG: "\x1b[48;2;168;48;76m", T: "x"},
				},
			},
		},
		{
			name: "tokenise multiline 2",
			input: strings.Join(
				[]string{
					"    |    \x1b[30m\x1b[1;87;87;87t| \x1b[36m\x1b[1;87;255;255t|\x1b[0m\x1b[0;0;0;0t\x1b[1;171;171;171t\x1b[11C\x1b[1;36m\x1b[1;87;255;255t____/ \x1b[30m\x1b[1;87;87;87t/¯.   ¯\\__________________.'  !·\x1b[0;35m\x1b[0;0;0;0t\x1b[1;171;0;171t.:\x1b[33m\x1b[1;171;87;0tY \x1b[1;31m\x1b[1;255;87;87tY   \x1b[0;33m\x1b[0;0;0;0t\x1b[1;171;87;0tj ! \x1b[1;31m\x1b[1;255;87;87t| \x1b[0;35m\x1b[0;0;0;0t\x1b[1;171;0;171t:\x1b[1;33m\x1b[1;255;255;87t|",
					"    |\x1b[41m\x1b[0;171;0;0t  \x1b[0;31m\x1b[0;0;0;0t\x1b[1;171;0;0t||\x1b[1;30m\x1b[1;87;87;87t| \x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t!_\x1b[1m\x1b[1;87;255;255t_\x1b[0;36m\x1b[0;0;0;0t\x1b[1;0;171;171t_\x1b[1m\x1b[1;87;255;255t_.-----· \x1b[30m\x1b[1;87;87;87t____/ \x1b[0;33m\x1b[0;0;0;0t\x1b[1;171;87;0t/\x1b[37m\x1b[1;171;171;171t\x1b[5C\x1b[1;30m\x1b[1;87;87;87t.·--------------------·\x1b[0;35m\x1b[0;0;0;0t\x1b[1;171;0;171t.:::\x1b[33m\x1b[1;171;87;0t' \x1b[1;31m\x1b[1;255;87;87t`. \x1b[0;33m\x1b[0;0;0;0t\x1b[1;171;87;0t~ / \x1b[1;31m\x1b[1;255;87;87t.'\x1b[0;35m\x1b[0;0;0;0t\x1b[1;171;0;171t.:\x1b[1;33m\x1b[1;255;255;87t|\r\n",
				},
				"\n",
			),
			expected: [][]convert.ANSILineToken{
				{
					{FG: "", BG: "", T: "    |    "},
					{FG: "\x1b[38;2;87;87;87m", BG: "", T: "| "},
					{FG: "\x1b[38;2;87;255;255m", BG: "\x1b[49m", T: "|"},
					{FG: "\x1b[38;2;171;171;171m", BG: "\x1b[40m", T: "           "},
					{FG: "\x1b[38;2;87;255;255m", BG: "\x1b[40m", T: "____/ "},
					{FG: "\x1b[38;2;87;87;87m", BG: "\x1b[40m", T: "/¯.   ¯\\__________________.'  !·"},
					{FG: "\x1b[38;2;171;0;171m", BG: "\x1b[40m", T: ".:"},
					{FG: "\x1b[38;2;171;87;0m", BG: "\x1b[40m", T: "Y "},
					{FG: "\x1b[38;2;255;87;87m", BG: "\x1b[40m", T: "Y   "},
					{FG: "\x1b[38;2;171;87;0m", BG: "\x1b[40m", T: "j ! "},
					{FG: "\x1b[38;2;255;87;87m", BG: "\x1b[40m", T: "| "},
					{FG: "\x1b[38;2;171;0;171m", BG: "\x1b[40m", T: ":"},
					{FG: "\x1b[38;2;255;255;87m", BG: "\x1b[40m", T: "|"},
				},
				{
					{FG: "\x1b[38;2;255;255;87m", BG: "\x1b[40m", T: "    |"},
					{FG: "\x1b[38;2;255;255;87m", BG: "\x1b[48;2;171;0;0m", T: "  "},
					{FG: "\x1b[38;2;171;0;0m", BG: "\x1b[40m", T: "||"},
					{FG: "\x1b[38;2;87;87;87m", BG: "\x1b[40m", T: "| "},
					{FG: "\x1b[38;2;0;171;171m", BG: "\x1b[40m", T: "!_"},
					{FG: "\x1b[1m\x1b[38;2;87;255;255m", BG: "\x1b[40m", T: "_"},
					{FG: "\x1b[1m\x1b[38;2;0;171;171m", BG: "\x1b[40m", T: "_"},
					{FG: "\x1b[1m\x1b[38;2;87;255;255m", BG: "\x1b[40m", T: "_.-----· "},
					{FG: "\x1b[1m\x1b[38;2;87;87;87m", BG: "\x1b[40m", T: "____/ "},
					{FG: "\x1b[1m\x1b[38;2;171;87;0m", BG: "\x1b[40m", T: "/"},
					{FG: "\x1b[1m\x1b[38;2;171;171;171m", BG: "\x1b[40m", T: "     "},
					{FG: "\x1b[1m\x1b[38;2;87;87;87m", BG: "\x1b[40m", T: ".·--------------------·"},
					{FG: "\x1b[1m\x1b[38;2;171;0;171m", BG: "\x1b[40m", T: ".:::"},
					{FG: "\x1b[1m\x1b[38;2;171;87;0m", BG: "\x1b[40m", T: "' "},
					{FG: "\x1b[1m\x1b[38;2;255;87;87m", BG: "\x1b[40m", T: "`. "},
					{FG: "\x1b[1m\x1b[38;2;171;87;0m", BG: "\x1b[40m", T: "~ / "},
					{FG: "\x1b[1m\x1b[38;2;255;87;87m", BG: "\x1b[40m", T: ".'"},
					{FG: "\x1b[1m\x1b[38;2;171;0;171m", BG: "\x1b[40m", T: ".:"},
					{FG: "\x1b[1m\x1b[38;2;255;255;87m", BG: "\x1b[40m", T: "|"},
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
				{},
				{},
				{{FG: "", BG: "", Control: "", T: "                              .                ."}},
				{{FG: "", BG: "", Control: "", T: "                              .       /\\       ."}},
				{{FG: "", BG: "", Control: "", T: "                              :     _/  \\_     :"}},
				{{FG: "", BG: "", Control: "", T: "                         ___  |____/  __  \\____|  ___"}},
				{{FG: "", BG: "", Control: "", T: "       _____ ___ ____ __/  \\\\ // .___/  \\___. \\\\ //  \\__ ____ ___ _____"}},
				{{FG: "", BG: "", Control: "", T: "      / _ __\\\\__\\\\___\\\\_.   \\\\/ <|_ _ /\\ _ _|> \\//   ._//___//__//__ _ \\"}},
				{{FG: "", BG: "", Control: "", T: "   ,_/   _______________| __ '    /Y \\\\// Y\\    ' __ |______________    \\_,"}},
				{{FG: "", BG: "", Control: "", T: " _ __ ___\\==============|/  \\____/ .  \\/  . \\____/  \\|=============/___ __ _"}},
				{{FG: "", BG: "", Control: "", T: " . _____________________________.  _________   _______    _______________. ."}},
				{{FG: "", BG: "", Control: "", T: " . \\_______      \\_________     |__\\___     \\_/   ___/____\\_________     | ."}},
				{{FG: "", BG: "", Control: "", T: " :  | .  |/       / .   | /     |.  ___      | .  \\|       | .   | /     | :"}},
				{{FG: "", BG: "", Control: "", T: " |  | .  ________/| .   |/      |.  \\ |      | .   |       | .   |/      | |"}},
				{{FG: "", BG: "", Control: "", T: " |  | :: \\      | | ::  /_______|::  \\|      | ::  '       | ::  /_______| |"}},
				{{FG: "", BG: "", Control: "", T: " |  | :::.\\     | | :: _      | :::.  \\      | :::.        | :: _      |   |"}},
				{{FG: "", BG: "", Control: "", T: " |  |______\\    | |_____\\     |________\\     |___          |_____\\     |   |"}},
				{{FG: "", BG: "", Control: "", T: " :     .  . \\___| _ _____\\____|         \\____|   \\_________|   .  \\____|   :"}},
				{{FG: "", BG: "", Control: "", T: " _ ___  \\  \\/  __  _ ___/                           \\___ _  __  \\/  .  ___ _"}},
				{{FG: "", BG: "", Control: "", T: " .._  \\  \\    .\\/. \\\\\\_      _____________________    _/// .\\/.    /  /  _.."}},
				{{FG: "", BG: "", Control: "", T: " || \\__\\__)  _|  |_   /    ./         \\_         /    \\   _|  |_  (__/__/ ||"}},
				{{FG: "", BG: "", Control: "", T: " :| (___\\   _\\    /_  \\.   |    ___    |     ___/___ ./  _\\    /_   /___) |:"}},
				{{FG: "", BG: "", Control: "", T: " . \\ \\   \\  \\  /\\  /  //   | .  | /    | .         / \\\\  \\  /\\  /  /   / / ."}},
				{{FG: "", BG: "", Control: "", T: " . \\\\ \\__ \\  \\ \\/ /  /.    | .  |/     | .    ____/   .\\  \\ \\/ /  / __/ // ."}},
				{{FG: "", BG: "", Control: "", T: " |\\ '\\__/ /  /    \\  \\     | :: '      | ::      |     /  /    \\  \\ \\__/' /|"}},
				{{FG: "", BG: "", Control: "", T: " ' _ ____/  /  /\\  \\  '    | :::.  ____| ::: ____|    '  /  /\\  \\  \\____ _ '"}},
				{{FG: "", BG: "", Control: "", T: ". _________/  /  \\  \\_ __  |______/    |____/       __ _/  /  \\  \\__________ ."}},
				{{FG: "", BG: "", Control: "", T: " \\\\  ________/    '     /                           \\     '    \\_________  //"}},
				{{FG: "", BG: "", Control: "", T: " // / ___   ______   ________ ________ ______  ._________________    ___ \\ \\\\"}},
				{{FG: "", BG: "", Control: "", T: " \\\\ \\ \\/  ./      \\_/      \\ \\\\       \\_     \\_|     \\________   \\_.  \\/ / //"}},
				{{FG: "", BG: "", Control: "", T: "  '\\ \\    | .    _   _      \\__________/.   .  |     ./.   | /     |    / /'"}},
				{{FG: "", BG: "", Control: "", T: "   / /___ | .     \\_/        | .      | .   |  |     | .   |/      | ___\\ \\"}},
				{{FG: "", BG: "", Control: "", T: "  /___  / | ::     |         | ::     | ::  |  |     | ::  |       | \\  ___\\"}},
				{{FG: "", BG: "", Control: "", T: "  __ / /  | :::.   |         | :::.   | ::: |  '     | ::__|       |  \\ \\ __"}},
				{{FG: "", BG: "", Control: "", T: " .\\_ \\ \\  |________|      ___|____    |_____|        |___\\ '      _|  / / _/."}},
				{{FG: "", BG: "", Control: "", T: "  _ __\\ \\___ __ _  |_____/        \\___|     |________|    \\______/  _/ /__ _"}},
				{{FG: "", BG: "", Control: "", T: " (_\\\\__   _                                                        _   __//_)"}},
				{{FG: "", BG: "", Control: "", T: "    '(_\\  \\  - -- ---+-[ iT's hARd tO fiND yOUr pEACe ]-+--- -- -  /  /_)'"}},
				{{FG: "", BG: "", Control: "", T: "       '\\__\\________________   ._.    __    ._.   ________________/__/'"}},
				{{FG: "", BG: "", Control: "", T: "         bHe'===============\\_.  |___ \\/ ___|  ._/==============='sE!"}},
				{{FG: "", BG: "", Control: "", T: "                   '----------|____  \\__/  ____|----------'"}},
				{{FG: "", BG: "", Control: "", T: "                              :    \\_    _/    :"}},
				{{FG: "", BG: "", Control: "", T: "                              :      \\  /      :"}},
				{{FG: "", BG: "", Control: "", T: "                              .       \\/       ."}},
			},
		},
		{
			name:     "File with erase codes",
			filepath: "../data/TSIENEXX.sample.ANS",
			expected: [][]convert.ANSILineToken{
				{
					{FG: "", BG: "", Control: "\x1b[2J", T: ""},
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
				{
					{FG: "", BG: "", Control: "\x1b[u", T: ""},
					{FG: "\x1b[32m", BG: "\x1b[40m", Control: "", T: "████"},
					{FG: "\x1b[33m", BG: "\x1b[42m", Control: "", T: "░▒▓███"},
					{FG: "\x1b[33m", BG: "\x1b[40m", Control: "", T: "▄▄▄▄▄▄                   "},
					{FG: "\x1b[36m", BG: "\x1b[40m", Control: "", T: "█"},
					{FG: "\x1b[36m", BG: "\x1b[47m", Control: "", T: "▒"},
					{FG: "\x1b[37m", BG: "\x1b[46m", Control: "", T: "▓"},
					{FG: "\x1b[36m", BG: "\x1b[47m", Control: "", T: "▒"},
					{FG: "\x1b[36m", BG: "\x1b[40m", Control: "", T: "█     "},
					{FG: "", BG: "", Control: "\x1b[s", T: ""},
				},
				{
					{FG: "", BG: "", Control: "\x1b[u", T: ""},
					{FG: "\x1b[31m", BG: "\x1b[40m", Control: "", T: "▐█▌ ▄███▒░"},
					{FG: "\x1b[30m", BG: "\x1b[41m", Control: "", T: "▒▓"},
					{FG: "\x1b[30m", BG: "\x1b[40m", Control: "", T: "████▄▄▄▄"},
					{FG: "\x1b[32m", BG: "\x1b[40m", Control: "", T: "██████████████████"},
					{FG: "", BG: "", Control: "\x1b[s", T: ""},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileData, err := os.ReadFile(tc.filepath)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tc.filepath, err)
			}
			encoding := parse.DetectEncoding(fileData)
			_, data, err := convert.SAUCERecord(fileData, encoding)
			if err != nil {
				t.Fatalf("Failed to read test file %s: %v", tc.filepath, err)
			}

			result := convert.TokeniseANSIString(data)
			test.PrintANSITestResults(data, tc.expected, result, t)
			test.Assert(tc.expected, result, t)
		})
	}
}
