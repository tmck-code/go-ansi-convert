package test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/tmck-code/go-ansi-convert/src/convert"
)

var (
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
			missing := make([]unicodeChar, 0)
			var found []rune

			for _, r := range CompleteSet {
				if _, exists := tc.mapping[r.char]; exists {
					found = append(found, r.char)
				} else {
					if !slices.Contains(tc.symmetricalExceptions, r.char) {
						if !slices.Contains(tc.nonMirroringExceptions, r.char) {
							missing = append(missing, r)
						}
					}
				}
			}

			for batch := range slices.Chunk(missing, 1) {
				for _, ch := range batch {
					fmt.Printf("U+%X '%c', ", ch.code, ch.char)
				}
				fmt.Println()
			}

			completeSetStrings := []string{}
			for _, ch := range CompleteSet {
				completeSetStrings = append(completeSetStrings, fmt.Sprintf("U+%X '%c'", ch.code, ch.char))
			}
			completeMissingStrings := []string{}
			for _, ch := range missing {
				completeMissingStrings = append(completeMissingStrings, fmt.Sprintf("U+%X '%c'", ch.code, ch.char))
			}
			completeFoundStrings := []string{}
			for _, ch := range found {
				completeFoundStrings = append(completeFoundStrings, fmt.Sprintf("U+%X '%c'", ch, ch))
			}

			PrintSimpleTestResults(
				strings.Join(completeSetStrings, ", "),
				"",
				strings.Join(completeMissingStrings, ", "),
			)
			Assert(missing, make([]unicodeChar, 0), t)
		})
	}
}
