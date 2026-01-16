package test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/convert"
	"github.com/tmck-code/go-ansi-convert/src/ansi-convert/parse"
)

// UnicodeChar is an exported struct for use in other packages
type UnicodeChar struct {
	Code int
	Char rune
}

var (
	successMark = "✓"
	failMark    = "❌"
)

// Checks if the test is running in Debug mode, i.e. has been run with the ENV var DEBUG=true.
// To do this, either first run `export DEBUG=true`, and then run the test command,
// or do it all at once with `DEBUG=true go test -v ./test“
func Debug() bool {
	debugValue := os.Getenv("DEBUG")
	return debugValue == "true" || debugValue == "1"
}

// Fails a test with a formatted message showing the expected vs. result. (These are both printed in %#v form)
func Fail(expected interface{}, result interface{}, t *testing.T) {
	t.Fatalf("\n\x1b[38;5;196m%s items don't match!\x1b[0m\n> expected:\t%#v\x1b[0m\n>   result:\t%#v\x1b[0m\n\n", failMark, expected, result)
}

// Takes in an expected & result object, of any type.
// Asserts that their Go syntax representations (%#v) are the same.
// Prints a message on success if the ENV var DEBUG is set to "true".
// Fails the test if this is not true.
// accepts an optional param "render" that acts like debug mode when set to true
func Assert(expected interface{}, result interface{}, t *testing.T) {
	expectedString, resultString := fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", result)
	if expectedString == resultString {
		if Debug() {
			t.Logf("\x1b[38;5;46m%s items match! expected/result:\x1b[0m\n\n%#v\x1b[0m\n\n", successMark, expected)
		}
		return
	}
	Fail(expectedString, resultString, t)
}

// Takes in an expected slice of objects and an 'item' object, of any type
// Asserts that the 'item' is contained within the slice.
// Prints a message on success if the ENV var DEBUG is set to "true".
// Fails the test if this is not true.
func AssertContains[T any](slice []T, item T, t *testing.T) {
	for _, el := range slice {
		if reflect.DeepEqual(el, item) {
			if Debug() {
				t.Logf("%s found expected item!\n>  item:\t%v\n> slice:\t%v\n", successMark, item, slice)
			}
			return
		}
	}
	Fail(slice, item, t)
}

// Flattens a given json string, removing all tabs, spaces and newlines
func FlattenJSON(json string) string {
	json = strings.Replace(json, "\n", "", -1)
	json = strings.Replace(json, "\t", "", -1)
	json = strings.Replace(json, " ", "", -1)
	return json
}

// AddBorder adds a box border around a multi-line string, handling ANSI escape codes.
func AddBorder(s string, pad bool) string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	maxLen := parse.LongestUnicodeLineLength(lines)
	result := make([]string, len(lines)+2)

	for i, line := range lines {
		if pad {
			visualLen := parse.UnicodeStringLength(line)
			padding := strings.Repeat(" ", maxLen-visualLen)
			result[i+1] = " " + line + padding + "┊"
		} else {
			result[i+1] = " " + line + "┊"
		}
	}
	return strings.Join(result, "\n") + "\n"
}

// TestTitleInput returns the formatted title string for "input" test sections.
func TestTitleInput() string {
	return "\x1b[44;30;1;3m ▶ input \x1b[0m\x1b[34m🭍🭑🬽\x1b[0m"
}

// TestTitleExpected returns the formatted title string for "expected" test sections.
func TestTitleExpected() string {
	return "\x1b[42;30;1;3m ✓ expected \x1b[0m\x1b[32m🭍🭑🬽\x1b[0m"
}

// TestTitleResult returns the formatted title string for "result" test sections.
func TestTitleResult() string {
	return "\x1b[43;30;1;3m ✭ result \x1b[0m\x1b[33m🭍🭑🬽\x1b[0m"
}

// PrintSimpleTestResults prints formatted test results for simple tests (with quoted output).
func PrintSimpleTestResults(input string, expected string, result string, t *testing.T) {
	t.Logf("%s\n%v\x1b[0m", TestTitleInput(), AddBorder(input, false))
	t.Logf("%s\n%v\x1b[0m", TestTitleExpected(), AddBorder(expected, false))
	t.Logf("%s\n%v\x1b[0m", TestTitleResult(), AddBorder(result, false))
}

// PrintANSITestResults prints formatted test results for ANSI tokenization and reversal tests.
func PrintANSITestResults(input string, expected, result [][]convert.ANSILineToken, t *testing.T, noPrintStrings ...bool) {
	if len(noPrintStrings) == 0 || !noPrintStrings[0] {
		t.Logf("%s\n%s\x1b[0m", TestTitleInput(), AddBorder(input, false))
		t.Logf("%s\n%s\x1b[0m", TestTitleExpected(), AddBorder(convert.BuildANSIString(expected, 0), false))
		t.Logf("%s\n%s\x1b[0m\n", TestTitleResult(), AddBorder(convert.BuildANSIString(result, 0), false))
	}

	if !reflect.DeepEqual(expected, result) {
		// Print detailed diff between expected and result
		PrintTokensForCopy(expected, result, t)
	}

	if Debug() {
		for i, line := range expected {
			t.Logf("%s %+v\x1b[0m\n", TestTitleExpected(), line)
			t.Logf("%s %+v\x1b[0m\n", TestTitleResult(), result[i])

			eb, err := json.MarshalIndent(line, "", "  ")
			if err != nil {
				t.Logf("error: %v", err)
			}
			t.Logf("%s %+v\x1b[0m\n", TestTitleExpected(), string(eb))
			rb, err := json.MarshalIndent(result[i], "", "  ")
			if err != nil {
				t.Logf("error: %v", err)
			}
			t.Logf("%s %+v\x1b[0m\n", TestTitleResult(), string(rb))
			for j, token := range result[i] {
				Assert(line[j], token, t)
				if (j + 1) < len(line) {
					break
				}
			}
			Assert(line, result[i], t)
		}
	}
}

