package test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
	"github.com/tmck-code/go-ansi-convert/test"
)

var (
	// from https://en.wikipedia.org/wiki/Box-drawing_characters#Symbols_for_Legacy_Computing
	// unicode chars from U+1FB00 to U+1FBFF
	CompleteSet = []test.UnicodeChar{
		{Code: 0x1fb00, Char: '🬀'}, {Code: 0x1fb01, Char: '🬁'}, {Code: 0x1fb02, Char: '🬂'}, {Code: 0x1fb03, Char: '🬃'},
		{Code: 0x1fb04, Char: '🬄'}, {Code: 0x1fb05, Char: '🬅'}, {Code: 0x1fb06, Char: '🬆'}, {Code: 0x1fb07, Char: '🬇'},
		{Code: 0x1fb08, Char: '🬈'}, {Code: 0x1fb09, Char: '🬉'}, {Code: 0x1fb0a, Char: '🬊'}, {Code: 0x1fb0b, Char: '🬋'},
		{Code: 0x1fb0c, Char: '🬌'}, {Code: 0x1fb0d, Char: '🬍'}, {Code: 0x1fb0e, Char: '🬎'}, {Code: 0x1fb0f, Char: '🬏'},
		{Code: 0x1fb10, Char: '🬐'}, {Code: 0x1fb11, Char: '🬑'}, {Code: 0x1fb12, Char: '🬒'}, {Code: 0x1fb13, Char: '🬓'},
		{Code: 0x1fb14, Char: '🬔'}, {Code: 0x1fb15, Char: '🬕'}, {Code: 0x1fb16, Char: '🬖'}, {Code: 0x1fb17, Char: '🬗'},
		{Code: 0x1fb18, Char: '🬘'}, {Code: 0x1fb19, Char: '🬙'}, {Code: 0x1fb1a, Char: '🬚'}, {Code: 0x1fb1b, Char: '🬛'},
		{Code: 0x1fb1c, Char: '🬜'}, {Code: 0x1fb1d, Char: '🬝'}, {Code: 0x1fb1e, Char: '🬞'}, {Code: 0x1fb1f, Char: '🬟'},
		{Code: 0x1fb20, Char: '🬠'}, {Code: 0x1fb21, Char: '🬡'}, {Code: 0x1fb22, Char: '🬢'}, {Code: 0x1fb23, Char: '🬣'},
		{Code: 0x1fb24, Char: '🬤'}, {Code: 0x1fb25, Char: '🬥'}, {Code: 0x1fb26, Char: '🬦'}, {Code: 0x1fb27, Char: '🬧'},
		{Code: 0x1fb28, Char: '🬨'}, {Code: 0x1fb29, Char: '🬩'}, {Code: 0x1fb2a, Char: '🬪'}, {Code: 0x1fb2b, Char: '🬫'},
		{Code: 0x1fb2c, Char: '🬬'}, {Code: 0x1fb2d, Char: '🬭'}, {Code: 0x1fb2e, Char: '🬮'}, {Code: 0x1fb2f, Char: '🬯'},
		{Code: 0x1fb30, Char: '🬰'}, {Code: 0x1fb31, Char: '🬱'}, {Code: 0x1fb32, Char: '🬲'}, {Code: 0x1fb33, Char: '🬳'},
		{Code: 0x1fb34, Char: '🬴'}, {Code: 0x1fb35, Char: '🬵'}, {Code: 0x1fb36, Char: '🬶'}, {Code: 0x1fb37, Char: '🬷'},
		{Code: 0x1fb38, Char: '🬸'}, {Code: 0x1fb39, Char: '🬹'}, {Code: 0x1fb3a, Char: '🬺'}, {Code: 0x1fb3b, Char: '🬻'},
		{Code: 0x1fb3c, Char: '🬼'}, {Code: 0x1fb3d, Char: '🬽'}, {Code: 0x1fb3e, Char: '🬾'}, {Code: 0x1fb3f, Char: '🬿'},
		{Code: 0x1fb40, Char: '🭀'}, {Code: 0x1fb41, Char: '🭁'}, {Code: 0x1fb42, Char: '🭂'}, {Code: 0x1fb43, Char: '🭃'},
		{Code: 0x1fb44, Char: '🭄'}, {Code: 0x1fb45, Char: '🭅'}, {Code: 0x1fb46, Char: '🭆'}, {Code: 0x1fb47, Char: '🭇'},
		{Code: 0x1fb48, Char: '🭈'}, {Code: 0x1fb49, Char: '🭉'}, {Code: 0x1fb4a, Char: '🭊'}, {Code: 0x1fb4b, Char: '🭋'},
		{Code: 0x1fb4c, Char: '🭌'}, {Code: 0x1fb4d, Char: '🭍'}, {Code: 0x1fb4e, Char: '🭎'}, {Code: 0x1fb4f, Char: '🭏'},
		{Code: 0x1fb50, Char: '🭐'}, {Code: 0x1fb51, Char: '🭑'}, {Code: 0x1fb52, Char: '🭒'}, {Code: 0x1fb53, Char: '🭓'},
		{Code: 0x1fb54, Char: '🭔'}, {Code: 0x1fb55, Char: '🭕'}, {Code: 0x1fb56, Char: '🭖'}, {Code: 0x1fb57, Char: '🭗'},
		{Code: 0x1fb58, Char: '🭘'}, {Code: 0x1fb59, Char: '🭙'}, {Code: 0x1fb5a, Char: '🭚'}, {Code: 0x1fb5b, Char: '🭛'},
		{Code: 0x1fb5c, Char: '🭜'}, {Code: 0x1fb5d, Char: '🭝'}, {Code: 0x1fb5e, Char: '🭞'}, {Code: 0x1fb5f, Char: '🭟'},
		{Code: 0x1fb60, Char: '🭠'}, {Code: 0x1fb61, Char: '🭡'}, {Code: 0x1fb62, Char: '🭢'}, {Code: 0x1fb63, Char: '🭣'},
		{Code: 0x1fb64, Char: '🭤'}, {Code: 0x1fb65, Char: '🭥'}, {Code: 0x1fb66, Char: '🭦'}, {Code: 0x1fb67, Char: '🭧'},
		{Code: 0x1fb68, Char: '🭨'}, {Code: 0x1fb69, Char: '🭩'}, {Code: 0x1fb6a, Char: '🭪'}, {Code: 0x1fb6b, Char: '🭫'},
		{Code: 0x1fb6c, Char: '🭬'}, {Code: 0x1fb6d, Char: '🭭'}, {Code: 0x1fb6e, Char: '🭮'}, {Code: 0x1fb6f, Char: '🭯'},
		{Code: 0x1fb70, Char: '🭰'}, {Code: 0x1fb71, Char: '🭱'}, {Code: 0x1fb72, Char: '🭲'}, {Code: 0x1fb73, Char: '🭳'},
		{Code: 0x1fb74, Char: '🭴'}, {Code: 0x1fb75, Char: '🭵'}, {Code: 0x1fb76, Char: '🭶'}, {Code: 0x1fb77, Char: '🭷'},
		{Code: 0x1fb78, Char: '🭸'}, {Code: 0x1fb79, Char: '🭹'}, {Code: 0x1fb7a, Char: '🭺'}, {Code: 0x1fb7b, Char: '🭻'},
		{Code: 0x1fb7c, Char: '🭼'}, {Code: 0x1fb7d, Char: '🭽'}, {Code: 0x1fb7e, Char: '🭾'}, {Code: 0x1fb7f, Char: '🭿'},
		{Code: 0x1fb80, Char: '🮀'}, {Code: 0x1fb81, Char: '🮁'}, {Code: 0x1fb82, Char: '🮂'}, {Code: 0x1fb83, Char: '🮃'},
		{Code: 0x1fb84, Char: '🮄'}, {Code: 0x1fb85, Char: '🮅'}, {Code: 0x1fb86, Char: '🮆'}, {Code: 0x1fb87, Char: '🮇'},
		{Code: 0x1fb88, Char: '🮈'}, {Code: 0x1fb89, Char: '🮉'}, {Code: 0x1fb8a, Char: '🮊'}, {Code: 0x1fb8b, Char: '🮋'},
		{Code: 0x1fb8c, Char: '🮌'}, {Code: 0x1fb8d, Char: '🮍'}, {Code: 0x1fb8e, Char: '🮎'}, {Code: 0x1fb8f, Char: '🮏'},
		{Code: 0x1fb90, Char: '🮐'}, {Code: 0x1fb91, Char: '🮑'}, {Code: 0x1fb92, Char: '🮒'}, {Code: 0x1fb93, Char: '🮓'},
		{Code: 0x1fb94, Char: '🮔'}, {Code: 0x1fb95, Char: '🮕'}, {Code: 0x1fb96, Char: '🮖'}, {Code: 0x1fb97, Char: '🮗'},
		{Code: 0x1fb98, Char: '🮘'}, {Code: 0x1fb99, Char: '🮙'}, {Code: 0x1fb9a, Char: '🮚'}, {Code: 0x1fb9b, Char: '🮛'},
		{Code: 0x1fb9c, Char: '🮜'}, {Code: 0x1fb9d, Char: '🮝'}, {Code: 0x1fb9e, Char: '🮞'}, {Code: 0x1fb9f, Char: '🮟'},
		{Code: 0x1fba0, Char: '🮠'}, {Code: 0x1fba1, Char: '🮡'}, {Code: 0x1fba2, Char: '🮢'}, {Code: 0x1fba3, Char: '🮣'},
		{Code: 0x1fba4, Char: '🮤'}, {Code: 0x1fba5, Char: '🮥'}, {Code: 0x1fba6, Char: '🮦'}, {Code: 0x1fba7, Char: '🮧'},
		{Code: 0x1fba8, Char: '🮨'}, {Code: 0x1fba9, Char: '🮩'}, {Code: 0x1fbaa, Char: '🮪'}, {Code: 0x1fbab, Char: '🮫'},
		{Code: 0x1fbac, Char: '🮬'}, {Code: 0x1fbad, Char: '🮭'}, {Code: 0x1fbae, Char: '🮮'}, {Code: 0x1fbaf, Char: '🮯'},
		{Code: 0x1fbb0, Char: '🮰'}, {Code: 0x1fbb1, Char: '🮱'}, {Code: 0x1fbb2, Char: '🮲'}, {Code: 0x1fbb3, Char: '🮳'},
		{Code: 0x1fbb4, Char: '🮴'}, {Code: 0x1fbb5, Char: '🮵'}, {Code: 0x1fbb6, Char: '🮶'}, {Code: 0x1fbb7, Char: '🮷'},
		{Code: 0x1fbb8, Char: '🮸'}, {Code: 0x1fbb9, Char: '🮹'}, {Code: 0x1fbba, Char: '🮺'}, {Code: 0x1fbbb, Char: '🮻'},
		{Code: 0x1fbbc, Char: '🮼'}, {Code: 0x1fbbd, Char: '🮽'}, {Code: 0x1fbbe, Char: '🮾'}, {Code: 0x1fbbf, Char: '🮿'},
		{Code: 0x1fbc0, Char: '🯀'}, {Code: 0x1fbc1, Char: '🯁'}, {Code: 0x1fbc2, Char: '🯂'}, {Code: 0x1fbc3, Char: '🯃'},
		{Code: 0x1fbc4, Char: '🯄'}, {Code: 0x1fbc5, Char: '🯅'}, {Code: 0x1fbc6, Char: '🯆'}, {Code: 0x1fbc7, Char: '🯇'},
		{Code: 0x1fbc8, Char: '🯈'}, {Code: 0x1fbc9, Char: '🯉'}, {Code: 0x1fbca, Char: '🯊'}, {Code: 0x1fbcb, Char: '🯋'},
		{Code: 0x1fbcc, Char: '🯌'}, {Code: 0x1fbcd, Char: '🯍'}, {Code: 0x1fbce, Char: '🯎'}, {Code: 0x1fbcf, Char: '🯏'},
		{Code: 0x1fbd0, Char: '🯐'}, {Code: 0x1fbd1, Char: '🯑'}, {Code: 0x1fbd2, Char: '🯒'}, {Code: 0x1fbd3, Char: '🯓'},
		{Code: 0x1fbd4, Char: '🯔'}, {Code: 0x1fbd5, Char: '🯕'}, {Code: 0x1fbd6, Char: '🯖'}, {Code: 0x1fbd7, Char: '🯗'},
		{Code: 0x1fbd8, Char: '🯘'}, {Code: 0x1fbd9, Char: '🯙'}, {Code: 0x1fbda, Char: '🯚'}, {Code: 0x1fbdb, Char: '🯛'},
		{Code: 0x1fbdc, Char: '🯜'}, {Code: 0x1fbdd, Char: '🯝'}, {Code: 0x1fbde, Char: '🯞'}, {Code: 0x1fbdf, Char: '🯟'},
		{Code: 0x1fbe0, Char: '🯠'}, {Code: 0x1fbe1, Char: '🯡'}, {Code: 0x1fbe2, Char: '🯢'}, {Code: 0x1fbe3, Char: '🯣'},
		{Code: 0x1fbe4, Char: '🯤'}, {Code: 0x1fbe5, Char: '🯥'}, {Code: 0x1fbe6, Char: '🯦'}, {Code: 0x1fbe7, Char: '🯧'},
		{Code: 0x1fbe8, Char: '🯨'}, {Code: 0x1fbe9, Char: '🯩'}, {Code: 0x1fbea, Char: '🯪'}, {Code: 0x1fbeb, Char: '🯫'},
		{Code: 0x1fbec, Char: '🯬'}, {Code: 0x1fbed, Char: '🯭'}, {Code: 0x1fbee, Char: '🯮'}, {Code: 0x1fbef, Char: '🯯'},
		{Code: 0x1fbf0, Char: '🯰'}, {Code: 0x1fbf1, Char: '🯱'}, {Code: 0x1fbf2, Char: '🯲'}, {Code: 0x1fbf3, Char: '🯳'},
		{Code: 0x1fbf4, Char: '🯴'}, {Code: 0x1fbf5, Char: '🯵'}, {Code: 0x1fbf6, Char: '🯶'}, {Code: 0x1fbf7, Char: '🯷'},
		{Code: 0x1fbf8, Char: '🯸'}, {Code: 0x1fbf9, Char: '🯹'}, {Code: 0x1fbfa, Char: '🯺'}, {Code: 0x1fbfb, Char: '🯻'},
		{Code: 0x1fbfc, Char: '🯼'}, {Code: 0x1fbfd, Char: '🯽'}, {Code: 0x1fbfe, Char: '🯾'}, {Code: 0x1fbff, Char: '🯿'},
	}
)

