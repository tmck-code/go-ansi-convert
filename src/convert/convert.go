package convert

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

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

// SanitiseUnicodeString cleans up ANSI codes in a string and optionally justifies lines.
// It ensures all lines end with a reset code and pads lines to equal width if justifyLines is true.
func SanitiseUnicodeString(s string, justifyLines bool) string {
	if s == "" {
		return s
	}

	tokenizedLines := TokeniseANSIString(s)
	if len(tokenizedLines) == 0 {
		return s
	}

	// Calculate max length if justification is needed
	maxLen := 0
	if justifyLines {
		maxLen = LongestUnicodeLineLength(strings.Split(s, "\n"))
	}

	var sanitised strings.Builder

	for i, tokens := range tokenizedLines {
		var lineBuilder strings.Builder
		lineLen := 0
		hasReset := false
		for _, token := range tokens {
			lineBuilder.WriteString(token.FG)
			lineBuilder.WriteString(token.BG)
			lineBuilder.WriteString(token.T)
			lineLen += UnicodeStringLength(token.T)
			if token.FG == "\x1b[0m" {
				hasReset = true
			} else if token.FG != "" || token.BG != "" {
				hasReset = false
			}
		}
		// Ensure line ends with reset before padding, if not already present
		if !hasReset {
			lineBuilder.WriteString("\x1b[0m")
		}
		// Add padding if justification is enabled
		if justifyLines && lineLen < maxLen {
			lineBuilder.WriteString(strings.Repeat(" ", maxLen-lineLen))
		}
		if i < len(tokenizedLines)-1 {
			lineBuilder.WriteString("\n")
		}
		sanitised.WriteString(lineBuilder.String())
	}
	return sanitised.String()
}

// ANSILineToken represents a segment of text with its associated ANSI formatting.
// - FG is the foreground color code
// - BG is the background color code, and
// - T is the text content.
type ANSILineToken struct {
	FG string
	BG string
	T  string
}

func OptimiseANSITokens(lines [][]ANSILineToken) [][]ANSILineToken {
	var optimisedLines [][]ANSILineToken

	for _, tokens := range lines {
		if len(tokens) == 0 {
			// Preserve empty lines to maintain line structure
			optimisedLines = append(optimisedLines, []ANSILineToken{})
			continue
		}
		var optimisedTokens []ANSILineToken
		var lastFG, lastBG string

		for _, tok := range tokens {
			// Ignore empty reset tokens
			if tok.FG == "\x1b[0m" && tok.T == "" {
				continue
			}
			if len(optimisedTokens) > 0 && tok.FG == lastFG && tok.BG == lastBG {
				optimisedTokens[len(optimisedTokens)-1].T += tok.T
			} else {
				if tok.FG != lastFG && tok.BG == lastBG {
					// Only FG changed
					optimisedTokens = append(optimisedTokens, ANSILineToken{FG: tok.FG, BG: "", T: tok.T})
				} else if tok.BG != lastBG && tok.FG == lastFG {
					// Only BG changed
					optimisedTokens = append(optimisedTokens, ANSILineToken{FG: "", BG: tok.BG, T: tok.T})
				} else {
					// Both changed or both new
					optimisedTokens = append(optimisedTokens, tok)
				}
				lastFG, lastBG = tok.FG, tok.BG
			}
		}
		optimisedLines = append(optimisedLines, optimisedTokens)
	}
	return optimisedLines
}

