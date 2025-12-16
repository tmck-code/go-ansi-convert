package ansi_flip

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// Returns the length of a string, taking into account Unicode characters and ANSI escape codes.
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

func SanitiseUnicodeString(s string) string {
	if s == "" {
		return s
	}

	tokenizedLines := TokeniseANSIString(s)
	if len(tokenizedLines) == 0 {
		return s
	}

	result := make([]string, 0, len(tokenizedLines))

	for _, tokens := range tokenizedLines {
		lineStr := ""
		for _, token := range tokens {
			lineStr += token.FG + token.BG + token.T
		}
		// Check if the last token has a reset
		hasReset := false
		if len(tokens) > 0 {
			lastToken := tokens[len(tokens)-1]
			hasReset = lastToken.FG == "\x1b[0m"
		}
		// Ensure line ends with reset if not already present
		if !hasReset {
			lineStr += "\x1b[0m"
		}
		result = append(result, lineStr)
	}
	return strings.Join(result, "\n")
}

func ReverseUnicodeString(s string) string {
	runes := []rune(s)
	reversed := make([]rune, len(runes))

	for i, r := range runes {
		reversed[len(runes)-1-i] = r
	}
	return string(reversed)
}

type ANSILineToken struct {
	FG string
	BG string
	T  string
	// Reset    bool
}

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
					if strings.Contains(colour, "[38") || strings.Contains(colour, "[39") {
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
				// if we are setting a bg colour, but the last token didn't have one
				// then add a background clear to the previous bg
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

func BuildANSIString(lines [][]ANSILineToken, padding int) string {
	s := ""

	for _, tokens := range lines {
		s += strings.Repeat(" ", padding) // add padding to the left
		for _, token := range tokens {
			s += token.FG + token.BG + token.T
		}
		s += "\x1b[0m\n"
	}
	return s
}

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
				T:  ReverseUnicodeString(tokens[i].T),
			})
		}
		linesRev[idx] = revTokens
	}
	return linesRev
}
