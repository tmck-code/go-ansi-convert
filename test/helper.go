package test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
)

var (
	successMark = "✓"
	failMark    = "❌"
)

// Checks if the test is running in Debug mode, i.e. has been run with the ENV var DEBUG=true.
// To do this, either first run `export DEBUG=true`, and then run the test command,
// or do it all at once with `DEBUG=true go test -v ./test“
func Debug() bool {
	return os.Getenv("DEBUG") == "true"
}

// Fails a test with a formatted message showing the expected vs. result. (These are both printed in %#v form)
func Fail(expected interface{}, result interface{}, test *testing.T) {
	test.Fatalf("\n\x1b[38;5;196m%s items don't match!\x1b[0m\n> expected:\t%#v\x1b[0m\n>   result:\t%#v\x1b[0m\n\n", failMark, expected, result)
}

// Takes in an expected & result object, of any type.
// Asserts that their Go syntax representations (%#v) are the same.
// Prints a message on success if the ENV var DEBUG is set to "true".
// Fails the test if this is not true.
func Assert(expected interface{}, result interface{}, test *testing.T) {
	expectedString, resultString := fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", result)
	if expectedString == resultString {
		if Debug() {
			fmt.Printf("\x1b[38;5;46m%s items match!\x1b[0m\n> expected:\t%#v\x1b[0m\n>   result:\t%#v\x1b[0m\n\n", successMark, expected, result)
		}
		return
	}
	Fail(expectedString, resultString, test)
}

// Takes in an expected slice of objects and an 'item' object, of any type
// Asserts that the 'item' is contained within the slice.
// Prints a message on success if the ENV var DEBUG is set to "true".
// Fails the test if this is not true.
func AssertContains[T any](slice []T, item T, test *testing.T) {
	for _, el := range slice {
		if reflect.DeepEqual(el, item) {
			if Debug() {
				fmt.Printf("%s found expected item!\n>  item:\t%v\n> slice:\t%v\n", successMark, item, slice)
			}
			return
		}
	}
	Fail(slice, item, test)
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

	maxLen := convert.LongestUnicodeLineLength(lines)

	result := make([]string, len(lines)+2)
	result[0] = "╭" + strings.Repeat("─", maxLen) + "╮"

	for i, line := range lines {
		if pad {
			visualLen := convert.UnicodeStringLength(line)
			padding := strings.Repeat(" ", maxLen-visualLen)
			result[i+1] = "│" + line + padding + "│"
		} else {
			result[i+1] = "│" + line + "│"
		}
	}
	result[len(result)-1] = "╰" + strings.Repeat("─", maxLen) + "╯"
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
func PrintSimpleTestResults(input string, expected string, result string) {
	fmt.Printf("%s\n%v\x1b[0m", TestTitleInput(), AddBorder(input, false))
	fmt.Printf("%s\n%v\x1b[0m", TestTitleExpected(), AddBorder(expected, false))
	fmt.Printf("%s\n%v\x1b[0m\n", TestTitleResult(), AddBorder(result, false))
}

// PrintANSITestResults prints formatted test results for ANSI tokenization and reversal tests.
func PrintANSITestResults(input string, expected, result [][]convert.ANSILineToken, test *testing.T) {
	fmt.Printf("%s\n%s\x1b[0m", TestTitleInput(), AddBorder(input, false))
	fmt.Printf("%s\n%s\x1b[0m", TestTitleExpected(), AddBorder(convert.BuildANSIString(expected, 0), false))
	fmt.Printf("%s\n%s\x1b[0m\n", TestTitleResult(), AddBorder(convert.BuildANSIString(result, 0), false))

	if Debug() {
		for i, line := range expected {
			fmt.Printf("%s %+v\x1b[0m\n", TestTitleExpected(), line)
			fmt.Printf("%s %+v\x1b[0m\n", TestTitleResult(), result[i])

			eb, err := json.MarshalIndent(line, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Printf("%s %+v\x1b[0m\n", TestTitleExpected(), string(eb))
			rb, err := json.MarshalIndent(result[i], "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Printf("%s %+v\x1b[0m\n", TestTitleResult(), string(rb))
			for j, token := range result[i] {
				Assert(line[j], token, test)
				if (j + 1) < len(line) {
					break
				}
			}
			Assert(line, result[i], test)
		}
	}
}