// TokeniseANSIString parses a string containing ANSI escape codes into structured tokens.
// It splits the input by lines and then tokenizes each line, extracting foreground/background
// colors and text segments.
// Returns a 2D slice where each inner slice represents tokens for one line.
func TokeniseANSIString(msg string) [][]ANSILineToken {
	isColour := false
	isReset := false
	fg := ""
	bg := ""
	lines := make([][]ANSILineToken, 0)
	lineSlice := strings.Split(msg, "\n")

	for i, line := range lineSlice {
		if i == len(lineSlice)-1 && (len(line) == 0 || line == "\x1b[0m") {
			continue
		}

		tokens := make([]ANSILineToken, 0)
		text := ""
		colour := ""
		isReset = false // Clear reset state at start of each line

		for _, ch := range line {
			// start of colour sequence detected!
			if ch == '\033' {
				isColour = true
				// if there is text in the current token buffer,
				if text != "" {
					// if we are setting a bg colour, but the last token didn't have one
					// then add a background clear to the previous bg
					if bg != "" && len(tokens) > 0 && !strings.Contains(bg, "[49m") && tokens[len(tokens)-1].BG == "" {
						prevToken := tokens[len(tokens)-1]
						tokens[len(tokens)-1] = ANSILineToken{prevToken.FG, "\x1b[49m", prevToken.T}
					}
					if isReset {
						tokens = append(tokens, ANSILineToken{"\x1b[0m", "", text})
						isReset = false
					} else {
						tokens = append(tokens, ANSILineToken{fg, bg, text})
					}
					colour = string(ch)
					text = ""
				} else {
					colour = string(ch)
				}
			} else if isColour {
				// keep building the current ANSI escape code if \033 was found earlier
				// disable the isColour bool if the end of the ANSI escape code is found
				colour += string(ch)
				if ch == 'm' {
					isColour = false
					// Check for 256-color or true color codes first (contains ;5; or ;2;)
					if strings.Contains(colour, ";5;") || strings.Contains(colour, ";2;") {
						// 256-color or true color format
						if strings.Contains(colour, "[38") || strings.Contains(colour, "[39") {
							fg = colour
							isReset = false
						} else if strings.Contains(colour, "[48") || strings.Contains(colour, "[49") {
							bg = colour
							isReset = false
						}
					} else if strings.Contains(colour, ";") {
						// Check if this is a combined code with both FG and BG (e.g., \x1b[0;31;40m)
						parts := strings.Split(strings.TrimPrefix(strings.TrimSuffix(colour, "m"), "\x1b["), ";")
						hasFG := false
						hasBG := false
						fgCode := ""
						bgCode := ""
						hasReset := false

						for _, part := range parts {
							if part == "0" {
								hasReset = true
							} else if len(part) >= 2 {
								// Check for FG codes (30-37, 90-97)
								if (part[0] == '3' && part[1] >= '0' && part[1] <= '7') ||
									(part[0] == '9' && part[1] >= '0' && part[1] <= '7') {
									hasFG = true
									fgCode = part
								}
								// Check for BG codes (40-47, 100-107)
								if (part[0] == '4' && part[1] >= '0' && part[1] <= '7') ||
									(len(part) == 3 && part[0] == '1' && part[1] == '0' && part[2] >= '0' && part[2] <= '7') {
									hasBG = true
									bgCode = part
								}
							}
						}

						// If it has both FG and BG codes, split them into separate codes
						if hasFG && hasBG {
							fg = "\x1b[" + fgCode + "m"
							bg = "\x1b[" + bgCode + "m"
							isReset = false
						} else if hasFG {
							fg = colour
							isReset = false
						} else if hasBG {
							bg = colour
							isReset = false
						} else if hasReset {
							isReset = true
							fg = ""
							bg = ""
							colour = ""
						}
					} else if strings.Contains(colour, "[3") || strings.Contains(colour, "[9") && (colour[len(colour)-2] >= '0' && colour[len(colour)-2] <= '7') {
						// 30m > 37m
						fg = colour
						isReset = false
					} else if strings.Contains(colour, "[4") || strings.Contains(colour, "[10") && (colour[len(colour)-2] >= '0' && colour[len(colour)-2] <= '7') {
						// 40m > 47m
						bg = colour
						isReset = false
					} else if strings.Contains(colour, "[0m") {
						isReset = true
						fg, bg, colour = "", "", ""
					} else {
					}
				}
			} else {
				text += string(ch)
			}
		}
		if colour != "" || text != "" {
			if isReset {
				tokens = append(tokens, ANSILineToken{"\x1b[0m", "", text})
				isReset = false
			} else if colour != "\x1b[0m" && len(tokens) > 0 && tokens[len(tokens)-1].FG == "\x1b[0m" && tokens[len(tokens)-1].T == "" {
				// if the previous token was a reset, but didn't have any text, and this token sets a new colour,
				// then replace the previous reset token with the new token
				tokens[len(tokens)-1] = ANSILineToken{fg, bg, text}
			} else {
				// If we are setting a bg colour, but the last token didn't have one
				// then add a background clear to the previous bg.
				// This makes it less of a nightmare to flip horizontally if required.
				if bg != "" && len(tokens) > 0 && !strings.Contains(bg, "[49m") && tokens[len(tokens)-1].BG == "" {
					prevToken := tokens[len(tokens)-1]
					tokens[len(tokens)-1] = ANSILineToken{prevToken.FG, "\x1b[49m", prevToken.T}
				}
				tokens = append(tokens, ANSILineToken{fg, bg, text})
			}
		}
		lines = append(lines, tokens)
	}
	return lines
}

