package test

import (
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
)

func TestFlipVertical(t *testing.T) {
    testCases := []struct {
        name     string
        input    []string
        expected []string
    }{
        {
            name: "ASCII vertical flip",
            input: []string{
                "VW",
                "AB",
            },
            expected: []string{
                "ꓯB",
                "ꓥM",
            },
        },
        {
            // "█▀▀▀▐▄▐█"
            // "█▄▄▄▐▄▐█"
            name: "Unicode vertical flip",
            input: []string{
                "█▀▀▀▐▄▐█", 
            },
			expected: []string{
				"█▄▄▄▐▀▐█",
			},
        },
        {
            name: "Unicode multi-line vertical flip",
            input: []string{
                "█████████████████▌▀▀█████",
                "██████▀██████▀ ▀▀▐▄ ▐████",
                "████▌▀▀▐█▌▄▌ ░▓▓ ▓▓▓ ████",
                "██▌▀▀▀▀▀▀▌█▌ ▐ ▀▄▓▐▄▀▀███",
            },
            expected: []string{
                "██▌▄▄▄▄▄▄▌█▌ ▐ ▄▀▓▐▀▄▄███",
                "████▌▄▄▐█▌▀▌ ░▓▓ ▓▓▓ ████",
                "██████▄██████▄ ▄▄▐▀ ▐████",
                "█████████████████▌▄▄█████",
            },
        },
    }
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            input := strings.Join(tc.input, "\n")
            result := convert.FlipVertical(input)
            expected := [][]convert.ANSILineToken{}
            for _, line := range tc.expected {
                expected = append(expected, []convert.ANSILineToken{{FG: "", BG: "", T: line}})
            }

			PrintANSITestResults(input, expected, result, t)
            Assert(expected, result, t)
        })
    }
}
