package test

import (
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/src/parse"
)

// Benchmark tests -------------------------------------------------------------

func BenchmarkUnicodeStringLength(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{"ASCII", "Hello World"},
		{"Box drawing", "┌─────┐│     │└─────┘"},
		{"Japanese", "こんにちは日本語"},
		{"ANSI colored ASCII", "\x1b[38;5;129mHello\x1b[0m"},
		{"Complex ANSI mixed", "\x1b[38;5;129mHello\x1b[0m \x1b[38;5;160m世界\x1b[0m \x1b[38;5;46m█▀\x1b[0m"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				parse.UnicodeStringLength(tc.input)
			}
		})
	}
}

func BenchmarkSanitiseUnicodeString(b *testing.B) {
	testCases := []struct {
		name    string
		input   string
		justify bool
	}{
		{"Single line no ANSI", "Hello World", false},
		{"Single line with ANSI", "\x1b[38;5;129mColored text", false},
		{"Multi-line with ANSI", "\x1b[38;5;129mLine 1\nLine 2", false},
		{"Multi-line with justify", "Short\nLonger line\nMed", true},
		{"Complex multi-line", "\x1b[38;5;160m▄\x1b[38;5;46m▄\n▄\x1b[38;5;190m▄", false},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				convert.SanitiseUnicodeString(tc.input, tc.justify)
			}
		})
	}
}

func BenchmarkTokeniseANSIString(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Single line no color", "         ▄▄          ▄▄"},
		{"Single line FG+BG", "\x1b[38;5;129mAAA\x1b[48;5;160mXX"},
		{"Multi-line", "\x1b[38;5;160m▄\x1b[38;5;46m▄\n▄\x1b[38;5;190m▄"},
		{"Complex egg", "    \x1b[49m   \x1b[38;5;16m▄▄\x1b[48;5;16m\x1b[38;5;142m▄▄▄\x1b[49m\x1b[38;5;16m▄▄\n     ▄\x1b[48;5;16m\x1b[38;5;58m▄\x1b[48;5;58m\x1b[38;5;70m▄\x1b[48;5;70m \x1b[48;5;227m    \x1b[48;5;237m\x1b[38;5;227m▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[49m\x1b[38;5;16m▄"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				convert.TokeniseANSIString(tc.input)
			}
		})
	}
}

func BenchmarkMirrorHorizontally(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Basic no color", "         ▄▄          ▄▄"},
		{"With trailing spaces", "         ▄▄          ▄▄      "},
		{"Japanese characters", "こんにちは日本語"},
		{"Box drawing", "┌─────┐│     │└─────┘"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				convert.MirrorHorizontally(tc.input)
			}
		})
	}
}

func BenchmarkFlipHorizontal(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Single line ANSI", "\x1b[38;5;129mAAA \x1b[48;5;160m XX \x1b[0m"},
		{"Multi-line ANSI", "\x1b[38;5;129mAAA    \x1b[48;5;160m XX \x1b[0m"},
		{
			"Pikachu no color",
			strings.Join([]string{
				"         ▄▄          ▄▄",
				"        ▄▄▄     ▄▄▄▄▄▄ ▄▄",
				"       ▄  ▄▀ ▄▄▄  ▄▄   ▄▀",
				"     ▄▄▄   ▄▄  ▄▄    ▄▀",
				"    ▄▄   ▄▄▄  ▄ ▀▄  ▄▄",
				"    ▀▄▄   ▄▄▄   ▄▄▄ ▄▀",
				" ▀▄▄▄▄▄ ▄▄   ▄▄▄▄▄",
				"           ▄▄▄▄  ▄▄▀",
				"       ▀▄▄▄    ▄▀",
				"           ▀▀▄▀",
			}, "\n"),
		},
		{
			"Pikachu with color",
			"    \x1b[49m     \x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;232m▄ \x1b[49m         \x1b[38;5;16m▄▄\n        ▄\x1b[48;5;16m\x1b[38;5;94m▄\x1b[48;5;232m▄\x1b[48;5;16m \x1b[49m    \x1b[38;5;16m▄▄▄▄\x1b[48;5;16m\x1b[38;5;214m▄\x1b[48;5;214m\x1b[38;5;94m▄\x1b[48;5;94m \x1b[48;5;16m▄\x1b[49m\x1b[38;5;16m▄\n       ▄\x1b[48;5;16m \x1b[48;5;94m \x1b[48;5;58m▄\x1b[49m▀ ▄\x1b[48;5;16m\x1b[38;5;214m▄▄\x1b[48;5;232m  \x1b[38;5;94m▄\x1b[48;5;214m▄\x1b[48;5;94m   \x1b[38;5;16m▄\x1b[49m▀",
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tokens := convert.TokeniseANSIString(tc.input)
				convert.FlipHorizontal(tokens)
			}
		})
	}
}