// BuildANSIString reconstructs an ANSI-formatted string from tokenized lines.
// It adds the specified padding (spaces) to the left of each line and ensures each line
// ends with a reset code.
func BuildANSIString(lines [][]ANSILineToken, padding int) string {
	var builder strings.Builder
	builder.Grow(500000) // Preallocate for large output

	paddingStr := strings.Repeat(" ", padding)
	for _, tokens := range lines {
		builder.WriteString(paddingStr) // add padding to the left
		for _, token := range tokens {
			builder.WriteString(token.FG)
			builder.WriteString(token.BG)
			builder.WriteString(token.T)
		}
		builder.WriteString("\x1b[0m\n")
	}
	return builder.String()
}

// FlipHorizontal horizontally flips tokenized ANSI lines while preserving formatting.
// It reverses the order of tokens on each line and the characters within each token's text.
// If mirrorMap is provided, it mirrors any characters found in the map as it reverses.
// All lines are padded on the left to maintain vertical alignment based on the widest line.
func FlipHorizontal(lines [][]ANSILineToken) [][]ANSILineToken {
	linesRev := make([][]ANSILineToken, len(lines))

	maxWidth := 0
	widths := make([]int, len(lines))
	for idx, l := range lines {
		lineWidth := 0
		for _, token := range l {
			lineWidth += UnicodeStringLength(token.T)
		}
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
		widths[idx] = lineWidth
	}

	for idx, tokens := range lines {
		revTokens := make([]ANSILineToken, 0)
		// ensure vertical alignment
		padding := maxWidth - widths[idx]
		
		// Reverse and mirror tokens
		for i := len(tokens) - 1; i >= 0; i-- {
			revTokens = append(revTokens, ANSILineToken{
				FG: tokens[i].FG,
				BG: tokens[i].BG,
				T:  MirrorHorizontally(tokens[i].T),
			})
		}
		
		// If padding is needed, prepend it to the first token if colors match, otherwise create new token
		if padding > 0 {
			paddingStr := strings.Repeat(" ", padding)
			if len(revTokens) > 0 && revTokens[0].FG == "" && revTokens[0].BG == "" {
				// First token has no colors, prepend padding to its text
				revTokens[0] = ANSILineToken{FG: "", BG: "", T: paddingStr + revTokens[0].T}
			} else {
				// First token has colors or no tokens exist, add padding as separate token
				revTokens = append([]ANSILineToken{{FG: "", BG: "", T: paddingStr}}, revTokens...)
			}
		}
		
		linesRev[idx] = revTokens
	}
	return linesRev
}

func FlipVertical(lines [][]ANSILineToken) [][]ANSILineToken {
	n := len(lines)
	flipped := make([][]ANSILineToken, n)

	for i, line := range lines {
		mirroredLine := make([]ANSILineToken, len(line))
		for j, tok := range line {
			mirroredLine[j] = ANSILineToken{
				FG: tok.FG, BG: tok.BG, T: MirrorVertically(tok.T),
			}
		}
		flipped[n-1-i] = mirroredLine
	}
	return flipped
}

func MirrorVertically(s string) string {
	runes := []rune(s)
	mirrored := make([]rune, len(runes))
	for i, r := range runes {
		mirrored[i] = getOrDefault(VerticalMirrorMap, r, r)
	}
	return string(mirrored)
}

// MirrorHorizontally reverses a string and mirrors any runes found in mirrorMap.
func MirrorHorizontally(s string) string {
	runes := []rune(s)
	mirrored := make([]rune, len(runes))
	for i, r := range runes {
		mirrored[len(runes)-1-i] = getOrDefault(HorizontalMirrorMap, r, r)
	}
	return string(mirrored)
}

// GetOrDefault returns the value for the key in the map,
// or the specified defaultValue if the key is not found.
func getOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

