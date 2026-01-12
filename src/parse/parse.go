package parse

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mattn/go-runewidth"
	"github.com/tmck-code/go-ansi-convert/src"
	"golang.org/x/text/encoding/charmap"
)

// DetectEncoding attempts to determine if file data is CP437 or ISO-8859-1
// Uses chardet library for detection, with fallback heuristics for ANSI art
func DetectEncoding(data []byte) string {
	pointsForCP437, pointsForISO := 0, 0

	// first, try CP437
	decoder := charmap.CodePage437.NewDecoder()
	CP437Translated, err := decoder.Bytes(data)
	if err != nil {
		pointsForISO += 1
	} else {
		// now try encoding the cp437 data as ISO-8859-1
		encoder := charmap.ISO8859_1.NewEncoder()
		_, err = encoder.Bytes(CP437Translated)
		// if that failed, then it's definitely cp437
		if err != nil {
			src.DebugFprintln("  - translation from CP437 to ISO-8859-1 failed!!")
			pointsForCP437 += 1
		}
	}

	// second, try ISO-8859-1, then re-encode as CP437
	decoder = charmap.ISO8859_1.NewDecoder()
	ISOTranslated, err := decoder.Bytes(data)
	if err != nil {
		pointsForCP437 += 1
	} else {
		encoder := charmap.CodePage437.NewEncoder()
		_, err = encoder.Bytes(ISOTranslated)
		// if that failed, then it's definitely iso-8859-1
		if err != nil {
			src.DebugFprintln("  - translation from ISO-8859-1 to CP437 failed!!")
			pointsForISO += 1
		}
	}

	// having more of these chars is usually bad
	for _, ch := range [][]byte{[]byte("»"), []byte("Ü"), []byte("╖")} {
		origCount := bytes.Count(data, ch)
		cp437Count := bytes.Count(CP437Translated, ch)
		isoCount := bytes.Count(ISOTranslated, ch)

		if cp437Count > isoCount {
			src.DebugFprintf("  - char %q: orig=%d, \x1b[91mcp437=%d\x1b[0m, iso=%d\n", ch, origCount, cp437Count, isoCount)
			pointsForISO += 1
		} else if isoCount > cp437Count {
			src.DebugFprintf("  - char %q: orig=%d, cp437=%d, \x1b[91miso=%d\x1b[0m\n", ch, origCount, cp437Count, isoCount)
			pointsForCP437 += 1
		} else {
			src.DebugFprintf("  - char %q: orig=%d, cp437=%d, iso=%d\n", ch, origCount, cp437Count, isoCount)
		}
	}
	// having more of these chars is usually good
	// other candidates: []byte("²")
	for _, ch := range [][]byte{[]byte("█"), []byte("¯"), []byte("░"), []byte("┌")} {
		origCount := bytes.Count(data, ch)
		cp437Count := bytes.Count(CP437Translated, ch)
		isoCount := bytes.Count(ISOTranslated, ch)

		if cp437Count > isoCount {
			src.DebugFprintf("  - char %q: orig=%d, \x1b[92mcp437=%d\x1b[0m, iso=%d\n", ch, origCount, cp437Count, isoCount)
			pointsForCP437 += 1
		} else if isoCount > cp437Count {
			src.DebugFprintf("  - char %q: orig=%d, cp437=%d, \x1b[92miso=%d\x1b[0m\n", ch, origCount, cp437Count, isoCount)
			pointsForISO += 1
		} else {
			src.DebugFprintf("  - char %q: orig=%d, cp437=%d, iso=%d\n", ch, origCount, cp437Count, isoCount)
		}
	}
	src.DebugFprintf("Points:\n- CP437:      %d\n- ISO-8859-1: %d\n", pointsForCP437, pointsForISO)

	if pointsForCP437 > pointsForISO {
		return "cp437"
	} else if pointsForISO > pointsForCP437 {
		return "iso-8859-1"
	}
	src.DebugFprintln("Unable to detect encoding, assuming ASCII")
	return "ascii"
}

func DecodeFileContents(data []byte, encoding string) (string, error) {
	var charMap *charmap.Charmap

	switch encoding {
	case "cp437":
		charMap = charmap.CodePage437
	case "iso-8859-1":
		charMap = charmap.ISO8859_1
	case "ascii":
		return string(data), nil
	default:
		return "", fmt.Errorf("unknown encoding: %s", encoding)
	}

	reader := charMap.NewDecoder().Reader(bytes.NewReader(data))
	decodedData, err := io.ReadAll(reader)
	if err != nil {
		src.DebugFprintln("Error decoding data:", err)
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
