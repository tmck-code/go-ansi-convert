package parse

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/log"
	"golang.org/x/text/encoding/charmap"
)

var (
	// '░', '▒', '█', '▄', '▐', '▀'
	blockChars = [][]byte{[]byte("░"), []byte("▒"), []byte("█"), []byte("▄"), []byte("▐"), []byte("▀")}
)

// DetectEncoding attempts to determine if file data is CP437 or ISO-8859-1
// Uses chardet library for detection, with fallback heuristics for ANSI art
func DetectEncoding(data []byte) string {
	pointsForCP437, pointsForISO := 0, 0

	// first, try CP437
	decoder := charmap.CodePage437.NewDecoder()
	CP437Translated, _ := decoder.Bytes(data)

	// second, try ISO-8859-1, then re-encode as CP437
	decoder = charmap.ISO8859_1.NewDecoder()
	ISOTranslated, _ := decoder.Bytes(data)
	if utf8.Valid(data) {
		log.DebugFprintln("encoding is UTF-8:")
		return "utf-8"
	}

	// having more of these chars is usually bad
	for _, ch := range [][]byte{[]byte("»"), []byte("Ü"), []byte("╖")} {
		cp437Count := bytes.Count(CP437Translated, ch)
		isoCount := bytes.Count(ISOTranslated, ch)

		if cp437Count > isoCount {
			log.DebugFprintf("  - char %q: \x1b[91mcp437=%d\x1b[0m, iso=%d\n", ch, cp437Count, isoCount)
			pointsForISO += 1
		} else if isoCount > cp437Count {
			log.DebugFprintf("  - char %q: cp437=%d, \x1b[91miso=%d\x1b[0m\n", ch, cp437Count, isoCount)
			pointsForCP437 += 1
		}
	}

	blockCharCounts := make(map[string]int)
	for _, ch := range blockChars {
		cp437Count := bytes.Count(CP437Translated, ch)
		if cp437Count > 0 {
			blockCharCounts[string(ch)] = cp437Count
		}
	}
	log.DebugFprintf("  - detected %d different block characters: %+v\n", len(blockCharCounts), blockCharCounts)
	if len(blockCharCounts) > 1 {
		pointsForCP437 += len(blockCharCounts) + 1
	} else if len(blockCharCounts) == 0 {
		pointsForISO += 3
	}

	// having more of these chars is usually good
	for _, ch := range [][]byte{[]byte("█"), []byte("¯"), []byte("░"), []byte("┌")} {
		cp437Count := bytes.Count(CP437Translated, ch)
		isoCount := bytes.Count(ISOTranslated, ch)

		if cp437Count > isoCount {
			log.DebugFprintf("  - char %q: \x1b[92mcp437=%d\x1b[0m, iso=%d\n", ch, cp437Count, isoCount)
			pointsForCP437 += 1
		} else if isoCount > cp437Count {
			log.DebugFprintf("  - char %q: cp437=%d, \x1b[92miso=%d\x1b[0m\n", ch, cp437Count, isoCount)
			pointsForISO += 1
		}
	}

	if pointsForCP437 > pointsForISO {
		log.DebugFprintf("Points:\n- \x1b[91mCP437\x1b[0m:      %d\n- ISO-8859-1: %d\n", pointsForCP437, pointsForISO)
		return "cp437"
	} else if pointsForISO > pointsForCP437 {
		log.DebugFprintf("Points:\n- CP437:      %d\n- \x1b[91mISO-8859-1\x1b[0m: %d\n", pointsForCP437, pointsForISO)
		return "iso-8859-1"
	}
	log.DebugFprintln("Unable to detect encoding, assuming ASCII")
	return "ascii"
}

func DecodeFileContents(data []byte, encoding string) (string, error) {
	var charMap *charmap.Charmap

	switch encoding {
	case "cp437":
		charMap = charmap.CodePage437
	case "iso-8859-1":
		charMap = charmap.ISO8859_1
	case "ascii", "utf-8":
		return string(data), nil
	default:
		return "", fmt.Errorf("unknown encoding: %s", encoding)
	}

	reader := charMap.NewDecoder().Reader(bytes.NewReader(data))
	decodedData, err := io.ReadAll(reader)
	if err != nil {
		log.DebugFprintln("Error decoding data:", err)
		return "", err
	}
	return string(decodedData), nil
}

// UnicodeStringLength calculates the display length of a string, accounting for:
// - Unicode characters that are double-width (e.g., CJK characters, emojis)
// - ANSI escape codes (which don't contribute to display width)
// Returns the total display width of the string.
func UnicodeStringLength(s string) int {
	nRunes, totalLen, ansiCode := len(s), 0, false

	for i, r := range s {
		if i < nRunes-1 {
			// detect the beginning of an ANSI escape code
			// e.g. "\x1b[38;5;196m"
			//       ^^^ start    ^ end
			if s[i:i+2] == "\x1b[" {
				ansiCode = true
			}
		}
		if ansiCode {
			// detect the end of an ANSI escape code
			if r == 'm' {
				ansiCode = false
			}
		} else {
			if r < 128 {
				// if ascii, then use width of 1. this saves some time
				totalLen++
			} else {
				totalLen += runewidth.RuneWidth(r)
			}
		}
	}
	return totalLen
}

func UnicodeLineLengths(lines []string) []int {
	lengths := make([]int, len(lines))
	for i, line := range lines {
		lengths[i] = UnicodeStringLength(line)
	}
	return lengths
}

// LongestUnicodeLineLength finds the maximum display length among all lines.
func LongestUnicodeLineLength(lines []string) int {
	maxLen := 0
	for _, line := range UnicodeLineLengths(lines) {
		if line > maxLen {
			maxLen = line
		}
	}
	return maxLen
}