// ConvertAns converts legacy ANS format ANSI codes to modern UTF-8 format.
// It removes SAUCE metadata and adds reset codes before line endings for clean display.
// Lines are padded to the character width specified in SAUCE (or 80 by default).
// Long lines are wrapped at the character width boundary.
// The ANSI codes are passed through unchanged (CP437 decoding is done in main.go).
func ConvertAns(s string) string {
	// Default character width for ANSI art
	charWidth := 80

	// Check for SAUCE metadata and extract character width if present
	if idx := strings.Index(s, "SAUCE00"); idx != -1 {
		// SAUCE record is 128 bytes, character width is at offset 0x60 (96) from SAUCE00
		// The width is stored as a 2-byte little-endian integer
		if idx+96+2 <= len(s) {
			// Read the 2-byte width value (little-endian)
			widthBytes := []byte(s[idx+96 : idx+96+2])
			width := int(widthBytes[0]) | (int(widthBytes[1]) << 8)
			if width > 0 && width <= 1000 { // Sanity check
				charWidth = width
			}
		}
		// Remove SAUCE metadata
		s = s[:idx]
	}

	// Remove the EOF marker (SUB character, 0x1A) that precedes SAUCE
	s = strings.TrimRight(s, "\x1a")

	// Normalize line endings to LF only
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// Split into lines
	lines := strings.Split(s, "\n")
	var builder strings.Builder
	builder.Grow(len(s) + len(lines)*charWidth) // Pre-allocate

	defaultColors := "\x1b[37m\x1b[40m"

	// Process each line, wrapping if it exceeds charWidth
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Tokenize the line to separate ANSI codes from text
		tokens := TokeniseANSIString(line)
		if len(tokens) == 0 || len(tokens[0]) == 0 {
			continue
		}

		lineTokens := tokens[0] // TokeniseANSIString returns [][]ANSILineToken, we want the first line
		currentLineWidth := 0

		for tokenIdx, token := range lineTokens {
			// Start of line - add default colors
			if currentLineWidth == 0 {
				builder.WriteString(defaultColors)
			}

			// Write the ANSI codes (don't count toward width)
			// Skip reset codes and background resets at line start since we just set defaults
			if token.FG != "" && !(currentLineWidth == 0 && (token.FG == "\x1b[0m" || token.FG == "\x1b[49m")) {
				builder.WriteString(token.FG)
			}
			if token.BG != "" && !(currentLineWidth == 0 && (token.BG == "\x1b[0m" || token.BG == "\x1b[49m")) {
				builder.WriteString(token.BG)
			}

			// Process text character by character
			text := token.T
			for len(text) > 0 {
				// Get the first rune and its display width
				r, size := utf8.DecodeRuneInString(text)
				runeWidth := runewidth.RuneWidth(r)

				// Check if adding this character would exceed the line width
				if currentLineWidth+runeWidth > charWidth {
					// Pad remaining space on current line
					if currentLineWidth < charWidth {
						builder.WriteString(strings.Repeat(" ", charWidth-currentLineWidth))
					}
					// End current line
					builder.WriteString("\x1b[0m\n")
					// Start new line
					builder.WriteString(defaultColors)
					// Re-apply current colors for continuation
					if token.FG != "" && token.FG != "\x1b[0m" {
						builder.WriteString(token.FG)
					}
					if token.BG != "" {
						builder.WriteString(token.BG)
					}
					currentLineWidth = 0
				}

				// Add the character
				builder.WriteString(text[:size])
				currentLineWidth += runeWidth
				text = text[size:]

				// If we've reached the line width exactly, wrap to next line
				if currentLineWidth == charWidth && (len(text) > 0 || tokenIdx < len(lineTokens)-1) {
					builder.WriteString("\x1b[0m\n")
					builder.WriteString(defaultColors)
					// Re-apply current colors for continuation
					if token.FG != "" && token.FG != "\x1b[0m" {
						builder.WriteString(token.FG)
					}
					if token.BG != "" {
						builder.WriteString(token.BG)
					}
					currentLineWidth = 0
				}
			}
		}

		// Finish the current line
		if currentLineWidth > 0 {
			if currentLineWidth < charWidth {
				builder.WriteString(strings.Repeat(" ", charWidth-currentLineWidth))
			}
			builder.WriteString("\x1b[0m\n")
		}
	}

	return builder.String()
}
