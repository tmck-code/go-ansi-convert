package test

import (
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
)

func TestUnicodeStringLength(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "ASCII text",
			input:    "Hello World",
			expected: 11,
		},
		{
			name:     "Block/box drawing characters",
			input:    "┌─────┐│     │└─────┘",
			expected: 21,
		},
		{
			name:     "Full block characters",
			input:    "█▀▄▌▐░▒▓",
			expected: 8,
		},
		{
			name:     "Box drawing with double lines",
			input:    "╔═══╗║   ║╚═══╝",
			expected: 15,
		},
		{
			name:     "Japanese hiragana",
			input:    "こんにちは", // "hello" in hiragana
			expected: 10,      // 5 chars × 2 width each
		},
		{
			name:     "Japanese katakana",
			input:    "カタカナ", // katakana
			expected: 8,      // 4 chars × 2 width each
		},
		{
			name:     "Japanese kanji",
			input:    "日本語", // "Japanese language"
			expected: 6,     // 3 chars × 2 width each
		},
		{
			name:     "Mixed ASCII and Japanese",
			input:    "Hello世界", // "Hello world"
			expected: 9,         // 5 ASCII + 2 kanji (4 width)
		},
		{
			name:     "ANSI colored ASCII",
			input:    "\x1b[38;5;129mHello\x1b[0m",
			expected: 5, // ANSI codes don't count
		},
		{
			name:     "ANSI colored Japanese",
			input:    "\x1b[38;5;160m日本語\x1b[0m",
			expected: 6, // 3 kanji × 2 width, ANSI codes don't count
		},
		{
			name:     "ANSI colored box characters",
			input:    "\x1b[38;5;196m█\x1b[48;5;16m▀▄\x1b[0m",
			expected: 3, // 3 block chars, ANSI codes don't count
		},
		{
			name:     "Complex ANSI with mixed characters",
			input:    "\x1b[38;5;129mHello\x1b[0m \x1b[38;5;160m世界\x1b[0m \x1b[38;5;46m█▀\x1b[0m",
			expected: 13,
			// Hello 世界 █▀
			// "Hello " (6) + "世界 " (5) + "█▀" (2) = 13
		},
		{
			name:     "Multiple ANSI codes in sequence",
			input:    "\x1b[38;5;160m\x1b[48;5;16mTest\x1b[0m",
			expected: 4,
		},
		{
			name:     "Japanese with FG and BG colors",
			input:    "\x1b[38;5;129m\x1b[48;5;160mこんにちは\x1b[0m",
			expected: 10, // 5 hiragana chars × 2 width
		},
		{
			name:     "Box art",
			input:    " ▄  █ ▄███▄   █    █    ████▄       ▄ ▄   ████▄ █▄▄▄▄ █     ██▄",
			expected: 63,
		},
		{
			name:     "Mixed width Unicode with ANSI",
			input:    "\x1b[38;5;46m┌──┐\x1b[0m \x1b[38;5;160m日本\x1b[0m \x1b[38;5;129mABC\x1b[0m",
			expected: 13, // "┌──┐" (4) + " " (1) + "日本" (4) + " " (1) + "ABC" (3) = 13
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "Only ANSI codes",
			input:    "\x1b[38;5;129m\x1b[48;5;160m\x1b[0m",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convert.UnicodeStringLength(tc.input)
			if result != tc.expected {
				t.Errorf("UnicodeStringLength(%q) = %d; want %d", tc.input, result, tc.expected)
			}
		})
	}
}
