package convert

import (
	"fmt"
	"os"
	"strings"

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

// SplitStringByWidth splits a string at a specific display width,
// returning (prefix, suffix) where prefix has the specified display width.
// This respects multibyte UTF-8 characters and double-width characters.
func SplitStringByWidth(s string, width int) (string, string) {
	if width <= 0 {
		return "", s
	}

	var currentWidth int
	var splitPos int

	for i, r := range s {
		charWidth := 1
		if r >= 128 {
			charWidth = runewidth.RuneWidth(r)
		}

		if currentWidth+charWidth > width {
			// Next character would exceed width, split here
			return s[:splitPos], s[splitPos:]
		}

		currentWidth += charWidth
		splitPos = i + len(string(r)) // Move splitPos to after this rune
	}
	return s, ""
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
				colour = string(ch)
			} else if isColour {
				// keep building the current ANSI escape code if \033 was found earlier
				// disable the isColour bool if the end of the ANSI escape code is found
				colour += string(ch)
				switch ch {
				case 'm':
					isColour = false
					// if there is text in the current token buffer, we need to finish it
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
						text = ""
					}
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
				case 'C':
					// Cursor forward - translate to spaces
					isColour = false
					numSpaces := 0
					fmt.Sscanf(colour, "\x1b[%dC", &numSpaces)
					text += strings.Repeat(" ", numSpaces)
				case 't':
					// this indicates a "true color" ANSI code!
					isColour = false
				default:
					// still in colour code
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

func AdjustANSILineWidths(lines [][]ANSILineToken, targetWidth int, targetLines int) ([][]ANSILineToken, error) {
	targetLinesKnown := targetLines != 0
	adjustedLines := make([][]ANSILineToken, 0)
	currLineN, currWidthN := 0, 0
	currTokenIdx, currTokenLineIdx := 0, 0

	splitTokenExists := false
	var splitToken ANSILineToken

	// Helper to check if we've consumed all input
	inputExhausted := func() bool {
		return !splitTokenExists && currTokenLineIdx >= len(lines)
	}

	for !inputExhausted() && (!targetLinesKnown || currLineN < targetLines) {
		// Ensure current line exists
		if currLineN >= len(adjustedLines) {
			adjustedLines = append(adjustedLines, make([]ANSILineToken, 0))
		}

		// Process tokens until we fill the current line
		for currWidthN < targetWidth {
			var currToken ANSILineToken
			var currTokenLen int

			// Get next token (either from split or from input)
			if splitTokenExists {
				currToken = splitToken
				splitTokenExists = false
				currTokenLen = UnicodeStringLength(splitToken.T)
			} else {
				// Check if we have more input
				if currTokenLineIdx >= len(lines) {
					return nil, fmt.Errorf("Not enough input to fill lines, current: %d, total: %d", currLineN, len(lines))
				}
				if currTokenIdx >= len(lines[currTokenLineIdx]) {
					currTokenLineIdx++
					currTokenIdx = 0
					if currTokenLineIdx >= len(lines) {
						break
					}
				}
				currToken = lines[currTokenLineIdx][currTokenIdx]
				currTokenLen = UnicodeStringLength(currToken.T)
				currTokenIdx++
			}

			// Check for too much input
			if targetLinesKnown && currLineN >= targetLines {
				return nil, fmt.Errorf("Too many characters for length %d and lines %d", targetWidth, targetLines)
			}

			// Calculate remaining space on current line
			remainingWidth := targetWidth - currWidthN

			// Can the current token fit in the current line?
			if currTokenLen <= remainingWidth {
				// Token fits completely
				adjustedLines[currLineN] = append(adjustedLines[currLineN], currToken)
				currWidthN += currTokenLen
			} else {
				// Token must be split - use width-aware splitting
				prefix, suffix := SplitStringByWidth(currToken.T, remainingWidth)

				// First part goes to current line
				adjustedLines[currLineN] = append(adjustedLines[currLineN], ANSILineToken{
					FG: currToken.FG, BG: currToken.BG, T: prefix,
				})
				// Remaining part will be processed in next line
				splitToken = ANSILineToken{
					FG: currToken.FG, BG: currToken.BG, T: suffix,
				}
				splitTokenExists = true
				currWidthN = targetWidth // Line is now full
			}

			// If line is full, move to next line
			if currWidthN >= targetWidth {
				break
			}
			if currTokenIdx >= len(lines[currTokenLineIdx]) && !splitTokenExists {
				break // Move to next line - will pad below if needed
			}
		}

		// Pad the line if it's not full (calculate actual width)
		actualWidth := 0
		for _, token := range adjustedLines[currLineN] {
			actualWidth += UnicodeStringLength(token.T)
		}
		if actualWidth < targetWidth {
			adjustedLines[currLineN] = append(adjustedLines[currLineN], ANSILineToken{
				FG: "\x1b[0m", BG: "\x1b[0m", T: strings.Repeat(" ", targetWidth-actualWidth),
			})
		}

		// Move to next line
		currLineN++
		currWidthN = 0
	}

	// Validate we got the expected number of lines if specified
	if targetLinesKnown && len(adjustedLines) != targetLines {
		return nil, fmt.Errorf("Expected %d lines but got %d", targetLines, len(adjustedLines))
	}

	return adjustedLines, nil
}

// ConvertAns converts legacy ANS format ANSI codes to modern UTF-8 format.
// It removes SAUCE metadata and adds reset codes before line endings for clean display.
// Lines are padded to the character width specified in SAUCE (or 80 by default).
// Long lines are wrapped at the character width boundary.
// The ANSI codes are passed through unchanged (CP437 decoding is done in main.go).
func ConvertAns(s string, info SAUCE) string {
	charWidth := 80 // Default character width for ANSI art
	fileLines := -1
	if info.TInfo1.Value > 0 {
		charWidth = int(info.TInfo1.Value) // Use character width from SAUCE if available
	}
	if info.TInfo2.Value > 0 {
		fileLines = int(info.TInfo2.Value) // Use number of lines from SAUCE if available
	}

	// Remove carriage returns (\r) from DOS line endings
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "")

	// Tokenise the input
	lines := TokeniseANSIString(s)

	// Adjust line widths (wrap/pad to match target width and lines)
	lines, err := AdjustANSILineWidths(lines, charWidth, fileLines)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adjusting line widths: %v\n", err)
		return ""
	}

	// Build the output
	var builder strings.Builder
	builder.Grow(len(s) + len(lines)*charWidth) // Pre-allocate

	for _, line := range lines {
		// Write the tokens for this line
		for i, token := range line {
			fg, bg := token.FG, token.BG

			// At line start, convert reset codes to explicit default colors
			if i == 0 {
				if fg == "\x1b[0m" {
					fg = "\x1b[37m" // White foreground
				}
				if bg == "\x1b[0m" || bg == "\x1b[49m" || bg == "" {
					bg = "\x1b[40m" // Black background
				}
			}

			// Don't write double reset at end of line (from padding)
			if i == len(line)-1 && fg == "\x1b[0m" && bg == "\x1b[0m" {
				// This is padding token, just write text
				builder.WriteString(token.T)
				continue
			}

			builder.WriteString(fg)
			builder.WriteString(bg)
			builder.WriteString(token.T)
		}

		// End each line with reset
		builder.WriteString("\x1b[0m\n")
	}

	return builder.String()
}
