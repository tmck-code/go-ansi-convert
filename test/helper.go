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

type unicodeChar struct {
	code int
	char rune
}

var (
	successMark = "✓"
	failMark    = "❌"
	// from https://en.wikipedia.org/wiki/Box-drawing_characters#Symbols_for_Legacy_Computing
	// unicode chars from U+1FB00 to U+1FBFF
	CompleteSet = []unicodeChar{
		{code: 0x1fb00, char: '🬀'}, {code: 0x1fb01, char: '🬁'}, {code: 0x1fb02, char: '🬂'}, {code: 0x1fb03, char: '🬃'},
		{code: 0x1fb04, char: '🬄'}, {code: 0x1fb05, char: '🬅'}, {code: 0x1fb06, char: '🬆'}, {code: 0x1fb07, char: '🬇'},
		{code: 0x1fb08, char: '🬈'}, {code: 0x1fb09, char: '🬉'}, {code: 0x1fb0a, char: '🬊'}, {code: 0x1fb0b, char: '🬋'},
		{code: 0x1fb0c, char: '🬌'}, {code: 0x1fb0d, char: '🬍'}, {code: 0x1fb0e, char: '🬎'}, {code: 0x1fb0f, char: '🬏'},
		{code: 0x1fb10, char: '🬐'}, {code: 0x1fb11, char: '🬑'}, {code: 0x1fb12, char: '🬒'}, {code: 0x1fb13, char: '🬓'},
		{code: 0x1fb14, char: '🬔'}, {code: 0x1fb15, char: '🬕'}, {code: 0x1fb16, char: '🬖'}, {code: 0x1fb17, char: '🬗'},
		{code: 0x1fb18, char: '🬘'}, {code: 0x1fb19, char: '🬙'}, {code: 0x1fb1a, char: '🬚'}, {code: 0x1fb1b, char: '🬛'},
		{code: 0x1fb1c, char: '🬜'}, {code: 0x1fb1d, char: '🬝'}, {code: 0x1fb1e, char: '🬞'}, {code: 0x1fb1f, char: '🬟'},
		{code: 0x1fb20, char: '🬠'}, {code: 0x1fb21, char: '🬡'}, {code: 0x1fb22, char: '🬢'}, {code: 0x1fb23, char: '🬣'},
		{code: 0x1fb24, char: '🬤'}, {code: 0x1fb25, char: '🬥'}, {code: 0x1fb26, char: '🬦'}, {code: 0x1fb27, char: '🬧'},
		{code: 0x1fb28, char: '🬨'}, {code: 0x1fb29, char: '🬩'}, {code: 0x1fb2a, char: '🬪'}, {code: 0x1fb2b, char: '🬫'},
		{code: 0x1fb2c, char: '🬬'}, {code: 0x1fb2d, char: '🬭'}, {code: 0x1fb2e, char: '🬮'}, {code: 0x1fb2f, char: '🬯'},
		{code: 0x1fb30, char: '🬰'}, {code: 0x1fb31, char: '🬱'}, {code: 0x1fb32, char: '🬲'}, {code: 0x1fb33, char: '🬳'},
		{code: 0x1fb34, char: '🬴'}, {code: 0x1fb35, char: '🬵'}, {code: 0x1fb36, char: '🬶'}, {code: 0x1fb37, char: '🬷'},
		{code: 0x1fb38, char: '🬸'}, {code: 0x1fb39, char: '🬹'}, {code: 0x1fb3a, char: '🬺'}, {code: 0x1fb3b, char: '🬻'},
		{code: 0x1fb3c, char: '🬼'}, {code: 0x1fb3d, char: '🬽'}, {code: 0x1fb3e, char: '🬾'}, {code: 0x1fb3f, char: '🬿'},
		{code: 0x1fb40, char: '🭀'}, {code: 0x1fb41, char: '🭁'}, {code: 0x1fb42, char: '🭂'}, {code: 0x1fb43, char: '🭃'},
		{code: 0x1fb44, char: '🭄'}, {code: 0x1fb45, char: '🭅'}, {code: 0x1fb46, char: '🭆'}, {code: 0x1fb47, char: '🭇'},
		{code: 0x1fb48, char: '🭈'}, {code: 0x1fb49, char: '🭉'}, {code: 0x1fb4a, char: '🭊'}, {code: 0x1fb4b, char: '🭋'},
		{code: 0x1fb4c, char: '🭌'}, {code: 0x1fb4d, char: '🭍'}, {code: 0x1fb4e, char: '🭎'}, {code: 0x1fb4f, char: '🭏'},
		{code: 0x1fb50, char: '🭐'}, {code: 0x1fb51, char: '🭑'}, {code: 0x1fb52, char: '🭒'}, {code: 0x1fb53, char: '🭓'},
		{code: 0x1fb54, char: '🭔'}, {code: 0x1fb55, char: '🭕'}, {code: 0x1fb56, char: '🭖'}, {code: 0x1fb57, char: '🭗'},
		{code: 0x1fb58, char: '🭘'}, {code: 0x1fb59, char: '🭙'}, {code: 0x1fb5a, char: '🭚'}, {code: 0x1fb5b, char: '🭛'},
		{code: 0x1fb5c, char: '🭜'}, {code: 0x1fb5d, char: '🭝'}, {code: 0x1fb5e, char: '🭞'}, {code: 0x1fb5f, char: '🭟'},
		{code: 0x1fb60, char: '🭠'}, {code: 0x1fb61, char: '🭡'}, {code: 0x1fb62, char: '🭢'}, {code: 0x1fb63, char: '🭣'},
		{code: 0x1fb64, char: '🭤'}, {code: 0x1fb65, char: '🭥'}, {code: 0x1fb66, char: '🭦'}, {code: 0x1fb67, char: '🭧'},
		{code: 0x1fb68, char: '🭨'}, {code: 0x1fb69, char: '🭩'}, {code: 0x1fb6a, char: '🭪'}, {code: 0x1fb6b, char: '🭫'},
		{code: 0x1fb6c, char: '🭬'}, {code: 0x1fb6d, char: '🭭'}, {code: 0x1fb6e, char: '🭮'}, {code: 0x1fb6f, char: '🭯'},
		{code: 0x1fb70, char: '🭰'}, {code: 0x1fb71, char: '🭱'}, {code: 0x1fb72, char: '🭲'}, {code: 0x1fb73, char: '🭳'},
		{code: 0x1fb74, char: '🭴'}, {code: 0x1fb75, char: '🭵'}, {code: 0x1fb76, char: '🭶'}, {code: 0x1fb77, char: '🭷'},
		{code: 0x1fb78, char: '🭸'}, {code: 0x1fb79, char: '🭹'}, {code: 0x1fb7a, char: '🭺'}, {code: 0x1fb7b, char: '🭻'},
		{code: 0x1fb7c, char: '🭼'}, {code: 0x1fb7d, char: '🭽'}, {code: 0x1fb7e, char: '🭾'}, {code: 0x1fb7f, char: '🭿'},
		{code: 0x1fb80, char: '🮀'}, {code: 0x1fb81, char: '🮁'}, {code: 0x1fb82, char: '🮂'}, {code: 0x1fb83, char: '🮃'},
		{code: 0x1fb84, char: '🮄'}, {code: 0x1fb85, char: '🮅'}, {code: 0x1fb86, char: '🮆'}, {code: 0x1fb87, char: '🮇'},
		{code: 0x1fb88, char: '🮈'}, {code: 0x1fb89, char: '🮉'}, {code: 0x1fb8a, char: '🮊'}, {code: 0x1fb8b, char: '🮋'},
		{code: 0x1fb8c, char: '🮌'}, {code: 0x1fb8d, char: '🮍'}, {code: 0x1fb8e, char: '🮎'}, {code: 0x1fb8f, char: '🮏'},
		{code: 0x1fb90, char: '🮐'}, {code: 0x1fb91, char: '🮑'}, {code: 0x1fb92, char: '🮒'}, {code: 0x1fb93, char: '🮓'},
		{code: 0x1fb94, char: '🮔'}, {code: 0x1fb95, char: '🮕'}, {code: 0x1fb96, char: '🮖'}, {code: 0x1fb97, char: '🮗'},
		{code: 0x1fb98, char: '🮘'}, {code: 0x1fb99, char: '🮙'}, {code: 0x1fb9a, char: '🮚'}, {code: 0x1fb9b, char: '🮛'},
		{code: 0x1fb9c, char: '🮜'}, {code: 0x1fb9d, char: '🮝'}, {code: 0x1fb9e, char: '🮞'}, {code: 0x1fb9f, char: '🮟'},
		{code: 0x1fba0, char: '🮠'}, {code: 0x1fba1, char: '🮡'}, {code: 0x1fba2, char: '🮢'}, {code: 0x1fba3, char: '🮣'},
		{code: 0x1fba4, char: '🮤'}, {code: 0x1fba5, char: '🮥'}, {code: 0x1fba6, char: '🮦'}, {code: 0x1fba7, char: '🮧'},
		{code: 0x1fba8, char: '🮨'}, {code: 0x1fba9, char: '🮩'}, {code: 0x1fbaa, char: '🮪'}, {code: 0x1fbab, char: '🮫'},
		{code: 0x1fbac, char: '🮬'}, {code: 0x1fbad, char: '🮭'}, {code: 0x1fbae, char: '🮮'}, {code: 0x1fbaf, char: '🮯'},
		{code: 0x1fbb0, char: '🮰'}, {code: 0x1fbb1, char: '🮱'}, {code: 0x1fbb2, char: '🮲'}, {code: 0x1fbb3, char: '🮳'},
		{code: 0x1fbb4, char: '🮴'}, {code: 0x1fbb5, char: '🮵'}, {code: 0x1fbb6, char: '🮶'}, {code: 0x1fbb7, char: '🮷'},
		{code: 0x1fbb8, char: '🮸'}, {code: 0x1fbb9, char: '🮹'}, {code: 0x1fbba, char: '🮺'}, {code: 0x1fbbb, char: '🮻'},
		{code: 0x1fbbc, char: '🮼'}, {code: 0x1fbbd, char: '🮽'}, {code: 0x1fbbe, char: '🮾'}, {code: 0x1fbbf, char: '🮿'},
		{code: 0x1fbc0, char: '🯀'}, {code: 0x1fbc1, char: '🯁'}, {code: 0x1fbc2, char: '🯂'}, {code: 0x1fbc3, char: '🯃'},
		{code: 0x1fbc4, char: '🯄'}, {code: 0x1fbc5, char: '🯅'}, {code: 0x1fbc6, char: '🯆'}, {code: 0x1fbc7, char: '🯇'},
		{code: 0x1fbc8, char: '🯈'}, {code: 0x1fbc9, char: '🯉'}, {code: 0x1fbca, char: '🯊'}, {code: 0x1fbcb, char: '🯋'},
		{code: 0x1fbcc, char: '🯌'}, {code: 0x1fbcd, char: '🯍'}, {code: 0x1fbce, char: '🯎'}, {code: 0x1fbcf, char: '🯏'},
		{code: 0x1fbd0, char: '🯐'}, {code: 0x1fbd1, char: '🯑'}, {code: 0x1fbd2, char: '🯒'}, {code: 0x1fbd3, char: '🯓'},
		{code: 0x1fbd4, char: '🯔'}, {code: 0x1fbd5, char: '🯕'}, {code: 0x1fbd6, char: '🯖'}, {code: 0x1fbd7, char: '🯗'},
		{code: 0x1fbd8, char: '🯘'}, {code: 0x1fbd9, char: '🯙'}, {code: 0x1fbda, char: '🯚'}, {code: 0x1fbdb, char: '🯛'},
		{code: 0x1fbdc, char: '🯜'}, {code: 0x1fbdd, char: '🯝'}, {code: 0x1fbde, char: '🯞'}, {code: 0x1fbdf, char: '🯟'},
		{code: 0x1fbe0, char: '🯠'}, {code: 0x1fbe1, char: '🯡'}, {code: 0x1fbe2, char: '🯢'}, {code: 0x1fbe3, char: '🯣'},
		{code: 0x1fbe4, char: '🯤'}, {code: 0x1fbe5, char: '🯥'}, {code: 0x1fbe6, char: '🯦'}, {code: 0x1fbe7, char: '🯧'},
		{code: 0x1fbe8, char: '🯨'}, {code: 0x1fbe9, char: '🯩'}, {code: 0x1fbea, char: '🯪'}, {code: 0x1fbeb, char: '🯫'},
		{code: 0x1fbec, char: '🯬'}, {code: 0x1fbed, char: '🯭'}, {code: 0x1fbee, char: '🯮'}, {code: 0x1fbef, char: '🯯'},
		{code: 0x1fbf0, char: '🯰'}, {code: 0x1fbf1, char: '🯱'}, {code: 0x1fbf2, char: '🯲'}, {code: 0x1fbf3, char: '🯳'},
		{code: 0x1fbf4, char: '🯴'}, {code: 0x1fbf5, char: '🯵'}, {code: 0x1fbf6, char: '🯶'}, {code: 0x1fbf7, char: '🯷'},
		{code: 0x1fbf8, char: '🯸'}, {code: 0x1fbf9, char: '🯹'}, {code: 0x1fbfa, char: '🯺'}, {code: 0x1fbfb, char: '🯻'},
		{code: 0x1fbfc, char: '🯼'}, {code: 0x1fbfd, char: '🯽'}, {code: 0x1fbfe, char: '🯾'}, {code: 0x1fbff, char: '🯿'},
	}
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

	for i, line := range lines {
		if pad {
			visualLen := convert.UnicodeStringLength(line)
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