func PrintSAUCETestResults(input string, expected, result *convert.SAUCE, t *testing.T) {
	if Debug() {
		eb, err := json.MarshalIndent(expected, "", "  ")
		if err != nil {
			t.Logf("error: %v", err)
		}
		t.Logf("%s\n%+v\x1b[0m\n", TestTitleExpected(), string(eb))
		rb, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			t.Logf("error: %v", err)
		}
		t.Logf("%s\n%+v\x1b[0m\n", TestTitleResult(), string(rb))
		Assert(expected, result, t)
	}
}

// FormatTokensForTest formats [][]ANSILineToken in a copy-paste ready format for test code.
// This makes it easy to copy the output and paste it directly into test expected values.
func FormatTokensForTest(tokens [][]convert.ANSILineToken) string {
	var sb strings.Builder
	sb.WriteString("[][]convert.ANSILineToken{\n")

	for _, line := range tokens {
		sb.WriteString("\t{\n")
		for _, token := range line {
			sb.WriteString(fmt.Sprintf("\t\t{FG: %q, BG: %q, Control: %q, T: %q},\n",
				token.FG, token.BG, token.Control, token.T))
		}
		sb.WriteString("\t},\n")
	}

	sb.WriteString("}")
	return sb.String()
}

// PrintTokensForCopy prints a git-style diff between expected and result tokens.
// Use this in test failures to understand what differs and generate the correct expected output.
func PrintTokensForCopy(expected, result [][]convert.ANSILineToken, t *testing.T) {
	t.Logf("\n\x1b[1;37mdiff expected/result\x1b[0m")
	t.Logf("\x1b[1;31m--- expected\x1b[0m")
	t.Logf("\x1b[1;32m+++ result\x1b[0m")

	maxLines := len(expected)
	if len(result) > maxLines {
		maxLines = len(result)
	}

	hasDifferences := false

	for i := 0; i < maxLines; i++ {
		if i >= len(expected) {
			// Extra line in result
			hasDifferences = true
			t.Logf("\x1b[36m@@ Line %d @@\x1b[0m", i)
			for _, token := range result[i] {
				t.Logf("\x1b[32m+\t\t{FG: %q, BG: %q, Control: %q, T: %q},\x1b[0m",
					token.FG, token.BG, token.Control, token.T)
			}
			continue
		}
		if i >= len(result) {
			// Missing line in result
			hasDifferences = true
			t.Logf("\x1b[36m@@ Line %d @@\x1b[0m", i)
			for _, token := range expected[i] {
				t.Logf("\x1b[31m-\t\t{FG: %q, BG: %q, Control: %q, T: %q},\x1b[0m",
					token.FG, token.BG, token.Control, token.T)
			}
			continue
		}

		expectedLine := expected[i]
		resultLine := result[i]

		if reflect.DeepEqual(expectedLine, resultLine) {
			// Lines match - show in debug mode only
			if Debug() {
				t.Logf("\x1b[90m Line %d (no change)\x1b[0m", i)
			}
			continue
		}

		// Lines differ - show as git-style diff with context
		hasDifferences = true
		t.Logf("\x1b[36m@@ Line %d @@\x1b[0m", i)

		// Compare tokens and show context + differences
		maxTokens := len(expectedLine)
		if len(resultLine) > maxTokens {
			maxTokens = len(resultLine)
		}

		for j := 0; j < maxTokens; j++ {
			if j >= len(expectedLine) {
				// Extra token in result
				t.Logf("\x1b[32m+\t\t{FG: %q, BG: %q, Control: %q, T: %q},\x1b[0m",
					resultLine[j].FG, resultLine[j].BG, resultLine[j].Control, resultLine[j].T)
			} else if j >= len(resultLine) {
				// Missing token in result
				t.Logf("\x1b[31m-\t\t{FG: %q, BG: %q, Control: %q, T: %q},\x1b[0m",
					expectedLine[j].FG, expectedLine[j].BG, expectedLine[j].Control, expectedLine[j].T)
			} else if !reflect.DeepEqual(expectedLine[j], resultLine[j]) {
				// Token differs - show both
				t.Logf("\x1b[31m-\t\t{FG: %q, BG: %q, Control: %q, T: %q},\x1b[0m",
					expectedLine[j].FG, expectedLine[j].BG, expectedLine[j].Control, expectedLine[j].T)
				t.Logf("\x1b[32m+\t\t{FG: %q, BG: %q, Control: %q, T: %q},\x1b[0m",
					resultLine[j].FG, resultLine[j].BG, resultLine[j].Control, resultLine[j].T)
			} else {
				// Token matches - show as context
				t.Logf("\t\t{FG: %q, BG: %q, Control: %q, T: %q},",
					expectedLine[j].FG, expectedLine[j].BG, expectedLine[j].Control, expectedLine[j].T)
			}
		}
	}

	if hasDifferences {
		t.Logf("\n\x1b[1;36m━━━ Copy-paste ready Result tokens ━━━\x1b[0m\n%s\n", FormatTokensForTest(result))
	}
}