func TestMapCompleteness(t *testing.T) {
	// Check that character exists in the vertical mirror map

	testCases := []struct {
		name                   string
		mapping                map[rune]rune
		symmetricalExceptions  []rune
		nonMirroringExceptions []rune
	}{
		{
			name:                   "vertical mirror map completeness",
			mapping:                convert.VerticalMirrorMap,
			symmetricalExceptions:  convert.VerticalSymmetricalRunes,
			nonMirroringExceptions: convert.VerticalNonMirroringRunes,
		},
		{
			name:                   "horizontal mirror map completeness",
			mapping:                convert.HorizontalMirrorMap,
			symmetricalExceptions:  convert.HorizontalSymmetricalRunes,
			nonMirroringExceptions: convert.HorizontalNonMirroringRunes,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			missing := make([]test.UnicodeChar, 0)
			var found []rune

			for _, r := range CompleteSet {
				if _, exists := tc.mapping[r.Char]; exists {
					found = append(found, r.Char)
				} else {
					if !slices.Contains(tc.symmetricalExceptions, r.Char) {
						if !slices.Contains(tc.nonMirroringExceptions, r.Char) {
							missing = append(missing, test.UnicodeChar{Code: r.Code, Char: r.Char})
						}
					}
				}
			}

			for batch := range slices.Chunk(missing, 1) {
				for _, ch := range batch {
					fmt.Printf("U+%X '%c', ", ch.Code, ch.Char)
				}
				fmt.Println()
			}

			completeSetStrings := []string{}
			for _, ch := range CompleteSet {
				completeSetStrings = append(completeSetStrings, fmt.Sprintf("U+%X '%c'", ch.Code, ch.Char))
			}
			completeMissingStrings := []string{}
			for _, ch := range missing {
				completeMissingStrings = append(completeMissingStrings, fmt.Sprintf("U+%X '%c'", ch.Code, ch.Char))
			}
			completeFoundStrings := []string{}
			for _, ch := range found {
				completeFoundStrings = append(completeFoundStrings, fmt.Sprintf("U+%X '%c'", ch, ch))
			}

			test.PrintSimpleTestResults(
				strings.Join(completeSetStrings, ", "),
				"",
				strings.Join(completeMissingStrings, ", "),
			)
			test.Assert(missing, make([]test.UnicodeChar, 0), t)
		})
	}
}
