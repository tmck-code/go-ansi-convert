package convert

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

var (
	HorizontalMirrorMap = map[rune]rune{
		'<': '>', '>': '<',
		'(': ')', ')': '(',
		'[': ']', ']': '[',
		'{': '}', '}': '{',
		'/': '\\', '\\': '/',
		'b': 'd', 'd': 'b',
		'p': 'q', 'q': 'p',
		'B': 'ᗺ', 'ᗺ': 'B',
		'C': 'Ɔ', 'Ɔ': 'C',
		'D': 'ᗡ', 'ᗡ': 'D',
		'E': 'Ǝ', 'Ǝ': 'E',
		'F': 'ꟻ', 'ꟻ': 'F',
		'G': 'ວ', 'ວ': 'G',
		'J': 'ᒐ', 'ᒐ': 'J',
		'K': 'ꓘ', 'ꓘ': 'K',
		'L': '⅃', '⅃': 'L',
		'N': 'И', 'И': 'N',
		'O': 'O',
		'P': 'ᑫ', 'ᑫ': 'P',
		'Q': 'Ϙ', 'Ϙ': 'Q',
		'R': 'Я', 'Я': 'R',
		'S': 'Ƨ', 'Ƨ': 'S',
		'a': 'ɒ', 'ɒ': 'a',
		'c': 'ɔ', 'ɔ': 'c',
		'e': 'ɘ', 'ɘ': 'e',
		'f': 'ᆿ', 'ᆿ': 'f',
		'g': 'ϱ', 'ϱ': 'g',
		'h': '⑁', '⑁': 'h',
		'j': 'ᒑ', 'ᒑ': 'j',
		'k': 'ʞ', 'ʞ': 'k',
		'r': 'ɿ', 'ɿ': 'r',
		's': 'ƨ', 'ƨ': 's',
		't': 'ɟ', 'ɟ': 't',
		'y': 'γ', 'γ': 'y',
		'┌': '┐', '┐': '┌',
		'┍': '┑', '┑': '┍',
		'┎': '┒', '┒': '┎',
		'┏': '┓', '┓': '┏',
		'└': '┘', '┘': '└',
		'┕': '┙', '┙': '┕',
		'┖': '┚', '┚': '┖',
		'┗': '┛', '┛': '┗',
		'├': '┤', '┤': '├',
		'┝': '┥', '┥': '┝',
		'┞': '┦', '┦': '┞',
		'┟': '┧', '┧': '┟',
		'┠': '┨', '┨': '┠',
		'┡': '┩', '┩': '┡',
		'┢': '┪', '┪': '┢',
		'┣': '┫', '┫': '┣',
		'╒': '╕', '╕': '╒',
		'╓': '╖', '╖': '╓',
		'╔': '╗', '╗': '╔',
		'╘': '╛', '╛': '╘',
		'╙': '╜', '╜': '╙',
		'╚': '╝', '╝': '╚',
		'╞': '╡', '╡': '╞',
		'╟': '╢', '╢': '╟',
		'╠': '╣', '╣': '╠',
		'╭': '╮', '╮': '╭',
		'╰': '╯', '╯': '╰',
		'╴': '╶', '╶': '╴',
		'╸': '╺', '╺': '╸',
		'╼': '╾', '╾': '╼',
		'▖': '▗', '▗': '▖',
		'▘': '▝', '▝': '▘',
		'▌': '▐', '▐': '▌',
		'▙': '▜', '▜': '▙',
		'▚': '▞', '▞': '▚',
	}

	VerticalMirrorMap = map[rune]rune{
		'!': '¡', '¡': '!',
		'A': 'ꓯ', 'ꓯ': 'A',
		'M': 'ꟽ', 'ꟽ': 'M',
		'Q': 'Ό', 'Ό': 'Q',
		'T': 'ꓕ', 'ꓕ': 'T',
		'Y': '⅄', '⅄': 'Y',
		'U': 'ꓵ', 'ꓵ': 'U',
		'V': 'ꓥ', 'ꓥ': 'V',
		'W': 'M', '^': 'v',
		'v': '^', 'w': 'm',
		'm': 'w', 'u': 'n',
		'n': 'u', 'i': 'ᴉ',
		'ᴉ': 'i', 'ɯ': 'm',
		'ʌ': 'v', 'ʍ': 'w',
		'y': 'ʎ', 'ʎ': 'y',
		'┌': '└', '└': '┌',
		'┍': '┕', '┕': '┍',
		'┎': '┖', '┖': '┎',
		'┏': '┗', '┗': '┏',
		'┐': '┘', '┘': '┐',
		'┑': '┙', '┙': '┑',
		'┒': '┚', '┚': '┒',
		'┓': '┛', '┛': '┓',
		'┬': '┴', '┴': '┬',
		'┭': '┵', '┵': '┭',
		'┮': '┶', '┶': '┮',
		'┯': '┷', '┷': '┯',
		'┰': '┸', '┸': '┰',
		'┱': '┹', '┹': '┱',
		'┲': '┺', '┺': '┲',
		'┳': '┻', '┻': '┳',
		'╒': '╘', '╘': '╒',
		'╓': '╙', '╙': '╓',
		'╔': '╚', '╚': '╔',
		'╕': '╛', '╛': '╕',
		'╖': '╜', '╜': '╖',
		'╗': '╝', '╝': '╗',
		'╤': '╧', '╧': '╤',
		'╥': '╨', '╨': '╥',
		'╦': '╩', '╩': '╦',
		'╭': '╰', '╰': '╭',
		'╮': '╯', '╯': '╮',
		'╵': '╷', '╷': '╵',
		'╹': '╻', '╻': '╹',
		'╽': '╿', '╿': '╽',
		'▀': '▄', '▄': '▀',
		'▔': '▁', '▁': '▔',
		'▖': '▘', '▘': '▖',
		'▗': '▝', '▝': '▗',
		'▙': '▛', '▛': '▙',
		'▚': '▞', '▞': '▚',
		'H': 'H', 'I': 'I',
		'O': 'O', 'X': 'X',
		'l': 'l', 'x': 'x',
		'├': '├', '┝': '┝',
		'┠': '┠', '┣': '┣',
		'┤': '┤', '┥': '┥',
		'┨': '┨', '┫': '┫',
		'╞': '╞', '╟': '╟',
		'╠': '╠', '╡': '╡',
		'╢': '╢', '╣': '╣',
		'╴': '╴', '╶': '╶',
		'╸': '╸', '╺': '╺',
		'╼': '╼', '╾': '╾',
		'─': '─', '━': '━',
		'│': '│', '┃': '┃',
		'═': '═', '║': '║',
		'┼': '┼', '┽': '┽',
		'┾': '┾', '┿': '┿',
		'╂': '╂', '╉': '╉',
		'╊': '╊', '╋': '╋',
		'╪': '╪', '╫': '╫',
		'╬': '╬', '╳': '╳',
		'▌': '▌', '▍': '▍',
		'▎': '▎', '▏': '▏',
		'▐': '▐', '▕': '▕',
		'█': '█', '▬': '▬',
		'░': '░', '▒': '▒',
		'▓': '▓',
	}
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

// LongestUnicodeLineLength finds the maximum display length among all lines.
// It uses UnicodeStringLength to properly measure each line's width.
func LongestUnicodeLineLength(lines []string) int {
	maxLen := 0
	for _, line := range lines {
		length := UnicodeStringLength(line)
		if length > maxLen {
			maxLen = length
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

	result := make([]string, len(tokenizedLines))

	// Calculate max length if justification is needed
	maxLen := 0
	if justifyLines {
		maxLen = LongestUnicodeLineLength(strings.Split(s, "\n"))
	}

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
		result[i] = lineBuilder.String()
	}
	return strings.Join(result, "\n")
}

// ReverseUnicodeString reverses the order of characters in a string.
// It works on a rune-by-rune basis to properly handle multi-byte Unicode characters.
func ReverseUnicodeString(s string) string {
	runes := []rune(s)
	reversed := make([]rune, len(runes))

	for i, r := range runes {
		reversed[len(runes)-1-i] = getOrDefault(HorizontalMirrorMap, r, r)
	}
	return string(reversed)
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

	for _, line := range strings.Split(msg, "\n") {
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
					} else if strings.Contains(colour, "[38") || strings.Contains(colour, "[39") {
						fg = colour
						isReset = false
					} else if strings.Contains(colour, "[48") || strings.Contains(colour, "[49") {
						bg = colour
						isReset = false
					} else if strings.Contains(colour, "[0m") {
						isReset = true
						fg = ""
						bg = ""
						colour = ""
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

// ReverseANSIString horizontally flips tokenized ANSI lines while preserving formatting.
// It reverses the order of tokens on each line and the characters within each token's text.
// If mirrorMap is provided, it mirrors any characters found in the map as it reverses.
// All lines are padded on the left to maintain vertical alignment based on the widest line.
func ReverseANSIString(lines [][]ANSILineToken) [][]ANSILineToken {
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
		revTokens := make([]ANSILineToken, 1)
		// ensure vertical alignment
		revTokens[0] = ANSILineToken{FG: "", BG: "", T: strings.Repeat(" ", maxWidth-widths[idx])}
		for i := len(tokens) - 1; i >= 0; i-- {
			revTokens = append(revTokens, ANSILineToken{
				FG: tokens[i].FG,
				BG: tokens[i].BG,
				T:  reverseAndMirrorUnicodeString(tokens[i].T),
			})
		}
		linesRev[idx] = revTokens
	}
	return linesRev
}

func FlipVertical(msg string) [][]ANSILineToken {
	tokenized := TokeniseANSIString(msg)
	n := len(tokenized)
	flipped := make([][]ANSILineToken, n)
	for i, line := range tokenized {
		mirroredLine := make([]ANSILineToken, len(line))
		for j, tok := range line {
			// Mirror each rune in the token's text
			runes := []rune(tok.T)
			for k, r := range runes {
				runes[k] = getOrDefault(VerticalMirrorMap, r, r)
			}
			tok.T = string(runes)
			mirroredLine[j] = tok
		}
		flipped[n-1-i] = mirroredLine
	}
	return flipped
}

// reverseAndMirrorUnicodeString reverses a string and mirrors any runes found in mirrorMap.
func reverseAndMirrorUnicodeString(s string) string {
	runes := []rune(s)
	reversed := make([]rune, len(runes))
	for i, r := range runes {
		reversed[len(runes)-1-i] = getOrDefault(HorizontalMirrorMap, r, r)
	}
	return string(reversed)
}

// GetOrDefault returns the value for the key in the map,
// or the specified defaultValue if the key is not found.
func getOrDefault[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultValue
}

